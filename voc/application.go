package voc

type Application struct {

	Compute	[]ResourceID `json:"Compute"`
	Functionality	*[]Functionality `json:"functionality"`
	ProgrammingLanguage	string `json:"programmingLanguage"`
	TranslationUnits	[]string `json:"translationUnits"`
}

