# Set Expiry Time for Secrets

Checks if all secrets in Azure Key Vault have an expiry time set and is less than a year.

```ccl
KeyVault has attributes.exp before 365 days in all secrets
```

condition: 
controls:
  - "CIS Microsoft Azure Foundations Benchmark/Azure 8.2"
