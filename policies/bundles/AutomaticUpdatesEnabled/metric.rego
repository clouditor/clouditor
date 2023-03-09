# METADATA
# title: MFA_Enabled
package clouditor.metrics.automatic_updates_enabled

import data.clouditor.compare
import input.automaticUpdates as am

default applicable = false

default compliant = false

# METADATA
# title: MFA_Enabled
applicable {
	am
}

# METADATA
# title: MFA_Enabled
compliant {
    annotations := rego.metadata.title()
    decision := annotations.title
	compare(data.operator, data.target_value, am.enabled)
}
