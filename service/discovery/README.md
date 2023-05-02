# Discovery Status

✅: Discovered <br />
❌: Not Discovered <br />
🚫: Not available <br />

### Compute
<details>
<summary>Expand</summary>

### Function

| Evidence        | Azure | AWS |
|-----------------|-------|-----|
| Compute         | ✅     | ✅   |
| RuntimeLanguage | ✅    | ❌   |
| RuntimeVersion  | ✅    | ❌   |

### VirtualMachine

| Evidence          | Azure | AWS |
|-------------------|-------|-----|
| Compute           | ✅     | ✅   |
| BlockStorage      | ✅     | ✅   |
| MalwareProtection | ✅     | ❌   |
| BootLogging       | ✅     | ✅   |
| OSLogging         | ✅     | ✅   |
| AutomaticUpdates  | ✅     | ❌   |

#### Compute
| Evidence          | Azure | AWS |
|-------------------|-------|-----|
| Resource          | ✅     | ✅   |
| NetworkInterfaces | ✅     | ✅   |
| ResourceLogging  | ✅     |    |
| Backup  | ✅    |    |

#### Resource
| Evidence    | Azure | AWS |
|-------------|-------|-----|
| ID          | ✅     | ✅   |
| Name        | ✅     | ✅   |
| Type        | ✅     | ✅   |
| GeoLocation | ✅     | ✅   |
| Labels      | ✅     | ✅   |

#### OSLogging
| Evidence        | Azure | AWS |
|-----------------|-------|-----|
| Auditing        | ✅     | 🚫  |
| SecurityFeature | ✅     | 🚫  |
| Enabled         | ✅     | ❌   |
| LoggingService  | ✅     | 🚫  |
| RetentionPeriod | ✅     | 🚫  |

#### BootLogging
| Evidence        | Azure | AWS |
|-----------------|-------|-----|
| Auditing        | ✅     | 🚫  |
| SecurityFeature | ✅     | 🚫  |
| Enabled         | ✅     | ❌   |
| LoggingService  | ✅     | 🚫  |
| RetentionPeriod | ✅     | 🚫  |

#### ResourceLogging
| Evidence                  | Azure | AWS |
|-----------------          |-------|-----|
| MonitoringLogDataEnabled  | ✅     |   |
| SecurityAlertsEnabled     | ✅     |   |


### BlockStorage

| Evidence          | Azure | AWS |
|-------------------|-------|-----|
| Resource          | ✅     | ✅   |
| AtRestEncryption  | ✅     | ✅   |
| Immutability      | ✅     | ❌   |

#### ManagedKeyEncryption
| Evidence  | Azure | AWS |
|-----------|-------|-----|
| Enabled   | ✅     | ❌   |
| Algorithm | ✅     | ❌   |

#### CustomerKeyEncryption
| Evidence  | Azure | AWS |
|-----------|-------|-----|
| Enabled   | ✅     | ❌   |
| Algorithm | ❌     | ❌   |
| KeyUrl    | ✅     | ❌   |

</details>

### Network
<details>
<summary>Expand</summary>


### LoadBalancer
| Evidence          | Azure | AWS |
|-------------------|-------|-----|
| Networkservice    | ✅     | ❌   |
| AccessRestriction | ✅     | ❌   |
| HttpEndpoints     | ✅     | ❌   |
| Networkservices   | ✅     | ❌   |
| Urls              | ✅     | ❌   |

#### Networkservice
| Evidence             | Azure | AWS |
|----------------------|-------|-----|
| Networking           | ✅     | ❌   |
| Authenticity         | ✅     | ❌   |
| Compute              | ✅     | ❌   |
| TransportEncryption  | ✅     | ❌   |
| Ips                  | ✅     | ❌   |
| Ports                | ✅     | ❌   |

### Networkinterfaces
| Evidence          | Azure | AWS |
|-------------------|-------|-----|
| Networking        | ✅     | ✅   |
| Networkservice    | ❌     | ❌   |
| AccessRestriction | partly     | ❌   |
</details>

### Storage
<details>
<summary>Expand</summary>

### ObjectStorage
| Evidence     | Azure | AWS |
|--------------|-------|-----|
| Storage      | ✅     | ✅   |
| PublicAccess | ✅     | ❌   |

#### Storage
| Evidence         | Azure | AWS |
|------------------|-------|-----|
| Resource         | ✅     | ✅   |
| AtRestEncryption | ✅     | ✅   |
| Immutability     | ✅     | ❌   |
| ResourceLogging  | ✅     |    |
| Backup  | ✅     |    |

### ObjectStorageService
| Evidence       | Azure | AWS |
|----------------|-------|-----|
| NetworkService | ✅     | ✅   |
| HttpEndpoint   | ✅     | ✅   |

#### Networkservice
| Evidence             | Azure | AWS |
|----------------------|-------|-----|
| Networking           | ✅     | ✅   |
| Authenticity         | ❌     | ❌   |
| Compute              | ❌     | ❌   |  
| TransportEncryption  | ✅     | ✅   |
| Ips                  | ❌     | ❌   |
| Ports                | ❌     | ❌   |

#### HttpEndpoint
| Evidence            | Azure    | AWS |
|---------------------|----------|-----|
| Url                 | ✅        | ✅   |
| TransportEncryption | ✅        | ✅   |

### FileStorage
| Evidence          | Azure | AWS |
|-------------------|-------|-----|
| Storage           | ✅     | ❌   |

#### ManagedKeyEncryption
| Evidence  | Azure | AWS |
|-----------|-------|-----|
| Enabled   | ✅     | ✅   |
| Algorithm | ✅     | ✅   |

#### CustomerKeyEncryption
| Evidence   | Azure | AWS |
|------------|-------|-----|
| Enabled    | ✅     | ✅   |
| Algorithm  | ❌     | ❌   |
| KeyUrl     | ✅     | ✅   |

</details>

# Azure Backup
There are 2 different backup solutions for different resources
- Backup Vaults and
- Recovery Services Vault.

| Resource   | Backup Vaults | Recovery Services Vault |
|------------|-------|-----|
| Azure Virtual Machine | | ✅ |
| Azure Storage (Files)| | ✅ |
| Azure Backup Agent| | ✅ |
| Azure Backup Server| | ✅ |
| DPM| | ✅ |
| SQL in Azure VM | | ✅ |
| SAP HANA in Azure VM | | ✅ |
| Azure Storage (Blobs) | ✅ | |
| Azure disks | ✅ | |
| Azure Database for PostgreSQL servers | ✅ | |
| Kubernetes Services | ✅ | |
