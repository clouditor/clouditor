package clouditor

default compliant = false

# this is an implementation of metric AnomalyDetectionEnabled

compliant {
	ad := input.anomalyDetection
	tls.enabled == true
}