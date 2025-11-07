# Wave 4 Compilation Fixes - IN PROGRESS

**Status:** Fixing remaining compilation errors

## Completed Fixes:

1. ✅ Created `folding/helpers.go` with public wrappers:
   - `NewProteinFromSequence()` - creates extended chain
   - `CloneProtein()` - deep copy
   - `CalculateEnergy()` - wrapper for physics.CalculateTotalEnergy
   - `GetSequence()` - extract sequence from residues

2. ✅ Created `parser/protein_methods.go` with:
   - `Protein.Copy()` method
   - `Protein.Sequence()` method

3. ✅ Created `sampling/wrappers.go` with simplified interfaces:
   - `FibonacciSphereBasins()`
   - `VedicMonteCarlo()`
   - `GenerateFragmentStructures()` (renamed to avoid conflict)
   - `BasinExplorer()`

4. ✅ Fixed full_pipeline/main.go:
   - Changed `nativeProtein.Sequence` to `nativeProtein.Sequence()`
   - Changed `*folding.Protein` to `*parser.Protein`

5. ✅ Fixed phase2_to_3/main.go:
   - Changed `*folding.Protein` to `*parser.Protein`

6. ✅ Fixed lbfgs_benchmark/main.go:
   - Changed `nativeProtein.Sequence` to `nativeProtein.Sequence()`

7. ✅ Fixed lbfgs_tuning.go:
   - Added error handling for `validation.CalculateRMSD()` (returns 2 values)
   - Changed `QuaternionLBFGS()` to `MinimizeLBFGS()`
   - Removed duplicate LBFGSConfig type alias

8. ✅ Fixed diversity.go:
   - Changed to use CA distance instead of Phi/Psi fields
   - Added folding import

9. ✅ Fixed ensemble.go:
   - Changed QuaternionDiversity to use CA distance

## Remaining Issues to Fix:

### 1. LBFGSConfig field mismatches
Need to map LBFGSTuningConfig fields to actual LBFGSConfig fields:
- `config.GradientTol` → `config.GradientTolerance`
- `config.StepSize` → NO DIRECT EQUIVALENT (remove or find alternative)

### 2. BasinExplorerConfig field mismatches in wrappers.go
- `config.NumSamplesPerBasin` → `config.SamplesPerBasin`
- `config.Basins` → DOESN'T EXIST (need to check how basins are configured)

### 3. CalculateRMSD error handling in ensemble.go
Multiple places need to handle the (float64, error) return:
- Line 118, 225, 230, 338

### 4. FragmentAssembly call mismatch in ensemble.go
Line 77 calls `FragmentAssembly(sequence, numFrag)` but it needs:
`FragmentAssembly(sequence, *FragmentLibrary, FragmentAssemblyConfig)`

## Next Steps:

1. Check if LBFGSConfig supports step size configuration
2. Rewrite wrappers.go to use correct BasinExplorerConfig fields
3. Add error handling to all CalculateRMSD calls in ensemble.go
4. Fix FragmentAssembly call in ensemble.go to use GenerateFragmentEnsemble or proper library

## Critical Decision Needed:

The Wave 4 code assumes interfaces that don't match the actual implementation.
Two options:
A. Continue fixing mismatches (may take 30+ more minutes)
B. Create minimal stub wrappers that compile but may not work perfectly

**Recommendation:** Continue with fixes since we're 80% done.
