![clouditor](images/claudi.png "Clouditor")

# Clouditor Community Edition
![build](https://github.com/clouditor/clouditor/workflows/build/badge.svg) 
[![](https://godoc.org/clouditor.io/clouditor?status.svg)](https://pkg.go.dev/clouditor.io/clouditor)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=clouditor_clouditor&metric=alert_status)](https://sonarcloud.io/dashboard?id=clouditor_clouditor) 
[![Docker Pulls](https://img.shields.io/docker/pulls/clouditor/clouditor.svg)](https://hub.docker.com/r/clouditor/clouditor)
[![Coverage](https://sonarcloud.io/api/project_badges/measure?project=clouditor_clouditor&metric=coverage)](https://sonarcloud.io/dashboard?id=clouditor_clouditor) 
[![Bugs](https://sonarcloud.io/api/project_badges/measure?project=clouditor_clouditor&metric=bugs)](https://sonarcloud.io/dashboard?id=clouditor_clouditor) 
[![Vulnerabilities](https://sonarcloud.io/api/project_badges/measure?project=clouditor_clouditor&metric=vulnerabilities)](https://sonarcloud.io/dashboard?id=clouditor_clouditor)


## Introduction

Clouditor is a tool which supports continuous cloud assurance. Its main goal is to continuously evaluate if a cloud-based application (built using, e.g., Amazon Web Services (AWS) or Microsoft Azure) is configured in a secure way and thus complies with security requirements defined by, e.g., Cloud Computing Compliance Controls Catalogue (C5) issued by the German Office for Information Security (BSI) or the Cloud Control Matrix (CCM) published by the Cloud Security Alliance (CSA).

## Features

Clouditor currently supports over 60 checks for Amazon Web Services (AWS), Microsoft Azure and OpenStack. Results of these checks are evaluated against security requirements of the BSI C5 and CSA CCM.

Key features are:

* automated compliance rules for AWS and MS Azure
* granular report of detected non-compliant configurations
* quick and adaptive integration with existing service through automated service discovery
* descriptive development of custom rules using [Cloud Compliance Language (CCL)](clouditor-engine-azure/src/main/resources/rules/azure/compute/vm-data-encryption.md) to support individual evaluation scenarios
* integration of custom security requirements and mapping to rules

## Usage

To run the Clouditor in a demo-like mode, with no persisted database:

```
docker run -p 9999:9999 clouditor/clouditor
```

To enable auto-discovery for AWS or Azure credentials stored in your home folder, you can use:

```
docker run -v $HOME/.aws:/root/.aws -v $HOME/.azure:/root/.azure -p 9999:9999 clouditor/clouditor
```

Then open a web browser at http://localhost:9999. Login with user `clouditor` and the default password `clouditor`.


## Screenshots

#### Configuring an account
![Account configuration](images/Accounts.png "Accounts")

#### Discovering resources of cloud-based application

![Discovery view](/images/Discovery.png "Discovery")

#### Overview of rule-based assessment 

![Rule assessment](images/Rules.png "Assessment")

#### View details of rules

![Rule assessment](images/Assessment.png "Assessment")

#### Load and map compliance requirements

![Compliance overview](images/Compliance.png "Compliance")

## Development

### Code Style

We use [Google Java Style](https://github.com/google/google-java-format) as a formatting. Please install the appropriate plugin for your IDE.

### Git Hooks

You can use the hook in `style/pre-commit` to check for formatting errors:
```
cp style/pre-commit .git/hooks
```

### Build (gradle)

To build the Clouditor, you can use the following gradle commands:

```
./gradlew clean build
```

### Build (Docker)

To build all necessary docker images, run the following command:

```
./gradlew docker
```

### Build (Go components) - Experimental

Install necessary protobuf tools.

```
go install google.golang.org/protobuf/cmd/protoc-gen-go \
google.golang.org/grpc/cmd/protoc-gen-go-grpc \
github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway
```

Also make sure that `$HOME/go/bin` is on your `$PATH` and build:

```
go generate ./...
go build ./...
```

To test, start the engine with an in-memory DB

```
./engine --db-in-memory
```

Alternatively, be sure to start a postgre DB:

```
docker run -e POSTGRES_HOST_AUTH_METHOD=trust -d -p 5432:5432 postgres 
```

### Clouditor CLI

The Go components contain a basic CLI command called `cl`. It can be installed using `go install cmd/cli/cl.go`. Make sure that your `~/go/bin` is within your $PATH. Afterwards the binary can be used to connect to a Clouditor instance.

```bash
cl login <host:grpcPort>
```

#### Command Completion

The CLI offers command completion for most shells using the `cl completion` command. For example, `zsh` can be configured in the following way.

```bash
echo """
autoload -Uz compinit && compinit -C
source <(cl completion zsh)
compdef _cl cl
""" >> ~/.zshrc
```
