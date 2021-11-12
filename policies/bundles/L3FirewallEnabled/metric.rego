package clouditor

default applicable = false

default compliant = false

l3f := input.l3Firewall

applicable {
	l3f
}

compliant {
	compare(data.operator, data.target_value, l3f.enabled)
}
