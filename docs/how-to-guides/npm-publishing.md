# How to Publish Vibium to npm

## Prerequisites

- npm account with access to `vibium` package
- Member of `@vibium` org on npm
- Logged in (`npm login`)

## Version Bumping

Before publishing a new version, update the version across all packages:

```bash
make set-version VERSION=x.x.x
```

This updates all package.json files, optionalDependencies, Python packages, and regenerates package-lock.json.

## Build, Test & Package

```bash
make build
make test
make package
```

## Local Testing (Before Publishing)

Always test locally before publishing:

```bash
# Pack the main package
cd packages/vibium && npm pack

# Test in a fresh directory
mkdir /tmp/vibium-test && cd /tmp/vibium-test
npm init -y
npm install /path/to/vibium/packages/vibium/vibium-*.tgz

# Verify it works
node -e "const { browser } = require('vibium'); console.log('OK')"
npx vibium  # Should start MCP server
```

## Publishing

**Important:** Publish platform packages first, then main package.

```bash
# Platform packages (all must succeed before publishing main)
(cd packages/linux-x64 && npm publish --access public)
(cd packages/linux-arm64 && npm publish --access public)
(cd packages/darwin-x64 && npm publish --access public)
(cd packages/darwin-arm64 && npm publish --access public)
(cd packages/win32-x64 && npm publish --access public)

# Main package (after all platform packages are live)
(cd packages/vibium && npm publish)
```

## Troubleshooting

### "You must be logged in to publish"
```bash
npm login
npm whoami  # Verify you're logged in
```

### "Package name too similar to existing package"
The `@vibium` scope prevents this. Ensure you're using scoped names.

### "Cannot publish over previously published version"
Bump the version number. npm doesn't allow republishing the same version.
