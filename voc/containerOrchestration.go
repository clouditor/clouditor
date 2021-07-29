package voc

type ContainerOrchestration struct {
	*CloudResource
	Container	[]ResourceID `json:"container"`
	ResourceLogging	[]ResourceID `json:"resourceLogging"`
	ManagementUrl	string `json:"managementUrl"`
}

