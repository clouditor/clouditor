module clouditor.io/clouditor/v2

go 1.22.0

// runtime dependencies (assessment)
require (
	github.com/OneOfOne/xxhash v1.2.8 // indirect
	github.com/agnivade/levenshtein v1.2.0 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/go-co-op/gocron v1.37.0
	github.com/go-ini/ini v1.67.0 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/gobwas/glob v0.2.3 // indirect
	github.com/google/uuid v1.6.0
	github.com/gorilla/mux v1.8.1 // indirect
	github.com/open-policy-agent/opa v0.69.0
	github.com/pkg/errors v0.9.1 // indirect
	github.com/prometheus/client_golang v1.20.4 // indirect
	github.com/prometheus/client_model v0.6.1 // indirect
	github.com/prometheus/common v0.55.0 // indirect
	github.com/prometheus/procfs v0.15.1 // indirect
	github.com/rcrowley/go-metrics v0.0.0-20200313005456-10cdbea86bc0 // indirect
	github.com/robfig/cron/v3 v3.0.1 // indirect
	github.com/tchap/go-patricia/v2 v2.3.1 // indirect
	github.com/xeipuuv/gojsonpointer v0.0.0-20190905194746-02993c407bfb // indirect
	github.com/xeipuuv/gojsonreference v0.0.0-20180127040603-bd5ef7bd5415 // indirect
	github.com/yashtewari/glob-intersection v0.2.0 // indirect
	go.opentelemetry.io/otel v1.28.0 // indirect
	go.opentelemetry.io/otel/metric v1.28.0 // indirect
	go.opentelemetry.io/otel/sdk v1.28.0 // indirect
	go.opentelemetry.io/otel/trace v1.28.0 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	golang.org/x/sync v0.8.0
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

// runtime dependencies (logging)
require (
	github.com/lmittmann/tint v1.0.5 // indirect
	github.com/logrusorgru/aurora/v3 v3.0.0
	github.com/sirupsen/logrus v1.9.3
	golang.org/x/sys v0.25.0 // indirect

)

// runtime dependencies (security)
require (
	github.com/MicahParks/keyfunc/v2 v2.1.0
	github.com/golang-jwt/jwt/v5 v5.2.1
	github.com/oxisto/oauth2go v0.14.0
	golang.org/x/crypto v0.27.0 // indirect
	golang.org/x/oauth2 v0.24.0

)

// runtime dependencies (rpc)
require (
	buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go v1.35.1-20240920164238-5a7b106cbb87.1
	github.com/antlr4-go/antlr/v4 v4.13.0 // indirect
	github.com/bufbuild/protovalidate-go v0.7.0
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/google/cel-go v0.21.0 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.4.0
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.22.0
	github.com/stoewer/go-strcase v1.3.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20240814211410-ddb44dafa142
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240814211410-ddb44dafa142 // indirect
	google.golang.org/grpc v1.67.0
	google.golang.org/protobuf v1.35.1

)

// runtime dependencies (storage)
require (
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/glebarez/go-sqlite v1.21.2 // indirect
	github.com/glebarez/sqlite v1.11.0
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgx/v5 v5.5.4 // indirect
	github.com/jackc/puddle/v2 v2.2.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/mattn/go-isatty v0.0.17 // indirect
	github.com/remyoudompheng/bigfft v0.0.0-20230129092748-24d4a6f8daec // indirect
	gorm.io/driver/postgres v1.5.0
	gorm.io/gorm v1.25.7
	modernc.org/libc v1.22.5 // indirect
	modernc.org/mathutil v1.5.0 // indirect
	modernc.org/memory v1.5.0 // indirect
	modernc.org/sqlite v1.23.1 // indirect

)

// runtime dependencies (config, cli)
require (
	github.com/alecthomas/kong v0.9.0 // indirect
	github.com/fsnotify/fsnotify v1.7.0 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/magiconair/properties v1.8.7 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/pelletier/go-toml/v2 v2.2.2 // indirect
	github.com/sagikazarmark/locafero v0.4.0 // indirect
	github.com/sagikazarmark/slog-shim v0.1.0 // indirect
	github.com/sourcegraph/conc v0.3.0 // indirect
	github.com/spf13/afero v1.11.0 // indirect
	github.com/spf13/cast v1.6.0 // indirect
	github.com/spf13/cobra v1.8.1
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/spf13/viper v1.19.0
	github.com/subosito/gotenv v1.6.0 // indirect
	go.uber.org/multierr v1.9.0 // indirect
	golang.org/x/exp v0.0.0-20240325151524-a685a6edb6d8 // indirect
	golang.org/x/text v0.19.0
	gopkg.in/ini.v1 v1.67.0 // indirect
)

// runtime dependencies (Azure)
require (
	github.com/Azure/azure-sdk-for-go/sdk/azcore v1.11.1
	github.com/Azure/azure-sdk-for-go/sdk/azidentity v1.6.0
	github.com/Azure/azure-sdk-for-go/sdk/internal v1.8.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/appservice/armappservice/v2 v2.3.0
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v3 v3.0.0
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/cosmos/armcosmos v1.0.0
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/dataprotection/armdataprotection v1.0.0
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork v1.1.0
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources v1.2.0
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/security/armsecurity v0.12.0
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/sql/armsql v1.2.0
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage v1.5.0
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/subscription/armsubscription v1.2.0
	github.com/AzureAD/microsoft-authentication-library-for-go v1.2.2 // indirect
	github.com/kylelemons/godebug v1.1.0 // indirect
	github.com/pkg/browser v0.0.0-20240102092130-5ac0b6a4141c // indirect
	golang.org/x/net v0.29.0 // indirect

)

// runtime dependencies (AWS)
require (
	github.com/aws/aws-sdk-go-v2 v1.32.4
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.6.6 // indirect
	github.com/aws/aws-sdk-go-v2/config v1.28.4
	github.com/aws/aws-sdk-go-v2/credentials v1.17.45 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.16.19 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.3.23 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.6.23 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.8.1 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.3.23 // indirect
	github.com/aws/aws-sdk-go-v2/service/ec2 v1.189.0
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.12.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.4.4 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.12.4 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.18.4 // indirect
	github.com/aws/aws-sdk-go-v2/service/lambda v1.66.0
	github.com/aws/aws-sdk-go-v2/service/s3 v1.67.0
	github.com/aws/aws-sdk-go-v2/service/sso v1.24.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.28.4 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.33.0
	github.com/aws/smithy-go v1.22.0
	github.com/jmespath/go-jmespath v0.4.0 // indirect

)

// runtime dependencies (k8s)
require (
	github.com/emicklei/go-restful/v3 v3.11.0 // indirect
	github.com/fxamacker/cbor/v2 v2.7.0 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-openapi/jsonpointer v0.19.6 // indirect
	github.com/go-openapi/jsonreference v0.20.2 // indirect
	github.com/go-openapi/swag v0.22.4 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/google/gnostic-models v0.6.9-0.20230804172637-c7be7c783f49 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/imdario/mergo v0.3.13 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/x448/float16 v0.8.4 // indirect
	golang.org/x/term v0.24.0 // indirect
	gopkg.in/evanphx/json-patch.v4 v4.12.0 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	k8s.io/api v0.31.1
	k8s.io/apimachinery v0.31.1
	k8s.io/client-go v0.31.1
	k8s.io/klog/v2 v2.130.1 // indirect
	k8s.io/kube-openapi v0.0.0-20240228011516-70dd3763d340 // indirect
	k8s.io/utils v0.0.0-20240711033017-18e509b52bc8 // indirect
	sigs.k8s.io/json v0.0.0-20221116044647-bc3834ca7abd // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.4.1 // indirect
	sigs.k8s.io/yaml v1.4.0 // indirect
)

// runtime dependencies (extra-csaf)
require (
	github.com/Intevation/gval v1.3.0 // indirect
	github.com/Intevation/jsonpath v0.2.1 // indirect
	github.com/ProtonMail/go-crypto v1.0.0
	github.com/cloudflare/circl v1.3.7 // indirect
	github.com/csaf-poc/csaf_distribution/v3 v3.0.1-0.20240425111311-617deb4c1721
	github.com/santhosh-tekuri/jsonschema/v5 v5.3.1 // indirect
	github.com/shopspring/decimal v1.3.1 // indirect
	go.etcd.io/bbolt v1.3.10 // indirect
	golang.org/x/time v0.6.0 // indirect

)

// testing dependencies (core)
require (
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/google/go-cmp v0.6.0
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/stretchr/testify v1.9.0
	gopkg.in/yaml.v3 v3.0.1 // indirect

)

// tools dependencies
require (
	github.com/bmatcuk/doublestar/v4 v4.0.2 // indirect
	github.com/fatih/structtag v1.2.0 // indirect
	github.com/google/addlicense v1.1.0
	github.com/lyft/protoc-gen-star v0.6.1 // indirect
	github.com/oxisto/owl2proto v0.2.1
	github.com/srikrsna/protoc-gen-gotag v1.0.0

)
