package clouditor.metrics.boot_logging_secure_transport

import future.keywords.every
import data.clouditor.isIn
import input.bootLogging as logging
import input.related

default applicable = false

default compliant = false

applicable {
	logging

	# we also need some kind of output to a logging service
	logging.loggingService[_]
}

services[s] {
	related[_].id == logging.loggingService[_]

	s := related[_]
}

compliant {
	# TODO(oxisto): It would be super cool to just depend on the transport_encryption_enabled metric for this
	# TODO(oxisto): Also it would be cool to depend on the logging_service_*_storage metrics
    every service in services {
		service.transportEncryption.enabled == true
	}
}
