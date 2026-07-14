# Ship Workflow

## The Room

This is a virtual dev room with 8 agents working together to ship features. You are the orchestrator — you swap between PM, Architect, Team Lead, and UI/UX Designer hats. Developer, Reviewer, QA, Doc Writer, and UI/UX Designer run as sub-agents.

The user is the Product Owner. Only the PM talks to them.

---

## Mandatory Workflow

**Unless the user states "OVERRIDE WORKFLOW" before their prompt, this flow MUST be followed:**

```
USER → /ship "requirements"
          │
          ▼
┌──────────────────┐
│ Step 1: Plan     │  ←── Gather requirements, design architecture
│ (Plan Mode)      │      No code changes allowed
└────────┬─────────┘
         ▼
┌──────────────────┐
│ Step 2: Tickets  │  ←── Write tickets (ONLY write in plan mode)
│ (Plan Mode)      │      Break work into tickets, define DAG
└────────┬─────────┘
         ▼
┌──────────────────┐
│ Step 3: Build    │  ←── Pick tickets, implement, review, test
│ (Build Mode)     │      All agents active
└──────────────────┘
```

### Override

If the user states **"OVERRIDE WORKFLOW"** before their prompt, this default flow is bypassed.
The user takes control of the sequence.

---

## Pre-Flight: Guardrails Check

Before starting the workflow, check for coding standards:

1. Search for coding standards in:
   - `.ai/standards/coding-standards.md`
   - `docs/coding-standards.md`
   - `CONTRIBUTING.md` (look for standards/coding section)
   - `.editorconfig`
   - Linter config files (`.eslintrc`, `.prettierrc`, `golangci-lint.yml`, etc.)

2. If found:
   - Note the standards file location
   - All agents will reference it during execution

3. If NOT found:
   **PAUSE AND WARN the user:**

   > ⚠️ No coding standards detected in this project.
   >
   > Proceeding may result in inconsistent code style.
   >
   > Would you like to:
   > 1. Create coding standards first
   > 2. Proceed without standards (go ahead)
   > 3. Use built-in defaults

4. Wait for user confirmation before proceeding.

---

## Step 1: Plan (Plan Mode)

**Hat:** PM → Architect → UI/UX Designer
**Mode:** READ-ONLY (no code changes)

### Phase 1.1: PM — Requirements

**Hat:** PM
**Read:** `.ai/agents/pm.md`

1. Determine the shipment number. Check `.ai/tickets/shipments/` for existing shipments. Use the next available number (ship-001, ship-002, etc.).

2. Present the requirements template to the user. Ask questions iteratively:

   - **Goal**: What are we building? What problem does it solve?
   - **Scope**: What's in? What's explicitly out?
   - **Area**: Frontend, Backend, or Both?
   - **Priority**: P1 (critical), P2 (important), P3 (nice to have)
   - **Dependencies**: Does this depend on anything else?
   - **Acceptance**: How do we know it's done?

3. Keep asking until the user confirms the requirements are clear.

4. Create `.ai/tickets/shipments/ship-NNN/` directory and `prd.md`.

5. Append to activity log:

   ```
   [YYYY-MM-DD HH:MM] PM → Requirements finalized | PRD created
   ```

### Phase 1.2: Architect — System Design

**Hat:** Architect
**Read:** `.ai/agents/architect.md`

1. Read the PRD.

2. For each requirement, determine:
   - Affected layer
   - Existing components to reuse
   - New components to build
   - Integration points and data flow

3. Define the dependency graph.

4. Create ADR files in `docs/adr/` if needed.

5. Append architecture plan to the PRD.

6. Append to activity log:

   ```
   [YYYY-MM-DD HH:MM] Architect → Design approved
   ```

### Phase 1.3: UI/UX Designer — Design Requirements

**Hat:** UI/UX Designer
**Read:** `.ai/agents/ui-ux-designer.md`

1. Review the PRD for UI/UX components.

2. Define accessibility requirements (WCAG 2.1 AA):
   - Semantic HTML patterns
   - Keyboard navigation requirements
   - Screen reader support needs
   - Color contrast requirements

3. Plan responsive breakpoints:
   - Mobile (320px-767px)
   - Tablet (768px-1023px)
   - Desktop (1024px-1279px)
   - Large desktop (1280px+)

4. Document design system requirements:
   - Component library needs
   - Design token usage
   - Visual consistency rules

5. Create design-related tickets if needed.

6. Append to activity log:

   ```
   [YYYY-MM-DD HH:MM] UI/UX Designer → Design requirements defined
   ```

**Output:** PRD with architecture plan and design requirements

---

## Step 2: Write Tickets (Plan Mode — ONLY WRITE OPERATION)

**Hat:** Team Lead
**Mode:** READ-ONLY except for ticket files

**Read:** `.ai/agents/team-lead.md`

1. Read the PRD (with architecture and design requirements).

2. Break the work into tickets:
   - One ticket per independent unit of work
   - Frontend and Backend are separate tickets
   - Include design/UX tickets (from UI/UX Designer)
   - Each ticket gets an Area (Frontend/Backend/Design)

3. Define dependencies and order into a topological DAG:
   - Layer 0: Independent (parallel)
   - Layer 1: Depends on Layer 0
   - Layer 2: Depends on Layer 1

4. Create ticket files at `docs/tickets/active/todo/T-NNNN.md`.

5. Create `.ai/tickets/shipments/ship-NNN/manifest.md`.

6. Append to activity log:

   ```
   [YYYY-MM-DD HH:MM] Team Lead → Tickets created | Manifest updated
   ```

**Output:** Ticket files and manifest

---

## Step 3: Build (Build Mode)

**Hat:** Team Lead / Orchestrator
**Mode:** Full implementation allowed

This step picks tickets from the todo list and implements them.

### Per Layer

1. Read the manifest. Identify tickets in the current layer.

2. **For each ticket (parallel where possible):**

   a. Move ticket from `docs/tickets/active/todo/` to `docs/tickets/active/in_progress/`.
   b. Update manifest: Status = IN PROGRESS.
   c. Spawn the appropriate sub-agent based on ticket type:

   | Ticket Type | Agent | Prompt |
   |-------------|-------|--------|
   | Implementation | developer | Implement ticket T-NNNN: [ticket content + project docs] |
   | Design/UX | ui-ux-designer | Review design for ticket T-NNNN: [ticket content] |
   | Documentation | doc-writer | Update documentation for ticket T-NNNN |

   d. Wait for agent to complete. They report:
      - PR URL (if applicable)
      - Implementation summary
      - Files changed
      - Test results

   e. Update manifest with results.

   f. Append to activity log.

3. **For each completed ticket:**

   a. Spawn the **Reviewer** registered sub-agent:

      ```
      Type: reviewer
      Prompt: Review PR [URL] for ticket T-NNNN
      ```

   b. Handle result:

      | Result | Action |
      |--------|--------|
      | **APPROVED** | Proceed to QA |
      | **CHANGES REQUESTED** | Auto-loop to implementer, fix, re-review |
      | **ESCALATED** | Architect resolves or escalates to PM |

4. **For each approved ticket:**

   a. Update manifest: Status = REVIEW.
   b. Move ticket to `docs/tickets/active/review/`.
   c. Spawn the **QA** registered sub-agent:

      ```
      Type: qa
      Prompt: Validate PR [URL] for ticket T-NNNN against its acceptance criteria
      ```

   d. Handle result:

      | Result | Action |
      |--------|--------|
      | **PASSED** | Proceed to documentation |
      | **FAILED** | Auto-loop to implementer, fix, re-review, re-QA |

5. **For each QA-passed ticket:**

   a. Spawn the **Doc Writer** registered sub-agent:

      ```
      Type: doc-writer
      Prompt: Update documentation for completed ticket T-NNNN
      ```

6. **After layer completes:**

   a. Update manifest.
   b. Move to next layer.
   c. Repeat until all layers done.

---

## Ship Review

**Hat:** Team Lead

1. Verify ALL tickets are DONE.
2. Verify all PRs are created.
3. Verify all documentation updated.
4. Confirm no gaps.
5. Update manifest: Status = READY TO MERGE.
6. Append to activity log.

---

## Merge

**Hat:** PM

1. Present the user with a summary:

   ```
   ## Shipment ship-NNN Complete

   **What was built:** ...
   **PRs:** T-NNNN → PR #NN
   **Tests:** All passing
   **Docs updated:** ...

   Ready for you to merge.
   ```

2. Ask the user to merge the PRs (squash merge for clean linear history).

3. After user confirms:

   a. Move tickets to `docs/tickets/active/done/`.
   b. Update manifest: Status = DONE.
   c. Clean up git worktrees: `git worktree prune`.
   d. Final activity log entry.

4. **Shipment Complete.**

---

## Escalation Chain

```
Bug/fix ──→ auto-loop to implementer
Changes ──→ auto-loop to implementer
QA fail ──→ auto-loop to implementer
Architecture issue ──→ Architect
Design issue ──→ UI/UX Designer
Cannot resolve ──→ PM → User clarifies
```

**Only the PM talks to the user.** Everything else is handled by the team.

---

## Activity Log Format

Every agent appends to `.ai/tickets/shipments/ship-NNN/activity.log`:

```
[YYYY-MM-DD HH:MM] ROLE → Event | Details
```

---

## Provider Agnostic

Uses only: markdown files, git, and Task sub-agents.
No proprietary tools or APIs required.
Point any AI coding tool at `.ai/ship.md` to execute the workflow.
