package clouditor.metrics.anomaly_detection_output

import data.clouditor.compare

default applicable = false

default compliant = false

output := input.anomalyDetection.applicationLogging.loggingService

applicable {
	output != null
}

compliant {
	compare(data.operator, data.target_value, output)
}
