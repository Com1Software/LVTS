package main

import (
	"fmt"
	"log"
	"strings"

	"go.bug.st/serial"
)

func main() {
	fmt.Println("Pi Test Vehicle")
	gpsport := ""
	rcport := ""
	senport := ""
	gtg := 3
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
	fmt.Println(dp)
	for x := 0; x < len(ports); x++ {
		if ports[x][8:11] == "ACM" {
			port, err := serial.Open(ports[x], mode)
			if err != nil {
				log.Fatal(err)
			}
			dp := fmt.Sprintf("Detecting Port %d", x+1)
			fmt.Println(dp)
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

	}

	if len(gpsport) > 0 {
		fmt.Printf("GPS Port %s\n", gpsport)
		gtg--
	} else {
		fmt.Printf("GPS Port Not Found\n")
	}
	if len(rcport) > 0 {
		fmt.Printf("RC Port %s\n", rcport)
		gtg--
	} else {
		fmt.Printf("RC Port Not Found\n")
	}
	if len(senport) > 0 {
		fmt.Printf("Sensor Port %s\n", senport)
		gtg--
	} else {
		fmt.Printf("Sensor Port Not Found\n")
	}
	if gtg == 0 {
		fmt.Printf("Good to go\n")
		fmt.Printf("All ports found\n")
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

			}
		}

	} else {
		fmt.Printf("Init Failure  %d ports not found\n", gtg)
	}
}
