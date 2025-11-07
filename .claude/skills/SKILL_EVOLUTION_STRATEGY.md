# Skill Evolution Strategy
## How Skills Improve, Adapt, and Stay Relevant Over Time

**Purpose:** Ensure skill library evolves with project needs, incorporates learnings, deprecates obsolete patterns, and maintains high quality standards.

**Philosophy:** "Skills that don't evolve become technical debt. Living documentation adapts with experience. Dead documentation lies forgotten in dusty folders."

---

## The Four Evolution Mechanisms

### MECHANISM 1: Usage Tracking

**What to track:**
```yaml
# Automatically updated in skill frontmatter
last_used: YYYY-MM-DD           # Most recent application
success_rate: X/Y applications  # Successful uses / Total attempts
total_applications: N           # Lifetime usage count
```

**How to track:**

1. **After agent applies skill successfully:**
   ```
   - Increment success counter: 5/5 → 6/6
   - Update last_used timestamp
   - Log application context (task type, outcome)
   ```

2. **After agent applies skill unsuccessfully:**
   ```
   - Increment total attempts: 5/5 → 5/6
   - Update last_used timestamp
   - Log failure reason (why didn't pattern work?)
   ```

3. **Periodic review (monthly):**
   ```
   - Skills with 0 uses in 90 days → Flag for review
   - Skills with < 60% success rate → Flag for refinement
   - Skills with > 95% success rate + 10+ uses → Mark as PROVEN
   ```

**Usage Dashboard (conceptual):**

```markdown
## Skill Library Health Dashboard

**Total Skills:** 25
**Active (used in last 30 days):** 18 (72%)
**Dormant (unused 30-90 days):** 5 (20%)
**Obsolete (unused > 90 days):** 2 (8%)

**Top 5 Most Used:**
1. ananta-reasoning (14 applications, 100% success)
2. api-integration-complete-pattern (5 applications, 100% success)
3. migration-safety-rename-not-recreate (2 applications, 100% success)
4. asymmsocket-response-pattern (115 applications, 100% success)
5. five-timbres-quality-validation (14 applications, 100% success)

**Needs Attention:**
- database-query-optimization.md (3/7 success = 43%) → REFINE
- frontend-state-management.md (0 uses in 95 days) → REVIEW
```

---

### MECHANISM 2: Feedback Collection

**What feedback to collect:**

1. **What worked well:**
   ```
   Example application log:
   Date: 2025-11-01
   Skill: migration-safety-rename-not-recreate.md
   Task: Rename Customer table to customers
   Outcome: SUCCESS
   What worked: RENAME preserved 1,390 rows, < 1 second execution
   Notes: Foreign key update step was crucial (would have broken without it)
   ```

2. **What didn't work:**
   ```
   Date: 2025-11-02
   Skill: api-integration-complete-pattern.md
   Task: Connect Orders page to backend
   Outcome: PARTIAL (worked but slow)
   Issue: Pattern didn't cover optimistic updates for slow networks
   Suggestion: Add section on optimistic UI updates with rollback
   ```

3. **What was confusing:**
   ```
   Date: 2025-11-03
   Skill: five-timbres-quality-validation.md
   Task: Assess Wave 4 quality
   Outcome: SUCCESS (but took longer than expected)
   Confusion: "Synergy" timbre calculation unclear, needed examples
   Suggestion: Add 3 worked examples with specific synergy calculations
   ```

**Feedback Capture Methods:**

**Automatic (built into agent workflow):**
```
Agent applies skill → Tracks success/failure → Logs to skill metadata
```

**Manual (user or agent adds note):**
```
User: "The migration skill worked great, but step 3 was confusing"
Agent: Updates migration-safety-rename-not-recreate.md:
  - Adds clarification to Step 3
  - Increments version: 1.0.0 → 1.1.0
  - Logs change in version history
```

---

### MECHANISM 3: Pattern Refinement

**When to refine:**

1. **Success rate < 80%** (skill not working reliably)
2. **Common pitfalls emerge** (same mistake repeated 3+ times)
3. **New use cases discovered** (skill applies beyond original scope)
4. **Better approach found** (new pattern supersedes old)

**Refinement Process:**

#### Step 1: Identify Improvement Opportunity

```
Trigger: api-integration-complete-pattern.md success rate = 4/7 (57%)

Analysis:
- 3 failures all had same root cause: "Type mismatch errors"
- Root cause: Backend uses snake_case, frontend uses camelCase
- Current skill doesn't mention field name mapping
```

#### Step 2: Draft Update

```
New section to add:

### Common Mistake 2: Type Mismatch (Backend snake_case vs Frontend camelCase)

**What happens:**
Backend returns `customer_code`, frontend expects `customerCode`,
field appears as `undefined`, breaks UI.

**Why it happens:**
Backend uses snake_case (Rust convention), frontend uses camelCase
(JavaScript convention), no transformation layer.

**How to avoid:**
Map backend fields to frontend types in API client:
[Code example...]
```

#### Step 3: Test Update

```
Apply updated skill to new task → Verify improvement
If success: Proceed to Step 4
If still failing: Iterate on update
```

#### Step 4: Version Bump and Deploy

```yaml
# OLD frontmatter:
version: 1.0.0

# NEW frontmatter:
version: 1.1.0  # Minor version bump (backward compatible)

# If breaking change (pattern fundamentally changed):
version: 2.0.0  # Major version bump
```

```markdown
## Version History

### v1.1.0 (2025-11-05)
- Added Common Mistake 2: Type mismatch handling
- Added mapping function example for snake_case → camelCase
- Improved clarity of Step 7 (form submission)
- Success rate improved: 4/7 → 7/7 (100%)

### v1.0.0 (2025-11-01)
- Initial extraction from Wave 3C Frontend Integration
```

---

### MECHANISM 4: Deprecation Protocol

**When to deprecate:**

1. **Better pattern discovered** (new skill supersedes old)
2. **Technology changed** (framework/library upgraded, pattern obsolete)
3. **Unused for 180+ days** (no applications, likely irrelevant)
4. **Success rate < 40%** (pattern fundamentally flawed)

**Deprecation Steps:**

#### Step 1: Mark as Deprecated

```yaml
# In skill frontmatter:
deprecated: true
deprecated_date: 2025-12-01
replacement: new-skill-name.md
deprecation_reason: "Superseded by improved pattern with 95% success rate (vs 40%)"
```

#### Step 2: Add Deprecation Notice

```markdown
# [Skill Name] (DEPRECATED)

⚠️ **DEPRECATED as of 2025-12-01**

**Reason:** Superseded by improved pattern with higher success rate.

**Use instead:** [new-skill-name.md](./new-skill-name.md)

**Why deprecated:**
- Old pattern: 40% success rate (too low for production)
- New pattern: 95% success rate (proven across 20+ applications)
- New pattern handles edge cases old pattern missed

**If you must use this (legacy code):**
[Keep original content for reference...]
```

#### Step 3: Update README.md

```markdown
## Deprecated Skills

### `old-pattern-name` (Deprecated 2025-12-01)
**Replaced by:** [new-pattern-name.md](./category/new-pattern-name.md)
**Reason:** Low success rate (40%), better alternative available
**Keep for:** Legacy code reference only
```

#### Step 4: Archive (Don't Delete)

```
Move to archive folder for historical reference:

.claude/skills/
├── database/
│   └── migration-safety-rename-not-recreate.md  (active)
│
└── _archived/
    └── old-migration-pattern.md  (deprecated, kept for reference)
```

**Why keep deprecated skills:**
- Historical reference (understand project evolution)
- Legacy code maintenance (old code may still use pattern)
- Learning resource (why was old pattern insufficient?)

---

## Version Management

### Semantic Versioning

```
version: MAJOR.MINOR.PATCH

MAJOR: Breaking change (pattern fundamentally different)
  Example: 1.5.3 → 2.0.0 (changed from Class to Hook pattern)

MINOR: Backward-compatible enhancement (added sections, examples)
  Example: 1.5.3 → 1.6.0 (added new pitfall, improved example)

PATCH: Typos, clarifications, no functional change
  Example: 1.5.3 → 1.5.4 (fixed code typo, clarified step 3)
```

### Version History Format

```markdown
## Version History

### v2.0.0 (YYYY-MM-DD) - Breaking Change
**Changed:** Complete rewrite using Hook pattern instead of Class
**Reason:** Hooks more idiomatic in React 16.8+
**Migration:** See [migration guide](./migration-to-v2.md)
**Impact:** All applications must update code

### v1.6.0 (YYYY-MM-DD) - Enhancement
**Added:** Common Mistake 4 (race condition handling)
**Added:** Example 3 (pagination with search)
**Improved:** Clarified Step 5 (loading state conditions)
**Impact:** Backward compatible, optional adoption

### v1.5.4 (YYYY-MM-DD) - Patch
**Fixed:** Typo in code example (line 23)
**Clarified:** Step 3 wording (reduced ambiguity)
**Impact:** No functional change
```

---

## Quality Gates for Updates

**Before committing skill update:**

- [ ] Version number incremented correctly
- [ ] Version history updated with changes
- [ ] Examples tested (code actually works)
- [ ] Success indicators still accurate
- [ ] No broken links (cross-references valid)
- [ ] Frontmatter metadata updated (last_used, etc.)

**Quality threshold for new skills:**
- Success rate: Must demonstrate ≥ 1 successful application
- Quality score: Must achieve ≥ 8.0/10 in source wave
- Documentation: All 9 sections complete (no TODOs)
- Evidence: Real examples from actual project (not hypothetical)

---

## Review Schedule

**Monthly Review (1st of month):**
```
1. Check usage metrics (identify dormant skills)
2. Review success rates (flag skills < 80%)
3. Collect feedback from recent applications
4. Update last_reviewed timestamp in README
```

**Quarterly Deep Dive (Jan 1, Apr 1, Jul 1, Oct 1):**
```
1. Refine skills with < 80% success rate
2. Deprecate skills unused for 180+ days
3. Extract new patterns from recent waves
4. Update skill library health dashboard
```

**Annual Audit (Jan 1):**
```
1. Review entire library (25 skills → comprehensive assessment)
2. Consolidate duplicate/overlapping skills
3. Major version bumps for outdated patterns
4. Archive deprecated skills (move to _archived/)
5. Calculate year-over-year metrics (growth, quality trends)
```

---

## Evolution Metrics

**Track skill library improvement over time:**

```
Skills Added: N new skills per quarter
Skills Deprecated: N obsolete skills per quarter
Average Success Rate: X% (target: ≥ 85%)
Active Skill Ratio: Active / Total (target: ≥ 70%)
Quality Score: Average quality of all skills (target: ≥ 8.5/10)
```

**Goal:** Library improves continuously
- Success rates increase (better patterns emerge)
- Active ratio stays high (relevant skills)
- Quality scores rise (refinement over time)

---

## Continuous Improvement Loop

```
APPLY skill → TRACK usage → COLLECT feedback → REFINE pattern
    ↓                                              ↑
    └──────────────── LOOP ─────────────────────┘
```

**Each iteration:**
- Success rate increases (pattern refined based on failures)
- Documentation clarifies (based on confusion points)
- Examples expand (based on new use cases)
- Pitfalls documented (based on mistakes made)

**Result:** Skills compound in value. Wave 10 benefits from learnings of Waves 1-9.

---

## Example Evolution Timeline

**Month 1 (Initial Extraction):**
```
migration-safety-rename-not-recreate.md
version: 1.0.0
success_rate: 1/1 (100%)
applications: 1
status: NEW
```

**Month 2 (Second Application):**
```
version: 1.0.0
success_rate: 2/2 (100%)
applications: 2
feedback: "Step 3 unclear, added clarification"
status: VALIDATED
```

**Month 3 (Failure Reveals Gap):**
```
version: 1.0.0
success_rate: 2/3 (67%)
applications: 3
issue: "Forgot to update foreign keys, caused errors"
action: "Add explicit foreign key section"
status: NEEDS REFINEMENT
```

**Month 4 (Pattern Refined):**
```
version: 1.1.0 (MINOR bump)
success_rate: 3/3 (100% since update)
applications: 3 new (6 total)
changes: "Added foreign key update steps, examples"
status: IMPROVED
```

**Month 12 (Proven Pattern):**
```
version: 1.2.0
success_rate: 15/15 (100%)
applications: 15
status: PROVEN (used across 3 projects, zero failures)
```

---

## Summary: Evolution Principles

1. **Track everything** - Usage, success rate, feedback
2. **Refine relentlessly** - < 80% success rate → improve
3. **Deprecate gracefully** - Mark obsolete, suggest replacement, archive (don't delete)
4. **Version semantically** - Breaking changes = major, enhancements = minor, fixes = patch
5. **Review regularly** - Monthly metrics, quarterly deep dive, annual audit

**Outcome:** Skill library becomes more valuable over time. Knowledge compounds. Each wave makes future waves faster.

---

**Strategy Version:** 1.0.0
**Created:** November 1, 2025
**Last Updated:** November 1, 2025
**Maintained By:** Dr. Amara Osei (Knowledge Management Specialist)
