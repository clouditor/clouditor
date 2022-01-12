package clouditor


operator = data.operator
target_value = data.target_value

compare(operator, target_value, actual_value) {
    operator == "=="
    target_value == actual_value
}{
    operator == ">="
    actual_value >= target_value
}{
     operator == "<="
     actual_value <= target_value
}{
     operator == "<"
     actual_value < target_value
}{
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