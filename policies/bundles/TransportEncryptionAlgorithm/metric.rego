package clouditor.metrics.transport_encryption_algorithm

import data.clouditor.compare
import input.transportEncryption as enc

default compliant = false

default applicable = false

algorithm := endpoint.transportEncryption.algorithm

applicable {
<<<<<<< HEAD
	enc
}

compliant {
	compare(data.operator, data.target_value, enc.algorithm)
=======
	algorithm != null
}

compliant {
	compare(data.operator, data.target_value, algorithm)
>>>>>>> main
}
