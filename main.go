package main

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/gvalkov/golang-evdev"
)

type Message struct {
	Device *evdev.InputDevice
	Events []evdev.InputEvent
}

func listenEvents(kbd *evdev.InputDevice, replyTo chan Message) {
	for {
		events, _ := kbd.Read()
		replyTo <- Message{Device: kbd, Events: events}
	}
}

func main() {
	inbox := make(chan Message, 8)
	devs, _ := evdev.ListInputDevices("/dev/input/event3")
	go listenEvents(devs[0], inbox)
	ctrl := false
	for {
		select {
		case msg := <-inbox:
			for _, ev := range msg.Events {
				switch ev.Type {
				case evdev.EV_KEY:
					switch ev.Value {
					case 1: // key down
						switch ev.Code {
						case 29:
							ctrl = true
						case 36:
							if ctrl {
								cmd := exec.Command("xdotool", "getwindowfocus", "getwindowname")
								wname, _ := cmd.Output()
								name := strings.TrimSpace(string(wname))
								if name != "Yakuake" && !strings.Contains(name, "Konsole") {
									fmt.Println(name)
									cmd = exec.Command("xdotool", "key", "Return")
									cmd.Run()
								}
							}
						}
					case 0:
						switch ev.Code {
						case 29:
							ctrl = false
						}
					}
				}
			}
		}
	}
}
