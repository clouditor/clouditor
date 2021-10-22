package clouditor

default applicable = false
default compliant = false

# this is an implementation of metric OSLoggingRetention

name := "OSLoggingRetention"
metricID := 13

ad := input.anomalyDetection

applicable {
    ad
}

compliant {
    data.operator == "=="
	ad.enabled == data.target_value
}