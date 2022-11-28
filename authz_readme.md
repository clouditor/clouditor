# MS configuration file for Clouditor - Steps

- az login
- az ad sp create-for-rbac --sdk-auth > azure.auth
	- azure.auth must be protected and not added to repo
- Set AZURE_AUTH_LOCATION to the azure.auth file
	- in IntelliJ via the configuration
	- Linux via export AZURE_AUTH_LOCATION=/XXX/YYY/.../ZZZ/azure.auth
