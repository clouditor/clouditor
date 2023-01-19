package clouditor

# operator and target_value are declared here to add them to the output of each single policy (so assessment can use it)
operator = data.operator

target_value = data.target_value

# we also expose the whole metric configuration as well. In the future we then could get rid of operator and target
# value as individual variables
config = data.config

compare(operator, target_value, actual_value) {
	operator == "=="
	target_value == actual_value
}

compare(operator, target_value, actual_value) {
	operator == ">="
	actual_value >= target_value
}

compare(operator, target_value, actual_value) {
	operator == "<="
	actual_value <= target_value
}

compare(operator, target_value, actual_value) {
	operator == "<"
	actual_value < target_value
}

compare(operator, target_value, actual_value) {
	operator == ">"
	actual_value > target_value
}

# Params: target_values (multiple target values), actual_value (single actual value)
isIn(target_values, actual_value) {
	# Assess actual value with each compliant value in target values
	actual_value == target_values[_]
}

# Params: target_values (multiple target values), actual_values (multiple actual values)
isIn(target_values, actual_values) {
	# Current implementation: It is enough that one output is one of target_values
	actual_values[_] == target_values[_]
}

# // TODO: Add additional with target_values
# Params: actual_values (multiple actual values), target_value (single target value)
has_key(actual_values, target_value) {
	# the _ is needed, otherwise the following returns false: has_key({"foo": false}, "foo")
	 _ = actual_values[target_value]
}