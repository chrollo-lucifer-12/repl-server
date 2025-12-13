package docker

import (
	"context"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/client"
)

type FileInfo struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Size int64  `json:"size"`
	Mode string `json:"mode"`
}

type ContainerInfo struct {
	createdAt time.Time
	userId    string
}

type DockerClient struct {
	dockerClient *client.Client
	containers   sync.Map
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

func (d *DockerClient) StartContainer(ctx context.Context, outputWriter io.Writer, userId string) string {
	imageName := "node:20-bullseye"
	out, err := d.dockerClient.ImagePull(ctx, imageName, client.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	defer out.Close()
	io.Copy(os.Stdout, out)

	hostDir := "/var/repl/users/" + userId
	os.MkdirAll(hostDir, 0755)
	containerDir := "/home/" + userId

	resp, err := d.dockerClient.ContainerCreate(ctx, client.ContainerCreateOptions{
		Image:  imageName,
		Config: &container.Config{Tty: true, OpenStdin: true, AttachStdin: true, AttachStdout: true, AttachStderr: true, Cmd: []string{"sh"}, WorkingDir: containerDir},
		HostConfig: &container.HostConfig{
			Binds: []string{
				hostDir + ":" + containerDir,
			},
			Resources: container.Resources{
				Memory:     512 * 1024 * 1024,
				MemorySwap: 512 * 1024 * 1024,
				CPUShares:  512,
				CPUQuota:   50000,
				CPUPeriod:  100000,
			},
		},
	})
	if err != nil {
		panic(err)
	}

	if _, err := d.dockerClient.ContainerStart(ctx, resp.ID, client.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	fmt.Println(resp.ID)

	d.containers.Store(userId, resp.ID)

	return resp.ID
}
