# Enable Flow Logs in Network Security Groups

Checks if Network Security Group Flow Logs are enabled and retention period is set to greater than or equal to 90 days.

```ccl
NetworkSecurityGroup has flowLogSettings.enabled == true
NetworkSecurityGroup has flowLogSettings.retentionPolicy.enabled == true
NetworkSecurityGroup has flowLogSettings.retentionPolicy.days >= 90
```

## Controls

* CIS Microsoft Azure Foundations Benchmark/Azure 6.4
* BSI C5/OPS-10
* BSI C5/OPS-13
