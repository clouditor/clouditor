package clouditor

default applicable = false
default compliant = false

# this is an implementation of metric AnomalyDetectionOutput

name := "AnomalyDetectionOuput"

ad := input.anomalyDetection

applicable {
    ad
}

compliant {
    data.operator == "=="
    ad.output == data.target_value
}