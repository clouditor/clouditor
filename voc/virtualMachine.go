package voc

type VirtualMachine struct {
	*Compute
	NetworkInterface	[]ResourceID `json:"networkInterface"`
	Log	*Log `json:"log"`
	BlockStorage	[]ResourceID `json:"blockStorage"`
}

