# Set Expiry Time for Secrets

Checks if all secrets in Azure Key Vault have an expiry time set and is less than a year.

```ccl
KeyVault has attributes.exp before 365 days in all secrets
```

## Controls

* CIS Microsoft Azure Foundations Benchmark/Azure 8.2
* BSI C5/CRY-04