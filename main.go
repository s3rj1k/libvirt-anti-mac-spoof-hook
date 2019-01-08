package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {

	// GracefullExit logs error to defined logger and exits gracefully
	GracefullExit := func(err error) {

		// log non-empty errors
		if err != nil {
			log.Println(err)
		}

		// exit with 0 code, else libvirt daemon will fail to start VM
		os.Exit(0)
	}

	// validate number of arguments
	if len(os.Args) != 5 {
		GracefullExit(fmt.Errorf("incorrect number of arguments provided, must equal 5, supplied %d", len(os.Args)))
	}

	// switch on: `qemu vm1 {prepare} begin -`
	switch os.Args[2] {
	case "prepare":

		// switch on: `qemu vm1 prepare {begin} -`
		switch os.Args[3] {
		case "begin":

			// get Libvirt Domain XML as object
			domCfg, err := GetDomainXML(os.Stdin)
			if err != nil {
				GracefullExit(err)
			}

			// parse Libvirt Domain XML to get MAC Anti-Spoof Config
			cfg, err := GetAntiSpoofConfig(domCfg)
			if err != nil {
				GracefullExit(err)
			}

			// check debug flag
			if !strings.EqualFold(os.Getenv("DEBUG"), "true") {

				// apply MAC Anti-Spoof Config
				err = ConfigMacAntiSpoof(cfg)
				if err != nil {
					GracefullExit(err)
				}

			} else {
				// do debug printing (virsh dumpxml vm1 | DEBUG=true ./qemu vm1 start begin -)
				for _, el := range cfg {
					fmt.Printf("MAC '%s', Parent device '%s', Child device '%s'\n", el.MAC, el.ParentDevice, el.ChildDevice)
				}
			}

		default:
			GracefullExit(nil)
		}

	default:
		GracefullExit(nil)
	}

	GracefullExit(nil)
}
