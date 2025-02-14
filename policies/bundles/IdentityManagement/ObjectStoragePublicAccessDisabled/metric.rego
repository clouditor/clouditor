package clouditor.metrics.object_storage_public_access_disabled

import data.clouditor.compare
import input as storage

default compliant = false

default applicable = false

applicable if {
	# the resource type should be an ObjectStorage
	storage.type[_] == "ObjectStorage"
}

compliant if {
	compare(data.operator, data.target_value, storage.publicAccess)
}
