package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {

	// closing logfile
	defer func(fd *os.File) {
		err := fd.Close()
		if err != nil {
			log.Fatalf("error closing log file: %v", err)
		}
	}(Fd)

	// GracefullExit logs error to defined logger and exits gracefully
	GracefullExit := func(err error) {

		// log hook exit
		if err != nil {
			Logger.Println("graceful exit for libvirt, but error occurred")
		} else {
			Logger.Println("graceful exit for libvirt, no errors occurred")
		}

		// exit with 0 code, else libvirt daemon will fail to start VM
		os.Exit(0)
	}

	switch os.Args[2] {

	// switch on: `qemu vm1 {prepare} begin -`
	case "prepare":

		switch os.Args[3] {

		// switch on: `qemu vm1 prepare {begin} -`
		case "begin":

			Logger.Println("hook: started, begin")

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
				// do debug printing (virsh dumpxml vm1 | DEBUG=true ./qemu vm1 prepare begin -)
				for _, el := range cfg {
					fmt.Printf("MAC '%s', Parent device '%s', Child device '%s'\n", el.MAC, el.ParentDevice, el.ChildDevice)
				}
			}

		default:
			GracefullExit(nil)
		}

	// switch on: `qemu vm1 {stopped} end -`
	case "stopped":

		switch os.Args[3] {

		// switch on: `qemu vm1 stopped {end} -`
		case "end":

			Logger.Println("hook: stopped, end")

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
				// de-apply MAC Anti-Spoof Config
				err = UnConfigMacAntiSpoof(cfg)
				if err != nil {
					GracefullExit(err)
				}
			} else {
				// do debug printing (virsh dumpxml vm1 | DEBUG=true ./qemu vm1 stopped end -)
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
