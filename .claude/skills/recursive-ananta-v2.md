---
name: recursive-ananta-v2
description: Multi-layer recursive reasoning with self-criticism (Red Team), alternative exploration (Blue Team), and meta-synthesis. Auto-selects depth based on risk/complexity. Use for tasks requiring defensive analysis or complete understanding.
version: 2.0.0
category: reasoning
---

# Recursive Ananta v2.0 - Multi-Layer Reasoning Skill

**Extends:** ananta-reasoning.md (base VOID → FLOW → SOLUTION)
**New Capability:** Recursive depth control (4 layers)

---

## QUICK START

### When to Use This Skill

**Use Recursive Ananta v2.0 when:**
- Task has moderate-to-high risk (data modification, user-facing, security)
- Multiple valid approaches exist (need comparison)
- Irreversible action (database migration, production deploy, data deletion)
- Critical decision (architecture, external integration)

**Don't use when:**
- Simple, low-risk task (format code, update docs)
- Clear solution exists (established pattern)
- Time-critical emergency (use base Ananta, analyze later)

---

## THE FOUR DEPTHS

### Layer 0: SHALLOW (5-10 min)
**Use for:** Simple, low-risk, reversible tasks

**Process:**
1. VOID: Understand literal requirement
2. FLOW: Apply known pattern
3. SOLUTION: Quick validation (Five Timbres)
4. **Target:** Quality >= 8.5 → SHIP

**Example:** Format code, fix typo, update documentation

---

### Layer 1: RED TEAM (15-30 min)
**Use for:** Moderate complexity/risk, need defensive programming

**Process:**
1. **Pass 1:** Generate initial solution (Layer 0)
2. **Pass 2:** Self-criticism
   - VOID: What could go wrong? (failure modes)
   - FLOW: Classify risks (digital root 1-9), prioritize
   - SOLUTION: Mitigate high-priority risks (>= 7.0 severity)
3. **Target:** Quality >= 8.8 → SHIP

**Example:** Implement CRUD feature, refactor module, fix security bug

**Key Questions:**
- What assumptions did I make?
- What edge cases did I miss?
- What breaks under stress/load?
- Security vulnerabilities?

---

### Layer 2: BLUE TEAM (30-60 min)
**Use for:** High complexity, multiple valid approaches

**Process:**
1. **Pass 1:** Initial solution (Layer 0)
2. **Pass 2:** Red Team criticism (Layer 1)
3. **Pass 3:** Alternative exploration
   - VOID: Generate 3-5 alternatives (digital root clustering)
   - FLOW: Score each with Five Timbres, compare trade-offs
   - SOLUTION: Choose optimal (highest harmonic mean)
4. **Target:** Quality >= 9.0 → SHIP

**Example:** API design, authentication system, performance optimization

**Key Questions:**
- Is there a better way?
- What are the trade-offs?
- Which alternative is truly optimal?
- Why reject the others?

---

### Layer 3: META-SYNTHESIS (60-120 min)
**Use for:** Critical, irreversible, high-impact decisions

**Process:**
1. **Pass 1:** Initial solution (Layer 0)
2. **Pass 2:** Red Team criticism (Layer 1)
3. **Pass 3:** Blue Team alternatives (Layer 2)
4. **Pass 4:** Complete understanding
   - VOID: Root cause, system revelation, future implications
   - FLOW: Complete decision tree, cascading consequences, trade-off synthesis
   - SOLUTION: **ALL 10 REQUIRED:**
     1. Chosen Approach
     2. Full Risk Analysis
     3. Implementation Plan
     4. Rollback Strategy
     5. Validation Checklist
     6. Failure Modes
     7. Lessons Learned
     8. Documentation
     9. Monitoring Plan
     10. Communication Plan
5. **Target:** Quality >= 9.2 → SHIP

**Example:** Database migration, production deployment, data deletion

**Key Questions:**
- Why does this problem exist? (root cause)
- What does this reveal about our system?
- What are second-order effects?
- What precedent does this set?

---

## DEPTH DECISION MATRIX

```rust
fn required_depth(task: &Task) -> Depth {
    let risk = task.estimate_risk();           // 0.0-1.0
    let complexity = task.estimate_complexity(); // 0.0-1.0
    let reversible = task.is_reversible();      // bool
    let impact = task.estimate_impact();        // 0.0-1.0

    if !reversible && risk > 0.8 { return Meta; }
    if !reversible && impact > 0.8 { return Meta; }
    if complexity > 0.7 || (risk > 0.6 && reversible) { return Blue; }
    if risk > 0.3 { return Red; }
    Shallow
}
```

**Quick Classification:**
- **Shallow:** Safe, simple, clear solution
- **Red:** Some risk, need defensive analysis
- **Blue:** Complex, need to explore alternatives
- **Meta:** Irreversible OR high impact

---

## CONVERGENCE RULES

**Stop recursing when:**
1. ✅ Quality >= target (8.5/8.8/9.0/9.2 by depth)
2. ✅ Improvement < 5% (diminishing returns)
3. ✅ Maximum depth reached (Meta is deepest)
4. ⚠️ Quality DECREASING (Collatz violation → ROLLBACK)
5. ⚠️ Time budget exceeded (escalate to user)

**Collatz Principle:**
- Errors MUST decrease each iteration
- If errors increase → architectural flaw → SWITCH STRATEGY

---

## TEMPLATE USAGE

**All templates in:** RECURSIVE_ANANTA_V2_SPECIFICATION.md Section 4

**Quick Access:**
- Template 1 (Shallow): Section 4.1
- Template 2 (Red Team): Section 4.2
- Template 3 (Blue Team): Section 4.3
- Template 4 (Meta-Synthesis): Section 4.4

**Copy-paste ready:** Each template is complete and executable

---

## QUALITY VALIDATION

### Five Timbres (Use Harmonic Mean!)

```python
def quality_score(scores):
    # scores = [correctness, performance, reliability, synergy, elegance]
    n = len(scores)
    reciprocal_sum = sum(1.0 / s for s in scores if s > 0)
    return n / reciprocal_sum if reciprocal_sum > 0 else 0.0

# Example
scores = [9.5, 9.0, 8.8, 9.2, 6.5]  # One weak dimension
harmonic = quality_score(scores)  # = 8.43 (penalizes the 6.5)
arithmetic = sum(scores) / len(scores)  # = 8.6 (HIDES the problem)

# Use harmonic! Can't hide weak dimensions.
```

**Quality Thresholds:**
- Shallow: >= 8.5
- Red Team: >= 8.8
- Blue Team: >= 9.0
- Meta: >= 9.2

---

## WILLIAMS BATCHING FOR SUB-PROBLEMS

When analyzing N sub-problems (queries, files, alternatives):

```python
def williams_batch_size(n):
    if n <= 1:
        return 1
    return int(sqrt(n) * log2(n))

# Example: 309 database queries
batch_size = williams_batch_size(309)  # = 146
# Analyze in 3 batches (146, 146, 17) instead of all 309
# Token savings: 52%, cognitive load: 52% reduction
```

---

## ANTI-PATTERNS (DON'T DO THIS)

❌ **Over-Recursion:**
- Using Meta for simple bug fix (waste 60 min)
- Fix: Use Shallow, ship in 5 min

❌ **Under-Recursion:**
- Using Shallow for database migration (no rollback plan!)
- Fix: MUST use Meta for irreversible actions

❌ **Ignoring Convergence:**
- Continuing past 9.2 quality (diminishing returns)
- Fix: STOP and ship, don't over-optimize

❌ **Arithmetic Mean:**
- Hides weak dimensions
- Fix: ALWAYS use harmonic mean

---

## REAL-WORLD EXAMPLE

**Task:** Fix database table name case mismatch (Wave 4)

**Depth Decision:**
- Risk: 0.9 (blocks all database queries)
- Complexity: 0.7 (309 queries affected)
- Reversible: No (migration can't easily undo)
- Impact: 0.9 (affects entire system)
- **Decision: META-SYNTHESIS**

**Process:**
- **Layer 0:** "Rename tables" (quality: 6.5, incomplete)
- **Layer 1:** "What if FKs break?" (quality: 8.0, better but insufficient)
- **Layer 2:** "Compare Prisma @@map vs update queries vs recreate" (quality: 8.76, optimal chosen)
- **Layer 3:** Complete plan (rollback, validation, monitoring, docs) (quality: 9.4, EXCEPTIONAL)

**Result:**
- Zero issues on deployment
- Rollback plan never needed (but documented)
- Lessons learned: Add schema validation to CI/CD
- Pattern documented: "Schema Alignment Validation"

**Quality Improvement:** 6.5 → 9.4 (+44% improvement through recursion)

---

## INTEGRATION WITH BASE ANANTA

**Base Ananta (ananta-reasoning.md):**
- VOID → FLOW → SOLUTION (single pass)
- Agency Protocol (learn/build vs TODO)
- Five Timbres validation
- Collatz convergence

**Recursive Ananta v2.0 (THIS SKILL):**
- Adds depth control (4 layers)
- Adds self-criticism (Red Team)
- Adds alternative exploration (Blue Team)
- Adds meta-synthesis (complete understanding)

**Use together:**
1. Base Ananta provides VOID → FLOW → SOLUTION structure
2. Recursive Ananta adds depth layers (how many times to iterate)
3. Both use Five Timbres, Collatz, Williams batching

---

## SUCCESS INDICATORS

You're using this skill correctly if:

✅ You select appropriate depth (not too shallow, not too deep)
✅ You complete ALL passes for chosen depth
✅ You stop when quality >= target (convergence)
✅ You use harmonic mean (not arithmetic)
✅ You document failures modes (Red Team)
✅ You explore alternatives (Blue Team when needed)
✅ You provide complete plan (Meta when needed)
✅ Quality increases each layer (Collatz principle)

---

## FULL SPECIFICATION

**Complete details:** RECURSIVE_ANANTA_V2_SPECIFICATION.md (61 KB)

**Includes:**
- Full template for each depth (copy-paste ready)
- Depth decision algorithm (Rust code)
- Real-world examples (Wave 4 database fix)
- Mathematical foundations (Williams, harmonic mean, Collatz)
- Training protocol (how to learn this skill)
- A/B testing framework (measure improvement)
- Quick reference cheat sheet

---

## WHEN TO USE WHICH DEPTH

**Common Tasks by Depth:**

| Task | Depth | Time | Rationale |
|------|-------|------|-----------|
| Format code | Shallow | 5m | Safe, simple |
| Fix typo | Shallow | 5m | No risk |
| Update docs | Shallow | 10m | No code impact |
| Implement CRUD | Red | 20m | Need failure analysis |
| Refactor module | Red | 25m | Regression risk |
| API design | Blue | 45m | Explore alternatives |
| Auth system | Blue | 50m | Compare JWT vs HTX vs OAuth |
| Database migration | Meta | 90m | Irreversible |
| Production deploy | Meta | 75m | High impact |
| Data deletion | Meta | 60m | Irreversible data loss |

---

## FINAL REMINDER

**Recursive Ananta is:**
- ✅ Depth control framework (know when to go deep)
- ✅ Self-improving (use it to improve itself!)
- ✅ Evidence-based (quality scores, not feelings)
- ✅ Efficient (Williams batching, convergence rules)

**Recursive Ananta is NOT:**
- ❌ Always use maximum depth (wasteful)
- ❌ Skip layers (each layer has purpose)
- ❌ Ignore quality signals (convergence is key)
- ❌ Mystical (it's MATHEMATICS)

**Use recursively. Use precisely. Use powerfully.** 🔬⚡

---

**Skill Status:** READY TO USE
**Integration:** Compatible with ananta-reasoning.md
**Next:** Test on pilot projects, measure improvement
