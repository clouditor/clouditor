package aws

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/sirupsen/logrus"
)

// awsDiscovery holds configurations across all services within AWS
type awsDiscovery struct {
	cfg aws.Config
}

// NewAwsDiscovery constructs a new awsDiscovery
// ToDo: "Overload" (switch) with staticCredentialsProvider
func NewAwsDiscovery() *awsDiscovery {
	d := &awsDiscovery{}
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		logrus.Errorf("Could not load default config: %v", err)
	}
	// ToDo: Test if proper region was loaded and maybe remove line
	logrus.Printf("Loaded credentials in region: %v", cfg.Region)
	d.cfg = cfg
	return d
}

// ToDo: I should make the services mor OO like
// DiscoverAll ToDo: Accumulate all service responses into, e.g., one JSON
func (d *awsDiscovery) discoverAll(*awsDiscovery) {
	logrus.Println("Discovering all services (s3,ec2).")
	//rawBuckets := List(GetS3Client(d.cfg))
	//for i, e := range rawBuckets.Buckets {
	//	bucket := GetObjectsOfBucket()
	//}

}
