## 1. Plan and Reproduce

- [x] 1.1 Record the regression-test plan, expected pre-fix failure, exact command, and generic array-backed test model before editing production code
- [x] 1.2 Add only `Test_Issue4786` and its private array-backed `UnmarshalText` test type to `util/gconv/gconv_z_unit_issue_test.go`
- [x] 1.3 Run the targeted test against the unchanged converter and record evidence that it fails because `reflect.Value.IsNil` is called on an array

## 2. Implement the Minimal Fix

- [x] 2.1 Route addressable `reflect.Array` values through the existing common-interface unmarshalling helper while leaving nil-capable kinds on their current path
- [x] 2.2 Rerun the exact targeted regression test and confirm it passes

## 3. Verify and Review

- [x] 3.1 Run `go test -count=1 -race ./util/gconv/...`
- [x] 3.2 Run relevant diagnostics and repository lint, resolving findings within the change scope
- [x] 3.3 Record red-green evidence, verify OpenSpec compliance, and run the required `gf-review`
- [x] 3.4 Confirm the intended PR scope excludes local agent and Trellis metadata
