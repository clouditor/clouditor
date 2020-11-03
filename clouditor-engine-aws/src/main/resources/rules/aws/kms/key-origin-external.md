# Key Origin (External)

Checks if the KMS keys have the correct origin (default: 'external'). Master keys are exempted from this check.

```ccl
KeyMetadata has originAsString == "EXTERNAL"
```

## Controls
* BSI C5/CRY-04

[comment]: # TODO merge together all key origin checks and parametrize the rule
