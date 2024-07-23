package iotalerter

import (
	"log"
	"time"

	"github.com/stianeikeland/go-rpio"
	"github.com/talkkonnect/gpio"
)

type structIO struct {
	Item               string        "xml:\"item,attr\""
	Gpio               uint          "xml:\"gpio,attr\""
	Direction          string        "xml:\"direction,attr\""
	Devicetype         string        "xml:\"devicetype,attr\""
	Name               string        "xml:\"name,attr\""
	Pulseleadingmsecs  time.Duration "xml:\"pulseleadingmsecs,attr\""
	Pulsemsecs         time.Duration "xml:\"pulsemsecs,attr\""
	Pulsetrailingmsecs time.Duration "xml:\"pulsetrailingmsecs,attr\""
	Inverted           bool          "xml:\"inverted,attr\""
	CurrentState       bool
	Blocking           bool
	Enabled            bool "xml:\"enabled,attr\""
}

// Variables for Input Buttons/Switches
var (
	Opto1Used  bool
	Opto1      gpio.Pin
	Opto1Pin   uint
	Opto1State uint

	Opto2Used  bool
	Opto2      gpio.Pin
	Opto2Pin   uint
	Opto2State uint
)

func initGPIO() {

	if !Config.Global.Gpio.Enabled {
		return
	}

	if err := rpio.Open(); err != nil {
		log.Println("error: GPIO Error, ", err)
		GPIOEnabled = false
		return
	}

	//initially turn all relays off on initalization
	for no, io := range Config.Global.Gpio.Mapping.Map {
		if io.Enabled && io.Direction == "output" && io.Devicetype == "relay" {
			gpio.NewOutput(io.Gpio, true)
			Config.Global.Gpio.Mapping.Map[no].CurrentState = true
		}
	}

	// handle inputs on RPI GPIO
	for _, io := range Config.Global.Gpio.Mapping.Map {
		if io.Enabled && io.Direction == "input" && io.Devicetype == "opto" {
			if io.Name == "opto1" && io.Gpio > 0 {
				log.Printf("debug: GPIO Gpio Input Name %v PinNo %v", io.Name, io.Gpio)
				Opto1PinPullUp := rpio.Pin(io.Gpio)
				Opto1PinPullUp.PullUp()
				Opto1Used = true
				Opto1Pin = io.Gpio
			}
		}
		if io.Name == "opto2" && io.Gpio > 0 {
			log.Printf("debug: GPIO Gpio Input Name %v PinNo %v", io.Name, io.Gpio)
			Opto2PinPullUp := rpio.Pin(io.Gpio)
			Opto2PinPullUp.PullUp()
			Opto2Used = true
			Opto2Pin = io.Gpio
		}
	}
}

func GPIOOutPinByName(name string, command string) {
	if !Config.Global.Gpio.Enabled {
		return
	}

	for no, io := range Config.Global.Gpio.Mapping.Map {

		if io.Enabled && io.Direction == "output" && io.Name == name {
			if command == "off" {
				if io.Blocking {
					pinOn(no, io)
				} else {
					go pinOn(no, io)
				}
				break
			}

			if command == "on" {
				if io.Blocking {
					pinOff(no, io)
				} else {
					go pinOff(no, io)
				}
				break
			}

			if command == "toggle" {
				if io.Blocking {
					pinToggle(no, io)
				} else {
					go pinToggle(no, io)
				}
				break
			}

			if command == "pulse" {
				if io.Blocking {
					pinPulse(no, io)
				} else {
					go pinPulse(no, io)
				}
				break
			}
		}
	}

}

func GPIOOutPinByItem(item string, command string) {
	if !Config.Global.Gpio.Enabled {
		return
	}

	for no, io := range Config.Global.Gpio.Mapping.Map {
		if io.Enabled && io.Direction == "output" && io.Item == item {
			if command == "off" {
				if io.Blocking {
					pinOn(no, io)
				} else {
					go pinOn(no, io)
				}
				break
			}

			if command == "on" {
				if io.Blocking {
					pinOff(no, io)
				} else {
					go pinOff(no, io)
				}
				break
			}

			if command == "toggle" {
				if io.Blocking {
					pinToggle(no, io)
				} else {
					go pinToggle(no, io)
				}
				break
			}

			if command == "pulse" {
				if io.Blocking {
					pinPulse(no, io)
				} else {
					go pinPulse(no, io)
				}
				break
			}
		}
	}
}

func GPIOOutPinByGPIO(pin uint, command string) {
	if !Config.Global.Gpio.Enabled {
		return
	}

	var io structIO
	io.Name = "heartbeat"

	if command == "off" {
		if io.Blocking {
			pinOn(int(pin), io)
		} else {
			go pinOn(int(pin), io)
		}
	}

	if command == "on" {
		if io.Blocking {
			pinOff(int(pin), io)
		} else {
			go pinOff(int(pin), io)
		}
	}

}

func GPIOOutAll(name string, command string) {
	if !Config.Global.Gpio.Enabled {
		return
	}

	for _, io := range Config.Global.Gpio.Mapping.Map {
		if io.Enabled && io.Direction == "output" && io.Devicetype == "led/relay" {
			if command == "off" {
				if io.Inverted {
					log.Printf("debug: Turning On %v Output GPIO (Inverted)\n", io.Name)
					gpio.NewOutput(io.Gpio, false)
				} else {
					log.Printf("debug: Turning On %v Output GPIO (Not-Inverted)\n", io.Name)
					gpio.NewOutput(io.Gpio, true)
				}
			}
			if command == "on" {
				if io.Inverted {
					log.Printf("debug: Turning Off %v Output GPIO (Inverted)\n", io.Name)
					gpio.NewOutput(io.Gpio, true)
				} else {
					log.Printf("debug: Turning Off %v Output GPIO (Not-Inverted)\n", io.Name)
					gpio.NewOutput(io.Gpio, false)
				}
			}
		}
	}
}

func heartBeat() {
	if Config.Global.Gpio.Heartbeat.Enabled {
		HeartBeat := time.NewTicker(time.Duration(Config.Global.Gpio.Heartbeat.Periodmsecs) * time.Millisecond)

		for range HeartBeat.C {
			timer1 := time.NewTimer(time.Duration(Config.Global.Gpio.Heartbeat.Ledonmsecs) * time.Millisecond)
			timer2 := time.NewTimer(time.Duration(Config.Global.Gpio.Heartbeat.Ledoffmsecs) * time.Millisecond)
			<-timer1.C
			if Config.Global.Gpio.Heartbeat.Enabled {
				GPIOOutPinByName("heartbeat", "on")
			}
			<-timer2.C
			if Config.Global.Gpio.Heartbeat.Enabled {
				GPIOOutPinByName("heartbeat", "off")
			}
			if KillHeartBeat {
				HeartBeat.Stop()
				return
			}
		}
	}
}

func pinOn(no int, io structIO) {
	if !io.Inverted {
		if io.Name != "heartbeat" {
			log.Printf("debug: Turning On Item %v Name %v at GPIO %v Output GPIO (Non-Inverting)\n", io.Item, io.Name, io.Gpio)
		}
		gpio.NewOutput(io.Gpio, true)
		Config.Global.Gpio.Mapping.Map[no].CurrentState = true
	} else {
		if io.Name != "heartbeat" {
			log.Printf("debug: Turning On Item %v Name %v at GPIO %v Output GPIO (Inverting)\n", io.Item, io.Name, io.Gpio)
		}
		gpio.NewOutput(io.Gpio, false)
		Config.Global.Gpio.Mapping.Map[no].CurrentState = false
	}
}

func pinOff(no int, io structIO) {
	if !io.Inverted {
		if io.Name != "heartbeat" {
			log.Printf("debug: Turning Off Item %v Name %v at GPIO %v Output GPIO (Non-Inverting)\n", io.Item, io.Name, io.Gpio)
		}
		gpio.NewOutput(io.Gpio, false)
		Config.Global.Gpio.Mapping.Map[no].CurrentState = false
	} else {
		if io.Name != "heartbeat" {
			log.Printf("debug: Turning Off Item %v Name %v at GPIO %v Output GPIO (Inverting)\n", io.Item, io.Name, io.Gpio)
		}
		gpio.NewOutput(io.Gpio, true)
		Config.Global.Gpio.Mapping.Map[no].CurrentState = true
	}
}

func pinToggle(no int, io structIO) {
	if Config.Global.Gpio.Mapping.Map[no].CurrentState {
		log.Printf("debug: Turning On Item %v Name %v at GPIO %v Output GPIO (Non-Inverting)\n", io.Item, io.Name, io.Gpio)
	}

	if !Config.Global.Gpio.Mapping.Map[no].CurrentState {
		log.Printf("debug: Turning Off Item %v Name %v at GPIO %v Output GPIO (Non-Inverting)\n", io.Item, io.Name, io.Gpio)
	}

	Config.Global.Gpio.Mapping.Map[no].CurrentState = !Config.Global.Gpio.Mapping.Map[no].CurrentState
	gpio.NewOutput(io.Gpio, Config.Global.Gpio.Mapping.Map[no].CurrentState)

}

func pinPulse(no int, io structIO) {
	log.Printf("debug: Item No %v Pulsing Item %v Name %v at GPIO %v Output GPIO\n", no, io.Item, io.Name, io.Gpio)
	if io.Inverted {
		gpio.NewOutput(io.Gpio, false)
		time.Sleep(io.Pulseleadingmsecs * time.Millisecond)
		gpio.NewOutput(io.Gpio, true)
		time.Sleep(io.Pulsemsecs * time.Millisecond)
		gpio.NewOutput(io.Gpio, false)
		time.Sleep(io.Pulsetrailingmsecs * time.Millisecond)
	}
	if !io.Inverted {
		gpio.NewOutput(io.Gpio, true)
		time.Sleep(io.Pulseleadingmsecs * time.Millisecond)
		gpio.NewOutput(io.Gpio, false)
		time.Sleep(io.Pulsemsecs * time.Millisecond)
		gpio.NewOutput(io.Gpio, true)
		time.Sleep(io.Pulsetrailingmsecs * time.Millisecond)
	}
}
