# Skill Extraction Workflow
## From Ephemeral Agent Learnings to Permanent Cognitive Patterns

**Purpose:** Transform one-time agent solutions into reusable institutional knowledge.

**Philosophy:** Every solved problem is a pattern waiting to be extracted. Every mistake avoided is wisdom waiting to be encoded. Skill extraction converts ephemeral brilliance into permanent advantage.

---

## The Four-Phase Extraction Process

### PHASE 1: SKILL DISCOVERY (After Each Wave)

**When to Extract:**
- Wave completes successfully (quality ≥ 8.0/10)
- Novel pattern emerged (not already documented)
- Solution worked exceptionally well (validated with evidence)
- Mistake avoided through expertise (prevented failure)
- Future agents would benefit (generalizable beyond one task)

**Discovery Questions:**

1. **What pattern emerged?**
   - Technical pattern (code structure, algorithm, architecture)
   - Process pattern (workflow, methodology, decision framework)
   - Domain pattern (business logic, requirements interpretation)
   - Meta-pattern (thinking approach, problem decomposition)

2. **Why did it work so well?**
   - Mathematical optimality (provable efficiency gain)
   - Architectural elegance (reduced complexity, improved maintainability)
   - Error prevention (avoided common failure modes)
   - Performance gain (measurable speedup or resource savings)
   - Quality improvement (higher Five Timbres scores)

3. **What makes it reusable?**
   - Applies to multiple similar problems (N > 3)
   - Technology/domain independent (transferable)
   - Clear triggering conditions (when to use obvious)
   - Teachable in < 5 minutes (concise, not complex)

4. **What would have happened without it?**
   - Time wasted (estimate hours saved)
   - Quality compromised (errors that would have occurred)
   - Architecture debt (shortcuts that would have been taken)
   - Knowledge lost (would this pattern have been rediscovered?)

**Evidence Requirements:**

Every skill extraction MUST include:
- **Source:** Which wave/agent discovered this pattern
- **Problem:** What specific challenge was solved
- **Solution:** Exact pattern that worked
- **Evidence:** Quality scores, metrics, test results
- **Impact:** What improved (time, quality, reliability, etc.)

**Output:** List of candidate skills with evidence

---

### PHASE 2: SKILL VALIDATION (Quality Filter)

**Validation Criteria:**

#### 1. Reusability Test
- **Question:** Will this pattern apply to 3+ future scenarios?
- **Pass:** Yes, generalizable across problems
- **Fail:** One-off solution specific to single context

**Example Pass:**
```
Pattern: "Database Migration Safety - Use RENAME not DROP/CREATE"
Reusable: YES (applies to all production schema changes)
Scenarios: Customer table rename, adding indexes, column renames
```

**Example Fail:**
```
Pattern: "Fix Modal.svelte duplicate script tag"
Reusable: NO (specific to one broken file)
Impact: One-time fix, not a pattern
```

#### 2. Teachability Test
- **Question:** Can an agent learn this pattern in < 5 minutes?
- **Pass:** Clear, concise, actionable
- **Fail:** Requires extensive context or training

**Measure:**
```
Word count: < 500 words for core pattern explanation
Steps: ≤ 7 distinct steps (Miller's Law - working memory limit)
Examples: ≥ 2 concrete examples provided
Success criteria: Objective and measurable
```

#### 3. Success Criteria Test
- **Question:** Can we objectively measure if this skill was applied correctly?
- **Pass:** Clear success/failure indicators
- **Fail:** Subjective or vague outcomes

**Example Pass:**
```
Skill: "Harmonic Mean Quality Validation"
Success: Quality score ≥ 8.0 AND no dimension < 6.0
Failure: Quality < 8.0 OR any dimension < 6.0
Objective: Mathematical calculation, no ambiguity
```

**Example Fail:**
```
Skill: "Write clean code"
Success: ??? (undefined)
Failure: ??? (undefined)
Objective: NO (subjective judgment)
```

#### 4. Impact Test
- **Question:** Does applying this skill measurably improve quality or efficiency?
- **Pass:** ≥10% improvement in at least one dimension
- **Fail:** Marginal or unmeasurable benefit

**Dimensions to Measure:**
```
Time saved: ≥ 1 hour per application
Error prevention: ≥ 1 class of bugs eliminated
Quality improvement: ≥ 0.5 point increase in Five Timbres score
Performance gain: ≥ 2× speedup
Reliability boost: ≥ 50% reduction in failure rate
```

**Validation Decision Tree:**

```
Reusable (3+ scenarios)?
├─ NO → REJECT (Document as case study, not skill)
└─ YES → Continue

Teachable (< 5 min)?
├─ NO → SIMPLIFY or REJECT
└─ YES → Continue

Measurable Success?
├─ NO → REFINE criteria or REJECT
└─ YES → Continue

Measurable Impact (≥10%)?
├─ NO → REJECT (Nice-to-have, not critical skill)
└─ YES → ACCEPT (Proceed to documentation)
```

**Output:** Validated skills ready for documentation

---

### PHASE 3: SKILL DOCUMENTATION (Standard Format)

**Template Location:** `.claude/skills/SKILL_TEMPLATE.md`

**Required Sections:**

#### 1. Frontmatter (YAML)
```yaml
---
name: skill-name-kebab-case
description: Brief description (200-1024 chars). Used by Claude for skill selection.
category: [database | frontend | backend | testing | architecture | performance | security]
version: 1.0.0
created: YYYY-MM-DD
last_used: YYYY-MM-DD
success_rate: X/Y applications
source_wave: Wave X, Agent Name
---
```

#### 2. Purpose
- **1 sentence:** What this skill does
- **Why it exists:** Problem it solves
- **When discovered:** Wave/agent that found it

#### 3. When to Use (Triggering Conditions)
- **Specific scenarios:** List 3-5 concrete situations
- **Anti-indicators:** When NOT to use this skill
- **Prerequisites:** What must be true before applying

#### 4. The Pattern (Step-by-Step)
- **Numbered steps:** Clear, actionable process
- **Code examples:** Concrete implementations
- **Decision points:** What to do when X vs Y

#### 5. Why It Works (Underlying Principle)
- **Mathematical basis:** Algorithm, formula, proof
- **Architectural reason:** Why this structure is optimal
- **Domain logic:** Business rule that necessitates this pattern

#### 6. Evidence (Validation)
- **Quality Score:** X/10 (Five Timbres breakdown)
- **Performance Impact:** X% improvement, Y seconds saved
- **Risk Reduction:** Z class of errors eliminated
- **Source:** Link to wave report where this was proven

#### 7. Examples (Real Usage)
- **Minimum 2 examples** from actual project
- **Before/After:** What changed when skill applied
- **Outcome:** Measurable results

#### 8. Related Skills
- **Complementary:** Skills that work well together
- **Prerequisites:** Skills to master first
- **Alternatives:** Other approaches to same problem

#### 9. Pitfalls (Common Mistakes)
- **What can go wrong:** 3-5 common errors
- **How to avoid:** Specific prevention tactics
- **How to detect:** Warning signs skill misapplied

**Documentation Standards:**

```
Clarity: Any agent should understand without additional context
Conciseness: Core pattern in ≤ 500 words
Completeness: All 9 sections filled with evidence
Examples: Real code from actual waves
Metrics: Objective measurements, not opinions
```

**Output:** Complete skill document ready for library

---

### PHASE 4: SKILL ORGANIZATION (Retrieval Strategy)

**Directory Structure:**
```
.claude/skills/
├── README.md                          # Index of all skills
├── SKILL_EXTRACTION_WORKFLOW.md       # This document
├── SKILL_TEMPLATE.md                  # Blank template
│
├── database/                          # Database skills
│   ├── migration-safety.md
│   ├── index-optimization.md
│   └── transaction-patterns.md
│
├── frontend/                          # Frontend skills
│   ├── api-integration-pattern.md
│   ├── loading-states.md
│   └── error-handling-ux.md
│
├── backend/                           # Backend skills
│   ├── asymmsocket-response.md
│   ├── route-handler-pattern.md
│   └── middleware-composition.md
│
├── testing/                           # Testing skills
│   ├── five-timbres-validation.md
│   ├── contract-test-design.md
│   └── load-test-methodology.md
│
├── architecture/                      # Architecture skills
│   ├── appstate-dependency-injection.md
│   ├── error-propagation-strategy.md
│   └── type-safety-patterns.md
│
├── performance/                       # Performance skills
│   ├── williams-optimizer.md         # Already exists
│   ├── vedic-batching.md
│   └── quaternion-matching.md
│
└── meta/                              # Meta-cognitive skills
    ├── ananta-reasoning.md           # Already exists
    ├── problem-decomposition.md
    └── pattern-extraction.md
```

**Categorization Rules:**

1. **Primary category:** Choose ONE that best fits
2. **Cross-references:** Link to related skills in other categories
3. **Metadata tags:** Add keywords for search

**Naming Conventions:**

```
File name: lowercase-kebab-case.md
Skill name (in YAML): matches filename without .md
Max length: 64 characters
Format: [problem]-[solution-approach].md

Examples:
✅ database-migration-safety.md
✅ api-integration-pattern.md
✅ five-timbres-validation.md
❌ DatabaseMigrationSafety.md (wrong case)
❌ skill_for_database_migrations.md (wrong separator)
❌ migration.md (too vague)
```

**README.md Updates:**

Every new skill MUST be added to `.claude/skills/README.md`:

```markdown
### `skill-name` - Brief Title

**Purpose:** One-sentence description

**Category:** [category]

**When to use:**
- Scenario 1
- Scenario 2
- Scenario 3

**Success rate:** X/Y applications (Z%)

**Quality impact:** +N points (Five Timbres average)

**Source:** Wave X, Agent Name (Date)
```

**Search Strategy:**

Agents should search skills by:
1. **Category:** Browse relevant category folder
2. **Keyword:** Search README.md for terms
3. **Problem type:** Match triggering conditions
4. **Related skills:** Follow cross-references

---

## Summary: Four Phases at a Glance

```
DISCOVERY → VALIDATION → DOCUMENTATION → ORGANIZATION
    ↓            ↓              ↓               ↓
 Identify    Filter for    Write complete   Categorize
 patterns    reusability   skill document   and index

Input:      Input:        Input:           Input:
Wave        Candidate     Validated        Skill
reports     skills        skills           documents

Output:     Output:       Output:          Output:
Candidate   Validated     Skill docs       Searchable
skills      skills                         library
```

---

## Quality Metrics (Skill Library Health)

**Library-Level Metrics:**

```
Total Skills: N
Categories: 7 (database, frontend, backend, testing, architecture, performance, meta)
Average Quality Score: X/10 (from source waves)
Success Rate: Y% (applications where skill helped)
Time Saved: Z hours (cumulative across all applications)
```

**Per-Skill Metrics:**

```
Applications: N times used
Success Rate: X/N successful applications
Average Quality Impact: +Y points (Five Timbres)
Time Saved per Use: Z minutes average
Last Updated: YYYY-MM-DD
```

**Health Indicators:**

```
✅ Healthy Library:
- ≥ 10 skills documented
- ≥ 80% success rate across skills
- All skills used ≥ 1 time
- README updated within 7 days

⚠️ Needs Attention:
- < 10 skills (insufficient coverage)
- < 60% success rate (poor quality filter)
- Skills unused for > 30 days (relevance check needed)
- README outdated > 14 days (maintenance lag)

❌ Critical Issues:
- Skills with 0% success rate (remove or fix)
- Duplicate skills (consolidate)
- Missing categories (gaps in coverage)
```

---

## Continuous Improvement

**Skill Evolution:**

1. **Usage Tracking:** Log each application (success/failure)
2. **Feedback Collection:** What worked, what didn't
3. **Pattern Refinement:** Update skill based on learnings
4. **Version Bumping:** Increment version when updated
5. **Deprecation:** Mark obsolete skills, suggest replacements

**Version Management:**

```yaml
version: 1.0.0  # Initial extraction
version: 1.1.0  # Minor refinement (examples added)
version: 1.2.0  # Minor update (new pitfall documented)
version: 2.0.0  # Major update (pattern significantly changed)
```

**Deprecation Protocol:**

When a skill becomes obsolete:

1. Add `deprecated: true` to frontmatter
2. Add `replacement: new-skill-name.md` to frontmatter
3. Keep file for historical reference
4. Update README with deprecation notice

---

## Success Indicators

**Skill extraction working well when:**

✅ Agents reference skills before starting work
✅ Quality scores improve over time (learning compound effect)
✅ Common mistakes eliminated (pitfalls prevent errors)
✅ Time-to-completion decreases (patterns accelerate work)
✅ Skills are reused ≥ 3 times each (validation of reusability)
✅ New skills emerge from recent waves (continuous extraction)
✅ Deprecated skills replaced (evolution, not stagnation)

**Skill extraction needs work when:**

❌ Agents don't use skill library (visibility problem)
❌ Skills rarely reused (poor reusability filter)
❌ No new skills added in 30 days (extraction stopped)
❌ Success rate declining (quality filter too loose)
❌ Duplicate patterns emerging (consolidation needed)

---

## Final Checklist: Creating a New Skill

Before committing a new skill to the library:

- [ ] Validated: Passes all 4 validation criteria
- [ ] Documented: All 9 sections complete with evidence
- [ ] Categorized: Placed in correct category folder
- [ ] Named: Follows kebab-case naming convention
- [ ] Indexed: Added to README.md with metadata
- [ ] Tested: Applied successfully ≥ 1 time
- [ ] Reviewed: Cross-referenced with related skills
- [ ] Proven: Quality score ≥ 8.0 in source wave

---

**Philosophy:** "Methods that aren't documented die with their creators. Skills that aren't extracted are relearned expensively. Institutional memory is the difference between perpetual novices and compounding experts."

**Outcome:** Every wave makes future waves faster, smarter, and more successful. Knowledge accumulates instead of evaporating.

---

**Workflow Version:** 1.0.0
**Created:** November 1, 2025
**Last Updated:** November 1, 2025
**Maintained By:** Dr. Amara Osei (Knowledge Management Specialist)
