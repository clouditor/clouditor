# Enable Data Disk Encryption in Virtual Machines

Data disks in Virtual machines should be encrypted.

```ccl
VirtualMachine has dataDiskEncryption == "Encrypted"
```

## Remediation

* Follow the guide at https://docs.microsoft.com/en-us/azure/security/azure-security-disk-encryption-linux. 

## Controls

CIS Microsoft Azure Foundations Benchmark/Azure 7.3
