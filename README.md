# Clouditor Community Edition [![CircleCI](https://circleci.com/gh/clouditor/clouditor.svg?style=svg)](https://circleci.com/gh/clouditor/clouditor) [![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=clouditor_clouditor&metric=alert_status)](https://sonarcloud.io/dashboard?id=clouditor_clouditor) [![Coverage](https://sonarcloud.io/api/project_badges/measure?project=clouditor_clouditor&metric=coverage)](https://sonarcloud.io/dashboard?id=clouditor_clouditor) [![Bugs](https://sonarcloud.io/api/project_badges/measure?project=clouditor_clouditor&metric=bugs)](https://sonarcloud.io/dashboard?id=clouditor_clouditor) [![Vulnerabilities](https://sonarcloud.io/api/project_badges/measure?project=clouditor_clouditor&metric=vulnerabilities)](https://sonarcloud.io/dashboard?id=clouditor_clouditor)

The Clouditor is a tool to support continuous cloud assurance.

# Development

## Code Style

We use [Google Java Style](https://github.com/google/google-java-format) as a formatting. Please install the appropriate plugin for your IDE.

## Git Hooks

You can use the hook in `style/pre-commit` to check for formatting errors:
```
cp style/pre-commit .git/hooks
```

# Build (gradle)

To build the Clouditor, you can use the following gradle commands:

```
./gradlew clean build
```

# Build (Docker)

To build all necessary docker images, run the following command:

```
./gradlew docker
```

# Run (Docker)

To run the Clouditor in a demo-like mode, with no persisted database:

```
docker run -p 9999:9999 clouditor/clouditor
```

To enable auto-discovery for AWS or Azure credentials stored in your home folder, you can use:

```
docker run -v $HOME/.aws:/root/.aws -v $HOME/.azure:/root/.azure -p 9999:9999 clouditor/clouditor
```

Then open a web browser at http://localhost:9999.
