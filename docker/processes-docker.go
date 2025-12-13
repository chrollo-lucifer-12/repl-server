package docker

import (
	"context"
	"fmt"
	"io"

	"github.com/moby/moby/client"
)

func (d *DockerClient) ExecCommand(ctx context.Context, containerId string, cmd []string, outputWriter io.Writer) error {
	execResp, err := d.dockerClient.ExecCreate(ctx, containerId, client.ExecCreateOptions{
		Cmd:          cmd,
		AttachStdout: true,
		AttachStdin:  true,
		TTY:          true,
	})
	if err != nil {
		return err
	}
	resp, err := d.dockerClient.ExecAttach(ctx, execResp.ID, client.ExecAttachOptions{TTY: true})
	if err != nil {
		panic(err)
		return err
	}
	defer resp.Close()
	if outputWriter == nil {
		_, err = io.Copy(io.Discard, resp.Reader)
	} else {
		_, err = io.Copy(outputWriter, resp.Reader)
	}
	return err
}

func (d *DockerClient) StartInteractiveRepl(ctx context.Context, userId string, input io.Reader, output io.Writer) error {
	containerId, ok := d.containers.Load(userId)
	if !ok {
		return fmt.Errorf("container was deleted")
	}
	execResp, err := d.dockerClient.ExecCreate(ctx, containerId.(string), client.ExecCreateOptions{
		Cmd:          []string{"node"},
		AttachStdout: true,
		AttachStdin:  true,
		TTY:          true,
	})
	if err != nil {
		return err
	}

	hijackedResp, err := d.dockerClient.ExecAttach(ctx, execResp.ID, client.ExecAttachOptions{
		TTY: true,
	})
	defer hijackedResp.Close()

	go func() {
		if input != nil {
			io.Copy(hijackedResp.Conn, input)
		}
	}()

	if output != nil {
		_, err = io.Copy(output, hijackedResp.Conn)
	}

	inspectResp, err := d.dockerClient.ExecInspect(ctx, execResp.ID, client.ExecInspectOptions{})
	if inspectResp.ExitCode != 0 {
		return fmt.Errorf("node exited with code %d", inspectResp.ExitCode)
	}
	return nil
}

func (d *DockerClient) StartLongRunningProcess(ctx context.Context, userId string, cmd []string, outputWriter io.Writer) (string, error) {
	containerId, ok := d.containers.Load(userId)
	if !ok {
		return "", fmt.Errorf("container was deleted")
	}
	execResp, err := d.dockerClient.ExecCreate(ctx, containerId.(string), client.ExecCreateOptions{
		Cmd:          cmd,
		AttachStdout: true,
		AttachStderr: true,
		TTY:          false,
	})

	if err != nil {
		return "", err
	}
	hijackedResp, err := d.dockerClient.ExecAttach(ctx, execResp.ID, client.ExecAttachOptions{})

	go func() {
		defer hijackedResp.Close()
		if outputWriter != nil {
			io.Copy(outputWriter, hijackedResp.Reader)
		}
	}()

	return execResp.ID, nil
}
