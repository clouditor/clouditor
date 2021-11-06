package clouditor

default applicable = false
default compliant = false

# this is an implementation of metric L3FirewallRestrictedPorts

name := "L3FirewallRestrictedPorts"

l3f := input.l3Firewall

applicable {
    l3f
}

# TODO(all): Maybe change restrictet ports to array of strings. See comment in Ontology.
compliant {
    compare(data.operator, data.target_value, l3f.restrictedPorts)
}