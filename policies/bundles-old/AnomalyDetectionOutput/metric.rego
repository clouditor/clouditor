package clouditor.anomaly.detection.output

import data.clouditor.isIn

default applicable = false

default compliant = false

ad := input.anomalyDetection

applicable {
	ad
}

compliant {
	isIn(data.target_value, ad.output)
}
