package clouditor.virtual.machine.disk.encryption.enabled

import input.related
import input.virtualMachine as vm

default applicable = false

default compliant = false

applicable {
	vm
}

disks[d] {
	related[_].id == vm.blockStorage[_]

	d := related[_]
}

compliant {
	disks[_].atRestEncryption.enabled == true
}
