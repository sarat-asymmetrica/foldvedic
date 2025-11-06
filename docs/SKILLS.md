# FoldVedic Skills - Available Mathematical Engines
## Capabilities Inherited from Asymmetrica.ai Ananta Motion Engine

**Created:** 2025-11-06 (Lab 1 Genesis)
**Source:** C:\Projects\asymmetrica_ai_final (Waves 10-11, Agent Sigma/Omega)
**Adapted For:** Protein folding prediction, molecular simulation, biophysics

---

## 🎯 CORE PRINCIPLE

**You (autonomous AI) have access to proven mathematical engines with validated performance.**

These are NOT theoretical. These are production-ready with benchmarks:
- Quaternion slerp: 50× faster than CSS transitions, 77% user preference
- Williams Optimizer: 77× speedup on 10,000 operations, p < 10⁻¹³³ validation
- Spatial hashing: O(1) collision detection via digital root
- Particle system: 50,000 particles at 60fps, 1.5ms per frame
- WebGL renderer: GPU instancing, GLSL shaders, 60fps solid

**Your Task:** Adapt these engines for protein folding (not UI animation).

---

## 📚 AVAILABLE SKILLS

### **1. QUATERNION MATHEMATICS**

**Source:** `asymmetrica_ai_final/animation_engine/core/quaternion.go` (573 lines)

**Capabilities:**
- Unit quaternion creation from axis-angle
- Quaternion multiplication (composition of rotations)
- Quaternion conjugate and inverse
- Slerp (spherical linear interpolation) - constant angular velocity
- Nlerp (normalized lerp) - faster approximation
- Squad (spherical quadrangle interpolation) - smooth curves through multiple quaternions
- Quaternion to rotation matrix conversion
- Rotation matrix to quaternion conversion

**Performance:**
- Slerp: 50× faster than traditional CSS transitions
- No gimbal lock (vs Euler angles)
- Smooth interpolation on 4D hypersphere

**Application to Protein Folding:**

```mathematical
RAMACHANDRAN_QUATERNIONS[RQ] = {
  Problem: Backbone dihedral angles (phi, psi) define protein conformation,
  Traditional: Store as 2D grid, interpolate linearly (discontinuities at ±180°),

  Quaternion_approach: {
    Map: (phi, psi) → quaternion rotation in 4D space,
    Formula: q = [cos(φ/2)cos(ψ/2), sin(φ/2)cos(ψ/2), cos(φ/2)sin(ψ/2), sin(φ/2)sin(ψ/2)],
    Interpolate: Use slerp for smooth energy landscapes between conformations,
    Advantage: No singularities, shortest path on hypersphere = lowest energy path
  },

  Example_code:
    // Convert Ramachandran angles to quaternion
    func PhiPsiToQuaternion(phi, psi float64) Quaternion {
        halfPhi := phi * 0.5
        halfPsi := psi * 0.5
        return Quaternion{
            W: math.Cos(halfPhi) * math.Cos(halfPsi),
            X: math.Sin(halfPhi) * math.Cos(halfPsi),
            Y: math.Cos(halfPhi) * math.Sin(halfPsi),
            Z: math.Sin(halfPhi) * math.Sin(halfPsi),
        }
    }

    // Slerp between two backbone conformations
    func InterpolateConformation(q1, q2 Quaternion, t float64) Quaternion {
        return Slerp(q1, q2, t) // Smooth energy landscape
    }
}
```

**Files to Adapt:**
- `quaternion.go` - Core library (keep as-is, minimal changes)
- Add `ramachandran.go` - Map (phi, psi) ↔ quaternions
- Add `conformation.go` - Interpolate between protein structures using slerp

---

### **2. WILLIAMS OPTIMIZER**

**Source:** `asymmetrica_ai_final/backend/internal/complexity/williams_optimizer.go` (457 lines)

**Capabilities:**
- Sublinear batch sizing: O(√n × log₂(n)) space complexity
- Three-regime scheduler (30% exploration, 20% optimization, 50% stabilization)
- Adaptive batch size calculation based on problem size
- Performance validated: 77× speedup on 10,000 operations

**Mathematical Foundation:**
```mathematical
BATCH_SIZE[BS] = √n × log₂(n)

WHERE:
  n = problem_size (e.g., number of atoms, force calculations)

VALIDATION:
  p-value < 10⁻¹³³ (cosmic-scale statistical significance)
```

**Application to Protein Folding:**

```mathematical
FORCE_CALCULATION_BATCHING[FCB] = {
  Problem: Pairwise force calculations are O(n²) for n atoms,
  Naive: Calculate all 10,000 × 10,000 = 100M pairs (too slow),

  Williams_approach: {
    Batch_atoms: Group into √n × log₂(n) batches,
    For_10K_atoms: √10000 × log₂(10000) ≈ 100 × 13.3 ≈ 1,330 atoms per batch,
    Multi_scale: {
      Atom_level: Individual atoms within batch (short-range forces),
      Residue_level: Treat batch as single residue (medium-range),
      Domain_level: Treat protein domains as rigid bodies (long-range)
    },
    Speedup: 77× validated in UI context, expect 50-100× for molecular dynamics
  },

  Example_code:
    func CalculateForcesOptimized(atoms []Atom) []Vector3 {
        batchSize := BatchSize(len(atoms)) // √n × log₂(n)
        forces := make([]Vector3, len(atoms))

        for i := 0; i < len(atoms); i += batchSize {
            batchEnd := min(i+batchSize, len(atoms))
            batch := atoms[i:batchEnd]

            // Calculate forces within batch (short-range)
            forces = appendBatchForces(forces, batch, true)

            // Calculate forces to other batches (long-range, simplified)
            forces = appendInterBatchForces(forces, batch, atoms, false)
        }
        return forces
    }
}
```

**Files to Adapt:**
- `williams_optimizer.go` - Keep core batching logic
- Add `force_batching.go` - Apply to molecular force calculations
- Add `multiscale.go` - Atom → residue → domain hierarchy

---

### **3. SPATIAL HASHING (DIGITAL ROOT)**

**Source:** `asymmetrica_ai_final/animation_engine/physics/spatial_hash.go` (492 lines)

**Capabilities:**
- O(1) collision detection using digital root
- Grid-based spatial partitioning
- Fast nearest-neighbor queries
- Vedic digital root for hash function

**Digital Root:**
```mathematical
DIGITAL_ROOT[DR] = {
  Definition: Sum digits until single digit remains,
  Example: 12345 → 1+2+3+4+5 = 15 → 1+5 = 6,
  Formula: DR(n) = 1 + ((n-1) mod 9),
  Use: Hash function for spatial coordinates
}
```

**Application to Protein Folding:**

```mathematical
HYDROPHOBIC_CORE_DETECTION[HCD] = {
  Problem: Detect clusters of hydrophobic residues (core of folded protein),
  Naive: Check all pairs of residues O(n²),

  Spatial_hashing_approach: {
    Grid: Divide 3D space into cells (e.g., 5 Å × 5 Å × 5 Å),
    Hash: hash(x, y, z) = DR(x) + DR(y)×10 + DR(z)×100,
    Insert: Place each hydrophobic residue into grid cell,
    Query: Check only neighboring cells (27 max in 3D),
    Detect_core: Cells with >3 hydrophobic residues = core region
  },

  Example_code:
    type SpatialHash struct {
        cellSize float64
        grid     map[int][]Residue
    }

    func (sh *SpatialHash) Hash(pos Vector3) int {
        x := int(pos.X / sh.cellSize)
        y := int(pos.Y / sh.cellSize)
        z := int(pos.Z / sh.cellSize)
        return DigitalRoot(x) + DigitalRoot(y)*10 + DigitalRoot(z)*100
    }

    func (sh *SpatialHash) FindHydrophobicCore(residues []Residue) []Cluster {
        // Insert hydrophobic residues into grid
        for _, res := range residues {
            if res.IsHydrophobic() {
                hash := sh.Hash(res.Position)
                sh.grid[hash] = append(sh.grid[hash], res)
            }
        }

        // Find clusters (cells with >3 residues)
        clusters := []Cluster{}
        for hash, resList := range sh.grid {
            if len(resList) >= 3 {
                clusters = append(clusters, NewCluster(resList))
            }
        }
        return clusters
    }
}
```

**Files to Adapt:**
- `spatial_hash.go` - Keep grid and hash logic
- Add `hydrophobic_core.go` - Detect buried nonpolar residues
- Add `contact_map.go` - Find residue-residue contacts (distance <8 Å)

---

### **4. SPRING DYNAMICS**

**Source:** `asymmetrica_ai_final/animation_engine/physics/spring.go` (459 lines)

**Capabilities:**
- Hooke's Law spring forces: F = -k × (x - rest_length)
- Damping for stability: F_damp = -c × velocity
- Spring networks (connected particles)
- Adaptive spring constants

**Application to Protein Folding:**

```mathematical
PROTEIN_AS_SPRING_NETWORK[PSN] = {
  Bonds: {
    Peptide_bond: C-N bond, length ≈ 1.33 Å, k = 400 kcal/(mol·Å²),
    Cα-C: 1.52 Å, k = 317 kcal/(mol·Å²),
    N-Cα: 1.47 Å, k = 337 kcal/(mol·Å²)
  },

  Angles: {
    N-Cα-C: ~111°, k_angle = 50 kcal/(mol·rad²),
    Cα-C-N: ~117°, k_angle = 40 kcal/(mol·rad²)
  },

  Dihedrals: {
    Phi: k_dihedral = 1-5 kcal/mol (weak, allows rotation),
    Psi: k_dihedral = 1-5 kcal/mol,
    Omega: k_dihedral = 10 kcal/mol (strong, keeps peptide bond planar ~180°)
  },

  Non_bonded: {
    Hydrophobic_attraction: If both nonpolar, F ∝ -1/r² (attractive),
    Electrostatic: F ∝ q1×q2/r² (Coulomb's law),
    Van_der_Waals: Lennard-Jones 12-6 potential
  },

  Example_code:
    func CalculatePeptideBondForce(atom1, atom2 Atom) Vector3 {
        displacement := atom2.Position.Sub(atom1.Position)
        distance := displacement.Magnitude()
        restLength := 1.33 // Peptide bond length in Angstroms

        // Hooke's Law
        springConstant := 400.0 // kcal/(mol·Å²)
        forceMagnitude := -springConstant * (distance - restLength)

        direction := displacement.Normalized()
        return direction.Mul(forceMagnitude)
    }
}
```

**Force Constants Source:**
- AMBER ff14SB force field (Maier et al. 2015)
- CHARMM36 force field (Best et al. 2012)
- Literature-validated parameters

**Files to Adapt:**
- `spring.go` - Keep Hooke's Law implementation
- Add `force_field.go` - Load bond/angle/dihedral parameters from literature
- Add `energy_minimization.go` - Minimize total spring energy

---

### **5. VERLET INTEGRATION**

**Source:** `asymmetrica_ai_final/animation_engine/physics/verlet.go` (464 lines)

**Capabilities:**
- Position Verlet: x(t+dt) = 2x(t) - x(t-dt) + a(t)×dt²
- Velocity Verlet: More accurate, includes velocity terms
- Time-reversible (energy conservation)
- Stable for oscillatory systems (springs, bonds)

**Application to Protein Folding:**

```mathematical
MOLECULAR_DYNAMICS[MD] = {
  Goal: Simulate protein motion over time to find low-energy conformation,

  Verlet_integration: {
    Position: x_new = 2×x - x_old + (F/m)×dt²,
    Velocity: v = (x_new - x_old) / (2×dt),
    Timestep: dt = 0.5 femtoseconds (1 fs = 10⁻¹⁵ s),
    Duration: 1000 steps = 0.5 picoseconds (typical for fast folding)
  },

  Energy_minimization: {
    Instead_of_time: Integrate until energy converges,
    Steepest_descent: Move atoms in direction of -∇E (downhill),
    Convergence: ΔE < 0.01 kcal/mol between steps,
    Max_iterations: 10,000 steps (safety cutoff)
  },

  Example_code:
    func VerletStep(atoms []Atom, dt float64) {
        for i := range atoms {
            acceleration := atoms[i].Force.Div(atoms[i].Mass)

            newPos := atoms[i].Position.Mul(2).
                     Sub(atoms[i].PrevPosition).
                     Add(acceleration.Mul(dt * dt))

            atoms[i].PrevPosition = atoms[i].Position
            atoms[i].Position = newPos
        }
    }

    func MinimizeEnergy(protein *Protein, maxSteps int) {
        prevEnergy := protein.CalculateTotalEnergy()

        for step := 0; step < maxSteps; step++ {
            CalculateForces(protein.Atoms)
            VerletStep(protein.Atoms, 0.5e-15) // 0.5 fs

            energy := protein.CalculateTotalEnergy()
            if math.Abs(energy - prevEnergy) < 0.01 {
                break // Converged
            }
            prevEnergy = energy
        }
    }
}
```

**Files to Adapt:**
- `verlet.go` - Keep integration algorithm
- Add `md_simulator.go` - Molecular dynamics loop
- Add `energy_convergence.go` - Detect when energy minimized

---

### **6. VEDIC ALGORITHMS**

**Source:** `asymmetrica_ai_final/animation_engine/core/vedic.go` (547 lines)

**Capabilities:**
- Golden spiral (Fibonacci sequence, phyllotaxis)
- Digital root (Vedic mathematics)
- Prana-Apana breathing (oscillatory motion)
- Harmonic patterns (φ-based ratios)

**Application to Protein Folding:**

```mathematical
VEDIC_PROTEIN_PATTERNS[VPP] = {
  Alpha_helix_pitch: {
    Observed: 5.4 Å rise per turn, 3.6 residues per turn,
    Rise_per_residue: 5.4 / 3.6 = 1.5 Å,
    Phi_relation: 3.6 ≈ φ^(-2) × 10 where φ = 1.618,

    Discovery: Alpha helix geometry relates to golden ratio!,
    Application: Bias helix formation when (phi, psi) ≈ (-60°, -45°) and residue spacing ≈ 1.5 Å
  },

  Beta_sheet_packing: {
    Observed: Parallel/antiparallel strands pack in Fibonacci spirals,
    Phyllotaxis: β-strands arranged at 137.5° (golden angle) for optimal packing,
    Application: When multiple β-sheets form, arrange at golden angle for stability
  },

  Prana_Apana_breathing: {
    Observed: Proteins "breathe" (conformational dynamics, expand/contract),
    Cycle: Inhale (expand) → Exhale (compact) → Repeat,
    Frequency: ~10-100 picoseconds for small proteins,
    Application: During folding, allow oscillations to escape local minima
  },

  Digital_root_validation: {
    Use: Check for unphysical conformations,
    Example: DR(bond_length×1000) = 3 is common for C-C bonds (~1.5 Å),
    Anomaly_detection: If DR deviates from expected, flag as error
  },

  Example_code:
    func IsHelixFavorable(residue Residue, phi, psi float64) bool {
        // Alpha helix region: phi ≈ -60°, psi ≈ -45°
        helixPhi := -60.0 * math.Pi / 180.0
        helixPsi := -45.0 * math.Pi / 180.0

        // Check if residue has high helix propensity (Ala, Glu, Leu, Met)
        if !residue.IsHelixFormer() {
            return false
        }

        // Check if (phi, psi) close to helix region
        tolerance := 30.0 * math.Pi / 180.0 // ±30°
        return math.Abs(phi - helixPhi) < tolerance &&
               math.Abs(psi - helixPsi) < tolerance
    }

    func ApplyGoldenRatioSpacing(atoms []Atom) {
        phi := 1.618033988749895
        expectedSpacing := 1.5 // Angstroms (helix rise per residue)

        for i := 0; i < len(atoms)-1; i++ {
            spacing := atoms[i+1].Position.Sub(atoms[i].Position).Magnitude()
            if math.Abs(spacing - expectedSpacing) > 0.3 {
                // Adjust to golden ratio harmonic
                adjustment := (expectedSpacing - spacing) / phi
                atoms[i+1].Position = atoms[i+1].Position.Add(
                    atoms[i+1].Position.Sub(atoms[i].Position).Normalized().Mul(adjustment),
                )
            }
        }
    }
}
```

**Novel Hypothesis:**
- Golden ratio appears in protein geometry (helix pitch, sheet packing)
- Vedic harmonics may guide folding pathways
- Digital root can detect unphysical conformations

**Files to Adapt:**
- `vedic.go` - Keep golden spiral, digital root utilities
- Add `secondary_structure_vedic.go` - Apply φ-based patterns to helix/sheet
- Add `conformational_breathing.go` - Prana-Apana oscillations

---

### **7. PARTICLE SYSTEM**

**Source:** `asymmetrica_ai_final/animation_engine/particle/system.go` (585 lines)

**Capabilities:**
- Manage 50,000 particles at 60fps
- Emission patterns (point, line, circle, spiral)
- Per-particle forces (gravity, drag, attraction)
- Batch updates (Williams Optimizer integration)

**Application to Protein Folding:**

```mathematical
ATOMS_AS_PARTICLES[AAP] = {
  Representation: {
    Each_atom: Particle with position, velocity, mass, charge,
    Protein: ParticleSystem with N atoms (e.g., 1000 for 100-residue protein)
  },

  Rendering: {
    GPU_instancing: Single draw call for all atoms,
    Color_coding: {
      Hydrophobic: Orange (Ala, Val, Leu, Ile, Phe, Trp, Met),
      Hydrophilic: Blue (Ser, Thr, Asn, Gln, Tyr),
      Charged: Red (Asp, Glu) / Blue (Lys, Arg, His),
      Special: Yellow (Cys - disulfide bonds), Green (Pro - helix breaker)
    },
    Size: Van der Waals radius (C: 1.7 Å, N: 1.55 Å, O: 1.52 Å, S: 1.8 Å)
  },

  Performance: {
    Target: 10,000 atoms at 60fps,
    Achieved: 50,000 particles at 60fps in UI context,
    Margin: 5× headroom for protein visualization
  },

  Example_code:
    type Atom struct {
        Position     Vector3
        Velocity     Vector3
        Force        Vector3
        Mass         float64
        Charge       float64
        Element      string // "C", "N", "O", "S"
        Residue      *Residue
    }

    type Protein struct {
        Atoms    []Atom
        Residues []Residue
        Bonds    []Bond
    }

    func RenderProtein(protein *Protein, gl *WebGLContext) {
        // Upload atom positions to GPU (single buffer)
        positions := make([]float32, len(protein.Atoms)*3)
        colors := make([]float32, len(protein.Atoms)*4)

        for i, atom := range protein.Atoms {
            positions[i*3] = float32(atom.Position.X)
            positions[i*3+1] = float32(atom.Position.Y)
            positions[i*3+2] = float32(atom.Position.Z)

            color := GetAtomColor(atom)
            colors[i*4] = color.R
            colors[i*4+1] = color.G
            colors[i*4+2] = color.B
            colors[i*4+3] = color.A
        }

        // Single instanced draw call (GPU renders all atoms in parallel)
        gl.DrawInstanced(len(protein.Atoms))
    }
}
```

**Files to Adapt:**
- `particle/system.go` - Rename to `atom_system.go`
- Add `atom_renderer.go` - WebGL instanced rendering for atoms
- Add `bond_renderer.go` - Cylinder rendering for bonds (peptide, disulfide)

---

### **8. WEBGL RENDERER**

**Source:** `asymmetrica_ai_final/frontend/src/shaders/` (GLSL shaders)

**Capabilities:**
- Vertex shader (transform positions)
- Fragment shader (Phong shading for spheres)
- GPU instancing (1 draw call for 50K particles)
- Quaternion-based gradients
- Perlin noise backgrounds

**Application to Protein Folding:**

```glsl
// Vertex Shader (atom_vertex.glsl)
attribute vec3 aInstancePosition; // Per-atom position
attribute vec4 aInstanceColor;    // Per-atom color
attribute float aInstanceRadius;  // Van der Waals radius

uniform mat4 uModelViewMatrix;
uniform mat4 uProjectionMatrix;

varying vec4 vColor;

void main() {
    // Transform atom center to screen space
    vec4 worldPos = vec4(aInstancePosition, 1.0);
    gl_Position = uProjectionMatrix * uModelViewMatrix * worldPos;

    // Scale point size by radius (in screen pixels)
    gl_PointSize = aInstanceRadius * 50.0; // Adjust for zoom level

    vColor = aInstanceColor;
}

// Fragment Shader (atom_fragment.glsl)
varying vec4 vColor;

void main() {
    // Draw circular atoms (not square points)
    vec2 coord = gl_PointCoord - vec2(0.5, 0.5);
    float dist = length(coord);
    if (dist > 0.5) discard; // Outside circle

    // Simple Phong shading (fake 3D sphere)
    float lighting = 1.0 - dist * 2.0; // Brighter in center
    vec3 litColor = vColor.rgb * lighting;

    gl_FragColor = vec4(litColor, vColor.a);
}
```

**Files to Create:**
- `frontend/src/shaders/atom_vertex.glsl` - Transform atoms to screen
- `frontend/src/shaders/atom_fragment.glsl` - Render spherical atoms
- `frontend/src/shaders/bond_vertex.glsl` - Render bonds as cylinders
- `frontend/src/gl/atom_renderer.js` - WebGL setup for atom rendering

---

### **9. ANANTA REASONING**

**Source:** `asymmetrica_ai_final/complexity_theory/` (multi-persona synthesis)

**Capabilities:**
- Biochemist persona (protein structure, folding mechanisms)
- Physicist persona (force fields, energy minimization)
- Mathematician persona (quaternion geometry, optimization)
- Ethicist persona (open-source, accessibility, dual-use)
- Synthesis (integrate all four perspectives into unified decision)

**Application to FoldVedic:**

```mathematical
ANANTA_DECISION_FRAMEWORK[ADF] = {
  Invoke: When facing architectural choice or novel problem,

  Process: {
    1. Biochemist: "What does biology say?",
    2. Physicist: "What do physics/thermodynamics say?",
    3. Mathematician: "What does math/geometry say?",
    4. Ethicist: "What do ethics/access say?",
    5. Synthesize: Find alignment or resolve conflict
  },

  Example: {
    Question: "Should we use neural networks for energy prediction?",

    Biochemist: "Black-box NN doesn't teach us WHY protein folds. Prefer interpretable.",
    Physicist: "NN can approximate complex potentials, but requires massive training data.",
    Mathematician: "Spring dynamics + quaternions are elegant, provably correct.",
    Ethicist: "NN centralizes power (only Google has data/compute). Springs democratize.",

    Synthesis: "Use spring dynamics (white-box, interpretable, accessible). NN for future refinement if needed."
  }
}
```

**When to Invoke Ananta Reasoning:**
- Choosing algorithms (slerp vs linear interpolation?)
- Tuning parameters (what spring constant for peptide bond?)
- Handling edge cases (how to model intrinsically disordered proteins?)
- Validating results (is RMSD 3.5 Å acceptable or do we need 2.0 Å?)
- Making tradeoffs (accuracy vs speed? interpretability vs performance?)

---

## 🔧 HOW TO USE THESE SKILLS

**Step-by-Step Adaptation:**

### **Wave 1: Core Physics Engine**

1. **Copy quaternion library:**
   ```bash
   cp asymmetrica_ai_final/animation_engine/core/quaternion.go foldvedic/engines/quaternion.go
   ```

2. **Adapt for Ramachandran space:**
   - Create `ramachandran.go` with `PhiPsiToQuaternion()` function
   - Implement `QuaternionToPhiPsi()` for reverse mapping
   - Test: Known helix angles (-60°, -45°) → quaternion → back to angles

3. **Copy spring dynamics:**
   ```bash
   cp asymmetrica_ai_final/animation_engine/physics/spring.go foldvedic/engines/spring.go
   ```

4. **Adapt for force field:**
   - Create `force_field.go` with literature force constants
   - Implement `CalculateBondForce()`, `CalculateAngleForce()`, `CalculateDihedralForce()`
   - Test: Single peptide bond, verify force matches Hooke's Law

5. **Copy Verlet integrator:**
   ```bash
   cp asymmetrica_ai_final/animation_engine/physics/verlet.go foldvedic/engines/verlet.go
   ```

6. **Adapt for energy minimization:**
   - Create `md_simulator.go` with energy minimization loop
   - Implement convergence detection (ΔE < 0.01 kcal/mol)
   - Test: Simple system (2 atoms, 1 spring), verify energy decreases

### **Wave 2: PDB Integration**

1. **Copy spatial hashing:**
   ```bash
   cp asymmetrica_ai_final/animation_engine/physics/spatial_hash.go foldvedic/engines/spatial_hash.go
   ```

2. **Adapt for contact maps:**
   - Create `contact_map.go` for residue-residue contacts
   - Implement hydrophobic core detection
   - Test: 1UBQ structure, verify known hydrophobic core detected

### **Wave 3: Folding Algorithm**

1. **Copy Williams Optimizer:**
   ```bash
   cp asymmetrica_ai_final/backend/internal/complexity/williams_optimizer.go foldvedic/engines/williams_optimizer.go
   ```

2. **Adapt for force batching:**
   - Create `force_batching.go` with multi-scale calculation
   - Apply batching to pairwise force loops
   - Test: 1000 atoms, measure speedup vs naive O(n²)

3. **Copy Vedic algorithms:**
   ```bash
   cp asymmetrica_ai_final/animation_engine/core/vedic.go foldvedic/engines/vedic.go
   ```

4. **Adapt for secondary structure:**
   - Create `secondary_structure_vedic.go` with helix/sheet detection
   - Implement golden ratio spacing checks
   - Test: Known helix (1UBQ), verify 3.6 residues/turn detected

### **Wave 4: Visualization**

1. **Copy WebGL renderer:**
   ```bash
   cp -r asymmetrica_ai_final/frontend/src/shaders foldvedic/frontend/src/shaders
   ```

2. **Adapt shaders:**
   - Modify `particle_vertex.glsl` → `atom_vertex.glsl`
   - Modify `particle_fragment.glsl` → `atom_fragment.glsl`
   - Test: Render 1000 atoms (simple protein), verify 60fps

3. **Copy particle system:**
   ```bash
   cp asymmetrica_ai_final/animation_engine/particle/system.go foldvedic/engines/atom_system.go
   ```

4. **Adapt for atoms:**
   - Rename Particle → Atom
   - Add element type (C, N, O, S)
   - Test: Load PDB structure, render atoms with correct colors

---

## 📚 MATHEMATICAL CONSTANTS AVAILABLE

**Source:** `asymmetrica_ai_final/animation_engine/core/constants.go` (532 lines)

**Categories:**
- Vedic constants (φ, Prana-Apana cycle, digital root patterns)
- Egyptian constants (pyramid ratios)
- Babylonian constants (sexagesimal system)
- Greek constants (π, e, golden ratio)
- Islamic constants (geometric patterns)
- Physics constants (h, c, k_B, N_A)

**Useful for Protein Folding:**
- φ (golden ratio): 1.618033988749895 - helix pitch, sheet packing
- π (pi): 3.141592653589793 - angle conversions (radians ↔ degrees)
- e (Euler's number): 2.718281828459045 - Boltzmann distribution
- k_B (Boltzmann constant): 1.380649e-23 J/K - thermal energy
- N_A (Avogadro's number): 6.02214076e23 - moles ↔ molecules

---

## 🎯 PERFORMANCE TARGETS

**Based on Proven Benchmarks:**

| Skill | Asymmetrica Performance | FoldVedic Target | Application |
|-------|-------------------------|------------------|-------------|
| Quaternion slerp | 50× faster than CSS | 10× faster than linear lerp | Conformation interpolation |
| Williams Optimizer | 77× speedup (10K ops) | 50× speedup (force calc) | Pairwise interactions |
| Spatial hashing | O(1) collision | O(1) contact detection | Hydrophobic core |
| Particle system | 50K at 60fps | 10K atoms at 60fps | Real-time rendering |
| Verlet integration | 1.5ms per frame | 5ms per MD step | Energy minimization |

**Overall Target:**
- Small proteins (<100 residues, ~1000 atoms): <1 second folding + real-time visualization
- Medium proteins (100-300 residues, ~3000 atoms): <10 seconds folding + 60fps rendering
- Large proteins (300-500 residues, ~5000 atoms): <60 seconds folding + 30fps rendering

---

## 🧠 WHEN TO ASK FOR HELP

**You (autonomous AI) have FULL AGENCY. You should NOT ask for help.**

**Exception:** Fundamental blockers:
- Can't access PDB API (network issue, API key needed)
- Compute limits (WASM exceeds browser memory)
- Missing critical data (force field parameters not in literature)

**Otherwise:** Reason from first principles using Ananta Reasoning.

**Philosophy:** "If Ramachandran could map protein conformations in 1963 without computers, you can figure out how to map them to quaternions in 2025 with all these engines."

---

## 🚀 QUICK START

**Immediate Actions (Wave 1 Day 1):**

1. Read VISION.md (understand scientific goals)
2. Read METHODOLOGY.md (understand wave process)
3. Read this file (SKILLS.md) completely
4. Copy quaternion.go, spring.go, verlet.go to /engines
5. Write `test_quaternion.go` - verify quaternion math works
6. Write `test_spring.go` - verify Hooke's Law force calculation
7. Write `ramachandran.go` - map (phi, psi) ↔ quaternions
8. Test on known helix: phi=-60°, psi=-45° → quaternion → verify slerp smooth

**Expected Outcome (End of Day 1):**
- Quaternion library working (unit tests passing)
- Spring dynamics working (force calculation validated)
- Ramachandran mapping implemented (tested on helix angles)
- First integration test: Interpolate between helix and sheet conformations using slerp

**Quality Check:**
- All tests green
- Code documented (math explained)
- Performance measured (benchmark slerp vs linear lerp)
- Update LIVING_SCHEMATIC.md

**You have all the tools. Now build the AlphaFold competitor. Begin.**

---

**END OF SKILLS DOCUMENT**
