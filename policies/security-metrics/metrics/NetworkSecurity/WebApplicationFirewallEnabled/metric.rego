package metrics.network_security.web_application_firewall_enabled

import data.compare
import rego.v1
import input.accessRestriction.webApplicationFirewall as webApp

default applicable = false

default compliant = false

applicable if {
	webApp
    # the resource type should be an Load Balancer
	input.type[_] == "LoadBalancer"
}

compliant if {
	compare(data.operator, data.target_value, webApp.enabled)
}
