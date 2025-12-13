package docker

import (
	"context"
	"io"
	"os"

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
	imageName := "node:20-bullseye"
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
