package clouditor
import future.keywords.every # includes also in

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

# Checks if the actual_value (string) exists in target_values (array)
compare(operator, target_values, actual_value) {
	operator == "isIn"
    # Check if the input value actual_value is a string, otherwise the compare function for array must be used
    is_string(actual_value)
	actual_value in target_values
}

# Checks if one element of actual_values (array) exists in target_values (array)
compare(operator, target_values, actual_values) {
	operator == "isIn"
    is_array(actual_values)
    some act_val in actual_values 
    act_val in target_values
}

# Checks if the actual_value (string) exists in target_values (array)
compare(operator, target_values, actual_value) {
	operator == "allIn"
	# Check if the input value actual_value is a string, otherwise the compare function for array must be used
    is_string(actual_value)
    actual_value in target_values
}

# Checks if all elements of actual_values (array) exists in target_values (array)
compare(operator, target_values, actual_values) {
	operator == "allIn"
    is_array(actual_values)
    every act_val in actual_values {
    	act_val in target_values
    }
}