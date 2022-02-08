package clouditor.os.logging.enabled

import data.clouditor.compare

default applicable = false

default compliant = false

OSLog := input.oSLog

applicable {
	OSLog
}

compliant {
	compare(data.operator, data.target_value, OSLog.enabled)
}
