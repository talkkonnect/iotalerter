package iotalerter

import (
	"errors"
	"log"
	"net"
	"os"
	"reflect"
	"time"

	term "github.com/talkkonnect/termbox-go"
)

func reset() {
	term.Sync()
}

func localAddresses() {
	ifaces, err := net.Interfaces()
	if err != nil {
		log.Printf("error: localAddresses %v\n", err.Error())
		return
	}

	for _, i := range ifaces {
		addrs, err := i.Addrs()

		if err != nil {
			log.Printf("error: localAddresses %v\n", err.Error())
			continue
		}

		for _, a := range addrs {
			if i.Name != "lo" {
				log.Printf("info: %v %v\n", i.Name, a)
			}
		}
	}
}

func FileExists(filepath string) bool {

	fileinfo, err := os.Stat(filepath)

	if os.IsNotExist(err) {
		return false
	}

	return !fileinfo.IsDir()
}

func FatalCleanUp(message string) {
	log.Println("alert: " + message)
	log.Println("alert: iotalerter Terminated Abnormally with the Error(s) As Described.")
	log.Println("info: This Screen will close in 5 seconds")
	time.Sleep(5 * time.Second)
	term.Close()
	os.Exit(1)
}

func Call(m map[string]interface{}, name string, params ...interface{}) (result []reflect.Value, err error) {
	f := reflect.ValueOf(m[name])
	if len(params) != f.Type().NumIn() {
		err = errors.New("the number of params is not adapted")
		return
	}
	in := make([]reflect.Value, len(params))
	for k, param := range params {
		in[k] = reflect.ValueOf(param)
	}
	log.Printf("alert: Calling %v\n", name)
	result = f.Call(in)
	return
}
