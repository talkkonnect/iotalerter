package iotalerter

import (
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/talkkonnect/colog"
	term "github.com/talkkonnect/termbox-go"
)

func Init(configFile string) {

	err := term.Init()
	if err != nil {
		log.Println("alert: Cannot Initalize Terminal Error: ", err)
		FatalCleanUp("Cannot Initialize Terminal Error: " + err.Error())
	}

	defer term.Close()

	colog.Register()
	colog.SetOutput(os.Stdout)

	err = readxmlconfig(configFile)

	if err != nil {
		log.Println("alert: Problem opening iotalerter.log file Error: ", err)
		FatalCleanUp("alert: Exiting iotalerter! ...... bye!\n")
	}

	if Config.Global.Settings.Logging == "screen" {
		colog.SetFlags(log.Ldate | log.Ltime)
	}

	if Config.Global.Settings.Logging == "screenwithlineno" || Config.Global.Settings.Logging == "screenandfilewithlineno" {
		colog.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	}

	switch Config.Global.Settings.Loglevel {
	case "trace":
		colog.SetMinLevel(colog.LTrace)
		log.Println("info: Loglevel Set to Trace")
	case "debug":
		colog.SetMinLevel(colog.LDebug)
		log.Println("info: Loglevel Set to Debug")
	case "info":
		colog.SetMinLevel(colog.LInfo)
		log.Println("info: Loglevel Set to Info")
	case "warning":
		colog.SetMinLevel(colog.LWarning)
		log.Println("info: Loglevel Set to Warning")
	case "error":
		colog.SetMinLevel(colog.LError)
		log.Println("info: Loglevel Set to Error")
	case "alert":
		colog.SetMinLevel(colog.LAlert)
		log.Println("info: Loglevel Set to Alert")
	default:
		colog.SetMinLevel(colog.LInfo)
		log.Println("info: Default Loglevel unset in XML config automatically loglevel to Info")
	}

	f, err := os.OpenFile(Config.Global.Settings.Logfilenameandpath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	log.Println("info: Trying to Open File ", Config.Global.Settings.Logfilenameandpath)
	if err != nil {
		FatalCleanUp("Problem Opening iotalerter.log file " + err.Error())
	}

	if Config.Global.Settings.Logging == "screenandfile" {
		log.Println("info: Logging is set to: ", Config.Global.Settings.Logging)
		wrt := io.MultiWriter(os.Stdout, f)
		colog.SetFlags(log.Ldate | log.Ltime)
		colog.SetOutput(wrt)
	}

	if Config.Global.Settings.Logging == "screenandfilewithlineno" {
		log.Println("info: Logging is set to: ", Config.Global.Settings.Logging)
		wrt := io.MultiWriter(os.Stdout, f)
		colog.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
		colog.SetOutput(wrt)
	}

	if Config.Global.Gpio.Enabled {
		log.Println("info: GPIO Enabled ")
		if Config.Global.Gpio.GPIOOffset > 0 {
			for item, pins := range Config.Global.Gpio.Mapping.Map {
				if pins.Enabled {
					newPinNo := Config.Global.Gpio.GPIOOffset + pins.Gpio
					log.Printf("info: Offsetting GPIO PinNo=%v -> %v Name=%v\n", pins.Gpio, newPinNo, pins.Name)
					Config.Global.Gpio.Mapping.Map[item].Gpio = newPinNo
				}
			}
		}
		initGPIO()
		//PCB Specific Settings
		GPIOOutPinByItem("1", "off")
		GPIOOutPinByItem("2", "off")
		GPIOOutPinByItem("3", "off")

	} else {
		log.Println("info: Target Board Set as PC (gpio disabled) ")
	}

	t := time.Now()
	log.Printf("info: iotalerter Application Started at %s\n", t.Format("2006-01-02 15:04:05"))

	go heartBeat()

	iotRelayBoardBanner("\x1b[0;44m")

	if Config.Global.Communication.Mqtt.Enabled {
		mqttsubscribe()
	}

	if Config.Global.Communication.HTTP.Enabled {
		go func() {
			http.HandleFunc("/", httpAPI)
			if err := http.ListenAndServe(":"+Config.Global.Communication.HTTP.Settings.Listenport, nil); err != nil {
				FatalCleanUp("Problem Starting HTTP API Server " + err.Error())
			}
		}()
	}

keyPressListenerLoop:
	for {
		switch ev := term.PollEvent(); ev.Type {
		case term.EventKey:
			switch ev.Key {
			case term.KeyEsc:
				log.Println("info: --")
				log.Println("warn: ESC Key is Invalid")
				reset()
				break keyPressListenerLoop
			case term.KeyDelete:
				log.Println("info: --")
				log.Println("info: Del Key Pressed Menu Requested")
				iotRelayBoardBannerAddColor()
				log.Println("info: --")
			case term.KeyF1:
				log.Println("info: --")
				log.Println("info: F1 Pressed Relay 1 Togle")
				GPIOOutPinByName("relay1", "toggle")
				log.Println("info: --")
			case term.KeyF2:
				log.Println("info: --")
				log.Println("Info: F2 Pressed Relay 2 Toggle")
				GPIOOutPinByName("relay2", "toggle")
				log.Println("info: --")
			case term.KeyF3:
				log.Println("info: --")
				log.Println("Info: F3 Pressed Relay 3 Toggle")
				GPIOOutPinByName("relay3", "toggle")
				log.Println("info: --")
			case term.KeyCtrlC:
				log.Println("info: --")
				log.Println("info: Ctrl-C Terminate Program Requested")
				log.Println("warn: SIGHUP Termination of Program Requested...shutting down...bye!")
				iotalerterAcknowledgements()
				iotRelayBoardQuit()
			case term.KeyCtrlE:
				log.Println("info: --")
				log.Println("info: Ctrl-E Banner Requested")
				iotRelayBoardBannerAddColor()
				log.Println("info: --")
			case term.KeyCtrlG:
				log.Println("info: --")
				log.Println("info: Ctrl-G Pressed")
				log.Println("info: --")
			case term.KeyCtrlL:
				log.Println("info: --")
				log.Println("info: Ctrl-L Clear Screen Requested")
				reset()
				log.Println("info: --")
			case term.KeyCtrlH:
				log.Println("info: --")
				log.Println("info: Ctrl-H Pressed")
				log.Println("info: --")
			case term.KeyCtrlT:
				log.Println("info: --")
				log.Println("info: Ctrl-T Licensing Requested")
				iotRelayBoardLicensing()
				log.Println("info: --")
			case term.KeyCtrlV:
				log.Println("info: --")
				log.Println("info: Ctrl-V Version Requested")
				iotalerterLicensing()
				log.Println("info: --")
			case term.KeyCtrlX:
				log.Println("info: --")
				log.Println("info: Ctrl-X Config Request")
				iotRelayBoardDumpConfig()
				log.Println("info: --")
			case term.KeyCtrlY:
				log.Println("info: --")
				log.Println("info: Ctrl-Y Pressed")
			default:
				log.Println("info: --")
				if ev.Ch == 49 {
					log.Println("Info: 1 Pressed Relay 1 Toggle")
					GPIOOutPinByName("relay1", "toggle")
					log.Println("info: --")
					break
				}
				if ev.Ch == 50 {
					log.Println("Info: 2 Pressed Relay 2 Toggle")
					GPIOOutPinByName("relay2", "toggle")
					log.Println("info: --")
					break
				}
				if ev.Ch == 51 {
					log.Println("Info: 3 Pressed Relay 3 Toggle")
					GPIOOutPinByName("relay3", "toggle")
					log.Println("info: --")
					break
				}
				if ev.Ch == 52 {
					log.Println("Info: 4 Pressed Relay 4 Toggle")
					GPIOOutPinByName("relay4", "toggle")
					log.Println("info: --")
					break
				}
				if ev.Ch == 53 {
					log.Println("Info: 5 Pressed Relay 5 Toggle")
					GPIOOutPinByName("relay5", "toggle")
					log.Println("info: --")
					break
				}
				if ev.Ch == 54 {
					log.Println("Info: 6 Pressed Relay 6 Toggle")
					GPIOOutPinByName("relay6", "toggle")
					log.Println("info: --")
					break
				}
				if ev.Ch == 55 {
					log.Println("Info: 7 Pressed Relay 7 Toggle")
					GPIOOutPinByName("relay7", "toggle")
					log.Println("info: --")
					break
				}
				if ev.Ch == 56 {
					log.Println("Info: 8 Pressed Relay 8 Toggle")
					GPIOOutPinByName("relay8", "toggle")
					log.Println("info: --")
					break
				}
				if ev.Ch == 33 {
					log.Println("Info: Shift-1 Pressed Relay 1 Pulse")
					GPIOOutPinByName("relay1", "pulse")
					log.Println("info: --")
					break
				}
				if ev.Ch == 64 {
					log.Println("Info: Shift-2 Pressed Relay 2 Pulse")
					GPIOOutPinByName("relay2", "pulse")
					log.Println("info: --")
					break
				}
				if ev.Ch == 35 {
					log.Println("Info: Shift-3 Pressed Relay 3 Pulse")
					GPIOOutPinByName("relay3", "pulse")
					log.Println("info: --")
					break
				}
				if ev.Ch == 36 {
					log.Println("Info: Shift-4 Pressed Relay 4 Pulse")
					GPIOOutPinByName("relay4", "pulse")
					log.Println("info: --")
					break
				}
				if ev.Ch == 37 {
					log.Println("Info: Shift-5 Pressed Relay 5 Pulse")
					GPIOOutPinByName("relay5", "pulse")
					log.Println("info: --")
					break
				}
				if ev.Ch == 94 {
					log.Println("Info: Shift-6 Pressed Relay 6 Pulse")
					GPIOOutPinByName("relay6", "pulse")
					log.Println("info: --")
					break
				}
				if ev.Ch == 38 {
					log.Println("Info: Shift-7 Pressed Relay 7 Pulse")
					GPIOOutPinByName("relay7", "pulse")
					log.Println("info: --")
					break
				}
				if ev.Ch == 42 {
					log.Println("Info: Shift-8 Pressed Relay 8 Pulse")
					GPIOOutPinByName("relay8", "pulse")
					log.Println("info: --")
					break
				}

				if ev.Ch != 0 {
					log.Println("warn: Invalid Keypress ASCII", ev.Ch)
				} else {
					log.Println("warn: Key Not Mapped")
				}
				log.Println("info: --")
			}
		case term.EventError:
			log.Println("alert: Terminal Error: ", ev.Err)
			log.Fatal("alert: Exiting iotalerter! ...... bye!\n")
		}

	}

}
