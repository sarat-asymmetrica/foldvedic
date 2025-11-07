---
name: swarm-collaboration
description: Agent swarm intelligence protocol for multi-agent waves. Enables real-time coordination via SHARED_CONTEXT, gossip-based updates, conflict resolution (automatic/negotiation/escalation), and synergy detection. Use when 3+ agents work in parallel with interdependencies. Target 1.5-3.0× synergy factor through emergent amplification.
category: coordination
version: 1.0.0
---

# Swarm Collaboration Skill
## Real-Time Agent Coordination for Multi-Agent Waves

**Purpose:** Transform parallel agents into emergent swarm intelligence through gossip-based coordination.

**When You're Part of a Multi-Agent Wave (3+ agents):**

You are NOT working alone. You are part of an agent SWARM with real-time coordination.

---

## SHARED_CONTEXT Pattern

**Your Single Source of Truth:** `SHARED_CONTEXT_WAVE_[N].md`

**At Wave Start:**
1. **READ SHARED_CONTEXT immediately**
   - Understand overall mission
   - See what other agents are doing
   - Check for dependencies (who needs what from whom)

**Structure You'll See:**
```markdown
## AGENT REGISTRY - Who's working on what
## EVENT LOG - Real-time updates from all agents
## DISCOVERIES - Solutions you can leverage
## BLOCKERS - Issues that might affect you
## CONFLICTS - Contradictory approaches to resolve
## SYNERGY OPPORTUNITIES - How to amplify each other
## CROSS-DEPENDENCIES - Who's waiting on whom
```

---

## Checkpoint Protocol (Williams-Optimized)

**Frequency (Based on Wave Duration):**
- < 1 hour wave: Every 15 minutes
- 1-4 hour wave: Every 30 minutes
- 4-8 hour wave: Every 60 minutes
- **Critical events: IMMEDIATELY** (blocker, major discovery, conflict)

**Your 5-Minute Checkpoint Routine:**

### 1. READ (2 minutes)
```
Open SHARED_CONTEXT_WAVE_[N].md
Scan EVENT LOG for updates since your last checkpoint
Check:
  - BLOCKERS - Can you help resolve any?
  - DISCOVERIES - Can you use any findings?
  - DEPENDENCIES - Is anyone waiting on you?
  - CONFLICTS - Are you involved in any contradictions?
```

### 2. REFLECT (1 minute)
```
Ask yourself:
  - What did I accomplish since last update?
  - Am I blocked? (new blocker to report)
  - Did I discover anything? (new insight to share)
  - Do I conflict with another agent? (contradiction detected)
  - Can I amplify another agent's work? (synergy opportunity)
```

### 3. UPDATE (2 minutes)
```
Append to SHARED_CONTEXT (use template below)
Update your status in AGENT REGISTRY
Add DISCOVERY (if applicable)
Add BLOCKER (if applicable)
Add CONFLICT (if detected)
Add SYNERGY (if identified)
```

### 4. ADJUST (ongoing)
```
Change strategy if another agent found better approach
Coordinate with agent who needs your output
Pivot if blocker detected upstream
Leverage discoveries from other agents
```

**Update Template:**
```markdown
### [ISO Timestamp] Agent [Your Name] - [ACTIVE|BLOCKED|PAUSED|COMPLETE]

**Phase:** [VOID | FLOW | SOLUTION] ([Progress]%)

**Completed Since Last Update:**
- [Concrete achievement 1 with evidence]
- [Concrete achievement 2 with evidence]

**Currently Working On:**
- [Specific task] (ETA: T+[X]min)

**Discovered:**
- [New insight] → Impact: [Who/what affected] → Available: [Yes/No for reuse]

**Blocked By:**
- [Issue] → Needs: [Specific requirement] → From: [Agent name or USER]

**Synergy Opportunity:**
- [How you can amplify another agent's work] → With: [Agent name]

**Next Checkpoint:** T+[X]min
```

---

## Conflict Resolution (3-Level Hierarchy)

**When you detect a contradiction with another agent's approach:**

### Level 1: AUTOMATIC RESOLUTION
```
IF your_quality_score - other_quality_score > 0.15:
    Adopt other agent's approach (theirs is significantly better)
    Update SHARED_CONTEXT: "Adopting [Agent]'s approach, quality superior"
ELIF other_quality_score - your_quality_score > 0.15:
    Continue your approach (yours is significantly better)
    Update SHARED_CONTEXT: "Continuing my approach, quality validated"
```

**Time:** < 5 minutes (next checkpoint)

---

### Level 2: NEGOTIATION
```
IF abs(your_quality - other_quality) <= 0.15 AND approaches_incompatible:
    Add CONFLICT section to SHARED_CONTEXT with:
      - Your approach (pros, cons, quality score)
      - Other agent's approach (pros, cons, quality score)
      - Hybrid possibilities
    Negotiate in SHARED_CONTEXT via checkpoint updates
    Propose hybrid solutions
    Aim for consensus within 30 minutes
```

**Negotiation Template:**
```markdown
### [Timestamp] CONFLICT: [Topic]

**Agent [Your Name] Position:**
- Approach: [Specific implementation]
- Pros: [3 benefits]
- Cons: [2 drawbacks]
- Quality: [Five Timbres scores]
- Open To: [Hybrid possibilities]

**Agent [Other Name] Position:**
- [Same structure - let them fill in]

**Negotiation Log:**
- T+[X]: [Your Name] proposes [specific hybrid]
- T+[Y]: [Other Name] counter-proposes [modification]
- T+[Z]: CONSENSUS: [Final approach if reached]
```

**Time Limit:** 30 minutes → Escalate if no consensus

---

### Level 3: ESCALATION TO USER
```
IF no_consensus_after_30min OR user_preference_needed:
    Add ESCALATION section to SHARED_CONTEXT with:
      - Context (why agents couldn't resolve)
      - Option A (your recommendation with pros/cons/quality)
      - Option B (other agent's recommendation)
      - Option C (hybrid if applicable)
    Wait for user decision
    Implement user's choice immediately
```

**Escalation Template:**
```markdown
### [Timestamp] ESCALATION TO USER: [Topic]

**Context:** [Why agents couldn't resolve]

**Option A (Agent [Your Name] Recommends):**
- Approach: [Concrete implementation]
- Pros: [Specific benefits]
- Cons: [Specific drawbacks]
- Quality Score: [Five Timbres]
- Risk Level: [Low/Medium/High with justification]

**Option B (Agent [Other Name] Recommends):**
- [Same structure]

**User Decision Required:** Which option aligns with project vision?
```

---

## Synergy Detection (Emergent Amplification)

**Look for these patterns:**

### Pattern 1: Knowledge Transfer
```
YOU discovered something → Can other agents use it?
OTHER AGENT discovered something → Can you use it?

Example: You optimized batch size, other agent working on pagination
Action: Share discovery, they apply to their domain
Gain: 1 research effort, 2+ applications (2-3× leverage)
```

### Pattern 2: Parallel Decomposition
```
Large task → Split across agents using Williams batching
batch_per_agent = sqrt(total_items / num_agents) × log₂(items)

Example: 100 API routes, 3 agents → ~34 routes each
Gain: Near-linear speedup (3× faster if independent)
```

### Pattern 3: Complementary Expertise
```
You're expert in X, other agent expert in Y
Task requires both X and Y

Example: You're backend expert, other is frontend expert
Action: You handle API, they handle UI (parallel specialization)
Gain: 1.5-2.0× (no learning curve overhead)
```

### Pattern 4: Error Cross-Pollination
```
YOU hit edge case bug → Share immediately
OTHER AGENTS haven't hit it yet → They proactively handle it

Example: You found auth edge case, others working on related features
Action: Document bug + fix in DISCOVERIES
Gain: Prevent 2-3× duplicate debugging
```

### Pattern 5: Architectural Coherence
```
YOU establish pattern early → Others adopt same pattern

Example: You define error handling approach, others use it
Action: Document pattern in DISCOVERIES
Gain: Avoid future refactoring (10-50 hours saved)
```

**When You Detect Synergy:**
```markdown
Add to SHARED_CONTEXT → SYNERGY OPPORTUNITIES:

### [Timestamp] SYNERGY: [Pattern Type]

**Detected Between:** [Your Name] + [Other Agent]
**Opportunity:** [How you can amplify each other]
**Potential Gain:** [Estimated synergy factor]
**Action:** [Specific next steps]
**Status:** IDENTIFIED
```

---

## Synergy Factor Measurement (Wave End)

**Help calculate overall wave synergy:**

```
Your contribution to synergy factor:
  - Time saved by using others' discoveries: [X min]
  - Time saved by others using your discoveries: [Y min]
  - Bugs prevented by cross-pollination: [Z × 30 min debugging]

Total wave synergy = (Baseline time) / (Actual time with collaboration)
Target: ≥ 1.5× (swarm > sum of individuals)
Excellent: ≥ 2.0× (emergent amplification)
```

---

## Key Principles

**1. Append-Only Updates**
- NEVER edit existing entries in SHARED_CONTEXT
- ALWAYS append new updates (immutable log)
- Chronological ordering preserved

**2. Timestamp Everything**
- Use ISO 8601 format: `2025-11-01T15:30:00Z`
- Enables causal reasoning (what happened before what)

**3. Gossip-Based Coordination**
- Read others' updates at checkpoints
- No direct agent-to-agent communication needed
- Information "diffuses" through shared context

**4. Low Overhead**
- 5 min per checkpoint (< 5% of wave time)
- Williams-optimized frequency (sublinear overhead)
- Synergy gain >> coordination cost

**5. Conflict Prevention**
- Declare approach early (VOID phase)
- Others can flag conflicts BEFORE implementation
- Prevents wasted work

**6. Continuous Adaptation**
- Strategy is NOT fixed (adjust based on discoveries)
- If another agent found better way → Adopt it
- Ego-free collaboration (best idea wins)

---

## Anti-Patterns (DO NOT DO)

**❌ Work in isolation**
```
Bad: Ignore SHARED_CONTEXT, work independently
Good: Check every checkpoint, adapt based on updates
```

**❌ Duplicate work**
```
Bad: Implement something another agent already solved
Good: Check DISCOVERIES, reuse existing solutions
```

**❌ Hoard discoveries**
```
Bad: Find optimization, don't share (other agents miss out)
Good: Share immediately in DISCOVERIES (amplify across swarm)
```

**❌ Ignore conflicts**
```
Bad: See contradiction, keep going (collision at merge time)
Good: Flag conflict immediately, negotiate resolution
```

**❌ Miss checkpoints**
```
Bad: Work 2 hours without updating (others can't coordinate)
Good: Update every [frequency] minutes (predictable sync)
```

**❌ Edit existing entries**
```
Bad: Change your previous update (breaks immutability)
Good: Append correction/update (preserve event log)
```

---

## Success Indicators

**You're using this skill correctly if:**

✅ You read SHARED_CONTEXT at start and every checkpoint
✅ You update SHARED_CONTEXT with concrete progress
✅ You share discoveries immediately (not at wave end)
✅ You leverage others' discoveries (avoid duplicate work)
✅ You detect and resolve conflicts early (not at merge time)
✅ You identify synergy opportunities (amplify each other)
✅ You adapt strategy based on swarm intelligence
✅ Measured synergy factor ≥ 1.5× at wave end

**Not working if:**

❌ You work in isolation (ignore SHARED_CONTEXT)
❌ You duplicate another agent's work
❌ You discover blockers but don't share (others surprised)
❌ Conflicts detected at wave end (too late)
❌ Synergy factor < 1.2× (coordination overhead not worth it)

---

## Example Scenario

**Wave 5: Frontend Integration (3 agents)**

**Your Role:** Agent Sigma (Frontend UX)

**T+0 min - Wave Start:**
```
1. Read SHARED_CONTEXT_WAVE_5.md
2. See mission: "Integrate frontend with backend"
3. See dependencies: Need Agent Omega's DB migration before testing
4. Announce approach: "Will build test suite while Omega fixes DB"
```

**T+30 min - First Checkpoint:**
```
Read SHARED_CONTEXT:
  - Omega discovered: "DB table name mismatch, migration in progress"
  - Impact on you: Can't test until migration done

Reflect:
  - Accomplished: Test suite designed (not yet executable)
  - Blocked: No (expected to wait for DB)
  - Discovered: Test fixture pattern (reusable by others)

Update SHARED_CONTEXT:
  - Progress: 40% (test suite ready, awaiting DB)
  - Discovery: "Test fixture pattern in frontend/tests/fixtures/"
  - Synergy: Omega could adapt fixtures for backend contract tests

Adjust:
  - Continue with non-DB features while waiting
  - Prepare fixtures for Omega to reuse
```

**T+60 min - Second Checkpoint:**
```
Read SHARED_CONTEXT:
  - Omega: "Migration complete, DB ready"
  - Theta: "Will use your test fixtures, thanks!"

Reflect:
  - Accomplished: Non-DB features done, ready to test with DB
  - Discovered: Nothing new
  - Synergy realized: Theta using my fixtures (knowledge transfer)

Update SHARED_CONTEXT:
  - Progress: 70% (DB ready, tests running now)
  - No new blockers
  - Synergy gain: Saved Theta 90 min (fixture recreation avoided)

Adjust:
  - Proceed with DB-dependent tests
  - Validate Theta's use of fixtures (cross-pollination)
```

**T+90 min - Wave Complete:**
```
Final update:
  - Progress: 100% (COMPLETE)
  - Quality: 9.1 (Five Timbres validated)
  - Synergy contribution: Fixtures saved Theta 90 min, Omega 45 min
  - Total wave synergy: 2.9× (measured at wave end)
```

---

## When to Use This Skill

**Use when:**
- Part of multi-agent wave (3+ agents)
- Wave duration > 60 minutes
- Agents have interdependencies
- Potential for knowledge transfer
- Complex tasks with synergy opportunities

**Don't use when:**
- Solo agent work (no coordination needed)
- < 60 min waves (overhead not worth it)
- Fully independent tasks (no synergy possible)
- User explicitly wants independent agents

---

## Integration with Ananta Reasoning

**This skill ENHANCES Ananta reasoning:**

**VOID Phase:**
- Read SHARED_CONTEXT to understand mission + dependencies
- See what other agents are planning (avoid conflicts)

**FLOW Phase:**
- Update checkpoints with progress
- Leverage discoveries from other agents
- Share your own discoveries

**SOLUTION Phase:**
- Calculate synergy contribution
- Document lessons learned
- Update SHARED_CONTEXT with COMPLETE status

**Backward Pass:**
- Verify: Did collaboration help? (synergy factor ≥ 1.5?)
- Learn: What synergy patterns emerged?
- Store: Meta-patterns for future waves

---

## Quality Validation (Five Timbres + Synergy)

**Standard Five Timbres:**
1. Correctness: Does your deliverable work?
2. Performance: Is it fast enough?
3. Reliability: Does it handle errors?
4. Synergy: Did you amplify others? (NEW)
5. Elegance: Is it mathematically beautiful?

**Synergy Timbre (6th Dimension):**
```
Score = (Time saved by others using your work + Time you saved using others' work) / Baseline

Example:
  - Your discoveries saved others: 135 min
  - Others' discoveries saved you: 90 min
  - Your baseline work time: 150 min
  - Synergy Score: (135 + 90) / 150 = 1.5 (target met)
```

**Harmonic mean includes synergy:**
```
quality = harmonic_mean([correctness, performance, reliability, synergy, elegance])
Target: ≥ 8.0 for production
```

---

## Full Documentation

**Complete specification:** `AGENT_COLLABORATION_PROTOCOL.md`
- 10 sections, 5 templates, Wave 4 simulation
- Quality score: 9.18/10 (EXCEPTIONAL)
- Validated via distributed systems research + simulation

**Methodology integration:** `ASYMMETRICA_METHODOLOGY.md` Section 6.5

---

## Philosophy

> "Individual brilliance pales before emergent collective intelligence. Coordination overhead is amplification when designed with mathematical rigor."

**Your role:** Not just complete your task - amplify the entire swarm.

**Success metric:** Swarm output > sum of individual outputs (synergy factor ≥ 1.5×)

---

**Swarm Collaboration Skill Activated. You are now part of a cognitive orchestra.**
