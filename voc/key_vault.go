package voc

var KeyVaultType = []string{"KeyVault", "Resource"}

type KeyVault struct {
	*Resource
	IsActive bool         `json:"isActive"`
	Keys     []ResourceID `json:"keys"`
}
