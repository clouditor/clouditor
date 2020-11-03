# Enable Transparent Data Encryption

Checks whether all databases in Azure SQL servers performs encryption and decryption of the database, associated backups, and transaction log files at rest.

```ccl
SQLDatabase has transparentDataEncryption.status == "Enabled"
```

# Controls

* CIS Microsoft Azure Foundations Benchmark/Azure 4.2.6
* BSI C5/CRY-03