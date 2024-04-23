package csaf

import "clouditor.io/clouditor/v2/api/ontology"

func documentValidationErrors(messages []string) (errs []*ontology.Error) {
	for _, m := range messages {
		errs = append(errs, &ontology.Error{Message: m})
	}
	return
}
