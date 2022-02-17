package processes

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

type Process struct {
	Cmd *exec.Cmd
}

func Command(exe string, args ...string) (*Process, error) {
	command := exec.Command(exe, args...)
	errChan := make(chan error)
	go func() {
		if err := command.Start(); err != nil {
			errChan <- err
		}
	}()

	for {
		select {
		case <-time.After(time.Millisecond * 50):
			return &Process{Cmd: command}, nil

		case err := <-errChan:
			return nil, err
		}
	}
}

func ProcessSnapshot() (map[string]string, error) {
	if runtime.GOOS == "windows" {
		return processSnapshotWindows()
	}

	return processSnapshotLinux()
}

func GetPid(serverName string) (string, error) {
	snapshot, err := ProcessSnapshot()
	if err != nil {
		return "", err
	}
	for k, v := range snapshot {
		if strings.Index(v, serverName) != -1 {
			return k, nil
		}
	}
	return "", fmt.Errorf("NOT FUND")
}

func KillByPid(pid string) error {
	if runtime.GOOS == "windows" {
		return killPidWindows(pid)
	}
	return killPidLinux(pid)
}

func killPidLinux(pid string) error {
	_, err := RunCommand(fmt.Sprintf("kill -9 %s", pid))
	return err
}

func killPidWindows(pid string) error {
	c := fmt.Sprintf(`taskkill /F /pid %s`, pid)
	_, err := RunCommand(c)
	return err
}

func processSnapshotLinux() (map[string]string, error) {
	rp := map[string]string{}
	a := `ps -ef |awk '{print $2,$8}'`
	command, err := RunCommand(a)
	if err != nil {
		return nil, err
	}
	split := strings.Split(command, "\n")
	for _, v := range split {
		sp2 := strings.Split(v, " ")
		if len(sp2) >= 2 {
			rp[sp2[0]] = sp2[1]
		}
	}
	return rp, nil
}

func processSnapshotWindows() (map[string]string, error) {
	rp := map[string]string{}
	a := "tasklist"
	command, err := RunCommand(a)
	if err != nil {
		return nil, err
	}
	split := strings.Split(command, "\n")
	for _, v := range split {
		sp2 := strings.Split(v, " ")
		sp2 = clearNone(sp2)
		if len(sp2) >= 2 {
			rp[sp2[1]] = sp2[0]
		}
	}
	return rp, nil
}

func runInWindows(cmd string) (string, error) {
	result, err := exec.Command("cmd", "/c", cmd).Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(result)), err
}

func RunCommand(cmd string) (string, error) {
	if runtime.GOOS == "windows" {
		return runInWindows(cmd)
	} else {
		return runInLinux(cmd)
	}
}

func runInLinux(cmd string) (string, error) {
	//fmt.Println("Running Linux cmd:" + cmd)
	result, err := exec.Command("/bin/sh", "-c", cmd).Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(result)), err
}

func clearNone(sp []string) []string {
	var result []string
	for _, v := range sp {
		if v != "" {
			result = append(result, v)
		}
	}
	return result
}
