# libvirt-anti-mac-spoof-hook

This hook adds MAC anti-spoofing functionality to libvirt

Docs: 
  - https://www.libvirt.org/hooks.html
  - https://libvirt.org/formatdomain.html#elementsMetadata

Install:
  - copy `qemu` to `/etc/libvirt/hooks/qemu`
  - restart libvirt daemon `systemctl restart libvirtd`

Debug:
```
virsh dumpxml vm1 | DEBUG=true ./qemu vm1 prepare begin -
```

XML inside virsh edit:
```xml
<domain>
...
  <metadata>
	<my:custom xmlns:my="abe43f05-b1b3-4dd2-ad47-b31967e45413">
      <my:network mac_address="52:54:00:9a:c9:01" parent_device="vlan220"/>
    </my:custom>
  </metadata>
...
</domain>
```

XML inside Domain description from libvirt-go-xml:
```xml
<my:custom xmlns:my=\"abe43f05-b1b3-4dd2-ad47-b31967e45413\"><my:network mac_address=\"52:54:00:9a:c9:01\" parent_device=\"vlan220\"/></my:custom>
```

Manual workflow inside hypervisor node:
```
ip link add link vlan220 name macvlan0 type macvlan mode source
ip link set link dev macvlan0 type macvlan macaddr set 52:54:00:34:e5:01
ip -d link show dev macvlan0
```

Chnage MAC inside VM, security check:
```
ip link set dev ens2 address 52:54:00:9a:c9:ff
```

XML inside Libvirt Domain:
```xml
<domain>
  ...
  <devices>
    ...
    <interface type='direct'>
      <mac address='52:54:00:9a:c9:01'/>
      <source dev='macvlan0' mode='bridge'/>
      <model type='virtio'/>
	  </interface>
	...
  </devices>
...
</domain>
```