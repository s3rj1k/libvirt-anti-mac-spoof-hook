package main

import (
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

	// get meta config, errors are logged within GetNetworkConfigFromMetadataXML
	mMeta, err := GetNetworkConfigFromMetadataXML(domCfg)
	if err != nil {
		return nil, err
	}

	// get interface config, errors are logged within GetSupportedInterfacesFromDomainXML
	mIfaces, err := GetSupportedInterfacesFromDomainXML(domCfg)
	if err != nil {
		return nil, err
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

	// no errors, YAY!
	return out, nil
}
