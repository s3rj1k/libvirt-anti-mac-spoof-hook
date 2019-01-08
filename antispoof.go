package main

import (
	"fmt"
	"strings"

	"github.com/libvirt/libvirt-go-xml"
)

// MacAntiSpoofConfig stores data for iproute2 based MAC Anti-Spoof config (using macvlan mode source)
type MacAntiSpoofConfig struct {
	MAC          string
	ParentDevice string
	ChildDevice  string
}

// GetAntiSpoofConfig returns config object for iproute2
func GetAntiSpoofConfig(domCfg *libvirtxml.Domain) ([]MacAntiSpoofConfig, error) {

	// declare output variable
	out := make([]MacAntiSpoofConfig, 0)

	// declare custom error strings
	customErrors := make([]string, 0)

	// get meta config
	mMeta, err := GetNetworkConfigFromMetadataXML(domCfg)
	if err != nil {
		customErrors = append(customErrors, err.Error())
	}

	// get interface config
	mIfaces, err := GetSupportedInterfacesFromDomainXML(domCfg)
	if err != nil {
		customErrors = append(customErrors, err.Error())
	}

	// merge two maps into proper config
	for mac, childDeviceName := range mIfaces {
		// check if MAC is defined in metadata
		if parentDeviceName, ok := mMeta[mac]; ok {

			// iproute2 config object
			var c MacAntiSpoofConfig

			c.MAC = mac
			c.ChildDevice = childDeviceName
			c.ParentDevice = parentDeviceName

			// add to output
			out = append(out, c)
		}
	}

	// return custom errors (for logging)
	if len(customErrors) != 0 {
		return out, fmt.Errorf("%s", strings.Join(customErrors, " ;"))
	}

	// no errors, YAY!
	return out, nil
}
