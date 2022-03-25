package clouditor.metrics.logging_service_encrypted_storage

import future.keywords.every
import input as ls
import input.related

default applicable = false

default compliant = false

applicable {
	ls

	# we also need some kind of storage
	ls.storage[_]
}

disks[s] {
	related[_].id == ls.storage[_]

	s := related[_]
}

compliant {
    every disk in disks {
		disk.atRestEncryption.enabled == true
	}
}
