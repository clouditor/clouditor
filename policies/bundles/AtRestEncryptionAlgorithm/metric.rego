package clouditor.metrics.at_rest_encryption_algorithm

import data.clouditor.compare

default applicable = false

default compliant = false

algorithm := input.atRestEncryption.algorithm

applicable {
	algorithm
}

compliant {
	compare(data.operator, data.target_value, algorithm)
}
