package metrics.application_security.strong_cryptographic_hash

import data.clouditor.compare
import rego.v1
import input as app

default applicable = false

default compliant = false

hashes := [func | func := app.functionalities[_]; func.cryptographicHash]

applicable if {
	#some i
	#functionalities[i].cryptographicHash

	# the resource type should be an application
	app.type[_] == "Application"
}

compliant if {
	count(violations) == 0
}

message := "The anaylzed resource uses strong cryptographic hashes." if {
	compliant
} else := "The anaylzed resource contains evidence that weak cryptographic hashes are used." if {
	not compliant
}

results := [
mapped |
	func := app.functionalities[_]
	mapped := {
		"property": "cryptographicHash.algorithm",
		"value": func.cryptographicHash.algorithm,
		"target_value": data.target_value,
		"operator": data.operator,
		"success": compare(data.operator, data.target_value, func.cryptographicHash.algorithm),
	}
]

violations := [x | y := results[_]; y.success == false; x = y]