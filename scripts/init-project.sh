#!/usr/bin/env bash
set -e

OLD_MODULE=$(go list -m 2>/dev/null || echo "github.com/bagusyanuar/go-be-template")
NEW_MODULE="$1"

if [ -z "$NEW_MODULE" ]; then
  echo "Usage: $0 <new-module-path>"
  echo "Example: $0 github.com/myusername/my-service"
  exit 1
fi

echo "🚀 Initializing project with module: $NEW_MODULE..."

# 1. Update module in all Go files
echo "📝 Updating Go import statements..."
if [[ "$OSTYPE" == "darwin"* ]]; then
  find . -type f -name "*.go" -exec sed -i '' "s|${OLD_MODULE}|${NEW_MODULE}|g" {} +
  sed -i '' "s|${OLD_MODULE}|${NEW_MODULE}|g" go.mod
else
  find . -type f -name "*.go" -exec sed -i "s|${OLD_MODULE}|${NEW_MODULE}|g" {} +
  sed -i "s|${OLD_MODULE}|${NEW_MODULE}|g" go.mod
fi

# 2. Setup .env from .env.example
if [ ! -f .env ] && [ -f .env.example ]; then
  echo "📄 Creating .env from .env.example..."
  cp .env.example .env
fi

# 3. Clean up and tidy dependencies
echo "📦 Running go mod tidy..."
go mod tidy

echo "✅ Project successfully initialized for $NEW_MODULE!"
echo "👉 You can now run 'make run' or 'docker-compose up' to start the application."
