package clouditor

default applicable = false
default compliant = false

# this is an implementation of metric LoggingEnabled

name := "LoggingEnabled"

log := input.log

applicable {
    log
}

compliant {
    data.operator == "=="
    log.activated == data.target_value
}