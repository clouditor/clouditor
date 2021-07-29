package voc

type Image struct {
	*CloudResource
	Application	*Application `json:"application"`
}

