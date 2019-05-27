# Clouditor Community Edition

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
docker run -v $HOME/.aws:/root/.aws -v $HOME/.azure:/root/.azure -p 9999:9999 clouditor/clouditor
```

Then open a web browser at http://localhost:9999.
