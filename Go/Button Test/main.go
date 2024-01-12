package main

import (
	"fmt"
	"time"

	"github.com/stianeikeland/go-rpio"
)

func main() {
	err := rpio.Open()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rpio.Close()

	ledpin := rpio.Pin(4)
	pushpin := rpio.Pin(17)

	ledpin.Output()
	pushpin.Input()

	ledpin.Low()
	for {
		ledpin.Toggle()
		time.Sleep(200 * time.Millisecond)
	}
}

