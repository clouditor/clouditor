# Configuration
All configurations for the service can be set using command line flags and can be displayed via `--help`. The configuration options will be described for each component individually.

## TODO/Questions
- [ ] I think we should only add the Discovery specific flags. The "problem" is that several flags are added for each component, e.g., db-X, but are not needed for all components. Are the only necessary flags the ones which can be found [here](../server/commands/discovery/discovery.go)
- [ ] Sollten nicht alle Service Optionen auch einen flag für den service haben? ZB WithoutEvidenceStore hat keinen.


## Discovery

```
Usage:
  discovery [flags]

Flags:
      --api-cors-allowed-headers stringArray   Specifies the headers allowed in CORS (default [Content-Type,Accept,Authorization])
      --api-cors-allowed-methods stringArray   Specifies the methods allowed in CORS (default [GET,POST,PUT,DELETE])
      --api-cors-allowed-origins stringArray   Specifies the origins allowed in CORS
      --api-default-password string            Specifies the default API password (default "clouditor")
      --api-default-user string                Specifies the default API username (default "clouditor")
      --api-grpc-port uint16                   Specifies the port used for the Clouditor gRPC API (default 9091)
      --api-http-port uint16                   Specifies the port used for the Clouditor HTTP API (default 8081)
      --api-jwks-url string                    Specifies the JWKS URL used to verify authentication tokens in the gRPC and HTTP API (default "http://localhost:8080/v1/auth/certs")
      --api-key-password string                Specifies the password used to protect the API private key (default "changeme")
      --api-key-path string                    Specifies the location of the API private key (default "~/.clouditor/api.key")
      --api-key-save-on-create                 Specifies whether the API key should be saved on creation. It will only created if the default location is used. (default true)
      --assessment-url string                  Specifies the Assessment URL (default "localhost:9093")
      --certification-target-id string         Specifies the Certification Target ID (default "00000000-0000-0000-0000-000000000000")
      --dashboard-callback-url string          The callback URL of the Clouditor Dashboard. If the embedded server is used, a public OAuth 2.0 client based on this URL will be added (default "http://localhost:8080/callback")
      --db-host string                         Provides address of database (default "localhost")
      --db-in-memory                           Uses an in-memory database which is not persisted at all
      --db-name string                         Provides name of database (default "postgres")
      --db-password string                     Provides password of database (default "postgres")
      --db-port uint16                         Provides port for database (default 5432)
      --db-ssl-mode string                     The SSL mode for the database (default "disable")
      --db-user-name string                    Provides user name of database (default "postgres")
      --discovery-auto-start                   Automatically start the discovery when engine starts
      --discovery-csaf-domain string           The domain to look for a CSAF provider, if the CSAF discovery is enabled
  -p, --discovery-provider strings             Providers to discover, separated by comma
      --discovery-resource-group string        Limit the scope of the discovery to a resource group (currently only used in the Azure discoverer
  -h, --help                                   help for discovery
      --log-level string                       The default log level (default "info")
      --service-oauth2-client-id string        Specifies the OAuth 2.0 client ID (default "clouditor")
      --service-oauth2-client-secret string    Specifies the OAuth 2.0 client secret (default "clouditor")
      --service-oauth2-token-endpoint string   Specifies the OAuth 2.0 token endpoint (default "http://localhost:8080/v1/auth/token")
```