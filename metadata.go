package main

import (
	"encoding/xml"
	"fmt"
	"net"
	"strings"

	libvirtxml "github.com/libvirt/libvirt-go-xml"
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
	// prefix for metadata errors logging
	const errPrefix = "metadata error:"

	// check Domain existence
	if domCfg == nil {
		e := fmt.Errorf("%s empty Domain XML", errPrefix)
		Logger.Println(e)

		return nil, e
	}

	// declare output map
	out := make(map[string]string)

	// declare metadata object
	metadata := new(DomainMetadata)

	// check metada existence
	if domCfg.Metadata == nil {
		e := fmt.Errorf("%s no metadata inside Domain '%s' XML", errPrefix, domCfg.Name)
		Logger.Println(e)

		return nil, e
	}

	// decode XML to metadata object
	err := xml.Unmarshal([]byte(domCfg.Metadata.XML), metadata)
	if err != nil {
		e := fmt.Errorf("%s Domain '%s' error: %s", errPrefix, domCfg.Name, err.Error())
		Logger.Println(e)

		return nil, e
	}

	// validate metadata XMLNS
	if !strings.EqualFold(MetaDataNameSpace, metadata.XMLName.Space) || !strings.EqualFold(MetaDataNameSpace, metadata.My) {
		e := fmt.Errorf("%s Domain '%s' XML namespace must be equal to %s", errPrefix, domCfg.Name, MetaDataNameSpace)
		Logger.Println(e)

		return nil, e
	}

	// validate XMLNS local name
	if !strings.EqualFold("custom", metadata.XMLName.Local) {
		e := fmt.Errorf("%s Domain '%s' XML namespace local part must be equal to 'custom'", errPrefix, domCfg.Name)
		Logger.Println(e)

		return nil, e
	}

	// get local interfaces list
	ifaces, err := net.Interfaces()
	if err != nil {
		e := fmt.Errorf("%s Domain '%s' error: %s", errPrefix, domCfg.Name, err.Error())
		Logger.Println(e)

		return nil, e
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
			Logger.Printf("%s Domain '%s' local interface '%s' does not exist\n", errPrefix, domCfg.Name, el.ParentDevice)

			continue
		}

		// validate MAC address
		if !strings.HasPrefix(el.MacAddress, MACAddressQemuPrefix) {
			Logger.Printf("%s Domain '%s' MAC '%s' is not valid\n", errPrefix, domCfg.Name, el.MacAddress)

			continue
		}

		// update output map
		out[el.MacAddress] = el.ParentDevice
	}

	// no errors, YAY!
	return out, nil
}
