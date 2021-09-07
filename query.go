package processes

import (
	"bufio"
	"bytes"
	"os/exec"
	"time"
)

type ProcessQuery struct {
	stdout *bytes.Buffer
	stderr *bytes.Buffer
}

func Query(exe string, args ...string) (*ProcessQuery, error) {
	process := ProcessQuery{
		stdout: bytes.NewBuffer([]byte{}),
		stderr: bytes.NewBuffer([]byte{}),
	}

	command := exec.Command(exe, args...)

	process.stdout = bytes.NewBuffer([]byte{})
	stdoutWrite := bufio.NewWriter(process.stdout)

	process.stderr = bytes.NewBuffer([]byte{})
	stderrWrite := bufio.NewWriter(process.stderr)

	command.Stdout = stdoutWrite
	command.Stderr = stderrWrite

	errChan := make(chan error)
	go func() {
		if err := command.Start(); err != nil {
			errChan <- err
		}
	}()

	for {
		select {
		case <-time.After(time.Millisecond * 50):
			return &process, nil

		case err := <-errChan:
			return nil, err
		}
	}
}
