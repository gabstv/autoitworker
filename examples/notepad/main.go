package main

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/gabstv/autoitworker/au3master"
	"github.com/gabstv/freeport"
)

func main() {
	ptcp, err := freeport.TCP()
	if err != nil {
		panic(err)
	}
	conn := fmt.Sprintf("127.0.0.1:%d", ptcp)
	sv := au3master.NewProduction(conn)
	go sv.RunHTTP()
	cmd := exec.Command("au3worker.exe", fmt.Sprintf("http://%s", conn))
	err = cmd.Start()
	if err != nil {
		panic(err)
	}
	time.Sleep(time.Second * 1)
	title := sv.WinGetTitle("[CLASS:Notepad]", "")
	fmt.Println("Titulo: ", title)
	sv.Shutdown()
}
