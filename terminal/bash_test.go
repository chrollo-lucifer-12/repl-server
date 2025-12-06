package terminal

import "testing"

func TestBash(t *testing.T) {
	term := NewBashTerminal()

	if err := term.Start(); err != nil {
		t.Fatalf("failed to start bash: %v", err)
	}

	msg, err := term.Run("ls")
	if err != nil {
		t.Fatalf("failed to run command: %v", err)
	}

	t.Logf("output:\n%s", msg)
}
