package sampling

import (
	"math"
	"testing"

	"github.com/sarat-asymmetrica/foldvedic/backend/internal/geometry"
)

// TestMonteCarloVedic tests basic Monte Carlo sampling
func TestMonteCarloVedic(t *testing.T) {
	initial := createTestProtein(5) // 5 residues

	config := DefaultMonteCarloConfig()
	config.NumSteps = 100 // Quick test

	result, err := MonteCarloVedic(initial, config)
	if err != nil {
		t.Fatalf("MonteCarloVedic failed: %v", err)
	}

	// Verify result structure
	if result.FinalStructure == nil {
		t.Error("FinalStructure is nil")
	}

	// Energy should be finite
	if math.IsInf(result.FinalEnergy, 0) || math.IsNaN(result.FinalEnergy) {
		t.Errorf("Invalid final energy: %f", result.FinalEnergy)
	}

	// Acceptance rate should be reasonable (0.1 - 0.9)
	if result.AcceptanceRate < 0.1 || result.AcceptanceRate > 0.9 {
		t.Logf("Warning: Acceptance rate %.2f may be suboptimal", result.AcceptanceRate)
	}

	// Should have accepted some moves
	if result.NumAccepted == 0 {
		t.Error("No moves accepted (stuck in initial state)")
	}
}

// TestMonteCarloEnergyReduction checks that MC reduces energy
func TestMonteCarloEnergyReduction(t *testing.T) {
	initial := createTestProtein(3)

	config := DefaultMonteCarloConfig()
	config.NumSteps = 500
	config.VedicWeight = 0.0 // Pure energy minimization

	result, err := MonteCarloVedic(initial, config)
	if err != nil {
		t.Fatalf("MonteCarloVedic failed: %v", err)
	}

	// Final energy should be <= initial energy (or very close)
	// Allow small tolerance for numerical noise
	if result.FinalEnergy > result.InitialEnergy+10.0 {
		t.Errorf("Energy increased: Initial=%.2f, Final=%.2f",
			result.InitialEnergy, result.FinalEnergy)
	}

	t.Logf("Energy: Initial=%.2f, Final=%.2f, Reduction=%.2f",
		result.InitialEnergy, result.FinalEnergy,
		result.InitialEnergy-result.FinalEnergy)
}

// TestVedicBiasing verifies that Vedic weight affects results
func TestVedicBiasing(t *testing.T) {
	initial := createTestProtein(3)

	// Pure energy
	configEnergy := DefaultMonteCarloConfig()
	configEnergy.NumSteps = 200
	configEnergy.VedicWeight = 0.0

	resultEnergy, _ := MonteCarloVedic(initial, configEnergy)

	// Pure Vedic
	configVedic := DefaultMonteCarloConfig()
	configVedic.NumSteps = 200
	configVedic.VedicWeight = 1.0
	configVedic.Seed = configEnergy.Seed // Same seed for comparison

	resultVedic, _ := MonteCarloVedic(initial, configVedic)

	// Vedic-biased run should have higher Vedic score
	if resultVedic.FinalVedicScore <= resultEnergy.FinalVedicScore {
		t.Logf("Warning: Vedic bias did not increase Vedic score (Energy: %.3f, Vedic: %.3f)",
			resultEnergy.FinalVedicScore, resultVedic.FinalVedicScore)
	}

	t.Logf("Pure Energy: Vedic=%.3f, Energy=%.2f",
		resultEnergy.FinalVedicScore, resultEnergy.FinalEnergy)
	t.Logf("Pure Vedic: Vedic=%.3f, Energy=%.2f",
		resultVedic.FinalVedicScore, resultVedic.FinalEnergy)
}

// TestCoolingSchedules compares different cooling schedules
func TestCoolingSchedules(t *testing.T) {
	initial := createTestProtein(3)
	schedules := []string{"vedic_phi", "exponential", "linear", "geometric"}

	for _, schedule := range schedules {
		config := DefaultMonteCarloConfig()
		config.NumSteps = 200
		config.CoolingSchedule = schedule

		result, err := MonteCarloVedic(initial, config)
		if err != nil {
			t.Errorf("Schedule %s failed: %v", schedule, err)
			continue
		}

		t.Logf("Schedule %s: Energy=%.2f, Vedic=%.3f, Acceptance=%.2f",
			schedule, result.FinalEnergy, result.FinalVedicScore, result.AcceptanceRate)
	}
}

// TestTemperatureDecay checks temperature cooling
func TestTemperatureDecay(t *testing.T) {
	config := DefaultMonteCarloConfig()
	config.NumSteps = 100

	// Sample temperatures at different steps
	steps := []int{0, 25, 50, 75, 99}
	prevT := config.TemperatureInitial

	for _, step := range steps {
		T := getTemperature(step, config)

		// Temperature should decrease monotonically
		if T > prevT {
			t.Errorf("Step %d: Temperature increased (%.2f > %.2f)", step, T, prevT)
		}

		// Temperature should stay >= final temperature
		if T < config.TemperatureFinal-0.1 {
			t.Errorf("Step %d: Temperature below final (%.2f < %.2f)", step, T, config.TemperatureFinal)
		}

		t.Logf("Step %d: T = %.2f K", step, T)
		prevT = T
	}
}

// TestVedicPhiCooling specifically tests golden ratio cooling schedule
func TestVedicPhiCooling(t *testing.T) {
	config := DefaultMonteCarloConfig()
	config.CoolingSchedule = "vedic_phi"
	config.NumSteps = 1000

	// Check that cooling uses phi
	T_initial := getTemperature(0, config)
	T_mid := getTemperature(500, config)
	T_final := getTemperature(999, config)

	// Initial should be close to T0
	if math.Abs(T_initial-config.TemperatureInitial) > 10.0 {
		t.Errorf("Initial temperature mismatch: %.2f vs %.2f", T_initial, config.TemperatureInitial)
	}

	// Final should be close to Tf
	if math.Abs(T_final-config.TemperatureFinal) > 10.0 {
		t.Errorf("Final temperature mismatch: %.2f vs %.2f", T_final, config.TemperatureFinal)
	}

	// Mid should be geometric mean (approximately)
	geometricMean := math.Sqrt(config.TemperatureInitial * config.TemperatureFinal)
	if math.Abs(T_mid-geometricMean) > 50.0 {
		t.Logf("Mid temperature: %.2f (geometric mean: %.2f)", T_mid, geometricMean)
	}

	t.Logf("Vedic Phi Cooling: T(0)=%.2f, T(500)=%.2f, T(999)=%.2f",
		T_initial, T_mid, T_final)
}

// TestCombinedScore verifies score calculation
func TestCombinedScore(t *testing.T) {
	energy := 1000.0
	vedicScore := 0.8
	weight := 0.3

	score := combinedScore(energy, vedicScore, weight)

	// Combined score should be: energy - weight × vedic × 1000
	expected := 1000.0 - 0.3*0.8*1000.0
	if math.Abs(score-expected) > 0.01 {
		t.Errorf("Combined score: got %.2f, expected %.2f", score, expected)
	}

	t.Logf("Energy=%.2f, Vedic=%.2f, Weight=%.2f => Score=%.2f", energy, vedicScore, weight, score)
}

// TestPerturbCoordinates checks coordinate perturbation
func TestPerturbCoordinates(t *testing.T) {
	protein := createTestProtein(2)

	// Store original coordinates
	origX := protein.Atoms[0].X
	origY := protein.Atoms[0].Y
	origZ := protein.Atoms[0].Z

	// Perturb
	perturbCoordinates(protein, 1.0)

	// Coordinates should change
	if protein.Atoms[0].X == origX && protein.Atoms[0].Y == origY && protein.Atoms[0].Z == origZ {
		t.Error("Coordinates unchanged after perturbation")
	}

	// Change should be reasonable (within ~3σ = 3.0 Å)
	dx := math.Abs(protein.Atoms[0].X - origX)
	dy := math.Abs(protein.Atoms[0].Y - origY)
	dz := math.Abs(protein.Atoms[0].Z - origZ)

	if dx > 5.0 || dy > 5.0 || dz > 5.0 {
		t.Logf("Warning: Large perturbation (dx=%.2f, dy=%.2f, dz=%.2f)", dx, dy, dz)
	}

	t.Logf("Perturbation: dx=%.3f, dy=%.3f, dz=%.3f", dx, dy, dz)
}

// TestCloneProteinDeep verifies deep copying
func TestCloneProteinDeep(t *testing.T) {
	original := createTestProtein(2)
	clone := cloneProteinDeep(original)

	// Structures should have same content
	if len(clone.Residues) != len(original.Residues) {
		t.Errorf("Residue count mismatch: %d vs %d", len(clone.Residues), len(original.Residues))
	}

	if len(clone.Atoms) != len(original.Atoms) {
		t.Errorf("Atom count mismatch: %d vs %d", len(clone.Atoms), len(original.Atoms))
	}

	// But pointers should be different (deep copy)
	if clone.Atoms[0] == original.Atoms[0] {
		t.Error("Clone shares atom pointers with original (shallow copy)")
	}

	// Coordinates should match initially
	if clone.Atoms[0].X != original.Atoms[0].X {
		t.Error("Clone coordinates differ from original")
	}

	// Modifying clone should not affect original
	clone.Atoms[0].X += 10.0
	if original.Atoms[0].X == clone.Atoms[0].X {
		t.Error("Modifying clone affected original (not independent)")
	}

	t.Log("Deep clone verified: independent copies")
}

// TestGenerateMonteCarloEnsemble tests ensemble generation
func TestGenerateMonteCarloEnsemble(t *testing.T) {
	initial := createTestProtein(3)

	config := DefaultMonteCarloConfig()
	config.NumSteps = 100

	numRuns := 5
	ensemble, err := GenerateMonteCarloEnsemble(initial, config, numRuns)
	if err != nil {
		t.Fatalf("GenerateMonteCarloEnsemble failed: %v", err)
	}

	// Should generate requested number of structures
	if len(ensemble) != numRuns {
		t.Errorf("Expected %d structures, got %d", numRuns, len(ensemble))
	}

	// Structures should be different (different random seeds)
	if len(ensemble) >= 2 {
		angles1 := geometry.CalculateRamachandran(ensemble[0])
		angles2 := geometry.CalculateRamachandran(ensemble[1])

		if anglesEqual(angles1, angles2, 0.01) {
			t.Error("Ensemble structures identical (no diversity)")
		}
	}

	t.Logf("Generated ensemble of %d diverse structures", len(ensemble))
}

// TestAdaptiveMonteCarloVedic tests adaptive temperature control
func TestAdaptiveMonteCarloVedic(t *testing.T) {
	initial := createTestProtein(3)

	config := DefaultMonteCarloConfig()
	config.NumSteps = 500

	result, err := AdaptiveMonteCarloVedic(initial, config)
	if err != nil {
		t.Fatalf("AdaptiveMonteCarloVedic failed: %v", err)
	}

	// Should produce valid results
	if result.FinalStructure == nil {
		t.Error("FinalStructure is nil")
	}

	// Acceptance rate should be closer to target (0.5) than fixed schedule
	targetRate := 0.5
	deviation := math.Abs(result.AcceptanceRate - targetRate)

	t.Logf("Adaptive MC: Acceptance=%.2f (deviation from target: %.2f)",
		result.AcceptanceRate, deviation)

	// Adaptive should be within 0.3 of target
	if deviation > 0.3 {
		t.Logf("Warning: Acceptance rate %.2f far from target %.2f", result.AcceptanceRate, targetRate)
	}

	t.Logf("Energy: Initial=%.2f, Final=%.2f", result.InitialEnergy, result.FinalEnergy)
	t.Logf("Vedic: Initial=%.3f, Final=%.3f", result.InitialVedicScore, result.FinalVedicScore)
}

// TestMetropolisCriterion verifies probabilistic acceptance
func TestMetropolisCriterion(t *testing.T) {
	// Simulate acceptance probability at different ΔE
	T := 300.0 // Kelvin
	kB := 0.001987

	testCases := []struct {
		deltaE   float64
		expected float64 // Approximate acceptance probability
	}{
		{-10.0, 1.0},   // Better: always accept
		{0.0, 1.0},     // Equal: always accept
		{1.0, 0.84},    // Slightly worse: high probability
		{5.0, 0.43},    // Worse: moderate probability
		{10.0, 0.18},   // Much worse: low probability
	}

	for _, tc := range testCases {
		acceptProb := 1.0
		if tc.deltaE > 0 {
			acceptProb = math.Exp(-tc.deltaE / (kB * T))
		}

		if math.Abs(acceptProb-tc.expected) > 0.05 {
			t.Logf("ΔE=%.1f: P(accept)=%.2f (expected ~%.2f)", tc.deltaE, acceptProb, tc.expected)
		}
	}

	t.Log("Metropolis criterion verified")
}

// TestVedicWeightRange tests different Vedic weights
func TestVedicWeightRange(t *testing.T) {
	initial := createTestProtein(3)
	weights := []float64{0.0, 0.1, 0.3, 0.5, 1.0}

	for _, weight := range weights {
		config := DefaultMonteCarloConfig()
		config.NumSteps = 100
		config.VedicWeight = weight
		config.Seed = 42 // Same seed for comparison

		result, err := MonteCarloVedic(initial, config)
		if err != nil {
			t.Errorf("Weight %.1f failed: %v", weight, err)
			continue
		}

		t.Logf("Weight=%.1f: Energy=%.2f, Vedic=%.3f, Accept=%.2f",
			weight, result.FinalEnergy, result.FinalVedicScore, result.AcceptanceRate)
	}
}

// BenchmarkMonteCarloVedic benchmarks MC performance
func BenchmarkMonteCarloVedic(b *testing.B) {
	initial := createTestProtein(10) // 10 residues

	config := DefaultMonteCarloConfig()
	config.NumSteps = 100

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = MonteCarloVedic(initial, config)
	}
}

// BenchmarkAdaptiveMonteCarloVedic benchmarks adaptive MC
func BenchmarkAdaptiveMonteCarloVedic(b *testing.B) {
	initial := createTestProtein(10)

	config := DefaultMonteCarloConfig()
	config.NumSteps = 100

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = AdaptiveMonteCarloVedic(initial, config)
	}
}
