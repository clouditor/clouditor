module clouditor.io/clouditor

go 1.15

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/golang/protobuf v1.5.2
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.4.0
	github.com/jackc/pgx/v4 v4.9.2 // indirect
	github.com/oxisto/go-httputil v0.3.7
	github.com/plgd-dev/kit v0.0.0-20201102152602-1e03187a6a3a
	github.com/spf13/cobra v1.1.1
	github.com/spf13/viper v1.7.0
	golang.org/x/crypto v0.0.0-20201117144127-c1f2f97bffc9 // indirect
	google.golang.org/genproto v0.0.0-20210426193834-eac7f76ac494
	google.golang.org/grpc v1.37.1
	google.golang.org/grpc/cmd/protoc-gen-go-grpc v1.0.1
	google.golang.org/protobuf v1.26.0
	gorm.io/driver/postgres v1.0.5
	gorm.io/driver/sqlite v1.1.3
	gorm.io/gorm v1.21.9
)
