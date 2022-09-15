package clouditor.metrics.anomaly_detection_output

import data.clouditor.isIn

default applicable = false

default compliant = false

output := input.anomalyDetection.applicationLogging.loggingService

applicable {
	output != null
}

compliant {
	isIn(data.target_value, output)
}
