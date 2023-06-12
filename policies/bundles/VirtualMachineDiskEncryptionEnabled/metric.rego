package clouditor.metrics.virtual_machine_disk_encryption_enabled

import future.keywords.every
import input as vm
import input.related

default applicable = false

default compliant = false

applicable {
	# the resource type should be a VM
	vm.type[_] == "VirtualMachine"

	# there should be at least any block storage
	vm.blockStorage[_]
}

disks[d] {
	related[_].id == vm.blockStorage[_]

	d := related[_]
}

compliant {
    every disk in disks {
		disk.atRestEncryption.enabled == true
	}
}
