package aws

import (
	"testing"
)

// ToDo: Re-writing

// ToDo: Works with my credentials -> Mock it
func Test_awsDiscovery_NewAwsDiscovery(t *testing.T) {
	testDiscovery := NewAwsDiscovery()
	if region := testDiscovery.cfg.Region; region != "eu-central-1" {
		t.Fatalf("Excpected eu-central-1. Got %v", region)
	}
}

// ToDo: Works with my credentials -> Mock it
func Test_discoverAll(t *testing.T) {
	testDiscovery := NewAwsDiscovery()
	discoverAll(testDiscovery)
	if region := testDiscovery.cfg.Region; region != "eu-central-1" {
		t.Fatalf("Excpected eu-central-1. Got %v", region)
	}

}
