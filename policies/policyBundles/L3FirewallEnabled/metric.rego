package clouditor

default applicable = false
default compliant = false

# this is an implementation of metric L3FirewallEnabled

name := "L3FirewallEnabled"

l3f := input.l3Firewall

applicable {
    l3f
}

compliant {
    compare(data.operator, data.target_value, l3f.enabled)
}