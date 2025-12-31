package docker

import (
	"context"
	"fmt"

	"github.com/moby/moby/client"
)

func (d *DockerClient) ResizeTerminal(ctx context.Context,
	userId string, rows int, cols int) error {
	containerId, ok := d.containers.Load(userId)
	if !ok {
		return fmt.Errorf("container was deleted")
	}

	_, err := d.dockerClient.ExecResize(ctx, containerId.(string), client.ExecResizeOptions{
		Height: uint(rows),
		Width:  uint(cols),
	})
	if err != nil {
		return nil
	}

	return nil
}
