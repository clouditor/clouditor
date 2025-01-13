#/bin/bash

# Make sure the UI submodules are up-to-date
git submodule init
git submodule update --remote

# Build engine. This requires Go to be installed
go build ./cmd/engine/engine.go

# Install cl binary. 
go install ./cmd/cli/cl.go

# Build and run UI (in background). This requires Node.js to be installed
pushd ui
npm install --legacy-peer-deps

# Check if .env.local exists, if not, create it and create a secret
if [ ! -f .env.local ]; then
    npm install auth --legacy-peer-deps && npx auth secret
fi

npm run dev &
popd

# Run engine with parameters optimized for the UI and an in-memory database.
./engine --db-in-memory \
    --dashboard-callback-url=http://localhost:3000/api/auth/callback/confirmate \
    --api-cors-allowed-origins=http://localhost:3000  \
    --discovery-auto-start \
    $*
