package iotalerter

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

type ConfigStruct struct {
	XMLName xml.Name `xml:"document"`
	Type    string   `xml:"type,attr"`
	Global  struct {
		Device struct {
			Uuid     string `xml:"uuid"`
			Sitename string `xml:"sitename"`
			Zone     string `xml:"zone"`
			Floor    string `xml:"floor"`
			Location string `xml:"location"`
			Type     string `xml:"type"`
			Remark   string `xml:"remark"`
		} `xml:"device"`
		Settings struct {
			Logging                 string `xml:"logging"`
			Loglevel                string `xml:"loglevel"`
			Logfilenameandpath      string `xml:"logfilenameandpath"`
			Outputdevice            string `xml:"outputdevice"`
			Outputdeviceshort       string `xml:"outputdeviceshort"`
			Outputvolcontroldevice  string `xml:"outputvolcontroldevice"`
			Outputmutecontroldevice string `xml:"outputmutecontroldevice"`
		} `xml:"settings"`
		Sounds struct {
			Sound []struct {
				Event    string `xml:"event,attr"`
				File     string `xml:"file,attr"`
				Volume   string `xml:"volume,attr"`
				Blocking bool   `xml:"blocking,attr"`
				Enabled  bool   `xml:"enabled,attr"`
			} `xml:"sound"`
		} `xml:"sounds"`
		Communication struct {
			HTTP struct {
				Enabled  bool `xml:"enabled,attr"`
				Settings struct {
					Listenport     string `xml:"listenport"`
					Httpapirequest struct {
						Enabled bool `xml:"enabled,attr"`
					} `xml:"httpapirequest"`
				} `xml:"settings"`
				Commands struct {
					Command []struct {
						Request string `xml:"request,attr"`
						Enabled bool   `xml:"enabled,attr"`
					} `xml:"command"`
				} `xml:"commands"`
			} `xml:"http"`
			Mqtt struct {
				Enabled  bool `xml:"enabled,attr"`
				Settings struct {
					Mqttsubtopic string `xml:"mqttsubtopic"`
					Mqttpubtopic string `xml:"mqttpubtopic"`
					Mqttbroker   string `xml:"mqttbroker"`
					Mqttpassword string `xml:"mqttpassword"`
					Mqttuser     string `xml:"mqttuser"`
					Mqttid       string `xml:"mqttid"`
					Cleansess    string `xml:"cleansess"`
					Qos          byte   `xml:"qos"`
					Num          string `xml:"num"`
					Payload      string `xml:"payload"`
					Action       string `xml:"action"`
					Store        string `xml:"store"`
					Retained     bool   `xml:"retained"`
				} `xml:"settings"`
				Commands struct {
					Command []struct {
						Request string `xml:"request,attr"`
						Enabled bool   `xml:"enabled,attr"`
					} `xml:"command"`
				} `xml:"commands"`
			} `xml:"mqtt"`
		} `xml:"communication"`
		Payloads struct {
			Payload []struct {
				Name   string `xml:"name,attr"`
				Action []struct {
					Type    string `xml:"type,attr"`
					Item    string `xml:"item,attr"`
					Param   string `xml:"param,attr"`
					Token   string `xml:"token,attr"`
					Method  string `xml:"method,attr"`
					URL     string `xml:"url,attr"`
					ID      string `xml:"id,attr"`
					Title   string `xml:"title,attr"`
					Body    string `xml:"body,attr"`
					Payload string `xml:"payload,attr"`
					Enabled bool   `xml:"enabled,attr"`
				} `xml:"action"`
			} `xml:"payload"`
		} `xml:"payloads"`
		Gpio struct {
			Enabled    bool `xml:"enabled,attr"`
			GPIOOffset uint `xml:"gpiooffset"`
			Heartbeat  struct {
				Enabled         bool          `xml:"enabled,attr"`
				HeartBeatLEDPin uint          `xml:"heartbeatledpin"`
				Periodmsecs     time.Duration `xml:"periodmsecs"`
				Ledonmsecs      time.Duration `xml:"ledonmsecs"`
				Ledoffmsecs     time.Duration `xml:"ledoffmsecs"`
			} `xml:"heartbeat"`
			Mapping struct {
				Map []struct {
					Item               string        `xml:"item,attr"`
					Gpio               uint          `xml:"gpio,attr"`
					Direction          string        `xml:"direction,attr"`
					Devicetype         string        `xml:"devicetype,attr"`
					Name               string        `xml:"name,attr"`
					Pulseleadingmsecs  time.Duration `xml:"pulseleadingmsecs,attr"`
					Pulsemsecs         time.Duration `xml:"pulsemsecs,attr"`
					Pulsetrailingmsecs time.Duration `xml:"pulsetrailingmsecs,attr"`
					Inverted           bool          `xml:"inverted,attr"`
					CurrentState       bool
					Blocking           bool
					Enabled            bool `xml:"enabled,attr"`
				} `xml:"map"`
			} `xml:"mapping"`
		} `xml:"gpio"`
	} `xml:"global"`
}

// Variables used by Program
var (
	KillHeartBeat bool
)

type EventSoundStruct struct {
	Enabled  bool
	FileName string
	Volume   string
	Blocking bool
}

// Generic Global Config Variables
var Config ConfigStruct
var ConfigXMLFile string
var GPIOEnabled bool

func readxmlconfig(file string) error {
	xmlFile, err := os.Open(file)
	if err != nil {
		return errors.New(fmt.Sprintln("cannot open configuration file iotalerter.xml", err))
	}

	log.Println("info: Successfully Opened file iotalerter.xml")
	defer xmlFile.Close()

	byteValue, _ := io.ReadAll(xmlFile)

	err = xml.Unmarshal(byteValue, &Config)
	if err != nil {
		return errors.New(fmt.Sprintln("File iotalerter.xml formatting error Please fix! ", err))
	}

	log.Println("Successfully loaded configuration file into memory")
	return nil
}
