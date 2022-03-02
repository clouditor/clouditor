package clouditor.metrics.anomaly_detection_enabled

import data.clouditor.compare

default applicable = false

default compliant = false

ad := input.anomalyDetection

applicable {
	ad
}

compliant {
	compare(data.operator, data.target_value, ad.enabled)
}
