package docker

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/chrollo-lucifer-12/repl/utils"
	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/client"
)

type FileInfo struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Size int64  `json:"size"`
	Mode string `json:"mode"`
}

type DockerClient struct {
	dockerClient *client.Client
}

func NewDockerClient() *DockerClient {
	apiClient, err := client.New(client.FromEnv)
	if err != nil {
		panic(err)
	}
	return &DockerClient{dockerClient: apiClient}
}

func (d *DockerClient) Stop() error {
	return d.dockerClient.Close()
}

func (d *DockerClient) StartContainer(ctx context.Context, outputWriter io.Writer) string {
	imageName := "node:lts-alpine3.23"
	out, err := d.dockerClient.ImagePull(ctx, imageName, client.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	defer out.Close()
	io.Copy(os.Stdout, out)

	hostDir := "/var/repl/users/" + "45"
	os.MkdirAll(hostDir, 0755)
	containerDir := "/home/" + "hi"

	resp, err := d.dockerClient.ContainerCreate(ctx, client.ContainerCreateOptions{
		Image:  imageName,
		Config: &container.Config{Tty: true, OpenStdin: true, AttachStdin: true, AttachStdout: true, AttachStderr: true, Cmd: []string{"sh"}, WorkingDir: containerDir},
		HostConfig: &container.HostConfig{
			Binds: []string{
				hostDir + ":" + containerDir,
			},
			NetworkMode: "none",
		},
	})
	if err != nil {
		panic(err)
	}

	if _, err := d.dockerClient.ContainerStart(ctx, resp.ID, client.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	return resp.ID
}

func (d *DockerClient) ExecCommand(ctx context.Context, containerId string, cmd []string, outputWriter io.Writer) error {
	execResp, err := d.dockerClient.ExecCreate(ctx, containerId, client.ExecCreateOptions{
		Cmd:          cmd,
		AttachStdout: true,
		AttachStdin:  true,
	})
	if err != nil {
		return err
	}
	resp, err := d.dockerClient.ExecAttach(ctx, execResp.ID, client.ExecAttachOptions{})
	if err != nil {
		return err
	}
	defer resp.Close()

	return utils.ReadDockerOutput(resp.Reader, outputWriter)
}

func (d *DockerClient) WriteFile(ctx context.Context, containerID, path string, content string, outputWriter io.Writer) error {
	cmd := []string{"sh", "-c", "printf \"" + content + "\" > " + path}
	if err := d.ExecCommand(ctx, containerID, cmd, outputWriter); err != nil {
		return err
	}
	return nil
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
	cmd := []string{"ls", "-lA", path}
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
			continue
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

	jsonBytes, _ := json.Marshal(files)
	if outputWriter != nil {
		outputWriter.Write(jsonBytes)
	}
	return nil
}

func (d *DockerClient) RemoveContainer(ctx context.Context, containerId string) error {
	if _, err := d.dockerClient.ContainerStop(ctx, containerId, client.ContainerStopOptions{}); err != nil {
		return err
	}
	return nil
}

func (d *DockerClient) RemoveAllContainers(ctx context.Context) {
	containers, err := d.dockerClient.ContainerList(ctx, client.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	for _, container := range containers.Items {
		fmt.Print("Stopping container ", container.ID[:10], "... ")
		noWaitTimeout := 0
		if _, err := d.dockerClient.ContainerStop(ctx, container.ID, client.ContainerStopOptions{Timeout: &noWaitTimeout}); err != nil {
			panic(err)
		}
		fmt.Println("Success")
	}
}
