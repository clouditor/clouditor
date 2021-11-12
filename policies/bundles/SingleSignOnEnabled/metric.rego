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
    compare(data.operator, data.target_value, sso.enabled)
}