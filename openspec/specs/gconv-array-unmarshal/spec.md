# gconv-array-unmarshal Specification

## Purpose

Define how `gconv` converts textual values into array-backed destination types through supported unmarshalling interfaces.

## Requirements

### Requirement: Array-backed values use supported unmarshalling interfaces

`gconv` SHALL convert textual source values into addressable array-backed destination values when a pointer to the destination type implements a supported unmarshalling interface.

#### Scenario: Convert a string slice into array-backed unmarshalling values

- **WHEN** `gconv.Struct` binds a source string slice to a destination slice whose array-backed element type implements `UnmarshalText` through a pointer receiver
- **THEN** conversion SHALL return no error, invoke the unmarshalling behavior for each element, and preserve the source element order

#### Scenario: Avoid nil reflection operations on arrays

- **WHEN** the converter evaluates an array-backed destination value
- **THEN** it SHALL NOT invoke reflection operations that are valid only for nil-capable kinds
