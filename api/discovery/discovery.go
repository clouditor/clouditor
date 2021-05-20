package discovery

import "clouditor.io/clouditor/voc"

type Discoverer interface {
	List() ([]voc.IsResource, error)
}
