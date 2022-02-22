package clouditor.metrics.at_rest_encryption_algorithm

import data.clouditor.compare

default applicable = false

default compliant = false

enc := input.atRestEncryption

applicable {
	enc
}

compliant {
	compare(data.operator, data.target_value, enc.algorithm)
}
