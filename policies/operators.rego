package clouditor
import future.keywords.every # includes also in

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

# Checks if the actual_values (array) contains the target_value (string)
compare(operator, target_value, actual_values) {
	operator == "isIn"
    # Check if the input value actual_value is a string, otherwise the compare function for array must be used
    is_string(target_value)
    is_array(actual_values)
	target_value in actual_values
}

# Checks if one element of actual_values (array) exists in target_values (array)
compare(operator, target_values, actual_values) {
	operator == "isIn"
    is_array(actual_values)
    some act_val in actual_values 
    act_val in target_values
}

# Checks if one element of target_values (array) exists in key of actual_values (object)
compare(operator, target_values, actual_values) {
	operator == "isIn"
    is_object(actual_values)
    # Get all keys from objects
    value := object.keys(actual_values)
    # Check if one the keys is in array of target_values
    some v in value
    v in target_values
}

# Checks if the target_value (string) exists in key of actual_values (object)
compare(operator, target_value, actual_values) {
	operator == "isIn"
    is_object(actual_values)
    # Get all keys from objects
    value := object.keys(actual_values)
    # Check if target_value exists in the set of object's keys 
    target_value in value
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