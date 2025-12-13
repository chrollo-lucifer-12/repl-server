package docker

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"
)

/*
------------------------------------
Helpers
------------------------------------
*/

func setupContainer(t *testing.T) (*DockerClient, context.Context, string) {
	t.Helper()

	d := NewDockerClient()
	ctx := context.Background()

	containerID := d.StartContainer(ctx, nil)
	if containerID == "" {
		t.Fatal("failed to start container")
	}

	t.Cleanup(func() {
		_ = d.DeleteContainer(ctx, containerID)
		_ = d.Stop()
	})

	return d, ctx, containerID
}

/*
------------------------------------
Tests
------------------------------------
*/

func TestCreateDir(t *testing.T) {
	d, ctx, containerID := setupContainer(t)

	var buf bytes.Buffer
	err := d.CreateDir(ctx, containerID, "/home/hi/testdir", &buf)
	if err != nil {
		t.Fatalf("CreateDir failed: %v", err)
	}

	buf.Reset()
	err = d.ExecCommand(ctx, containerID, []string{"ls", "/home/hi"}, &buf)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Contains(buf.Bytes(), []byte("testdir")) {
		t.Fatalf("directory not created, got: %s", buf.String())
	}
}

func TestWriteAndReadFile(t *testing.T) {
	d, ctx, containerID := setupContainer(t)

	var buf bytes.Buffer
	path := "/home/hi/hello.txt"
	content := "hello docker"

	err := d.WriteFile(ctx, containerID, path, content, &buf)
	if err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}

	buf.Reset()
	err = d.ReadFile(ctx, containerID, path, &buf)
	if err != nil {
		t.Fatalf("ReadFile failed: %v", err)
	}

	if buf.String() != content {
		t.Fatalf("unexpected content: %q", buf.String())
	}
}

// func TestListFiles(t *testing.T) {
// 	d, ctx, containerID := setupContainer(t)

// 	_ = d.WriteFile(ctx, containerID, "/home/hi/a.txt", "a", nil)
// 	_ = d.WriteFile(ctx, containerID, "/home/hi/b.txt", "b", nil)

// 	var buf bytes.Buffer
// 	err := d.ListFiles(ctx, containerID, "/home/hi", &buf)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	var files []FileInfo
// 	if err := json.Unmarshal(buf.Bytes(), &files); err != nil {
// 		t.Fatal(err)
// 	}

// 	if len(files) < 2 {
// 		t.Fatalf("expected files, got %+v", files)
// 	}
// }

func TestStatFile(t *testing.T) {
	d, ctx, containerID := setupContainer(t)

	_ = d.WriteFile(ctx, containerID, "/home/hi/stat.txt", "data", nil)

	var buf bytes.Buffer
	err := d.StatFile(ctx, containerID, "/home/hi/stat.txt", &buf)
	if err != nil {
		t.Fatal(err)
	}

	var info FileInfo
	if err := json.Unmarshal(buf.Bytes(), &info); err != nil {
		t.Fatal(err)
	}

	if info.Size == 0 {
		t.Fatal("file size should not be zero")
	}
}

func TestSearchInFile(t *testing.T) {
	d, ctx, containerID := setupContainer(t)

	content := "hello\nworld\nhello again"
	_ = d.WriteFile(ctx, containerID, "/home/hi/search.txt", content, nil)

	var buf bytes.Buffer
	err := d.SearchInFile(ctx, containerID, "/home/hi/search.txt", "hello", &buf)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Contains(buf.Bytes(), []byte("hello")) {
		t.Fatal("search term not found")
	}
}

func TestRenameFileDir(t *testing.T) {
	d, ctx, containerID := setupContainer(t)

	_ = d.WriteFile(ctx, containerID, "/home/hi/old.txt", "x", nil)

	var buf bytes.Buffer
	err := d.RenameFileDir(ctx, containerID, "/home/hi/old.txt", "/home/hi/new.txt", &buf)
	if err != nil {
		t.Fatal(err)
	}

	buf.Reset()
	_ = d.ExecCommand(ctx, containerID, []string{"ls", "/home/hi"}, &buf)

	if !bytes.Contains(buf.Bytes(), []byte("new.txt")) {
		t.Fatal("file rename failed")
	}
}

func TestRemoveFile(t *testing.T) {
	d, ctx, containerID := setupContainer(t)

	_ = d.WriteFile(ctx, containerID, "/home/hi/delete.txt", "bye", nil)

	var buf bytes.Buffer
	err := d.RemoveFile(ctx, "/home/hi/delete.txt", containerID, &buf)
	if err != nil {
		t.Fatal(err)
	}

	buf.Reset()
	_ = d.ExecCommand(ctx, containerID, []string{"ls", "/home/hi"}, &buf)

	if bytes.Contains(buf.Bytes(), []byte("delete.txt")) {
		t.Fatal("file was not deleted")
	}
}

func TestNodeJSProjectRun(t *testing.T) {
	d, ctx, containerID := setupContainer(t)

	projectDir := "/home/hi/nodeapp"
	var buf bytes.Buffer

	// 1. Create project directory
	if err := d.CreateDir(ctx, containerID, projectDir, &buf); err != nil {
		t.Fatalf("CreateDir failed: %v", err)
	}

	// 2. npm init -y
	buf.Reset()
	err := d.ExecCommand(ctx, containerID,
		[]string{"sh", "-c", "cd " + projectDir + " && npm init -y"},
		&buf,
	)
	if err != nil {
		t.Fatalf("npm init failed: %v\n%s", err, buf.String())
	}

	// 3. install lodash
	buf.Reset()
	err = d.ExecCommand(ctx, containerID,
		[]string{"sh", "-c", "cd " + projectDir + " && npm install lodash"},
		&buf,
	)
	if err != nil {
		t.Fatalf("npm install failed: %v\n%s", err, buf.String())
	}

	// 4. create index.js
	jsCode := `
console.log("hi");
`
	buf.Reset()
	err = d.WriteFile(ctx, containerID, projectDir+"/index.js", jsCode, &buf)
	if err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}

	// 5. run node index.js
	buf.Reset()
	err = d.ExecCommand(ctx, containerID,
		[]string{"sh", "-c", "cd " + projectDir + " && node index.js"},
		&buf,
	)
	if err != nil {
		t.Fatalf("node run failed: %v\n%s", err, buf.String())
	}

	// 6. verify output
	output := buf.String()
	expected := "RESULT: 2,4,6,8"

	if !bytes.Contains([]byte(output), []byte(expected)) {
		t.Fatalf("unexpected output:\n%s\nexpected to contain:\n%s", output, expected)
	}

	t.Log("Node app output:", output)
}
