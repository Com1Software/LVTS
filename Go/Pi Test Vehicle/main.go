package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/googolgl/go-i2c"
	"github.com/googolgl/go-pca9685"
	"go.bug.st/serial"
)

const (
	RegisterA    = 0
	RegisterB    = 0x01
	RegisterMode = 0x02
	XAxisH       = 0x03
	ZAxisH       = 0x05
	YAxisH       = 0x07
	Declination  = -0.00669
	pi           = 3.14159265359
)

var bus *i2c.Options

// -------------------------------------------------------------------------
func main() {
	agent := SSE()
	cmdctl := []CommandCtl{}
	cmdctla := CommandCtl{}
	xip := fmt.Sprintf("%s", GetOutboundIP())
	port := "8080"
	//--- tctl 0 = normal mode test
	//         1 = high speed mode test
	headingAngle := 1
	tctl := 1
	tc := 0
	psp := 0
	drivefile := "drive.ctl"
	fmt.Println("Pi Test Vehicle")
	fmt.Printf("Operating System : %s\n", runtime.GOOS)
	fmt.Printf("Outbound IP  : %s Port : %s\n", xip, port)
	if _, err := os.Stat(drivefile); err == nil {
		fmt.Println("Drive File Present")
		lines, err := readLines(drivefile)
		cmdctla.cmdlines = lines
		cmdctl = append(cmdctl, cmdctla)
		cmdctl[0].cmdpos = 0
		cmdctl[0].cmdstatus = "ready"
		cmdctl[0].throttle = 0
		cmdctl[0].steering = 0
		if err != nil {
			fmt.Printf("Drive File Load Error : ( %s )\n", err)
		} else {
			fmt.Printf("Drive File Present Loaded %d Lines)\n", len(cmdctl[0].cmdlines))
		}
	} else {
		fmt.Println("drive.ctl Drive File Not Present")
	}
	exefile := "IMUReceiver/test"
	cmd := exec.Command(exefile)
	if err := cmd.Run(); err != nil {
		fmt.Printf("Command %s \n Error: %s\n", cmd, err)
	}
	fmt.Println(cmd)
	//imufile := "IMUReceiver/rec.dat"
	imufile := "rec.dat"
	if _, err := os.Stat(drivefile); err == nil {
		lines, err := readLines(imufile)
		data := strings.Split(strings.Join(lines, " "), "=")
		if len(data) > 1 {
			if err != nil {
				fmt.Printf("IMU File Load Error : ( %s )\n", err)
			}
			i, err := strconv.Atoi(data[1])
			if err != nil {
				// ... handle error
				panic(err)
			}
			headingAngle = i
			fmt.Printf("IMU File Present Heading set to %d\n", i)
		}
	}

	if runtime.GOOS == "windows" {
		xip = "http://localhost"
	}

	i2cs, err := i2c.New(pca9685.Address, "/dev/i2c-1")
	if err != nil {
		log.Fatal(err)
	}
	pca0, err := pca9685.New(i2cs, nil)
	if err != nil {
		log.Fatal(err)
	}
	pca0.SetChannel(0, 0, 130)
	servo0 := pca0.ServoNew(0, nil)

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

		bus, _ = i2c.New(0x1c, "/dev/i2c-1")
		defer bus.Close()
		//	MagnetometerInit()
		fmt.Println("Reading Heading Angle")
		go func() {
			for {
				switch {
				case tctl == 0:
					time.Sleep(time.Second * 1)
				case tctl == 1:
					time.Sleep(time.Second * -1)
					tc++
				}

				cmd := exec.Command(exefile)
				if err := cmd.Run(); err != nil {
					fmt.Printf("Command %s \n Error: %s\n", cmd, err)
				}
				if _, err := os.Stat(drivefile); err == nil {
					lines, err := readLines(imufile)
					data := strings.Split(strings.Join(lines, " "), "=")
					if len(data) > 1 {
						if err != nil {
							fmt.Printf("IMU File Load Error : ( %s )\n", err)
						}
						i, err := strconv.Atoi(data[1])
						if err != nil {
							// ... handle error
							panic(err)
						}
						headingAngle = i
					}
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

							event := fmt.Sprintf("%s  latitude=%s  %s   longitude=%s %s knots=%s degrees=%s ", id, latitude, ns, longitude, ew, gpsspeed, degree)
							event = event + fmt.Sprintf("Heading Angle = %d \n", headingAngle)
							fmt.Printf(event)
							cmdctl[0].longitude = longitude
							cmdctl[0].latitude = latitude
							cmdctl[0].heading = headingAngle
							cmdctl = ProcessControl(cmdctl)
							drivectl := fmt.Sprintf(" CMD=%s CMDPos=%d CMDStatus=%s Steering=%d  Throttle=%d Command Lines=%d\n", cmdctl[0].cmd, cmdctl[0].cmdpos, cmdctl[0].cmdstatus, cmdctl[0].steering, cmdctl[0].throttle, len(cmdctl[0].cmdlines))
							event = event + drivectl
							fmt.Printf(drivectl)
							if psp != cmdctl[0].steering {
								if psp > cmdctl[0].steering {
									for i := psp; i < cmdctl[0].steering; i++ {
										servo0.Angle(i)
										time.Sleep(10 * time.Millisecond)
									}
									psp = cmdctl[0].steering
								} else {
									for i := cmdctl[0].steering; i < psp; i-- {
										servo0.Angle(i)
										time.Sleep(10 * time.Millisecond)
									}
									psp = cmdctl[0].steering

								}
								servo0.Angle(cmdctl[0].steering)
							}
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
						if i > -300 {
							if i > 140 {
								fmt.Println("Right")
								if i > 190 {
									fmt.Println("Far Right")
									if i > 220 {
										fmt.Println("Super Far Right")
									}

								}

							}
							if i < 100 {
								fmt.Println("Left")
								if i < 60 {
									fmt.Println("Far Left")
									if i < 10 {
										fmt.Println("Super Far Left")
									}

								}
							}
							i, e = strconv.Atoi(ch2)
							if e != nil {
								fmt.Println(e)
							}
							if i > 150 {
								fmt.Println("Forward")
								if i > 200 {
									fmt.Println("Fast Forward")
									if i > 250 {
										fmt.Println("Super Fast Forward")
									}
								}

							}
							if i < 100 {
								fmt.Println("Backward")
								if i < 50 {
									fmt.Println("Fast Backward")
									if i < 10 {
										fmt.Println("Super Fast Backward")
									}

								}
							}

							fmt.Printf("CH1=%s CH2=%s CH3=%s CH4=%s\n", ch1, ch2, ch3, ch4)
						} else {
							fmt.Println("Conroller is off")
						}
					case line[0:3] == "DIS":
						ok = false
						dis1, pos1, min1, max1, dis2, pos2, min2, max2 := getSenPosition(line)
						fmt.Printf("DIS1=%s POS1=%s MIN1=%s MAX1=%s DIS2=%s POS2=%s MIN2=%s MAX2=%s\n", dis1, pos1, min1, max1, dis2, pos2, min2, max2)

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

type CommandCtl struct {
	cmdlines       []string
	cmd            string
	cmdpos         int
	cmdstatus      string
	steering       int
	throttle       int
	heading        int
	latitude       string
	longitude      string
	priorlatitude  string
	priorlongitude string
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

func getSenPosition(sentence string) (string, string, string, string, string, string, string, string) {
	data := strings.Split(sentence, ",")
	dis1 := ""
	pos1 := ""
	min1 := ""
	max1 := ""
	dis2 := ""
	pos2 := ""
	min2 := ""
	max2 := ""

	if string(data[0][0:4]) == "DIS1" {
		data1 := strings.Split(data[0], "=")
		dis1 = data1[1]
		data1 = strings.Split(data[1], "=")
		pos1 = data1[1]
		data1 = strings.Split(data[2], "=")
		min1 = data1[1]
		data1 = strings.Split(data[3], "=")
		max1 = data1[1]
		data1 = strings.Split(data[4], "=")
		dis2 = data1[1]
		data1 = strings.Split(data[5], "=")
		pos2 = data1[1]
		data1 = strings.Split(data[6], "=")
		min2 = data1[1]
		data1 = strings.Split(data[7], "=")
		max2 = data1[1]

	}

	return dis1, pos1, min1, max1, dis2, pos2, min2, max2
}

func MagnetometerInit() {
	bus.WriteRegU8(RegisterA, 0x70)
	bus.WriteRegU8(RegisterB, 0xa0)
	bus.WriteRegU8(RegisterMode, 0)
}

func readRawData(addr byte) int {
	high, _ := bus.ReadRegU8(addr)
	low, _ := bus.ReadRegU8(addr + 1)
	value := int(int16(high)<<8 | int16(low))
	if value > 32768 {
		value -= 65536
	}
	return value
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func ProcessControl(cmdctl []CommandCtl) []CommandCtl {
	if len(cmdctl[0].cmdlines) > 0 {
		if cmdctl[0].cmdpos == 0 && cmdctl[0].cmdstatus == "ready" {
			cmdctl[0].cmd = cmdctl[0].cmdlines[cmdctl[0].cmdpos]
			cmdctl[0].cmdstatus = "loaded"

		}
		data := strings.Split(cmdctl[0].cmd, " ")
		switch {
		case data[0] == "align" && len(data) == 2:
			switch {
			case data[1] == "N":
				if cmdctl[0].heading == 0 {
					cmdctl[0].throttle = 0
					cmdctl[0].steering = 0
					cmdctl[0].cmdstatus = "align complete"
					cmdctl[0].cmdpos++
					cmdctl[0].cmd = cmdctl[0].cmdlines[cmdctl[0].cmdpos]
					cmdctl[0].cmdstatus = "loading next command"

				} else {
					cmdctl[0].cmdstatus = "running align"
					if cmdctl[0].heading > 180 {
						cmdctl[0].throttle = 10
						cmdctl[0].steering = 20
					} else {
						cmdctl[0].throttle = 10
						cmdctl[0].steering = 110
					}
				}
			case data[1] == "S":

			case data[1] == "E":

			case data[1] == "W":
			}

		case data[0] == "drive" && len(data) == 2:

		}
	}
	return cmdctl
}
