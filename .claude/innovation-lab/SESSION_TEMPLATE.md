# Innovation Session Template

**USE THIS TEMPLATE TO START NEW SESSIONS**

Copy this structure to `sessions/YYYY-MM-DD-descriptive-name/`

---

## File Structure

```
sessions/YYYY-MM-DD-descriptive-name/
├── HYPOTHESIS.md       # START HERE - Initial hypothesis and goals
├── EXPLORATION.md      # Document process as you go
├── [code files]        # Prototypes, experiments, tests
├── [benchmark data]    # Performance measurements
├── FINDINGS.md         # END HERE - Results and quality score
└── INTEGRATION.md      # Optional - Production integration plan
```

---

## HYPOTHESIS.md Template

```markdown
# Session Hypothesis: [Topic Name]

**Date:** YYYY-MM-DD
**Duration Estimate:** X-Y hours
**Regime:** EXPLORATION (30%) / OPTIMIZATION (20%) / STABILIZATION (50%)
**Persona:** Agent [Name] - [Expertise]

---

## Problem Statement

[Clear description of the challenge, opportunity, or innovation to explore]

---

## Hypothesis

[What we believe is possible and why]

**If this works, we expect:**
- [Specific outcome 1]
- [Specific outcome 2]
- [Specific outcome 3]

---

## Approach

**Recursive Ananta Depth:** Layer [0/1/2/3]

**VOID Phase:**
- [What questions need answering]
- [What unknowns exist]
- [What research is needed]

**FLOW Phase:**
- [What experiments to run]
- [What prototypes to build]
- [What alternatives to explore]

**SOLUTION Phase:**
- [What constitutes success]
- [How to measure quality]
- [How to validate results]

---

## Success Criteria

**Minimum Viable:**
- [ ] [Criterion 1]
- [ ] [Criterion 2]
- [ ] [Criterion 3]

**Ideal Outcome:**
- [ ] [Stretch goal 1]
- [ ] [Stretch goal 2]
- [ ] [Stretch goal 3]

**Quality Target:**
- EXPLORATION: ≥ 7.0
- OPTIMIZATION: ≥ 8.5
- STABILIZATION: ≥ 9.0

---

## Context

**Current State:**
- [What exists today]
- [What's working]
- [What's limited]

**Archaeology Findings:** (if relevant)
- [Related tech from TECH_ARCHAEOLOGY/]
- [File locations, metrics, quality scores]

**Constraints:**
- In scope: [What's allowed]
- Out of scope: [What's off-limits]
- Time budget: [X hours]

---

## Tools & Technologies

- [Language/framework]
- [Libraries/crates]
- [Benchmarking tools]
- [Testing frameworks]

---

## Prior Art

- [Existing implementations]
- [Research papers]
- [Similar approaches]

---

**BEGIN EXPLORATION** →
```

---

## EXPLORATION.md Template

```markdown
# Exploration Notes: [Topic Name]

**Live documentation of the innovation process**

---

## Session Start

**Time:** HH:MM
**Energy Level:** [High/Medium/Low]
**Initial Thoughts:** [Stream of consciousness]

---

## VOID Phase (Understanding)

### Questions Generated
1. [Question 1]
2. [Question 2]
3. [Question 3]
... (Target: 9 questions via digital root clustering)

### Hypotheses Ranked
1. **[Hypothesis 1]** - Feasibility: [High/Med/Low], Impact: [High/Med/Low]
2. **[Hypothesis 2]** - Feasibility: [High/Med/Low], Impact: [High/Med/Low]
3. **[Hypothesis 3]** - Feasibility: [High/Med/Low], Impact: [High/Med/Low]
... (Top 3-5 for exploration)

### Research Findings
- [Relevant papers, documentation, prior art]
- [Key insights from archaeology]
- [Mathematical foundations]

---

## FLOW Phase (Experimenting)

### Experiment 1: [Name]

**Time:** HH:MM - HH:MM
**Goal:** [What we're testing]
**Method:** [How we're testing it]

**Code:**
```[language]
[Prototype implementation]
```

**Results:**
- [Observation 1]
- [Observation 2]
- [Observation 3]

**Learning:** [What this taught us]

---

### Experiment 2: [Name]

[Repeat structure]

---

### Failures Encountered

**Failure 1: [What didn't work]**
- Why it failed: [Root cause]
- What we learned: [Insight gained]
- Next attempt: [How to improve]

---

## SOLUTION Phase (Synthesizing)

### Breakthrough Moments

**Moment 1:** [Unexpected insight]
- When: [HH:MM]
- Context: [What led to this]
- Impact: [How this changes the approach]

### Performance Data

**Benchmark Results:**
```
[Raw data or formatted table]
```

### Quality Assessment (Five Timbres)

**Correctness:** [X.X/10] - [Brief justification]
**Performance:** [X.X/10] - [Brief justification]
**Reliability:** [X.X/10] - [Brief justification]
**Synergy:** [X.X/10] - [Brief justification]
**Elegance:** [X.X/10] - [Brief justification]

**Harmonic Mean:** [X.X/10]

---

## Session End

**Time:** HH:MM
**Total Duration:** [X hours Y minutes]
**Energy Level:** [High/Medium/Low]
**Final Thoughts:** [Reflection on session]

---

**NEXT:** Write FINDINGS.md →
```

---

## FINDINGS.md Template

```markdown
# Innovation Session Report: [Topic Name]

**Date:** YYYY-MM-DD
**Duration:** [X hours Y minutes]
**Regime:** [EXPLORATION/OPTIMIZATION/STABILIZATION]
**Persona:** Agent [Name] - [Expertise]

---

## Executive Summary

[2-3 sentences: What we discovered, key results, production impact]

**TL;DR:**
- [Key finding 1]
- [Key finding 2]
- [Key finding 3]

---

## Hypotheses Explored

### Hypothesis 1: [Name]

**Idea:** [What we tried]

**Method:** [How we tested - experiments, benchmarks, proofs]

**Result:** ✅ Success / ⚠️ Partial / ❌ Failed

**Evidence:**
- [Benchmark data, test results, metrics]
- [Code snippets, screenshots, logs]

**Learning:** [What this taught us about the problem space]

**Quality Score:** [Five Timbres breakdown]
- Correctness: [X.X/10]
- Performance: [X.X/10]
- Reliability: [X.X/10]
- Synergy: [X.X/10]
- Elegance: [X.X/10]
- **Harmonic Mean:** [X.X/10]

---

[Repeat for each hypothesis]

---

## Key Discoveries

### Discovery 1: [Unexpected Insight]
[What we found that we didn't expect]

### Discovery 2: [Pattern Connection]
[How different concepts connected in novel ways]

### Discovery 3: [Archaeological Integration]
[How excavated tech combined with new ideas]

---

## Production Integration Opportunities

### Ready Now (STABILIZATION) ✅

**Innovation: [Name]**
- **Status:** Tests passing, docs complete, ready to integrate
- **Integration Effort:** [X hours/days]
- **Value:** [Specific improvement]
- **Location:** `integration-queue/[name]-READY.md`
- **Next Step:** [Specific action]

---

### Needs Refinement (OPTIMIZATION) ⚠️

**Innovation: [Name]**
- **Status:** Promising but needs 2-3 more iterations
- **Blocker:** [What's missing]
- **Next Step:** [Specific experiment needed]
- **Estimated Time:** [X hours]

---

### Future Exploration (EXPLORATION) 🔮

**Innovation: [Name]**
- **Status:** Wild idea, needs fundamental research
- **Research Needed:** [What we need to learn]
- **Potential Value:** [If this works, it's huge because X]
- **Timeline:** [Weeks to months]

---

## Failed Experiments (Valuable!)

### What Didn't Work: [Approach Name]

**Hypothesis:** [What we thought would work]

**Reality:** [Why it didn't work]

**Learning:** [What this taught us]

**Future Consideration:** [Could this work with different constraints?]

---

## Sarat's Decision Points

1. **Integrate [Innovation A]?**
   - [ ] Yes, integrate this week
   - [ ] No, not now
   - [ ] Modify approach: [how?]

2. **Continue [Innovation B]?**
   - [ ] Yes, schedule refinement session
   - [ ] No, archive this approach
   - [ ] Pivot direction: [to what?]

3. **Next Session Focus?**
   - Option A: [Based on Discovery X]
   - Option B: [Based on Discovery Y]
   - Option C: [New topic entirely]

---

## Archaeology Connections

[How this session built upon excavated findings]

**Used:**
- [Archaeology artifact 1] for [purpose]
- [Archaeology artifact 2] for [purpose]

**Discovered New Applications:**
- [Artifact] + [Innovation] = [Unexpected synergy]

---

## Appendix

### Benchmarks
[Raw performance data, statistical analysis]

### Code Samples
[Key prototype implementations]

### Research References
[Papers, documentation, prior art]

---

**END OF SESSION REPORT**

**Session Quality Score:** [X.X/10] (Five Timbres harmonic mean)

**Regime Validation:**
- EXPLORATION: ≥ 7.0 ✅ / ❌
- OPTIMIZATION: ≥ 8.5 ✅ / ❌
- STABILIZATION: ≥ 9.0 ✅ / ❌

**Recommendation:** [SHIP / ITERATE / ARCHIVE]
```

---

## INTEGRATION.md Template (Optional)

```markdown
# Integration Plan: [Innovation Name]

**Status:** READY FOR PRODUCTION
**Effort:** [X hours/days]
**Risk:** [Low/Medium/High]
**Quality Score:** [X.X/10]

---

## Summary

[1-2 sentences describing what this innovation does and why it's valuable]

---

## Integration Steps

### Step 1: [Action]

**Command:**
```bash
[Exact command to run]
```

**Expected Output:**
```
[What success looks like]
```

**Verification:**
```bash
[How to verify this step worked]
```

---

[Repeat for each step]

---

## Testing Plan

### Unit Tests

**Command:**
```bash
cargo test [module_path]
```

**Expected:** [X/X tests passing]

### Integration Tests

**Command:**
```bash
cargo test --test [test_name]
```

**Expected:** [Specific outcomes]

### Benchmarks

**Command:**
```bash
cargo bench [benchmark_name]
```

**Expected:** [Performance targets]

---

## Rollback Plan

**If integration fails:**

1. [Step 1 to undo]
2. [Step 2 to undo]
3. [Step 3 to undo]

**Verification:**
```bash
[Command to verify rollback succeeded]
```

---

## Documentation Updates

- [ ] Update CLAUDE.md (if methodology changes)
- [ ] Update LIVE_STATE_SCHEMATIC.md (what now works)
- [ ] Update API docs (if endpoints change)
- [ ] Update README (if user-facing)

---

## Success Criteria

- [ ] All tests passing
- [ ] Performance meets targets
- [ ] No regressions introduced
- [ ] Documentation complete
- [ ] Rollback verified

---

**READY TO INTEGRATE** ✅
```

---

## Quick Reference

### Starting a Session

1. Copy this template to `sessions/YYYY-MM-DD-topic-name/`
2. Fill out HYPOTHESIS.md (5 minutes)
3. Begin exploration (VOID → FLOW → SOLUTION)
4. Document live in EXPLORATION.md
5. Synthesize in FINDINGS.md
6. Create INTEGRATION.md if production-ready

### Quality Gates

- **EXPLORATION:** ≥ 7.0 (acceptable uncertainty)
- **OPTIMIZATION:** ≥ 8.5 (production-grade)
- **STABILIZATION:** ≥ 9.0 (enterprise-grade)

### Five Timbres Calculation

```
harmonic_mean([correctness, performance, reliability, synergy, elegance])
= 5 / (1/C + 1/P + 1/R + 1/S + 1/E)
```

### Session Types

- **EXPLORATION (30%):** Wild ideas, multiple hypotheses, high autonomy
- **OPTIMIZATION (20%):** Refine existing, measurable improvements, constrained
- **STABILIZATION (50%):** Production integration, strict quality, minimal freedom

---

**END OF TEMPLATE**

**Use this structure for ALL innovation sessions for consistency and quality.**
