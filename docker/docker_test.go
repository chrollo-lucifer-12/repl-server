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
	err := d.ExecCommand(ctx, containerID, []string{"echo", "hello world"}, &buf)
	if err != nil {
		t.Fatal("ExecCommand failed:", err)
	}

	// log the output
	t.Log("Command output:", buf.String())
}
