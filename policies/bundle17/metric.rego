package clouditor

default applicable = false
default compliant = false

# this is an implementation of metric MalwareProtectionEnabeld

name := "MalwareProtectionEnabeld"
metricID := 16

mp := input.malwareProtection

applicable {
    mp
}

compliant {
    data.operator == "=="
	mp.enabled == data.target_value
}