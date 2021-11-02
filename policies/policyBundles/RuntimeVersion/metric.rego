package clouditor

default applicable = false
default compliant = false

# this is an implementation of metric RuntimeVersion

name := "RuntimeVersion"

runtimeLanguage := input.runtimeLanguage
runtimeVersion := input.runtimeVersion

applicable {
    runtimeLanguage
    runtimeVersion
}

# TODO(all): Consider to put `operator` into list of target_values for more granularity
compliant {
    data.operator == "=="
	some i
	runtimeLanguage == data.target_value[i].runtimeLanguage
	runtimeVersion == data.target_value[i].runtimeVersion
}{
    data.operator == ">="
	some i
	runtimeLanguage == data.target_value[i].runtimeLanguage
	runtimeVersion >= data.target_value[i].runtimeVersion
}
