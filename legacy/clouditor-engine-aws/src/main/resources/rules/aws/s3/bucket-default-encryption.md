# Enable Default Encryption in Buckets

Buckets should have default encryption using a strong cipher, such as AES-256.

```ccl
Bucket has (not empty applyServerSideEncryptionByDefault.sseAlgorithm) in any bucketEncryption.rules
```

## Controls

* BSI C5/CRY-03
