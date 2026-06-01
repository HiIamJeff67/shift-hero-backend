[English](./CONTRIBUTING.md) | [繁體中文](./CONTRIBUTING.zh.md)

# Contributing to shift-hero

Thank you for contributing to **shift-hero**, a backend template architecture developed by **Notezy**.

## Repository Scope

These community guidelines apply to this repository itself (`shift-hero`), not to downstream projects created from this template.

## Ownership and License

- The architecture and implementation of this repository are developed and fully controlled by Notezy.
- Notezy publishes this repository under [`Apache-2.0`](../LICENSE).
- By submitting contributions, you agree your contribution can be distributed under Apache-2.0 in this repository.

## Development Requirements

1. `cp .env.example .env`
2. `go mod tidy`
3. `go build ./...`
4. Run tests relevant to your changes.

## Pull Request Workflow

1. Create a feature/fix branch from `main`.
2. Keep commits focused and logically grouped.
3. Open a PR using one of the templates under `.github/PULL_REQUEST_TEMPLATE/`.
4. Include:
   - Problem statement
   - Scope and design decisions
   - Validation evidence (build/test/log snippets)

## Quality Bar

- Keep naming and behavior template-neutral unless change is intentionally repo-specific.
- Avoid unrelated refactors in the same PR.
- Update docs when behavior or workflows change.
- Ensure `go build ./...` passes before requesting review.

## Security Disclosure

Do not disclose vulnerabilities publicly in issues/PRs.
Use the process in [`SECURITY.md`](./SECURITY.md).

## Contact

- Maintainer / repository contact: `thenotezy@gmail.com`
