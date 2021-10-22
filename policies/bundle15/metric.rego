package clouditor

default applicable = false
default compliant = false

# this is an implementation of metric AnomalyDetectionEnabled

name := "AnomalyDetectionEnabled"
metricID := 15

ad := input.anomalyDetection

applicable {
    ad
}

compliant {
    data.operator == "=="
	ad.enabled == data.target_value
}