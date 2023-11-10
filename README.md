<img src="images/claudi.png" width="60%"/>

# Clouditor Community Edition

[![build](https://github.com/clouditor/clouditor/actions/workflows/build.yml/badge.svg)](https://github.com/clouditor/clouditor/actions/workflows/build.yml)
[![](https://godoc.org/clouditor.io/clouditor?status.svg)](https://pkg.go.dev/clouditor.io/clouditor)
[![Go Report Card](https://goreportcard.com/badge/clouditor.io/clouditor)](https://goreportcard.com/report/clouditor.io/clouditor)
[![codecov](https://codecov.io/gh/clouditor/clouditor/branch/main/graph/badge.svg)](https://codecov.io/gh/clouditor/clouditor)

> :warning: Note: We are currently in the transition of re-implementing most of the core functionality of Clouditor in Go. This choice will allows us to build a more scalable, microservice-friendly version of Clouditor. While we aim to replace most of the code, we still want to provide the same look and feel as before, so we decided NOT to brand this as a v2 release, but we are rather targeting to have a `v1.5` or later release with most of the functionality done. This is an intential break with our semver approach, but we feel it is necessary to circumvent some of the pitfalls of Go's enforced SIV-style for `v2` and later.
>
> If you are looking for a stable version using only the Java code, please use the [1.2.0](https://github.com/clouditor/clouditor/releases/tag/v1.2.0) release.


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

## Build

Install necessary protobuf tools.

```
go install google.golang.org/protobuf/cmd/protoc-gen-go \
google.golang.org/grpc/cmd/protoc-gen-go-grpc \
github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
github.com/google/gnostic/cmd/protoc-gen-openapi
```

Also make sure that `$HOME/go/bin` is on your `$PATH` and build:

```
go generate ./...
go build -o ./engine cmd/engine/engine.go
```

## Usage

To test, start the engine with an in-memory DB

```
./engine --db-in-memory
```

Alternatively, be sure to start a postgre DB:

```
docker run -e POSTGRES_HOST_AUTH_METHOD=trust -d -p 5432:5432 postgres 
```

## Clouditor CLI

The Go components contain a basic CLI command called `cl`. It can be installed using `go install cmd/cli/cl.go`. Make sure that your `~/go/bin` is within your $PATH. Afterwards the binary can be used to connect to a Clouditor instance.

```bash
cl login <host:grpcPort>
```

The CLI can also be used to interact with the experimental resource graph, for example to add additional information about an application and its dependencies:

```bash
cl service discovery experimental update-resource \
'{"id": "log4j", "cloudServiceId": "00000000-0000-0000-0000-000000000000", "resourceType": "Library,Resource", "properties":{"name": "log4j", "groupId": "org.apache.logging.log4j", "artifactId": "log4j-core", "version": "2.17.0", "dependencyType": "maven", "url": "https://github.com/apache/logging-log4j2", "vulnerabilities": ["CVE-2021-44832"]}}'
cl service discovery experimental update-resource \
'{"id": "Main.java", "cloudServiceId": "00000000-0000-0000-0000-000000000000", "resourceType": "TranslationUnitDeclaration,Resource", "properties":{"name": "Main.java", "code": "class Main { public static void main(String[] args) { return; } }"}}'
cl service discovery experimental update-resource \
'{"id": "MyApplication", "cloudServiceId": "00000000-0000-0000-0000-000000000000", "resourceType": "Application,Resource", "properties":{"id:": "MyApplication", "name": "MyApplication","dependencies":["log4j"],"translationUnits":["Main.java"]}}'
cl service discovery experimental update-resource \
'{"id": "github.com/org/app", "cloudServiceId": "00000000-0000-0000-0000-000000000000", "resourceType": "CodeRepository,Resource", "properties":{"id:": "github.com/org/app", "name": "github.com/org/app", "parent": "MyApplication", "url": "github.com/org/app"}}'
```

### Command Completion

The CLI offers command completion for most shells using the `cl completion` command. Specific instructions to install the shell completions can be accessed using `cl completion --help`.
