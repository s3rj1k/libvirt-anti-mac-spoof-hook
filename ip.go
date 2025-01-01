package main

import (
	"bytes"
	"fmt"
)

// ConfigMacAntiSpoof configures MAC anti-spoofing according to input config
func ConfigMacAntiSpoof(cfg []MacAntiSpoofConfig) error {
	// prefix for antispoof config errors logging
	const errPrefix = "antispoof config error:"

	for _, c := range cfg {
		// create upper macvlan device with mode 'source'
		cmd := RunCommand("ip", "link", "add", "link", c.ParentDevice, "name", c.ChildDevice, "type", "macvlan", "mode", "source")
		if cmd.ReturnCode != 0 && cmd.ReturnCode != 2 { // return code 2 is for RTNETLINK answers: File exists
			e := fmt.Errorf("%s running command '%s' failed with exit code '%d', output '%s'", errPrefix, cmd.Command, cmd.ReturnCode, cmd.CombinedOutput)
			Logger.Println(e)

			return e
		}

		// set allowed MAC in upper device, this MAC is taken from Libvirt interface config
		cmd = RunCommand("ip", "link", "set", "link", "dev", c.ChildDevice, "type", "macvlan", "macaddr", "set", c.MAC)
		if cmd.ReturnCode != 0 {
			e := fmt.Errorf("%s running command '%s' failed with exit code '%d', output '%s'", errPrefix, cmd.Command, cmd.ReturnCode, cmd.CombinedOutput)
			Logger.Println(e)

			return e
		}
	}

	return nil
}

// UnConfigMacAntiSpoof removes upper macvlan (in mode source) interface from system
func UnConfigMacAntiSpoof(cfg []MacAntiSpoofConfig) error {
	// prefix for antispoof config errors logging
	const errPrefix = "antispoof config error:"

	for _, c := range cfg {
		// get extended information for defined parent interface
		cmd := RunCommand("ip", "-o", "-d", "l", "show", c.ChildDevice, "type", "macvlan")
		if cmd.ReturnCode != 0 {
			e := fmt.Errorf("%s running command '%s' failed with exit code '%d', output '%s'", errPrefix, cmd.Command, cmd.ReturnCode, cmd.CombinedOutput)
			Logger.Println(e)

			return e
		}

		// skip interfaces not in mode 'source'
		if !bytes.Contains(cmd.CombinedOutput, []byte("macvlan mode source")) {
			continue
		}

		// remove parent interface
		cmd = RunCommand("ip", "l", "del", c.ChildDevice, "type", "macvlan")
		if cmd.ReturnCode != 0 {
			e := fmt.Errorf("%s running command '%s' failed with exit code '%d', output '%s'", errPrefix, cmd.Command, cmd.ReturnCode, cmd.CombinedOutput)
			Logger.Println(e)

			return e
		}
	}

	return nil
}
