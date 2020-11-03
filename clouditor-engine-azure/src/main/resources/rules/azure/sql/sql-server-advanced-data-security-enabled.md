# Enable Advanced Data Security

Checks, whether the SQL server has Advanced Data Security enabled.

```ccl
SQLServer has securityAlertPolicy.state == "Enabled"
SQLServer has not empty securityAlertPolicy.storageEndpoint
```

## Cntrols

* CIS Microsoft Azure Foundations Benchmark/Azure 4.1.2
* BSI C5/OPS-10
* BSI C5/OPS-13