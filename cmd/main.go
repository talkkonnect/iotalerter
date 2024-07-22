package main

import (
	"flag"
	"fmt"
	"github.com/talkkonnect/iotalerter"
)

func main() {

	config := flag.String("config", "/etc/iotalerter.xml", "full path to iotalerter.xml configuration file")

	flag.Usage = iotalerterusage
	flag.Parse()

	iotalerter.Init(*config)
}

func iotalerterusage() {
	fmt.Println("------------------------------------------------------------------------------------------")
	fmt.Println("Usage: iotalerter [-config=[full path and file to iotalerter.xml configuration file]]")
	fmt.Println("By Suvir Kumar <suvir@dits.co.th>")
	fmt.Println("------------------------------------------------------------------------------------------")
	fmt.Println("-config=/home/talkkonnect/gocode/src/github.com/talkkonnect/iotalerter/iotalerter.xml")
	fmt.Println("-version for the version")
	fmt.Println("-help for this screen")
}
