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
| RuntimeLanguage | 🚫    | ❌   |
| RuntimeVersion  | 🚫    | ❌   |

### VirtualMachine

| Evidence          | Azure | AWS |
|-------------------|-------|-----|
| Compute           | ✅     | ✅   |
| BlockStorage      | ✅     | ✅   |
| MalwareProtection | ✅     | ❌   |
| BootLogging       | ✅     | ✅   |
| OSLogging         | ✅     | ✅   |
| AutomaticUpdates  | ❌     | ❌   |

#### Compute
| Evidence          | Azure | AWS |
|-------------------|-------|-----|
| Resource          | ✅     | ✅   |
| NetworkInterfaces | ✅     | ✅   |

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
| AccessRestriction | ❌     | ❌   |
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

### Storage
| Evidence          | Azure | AWS |
|-------------------|-------|-----|
| Resource          | ✅     | ✅   |
| AtRestEncryption  | ✅     | ✅   |

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

