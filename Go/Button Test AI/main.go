package main

import (
	"fmt"
	"time"

	"github.com/stianeikeland/go-rpio"
)

func main() {
	err := rpio.Open()
	if err != nil {
		fmt.Println("Error opening GPIO:", err)
		return
	}
	defer rpio.Close()

	button := rpio.Pin(2)
	button.Input()

	for {
		if button.Read() == rpio.High {
			fmt.Println("You pushed me")
			break
		}
		time.Sleep(time.Millisecond * 100)
	}
}

