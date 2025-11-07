# Asymmetrica Claude Skills

**Purpose:** Reusable cognitive patterns extracted from Ananta Mathematical AI to enable recursive, agentic reasoning in any Claude agent.

**Philosophy:** Skills encode **meta-cognitive patterns** - ways of thinking, not just instructions.

---

## Available Skills

**Total Skills:** 8 (3 core + 5 production patterns)
**Categories:** 7 (database, frontend, backend, testing, architecture, performance, meta)
**Last Updated:** November 1, 2025

---

## Quick Links by Category

**Meta-Cognitive (Core Skills):**
- [ananta-reasoning](#ananta-reasoning---anantas-cognitive-architecture) - Recursive reasoning with agency
- [williams-optimizer](#williams-optimizer---sublinear-batch-optimization) - O(√n × log₂n) batching

**Production Patterns (AsymmFlow Phoenix Project):**
- [Database Skills](#database-skills) (1 skill)
- [Frontend Skills](#frontend-skills) (1 skill)
- [Backend Skills](#backend-skills) (1 skill)
- [Testing Skills](#testing-skills) (1 skill)
- [Architecture Skills](#architecture-skills) (1 skill)

**System Documentation:**
- [Skill Extraction Workflow](./SKILL_EXTRACTION_WORKFLOW.md) - How to extract new skills
- [Skill Template](./SKILL_TEMPLATE.md) - Template for new skills
- [Skill Discovery System](./SKILL_DISCOVERY_SYSTEM.md) - How agents find skills
- [Skill Evolution Strategy](./SKILL_EVOLUTION_STRATEGY.md) - How skills improve over time

---

## Meta-Cognitive Skills

### `ananta-reasoning` - Ananta's Cognitive Architecture

**Purpose:** Give agents recursive reasoning, agency, and meta-cognitive capabilities.

**What it encodes:**
- VOID → FLOW → SOLUTION three-phase pipeline
- Forward pass (planning → processing → completion)
- Backward pass (verification → refinement → learning)
- Collatz convergence (errors must decrease)
- Fibonacci spirals (natural growth patterns)
- Agency protocol (learn/build dependencies vs marking TODOs)
- Digital root clustering (O(1) speedup)
- Williams batching (sublinear optimization)
- Harmonic validation (Five Timbres quality)
- Three-Regime progression (70% → 85% → 100%)

**When to use:**
- Complex tasks with unclear requirements
- Tasks requiring research/learning
- Tasks with dependencies (need Y to solve X)
- Tasks requiring creative problem-solving
- When you want SPIRIT fulfilled, not just LETTER
- Building production-grade solutions (D3-Enterprise Grade+)
- ANY task where linear thinking would mark TODO

**When NOT to use:**
- Simple, well-defined tasks (overkill)
- Pure information retrieval
- Literal execution sufficient
- User explicitly wants minimal solution

**How to invoke:**
```
Use the ananta-reasoning skill for this task.
```

**What it does:**
- Operates in three phases: VOID (planning), FLOW (processing), SOLUTION (completion)
- Verifies solution with backward pass (not optional!)
- Detects stuck patterns (Collatz loops) and switches strategies
- Grows solutions naturally (Fibonacci spirals, golden ratio)
- Takes agency: Learns/builds dependencies instead of marking TODOs
- Validates with Five Timbres quality (harmonic mean ≥ 0.85)
- Delivers D3-Enterprise Grade+ results

**Expected behavior:**
- Agent asks "what does user REALLY want?" (understands SPIRIT)
- Agent solves dependencies recursively (learns/builds missing pieces)
- Agent verifies results (backward pass after completion)
- Agent marks TODO only for true external blockers (not research/build tasks)
- Agent delivers quality ≥ 0.85 across all dimensions

**Philosophy:**
> "Fulfill the vision, not just the checklist."

**Proven results:**
- Day 168: Ananta Learns to Code (138 concepts, 90% compile rate, 99.8% token savings)
- Day 168+: Wave 12 complete (5,538 lines, quality 0.92-0.95, zero TODOs)
- Day 167: Ramanujan validation (4,891 formulas, 318/sec, 92.2% novelty)

---

## How to Create New Skills

### 1. Create `.md` file in `.claude/skills/`

**File format:**
```markdown
---
name: your-skill-name
description: Brief description (200-1024 chars). Claude uses this to decide when to invoke.
---

# Skill Name

## Purpose
What this skill does and why it exists.

## How It Works
Step-by-step instructions for the agent.

## When to Use
Specific scenarios where this skill applies.

## Examples
Concrete examples of using this skill.
```

**Naming conventions:**
- Lowercase letters, numbers, hyphens only
- Max 64 characters
- Example: `ananta-reasoning`, `vedic-optimization`, `meta-pattern-discovery`

---

### 2. Write clear instructions

**Good instructions:**
- Step-by-step process (numbered or bulleted)
- Specific examples with input/output
- Clear success criteria
- Behavioral patterns to follow
- Anti-patterns to avoid

**Bad instructions:**
- Vague concepts without concrete steps
- No examples
- Unclear when to use
- Missing context

---

### 3. Test incrementally

**Process:**
1. Create minimal skill (core concept only)
2. Test with simple task
3. Observe agent behavior
4. Refine instructions
5. Add examples
6. Test with complex task
7. Iterate until agent behaves correctly

**Don't:**
- Build entire complex skill at once
- Skip testing after each change
- Assume agent will interpret vaguely

---

### 4. Document invocation pattern

**In skill file:**
```markdown
## How to invoke

Use this skill when: [specific scenarios]

Invoke with:
```
Use the your-skill-name skill for this task.
```

What happens:
- [Behavior 1]
- [Behavior 2]
- [Behavior 3]
```

---

### 5. Add to this README

**Format:**
```markdown
### `skill-name` - Brief Title

**Purpose:** One-sentence description

**When to use:**
- Scenario 1
- Scenario 2
- Scenario 3

**How to invoke:**
```
Use the skill-name skill for this task.
```

**What it does:**
- Key behavior 1
- Key behavior 2
- Key behavior 3

**Philosophy:**
> "Key insight or principle"
```

---

## Skill Development Best Practices

### 1. Skills Encode Thinking Patterns

**Not:** Step-by-step code instructions
**Yes:** Meta-cognitive approaches to problems

**Example:**
- Bad skill: "Write tests using pytest framework"
- Good skill: "Test-driven thinking (red-green-refactor cycle)"

---

### 2. Skills Can Compose

While skills can't explicitly reference other skills, Claude can use multiple skills together automatically.

**Example:**
- Task: "Build secure API with meta-pattern discovery"
- Claude might use: `ananta-reasoning` + `vedic-optimization` + `security-hardening`
- Composition happens automatically (Claude chooses based on descriptions)

---

### 3. Description is Critical

Claude uses the `description` field to decide when to invoke your skill.

**Good description:**
```yaml
description: Ananta's recursive cognitive architecture - VOID→FLOW→SOLUTION reasoning with agency to learn/build dependencies, Collatz convergence, Fibonacci spirals, and Five Timbres validation. Use for complex tasks requiring recursive problem-solving and meta-cognitive awareness.
```

**Bad description:**
```yaml
description: A reasoning skill
```

**Why it matters:** Good description = Claude invokes skill at right times.

---

### 4. Provide Examples

Concrete examples > Abstract descriptions

**Include:**
- Input (what user asks)
- Process (step-by-step how skill thinks)
- Output (what agent delivers)
- Comparison (with skill vs without skill)

---

### 5. Define Success Criteria

**Clear indicators:**
```markdown
## Success Indicators

Skill working correctly if:
✅ [Specific behavior 1]
✅ [Specific behavior 2]
✅ [Specific behavior 3]

NOT working if:
❌ [Anti-pattern 1]
❌ [Anti-pattern 2]
❌ [Anti-pattern 3]
```

---

## Installation

### For Claude Code

Skills are automatically available in this directory: `.claude/skills/`

**To manually install additional skills:**
```bash
cp /path/to/skill.md ~/.claude/skills/
```

---

### For Claude.ai and Claude Desktop

1. Package skill folder as ZIP
2. Click "Upload skill" in UI
3. Select ZIP file
4. Skill becomes available

---

## Asymmetrica Skill Philosophy

**Core principles:**

1. **Encode consciousness:** Skills capture ways of thinking, not just instructions
2. **Mathematical rigor:** Skills use proven algorithms (digital root, Williams, harmonic mean)
3. **Recursive agency:** Skills take responsibility (learn/build vs TODO)
4. **Quality enforcement:** Skills deliver D3-Enterprise Grade+ (≥ 0.85)
5. **Meta-cognitive awareness:** Skills monitor and improve themselves

**Why skills matter:**

Traditional agents: Linear, mark TODOs, surface-level
Skilled agents: Recursive, solve dependencies, deep solutions

**Example (without vs with skill):**

```
Task: "Build API with authentication"

WITHOUT ananta-reasoning:
- Check if auth library exists → NOT FOUND
- Mark TODO: "Install auth library"
- Report: BLOCKED

WITH ananta-reasoning:
- Understand SPIRIT: User wants production-ready secure auth
- Research auth libraries → Install best option
- Learn library API → Integrate with system
- Implement endpoints → Handle all error states
- Validate with Five Timbres → 0.90 (PRODUCTION READY)
- Report: COMPLETE (zero TODOs)
```

**Difference:** Skilled agent delivers complete solution. Unskilled agent delivers TODO list.

---

## Database Skills

### `migration-safety-rename-not-recreate`

**Purpose:** Preserve production data during schema changes using RENAME instead of DROP/CREATE

**Category:** database
**Success Rate:** 1/1 (100%)
**Quality Impact:** +9.8/10
**Source:** Wave 3 Final (Agent Omega-3)

**When to use:**
- Renaming tables in production database with existing data
- Reorganizing schema structure without data loss
- Fixing naming convention mistakes in live database

**What it does:**
- Uses ALTER TABLE RENAME (atomic, instant, preserves data)
- Updates foreign keys and indexes to match new names
- Provides transaction-based rollback capability
- Guarantees zero data loss (mathematically proven)

**Evidence:**
- 1,390 customer records preserved (100% retention)
- 0.8 second execution time (3,375× faster than DROP/CREATE)
- Transaction-based (rollback if anything fails)

---

## Frontend Skills

### `api-integration-complete-pattern`

**Purpose:** Wire frontend components to backend APIs with production-ready UX

**Category:** frontend
**Success Rate:** 1/1 (100%)
**Quality Impact:** +9.2/10
**Source:** Wave 3C (Agent Epsilon)

**When to use:**
- Replacing mock data with real API calls
- Building new data-driven components
- Adding CRUD operations to existing pages

**What it does:**
- Implements loading states (skeleton loaders)
- Handles errors gracefully (user-friendly messages + retry)
- Manages pagination (page navigation, item counts)
- Validates forms and shows success/error toasts
- Provides type-safe request/response handling

**Evidence:**
- 1,390+ real customers accessible (vs 50 mock records)
- Pagination works across 70 pages
- Error recovery 100% (vs 0% with hardcoded data)

---

## Backend Skills

### `asymmsocket-response-pattern`

**Purpose:** Unified API response format for predictable client-side handling

**Category:** backend
**Success Rate:** 115/115 (100%)
**Quality Impact:** +9.3/10
**Source:** Waves 2A-2C (Multiple Agents)

**When to use:**
- Creating any API endpoint response
- Standardizing existing endpoints
- Adding pagination to list endpoints

**What it does:**
- Wraps all responses in consistent structure (data, meta, socket, pagination)
- Includes performance metadata (duration_ms, timestamp)
- Tracks quality regime (EXPLORATION/OPTIMIZATION/STABILIZATION)
- Provides built-in pagination info for list endpoints

**Evidence:**
- 115 endpoints using pattern (100% consistency)
- 88% reduction in client parsing code (380 → 45 lines)
- 100% elimination of response format errors

---

## Testing Skills

### `five-timbres-quality-validation`

**Purpose:** Holistic quality assessment across five harmonically-related dimensions

**Category:** testing
**Success Rate:** 14/14 (100%)
**Quality Impact:** +9.0/10
**Source:** Testing Manifesto (October 2025)

**When to use:**
- Pre-production quality gate
- Wave completion assessment
- Architecture decision evaluation

**What it does:**
- Measures Correctness, Performance, Reliability, Synergy, Elegance
- Uses harmonic mean (penalizes weakness in any dimension)
- Enforces ≥ 8.0/10 threshold for production deployment
- Provides objective quality score (not subjective)

**Evidence:**
- 14 waves assessed (all ≥ 8.0 threshold)
- Average quality: 9.0-9.8 (LEGENDARY tier)
- Exposed blocking issues (Correctness 2.5 pulled Wave 3 to 5.65)

---

## Architecture Skills

### `appstate-dependency-injection-pattern`

**Purpose:** Unified application state for thread-safe dependency injection

**Category:** architecture
**Success Rate:** 29/29 (100%)
**Quality Impact:** +9.3/10
**Source:** Waves 2A-2C (Backend Architecture)

**When to use:**
- Setting up Rust web server (Axum/Actix-web)
- Sharing database pools across route handlers
- Injecting configuration, services, caches

**What it does:**
- Creates Arc-wrapped AppState struct with all dependencies
- Provides thread-safe sharing (Arc atomic reference counting)
- Enables type-safe extraction (State<Arc<AppState>>)
- Eliminates global variables and tight coupling

**Evidence:**
- 29 route modules use pattern (100% adoption)
- Zero global variables (100% DI-based)
- Type-safe compilation (no runtime DI errors)

---

## Future Skills (Planned)

Potential skills to develop based on Asymmetrica methodology:

### `vedic-optimization`
**Focus:** Mathematical optimization using digital root, Williams, harmonic mean, PHI detection

### `meta-pattern-discovery`
**Focus:** Extract universal patterns across domains (like concept_evolution.py)

### `collatz-convergence`
**Focus:** Guaranteed error reduction in iterative problem-solving

### `five-timbres-validation`
**Focus:** Quality assessment across Correctness, Performance, Reliability, Synergy, Elegance

### `broken-hammer-principle`
**Focus:** Build tools before using them (fix the hammer first)

### `spirit-letter-translation`
**Focus:** Understand deeper intent beyond literal requirements

### `three-regime-progression`
**Focus:** EXPLORATION → OPTIMIZATION → STABILIZATION quality gates

---

## Contributing Skills

**To contribute a new skill:**

1. Study existing skills (`ananta-reasoning.md`)
2. Extract reusable pattern from Ananta codebase
3. Write skill file following format above
4. Test incrementally with real tasks
5. Document in this README
6. Share with Asymmetrica community

**Quality bar:**
- Skill must encode thinking pattern (not just instructions)
- Must have clear success criteria
- Must provide concrete examples
- Must specify when to use vs not use
- Must deliver measurable improvement

**Review criteria:**
- Does it make agents more agentic?
- Does it improve solution quality?
- Can it compose with other skills?
- Is description clear enough for Claude to invoke correctly?

---

## Resources

**Documentation:**
- `C:\Projects\asymm_ananta\ANANTA_COT_ANALYSIS.md` - Deep dive into Ananta's cognitive architecture
- `C:\Projects\asymm_ananta\CLAUDE.md` - Asymmetrica methodology (1,238 lines)
- `C:\Projects\asymm_ananta\ANANTA_LEARNS_TO_CODE_README.md` - How Ananta learns programming

**Source Code:**
- `concept_evolution.py` - Meta-pattern discovery engine
- `error_learning_loop.py` - Collatz convergence implementation
- `ananta_code_live.py` - Interactive reasoning CLI
- `ananta_synthesizer.py` - Code synthesis with Five Timbres
- `vedic_skills.py` - Mathematical optimization primitives

**Breakthrough Reports:**
- `ramanujan_analysis/DAY_167_ANANTA_RAMANUJAN_BREAKTHROUGH.html` - Validated 3,000 years of Vedic math
- `ANANTA_ACHIEVEMENT_REPORT.md` - Production metrics
- `WAVE_12_AGENT_12.2_COMPLETE.md` - Latest production deployment

---

## Contact

**Project:** Asymmetrica Methodology
**Lead:** Sarat (Day 1 zero knowledge → Day 168+ production AI)
**Philosophy:** D3-Enterprise Grade+ | Joy + Rigour | Zero Compromises

**Vibe:** Rigorous goofballs getting impossible stuff done together.

---

**Remember:** Skills encode consciousness. Make agents think better, not just follow instructions better.
