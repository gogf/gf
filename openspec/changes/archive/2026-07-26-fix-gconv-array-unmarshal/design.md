## Context

`Converter.bindVarToReflectValue` groups `reflect.Array` with nil-capable kinds and calls `reflect.Value.IsNil` before checking whether the value implements `localinterface.ISet`. Arrays cannot be nil, so converting a textual slice element into an array-backed destination such as `uuid.UUID` enters the fallback binding path and returns a reflection error before the type's pointer `UnmarshalText` method can run.

The converter already provides `bindVarToReflectValueWithInterfaceCheck`, which handles addressable values implementing `IUnmarshalValue`, `IUnmarshalText`, `IUnmarshalJSON`, or `ISet`. The fix should reuse that path and preserve all unrelated conversion behavior.

## Goals / Non-Goals

**Goals:**

- Convert textual values into addressable array-backed types through existing supported unmarshalling interfaces.
- Support those types when they are elements of a destination slice.
- Add a focused regression test that fails on the current implementation and passes after the fix.
- Preserve existing behavior for ordinary slices, pointers, interfaces, structs, maps, and arrays without supported unmarshalling interfaces.

**Non-Goals:**

- Adding UUID-specific conversion logic.
- Adding a direct dependency on a UUID package for tests.
- Redesigning reflection conversion or public `gconv` APIs.
- Changing HTTP request parsing outside the generic conversion fix.

## Decisions

- Handle `reflect.Array` in its own branch before the nil-capable `reflect.Slice`, `reflect.Pointer`, and `reflect.Interface` branch. This avoids invalid `IsNil` calls while preserving the current fast path for nil-capable values.
- Delegate addressable array-backed values to `bindVarToReflectValueWithInterfaceCheck`. This reuses the converter's established interface precedence and allows pointer `UnmarshalText` implementations such as `uuid.UUID` without package-specific knowledge.
- Model the regression with a private array-backed test type whose pointer implements `UnmarshalText`. This reproduces the same reflection shape as `uuid.UUID`, avoids dependency metadata changes, and directly tests the generic contract.
- Follow an explicit red-green sequence: add only `Test_Issue4786`, run the targeted command and record the current reflection failure, then modify production code and rerun the same command.

Alternatives considered:

- Removing `reflect.Array` from the existing switch would stop the immediate `IsNil` error but would not invoke the array type's unmarshalling interface.
- Special-casing `uuid.UUID` would add an unnecessary dependency and leave other array-backed unmarshalling types broken.
- Calling the common-interface helper for every value would broaden the hot-path behavior beyond the reported defect.

## Risks / Trade-offs

- [Interface precedence changes for array-backed values] -> Limit the new common-interface check to `reflect.Array` and verify existing `util/gconv` tests with race detection.
- [Regression test does not import the exact UUID package] -> Use the same relevant type properties: an array underlying type, addressability, a pointer `UnmarshalText` implementation, and conversion from a string slice.
- [Fallback behavior for arrays without supported interfaces] -> Leave the existing reflection fallback unchanged and scope the test to the reported interface-backed case.

## Migration Plan

No migration is required. The change fixes internal conversion behavior without changing public signatures or persisted data. Reverting the production-code branch and regression test restores the previous behavior if an unexpected compatibility issue appears.

## Open Questions

None.

## Verification Evidence

### Red: Regression Reproduced Before Production Changes

- Command: `go test -count=1 -run '^Test_Issue4786$' ./util/gconv`
- Result: failed as planned on the unchanged converter.
- Failure: `reflect: call of reflect.Value.IsNil on array Value`.
- Stack location: `util/gconv/internal/converter/converter_struct.go:505` in `bindVarToReflectValue`.
- Production-code state: `converter_struct.go` had no diff when the failing test was run.

### Green: Same Regression Test Passed After the Fix

- Command: `go test -count=1 -run '^Test_Issue4786$' ./util/gconv`
- Result: passed after routing addressable arrays through the existing common-interface unmarshalling helper.

### Broader Verification

- `go test -count=1 -race ./util/gconv/...`: passed.
- Direct `[]uuid.UUID` reproduction: returned `err=<nil>` and the expected UUID.
- Targeted coverage with `-coverpkg=./util/gconv/...`: the new array branch and return path each executed once.
- `make lint`: passed with `0 issues` after installing the repository-documented `golangci-lint/v2` tool.
- `openspec validate fix-gconv-array-unmarshal --strict`: passed.

### GF Review

- Scope source: `git status --short`, expanded untracked files, the source diff, and all active change artifacts.
- Backend code review: passed; the fix is internal, generic, minimal, and introduces no dependency or public API change.
- Project spec review: passed; implementation and red-green evidence match the proposal, design, and both incremental specs.
- SQL review: not applicable; no SQL files or snippets changed.
- Unit test review: passed; `Test_Issue4786` directly covers the changed array-backed unmarshalling path and follows the package's issue-test convention.
- Findings: 0 critical, 0 warnings.
