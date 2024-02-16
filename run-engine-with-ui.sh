#/bin/bash

# Make sure the UI submodules are up-to-date
git submodule update --remote

# Build engine. This requires Go to be installed
go build ./cmd/engine/engine.go

# Install cl binary. 
go install ./cmd/cli/cl.go

# Build and run UI (in background). This requires nodejs to be installed
pushd ui
npm install
npm run dev &
popd

# Run engine with parameters optimized for the UI and an in-memory database.
./engine --db-in-memory \
    --dashboard-url=http://localhost:5173 \
    --api-cors-allowed-origins=http://localhost:5173  \
    --discovery-auto-start
