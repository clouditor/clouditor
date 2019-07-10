# Enable Auditing

Checks whether SQL databases have auditing enabled.

```ccl
SQLDatabase has auditingPolicy.state == "Enabled"
SQLDatabase has auditingPolicy.retentionDays > 90
```

## Controls

* CIS Microsoft Azure Foundations Benchmark/Azure 4.2.1
* CIS Microsoft Azure Foundations Benchmark/Azure 4.2.7

[comment]: # TODO: retentionDays == 0 is also valid, since this means forever
