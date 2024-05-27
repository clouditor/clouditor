# Security Policy

The following file contains information about the security policy and procedures used in the Clouditor software.

## Supported Versions

We are currently in the process of moving towards a `v2` version and already released several pre-release versions of
`v2`. Please note, that development of `v2` is subject to change. If you are looking for a more stable version, please select one of the supported `v1` versions.

| Version   | Supported          |
| --------- | ------------------ |
| v1.10.1   | :white_check_mark: |
| v1.10.0   | :white_check_mark: |
| <= v1.9.5 | :x:                |

## Reporting a Vulnerability

Should you encounter a vulnerability in the Clouditor software, please use the possibility to privately report a vulnerability through GitHub using https://github.com/clouditor/clouditor/security/advisories/new.

We will then get in contact with you, assess the impact of the reported issue and try to fix it. After a fix is released, we will publish a Security Advisory (see below).

## Security Advisories

All fixed security issues will be accompanied by a security advisory. We aim to provide them in two formats

* Using GitHub's internal database (https://github.com/clouditor/clouditor/security/advisories), in order to inform GitHub users as soon as possible
* In the Clouditor repo itself in the folder [csaf](./csaf/) using the [CSAF](https://docs.oasis-open.org/csaf/csaf/v2.0/os/csaf-v2.0-os.html) standard. This allows also for a more fine-grained reporting of a security issue as well as the current status and possible affected components.