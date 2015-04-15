package notify

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func TestMockNotifier(t *testing.T) {
	m := NewManager()
	n := MockNotifier{}
	m.Register(n)
	if len(m.notifiers) != 1 {
		t.Errorf("Got (%v) Expected (1)", len(m.notifiers))
	}

	// Hijack stdout
	old := os.Stdout // keep backup of the real stdout
	r, w, err := os.Pipe()
	os.Stdout = w
	outC := make(chan string)

	// copy the output in a separate goroutine so printing can't block indefinitely
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	m.send("Test message")

	// back to normal state
	w.Close()
	os.Stdout = old // restoring the real stdout
	out := <-outC

	if err != nil || strings.TrimSpace(out) != "Test message" {
		t.Errorf("Did not set message. Got(%v)", out)
	}
}

func TestEmailNotifier(t *testing.T) {
	t.Skip()
	m := NewManager()
	n := EmailNotifier{}
	m.Register(n)
	if len(m.notifiers) != 1 {
		t.Errorf("Got (%v) Expected (1)", len(m.notifiers))
	}
	m.send("ERROR: test error.")
}
