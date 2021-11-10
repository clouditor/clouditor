module clouditor.io/clouditor

go 1.16

require (
	github.com/Azure/azure-sdk-for-go v55.8.0+incompatible
	github.com/Azure/go-autorest/autorest v0.11.19
	github.com/Azure/go-autorest/autorest/azure/auth v0.5.8
	github.com/Azure/go-autorest/autorest/to v0.4.0
	github.com/Azure/go-autorest/autorest/validation v0.3.1 // indirect
	github.com/alexedwards/argon2id v0.0.0-20210511081203-7d35d68092b8
	github.com/antlr/antlr4/runtime/Go/antlr v0.0.0-20210521184019-c5ad59b459ec
	github.com/aws/aws-sdk-go-v2 v1.11.0
	github.com/aws/aws-sdk-go-v2/config v1.10.0
	github.com/aws/aws-sdk-go-v2/service/ec2 v1.19.0
	github.com/aws/aws-sdk-go-v2/service/lambda v1.11.0
	github.com/aws/aws-sdk-go-v2/service/s3 v1.18.0
	github.com/aws/aws-sdk-go-v2/service/sts v1.9.0
	github.com/aws/smithy-go v1.9.0
	github.com/go-co-op/gocron v1.9.0
	github.com/gofrs/uuid v4.0.0+incompatible // indirect
	github.com/golang-jwt/jwt/v4 v4.1.0
	github.com/google/addlicense v0.0.0-20210428195630-6d92264d7170
	github.com/google/uuid v1.3.0
	github.com/googleapis/gnostic v0.5.6-0.20210520165051-0320d74b3646
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.6.0
	github.com/imdario/mergo v0.3.8 // indirect
	github.com/lib/pq v1.9.0 // indirect
	github.com/logrusorgru/aurora/v3 v3.0.0
	github.com/nxadm/tail v1.4.8 // indirect
	github.com/onsi/ginkgo v1.15.0 // indirect
	github.com/onsi/gomega v1.10.5 // indirect
	github.com/open-policy-agent/opa v0.33.1
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/cobra v1.2.1
	github.com/spf13/viper v1.9.0
	github.com/stretchr/testify v1.7.0
	google.golang.org/genproto v0.0.0-20210903162649-d08c68adba83
	google.golang.org/grpc v1.41.0
	google.golang.org/grpc/cmd/protoc-gen-go-grpc v1.1.0
	google.golang.org/protobuf v1.27.1
	gorm.io/driver/postgres v1.1.0
	gorm.io/driver/sqlite v1.1.4
	gorm.io/gorm v1.22.2
	k8s.io/api v0.22.0
	k8s.io/apimachinery v0.22.0
	k8s.io/client-go v0.22.0
)
