package clouditor

default applicable = false

default compliant = false

ad := input.anomalyDetection

applicable {
	ad
}

compliant {
	compare(data.operator, data.target_value, ad.enabled)
}
