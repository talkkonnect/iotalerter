package iotalerter

import (
	"crypto/tls"
	"log"
	"strings"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

var MQTTPublishPayload MQTT.Token

var MQTTClient MQTT.Client

func mqttsubscribe() {
	if Config.Global.Communication.Mqtt.Enabled {
		log.Printf("info: MQTT Subscription Information")
		log.Printf("info: MQTT Broker      : %s\n", Config.Global.Communication.Mqtt.Settings.Mqttbroker)
		log.Printf("debug: MQTT clientid    : %s\n", Config.Global.Communication.Mqtt.Settings.Mqttid)
		log.Printf("debug: MQTT user        : %s\n", Config.Global.Communication.Mqtt.Settings.Mqttuser)
		log.Printf("debug: MQTT password    : %s\n", Config.Global.Communication.Mqtt.Settings.Mqttpassword)
		log.Printf("info: Subscribed topic : %s\n", Config.Global.Communication.Mqtt.Settings.Mqttsubtopic)

		connOpts := MQTT.NewClientOptions().AddBroker(Config.Global.Communication.Mqtt.Settings.Mqttbroker).SetClientID(Config.Global.Communication.Mqtt.Settings.Mqttid).SetCleanSession(true)
		if Config.Global.Communication.Mqtt.Settings.Mqttuser != "" {
			connOpts.SetUsername(Config.Global.Communication.Mqtt.Settings.Mqttuser)
			if Config.Global.Communication.Mqtt.Settings.Mqttpassword != "" {
				connOpts.SetPassword(Config.Global.Communication.Mqtt.Settings.Mqttpassword)
			}
		}
		tlsConfig := &tls.Config{InsecureSkipVerify: true, ClientAuth: tls.NoClientCert}
		connOpts.SetTLSConfig(tlsConfig)

		connOpts.OnConnect = func(c MQTT.Client) {
			if token := c.Subscribe(Config.Global.Communication.Mqtt.Settings.Mqttsubtopic, byte(Config.Global.Communication.Mqtt.Settings.Qos), onMessageReceived); token.Wait() && token.Error() != nil {
				log.Println("error: MQTT Token Error!")
				return
			}
		}

		MQTTClient = MQTT.NewClient(connOpts)
		if token := MQTTClient.Connect(); token.Wait() && token.Error() != nil {
			log.Println("error: MQTT Token Error!")
			return
		} else {
			log.Printf("info: Connected to     : %s\n", Config.Global.Communication.Mqtt.Settings.Mqttbroker)
		}
	}
}

func MQTTPublish(mqttPayload string) {
	MQTTPublishPayload = MQTTClient.Publish(Config.Global.Communication.Mqtt.Settings.Mqttpubtopic, Config.Global.Communication.Mqtt.Settings.Qos, Config.Global.Communication.Mqtt.Settings.Retained, mqttPayload)
	go func() {
		<-MQTTPublishPayload.Done()
		if MQTTPublishPayload.Error() != nil {
			log.Println("error: ", MQTTPublishPayload.Error())
		} else {
			log.Printf("info: Successfully Published MQTT Topic %v Payload %v\n", Config.Global.Communication.Mqtt.Settings.Mqttpubtopic, mqttPayload)
			return
		}
	}()
}

func onMessageReceived(client MQTT.Client, message MQTT.Message) {
	var (
		CommandDefined bool
		PayLoad        string
	)

	funcs := map[string]interface{}{
		"showmenu":    iotalerterMenu,
		"showbanner":  iotRelayBoardBannerAddColor,
		"showversion": iotRelayBoardShowVersion,
		"showconfig":  iotRelayBoardDumpConfig,
		"listmqtt":    listMqttAPI}

	PayLoad = strings.ToLower(string(message.Payload()))
	log.Printf("info: Received MQTT message on topic: %s Payload: %s\n", message.Topic(), PayLoad)

	byteCommand := strings.Split(strings.ToLower(PayLoad), ":")
	stringCommand := strings.Join(byteCommand[:], "")
	Command := strings.Split(stringCommand, " ")

	//handle static commands
	for _, mqttcommand := range Config.Global.Communication.Mqtt.Commands.Command {
		if strings.TrimSuffix(PayLoad, "\n") == strings.ToLower(mqttcommand.Request) && mqttcommand.Enabled {
			CommandDefined = true
			break
		}
	}

	if CommandDefined {
		for _, mqttcommand := range Config.Global.Communication.Mqtt.Commands.Command {
			var err error
			if strings.Contains(Command[0], mqttcommand.Request) {
				if len(Command) == 1 {
					_, err = Call(funcs, mqttcommand.Request)
					if err == nil {
						log.Printf("MQTT Command %v Processed", strings.TrimSuffix(PayLoad, "\n"))
					} else {
						log.Printf("error: MQTT Command %v Failed", strings.TrimSuffix(PayLoad, "\n"))
					}
				}
			}
		}
		return
	}

	//handle dynamic commands
	for _, payloads := range Config.Global.Payloads.Payload {
		for _, payload := range payloads.Action {
			if payloads.Name == strings.TrimSuffix(PayLoad, "\n") {
				CommandDefined = true
				if payload.Type == "gpio" && payload.Enabled {
					log.Printf("MQTT Command GPIO item=%v param=%v", payload.Item, payload.Param)
					GPIOOutPinByItem(payload.Item, payload.Param)
				}
				if payload.Type == "http" && payload.Enabled {
					if payload.Method == "get" {
						HTTPAPIRequestGet(payload.Token, payload.URL)
					}
					if payload.Method == "post" {
						HTTPAPIRequestPost(payload.Token, payload.URL)
					}
				}
				if payload.Type == "mqtt" && payload.Enabled {
					log.Printf("payload=%v", payload.Payload)
				}
			}
		}
	}

	if !CommandDefined {
		log.Printf("error: MQTT Command %v Not Defined or Not Enabled\n", strings.TrimSuffix(PayLoad, "\n"))
		return
	}

}

func listMqttAPI() {
	for _, standardMqttCommand := range Config.Global.Communication.Mqtt.Commands.Command {
		log.Printf("info: MQTT Standard Commands available %v \n", standardMqttCommand.Request)
	}
	for _, udpayloads := range Config.Global.Payloads.Payload {
		for _, udpayload := range udpayloads.Action {
			if udpayload.Enabled {
				log.Printf("info: MQTT User Defined Payloads Commands available %v %v \n", udpayloads.Name, udpayload)
			}
		}
	}
}
