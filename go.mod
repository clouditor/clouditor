module clouditor.io/clouditor

go 1.16

require (
	github.com/Azure/azure-sdk-for-go v55.0.0+incompatible
	github.com/Azure/go-autorest/autorest v0.11.17
	github.com/Azure/go-autorest/autorest/azure/auth v0.5.7
	github.com/Azure/go-autorest/autorest/date v0.3.0
	github.com/Azure/go-autorest/autorest/to v0.4.0
	github.com/Azure/go-autorest/autorest/validation v0.3.1 // indirect
	github.com/alexedwards/argon2id v0.0.0-20210511081203-7d35d68092b8
	github.com/golang-jwt/jwt v0.0.0-20210529135444-590e8c64cd52
	github.com/google/addlicense v0.0.0-20210428195630-6d92264d7170
	github.com/googleapis/gnostic v0.5.6-0.20210520165051-0320d74b3646
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.4.0
	github.com/logrusorgru/aurora/v3 v3.0.0
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/cobra v1.1.3
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.7.0
	golang.org/x/crypto v0.0.0-20210513164829-c07d793c2f9a // indirect
	google.golang.org/genproto v0.0.0-20210524171403-669157292da3
	google.golang.org/grpc v1.38.0
	google.golang.org/grpc/cmd/protoc-gen-go-grpc v1.1.0
	google.golang.org/protobuf v1.26.0
	gorm.io/driver/postgres v1.1.0
	gorm.io/driver/sqlite v1.1.4
	gorm.io/gorm v1.21.10
)

replace github.com/dgrijalva/jwt-go => github.com/golang-jwt/jwt v0.0.0-20210527125010-9e96e9651418
