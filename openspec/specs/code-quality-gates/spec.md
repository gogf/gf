# code-quality-gates Specification

## Purpose

Define mandatory unit-test, new-code coverage, and red-green regression verification gates for submitted code changes.

## Requirements

### Requirement: Submitted code changes include unit tests

The project SHALL require every submitted code change to include unit tests that directly cover the affected logic and expected behavior of the changed code path in the affected package.

#### Scenario: Behavior-changing code is submitted

- **WHEN** a contribution adds or changes code in the repository
- **THEN** the submission SHALL include unit tests that exercise the affected logic or preserve the expected behavior of the changed code path before the change is considered complete

### Requirement: Newly added code meets the coverage baseline

The project SHALL require the newly added code in a submission to maintain unit-test coverage of at least 80%, and reviews SHALL treat 90% or above as the preferred target when that level is feasible without artificial or low-value tests.

#### Scenario: Coverage falls below the minimum

- **WHEN** the newly added code in a submission is covered below 80%
- **THEN** the change SHALL not satisfy the project quality requirement

#### Scenario: Coverage meets the minimum baseline

- **WHEN** the newly added code in a submission reaches 80% or higher coverage
- **THEN** the change SHALL satisfy the minimum coverage requirement

#### Scenario: Coverage reaches the preferred target

- **WHEN** the newly added code in a submission reaches 90% or higher coverage
- **THEN** the change SHALL satisfy the preferred coverage target for the project

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
