# Set Expiry Time for Keys

Checks if all keys in Azure Key Vault have an expiry time set and is less than a year.

```ccl
KeyVault has attributes.exp before 365 days in all keys
```

## Controls
  
* CIS Microsoft Azure Foundations Benchmark/Azure 8.1
* BSI C5/CRY-04