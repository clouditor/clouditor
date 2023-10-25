package voc

var KeyVaultType = []string{"KeyVault", "Resource"}

type KeyVault struct {
	*Resource
	IsActive bool
	Keys     []*Key `json:"keys"`
}
