package clouditor

default applicable = false

default compliant = false

sso := input.singleSignOn

applicable {
	sso
}

compliant {
	compare(data.operator, data.target_value, sso.enabled)
}
