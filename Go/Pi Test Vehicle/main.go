package main

import (
	"bufio"
	"encoding/hex"
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
var xh int64 = 0
var heading int64 = 180

// -------------------------------------------------------------------------
func main() {
	agent := SSE()
	cmdctl := []CommandCtl{}
	cmdctla := CommandCtl{}
	xip := fmt.Sprintf("%s", GetOutboundIP())
	port := "8080"
	headingAngle := 1
	imurec := false
	serialIMU := true
	exefile := "IMUReceiver/test"
	imufile := "rec.dat"
	//--- tctl 0 = normal mode test
	//         1 = high speed mode test
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
	if imurec {
		cmd := exec.Command(exefile)
		if err := cmd.Run(); err != nil {
			fmt.Printf("Command %s \n Error: %s\n", cmd, err)
		}
		fmt.Println(cmd)
		if _, err := os.Stat(drivefile); err == nil {
			lines, err := readLines(imufile)
			data := strings.Split(strings.Join(lines, " "), "=")
			if len(data) > 1 {
				if err != nil {
					fmt.Printf("IMU File Load Error : ( %s )\n", err)
				}
				i, err := strconv.Atoi(data[1])
				if err != nil {
					panic(err)
				}
				headingAngle = i
				fmt.Printf("IMU File Present Heading set to %d\n", i)
			}
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
		BaudRate: 115200,
		Parity:   serial.NoParity,
		DataBits: 8,
		StopBits: serial.OneStopBit,
	}
	imuport := ""
	gpsport := ""
	rcport := ""
	if serialIMU {
		for x := 0; x < len(ports); x++ {
			port, err := serial.Open(ports[x], mode)
			if err != nil {
				log.Fatal(err)
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

				src := []byte(string(buff))
				encodedStr := hex.EncodeToString(src)
				if encodedStr == "55" {
					imuport = ports[x]
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
				}

			}
		}
	}
	fmt.Printf("IMU %s - GPS %s - RC %s \n", imuport, gpsport, rcport)

	headingAngle = SerialIMU(imuport, xip, mode)

	if len(gpsport) > 0 {
		//for x := 0; x < len(ports); x++ {

		port, err := serial.Open(gpsport, mode)
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
				if imurec {
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
				}

				line := ""
				ok := false
				buff := make([]byte, 1)
				for {
					//------------------------
					//headingAngle = SerialIMU(imuport, xip, mode)
					//-------------------

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
							headingAngle = SerialIMU(imuport, xip, mode)
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

func GetHeading(p int64) int64 {
	switch {

	case p == 0:
		xh = 0
	case p == 1:
		xh = 360
	case p == 2:
		xh = 359
	case p == 3:
		xh = 358
	case p == 4:
		xh = 356
	case p == 5:
		xh = 354
	case p == 6:
		xh = 352
	case p == 7:
		xh = 350
	case p == 8:
		xh = 348
	case p == 9:
		xh = 346
	case p == 10:
		xh = 344
	case p == 11:
		xh = 342
	case p == 12:
		xh = 340
	case p == 13:
		xh = 338
	case p == 14:
		xh = 336
	case p == 15:
		xh = 334
	case p == 16:
		xh = 332
	case p == 17:
		xh = 330
	case p == 18:
		xh = 328
	case p == 19:
		xh = 326
	case p == 20:
		xh = 322
	case p == 21:
		xh = 330
	case p == 22:
		xh = 328
	case p == 23:
		xh = 326
	case p == 24:
		xh = 324
	case p == 25:
		xh = 322
	case p == 26:
		xh = 321
	case p == 27:
		xh = 320
	case p == 28:
		xh = 318
	case p == 29:
		xh = 316
	case p == 30:
		xh = 314
	case p == 31:
		xh = 312
	case p == 32:
		xh = 310
	case p == 33:
		xh = 308
	case p == 34:
		xh = 307
	case p == 35:
		xh = 306
	case p == 36:
		xh = 304
	case p == 37:
		xh = 302
	case p == 38:
		xh = 300
	case p == 39:
		xh = 298
	case p == 40:
		xh = 296
	case p == 41:
		xh = 294
	case p == 42:
		xh = 292
	case p == 43:
		xh = 290
	case p == 44:
		xh = 289
	case p == 45:
		xh = 288
	case p == 46:
		xh = 287
	case p == 47:
		xh = 286
	case p == 48:
		xh = 285
	case p == 49:
		xh = 284
	case p == 50:
		xh = 283
	case p == 51:
		xh = 282
	case p == 52:
		xh = 281
	case p == 53:
		xh = 280
	case p == 54:
		xh = 279
	case p == 55:
		xh = 278
	case p == 56:
		xh = 277
	case p == 57:
		xh = 276
	case p == 58:
		xh = 275
	case p == 59:
		xh = 274
	case p == 60:
		xh = 273
	case p == 61:
		xh = 272
	case p == 62:
		xh = 271
	case p == 63:
		xh = 270
	case p == 64:
		xh = 269
	case p == 65:
		xh = 268
	case p == 66:
		xh = 267
	case p == 67:
		xh = 266
	case p == 68:
		xh = 265
	case p == 69:
		xh = 264
	case p == 70:
		xh = 263
	case p == 71:
		xh = 262
	case p == 72:
		xh = 261
	case p == 73:
		xh = 260
	case p == 74:
		xh = 259
	case p == 75:
		xh = 258
	case p == 76:
		xh = 257
	case p == 77:
		xh = 256
	case p == 78:
		xh = 255
	case p == 79:
		xh = 254
	case p == 80:
		xh = 253
	case p == 81:
		xh = 252
	case p == 82:
		xh = 251
	case p == 83:
		xh = 250
	case p == 84:
		xh = 248
	case p == 85:
		xh = 246
	case p == 86:
		xh = 244
	case p == 87:
		xh = 242
	case p == 88:
		xh = 240
	case p == 89:
		xh = 239
	case p == 90:
		xh = 238
	case p == 91:
		xh = 237
	case p == 92:
		xh = 236
	case p == 93:
		xh = 235
	case p == 94:
		xh = 234
	case p == 95:
		xh = 233
	case p == 96:
		xh = 232
	case p == 97:
		xh = 231
	case p == 98:
		xh = 230
	case p == 99:
		xh = 218
	case p == 100:
		xh = 214
	case p == 101:
		xh = 212
	case p == 102:
		xh = 210
	case p == 103:
		xh = 209
	case p == 104:
		xh = 208
	case p == 105:
		xh = 207
	case p == 106:
		xh = 206
	case p == 107:
		xh = 205
	case p == 108:
		xh = 204
	case p == 109:
		xh = 203
	case p == 110:
		xh = 202
	case p == 111:
		xh = 201
	case p == 112:
		xh = 200
	case p == 113:
		xh = 198
	case p == 114:
		xh = 196
	case p == 115:
		xh = 194
	case p == 116:
		xh = 192
	case p == 117:
		xh = 190
	case p == 118:
		xh = 189
	case p == 119:
		xh = 188
	case p == 120:
		xh = 187
	case p == 121:
		xh = 186
	case p == 122:
		xh = 185
	case p == 123:
		xh = 184
	case p == 124:
		xh = 183
	case p == 125:
		xh = 182
	case p == 126:
		xh = 181
	case p == 127:
		xh = 180
	case p == 128:
		xh = 176
	case p == 129:
		xh = 174
	case p == 130:
		xh = 170
	case p == 131:
		xh = 169
	case p == 132:
		xh = 168
	case p == 133:
		xh = 167
	case p == 134:
		xh = 166
	case p == 135:
		xh = 165
	case p == 136:
		xh = 164
	case p == 137:
		xh = 163
	case p == 138:
		xh = 162
	case p == 139:
		xh = 161
	case p == 140:
		xh = 160
	case p == 141:
		xh = 158
	case p == 142:
		xh = 156
	case p == 143:
		xh = 154
	case p == 144:
		xh = 152
	case p == 145:
		xh = 150
	case p == 146:
		xh = 149
	case p == 147:
		xh = 148
	case p == 148:
		xh = 147
	case p == 149:
		xh = 146
	case p == 150:
		xh = 145
	case p == 151:
		xh = 144
	case p == 152:
		xh = 143
	case p == 153:
		xh = 142
	case p == 154:
		xh = 141
	case p == 155:
		xh = 140
	case p == 156:
		xh = 138
	case p == 157:
		xh = 136
	case p == 158:
		xh = 134
	case p == 159:
		xh = 132
	case p == 160:
		xh = 130
	case p == 161:
		xh = 129
	case p == 162:
		xh = 128
	case p == 163:
		xh = 127
	case p == 164:
		xh = 126
	case p == 165:
		xh = 125
	case p == 166:
		xh = 124
	case p == 167:
		xh = 123
	case p == 168:
		xh = 122
	case p == 169:
		xh = 121
	case p == 170:
		xh = 120
	case p == 171:
		xh = 118
	case p == 172:
		xh = 116
	case p == 173:
		xh = 114
	case p == 174:
		xh = 112
	case p == 175:
		xh = 110
	case p == 176:
		xh = 109
	case p == 177:
		xh = 108
	case p == 178:
		xh = 107
	case p == 179:
		xh = 106
	case p == 180:
		xh = 105
	case p == 181:
		xh = 104
	case p == 182:
		xh = 103
	case p == 183:
		xh = 102
	case p == 184:
		xh = 101
	case p == 185:
		xh = 100
	case p == 186:
		xh = 98
	case p == 187:
		xh = 96
	case p == 188:
		xh = 94
	case p == 189:
		xh = 92
	case p == 190:
		xh = 90
	case p == 191:
		xh = 88
	case p == 192:
		xh = 85
	case p == 193:
		xh = 84
	case p == 194:
		xh = 83
	case p == 195:
		xh = 82
	case p == 196:
		xh = 81
	case p == 197:
		xh = 80
	case p == 198:
		xh = 79
	case p == 199:
		xh = 78
	case p == 200:
		xh = 77
	case p == 201:
		xh = 76
	case p == 202:
		xh = 75
	case p == 203:
		xh = 74
	case p == 204:
		xh = 73
	case p == 205:
		xh = 72
	case p == 206:
		xh = 71
	case p == 207:
		xh = 70
	case p == 208:
		xh = 68
	case p == 209:
		xh = 64
	case p == 210:
		xh = 62
	case p == 211:
		xh = 60
	case p == 212:
		xh = 59
	case p == 213:
		xh = 58
	case p == 214:
		xh = 57
	case p == 215:
		xh = 56
	case p == 216:
		xh = 55
	case p == 217:
		xh = 54
	case p == 218:
		xh = 53
	case p == 219:
		xh = 52
	case p == 220:
		xh = 51
	case p == 221:
		xh = 50
	case p == 222:
		xh = 48
	case p == 223:
		xh = 46
	case p == 224:
		xh = 44
	case p == 225:
		xh = 42
	case p == 226:
		xh = 40
	case p == 227:
		xh = 39
	case p == 228:
		xh = 38
	case p == 229:
		xh = 37
	case p == 230:
		xh = 36
	case p == 231:
		xh = 35
	case p == 232:
		xh = 34
	case p == 233:
		xh = 33
	case p == 234:
		xh = 32
	case p == 235:
		xh = 31
	case p == 236:
		xh = 30
	case p == 237:
		xh = 28
	case p == 238:
		xh = 26
	case p == 239:
		xh = 24
	case p == 240:
		xh = 22
	case p == 241:
		xh = 20
	case p == 242:
		xh = 19
	case p == 243:
		xh = 18
	case p == 244:
		xh = 17
	case p == 245:
		xh = 16
	case p == 246:
		xh = 15
	case p == 247:
		xh = 14
	case p == 248:
		xh = 13
	case p == 249:
		xh = 12
	case p == 250:
		xh = 11
	case p == 251:
		xh = 10
	case p == 252:
		xh = 8
	case p == 253:
		xh = 6
	case p == 254:
		xh = 4
	case p == 255:
		xh = 2

	}
	return xh
}

func SerialIMU(imuport string, xip string, mode *serial.Mode) int {
	pos := 0
	mctl := 0
	headingctl := false
	headingAngle := 1

	porta, err := serial.Open(imuport, mode)
	if err != nil {
		log.Fatal(err)
	}
	for {
		buff := make([]byte, 1)
		n, err := porta.Read(buff)
		if err != nil {
			log.Fatal(err)
		}
		if n == 0 {
			fmt.Println("Cloed Serial Port")
			porta.Close()

		}
		src := []byte(string(buff))
		encodedStr := hex.EncodeToString(src)
		if encodedStr == "55" {
			mctl = 1
			pos = 0
		}
		if encodedStr == "53" && mctl == 1 {
			mctl = 2
		}
		if mctl == 2 {
			pos++
			decimal, err := strconv.ParseInt(encodedStr, 16, 32)
			if err != nil {
				fmt.Println(err)
			}
			if pos == 7 {
				if headingctl {
					heading = GetHeading(decimal)
				} else {
					heading = GetHeading(decimal)
				}
				fmt.Printf(" Using IP: %s Heading= %d\n", xip, heading)
				headingAngle = int(heading)
				break
			}
		}
	}
	porta.Close()
	return headingAngle
}
