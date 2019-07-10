# Key Origin (KMS)

Checks if the KMS keys have the correct origin (default: 'kms'). Master keys are exempted from this check.

```ccl
KeyMetadata has originAsString == "AWS_KMS"
```

[comment]: # TODO merge together all key origin checks and parametrize the rule
