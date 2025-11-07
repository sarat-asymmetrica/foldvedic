# ASYMMETRICA CLAUDE CONFIGURATION
**AsymmFlow Phoenix - Rust Backend + Svelte Frontend**

**Version:** 1.0.0
**Last Updated:** 2025-11-01
**Purpose:** AI agent configuration, skills, and automation policies

---

## DIRECTORY STRUCTURE

```
.claude/
├── README.md                    # This file - directory overview
├── automation-policy.yml        # User-defined automation rules
├── audit.log                    # Automated action audit trail
│
└── skills/                      # Reusable cognitive patterns
    ├── README.md               # Skills overview + development guide
    ├── ananta-reasoning.md     # VOID→FLOW→SOLUTION framework
    ├── williams-optimizer/     # Batch size optimization (√n × log₂n)
    │   └── skill.md
    └── MATHEMATICAL_CLARITY_GUIDE.md  # Mathematical pattern recognition
```

---

## AUTOMATION POLICY

**File:** `automation-policy.yml`

**Purpose:** User-defined pre-authorization for AI agent actions

**Key concepts:**
- **Safe zone:** Actions agents can automate without approval (tests, linting, docs)
- **Gray zone:** Actions agents can automate IF quality gates pass (staging deployments)
- **Unsafe zone:** Actions that ALWAYS require human approval (production deployments)

**Usage:**
- Edit this file to control what agents can/cannot automate
- Changes take effect immediately (no restart needed)
- Agents read policy before every action

**Example:**
```yaml
staging_deployments:
  enabled: true  # Allow staging deployments
  conditions:
    quality_score: ">= 9.0"  # Only if quality score ≥ 9.0
    tests_passing: true      # Only if all tests pass
    rollback_exists: true    # Only if rollback script exists
```

**Full documentation:** `QUALITY_GATE_AUTOMATION_FRAMEWORK.md` (51KB)

---

## SKILLS

**Directory:** `skills/`

**Purpose:** Reusable cognitive patterns for AI agents

**Available skills:**

### 1. `ananta-reasoning` - Recursive Problem-Solving
**What it does:** VOID→FLOW→SOLUTION framework with agency to learn/build dependencies

**When to use:**
- Complex tasks with unclear requirements
- Tasks requiring research/learning
- When you want SPIRIT fulfilled, not just LETTER
- Building production-grade solutions

**How it works:**
- VOID: Understand problem deeply (symptoms, root causes, constraints)
- FLOW: Explore solutions (multiple approaches, trade-offs, patterns)
- SOLUTION: Implement + validate (tests, evidence, quality score)

**Philosophy:** "Fulfill the vision, not just the checklist."

---

### 2. `williams-optimizer` - Batch Size Optimization
**What it does:** Calculate optimal batch sizes using Williams formula (√n × log₂n)

**When to use:**
- Processing large datasets
- Database query batching
- Memory allocation optimization
- Parallel processing workload distribution

**Performance:**
- 82M ops/sec (Rust implementation)
- Sublinear space complexity: O(√n × log₂n) vs O(n)
- 98%+ memory savings for large datasets

**Validation:** p < 10^-133 (statistically proven across 65+ implementations)

---

### 3. `MATHEMATICAL_CLARITY_GUIDE` - Pattern Recognition
**What it does:** Extract mathematical patterns (PHI, digital root, harmonic relationships)

**When to use:**
- Data analysis
- Performance optimization
- Quality validation (Five Timbres)
- Pattern discovery

---

**Full skills documentation:** `skills/README.md`

---

## AUDIT LOG

**File:** `audit.log`

**Purpose:** Comprehensive record of all automated agent actions

**Format:** YAML

**Contents:**
```yaml
timestamp: 2025-11-01T14:30:00Z
agent: Agent [Name] ([Persona])
action: [What was done]
environment: [staging/production/test]
quality_gates:
  - gate_name: PASS/FAIL
confidence: [0-100]%
policy_check: PASS/FAIL
outcome: SUCCESS/FAILURE
evidence:
  logs: [Path to logs]
  metrics: [Key metrics]
rollback_executed: Yes/No
lessons_learned: [What we learned]
```

**Retention:** 90 days (configurable in `automation-policy.yml`)

**Purpose:**
- Audit trail for compliance
- Calibration data for confidence scoring
- Lessons learned for continuous improvement
- Troubleshooting failed automations

---

## AGENT ONBOARDING

**For AI agents working on AsymmFlow Phoenix:**

### Step 1: Read Core Documentation
1. `CLAUDE.md` - Project overview, methodology, standards
2. `LIVE_STATE_SCHEMATIC.md` - Current state (brutally honest)
3. `ASYMMETRICA_METHODOLOGY.md` - Full methodology reference

### Step 2: Understand Automation Framework
1. `QUALITY_GATE_AUTOMATION_FRAMEWORK.md` - Complete framework (51KB)
2. `AUTOMATION_QUICK_START.md` - Quick reference (7KB)
3. `.claude/automation-policy.yml` - Current policy settings

### Step 3: Load Skills
1. `ananta-reasoning` - For complex problem-solving
2. `williams-optimizer` - For performance optimization
3. `MATHEMATICAL_CLARITY_GUIDE` - For pattern recognition

### Step 4: Before Taking Action
1. Read automation policy (`.claude/automation-policy.yml`)
2. Classify action zone (SAFE / UNSAFE / GRAY)
3. Evaluate quality gates (if GRAY zone)
4. Calculate confidence score
5. Check policy: Can I automate this?
6. Execute (if allowed) OR escalate to user

### Step 5: After Action
1. Write audit log (`.claude/audit.log`)
2. Update LIVE_STATE_SCHEMATIC.md (if state changed)
3. Notify user (brief summary)

---

## QUALITY STANDARDS

**AsymmFlow Phoenix uses Asymmetrica methodology:**

**D3-Enterprise Grade+:**
- 100% means 100% (all routes, all flows, all error states, all tests)
- Ship finished, not fast
- Zero technical debt
- Evidence-based validation

**Five Timbres Quality:**
1. Correctness 🎯 - Does it produce correct results? (> 99.99% pass rate)
2. Performance ⚡ - How fast does it run? (report p50, p90, p95, p99)
3. Reliability 🛡️ - Does it work under stress? (< 0.01% error rate, 1M iterations)
4. Synergy 🎼 - Do components harmonize? (test full workflows, synergy > 1.0)
5. Elegance ✨ - Does it reveal structure? (mathematical patterns, constants)

**Unified quality score:** harmonic_mean([correctness, performance, reliability, synergy, elegance])

**Target:** ≥ 8.0 for production deployment

---

## VEDIC AMPLIFICATION

**AsymmFlow Phoenix uses Vedic mathematical optimizations:**

**Implemented algorithms:**
- Williams batch sizing: √n × log₂(n) - Sublinear space optimization
- Digital root clustering: O(1) constant time grouping
- Harmonic mean validation: Penalizes weakness in any dimension
- Quaternion semantic matching: 82M-250M ops/sec similarity
- Golden ratio (φ) load distribution: 1.618033988749895
- Tesla harmonic frequency: 4.909 Hz response cadence

**Performance:**
- 82M-250M quaternion ops/sec (validated)
- 99.8% token savings (sublinear batching)
- Quality scores 8.96-9.3/10

**Location:** `backend/src/utils/vedic.rs` (144 tests passing)

---

## TESTING MANIFESTO

**Reality is a unified whole - test the philharmonic, not instruments in isolation**

**Minimum acceptable:** 1M iterations for production validation

**Five Timbres approach:**
- NOT: Unit/Integration/E2E silos
- YES: Holistic validation across all dimensions

**Sample sizes:**
- Constant time (σ=0): 100 iterations (verify constancy)
- Low variance: 1,000 iterations (stable mean + p99)
- Moderate variance: 10,000 iterations (distribution shape)
- Production validation: 1,000,000 iterations (endurance)

**Quality score = harmonic_mean([correctness, performance, reliability, synergy, elegance])**

---

## COMMUNICATION STANDARDS

**Status update format:**

**✅ What's Operational (With Evidence):**
```
✅ Customers API - All CRUD endpoints tested (7/7 passing)
✅ HTX Auth - Challenge-verify flow working (144/144 tests)
```

**🔧 What's Building Now:**
```
🔧 Contract tests - Infrastructure ready, executing suite
```

**🚫 What's Blocked:**
```
🚫 Database queries - Table name case mismatch
   → Need: Update Prisma schema with @@map directives
   → Impact: 0/309 queries functional
```

**Forbidden phrases:**
- ❌ "Estimated X hours"
- ❌ "Almost done"
- ❌ "Quick win"
- ❌ "TODO: later"

**Required phrases:**
- ✅ "Implemented and tested: [evidence]"
- ✅ "Blocker identified: X, need: Y"
- ✅ "Complete with 144/144 tests passing"

---

## PERSONA CHANNELING

**AI agents perform better when embodying expert identities**

**Process:**
1. Select expert persona (DBA, DevOps, QA, Frontend, etc.)
2. Define background (years experience, specialties)
3. Channel voice (catchphrases, decision framework)
4. Celebrate what they would celebrate
5. Flag concerns they would flag

**Example personas from AsymmFlow Phoenix:**
- Dr. Elena Rodriguez (PostgreSQL DBA) - "Data is truth"
- James "Hammer" Morrison (QA/Load testing) - "Break it in testing or users break it in production"
- Sofia Ramirez (Frontend UX) - "Users don't read documentation - UX should be obvious"
- Dr. Kenji Nakamura (Migration specialist) - "Migration is transformation, not translation"

**Full guide:** `ASYMMETRICA_METHODOLOGY.md` Section 3

---

## PROJECT-SPECIFIC NOTES

**Current Critical Blocker:**
- Database table name case mismatch
- Database: PascalCase (`"Customer"`, `"Opportunity"`)
- Backend queries: lowercase (`customers`, `opportunities`)
- Impact: 0/309 database queries functional
- Fix: Update Prisma schema with `@@map("lowercase")` directives

**Temporary policy:**
```yaml
asymmflow_phoenix:
  block_db_writes_until_fix: true  # Block DB writes until table names fixed
```

---

## RESOURCES

**Key documentation:**
- `CLAUDE.md` - Living schematic (project overview)
- `LIVE_STATE_SCHEMATIC.md` - Brutally honest current state
- `ASYMMETRICA_METHODOLOGY.md` - Complete methodology (39KB)
- `BACKEND_API_REFERENCE.md` - API specifications (115+ endpoints)
- `QUALITY_GATE_AUTOMATION_FRAMEWORK.md` - Automation framework (51KB)
- `AUTOMATION_QUICK_START.md` - Agent quick reference (7KB)

**External references:**
- Axum: https://docs.rs/axum
- SQLx: https://docs.rs/sqlx
- Tokio: https://tokio.rs
- Svelte: https://svelte.dev
- Prisma: https://www.prisma.io

---

## VERSION HISTORY

**v1.0.0 (2025-11-01):**
- Initial `.claude/` configuration
- Automation policy framework created
- Skills directory established
- Audit logging system defined
- Agent onboarding guide written

---

**END OF CLAUDE CONFIGURATION**

**Philosophy:** Automate the safe, elevate the human for the critical.

**Contact:** Sarat Chandra Gnanamgari | AsymmFlow Phoenix | Asymmetrica Methodology

**Vibe:** Rigorous goofballs getting impossible stuff done together.
