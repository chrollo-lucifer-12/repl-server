package docker

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/chrollo-lucifer-12/repl/utils"
)

func (d *DockerClient) WriteFile(
	ctx context.Context,
	containerID, path, content string,
	outputWriter io.Writer,
) error {

	cmd := []string{
		"sh",
		"-c",
		fmt.Sprintf("cat > %s << 'EOF'\n%s\nEOF", path, content),
	}
	return d.ExecCommand(ctx, containerID, cmd, outputWriter)
}

func (d *DockerClient) ReadFile(ctx context.Context, containerID, path string, outputWriter io.Writer) error {
	cmd := []string{"cat", path}
	if err := d.ExecCommand(ctx, containerID, cmd, outputWriter); err != nil {
		return err
	}
	return nil
}

func (d *DockerClient) CreateDir(ctx context.Context, containerID, path string, outputWriter io.Writer) error {
	cmd := []string{"sh", "-c", "mkdir -p " + path}
	return d.ExecCommand(ctx, containerID, cmd, outputWriter)
}

func (d *DockerClient) RemoveFile(ctx context.Context, path string, containerId string, outputWriter io.Writer) error {
	cmd := []string{"rm", "-f", path}
	return d.ExecCommand(ctx, containerId, cmd, outputWriter)
}

func (d *DockerClient) ListFiles(ctx context.Context, containerID, path string, outputWriter io.Writer) error {
	cmd := []string{"ls", "-lA", "--color=never", path}

	var buf bytes.Buffer
	if err := d.ExecCommand(ctx, containerID, cmd, &buf); err != nil {
		return err
	}

	lines := bytes.Split(buf.Bytes(), []byte("\n"))
	var files []FileInfo

	for _, line := range lines {
		if len(line) == 0 {
			continue
		}

		fields := bytes.Fields(line)
		if len(fields) < 9 {
			continue // skip invalid lines
		}

		mode := string(fields[0])
		size := utils.BytesToInt64(fields[4])
		name := string(fields[8])
		fileType := "file"
		if mode[0] == 'd' {
			fileType = "dir"
		}

		files = append(files, FileInfo{
			Name: name,
			Type: fileType,
			Size: size,
			Mode: mode,
		})
	}

	jsonBytes, _ := json.MarshalIndent(files, "", "  ")
	if outputWriter != nil {
		outputWriter.Write(jsonBytes)
	}
	return nil
}

func (d *DockerClient) StatFile(ctx context.Context, containerID, path string, outputWriter io.Writer) error {
	cmd := []string{"stat", "-c", "%F %s %a", path}
	var buf bytes.Buffer
	if err := d.ExecCommand(ctx, containerID, cmd, &buf); err != nil {
		return err
	}
	output := strings.TrimSpace(buf.String())
	parts := strings.Fields(output)
	if len(parts) < 3 {
		return fmt.Errorf("unexpected stat output: %q", output)
	}

	fileType := parts[0]
	size, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid file size: %v", err)
	}
	mode := parts[2]

	info := &FileInfo{
		Name: path,
		Type: fileType,
		Size: size,
		Mode: mode,
	}

	if outputWriter != nil {
		jsonBytes, _ := json.Marshal(info)
		outputWriter.Write(jsonBytes)
	}
	return nil
}

func (d *DockerClient) SearchInFile(ctx context.Context, containerID, filePath, search string, outputWriter io.Writer) error {
	cmd := []string{"grep", "-nF", search, filePath}
	return d.ExecCommand(ctx, containerID, cmd, outputWriter)
}

func (d *DockerClient) RenameFileDir(ctx context.Context, containerID, path string, newName string, outputWriter io.Writer) error {
	cmd := []string{"mv", path, newName}
	return d.ExecCommand(ctx, containerID, cmd, outputWriter)
}
