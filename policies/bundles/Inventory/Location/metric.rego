package clouditor.metrics.location

import data.clouditor.compare

default applicable = false

default compliant = false

location := input.geoLocation.region

applicable if {
	location != null
}

compliant if {
	compare(data.operator, data.target_value, location)
}