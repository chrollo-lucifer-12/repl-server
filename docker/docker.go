package docker

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/chrollo-lucifer-12/repl/utils"
	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/client"
)

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

	resp, err := d.dockerClient.ContainerCreate(ctx, client.ContainerCreateOptions{
		Image:  imageName,
		Config: &container.Config{Tty: true, OpenStdin: true, AttachStdin: true, AttachStdout: true, AttachStderr: true, Cmd: []string{"sh"}},
	})
	if err != nil {
		panic(err)
	}

	if _, err := d.dockerClient.ContainerStart(ctx, resp.ID, client.ContainerStartOptions{}); err != nil {
		panic(err)
	}
	setUp := []string{}
	d.ExecCommand(ctx, resp.ID, setUp, outputWriter)

	return resp.ID
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
