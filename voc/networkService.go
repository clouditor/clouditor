package voc

type NetworkService struct {
	*Networking
	Compute	[]ResourceID `json:"compute"`
	Ips	[]string `json:"ips"`
	Ports	[]int16 `json:"ports"`
}

