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
| Backups  | ❌ |    |

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
| Backups      | ✅     |  ❌  |

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
| Backups  | ✅     |    |

#### Storage
| Evidence         | Azure | AWS |
|------------------|-------|-----|
| Resource         | ✅     | ✅   |
| AtRestEncryption | ✅     | ✅   |
| Immutability     | ✅     | ❌   |
| ResourceLogging  | ✅     |    |
| Backups  | ✅     |    |

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
|Backups             |         |  ❌  |

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

### Database Storage
| Evidence     | Azure | AWS |
|--------------|-------|-----|
| Storage      | ✅     | ❌   |
| Parent       |  ✅    | ❌   |

### Database Service
| Evidence     | Azure | AWS |
|--------------|-------|-----|
| NetworkService      | ✅     | ❌   |
| AnomalyDetection       |  ✅    | ❌   |

#### Networkservice
| Evidence             | Azure | AWS |
|----------------------|-------|-----|
| Networking           | ✅     | ✅   |
| Authenticity         | ❌     | ❌   |
| Compute              | ❌     | ❌   |  
| TransportEncryption  | ❌     | ✅   |
| Ips                  | ❌     | ❌   |
| Ports                | ❌     | ❌   |
</details>

# Azure Backup
<details>
<summary>Expand</summary>

There are 2 different backup solutions for different resources
- Backup Vaults and
- Recovery Services Vault.

| Resource   | Backup Vaults | Recovery Services Vault |
|------------|-------|-----|
| Azure Virtual Machine | | x |
| Azure Storage (Files)| | x |
| Azure Backup Agent| | x |
| Azure Backup Server| | x |
| DPM| | x |
| SQL in Azure VM | | x |
| SAP HANA in Azure VM | | x |
| Azure Storage (Blobs) | x | |
| Azure disks | x | |
| Azure Database for PostgreSQL servers | x | |
| Kubernetes Services | x | |
</details>
