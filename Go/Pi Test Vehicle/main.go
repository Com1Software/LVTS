package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"runtime"
	"strings"
	"time"

	"go.bug.st/serial"
)

// -------------------------------------------------------------------------
func main() {
	agent := SSE()
	xip := fmt.Sprintf("%s", GetOutboundIP())
	port := "8080"
	//	s1 := rand.NewSource(time.Now().UnixNano())
	//	r1 := rand.New(s1)
	//
	//--- tctl 0 = normal mode test
	//         1 = high speed mode test
	//
	tctl := 1
	tc := 0
	fmt.Println("Test SSE Server")
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

				agent.Notifier <- []byte(line)
				fmt.Println(line)
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
