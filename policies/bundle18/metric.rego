package clouditor

default applicable = false
default compliant = false

# this is an implementation of metric MalwareProtectionEnabeld

name := "MalwareProtectionOutput"
metricID := 18

mp := input.malwareProtection

applicable {
    mp
}

compliant {
    data.operator == "=="
	mp.output == data.target_value
}