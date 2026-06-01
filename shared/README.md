# Shared Layer

This directory contains cross-cutting code reused by multiple app layers.

## Structure

- `shared/constants/`: project-level constants (rate limits, versions, tokens, URLs).
- `shared/lib/`: reusable utility packages with narrow responsibilities.
- `shared/types/`: shared types used across app packages.

## Design Rules

- Keep packages framework-light and dependency-minimal.
- Prefer pure helper logic over business/domain logic.
- If a package is no longer imported anywhere, remove it to keep template surface small.
