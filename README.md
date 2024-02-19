<img src="images/claudi.png" width="60%"/>

# Clouditor Community Edition

[![build](https://github.com/clouditor/clouditor/actions/workflows/build.yml/badge.svg)](https://github.com/clouditor/clouditor/actions/workflows/build.yml)
[![](https://godoc.org/clouditor.io/clouditor/v2?status.svg)](https://pkg.go.dev/clouditor.io/clouditor/v2)
[![Go Report Card](https://goreportcard.com/badge/clouditor.io/clouditor/v2)](https://goreportcard.com/report/clouditor.io/clouditor/v2)
[![codecov](https://codecov.io/gh/clouditor/clouditor/branch/main/graph/badge.svg)](https://codecov.io/gh/clouditor/clouditor)

> [!NOTE]
> Note: We are currently preparing a `v2` release of Clouditor, which will be somewhat incompatible with regards to storage to `v1`. The APIs will remain largely the same, but will be improved and cleaned. We will regularly release pre-release `v2` versions, but do not have a concrete time-frame for a stable `v2` yet.
>
> If you are looking for a stable version, please use the [v1.10.1](https://github.com/clouditor/clouditor/releases/tag/v10.10.1) release.

## Introduction

Clouditor is a tool which supports continuous cloud assurance. Its main goal is to continuously evaluate if a cloud-based application (built using, e.g., Amazon Web Services (AWS) or Microsoft Azure) is configured in a secure way and thus complies with security requirements defined by, e.g., Cloud Computing Compliance Controls Catalogue (C5) issued by the German Office for Information Security (BSI) or the Cloud Control Matrix (CCM) published by the Cloud Security Alliance (CSA).

## Features

Clouditor currently supports over 60 checks for Amazon Web Services (AWS), Microsoft Azure and OpenStack. Results of these checks are evaluated against security requirements of the BSI C5 and CSA CCM.

Key features are:

- automated compliance rules for AWS and MS Azure
- granular report of detected non-compliant configurations
- quick and adaptive integration with existing service through automated service discovery
- descriptive development of custom rules using [Cloud Compliance Language (CCL)](clouditor-engine-azure/src/main/resources/rules/azure/compute/vm-data-encryption.md) to support individual evaluation scenarios
- integration of custom security requirements and mapping to rules

## QuickStart with UI

In order to just build and run the Clouditor, without generating the protobuf file, one can use the `run-engine-with-ui.sh` script. This still requires Go and Node.js to be installed. For example, to run the engine in-memory with the Azure provider the following command can be used:

```
./run-engine-with-ui.sh --discovery-provider=azure
```

This will start the all-in-on-engine with all discoverers enabled and launches the UI on http://localhost:5173. The
default credentials are clouditor/clouditor.

## Build

Install necessary protobuf tools, including `buf`. Please refer to the [`buf` install guide](https://buf.build/docs/installation).

```
go install github.com/srikrsna/protoc-gen-gotag \
github.com/oxisto/owl2proto/cmd/owl2proto
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
'{"id": "MyApplication", "cloudServiceId": "00000000-0000-0000-0000-000000000000", "resourceType": "Application,Resource", "properties":{"@type":"type.googleapis.com/clouditor.ontology.v1.Application", "id:": "MyApplication", "name": "MyApplication","dependencies":["log4j"],"translationUnits":["Main.java"]}}'
cl service discovery experimental update-resource \
'{"id": "github.com/org/app", "cloudServiceId": "00000000-0000-0000-0000-000000000000", "resourceType": "CodeRepository,Resource", "properties":{"id:": "github.com/org/app", "name": "github.com/org/app", "parent": "MyApplication", "url": "github.com/org/app"}}'
```

### Command Completion

The CLI offers command completion for most shells using the `cl completion` command. Specific instructions to install the shell completions can be accessed using `cl completion --help`.
