package iotalerter

import (
	"log"
	"os"
	"os/exec"
)

func iotRelayBoardQuit() {
	c := exec.Command("reset")
	c.Stdout = os.Stdout
	c.Run()
	os.Exit(0)
}

func iotRelayBoardLicensing() {
	log.Println("info: The software, documentation and any hardware accompanying this License whether on the hardware, in read only memory, on any other media")
	log.Println("info: or in any other form (collectively the 'Software') are licensed, not sold, to you by Dynamic IT Solutions Co,. Ltd.('DITS')")
	log.Println("info: for use only under the terms of this License, and DITS reserves all rights not expressly granted to you. The rights granted")
	log.Println("info: herein are limited to Global’s intellectual property rights and do not include any other patents or intellectual")
	log.Println("info: property rights. You own the hardware on which the Application is running but Global licensor(s) retain ownership of the Software itself.")
}

func iotalerterLicensing() {
	log.Printf("info: iotalerter Version %v Released %v\n", version, released)
}

func iotRelayBoardDumpConfig() {
	var colorSet string = "\033[31m"
	var colorReset string = "\u001b[0m"
	log.Printf("info: ---------- Dump Settings Configuration ------- \n")
	log.Printf("info: "+colorSet+"Version            "+colorReset+"%v\n", version)
	log.Printf("info: "+colorSet+"Released           "+colorReset+"%v\n", released)
	log.Println("info: ---------------------------------------------- ")
}

func iotRelayBoardBannerAddColor() {
	iotRelayBoardBanner("\x1b[0;44m")
}
