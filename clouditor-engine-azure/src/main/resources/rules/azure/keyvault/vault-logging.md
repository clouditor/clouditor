## Enable Logging for Vaults

Checks if AuditEvent logging for Key Vault instances is enabled.

```ccl
# TODO: the category should be AuditEvent, and also 0 days are valid (forever) and also enabled should be true
KeyVault has retentionPolicy.days >= 180 in any logs
```

## Controls

* CIS Microsoft Azure Foundations Benchmark/Azure 5.13
* BSI C5/OPS-10