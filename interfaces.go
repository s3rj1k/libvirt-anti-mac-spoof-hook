package main

import (
	"fmt"
	"strings"

	libvirtxml "github.com/libvirt/libvirt-go-xml"
)

// GetSupportedInterfacesFromDomainXML parses Libvirt and returns supported list of interface
// for which anti-spoof security can be applied, map[string]string (key = MAC address, value = device)
func GetSupportedInterfacesFromDomainXML(domCfg *libvirtxml.Domain) (map[string]string, error) {
	// prefix for interface config errors logging
	const errPrefix = "inteface config error:"

	// declare output variable
	out := make(map[string]string)

	// check Domain existence
	if domCfg == nil {
		e := fmt.Errorf("%s empty Domain XML", errPrefix)
		Logger.Println(e)

		return nil, e
	}

	// check Devices existence
	if domCfg.Devices == nil {
		e := fmt.Errorf("%s Domain '%s' no devices in Domain XML", errPrefix, domCfg.Name)
		Logger.Println(e)

		return nil, e
	}

	// check Interface existence
	if domCfg.Devices.Interfaces == nil {
		e := fmt.Errorf("%s Domain '%s' no interfaces defined in Domain XML", errPrefix, domCfg.Name)
		Logger.Println(e)

		return nil, e
	}

	// loop-over interfaces
	for i := range domCfg.Devices.Interfaces {
		// skip not-supported interface types
		if domCfg.Devices.Interfaces[i].Source.User != nil ||
			domCfg.Devices.Interfaces[i].Source.Ethernet != nil ||
			domCfg.Devices.Interfaces[i].Source.VHostUser != nil ||
			domCfg.Devices.Interfaces[i].Source.Server != nil ||
			domCfg.Devices.Interfaces[i].Source.Client != nil ||
			domCfg.Devices.Interfaces[i].Source.MCast != nil ||
			domCfg.Devices.Interfaces[i].Source.Network != nil ||
			domCfg.Devices.Interfaces[i].Source.Bridge != nil ||
			domCfg.Devices.Interfaces[i].Source.Internal != nil ||
			domCfg.Devices.Interfaces[i].Source.Hostdev != nil ||
			domCfg.Devices.Interfaces[i].Source.UDP != nil {

			Logger.Printf("%s device type with MAC '%s' is not supported\n", errPrefix, domCfg.Devices.Interfaces[i].MAC.Address)

			continue
		}

		// check if 'Direct' device type is defined
		if domCfg.Devices.Interfaces[i].Source.Direct == nil {
			Logger.Printf("%s Domain '%s' device type with MAC '%s' must be of 'Direct'\n", errPrefix, domCfg.Name, domCfg.Devices.Interfaces[i].MAC.Address)

			continue
		}

		// check if 'Direct' device type has proper mode set (bridge or private)
		if !strings.EqualFold(domCfg.Devices.Interfaces[i].Source.Direct.Mode, "bridge") && !strings.EqualFold(domCfg.Devices.Interfaces[i].Source.Direct.Mode, "private") {
			Logger.Printf("%s Domain '%s' device of 'Direct' type with MAC '%s' must have mode set to 'bridge' or 'private'\n", errPrefix, domCfg.Name, domCfg.Devices.Interfaces[i].MAC.Address)

			continue
		}

		// validate MAC address
		if !strings.HasPrefix(domCfg.Devices.Interfaces[i].MAC.Address, MACAddressQemuPrefix) {
			Logger.Printf("%s Domain '%s' MAC '%s' is not valid\n", errPrefix, domCfg.Name, domCfg.Devices.Interfaces[i].MAC.Address)

			continue
		}

		// add to output
		out[domCfg.Devices.Interfaces[i].MAC.Address] = domCfg.Devices.Interfaces[i].Source.Direct.Dev
	}

	// no errors, YAY!
	return out, nil
}
