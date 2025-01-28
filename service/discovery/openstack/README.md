# Openstack Discovery
OpenStack discovery is a feature of Clouditor that discovers the information about OpenStack environments through API calls. It discovers storage, virtual machines and networks. With sufficient permissions it is also possible to discover domains and projects. 

# Limitations
## Application Credentials
In OpenStack, application credentials are created specifically for a designated project within a specific domain. This makes the discovery of domains and projects unnecessary. However, a limitation is that the domain ID/name and project ID/name cannot be discovered without the appropriate permissions. 
*NOTE:* At the time of writing, the necessary permissions for this discovery are not yet known. 

- Adding OS_TENANT_ID or OS_TENANT_NAME as environment variables is not possible. Attempting to do so results in an error message:
`Error authenticating with application credential: Application credentials cannot request a scope.`
- The domain ID and domain name must be provided by the environment variables: `OS_PROJECT_DOMAIN_ID` and `OS_USER_DOMAIN_NAME`