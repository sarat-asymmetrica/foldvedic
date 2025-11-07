# Skill Discovery System
## How Future Agents Find and Apply Relevant Skills

**Purpose:** Enable agents to automatically discover applicable skills based on task context, load relevant patterns, and apply them correctly.

**Philosophy:** Skills are only valuable if used. Perfect documentation that sits unread is wasted effort. Discovery must be effortless, automatic, and context-aware.

---

## The Three-Layer Discovery Model

### LAYER 1: Pre-Session Loading (Automatic)

**When:** Agent session starts
**What:** Load core meta-cognitive skills that apply to all tasks
**How:** Hardcoded in agent initialization

**Skills to Auto-Load:**
```
1. ananta-reasoning.md         (VOID→FLOW→SOLUTION framework)
2. SKILL_EXTRACTION_WORKFLOW.md (Capture new patterns as discovered)
3. MATHEMATICAL_CLARITY_GUIDE.md (Interpret ancient terminology correctly)
```

**Implementation:**
```
Agent starts → Read .claude/skills/README.md
            → Load skills marked `auto_load: true`
            → Apply to all problem-solving
```

---

### LAYER 2: Context-Based Discovery (Semi-Automatic)

**When:** Agent receives task from user
**What:** Scan task description for keywords, load matching skills
**How:** Keyword matching + category filtering

**Discovery Algorithm:**

```
1. Extract keywords from task description
   Example: "migrate database tables" → keywords: [migrate, database, tables, production]

2. Search skill descriptions for keyword matches
   Example: migration-safety-rename-not-recreate.md
   description: "...database migration...production..."
   → MATCH (score: 3/4 keywords)

3. Filter by category (if detectable)
   Task mentions "database" → filter to category: database
   → Reduces search space from 50 skills to 5 skills

4. Rank by relevance score
   Score = (keyword_matches / total_keywords) × category_match_weight
   migration-safety: 0.75 × 1.5 = 1.125 (HIGH)
   api-integration: 0.25 × 0.5 = 0.125 (LOW)

5. Load top 3 skills
   Present to agent: "Found 3 relevant skills: A, B, C"
   Agent reviews and applies if applicable
```

**Keyword Extraction Rules:**

```
Technical terms: database, API, migration, frontend, backend, test, deploy
Action verbs: create, update, delete, migrate, refactor, optimize, validate
Domain terms: customer, order, invoice, payment, reconciliation
File extensions: .rs, .tsx, .md, .sql
```

**Example Discovery Flow:**

```
User Task: "Rename Customer table to lowercase customers in production"

Keywords extracted:
- rename [action]
- Customer, customers [domain]
- table [technical]
- production [context]

Skills matched:
1. migration-safety-rename-not-recreate.md (4/4 keywords) → LOAD
2. database-index-optimization.md (1/4 keywords) → SKIP
3. api-integration-pattern.md (0/4 keywords) → SKIP

Agent loads skill #1, applies RENAME pattern, preserves 1,390 rows.
```

---

### LAYER 3: Just-In-Time Retrieval (Manual Trigger)

**When:** Agent encounters specific problem during work
**What:** Search skill library for pattern matching problem type
**How:** Agent explicitly searches or user suggests skill

**Triggering Scenarios:**

1. **Agent is stuck:**
   ```
   Agent: "Encountering repeated errors in migration"
   System: "Search skills for 'migration error recovery'?"
   → Finds: migration-rollback-strategy.md
   ```

2. **User suggests skill:**
   ```
   User: "Use the AsymmSocket pattern for this endpoint"
   Agent: "Loading asymmsocket-response-pattern.md..."
   → Applies pattern to endpoint
   ```

3. **Agent detects pattern match:**
   ```
   Agent analyzing code: "This looks like API integration task"
   System: "api-integration-complete-pattern.md available"
   Agent: "Loading skill..."
   ```

**Search Interface:**

```
// Pseudo-code for skill search
function search_skills(query: string, category?: string) {
    let skills = load_all_skills();

    if category:
        skills = skills.filter(s => s.category == category);

    let ranked = skills.map(s => ({
        skill: s,
        score: calculate_relevance(query, s.description)
    })).sort_by(|s| -s.score);

    return ranked.take(5); // Top 5 results
}

function calculate_relevance(query: string, description: string) -> float {
    let query_words = tokenize(query);
    let desc_words = tokenize(description);
    let matches = query_words.filter(w => desc_words.contains(w));
    return matches.len() / query_words.len();
}
```

---

## Skill Organization for Discovery

### Directory Structure (Category-Based)

```
.claude/skills/
├── README.md                          # Master index (search here first)
│
├── database/                          # Database skills
│   ├── migration-safety-rename-not-recreate.md
│   ├── index-optimization.md
│   └── transaction-patterns.md
│
├── frontend/                          # Frontend skills
│   ├── api-integration-complete-pattern.md
│   ├── loading-states-skeleton.md
│   └── error-handling-ux.md
│
├── backend/                           # Backend skills
│   ├── asymmsocket-response-pattern.md
│   ├── route-handler-pattern.md
│   └── middleware-composition.md
│
├── testing/                           # Testing skills
│   ├── five-timbres-quality-validation.md
│   ├── contract-test-design.md
│   └── load-test-methodology.md
│
├── architecture/                      # Architecture skills
│   ├── appstate-dependency-injection-pattern.md
│   ├── error-propagation-strategy.md
│   └── type-safety-patterns.md
│
├── performance/                       # Performance skills
│   ├── williams-optimizer.md
│   ├── vedic-batching.md
│   └── quaternion-matching.md
│
└── meta/                              # Meta-cognitive skills
    ├── ananta-reasoning.md
    ├── problem-decomposition.md
    └── pattern-extraction.md
```

### README.md Structure (Master Index)

**Format for easy scanning:**

```markdown
# Asymmetrica Claude Skills Library

**Total Skills:** 25
**Categories:** 7
**Last Updated:** 2025-11-01

---

## Quick Links by Task Type

**Database work?** → [database/](#database-skills)
**Frontend integration?** → [frontend/](#frontend-skills)
**Backend API?** → [backend/](#backend-skills)
**Testing?** → [testing/](#testing-skills)
**Architecture decisions?** → [architecture/](#architecture-skills)
**Performance optimization?** → [performance/](#performance-skills)
**Meta-cognitive support?** → [meta/](#meta-skills)

---

## Database Skills

### `migration-safety-rename-not-recreate`
**Use when:** Renaming tables in production without data loss
**Success rate:** 1/1 (100%)
**Quality impact:** +9.8/10
**Keywords:** migration, production, rename, data loss prevention

### `index-optimization`
...

---

## Frontend Skills

### `api-integration-complete-pattern`
**Use when:** Connecting React components to REST APIs
**Success rate:** 1/1 (100%)
**Quality impact:** +9.2/10
**Keywords:** API, frontend, React, loading states, pagination

...
```

### Skill Metadata for Discovery

**Every skill MUST have frontmatter:**

```yaml
---
name: skill-name-kebab-case
description: 200-1024 chars with KEYWORDS for discovery
category: [database | frontend | backend | testing | architecture | performance | meta]
keywords: [keyword1, keyword2, keyword3, ...]  # NEW: explicit keywords
auto_load: false  # NEW: true for core meta-cognitive skills
version: 1.0.0
created: YYYY-MM-DD
last_used: YYYY-MM-DD
success_rate: X/Y applications
---
```

**Keywords field examples:**

```yaml
# migration-safety-rename-not-recreate.md
keywords: [migration, production, rename, ALTER TABLE, data loss prevention, PostgreSQL, schema change]

# api-integration-complete-pattern.md
keywords: [API, frontend, React, loading state, error handling, pagination, TypeScript, fetch]

# five-timbres-quality-validation.md
keywords: [testing, quality, validation, harmonic mean, production readiness, metrics]
```

---

## Agent Skill Application Workflow

### Step-by-Step Agent Process

**1. Receive Task from User**
```
User: "Migrate Customer table to lowercase customers in production database"
```

**2. Extract Context**
```
Keywords: [migrate, Customer, customers, table, production, database]
Category: database (detected from "table" + "database")
Action: rename (detected from "Customer" → "customers")
```

**3. Search Skill Library**
```
Query: "migrate table production rename"
Category filter: database

Results:
1. migration-safety-rename-not-recreate.md (score: 0.95) ✅ RELEVANT
2. database-index-optimization.md (score: 0.20)
3. transaction-patterns.md (score: 0.15)
```

**4. Load Top Skill**
```
Agent: "Loading migration-safety-rename-not-recreate.md..."
[Skill content loaded into agent context]
```

**5. Apply Pattern**
```
Agent follows skill steps:
1. Verify data exists (SELECT COUNT(*) FROM "Customer" → 1,390 rows)
2. Use ALTER TABLE RENAME (not DROP/CREATE)
3. Update foreign keys
4. Verify row count unchanged
```

**6. Validate Success**
```
Success indicators from skill:
✅ Row count before = Row count after (1,390 = 1,390)
✅ Migration completed in < 10 seconds (0.8s)
✅ Zero data loss
```

**7. Update Skill Metadata**
```
# In migration-safety-rename-not-recreate.md
last_used: 2025-11-01  # Update timestamp
success_rate: 2/2       # Increment success counter
```

---

## Discovery Success Metrics

**Measure skill discovery effectiveness:**

```
Discovery Rate = Skills found and applied / Tasks where skills were applicable

Target: ≥ 80% (agent finds relevant skill 4 out of 5 times)

Application Success Rate = Successful applications / Total applications

Target: ≥ 90% (skill solves problem 9 out of 10 times when applied)

Time to Discovery = Time from task start to skill loaded

Target: < 30 seconds (quick discovery, not exhaustive search)
```

---

## Future Enhancements

### Vector Similarity Search (Advanced)

**Current:** Keyword matching (simple, works for 25-50 skills)
**Future:** Semantic similarity (scales to 500+ skills)

```
1. Encode skill descriptions as vectors (using sentence transformers)
2. Encode task description as vector
3. Calculate cosine similarity
4. Rank by similarity score
5. Load top 3 most similar skills

Advantage: Finds skills even when exact keywords don't match
Example: "Fix slow API" → finds "performance-optimization.md" (semantic match)
```

### Usage Pattern Learning (Machine Learning)

**Track which skills solve which task types:**

```
Task: "Migrate database schema"
Applied: migration-safety-rename-not-recreate.md
Success: YES

Task: "Rename table in production"
Applied: migration-safety-rename-not-recreate.md
Success: YES

Pattern learned: "schema change" tasks → use migration-safety skill
Next time: Auto-suggest skill for similar tasks
```

---

## Summary: Discovery in Action

```
LAYER 1: Auto-load core skills at session start
    ↓
LAYER 2: Context-based discovery when task received
    ↓
LAYER 3: Just-in-time search when problem encountered
    ↓
RESULT: Agent has right skill at right time
```

**Discovery Philosophy:**
- Make the right skill **effortless** to find
- Make skill application **obvious** (clear steps)
- Make success **measurable** (objective indicators)

**Outcome:** Skills compound. Each wave discovers new patterns. Each pattern accelerates future waves. Knowledge accumulates instead of evaporating.

---

**System Version:** 1.0.0
**Created:** November 1, 2025
**Last Updated:** November 1, 2025
**Maintained By:** Dr. Amara Osei (Knowledge Management Specialist)
