package clouditor

default applicable = false
default compliant = false

# this is an implementation of metric WebApplicationFirewallEnabled

name := "WebApplicationFirewallEnabled"

waf := input.webApplicationFirewall

applicable {
    waf
}

compliant {
    data.operator == "=="
	waf.enabled == data.target_value
}