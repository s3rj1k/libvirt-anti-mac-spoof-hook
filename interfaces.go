package main

import (
	"fmt"
	"strings"

	"github.com/libvirt/libvirt-go-xml"
)

// GetSupportedInterfacesFromDomainXML parses Libvirt and returns supported list of interface
// for which anti-spoof security can be applied, map[string]string (key = MAC address, value = device)
func GetSupportedInterfacesFromDomainXML(domCfg *libvirtxml.Domain) (map[string]string, error) {

	// declare custom error strings
	customErrors := make([]string, 0)

	// declare output variable
	out := make(map[string]string)

	// check Domain existence
	if domCfg == nil {
		return nil, fmt.Errorf("interface config error: empty Domain XML")
	}

	// check Devices existence
	if domCfg.Devices == nil {
		return nil, fmt.Errorf("interface config error: no devices in Domain XML")
	}

	// check Interface existence
	if domCfg.Devices.Interfaces == nil {
		return nil, fmt.Errorf("interface config error: no interfaces defined in Domain XML")
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

			customErrors = append(customErrors, fmt.Sprintf("device type with MAC '%s' is not supported", domCfg.Devices.Interfaces[i].MAC.Address))

			continue
		}

		// validate MAC address
		if !strings.HasPrefix(domCfg.Devices.Interfaces[i].MAC.Address, MACAddressQemuPrefix) {
			customErrors = append(customErrors, fmt.Sprintf("MAC '%s' is not valid", domCfg.Devices.Interfaces[i].MAC.Address))

			continue
		}

		// check if 'Direct' device type is defined
		if domCfg.Devices.Interfaces[i].Source.Direct == nil {
			customErrors = append(customErrors, fmt.Sprintf("device type with MAC '%s' must be of 'Direct'", domCfg.Devices.Interfaces[i].MAC.Address))

			continue
		}

		// check if 'Direct' device type has proper mode set (bridge or private)
		if !strings.EqualFold(domCfg.Devices.Interfaces[i].Source.Direct.Mode, "bridge") && !strings.EqualFold(domCfg.Devices.Interfaces[i].Source.Direct.Mode, "private") {
			customErrors = append(customErrors, fmt.Sprintf("device of 'Direct' type with MAC '%s' must have mode set to 'bridge' or 'private'", domCfg.Devices.Interfaces[i].MAC.Address))

			continue
		}

		// add to output
		out[domCfg.Devices.Interfaces[i].MAC.Address] = domCfg.Devices.Interfaces[i].Source.Direct.Dev
	}

	// return custom errors (for logging)
	if len(customErrors) != 0 {
		return out, fmt.Errorf("interface config error: %s", strings.Join(customErrors, " ,"))
	}

	// no errors, YAY!
	return out, nil
}
