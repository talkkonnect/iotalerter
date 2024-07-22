package iotalerter

import (
	"log"
)

func iotRelayBoardBanner(backgroundcolor string) {
	var backgroundreset string = "\u001b[0m"
	log.Println("info: " + backgroundcolor + "┌───────────────────────────────────────────────────────────────────────────────────┐" + backgroundreset)
	log.Println("info: " + backgroundcolor + "│██╗ ██████╗ ████████╗     █████╗ ██╗     ███████╗██████╗ ████████╗███████╗██████╗  │" + backgroundreset)
	log.Println("info: " + backgroundcolor + "│██║██╔═══██╗╚══██╔══╝    ██╔══██╗██║     ██╔════╝██╔══██╗╚══██╔══╝██╔════╝██╔══██╗ │" + backgroundreset)
	log.Println("info: " + backgroundcolor + "│██║██║   ██║   ██║       ███████║██║     █████╗  ██████╔╝   ██║   █████╗  ██████╔╝ │" + backgroundreset)
	log.Println("info: " + backgroundcolor + "│██║██║   ██║   ██║       ██╔══██║██║     ██╔══╝  ██╔══██╗   ██║   ██╔══╝  ██╔══██╗ │" + backgroundreset)
	log.Println("info: " + backgroundcolor + "│██║╚██████╔╝   ██║       ██║  ██║███████╗███████╗██║  ██║   ██║   ███████╗██║  ██║ │" + backgroundreset)
	log.Println("info: " + backgroundcolor + "│╚═╝ ╚═════╝    ╚═╝       ╚═╝  ╚═╝╚══════╝╚══════╝╚═╝  ╚═╝   ╚═╝   ╚══════╝╚═╝  ╚═╝ │" + backgroundreset)
	log.Println("info: " + backgroundcolor + "├───────────────────────────────────────────────────────────────────────────────────┤" + backgroundreset)
	log.Println("info: " + backgroundcolor + "│iotalerter A Smart Alerter System for Raspberry Pi                                 │" + backgroundreset)
	log.Println("info: " + backgroundcolor + "├───────────────────────────────────────────────────────────────────────────────────┤" + backgroundreset)
	log.Println("info: " + backgroundcolor + "│Created By : Suvir Kumar <suvir@dits.co.th>                                        │" + backgroundreset)
	log.Println("info: " + backgroundcolor + "│(C) 2024 Suvir Kumar All Rights Reserved by Author                                 │" + backgroundreset)
	log.Println("info: " + backgroundcolor + "├───────────────────────────────────────────────────────────────────────────────────┤" + backgroundreset)
	log.Println("info: " + backgroundcolor + "│Press the <Del> key for Menu or <Ctrl-c> to Quit                                   │" + backgroundreset)
	log.Println("info: " + backgroundcolor + "│Released under Proprietary License of Dynamic IT Solutions Co LTD                  │" + backgroundreset)
	log.Println("info: " + backgroundcolor + "└──────────────────────────────────────────────────────────────── ──────────────────┘" + backgroundreset)
	log.Printf("info: iotalerter Version %v Released %v", version, released)
	log.Printf("info: ")
}

func iotalerterAcknowledgements() {
	log.Println("info: ┌──────────────────────────────────────────────────────────────────────────────────────────────┐")
	log.Println("info: │Acknowledgements & Inspriation from the iotalerter developers, maintainers & testers          │")
	log.Println("info: │iotalerter is based on the works of many people and many projects and libraries               │")
	log.Println("info: ├──────────────────────────────────────────────────────────────────────────────────────────────┤")
	log.Println("info: │Thanks to Organizations :-                                                                    │")
	log.Println("info: │Raspberry Pi Foundation, Developers and Maintainers of Debian                                 │")
	log.Println("info: │The Creators and Maintainers of Golang and all the libraries available on github.com          │")
	log.Println("info: ├──────────────────────────────────────────────────────────────────────────────────────────────┤")
	log.Println("info: │iotalerter was created and Released under Properitary License (C) 2024 DITS <Suvir Kumar>     │")
	log.Println("info: └──────────────────────────────────────────────────────────────────────────────────────────────┘")
}

func iotalerterMenu() {
	log.Println("info: ┌──────────────────────────────────────────────────────────────┐")
	log.Println("info: │             ███╗░░░███╗███████╗███╗░░██╗██╗░░░██╗            │")
	log.Println("info: │             ████╗░████║██╔════╝████╗░██║██║░░░██║            │")
	log.Println("info: │             ██╔████╔██║█████╗░░██╔██╗██║██║░░░██║            │")
	log.Println("info: │             ██║╚██╔╝██║██╔══╝░░██║╚████║██║░░░██║            │")
	log.Println("info: │             ██║░╚═╝░██║███████╗██║░╚███║╚██████╔╝            │")
	log.Println("info: │             ╚═╝░░░░░╚═╝╚══════╝╚═╝░░╚══╝░╚═════╝░            │")
	log.Println("info: ├─────────────────────────────┬────────────────────────────────┤")
	log.Println("info: │ <Del> to Display this Menu  │ Ctrl-C to Quit iotalerter      │")
	log.Println("info: ├─────────────────────────────┼────────────────────────────────┤")
	log.Println("info: │ <F1> Loglevel (Debug/Info)  │ <F2>                           │")
	log.Println("info: ├─────────────────────────────┼────────────────────────────────┤")
	log.Println("info: │ <1>  Relay 1 (On/Off)       │ <Shift-1> Relay 1 (Pulse)      │")
	log.Println("info: │ <2>  Relay 2 (On/Off)       │ <Shift-2> Relay 2 (Pulse)      │")
	log.Println("info: │ <3>  Relay 3 (On/Off)       │ <Shift-3> Relay 3 (Pulse)      │")
	log.Println("info: ├─────────────────────────────┼────────────────────────────────┤")
	log.Println("info: │<Ctrl-E> Show Banner         │ <Ctrl-L> Clear Screen          │")
	log.Println("info: │<Ctrl-V> Display Version     │ <Ctrl-T> View Licensing        │")
	log.Println("info: │<Ctrl-X> Dump Config         │ <Ctrl-Y>                       │")
	log.Println("info: └──────────────────────────────────────────────────────────────┘")
	log.Println("info: IP Address & Session Information")
	localAddresses()
	iotRelayBoardShowVersion()
}

func iotRelayBoardShowVersion() {
	log.Printf("info: iotalerter Version %v Released %v\n", version, released)
}
