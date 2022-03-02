package clouditor.metrics.mutual_authentication

import data.clouditor.compare

default applicable = false

default compliant = false

cba := input.certificateBasedAuthentication

enc := input.httpEndpoint.transportEncryption

applicable {
	cba
	enc
}

# TODO(all): Actually, in this case, data.operator and data.target_value are for the overall metric. Not single checks.
# TODO(cont.): That would mean it is reather compliant = data.targetValue
# TODO(lebogg): Look if we can access other evaluated metrics within this policy, e.g. TransportEncryptionEnabled
compliant {
	compare(data.operator, data.target_value, cba.enabled)
	compare(data.operator, data.target_value, enc.enabled)
}
