package autoitx

import (
	"os/exec"
	"testing"
	"time"
)

func TestControlGetText(t *testing.T) {
	cm0 := exec.Command("notepad.exe", "example.txt")
	cm0.Start()
	time.Sleep(time.Second * 3)
	txt := ControlGetText("[CLASS:Notepad]", "", "[CLASS:Edit; INSTANCE:1]")
	if txt != "example text" {
		t.Fatal(txt)
	}
	cm0.Process.Kill()
}

func TestControlSetText(t *testing.T) {
	cm0 := exec.Command("notepad.exe", "example.txt")
	cm0.Start()
	time.Sleep(time.Second * 3)
	//
	ControlSetText("[CLASS:Notepad]", "", "[CLASS:Edit; INSTANCE:1]", "autoitx rules")
	time.Sleep(time.Millisecond * 500)
	txt := ControlGetText("[CLASS:Notepad]", "", "[CLASS:Edit; INSTANCE:1]")
	if txt != "autoitx rules" {
		t.Fatal(txt)
	}
	ControlSetText("[CLASS:Notepad]", "", "[CLASS:Edit; INSTANCE:1]", "example text")
	//
	cm0.Process.Kill()
}

func TestWinExists(t *testing.T) {
	cm0 := exec.Command("notepad.exe", "example.txt")
	cm0.Start()
	//
	time.Sleep(time.Second * 3)
	ok := WinExists("[CLASS:Notepad]", "")
	if !ok {
		t.Fatal("win does not exist")
	}
	cm0.Process.Kill()
}
