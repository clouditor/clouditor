# Storage Configuration for Advanced Data Security

Checks, whether the SQL server has proper storage configuration for Advanced Data Security

```ccl
SQLServer has not empty securityAlertPolicy.storageEndpoint
SQLServer has securityAlertPolicy.retentionDays >= 90
```

## Controls

* CIS Microsoft Azure Foundations Benchmark/Azure 4.1.7

[comment]: # TODO: retentionDays == 0 is also ok
