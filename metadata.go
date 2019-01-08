package main

import (
	"encoding/xml"
	"fmt"
	"net"
	"strings"

	"github.com/libvirt/libvirt-go-xml"
)

// DomainMetadata describes custom metadata for Libvirt domain
type DomainMetadata struct {
	XMLName xml.Name `xml:"custom"`
	Text    string   `xml:",chardata"`
	My      string   `xml:"my,attr"`
	Network []struct {
		// Text         string `xml:",chardata"`
		MacAddress   string `xml:"mac_address,attr"`
		ParentDevice string `xml:"parent_device,attr"`
	} `xml:"network"`
}

// GetNetworkConfigFromMetadataXML parses Libvirt domain XML and returns map[string]string (key = MAC address, value = parent device)
func GetNetworkConfigFromMetadataXML(domCfg *libvirtxml.Domain) (map[string]string, error) {

	// check Domain existence
	if domCfg == nil {
		return nil, fmt.Errorf("metadata error: empty Domain XML")
	}

	// declare custom error strings
	customErrors := make([]string, 0)

	// declare output map
	out := make(map[string]string)

	// declare metadata object
	metadata := new(DomainMetadata)

	// check metada existence
	if domCfg.Metadata == nil {
		return nil, fmt.Errorf("metadata error: no metadata inside Domain XML")
	}

	// decode XML to metadata object
	err := xml.Unmarshal([]byte(domCfg.Metadata.XML), metadata)
	if err != nil {
		return nil, fmt.Errorf("metadata error: %s", err.Error())
	}

	// validate metadata XMLNS
	if !strings.EqualFold(MetaDataNameSpace, metadata.XMLName.Space) || !strings.EqualFold(MetaDataNameSpace, metadata.My) {
		return nil, fmt.Errorf("metadata error: XML namespace must be equal to '%s'", MetaDataNameSpace)
	}

	// validate XMLNS local name
	if !strings.EqualFold("custom", metadata.XMLName.Local) {
		return nil, fmt.Errorf("metadata error: XML namespace local part must be equal to 'custom'")
	}

	// get local interfaces list
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("metadata error: %s", err.Error())
	}

	// prepare map of local interfaces, for easier lookups
	mIfaces := make(map[string]bool)

	// populate local interfaces map
	for _, iface := range ifaces {
		mIfaces[iface.Name] = true
	}

	// loop-over network elements
	for _, el := range metadata.Network {

		// check if local interface exists
		if _, ok := mIfaces[el.ParentDevice]; !ok {
			customErrors = append(customErrors, fmt.Sprintf("local interface '%s' does not exist", el.ParentDevice))

			continue
		}

		// validate MAC address
		if !strings.HasPrefix(el.MacAddress, MACAddressQemuPrefix) {
			customErrors = append(customErrors, fmt.Sprintf("MAC '%s' is not valid", el.MacAddress))

			continue
		}

		// update output map
		out[el.MacAddress] = el.ParentDevice
	}

	// return custom errors (for logging)
	if len(customErrors) != 0 {
		return out, fmt.Errorf("metadata error: %s", strings.Join(customErrors, " ,"))
	}

	// no errors, YAY!
	return out, nil
}
