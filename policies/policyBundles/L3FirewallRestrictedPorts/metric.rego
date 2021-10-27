package clouditor

default applicable = false
default compliant = false

# this is an implementation of metric L3FirewallRestrictedPorts

name := "L3FirewallRestrictedPorts"

l3f := input.l3Firewall

applicable {
    l3f
}

compliant {
    data.operator == "=="
	l3f.restrictedPorts == data.target_value
}