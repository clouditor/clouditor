package clouditor.single.sign.on.enabled

import data.clouditor.compare

default applicable = false

default compliant = false

sso := input.singleSignOn

applicable {
	sso
}

compliant {
	compare(data.operator, data.target_value, sso.enabled)
}
