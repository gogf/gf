## ADDED Requirements

### Requirement: Bug fixes demonstrate the regression before implementation

The project SHALL require a bug fix to define a focused regression-test plan and demonstrate that test failing for the reported defect before production code is modified.

#### Scenario: Regression test reproduces the reported defect

- **WHEN** a contributor prepares a bug fix
- **THEN** the contributor SHALL document the test input, expected behavior, expected pre-fix failure, and exact command, then run the added regression test against the unchanged implementation
- **AND** the test SHALL fail for the reported defect before production-code changes begin

#### Scenario: Regression test does not reproduce the reported defect

- **WHEN** the planned regression test passes before the fix or fails for an unrelated reason
- **THEN** implementation SHALL pause until the reproduction or test plan is corrected

### Requirement: Bug fixes verify the same regression test after implementation

The project SHALL require the exact regression test used for pre-fix reproduction to pass after the production-code fix before broader verification and review continue.

#### Scenario: Regression test passes after the fix

- **WHEN** the minimal production-code fix has been implemented
- **THEN** the contributor SHALL rerun the same regression-test command and confirm it passes before running package-level tests, race detection, lint, or broader checks
