package clouditor.waf.enabled

import data.clouditor.compare

default applicable = false

default compliant = false

waf := input.webApplicationFirewall

applicable {
	waf
}

compliant {
	compare(data.operator, data.target_value, waf.enabled)
}
