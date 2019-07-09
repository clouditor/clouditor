# Enable Default Encryption in Buckets

Buckets should have default encryption using a strong cipher, such as AES-256.

```ccl
Bucket has (not empty applyServerSideEncryptionByDefault.sseAlgorithmAsString) in any bucketEncryption.rules
```

## Controls

* BSI C5/KRY-03
