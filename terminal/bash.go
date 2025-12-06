package terminal

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"

	"github.com/creack/pty"
)

type BashTerminal struct {
	cmd  *exec.Cmd
	pty  *os.File
	buff []byte
}

func NewBashTerminal() Terminal {
	return &BashTerminal{
		cmd:  exec.Command("bash"),
		buff: make([]byte, 4096),
	}
}

func (t *BashTerminal) Start() error {
	ptmx, err := pty.Start(t.cmd)
	if err != nil {
		return err
	}
	t.pty = ptmx
	return nil
}

func (t *BashTerminal) Run(command string) (string, error) {
	if t.pty == nil {
		return "", fmt.Errorf("terminal not started")
	}

	_, err := io.WriteString(t.pty, command+"\n")
	if err != nil {
		return "", err
	}

	time.Sleep(150 * time.Millisecond)

	n, err := t.pty.Read(t.buff)
	if err != nil {
		return "", err
	}

	return string(t.buff[:n]), nil
}

func (t *BashTerminal) Close() error {
	if t.pty != nil {
		return t.pty.Close()
	}
	return nil
}
