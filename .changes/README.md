# Changelog Entries

This directory contains individual changelog entries that will be used to generate the main CHANGELOG.md.

## Format

Each file should contain one or more blocks in this format:

```
```release-note:TYPE
DESCRIPTION
```
```

Where TYPE can be:
- `feature`: New functionality
- `enhancement`: Improvement to existing functionality  
- `bug`: Bug fixes
- `breaking-change`: Breaking changes
- `deprecation`: Deprecated functionality
- `security`: Security-related changes
- `note`: General notes

## File Naming

- For releases: `vX.Y.Z.txt` (e.g., `v1.0.0.txt`)
- For PRs: `PR-123.txt` where 123 is the PR number
- For issues: `issue-456.txt` where 456 is the issue number

## Automation

This structure allows for:
- Conflict-free changelog management
- Automated changelog generation
- Integration with CI/CD pipelines
- Consistent formatting across releases