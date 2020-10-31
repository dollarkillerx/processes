package main

import (
	"github.com/dollarkillerx/processes"
	"log"
)

func main() {
	snapshot, err := processes.ProcessSnapshot()
	if err != nil {
		log.Fatalln(err)
	}
	for k,v := range snapshot {
		log.Printf("%s %s \n",k,v)
	}

	pid, err := processes.GetPid("plumber.exe")
	if err != nil {
		log.Fatalln(err)
	}
	if err := processes.KillByPid(pid); err != nil {
		log.Fatalln(err)
	}
}
