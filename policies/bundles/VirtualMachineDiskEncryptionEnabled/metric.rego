package clouditor.metrics.virtual_machine_disk_encryption_enabled

import rego.v1

import input as vm

default applicable := false

default compliant := false

applicable if {
	# the resource type should be a VM
	"VirtualMachine" in vm.type

	# there should be at least any block storage
	some _ in vm.blockStorageIds
}

disks[d] if {
	related := vm.related
	related[_].id == vm.blockStorageIds[_]

	d := related[_]
}

compliant if {
	every disk in disks {
		disk.atRestEncryption.enabled == true
	}
}
