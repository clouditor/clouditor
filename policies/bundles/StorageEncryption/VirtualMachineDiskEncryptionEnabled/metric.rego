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

# If all blockStorageIds are available in the related object array, than we store all related disks in d
# TODO(all): We should only store the related disks in d, that are available in blockStorageIds
disks[d] if {
	related := vm.related
	every x in vm.blockStorageIds {
    	x in related[_]
    }
	d := related[_]
}

compliant if {
	is_object(disks)

	# Check if the list of disks is not empty
	# It is possible, that blockStorageIds exist, but they are not in the related disk list available. In that case the disks are empty and the evaluation ends with compliant = true.
	count(disks) > 0
    
	# The list of disks is a key/value pair with the ontology information as key and a boolean as value. As we want to check the ontology information, we have to use the key instead of the value for the check.
	every disk, _ in disks {
		disk.atRestEncryption.customerKeyEncryption.enabled == true
	}
}
