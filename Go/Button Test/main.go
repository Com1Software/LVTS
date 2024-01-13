package main

import (
	"fmt"
	"time"

	"github.com/Com1Software/go-rpio"
)

func main() {
	err := rpio.Open()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rpio.Close()

	button := rpio.Pin(2)
	button.Input()

	for {
		if button.Read() == rpio.Low {
			fmt.Println("You pushed me")
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
}
