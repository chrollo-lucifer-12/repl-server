package docker

import (
	"context"
	"io"
	"os"

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

func (d *DockerClient) StartContainer(ctx context.Context) string {
	imageName := "bfirsh/reticulate-splines"
	out, err := d.dockerClient.ImagePull(ctx, imageName, client.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	defer out.Close()
	io.Copy(os.Stdout, out)

	resp, err := d.dockerClient.ContainerCreate(ctx, client.ContainerCreateOptions{
		Image: imageName,
	})
	if err != nil {
		panic(err)
	}

	if _, err := d.dockerClient.ContainerStart(ctx, resp.ID, client.ContainerStartOptions{}); err != nil {
		panic(err)
	}
	return resp.ID
}
