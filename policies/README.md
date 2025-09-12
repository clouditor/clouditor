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

