package main

import (
	"fmt"
)

// ConfigMacAntiSpoof configures MAC anti-spoofing according to input config
func ConfigMacAntiSpoof(cfg []MacAntiSpoofConfig) error {

	for _, c := range cfg {
		// create upper macvlan device with mode 'source'
		cmd := RunCommand("ip", "link", "add", "link", c.ParentDevice, "name", c.ChildDevice, "type", "macvlan", "mode", "source")
		if cmd.ReturnCode != 0 && cmd.ReturnCode != 2 { // return code 2 is for RTNETLINK answers: File exists
			return fmt.Errorf("antispoof config error: running command '%s' failed with exit code '%d', output '%s'", cmd.Command, cmd.ReturnCode, cmd.CombinedOutput)
		}

		// set allowed MAC in upper device, this MAC is taken from Libvirt interface config
		cmd = RunCommand("ip", "link", "set", "link", "dev", c.ChildDevice, "type", "macvlan", "macaddr", "set", c.MAC)
		if cmd.ReturnCode != 0 {
			return fmt.Errorf("antispoof config error: running command '%s' failed with exit code '%d', output '%s'", cmd.Command, cmd.ReturnCode, cmd.CombinedOutput)
		}
	}

	return nil
}
