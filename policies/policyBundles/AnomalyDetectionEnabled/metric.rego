package clouditor

default applicable = false
default compliant = false

# this is an implementation of metric AnomalyDetectionEnabled

name := "AnomalyDetectionEnabled"

ad := input.anomalyDetection

applicable {
    ad
}

compliant {
    compare(data.operator, data.target_value, ad.enabled)
}