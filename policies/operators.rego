package clouditor
import future.keywords.in
import future.keywords.if

# TODO(anatheka): https://play.openpolicyagent.org/p/iNn3sOVG0Q
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

# Checks if the actual_value is in the list of target_values
compare(operator, target_values, actual_value) {
	operator == "isIn"
	actual_value in target_values
}

# Checks if one element of actual_values exists in target_values
compare(operator, target_values, actual_values) {
	operator == "isIn"
    isIn(target_values, actual_values)
}

# TODO(all): Is it necessary, than we have to implement that.
# # Checks if all elements of actual_values exists in target_values
# compare(operator, target_values, actual_values) {
# 	operator == "allIn"
#     isIn(target_values, actual_values)
# }

# Params: target_values (multiple target values), actual_values (multiple actual values)
isIn(target_values, actual_values) {
	# Current implementation: It is enough that one output is one of target_values
	actual_values[_] == target_values[_]
}