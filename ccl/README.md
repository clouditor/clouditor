# Cloud Compliance Language (CCL)

It is a simple domain specific language to model rules that should apply to discovered resources. It has the following
syntax:
```ccl
<resourceType> has <expression>
```

The `<resourceType>` refers to an existing resource type in the vocabulary. See package voc.
The `<expression>` must evaluate to a boolean expression, which decide whether the resource is compliant or not

## Expressions

Several different expression exists.

### Comparison

For example, a simple comparison of two values can be achieved using

```ccl
<field> <operatorType> <literalValue>
```

In this case, <field> refers to a field in the resource type, defined the vocabulary. See package voc.
<operatorType> can either be `==`, `!=`, `<=`, `<`, `>`, `>=`, or the special `contains` keyword.
<literalValue> refers to a literal in either a string, integer, float or boolean format.
