# Restrict Network Access

Checks if no SQL servers allow ingress from the internet.

```ccl
SQLServer has not properties.endIpAddress == "0.0.0.0" in all firewallRules
SQLServer has not properties.startIpAddress == "0.0.0.0" in all firewallRules
```

## Controls

* CIS Microsoft Azure Foundations Benchmark/Azure 6.3

[comment]: # TODO: should be an or?
