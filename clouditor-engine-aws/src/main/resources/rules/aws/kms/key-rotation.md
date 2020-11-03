# Key Rotation

Checks if KMS keys have key rotation enabled. Only applies to non-external keys.

```ccl
KeyMetadata has keyRotationEnabled == true
```

## Controls:

* BSI C5/CRY-04

[comment]: # TODO: actually filter out external and master keys, since they cannot be rotated
[comment]: # filter: 'KeyMetadata has originAsType == "AWS_KMS"'
