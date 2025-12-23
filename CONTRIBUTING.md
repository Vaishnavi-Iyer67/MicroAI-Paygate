# Contributing to MicroAI Paygate

Thanks for considering a contribution! This project is open to issues, bug fixes, docs improvements, and new features that keep the system lean and secure.

## Getting Started
- Fork the repo and create a feature branch: `git checkout -b feature/my-change`
- Install prerequisites: Bun, Go 1.21+, Rust/Cargo, and Node (for Next.js).
- Copy `.env.example` to `.env` and fill in required keys (see README).

## Development Workflow
- Run the stack locally: `bun run stack`
- Run unit tests:
  - Gateway: `cd gateway && go test -v`
  - Verifier: `cd verifier && cargo test`
- Run E2E tests: `bun run test:e2e` (auto-starts gateway/verifier)

## Coding Standards
- Keep changes minimal and focused; avoid large, unrelated refactors.
- Add tests for new behavior; update existing tests if logic changes.
- Prefer clear, concise documentation alongside code changes.

## Commit Hygiene
- Use meaningful commit messages (e.g., `fix: ...`, `feat: ...`, `docs: ...`).
- Ensure `git status` is clean before opening a PR.

## Pull Requests
- Describe the problem, the solution, and testing performed.
- Link related issues if they exist.
- Be responsive to review feedback; small, incremental PRs are easier to merge.

## Reporting Issues
- Include steps to reproduce, expected vs actual behavior, logs, and environment details (OS, versions).

Thank you for helping improve MicroAI Paygate!
