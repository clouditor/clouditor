package clouditor

default applicable = false
default compliant = false

# this is an implementation of metric BootLoggingEnabled

name := "BootLoggingEnabled"

bootLog := input.bootLog

applicable {
    bootLog[_]
}

compliant {
    data.operator == "=="
	bootLog.enabled == data.target_value
}