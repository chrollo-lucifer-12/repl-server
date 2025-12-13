package docker

import (
	"context"
	"testing"
)

func TestDocker(t *testing.T) {
	d := NewDockerClient()
	ctx := context.Background()
	id := d.StartContainer(ctx)
	t.Log(id)
}
