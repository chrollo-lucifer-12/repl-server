package docker

import (
	"bytes"
	"context"
	"os"
	"testing"
)

func TestStartContainer(t *testing.T) {
	d := NewDockerClient()
	defer d.Stop()

	ctx := context.Background()

	containerID := d.StartContainer(ctx, os.Stdout)
	if containerID == "" {
		t.Fatal("Expected container ID, got empty string")
	}

	t.Log("Container started with ID:", containerID)
}

func TestExecCommand(t *testing.T) {
	d := NewDockerClient()
	defer d.Stop()

	ctx := context.Background()
	containerID := d.StartContainer(ctx, os.Stdout)
	if containerID == "" {
		t.Fatal("Failed to start container")
	}

	var buf bytes.Buffer

	// Create a file
	err := d.ExecCommand(ctx, containerID, []string{"touch", "hi.txt"}, &buf)
	if err != nil {
		t.Fatal("ExecCommand failed:", err)
	}
	t.Log("Created file hi.txt, output:", buf.String())

	// List files
	buf.Reset()
	err = d.ExecCommand(ctx, containerID, []string{"ls"}, &buf)
	if err != nil {
		t.Fatal("ExecCommand failed:", err)
	}
	t.Log("Files in container:", buf.String())
}

func TestFileManagement(t *testing.T) {
	d := NewDockerClient()
	defer d.Stop()

	ctx := context.Background()
	containerID := d.StartContainer(ctx, os.Stdout)
	if containerID == "" {
		t.Fatal("Failed to start container")
	}

	var buf bytes.Buffer

	// ---------------------------
	// Test CreateDir
	// ---------------------------
	dirPath := "/code/testdir"
	err := d.CreateDir(ctx, containerID, dirPath, &buf)
	if err != nil {
		t.Fatalf("CreateDir failed: %v", err)
	}
	t.Log("CreateDir output:", buf.String())

	// Verify directory exists
	buf.Reset()
	err = d.ExecCommand(ctx, containerID, []string{"ls", "/code"}, &buf)
	if err != nil {
		t.Fatalf("Listing /code failed: %v", err)
	}
	if !bytes.Contains(buf.Bytes(), []byte("testdir")) {
		t.Fatalf("Directory %s not found in /code: %s", dirPath, buf.String())
	}
	t.Log("/code contents:", buf.String())

	// ---------------------------
	// Test writing a file
	// ---------------------------
	filePath := dirPath + "/hello.txt"
	fileContent := "Hello, Docker!"

	// Write file using echo safely
	writeCmd := []string{"sh", "-c", "echo \"" + fileContent + "\" > " + filePath}
	buf.Reset()
	if err := d.ExecCommand(ctx, containerID, writeCmd, &buf); err != nil {
		t.Fatalf("Writing file failed: %v", err)
	}
	t.Log("WriteFile output:", buf.String())

	// ---------------------------
	// Test ReadFile
	// ---------------------------
	buf.Reset()
	err = d.ReadFile(ctx, containerID, filePath, &buf)
	if err != nil {
		t.Fatalf("ReadFile failed: %v", err)
	}
	if buf.String() != fileContent+"\n" {
		t.Fatalf("Unexpected file content. Got: %q, Expected: %q", buf.String(), fileContent+"\n")
	}
	t.Log("ReadFile output:", buf.String())

	// ---------------------------
	// Test RemoveFile
	// ---------------------------
	buf.Reset()
	err = d.RemoveFile(ctx, filePath, containerID, &buf)
	if err != nil {
		t.Fatalf("RemoveFile failed: %v", err)
	}
	t.Log("RemoveFile output:", buf.String())

	// Verify file removed
	buf.Reset()
	err = d.ExecCommand(ctx, containerID, []string{"ls", dirPath}, &buf)
	if err != nil {
		t.Fatalf("Listing directory after removal failed: %v", err)
	}
	if bytes.Contains(buf.Bytes(), []byte("hello.txt")) {
		t.Fatal("File was not removed")
	}
	t.Log(dirPath, "contents after removal:", buf.String())
}
