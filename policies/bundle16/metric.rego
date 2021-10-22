package clouditor

default applicable = false
default compliant = false

# this is an implementation of metric AnomalyDetectionOutput

name := "AnomalyDetectionOutput"
metricID := 16

ad := input.anomalyDetection

applicable {
    ad
}

compliant {
    data.operator == "=="
	ad.output == data.target_value
}