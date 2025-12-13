package docker

import (
	"context"
	"fmt"

	"github.com/moby/moby/client"
)

func (d *DockerClient) RemoveContainer(ctx context.Context, containerId string) error {
	_, ok := d.containers.Load(containerId)
	if !ok {
		return fmt.Errorf("container was deleted")
	}
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

func (d *DockerClient) DeleteContainer(ctx context.Context, containerId string) error {
	_, ok := d.containers.Load(containerId)
	if !ok {
		return fmt.Errorf("container was deleted")
	}
	timeout := 0
	if _, err := d.dockerClient.ContainerStop(ctx, containerId, client.ContainerStopOptions{
		Timeout: &timeout,
	}); err != nil {
		return err
	}

	_, err := d.dockerClient.ContainerRemove(ctx, containerId, client.ContainerRemoveOptions{
		Force: true,
	})

	return err
}
