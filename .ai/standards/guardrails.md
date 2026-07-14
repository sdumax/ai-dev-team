# Developer Guardrails

## Mandatory Workflow

Unless explicitly overridden with **"OVERRIDE WORKFLOW"** before the prompt, all work MUST follow this sequence:

```
Plan → Write Tickets → Pick Tickets to Build
```

### Step 1: Plan (Read-Only)
- Gather requirements (PM talks to user)
- Design system architecture (Architect)
- Define design requirements, accessibility, responsive breakpoints (UI/UX Designer)
- Document decisions in ADRs
- **No code changes allowed**

### Step 2: Write Tickets (Only Write in Plan Mode)
- Break work into tickets (Team Lead)
- Define dependencies and DAG layers
- Write ticket files to `docs/tickets/active/todo/T-NNNN.md`
- Create manifest at `.ai/tickets/shipments/ship-NNN/manifest.md`
- **This is the ONLY write operation allowed in plan mode**

### Step 3: Pick Tickets to Build (Build Mode)
- Select tickets from the todo list
- Move to implementation
- Execute tickets in dependency order
- All agents can work in parallel where possible

### Override
If the user states **"OVERRIDE WORKFLOW"** before their prompt, this default flow is bypassed.
The user takes control of the sequence.

---

## Coding Standards Detection

Before implementation begins, check for coding standards:

### Detection Order
1. `.ai/standards/coding-standards.md`
2. `docs/coding-standards.md`
3. `CONTRIBUTING.md` (look for standards/coding section)
4. `.editorconfig`
5. Linter configs (`.eslintrc`, `.prettierrc`, `golangci-lint.yml`, `.flake8`, etc.)

### If Standards Found
- All agents must follow the project's coding standards
- Reference the standards file in PR descriptions
- Reviewer checks compliance during code review
- Note the file location for reference

### If No Standards Found
**PAUSE AND WARN the user:**

> ⚠️ No coding standards found in this project.
>
> Proceeding without defined standards may result in:
> - Inconsistent code style
> - Architecture drift
> - Harder maintenance
>
> Options:
> 1. Create `.ai/standards/coding-standards.md` with your standards
> 2. Proceed without standards (go ahead)
> 3. Use built-in defaults

Wait for user confirmation before proceeding.

---

## Built-in Defaults

If user confirms "proceed without standards", use these sensible defaults:

### Universal Principles
- Keep functions small and focused (single responsibility)
- Use meaningful variable and function names
- Handle errors explicitly (no silent failures)
- Write tests for new functionality
- Follow existing patterns in the codebase

### Git Standards
- Conventional commits: `feat:`, `fix:`, `chore:`, `refactor:`, `test:`, `docs:`
- One logical change per commit
- Branch naming: `feature/`, `bugfix/`, `chore/` + description

### Code Review Checklist
- [ ] Code follows existing patterns
- [ ] Tests pass
- [ ] No security vulnerabilities introduced
- [ ] Documentation updated if needed
- [ ] No unrelated changes in the same PR

---

## Stack-Specific Standards (Optional Templates)

If the user wants to create standards, here are templates for common stacks:

### React / Next.js
- Component naming: PascalCase
- File structure: `components/`, `hooks/`, `utils/`, `lib/`
- State management: Context, Redux, Zustand, or Jotai
- Styling: CSS Modules, Tailwind, or styled-components
- Server Components vs Client Components clearly marked

### Go
- Package naming: lowercase, single word
- Error handling: explicit, wrapped with context
- Interface design: small, focused interfaces
- Test file naming: `*_test.go`
- Use `gofmt` and `golangci-lint`

### Python
- PEP 8 compliance
- Type hints on all public functions
- Docstrings: Google or NumPy style
- Test organization: `tests/` directory with `pytest`
- Use `ruff` or `flake8` for linting

### TypeScript / Node.js
- Strict mode enabled
- No `any` types
- Prefer `interface` over `type` for object shapes
- Use ESLint + Prettier
- ESM modules preferred

---

## Enforcement

### Plan Mode (Steps 1-2)
- Architect and Team Lead reference standards when designing
- UI/UX Designer checks design system compliance
- **Only ticket writing is allowed (no code changes)**

### Build Mode (Step 3)
- Developer follows standards during implementation
- Reviewer validates compliance during code review
- QA checks for style violations in test runs
- UI/UX Designer checks design and accessibility compliance
