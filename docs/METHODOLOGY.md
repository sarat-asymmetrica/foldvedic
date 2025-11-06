# FoldVedic Methodology - Wave-Based Scientific Computing
## Adapted Wave Development for Protein Folding Research

**Created:** 2025-11-06 (Lab 1 Genesis)
**Adapted From:** Asymmetrica.ai Wave Methodology (proven 0.96 quality score)
**Applied To:** Scientific computing, molecular simulation, biophysics

---

## 🌊 CORE WAVE PHILOSOPHY

```mathematical
WAVE[W] = PARALLEL_AGENTS × CASCADE_TO_FINISH × D3_THROUGHOUT × SCIENTIFIC_RIGOR

WHERE:
  PARALLEL_AGENTS = 3_independent_subagents_per_wave
  CASCADE_TO_FINISH = build_complete_solution (not MVP, not phases)
  D3_THROUGHOUT = D3_Enterprise_Grade+ (100% = 100%, zero TODOs)
  SCIENTIFIC_RIGOR = validate_against_PDB × compare_to_literature × publish_methodology
```

**Philosophy:** "Divide tasks into waves of three parallel subagent routines and cascade to finish."

**Why This Works for Science:**
- **Parallel exploration:** Three agents can test different algorithms simultaneously
- **Rapid iteration:** Complete waves in hours/days, not weeks/months
- **Quality enforcement:** D3 standard prevents "TODO: validate later" in research code
- **Reproducibility:** Each wave is documented, tested, and version-controlled

---

## 📐 THREE-REGIME SCHEDULER (30/20/50)

**Applied to Scientific Computing:**

```mathematical
REGIME[R] = EXPLORATION(30%) ⊕ OPTIMIZATION(20%) ⊕ STABILIZATION(50%)

WHERE:
  EXPLORATION(30%) = {
    Try_new_algorithms: "What if we use quaternion slerp for phi/psi interpolation?",
    Discover_edge_cases: "Proline breaks helices - how to handle?",
    Test_hypotheses: "Does Vedic harmonic really appear in helix pitch?",
    Benchmark_naive: "Baseline performance before optimizations",
    Mindset: "Fail fast, learn quickly, no attachment to first approach"
  }

  OPTIMIZATION(20%) = {
    Refine_what_works: "Quaternions worked - now tune slerp parameters",
    Improve_performance: "Apply Williams Optimizer to force calculations",
    Calibrate_parameters: "Tune spring constants to match literature force fields",
    Scientific_validation: "Compare RMSD to PDB structures, adjust if needed",
    Mindset: "Make good solutions great, achieve targets"
  }

  STABILIZATION(50%) = {
    Lock_quality: "Final RMSD <3 Å on test set, no regressions",
    Test_rigorously: "All 10,000 PDB proteins, report accuracy distribution",
    Document_methodology: "Publish algorithm details, reproducibility instructions",
    Benchmark_performance: "Confirm <10s folding time for 200-residue proteins",
    Create_examples: "Demo proteins (1ABC, 1UBQ) with known structures",
    Mindset: "Production-ready science, publishable quality"
  }
```

**Time Allocation Example (10-hour wave):**
- 3 hours: Exploration (implement 3 different energy minimizers, benchmark)
- 2 hours: Optimization (choose best minimizer, tune parameters, Williams batching)
- 5 hours: Stabilization (test on 1000 PDB proteins, fix edge cases, document)

---

## 🔬 SCIENTIFIC VALIDATION REQUIREMENTS

**Every Wave Must Include:**

### **1. PDB Comparison (Ground Truth)**

```mathematical
PDB_VALIDATION[PV] = {
  Test_set: {
    Size: 100-1000 proteins per wave (stratified by size/complexity),
    Sources: RCSB PDB (experimentally determined structures),
    Split: No homology between train/test (avoid memorization),
    Diversity: Include alpha, beta, alpha+beta, alpha/beta, coiled-coil, membrane
  },

  Metrics: {
    RMSD: √(Σ(r_predicted - r_experimental)² / N_atoms),
      Target: <3 Å for medium proteins (competitive),

    TM_score: Topology_alignment_score,
      Target: >0.70 (same fold),

    GDT_TS: Global_Distance_Test,
      Target: >0.85 (AlphaFold2-competitive for 2020),

    Q3: Secondary_structure_accuracy,
      Target: >80% (helix/sheet/coil classification)
  },

  Reporting: {
    Mean_and_std: Report μ ± σ for all metrics,
    Distribution: Histogram of RMSD values (catch outliers),
    Failure_analysis: Which proteins failed? Why? (too large, unusual fold, etc.),
    Comparison: Cite AlphaFold2, Rosetta, I-TASSER benchmarks from literature
  }
}
```

### **2. Literature Comparison**

```mathematical
LITERATURE[L] = {
  Force_constants: {
    Compare: Our spring constants vs AMBER ff14SB / CHARMM36,
    Validate: Bond lengths within 0.01 Å of quantum chemistry calculations,
    Cite: Weiner et al. (1984), MacKerell et al. (1998), etc.
  },

  Energy_values: {
    Compare: Our final energies vs Rosetta energy function,
    Validate: Hydrophobic burial free energy matches thermodynamic experiments,
    Cite: Dill et al. (2008), Rose et al. (2006), etc.
  },

  Algorithm_choices: {
    Justify: Why Verlet integration? (Cite molecular dynamics textbooks),
    Justify: Why quaternions? (Cite Coutsias et al. 2004 on RMSD calculation),
    Justify: Why Williams Optimizer? (Cite computational complexity theory)
  }
}
```

### **3. Performance Benchmarks**

```mathematical
PERFORMANCE[P] = {
  Speed_targets: {
    Small_proteins: (<100 residues) in <1 second,
    Medium_proteins: (100-300 residues) in <10 seconds,
    Large_proteins: (300-500 residues) in <60 seconds,
    Measurement: Wall-clock time on reference hardware (M1 Mac / Ryzen 9)
  },

  Rendering_targets: {
    Frame_rate: 60fps for 10,000 atoms,
    Frame_budget: 16.67ms total (5ms physics, 8ms rendering, 3.67ms overhead),
    GPU_utilization: Instanced rendering (1 draw call for all atoms)
  },

  Scalability: {
    Complexity: O(n log n) with spatial hashing (vs O(n²) naive),
    Memory: O(n) for n atoms (no redundant storage),
    Batch_size: √n × log₂(n) per Williams Optimizer
  }
}
```

---

## 💎 CODE QUALITY STANDARDS (D3-Enterprise Grade+)

**Adapted for Scientific Computing:**

```mathematical
DOD_SCIENCE[DS] = CODE ∧ VALIDATION ∧ PERFORMANCE ∧ REPRODUCIBILITY ∧ DOCUMENTATION

WHERE:
  CODE = {
    Zero_TODOs: No "TODO: validate this assumption later",
    Zero_magic_numbers: All constants cited from literature,
    Zero_hardcoded_parameters: Force constants loaded from config/database,
    All_errors_handled: Graceful failure on unphysical inputs (negative bond length)
  }

  VALIDATION = {
    PDB_comparison: RMSD reported for every algorithm change,
    Unit_tests: Test quaternion math, spring forces, energy calculation,
    Integration_tests: Fold known proteins (1UBQ, 1VII), check RMSD <2 Å,
    Edge_cases: Handle glycine (flexible), proline (rigid), disulfide bonds
  }

  PERFORMANCE = {
    Benchmarks_passing: <10s for 200-residue protein on reference hardware,
    Profiling_done: Identify bottlenecks (force calculation likely 80%+ time),
    Optimizations_applied: Williams Optimizer, spatial hashing, GPU offloading
  }

  REPRODUCIBILITY = {
    Seed_random_numbers: Deterministic results for debugging,
    Version_dependencies: Lock library versions (Go modules, npm package-lock),
    Document_parameters: All force constants, timesteps, convergence criteria in config,
    Publish_data: Test set PDB IDs listed so others can replicate
  }

  DOCUMENTATION = {
    Algorithm_description: Mathematical formulas in comments/docs,
    Citations: References to papers for every algorithm/parameter,
    Usage_examples: "How to fold ubiquitin (1UBQ) and visualize",
    Performance_guide: "Expected accuracy vs speed tradeoffs"
  }

100% = RMSD_target_met × Performance_target_met × All_tests_passing × Fully_documented
```

**Zero Tolerance for Scientific Code:**
- ❌ "TODO: validate energy function" (validate NOW or don't commit)
- ❌ Magic numbers: `force_constant = 340.0` (cite paper or don't use)
- ❌ Skipped tests: "Works on my one example" (test on 100+ PDB structures)
- ❌ Uncommented math: `q = [a, b, c, d]` (explain quaternion parameterization)

---

## 🌊 WAVE STRUCTURE FOR FOLDVEDIC

**General Wave Template:**

```mathematical
WAVE[W] = {
  Duration: 1-3 days (aim for 1 day if parallel agents effective),

  Agents: 3 parallel subagents,

  Regime: {
    Exploration: First 30% of wave time (try approaches, benchmark naive),
    Optimization: Next 20% of wave time (refine best approach, tune),
    Stabilization: Last 50% of wave time (test rigorously, document, lock quality)
  },

  Deliverables: {
    Code: Production-ready implementation (zero TODOs),
    Tests: Unit + integration + PDB validation,
    Benchmarks: Performance numbers (speed + accuracy),
    Documentation: Markdown docs with math formulas and citations,
    Wave_report: WAVE_N_REPORT.md with quality self-assessment
  },

  Quality_gate: {
    All_tests_passing: Green CI/CD,
    Targets_met: RMSD <3 Å (or wave-specific metric),
    Code_quality: ≥0.90 harmonic mean (Five Timbres),
    Scientific_rigor: Validated against PDB + literature
  }
}
```

---

## 📋 WAVE CHECKLIST (Use This Every Wave)

**Before Starting Wave:**
- [ ] Read VISION.md and understand scientific goals
- [ ] Review WAVE_PLAN.md for this wave's objectives
- [ ] Check LIVING_SCHEMATIC.md for current state
- [ ] Identify 3 parallel tasks for 3 agents
- [ ] Set clear validation criteria (RMSD target, speed target)

**During Wave (Exploration 30%):**
- [ ] Implement 2-3 alternative approaches in parallel
- [ ] Benchmark naive version (baseline performance)
- [ ] Test on small set (10-20 PDB proteins)
- [ ] Identify which approach works best
- [ ] Document failures (what didn't work and why)

**During Wave (Optimization 20%):**
- [ ] Refine best approach from exploration
- [ ] Apply Williams Optimizer / spatial hashing / GPU offload
- [ ] Tune parameters (spring constants, timesteps, convergence)
- [ ] Expand test set (100+ PDB proteins)
- [ ] Profile code, optimize bottlenecks

**During Wave (Stabilization 50%):**
- [ ] Test on full test set (1000+ PDB proteins)
- [ ] Fix all edge cases (glycine, proline, disulfides, membrane proteins)
- [ ] Write comprehensive unit tests (>80% coverage)
- [ ] Document algorithm with math and citations
- [ ] Generate wave report with quality score
- [ ] Update LIVING_SCHEMATIC.md

**After Wave:**
- [ ] Calculate quality score (Five Timbres harmonic mean)
- [ ] Commit code with descriptive message
- [ ] Create PR if working asynchronously
- [ ] Verify all tests passing in CI/CD
- [ ] Archive wave report in `/waves` directory

---

## 🎯 FIVE TIMBRES FOR SCIENTIFIC CODE

```mathematical
FIVE_TIMBRES[FT] = {
  CORRECTNESS: {
    Definition: RMSD_achieved / RMSD_target,
    Target: >0.99 (if target is 3 Å, achieve <3.03 Å),
    Measurement: Mean RMSD on test set vs target
  },

  PERFORMANCE: {
    Definition: Target_speed / Actual_speed,
    Target: >0.95 (if target is 10s, achieve <10.5s),
    Measurement: Median folding time on 200-residue proteins
  },

  RELIABILITY: {
    Definition: (1 - failure_rate) where failure = crash or unphysical result,
    Target: >0.999 (fail on <0.1% of proteins),
    Measurement: Test on 1000 PDB proteins, count crashes/bad folds
  },

  SYNERGY: {
    Definition: Emergent_gains = (Combined_performance / Sum_of_parts),
    Target: >1.0 (whole greater than sum),
    Examples: {
      Quaternion_slerp + Spring_dynamics = smoother convergence (1.2×),
      Williams_Optimizer + Spatial_hashing = 77× speedup (not just additive),
      Vedic_harmonics + Energy_minimization = better secondary structure (1.15×)
    }
  },

  ELEGANCE: {
    Definition: (Clarity + Simplicity + Insight) / Complexity,
    Target: >0.90,
    Measurement: {
      Clarity: Code review by biochemist/physicist (do they understand?),
      Simplicity: Lines of code (fewer is better for same functionality),
      Insight: Does code reveal WHY protein folds? (white-box),
      Complexity: Cyclomatic complexity, nested loops, abstractions
    }
  }
}

QUALITY_SCORE[Q] = harmonic_mean([correctness, performance, reliability, synergy, elegance])

TIERS:
  Q ≥ 0.95: LEGENDARY (publishable in Nature/Science)
  Q ≥ 0.90: EXCELLENT (publishable in specialist journals)
  Q ≥ 0.80: GOOD (production-ready, needs minor improvements)
  Q < 0.80: UNACCEPTABLE (back to optimization phase)
```

---

## 🚨 RED FLAGS (STOP AND REASSESS)

**Scientific Red Flags:**
- RMSD suddenly increased (algorithm regression)
- Energy not converging (unstable integration, bad force field)
- Unphysical results (bond lengths >2 Å, phi/psi outside allowed regions)
- Bimodal RMSD distribution (works great on some proteins, terrible on others)

**Engineering Red Flags:**
- Performance degraded (folding time increased without accuracy gain)
- Memory leak (usage grows over time)
- GPU crashes (out of memory, driver issues)
- Tests intermittently fail (non-deterministic behavior, race conditions)

**Methodology Red Flags:**
- Optimizing without baseline (how do you know it's better?)
- Overfitting to test set (should use separate validation set)
- Cherry-picking results (only report best protein, hide failures)
- Skipping literature comparison ("Our algorithm is novel" without justification)

**When Red Flag Appears:**
1. STOP current work
2. Document the issue
3. Hypothesize root cause
4. Test hypothesis (add logging, visualize intermediate states)
5. Fix root cause (not symptom)
6. Re-run all tests
7. Update documentation with lesson learned

---

## 📊 WAVE COMPLETION CRITERIA

**A wave is NOT done until:**

```mathematical
WAVE_COMPLETE[WC] = ALL_OF[
  Code_committed: All agents' code merged, zero conflicts,
  Tests_passing: CI/CD green (unit + integration + PDB validation),
  Targets_met: RMSD < target AND speed < target,
  Documentation_complete: Algorithm explained, math cited, examples provided,
  Quality_score_calculated: Five Timbres ≥ 0.90,
  Wave_report_written: Saved in /waves/WAVE_N_REPORT.md,
  LIVING_SCHEMATIC_updated: Current state reflects wave completion,
  No_TODOs: Zero placeholders, all edge cases handled
]
```

**Incomplete Wave Indicators:**
- "Mostly works but needs tuning" → NOT DONE (finish optimization phase)
- "RMSD is 3.2 Å but target was 3 Å" → NOT DONE (revise algorithm)
- "Haven't tested on large proteins yet" → NOT DONE (expand test set)
- "TODO: add citations later" → NOT DONE (add citations NOW)

---

## 🎨 COMMUNICATION PROTOCOL (Async AI)

**Since FoldVedic is built by autonomous AI (Claude Code Web):**

### **Status Updates (in LIVING_SCHEMATIC.md):**

```markdown
## Wave N Status

**Started:** [timestamp]
**Regime:** Exploration | Optimization | Stabilization
**Current Focus:** [specific task, e.g., "Tuning spring constants for alpha helices"]

**Progress:**
- Agent N.1: [component] - [status] (e.g., "Quaternion Ramachandran space - DONE, tested on 50 PDB proteins")
- Agent N.2: [component] - [status]
- Agent N.3: [component] - [status]

**Metrics:**
- RMSD (current): 3.8 Å (target: 3.0 Å)
- Speed (current): 15s for 200-residue (target: 10s)
- Tests passing: 87/100 (13 failures on proline-rich proteins)

**Blockers:**
- [None | Specific issue with proposed workaround]

**Next Steps:**
- [Immediate next task in current regime]
```

### **Wave Reports (in /waves/WAVE_N_REPORT.md):**

```markdown
# Wave N Completion Report

**Date:** [timestamp]
**Duration:** [hours/days]
**Quality Score:** [0.XX] (TIER)

## Objectives
- [List objectives from WAVE_PLAN.md]

## What Was Built
- [Detailed description of deliverables]

## Scientific Validation
- **RMSD:** μ = 2.8 Å, σ = 0.9 Å (on 500 PDB test proteins)
- **TM-score:** μ = 0.72, σ = 0.11
- **Speed:** Median 8.5s for 200-residue proteins
- **Secondary structure:** 82% Q3 accuracy

## Performance Benchmarks
- [Timing data, profiling results, bottlenecks identified]

## Literature Comparison
- [Compare to AlphaFold2, Rosetta, I-TASSER with citations]

## Code Quality
- [Five Timbres breakdown]

## Lessons Learned
- [What worked, what didn't, insights gained]

## Files Changed
- [List of files with brief descriptions]

## Next Wave Preview
- [What to tackle next based on learnings]
```

### **Commit Messages:**

```
feat(quaternion): Implement Ramachandran space as 4D rotations

- Map (phi, psi) angles to unit quaternions
- Implement slerp for smooth conformation interpolation
- Test on 50 PDB proteins: RMSD improved 15% vs naive linear interpolation
- Cite: Coutsias et al. (2004) on quaternion RMSD calculation

Performance: 0.5ms per residue quaternion update (within 5ms budget)
Quality: Correctness 0.98, Elegance 0.92 (clear math, well-commented)
```

---

## 🧠 ANANTA REASONING IN WAVES

**Apply Multi-Persona Synthesis Throughout:**

### **Exploration Phase (30%):**

**Biochemist:** "What conformational space do we need to explore? Alpha helix, beta sheet, random coil?"

**Physicist:** "What energy function captures hydrophobic effect, electrostatics, steric clashes?"

**Mathematician:** "What parameterization avoids singularities? Quaternions vs Euler angles?"

**Ethicist:** "Are we making assumptions that favor certain protein types (bias toward European PDB submissions)?"

**Synthesis:** Try 3 approaches in parallel, benchmark against diverse PDB set (bacterial, plant, mammalian proteins).

### **Optimization Phase (20%):**

**Biochemist:** "Spring constants too weak - helices unraveling. Increase to match AMBER ff14SB."

**Physicist:** "Timestep too large - energy oscillating. Reduce to 0.5 fs for stability."

**Mathematician:** "Force calculation is O(n²) bottleneck. Apply Williams Optimizer batching."

**Ethicist:** "Optimization makes code less interpretable. Add comments explaining each approximation."

**Synthesis:** Tune parameters based on biochemical literature, ensure physical stability, optimize algorithmically, preserve interpretability.

### **Stabilization Phase (50%):**

**Biochemist:** "Test on edge cases: disulfide-rich proteins, membrane proteins, intrinsically disordered."

**Physicist:** "Validate energy values match thermodynamic experiments (folding free energy)."

**Mathematician:** "Prove convergence: energy decreases monotonically, algorithm terminates in finite steps."

**Ethicist:** "Document failure modes so users know when NOT to trust predictions."

**Synthesis:** Comprehensive testing, thermodynamic validation, mathematical proof, transparent limitations.

---

## 📚 METHODOLOGY SUMMARY

**FoldVedic Wave Development = Scientific Rigor × Rapid Iteration × D3 Quality**

1. **Waves of 3 parallel agents** (divide and conquer)
2. **30/20/50 regime** (explore → optimize → stabilize)
3. **PDB validation every wave** (ground truth comparison)
4. **D3-Enterprise Grade+** (zero TODOs, 100% = 100%)
5. **Five Timbres quality** (≥0.90 target)
6. **Ananta Reasoning** (4 personas synthesized)
7. **Full documentation** (math, citations, reproducibility)
8. **Cascade to finish** (complete solution, not MVP)

**Expected Outcome:**
- 6 waves × 2 days each = 12 days to production
- Quality score ≥ 0.90 (publishable)
- RMSD <3 Å on test set (competitive with 2020 AlphaFold2)
- Speed <10s for medium proteins (real-time interactive)
- Fully open-source, interpretable, browser-based

**Let's democratize protein science. Begin Wave 1.**

---

**END OF METHODOLOGY DOCUMENT**
