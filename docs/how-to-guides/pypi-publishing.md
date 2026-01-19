# Publishing to PyPI

Guide for publishing Vibium Python packages to PyPI.

## Packages

Six packages total:

| Package | Description |
|---------|-------------|
| `vibium` | Main package (clients/python) |
| `vibium-darwin-arm64` | macOS Apple Silicon binary |
| `vibium-darwin-x64` | macOS Intel binary |
| `vibium-linux-x64` | Linux x64 binary |
| `vibium-linux-arm64` | Linux ARM64 binary |
| `vibium-win32-x64` | Windows x64 binary |

## Version Bumping

Before publishing a new version, update the version across all packages:

```bash
make set-version VERSION=x.x.x
```

This updates all pyproject.toml files, `__init__.py` versions, and dependency constraints.

## Build Wheels

```bash
make package-python
```

This:
1. Cross-compiles clicker for all platforms
2. Copies binaries to platform packages
3. Creates `.venv-publish` with `twine` if it doesn't exist
4. Builds wheels for all 6 packages

Output:
```
packages/python/vibium_darwin_arm64/dist/vibium_darwin_arm64-0.1.0-py3-none-any.whl
packages/python/vibium_darwin_x64/dist/vibium_darwin_x64-0.1.0-py3-none-any.whl
packages/python/vibium_linux_x64/dist/vibium_linux_x64-0.1.0-py3-none-any.whl
packages/python/vibium_linux_arm64/dist/vibium_linux_arm64-0.1.0-py3-none-any.whl
packages/python/vibium_win32_x64/dist/vibium_win32_x64-0.1.0-py3-none-any.whl
clients/python/dist/vibium-0.1.0-py3-none-any.whl
```

## Test on TestPyPI (Recommended First Time)

```bash
# Activate the venv (twine is installed here)
source .venv-publish/bin/activate

# Upload to TestPyPI
twine upload --repository testpypi packages/python/*/dist/*.whl
twine upload --repository testpypi clients/python/dist/*.whl

# Test install from TestPyPI
pip install --index-url https://test.pypi.org/simple/ vibium
```

## Publish to PyPI

```bash
# Upload platform packages first
twine upload packages/python/*/dist/*.whl

# Upload main package
twine upload clients/python/dist/*.whl
```

## Authentication

### Option 1: API Token (Recommended)

1. Create token at https://pypi.org/manage/account/token/
2. Configure `~/.pypirc`:

```ini
[pypi]
username = __token__
password = pypi-xxxxx...
```

### Option 2: Environment Variables

```bash
export TWINE_USERNAME=__token__
export TWINE_PASSWORD=pypi-xxxxx...
twine upload dist/*.whl
```

### Option 3: Interactive

```bash
twine upload dist/*.whl
# Prompts for username and password
```

## Verify Published Package

```bash
pip install vibium
python3 -c "from vibium import browser_sync; print('OK')"
```
