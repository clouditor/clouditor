package clouditor

default applicable = false
default compliant = false

# this is an implementation of metric L3FirewallEnabled

name := "L3FirewallEnabled"
metricID := 23

l3f := input.l3Firewall

applicable {
    l3f
}

compliant {
    data.operator == "=="
	l3f.enabled == data.target_value
}