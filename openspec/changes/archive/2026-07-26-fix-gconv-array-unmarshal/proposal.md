## Why

`gconv.Struct` currently attempts to call `reflect.Value.IsNil` for array-backed values while converting slice elements, causing requests containing values such as `[]uuid.UUID` to fail with `reflect: call of reflect.Value.IsNil on array Value`. Issue #4786 reproduces this behavior on the current `master` branch and requires a focused converter regression fix.

## What Changes

- Add regression coverage for converting string slices into slices of array-backed types that implement `encoding.TextUnmarshaler` through a pointer receiver.
- Route addressable array-backed values through the converter's existing common-interface unmarshalling path without calling `reflect.Value.IsNil` on arrays.
- Require bug fixes to document a test plan, demonstrate the regression test failing before production-code changes, and rerun the same test after the fix.
- Keep the change internal to conversion behavior with no new public API or dependency.

## Capabilities

### New Capabilities

- `gconv-array-unmarshal`: Defines conversion behavior for array-backed types that implement a supported unmarshalling interface, including values nested in slices.

### Modified Capabilities

- `code-quality-gates`: Requires an explicit red-green regression-test sequence for bug fixes before broader verification and review.

## Impact

- `util/gconv/internal/converter/converter_struct.go`
- `util/gconv/gconv_z_unit_issue_test.go`
- `openspec/specs/code-quality-gates/spec.md` after the change is archived
- No public API, module dependency, generated code, or external service impact
