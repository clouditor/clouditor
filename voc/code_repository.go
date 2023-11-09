package voc

var CodeRepositoryType = []string{"CodeRepository", "Resource"}

type CodeRepository struct {
	*Resource
	URL string `json:"url"`
}
