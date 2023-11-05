package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"time"

	"go.bug.st/serial"
)

// -------------------------------------------------------------------------
func main() {
	agent := SSE()
	xip := fmt.Sprintf("%s", GetOutboundIP())
	port := "8080"
	//--- tctl 0 = normal mode test
	//         1 = high speed mode test
	tctl := 1
	tc := 0
	fmt.Println("Pi Test Vehicle")
	fmt.Printf("Operating System : %s\n", runtime.GOOS)
	fmt.Printf("Outbound IP  : %s Port : %s\n", xip, port)
	if runtime.GOOS == "windows" {
		xip = "http://localhost"
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
	for x := 0; x < len(ports); x++ {
		port, err := serial.Open(ports[x], mode)
		if err != nil {
			log.Fatal(err)
		}

		go func() {
			for {
				switch {
				case tctl == 0:
					time.Sleep(time.Second * 1)
				case tctl == 1:
					time.Sleep(time.Second * -1)
					tc++
				}
				line := ""
				ok := false
				buff := make([]byte, 1)
				for {
					n, err := port.Read(buff)
					if err != nil {
						log.Fatal(err)
					}
					if n == 0 {
						port.Close()
						fmt.Println("closed Port")
						break
					}
					if strings.Contains(string(buff[:n]), "\n") {
						break
					} else {
						line = line + string(buff[:n])
					}
				}
				if len(line) > 2 {
					switch {
					case line[0:3] == "$GP":
						ok = false
						id, latitude, longitude, ns, ew, gpsspeed, degree := getGPSPosition(line)
						if len(id) > 0 {
							event := fmt.Sprintf("%s  latitude=%s  %s   longitude=%s %s knots=%s degrees=%s\n", id, latitude, ns, longitude, ew, gpsspeed, degree)
							fmt.Println(event)
							agent.Notifier <- []byte(event)
						}
					case line[0:3] == "CH1":
						ok = false
						ch1, ch2, ch3, ch4 := getCHPosition(line)
						fmt.Print("\033[u\033[K")
						i, e := strconv.Atoi(ch1)
						if e != nil {
							fmt.Println(e)
						}
						if i > 150 {
							fmt.Println("Right")
						}
						if i < 100 {
							fmt.Println("Left")
						}
						i, e = strconv.Atoi(ch2)
						if e != nil {
							fmt.Println(e)
						}
						if i > 150 {
							fmt.Println("Forward")
						}
						if i < 100 {
							fmt.Println("Backward")
						}
						fmt.Printf("CH1=%s CH2=%s CH3=%s CH4=%s\n", ch1, ch2, ch3, ch4)

					case line[1:4] == "DIS":
						ok = false
					case line[0:3] == "POS":
						ok = false

					}
				}
				if ok {
					agent.Notifier <- []byte(line)
					fmt.Println(line)
				}
				ok = false
			}
		}()
	}

	fmt.Printf("Listening at  : %s Port : %s\n", xip, port)
	if runtime.GOOS == "windows" {
		http.ListenAndServe(":"+port, agent)
	} else {
		http.ListenAndServe(xip+":"+port, agent)
	}
}

type Agent struct {
	Notifier    chan []byte
	newuser     chan chan []byte
	closinguser chan chan []byte
	user        map[chan []byte]bool
}

func SSE() (agent *Agent) {
	agent = &Agent{
		Notifier:    make(chan []byte, 1),
		newuser:     make(chan chan []byte),
		closinguser: make(chan chan []byte),
		user:        make(map[chan []byte]bool),
	}
	go agent.listen()
	return
}

func (agent *Agent) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	flusher, ok := rw.(http.Flusher)
	if !ok {
		http.Error(rw, "Error ", http.StatusInternalServerError)
		return
	}
	rw.Header().Set("Content-Type", "text/event-stream")
	rw.Header().Set("Cache-Control", "no-cache")
	rw.Header().Set("Connection", "keep-alive")
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	mChan := make(chan []byte)
	agent.newuser <- mChan
	defer func() {
		agent.closinguser <- mChan
	}()
	notify := req.Context().Done()
	go func() {
		<-notify
		agent.closinguser <- mChan
	}()
	for {
		fmt.Fprintf(rw, "%s", <-mChan)
		flusher.Flush()
	}

}

func (agent *Agent) listen() {
	for {
		select {
		case s := <-agent.newuser:
			agent.user[s] = true
		case s := <-agent.closinguser:
			delete(agent.user, s)
		case event := <-agent.Notifier:
			for userMChan, _ := range agent.user {
				userMChan <- event
			}
		}
	}

}

func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

func getGPSPosition(sentence string) (string, string, string, string, string, string, string) {
	data := strings.Split(sentence, ",")
	id := ""
	latitude := ""
	longitude := ""
	ns := ""
	ew := ""
	speed := ""
	degree := ""

	switch {
	case string(data[0]) == "$GPGGA":
		//	id = data[0]
		latitude = data[2]
		ns = data[3]
		longitude = data[4]
		ew = data[5]

	case string(data[0]) == "$GPGLL":
		//	id = data[0]
		latitude = data[1]
		ns = data[2]
		longitude = data[3]
		ew = data[4]

	case string(data[0]) == "$GPVTG":
		//	id = data[0]
		degree = data[1]

	case string(data[0]) == "$GPRMC":
		id = data[0]
		latitude = data[3]
		ns = data[4]
		longitude = data[5]
		ew = data[6]
		speed = data[7]
		degree = data[8]

	case string(data[0]) == "$GPGSA":
	//	id = data[0]

	case string(data[0]) == "$GPGSV":
	//	id = data[0]

	case string(data[0]) == "$GPTXT":
	//	id = data[0]

	default:
		fmt.Println("-- %s", data[0])

	}

	return id, latitude, longitude, ns, ew, speed, degree
}

func getCHPosition(sentence string) (string, string, string, string) {
	data := strings.Split(sentence, ",")
	ch1 := ""
	ch2 := ""
	ch3 := ""
	ch4 := ""
	if len(data) == 4 {
		ch1data := strings.Split(data[0], "=")
		ch2data := strings.Split(data[1], "=")
		ch3data := strings.Split(data[2], "=")
		ch4data := strings.Split(data[3], "=")

		if string(ch1data[0]) == "CH1" {
			ch1 = ch1data[1]
		}
		if string(ch2data[0]) == "CH2" {
			ch2 = ch2data[1]
		}
		if string(ch3data[0]) == "CH3" {
			ch3 = ch3data[1]
		}
		if string(ch4data[0]) == "CH4" {
			ch4 = ch4data[1]
		}
	}
	return ch1, ch2, ch3, ch4

}
