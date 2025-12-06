## test-release

```bash
# Basic usage - build and test unsigned:
./scripts/test-release

# With options:
./scripts/test-release --sign       # Sign with Developer ID
./scripts/test-release --fresh      # Clear config for fresh install test
./scripts/test-release --pkg        # Open PKG installer instead of app
./scripts/test-release --notarize   # Sign and notarize

# Combine options:
./scripts/test-release --sign --fresh --pkg
```
