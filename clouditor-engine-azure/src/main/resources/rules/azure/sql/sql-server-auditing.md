# Enable Auditing

Checks whether SQL servers have auditing enabled.

```ccl
SQLServer has auditingPolicy.state == "Enabled"
SQLServer has auditingPolicy.retentionDays >= 90
```

# Controls

* CIS Microsoft Azure Foundations Benchmark/Azure 4.1.1
* CIS Microsoft Azure Foundations Benchmark/Azure 4.1.6

[comment]: # # TODO: retentionDays == 0 is also valid, since this means forever
