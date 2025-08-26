<img src="images/claudi.png" width="60%" alt="Clouditor logo"/>

# Clouditor Community Edition

[![build](https://github.com/clouditor/clouditor/actions/workflows/build.yml/badge.svg)](https://github.com/clouditor/clouditor/actions/workflows/build.yml)
[![](https://godoc.org/clouditor.io/clouditor/v2?status.svg)](https://pkg.go.dev/clouditor.io/clouditor/v2)
[![Go Report Card](https://goreportcard.com/badge/clouditor.io/clouditor/v2)](https://goreportcard.com/report/clouditor.io/clouditor/v2)
[![codecov](https://codecov.io/gh/clouditor/clouditor/branch/main/graph/badge.svg)](https://codecov.io/gh/clouditor/clouditor)
[![OpenSSF Best Practices](https://www.bestpractices.dev/projects/8691/badge)](https://www.bestpractices.dev/projects/8691)

> [!NOTE]
> Note: We are currently preparing a `v2` release of Clouditor, which will be somewhat incompatible with regards to storage to `v1`. The APIs will remain largely the same, but will be improved and cleaned. We will regularly release pre-release `v2` versions, but do not have a concrete time-frame for a stable `v2` yet.
>
> If you are looking for a stable version, please use the [v1.10.1](https://github.com/clouditor/clouditor/releases/tag/v1.10.1) release.

## Introduction

Clouditor is a tool which supports continuous cloud assurance. Its main goal is to continuously evaluate if a cloud-based application (built using, e.g., Amazon Web Services (AWS) or Microsoft Azure) is configured in a secure way and thus complies with security requirements defined by, e.g., Cloud Computing Compliance Controls Catalogue (C5) issued by the German Office for Information Security (BSI) or the Cloud Control Matrix (CCM) published by the Cloud Security Alliance (CSA).

## Features

Clouditor currently supports over 60 checks for Amazon Web Services (AWS), Microsoft Azure and OpenStack. Results of these checks are evaluated against security requirements of the BSI C5 and CSA CCM.

Key features are:

- automated compliance rules for AWS and MS Azure
- granular report of detected non-compliant configurations
- quick and adaptive integration with existing service through automated service discovery
- curated security metrics integrated from the external Security Metrics repository
- integration of custom security requirements and mapping to rules

## Quick start with UI

To quickly build and run Clouditor without generating protobuf files, you can use the `run-engine-with-ui.sh` script. This still requires Go and Node.js to be installed. For example, to run the engine in memory with the Azure provider, use:

```
./run-engine-with-ui.sh --discovery-provider=azure
```

This starts the all-in-one engine with all discoverers enabled and launches the UI at http://localhost:3000.
The default credentials are clouditor/clouditor.

### Using the extra discoverers (e.g., CSAF)

In addition to the regular cloud provider discoverers, Clouditor also includes a set of extra discoverers for dedicated protocols, for example CSAF. The CSAF discoverer allows the conformance check of a CSAF (trusted) provider.

It can be used with the following command:

```
./run-engine-with-ui.sh --discovery-provider=csaf --discovery-csaf-domain=clouditor.io
```

The domain `clouditor.io` can be replaced with your actual domain.

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

## Security metrics integration

Clouditor integrates security metrics from an external, community-driven repository to decouple metric content from the engine and enable independent versioning and collaboration.

- Repository: https://github.com/Cybersecurity-Certification-Hub/security-metrics
- Integration method: Git submodule located at `policies/security-metrics`
- Default load path: Clouditor loads metric definitions and implementations from `./policies/security-metrics/metrics` at startup.

Getting the metrics after cloning this repository:

```bash
# Initialize submodules (only needed once per fresh clone)
git submodule update --init --recursive
```

Updating to the latest metrics:

```bash
# Update the security-metrics submodule to the latest commit on its default branch
git submodule update --remote policies/security-metrics
# Commit the updated submodule pointer in this repo
git add policies/security-metrics && git commit -m "chore: bump security-metrics"
```

Notes:
- Submodules pin an exact commit. Updating the metrics in Clouditor requires committing the new submodule reference.
- The metrics follow a defined schema (see `policies/security-metrics/metric_schema.json`) and are implemented using OPA/Rego files and metadata.
- Custom or experimental metrics can be developed by contributing to the external repository above.

### About the Security Metrics repository

The external security-metrics repository collects reusable security metrics for continuous certification. It is organized as follows:

- api: preliminary place to define the metric data format programmatically (e.g., via protobuf)
- catalogs: definitions of certification catalogs/benchmarks and mappings from catalog requirements to metrics
- metrics: the actual metrics, grouped by domain; each metric folder typically contains:
  - metric.yml: metadata and configuration for the metric (see structure below)
  - metric.rego: the Rego implementation evaluable by OPA
- ontology: the domain ontology (possibly multiple versions) that underpins metric descriptions (resources, security properties, etc.)

Metric data structure

- Metadata (static)
  - id: human-readable unique identifier
  - description: must reference ontology terms and config parameters in brackets, e.g., [BlockStorage], [p1:AtRestEncryption]
  - version: metric version
  - comments: reasoning, applicability, and examples
- Configuration (context-dependent)
  - interval: integer hours for evidence collection (e.g., 24)
  - operator: one of the basic operators (==, >=, <, ...)
  - targetValue: the value to compare against (e.g., true for AtRestEncryptionEnabled)

Ontology overview

The ontology harmonizes evidence gathering and assessment across providers and environments. It includes taxonomies for resources (e.g., Compute, Networking; VMs, containers, functions, etc.) and security properties organized by STRIDE categories (Authentication, Integrity, Non-repudiation, Confidentiality, Availability, Authorization). This harmonization enables provider-agnostic metrics and reusable mappings between catalog requirements and ontological concepts. The ontology is authored in Protégé and published on WebProtégé.

## Usage

To test, start the engine with an in-memory database:

```
./engine --db-in-memory
```

Alternatively, start a PostgreSQL database:

```
docker run -e POSTGRES_HOST_AUTH_METHOD=trust -d -p 5432:5432 postgres
```


## Clouditor CLI

The Go components contain a basic CLI command called `cl`. It can be installed using `go install ./cmd/cli`. Make sure that your `~/go/bin` is within your `$PATH`. Afterwards, the binary can be used to connect to a Clouditor instance.

```bash
cl login <host:grpcPort>
```

The CLI can also be used to interact with the experimental resource graph, for example to add additional information about an application and its dependencies:

```bash
cl service discovery experimental update-resource \
'{"id": "log4j", "targetOfEvaluationId": "00000000-0000-0000-0000-000000000000", "resourceType": "Library,Resource", "properties":{"name": "log4j", "groupId": "org.apache.logging.log4j", "artifactId": "log4j-core", "version": "2.17.0", "dependencyType": "maven", "url": "https://github.com/apache/logging-log4j2", "vulnerabilities": ["CVE-2021-44832"]}}'
cl service discovery experimental update-resource \
'{"id": "Main.java", "targetOfEvaluationId": "00000000-0000-0000-0000-000000000000", "resourceType": "TranslationUnitDeclaration,Resource", "properties":{"name": "Main.java", "code": "class Main { public static void main(String[] args) { return; } }"}}'
cl service discovery experimental update-resource \
'{"id": "MyApplication", "targetOfEvaluationId": "00000000-0000-0000-0000-000000000000", "resourceType": "Application,Resource", "properties":{"@type":"type.googleapis.com/clouditor.ontology.v1.Application", "id:": "MyApplication", "name": "MyApplication","dependencies":["log4j"],"translationUnits":["Main.java"]}}'
cl service discovery experimental update-resource \
'{"id": "github.com/org/app", "targetOfEvaluationId": "00000000-0000-0000-0000-000000000000", "resourceType": "CodeRepository,Resource", "properties":{"id:": "github.com/org/app", "name": "github.com/org/app", "parent": "MyApplication", "url": "github.com/org/app"}}'
```

### Command completion

The CLI offers command completion for most shells using the `cl completion` command. Specific instructions to install the shell completions can be accessed using `cl completion --help`.

