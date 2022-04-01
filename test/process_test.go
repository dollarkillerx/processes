package main

import (
	"context"
	"fmt"
	"github.com/dollarkillerx/processes"
	"log"
	"os/exec"
	"strings"
	"testing"
	"time"
)

func TestStrSp(t *testing.T) {
	a := `ps -ef |awk '{print $2,$8}'`
	command, err := processes.RunCommand(a)
	if err != nil {
		log.Fatalln(err)
	}
	split := strings.Split(command, "\n")
	log.Println(len(split))
}

func TestSnapProcess(t *testing.T) {
	snapshot, err := processes.ProcessSnapshot()
	if err != nil {
		log.Fatalln(err)
	}
	for k, v := range snapshot {
		log.Printf("%s : %s \n", k, v)
	}
}

func TestP1(t *testing.T) {
	linux, err := processes.GetPid("rcu_gp")
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(linux)
}

func TestKill(t *testing.T) {
	pid, err := processes.GetPid("plumber")
	if err != nil {
		log.Fatalln(err)
	}
	if err := processes.KillByPid(pid); err != nil {
		log.Fatalln(err)
	}
}

func TestExec(t *testing.T) {
	cmd := processes.NewExecLinux()
	exec, err := cmd.Exec("pwd")
	if err != nil {
		return
	}

	fmt.Println(exec)

	cmd.Exec("cd ../../")

	exec, err = cmd.Exec("pwd")
	if err != nil {
		return
	}

	fmt.Println(exec)
}

func TestExec2(t *testing.T) {
	cmd := processes.NewExecLinuxGen("", "bash", time.Second)
	exec, err := cmd.Exec("sleep 3")
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println(exec)
}

func TestExec3(t *testing.T) {
	ctxt, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	command := exec.CommandContext(ctxt, "bash", "-c", "sleep 3")
	path, err := command.CombinedOutput()
	if err != nil {
		panic(err)
	}
	fmt.Println(path)
}
