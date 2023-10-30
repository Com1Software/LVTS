

package main

import (
	"fmt"
	"log"
	"strings"

	lcd "github.com/mskrha/rpi-lcd"
	"go.bug.st/serial"
)

func main() {
	fmt.Println("Pi Test Vehicle")
	l := lcd.New(lcd.LCD{Bus: "/dev/i2c-1", Address: 0x27, Rows: 4, Cols: 20, Backlight: true})
	if err := l.Init(); err != nil {
		panic(err)
	}
	gpsport := ""
	rcport := ""
	senport := ""
	gtg := 3
	if err := l.Print(1, 2, "Pi Test Vehicle"); err != nil {
		fmt.Println(err)
		return
	}
	ports, err := serial.GetPortsList()
	if err != nil {
		log.Fatal(err)
	}
	if len(ports) == 0 {
		log.Fatal("No Serial ports found!")
	}
	for _, port := range ports {
		fmt.Printf("Found port: %v\n", port)
	}
	mode := &serial.Mode{
		BaudRate: 9600,
		Parity:   serial.NoParity,
		DataBits: 8,
		StopBits: serial.OneStopBit,
	}
	dp := fmt.Sprintf("Detected %d Ports", len(ports))
	if err := l.Print(2, 2, dp); err != nil {
		fmt.Println(err)
		return
	}

	for x := 0; x < len(ports); x++ {
		port, err := serial.Open(ports[x], mode)
		if err != nil {
			log.Fatal(err)
		}
		if err := l.Print(3, 2, "                  "); err != nil {
			fmt.Println(err)
			return
		}
		dp := fmt.Sprintf("Detecting Port %d", x+1)
		if err := l.Print(3, 2, dp); err != nil {
			fmt.Println(err)
			return
		}
		line := ""
		buff := make([]byte, 1)
		for {
			n, err := port.Read(buff)
			if err != nil {
				log.Fatal(err)
			}
			if n == 0 {
				port.Close()
				break
			}
			line = line + string(buff[:n])
			if strings.Contains(string(buff[:n]), "\n") {
				port.Close()
				break
			}

		}
		if len(line) > 3 {
			switch {
			case line[0:3] == "$GP":
				gpsport = ports[x]
			case line[0:3] == "CH1":
				rcport = ports[x]
			case line[0:3] == "DIS":
				senport = ports[x]

			}

		}

	}
	if err := l.Print(2, 2, "                  "); err != nil {
		fmt.Println(err)
		return
	}
	if err := l.Print(3, 2, "                  "); err != nil {
		fmt.Println(err)
		return
	}

	if len(gpsport) > 0 {
		fmt.Printf("GPS Port %s\n", gpsport)
		gtg--
		if err := l.Print(2, 2, "GPS Port Found"); err != nil {
			fmt.Println(err)
			return
		}
	} else {
		fmt.Printf("GPS Port Not Found\n")
		if err := l.Print(2, 2, "GPS Port Fail"); err != nil {
			fmt.Println(err)
			return
		}
	}
	if len(rcport) > 0 {
		fmt.Printf("RC Port %s\n", rcport)
		gtg--
		if err := l.Print(3, 2, "RC Port Found"); err != nil {
			fmt.Println(err)
			return
		}
	} else {
		fmt.Printf("RC Port Not Found\n")
		if err := l.Print(3, 2, "RC Port Fail"); err != nil {
			fmt.Println(err)
			return
		}
	}
	if len(senport) > 0 {
		fmt.Printf("Sensor Port %s\n", senport)
		gtg--
		if err := l.Print(4, 2, "Sensor Port Found"); err != nil {
			fmt.Println(err)
			return
		}
	} else {
		fmt.Printf("Sensor Port Not Found\n")
		if err := l.Print(4, 2, "Sensor Port Fail"); err != nil {
			fmt.Println(err)
			return
		}
	}
	if gtg == 0 {
		fmt.Printf("Good to go\n")
		fmt.Printf("All ports found\n")
		if err := l.Print(1, 2, "                  "); err != nil {
			fmt.Println(err)
			return
		}
		if err := l.Print(2, 2, "                  "); err != nil {
			fmt.Println(err)
			return
		}
		if err := l.Print(3, 2, "                  "); err != nil {
			fmt.Println(err)
			return
		}
		if err := l.Print(4, 2, "                  "); err != nil {
			fmt.Println(err)
			return
		}
		if err := l.Print(2, 1, "Vehicle Good to Go"); err != nil {
			fmt.Println(err)
			return
		}
		if err := l.Print(3, 2, "All Ports Found"); err != nil {
			fmt.Println(err)
			return
		}
		if err := l.Print(1, 2, "                  "); err != nil {
			fmt.Println(err)
			return
		}
		if err := l.Print(2, 2, "                  "); err != nil {
			fmt.Println(err)
			return
		}
		if err := l.Print(3, 2, "                  "); err != nil {
			fmt.Println(err)
			return
		}
		if err := l.Print(4, 2, "                  "); err != nil {
			fmt.Println(err)
			return
		}
		porta, err := serial.Open(gpsport, mode)
		if err != nil {
			log.Fatal(err)
		}
		portb, err := serial.Open(senport, mode)
		if err != nil {
			log.Fatal(err)
		}
		portc, err := serial.Open(rcport, mode)
		if err != nil {
			log.Fatal(err)
		}

		for {

			linea := ""
			lineb := ""
			linec := ""
			la := true
			lb := true
			lc := true
			buff := make([]byte, 1)
			on := true
			for on != false {
				linea = ""
				lineb = ""
				linec = ""
				if la {
					for {
						n, err := porta.Read(buff)
						if err != nil {
							log.Fatal(err)
						}
						if n == 0 {
							fmt.Println("\nEOF")
							break
						}
						linea = linea + string(buff[:n])
						if strings.Contains(string(buff[:n]), "\n") {
							break
						}

					}
				}
				if lb {
					for {
						n, err := portb.Read(buff)
						if err != nil {
							log.Fatal(err)
						}
						if n == 0 {
							fmt.Println("\nEOF")
							break
						}
						lineb = lineb + string(buff[:n])
						if strings.Contains(string(buff[:n]), "\n") {
							break
						}

					}
				}
				if lc {
					for {
						n, err := portc.Read(buff)
						if err != nil {
							log.Fatal(err)
						}
						if n == 0 {
							fmt.Println("\nEOF")
							break
						}
						linec = linec + string(buff[:n])
						if strings.Contains(string(buff[:n]), "\n") {
							break
						}

					}
				}

				fmt.Printf("%s  - %s - %s \n", linea[0:18], lineb, linec)
				if err := l.Print(2, 0, "                  "); err != nil {
					fmt.Println(err)
					return
				}
				if err := l.Print(3, 0, "                  "); err != nil {
					fmt.Println(err)
					return
				}
				if err := l.Print(4, 0, "                  "); err != nil {
					fmt.Println(err)
					return
				}
				if la {
					if err := l.Print(2, 0, linea[0:18]); err != nil {
						fmt.Println(err)
						return
					}
				}
				if lb {
					if err := l.Print(3, 0, lineb); err != nil {
						fmt.Println(err)
						return
					}
				}
				if lc {
					if err := l.Print(4, 0, linec[0:18]); err != nil {
						fmt.Println(err)
						return
					}
				}
			}
		}

	} else {
		fmt.Printf("Init Failure  %d ports not found\n", gtg)

		if err := l.Print(1, 0, "                    "); err != nil {
			fmt.Println(err)
			return
		}
		if err := l.Print(1, 0, "Vehicle Failed Init"); err != nil {
			fmt.Println(err)
			return
		}
	}
}

