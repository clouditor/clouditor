package clouditor.metrics.virtual_machine_disk_encryption_enabled

import input as vm
import input.related

default applicable = false

default compliant = false

applicable {
	vm.type[_] == "VirtualMachine"
}

disks[d] {
	related[_].id == vm.blockStorage[_]

	d := related[_]
}

compliant {
	not {
		disks[_].atRestEncryption.enabled == false
	}
}
