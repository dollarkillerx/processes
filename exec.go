package processes

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"
)

type ExecLinux struct {
	Path    string
	bash    string
	timeout time.Duration
}

func NewExecLinux() *ExecLinux {
	return NewExecLinuxGen("", "", time.Second*60)
}

func NewExecLinuxGen(path string, bash string, timeout time.Duration) *ExecLinux {
	path = strings.TrimSpace(path)
	bash = strings.TrimSpace(bash)
	if path == "" {
		file, _ := exec.LookPath(os.Args[0])
		path, _ := filepath.Abs(file)
		index := strings.LastIndex(path, string(os.PathSeparator))
		path = path[:index]
	}

	if bash == "" {
		bash = "/bin/sh"
	}

	return &ExecLinux{
		Path:    path,
		bash:    bash,
		timeout: timeout,
	}
}

func (e *ExecLinux) Exec(cmd string) (string, error) {
	cmd = strings.TrimSpace(cmd)
	if cmd == "" {
		return "", errors.New("cmd is null")
	}
	ok := e.execCd(cmd)
	if ok {
		return "", nil
	}

	ctxt, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()

	command := exec.CommandContext(ctxt, "/bin/sh", "-c", cmd)
	command.Dir = e.Path
	output, err := command.CombinedOutput()
	if err != nil {
		if len(output) != 0 {
			err = fmt.Errorf("%s err: %s", output, err.Error())
		}
		return "", err
	}

	return strings.TrimSpace(string(output)), err
}

func (e *ExecLinux) execCd(cmd string) bool {
	split := strings.Split(cmd, " ")
	if len(split) != 2 {
		return false
	}
	if split[0] != "cd" {
		e.Path = path.Join(e.Path)
		return false
	}

	ctxt, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()

	command := exec.CommandContext(ctxt, e.bash, "-c", fmt.Sprintf("%s && pwd", cmd))
	command.Dir = e.Path
	path, err := command.CombinedOutput()
	if err != nil {
		fmt.Println(err)
		return false
	}

	e.Path = strings.TrimSpace(string(path))
	return true
}
