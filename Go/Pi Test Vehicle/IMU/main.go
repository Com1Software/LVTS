package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"runtime"
	"strconv"

	"go.bug.st/serial"
)

var xh int64 = 0
var heading int64 = 180

func main() {
	fmt.Println("Test Multiport Serial Controller")
	fmt.Printf("Operating System : %s\n", runtime.GOOS)

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
	mctl := 0

	headingctl := false
	po := false
	pos := 0
	for x := 0; x < len(ports); x++ {
		openport := ports[x]
		if runtime.GOOS == "linux" {
			openport = "/dev/ttyUSB0"
		}
		port, err := serial.Open(openport, mode)
		if err != nil {
			fmt.Println(err)
			po = false
		} else {
			po = true
		}
		if po {
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
					mctl = 1
					pos = 0
				}
				if encodedStr == "53" && mctl == 1 {
					fmt.Println("")
					mctl = 2
				}

				if mctl == 2 {
					pos++

					decimal, err := strconv.ParseInt(encodedStr, 16, 32)
					if err != nil {
						fmt.Println(err)
					}

					if pos == 7 {
						fmt.Printf(" zL= %d", decimal)
						if headingctl {

							heading = GetHeading(decimal)
						} else {
							heading = GetHeading(decimal)
						}

						fmt.Printf(" Heading= %d", heading)

					}
				}
			}
		}
	}

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
