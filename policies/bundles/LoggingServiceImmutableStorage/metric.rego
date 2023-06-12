package clouditor.metrics.logging_service_immutable_storage

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

storages[s] {
	related[_].id == ls.storage[_]

	s := related[_]
}

compliant {
    every storage in storages {
		storage.immutability.enabled == true
	}
}
