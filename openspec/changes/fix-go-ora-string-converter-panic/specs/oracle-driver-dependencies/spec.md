## ADDED Requirements

### Requirement: Oracle driver dependency safety
The Oracle driver SHALL depend on a go-ora version whose string converter safely handles truncated multibyte input without panicking.

#### Scenario: GBK decode receives a trailing lead byte
- **WHEN** the go-ora GBK string converter decodes input whose last byte is greater than `0x80` and no following byte is available
- **THEN** decoding SHALL return without panicking
