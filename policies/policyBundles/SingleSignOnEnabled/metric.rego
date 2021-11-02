package clouditor

default applicable = false
default compliant = false

# this is an implementation of metric SingleSignOnEnabled

name := "SingleSignOnEnabled"

sso := input.singleSignOn

applicable {
    sso
}

compliant {
    data.operator == "=="
	sso.enabled == data.target_value
}