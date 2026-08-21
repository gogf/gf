## Why

The Oracle driver currently depends on `github.com/sijms/go-ora/v2 v2.7.10`, whose GBK string converter can panic while decoding malformed or truncated multibyte data. A previously applied manual source edit targeted the vendored `StringConverter.Decode` implementation, but this repository does not commit that vendor path, so the fix would not reliably reach the mainline module build.

## What Changes

- Upgrade the Oracle driver's go-ora dependency to the smallest upstream version that includes the `StringConverter.Decode` bounds check fix.
- Add a focused regression test that exercises the GBK converter with a trailing lead byte and verifies decoding does not panic.
- Document the root cause so future maintainers do not reintroduce a manual vendor-source workaround.

## Capabilities

### New Capabilities
- `oracle-driver-dependencies`: Keeps the Oracle driver on a go-ora version that handles truncated GBK decode input without panicking.

### Modified Capabilities
- None.

## Impact

- `contrib/drivers/oracle/go.mod`
- `contrib/drivers/oracle/go.sum`
- Oracle driver dependency verification tests
