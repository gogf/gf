## ADDED Requirements

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
