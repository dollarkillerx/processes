package main

import (
	"github.com/dollarkillerx/processes"
	"log"
	"strings"
	"testing"
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