package clouditor

default applicable = false
default compliant = false

# this is an implementation of metric AnomalyDetectionOutput

name := "AnomalyDetectionOuput"
metricID := 12

ad := input.anomalyDetection

applicable {
    ad
}

compliant {
    # ToDo(lebogg): Check if 'input.osLog' is generated (in JSON) or, e.g., 'input.OSLog'

    data.operator == "=="
    ad.output == data.target_value
}