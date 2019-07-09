# eMail Configuration for Advanced Data Security"

Checks, whether the SQL server has proper email configuration for Advanced Data Security

```ccl
SQLServer has not empty securityAlertPolicy.emailAddresses
SQLServer has securityAlertPolicy.emailAccountAdmins == true
```

# Controls

* CIS Microsoft Azure Foundations Benchmark/Azure 4.1.4
* CIS Microsoft Azure Foundations Benchmark/Azure 4.1.5
