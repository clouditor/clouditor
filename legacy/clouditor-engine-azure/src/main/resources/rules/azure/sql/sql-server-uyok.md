# Use Your Own Key

Checks, whether certain Azure SQL servers have key vault keys stored in their encryption profile.

```ccl
SQLServer has encryptionProtectors.serverKeyType == "AzureKeyVault"
```

## Controls
* BSI C5/CRY-04