package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/libvirt/libvirt-go-xml"
)

// GetDomainXML acquires Libvirt Domain XML
func GetDomainXML(stdin io.Reader) (*libvirtxml.Domain, error) {

	// prepare scanner
	scanner := bufio.NewScanner(stdin)

	// declare variable for parsed XML lines
	lines := make([]string, 0)

	// loop-over scanner, process XML lines
	for scanner.Scan() {
		// clean XML line
		lines = append(lines, strings.TrimSpace(scanner.Text()))
	}

	// fail on scanner error, not EOF
	err := scanner.Err()
	if err != nil {
		e := fmt.Errorf("domain XML error: %s", err.Error())
		Logger.Println(e)

		return nil, e
	}

	// declare Libvirt Domain object
	domCfg := new(libvirtxml.Domain)

	// decode XML to Libvirt Domain object
	err = domCfg.Unmarshal(strings.Join(lines, ""))
	if err != nil {
		e := fmt.Errorf("domain XML error: %s", err.Error())
		Logger.Println(e)

		return nil, e
	}

	return domCfg, nil
}
