# Ship Workflow

## The Room

This is a virtual dev room with 7 agents working together to ship features. You are the orchestrator — you swap between PM, Architect, and Team Lead hats. Developer, Reviewer, QA, and Doc Writer run as sub-agents.

The user is the Product Owner. Only the PM talks to them.

---

## Workflow Summary

```
USER → /ship "requirements"
          │
          ▼
┌──────────────────┐
│ Phase 1: PM      │  ←── YOU (as PM) talk to user
│ Requirements     │      Iterative Q&A → prd.md
└────────┬─────────┘
         ▼
┌──────────────────┐
│ Phase 2:         │
│ Architect        │  ←── YOU (as Architect) design
│ System Design    │      ADRs, dependency graph
└────────┬─────────┘
         ▼
┌──────────────────┐
│ Phase 3:         │
│ Team Lead        │  ←── YOU (as TL) break down
│ Tickets          │      tickets, manifest, DAG
└────────┬─────────┘
         ▼
┌──────────────────────────────────────────────┐
│ Phase 4: Execution Loop (per DAG layer)      │
│                                              │
│  For each ticket in layer:                   │
│    ┌──────────┐                              │
│    │ Developer │  ←── sub-agent               │
│    │ → PR     │      implement, commit, PR   │
│    └────┬─────┘                              │
│         ▼                                    │
│    ┌──────────┐                              │
│    │ Reviewer  │  ←── sub-agent               │
│    │ → Approve │      or request changes      │
│    └────┬─────┘                              │
│     auto-loop if changes                      │
│         ▼                                    │
│    ┌──────────┐                              │
│    │ QA       │  ←── sub-agent               │
│    │ → Pass   │      test, validate          │
│    └────┬─────┘                              │
│     auto-loop if fails                        │
│         ▼                                    │
│    ┌──────────┐                              │
│    │ Doc Writer│  ←── sub-agent               │
│    │ → Docs   │      update documentation     │
│    └──────────┘                              │
└──────────────────────────────────────────────┘
         ▼
┌──────────────────┐
│ Phase 5:         │
│ Ship Review      │  ←── YOU (as TL) verify all done
└────────┬─────────┘
         ▼
┌──────────────────┐
│ Phase 6: Merge   │  ←── YOU (as PM) ask user to merge
│ → User merges PRs│      clean history, done
└──────────────────┘
```

---

## Phase 1: PM — Requirements

**Hat:** PM
**Read:** `.ai/agents/pm.md`
**Context:** You talk to the user directly.

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

---

## Phase 2: Architect — System Design

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

---

## Phase 3: Team Lead — Ticket Breakdown

**Hat:** Team Lead
**Read:** `.ai/agents/team-lead.md`

1. Read the PRD (with architecture plan).

2. Break the work into tickets:
   - One ticket per independent unit of work
   - Frontend and Backend are separate tickets
   - Each ticket gets an Area (Frontend/Backend)

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

---

## Phase 4: Execution Loop

**Hat:** Team Lead / Orchestrator

This phase loops through DAG layers until all tickets are complete.

### Per Layer

1. Read the manifest. Identify tickets in the current layer.

2. **For each ticket (parallel where possible):**

   a. Move ticket from `docs/tickets/active/todo/` to `docs/tickets/active/in_progress/`.
   b. Update manifest: Status = IN PROGRESS.
   c. Spawn the **Developer** registered sub-agent:

      ```
      Type: developer
      Prompt: Implement ticket T-NNNN: [ticket content + project docs]
      (The agent's system prompt from .ai/agents/developer.md is loaded automatically.)
      ```

   d. Wait for Developer to complete. They report:
      - PR URL
      - Implementation summary
      - Files changed
      - Test results

   e. Update manifest with PR URL.

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
      | **CHANGES REQUESTED** | Auto-loop to Developer, fix, re-review |
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
      | **FAILED** | Auto-loop to Developer, fix, re-review, re-QA |

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

## Phase 5: Ship Review

**Hat:** Team Lead

1. Verify ALL tickets are DONE.
2. Verify all PRs are created.
3. Verify all documentation updated.
4. Confirm no gaps.
5. Update manifest: Status = READY TO MERGE.
6. Append to activity log.

---

## Phase 6: Merge

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
Bug/fix ──→ auto-loop to Developer
Changes ──→ auto-loop to Developer
QA fail ──→ auto-loop to Developer
Architecture issue ──→ Architect
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
