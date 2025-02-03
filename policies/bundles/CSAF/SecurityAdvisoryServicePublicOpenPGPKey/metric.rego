package clouditor.metrics.security_advisory_service_public_open_pgp_key

import input as service
import rego.v1

default applicable := false

default compliant := false

key := service.related[service.keyIds[0]]

applicable if {
	# check resource type
	"SecurityAdvisoryService" in service.type
}

compliant if {
	# Public part of PGP key must be available, i.e. 1st) there is a public key und 2nd) it is (publicly) available
	# 1) There is a public key
	key_ids := service.keyIds
	count(key_ids) > 0

	# 2) They are publicly availabe
	every key_id in key_ids {
		key := service.related[key_id]
		key.internetAccessibleEndpoint == true
	}
}
