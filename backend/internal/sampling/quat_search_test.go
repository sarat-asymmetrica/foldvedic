package sampling

import (
	"math"
	"testing"

	"github.com/sarat-asymmetrica/foldvedic/backend/internal/geometry"
	"github.com/sarat-asymmetrica/foldvedic/backend/internal/parser"
)

// TestQuaternionGuidedSearch tests basic conformational sampling
func TestQuaternionGuidedSearch(t *testing.T) {
	// Create simple test protein (3 residues)
	initial := createTestProtein(3)

	config := QuaternionSearchConfig{
		NumSamples:         5,
		SlerpSteps:         3,
		PerturbRadius:      0.5,
		Seed:               42,
		UseFibonacciSphere: true,
	}

	ensemble, err := QuaternionGuidedSearch(initial, config)
	if err != nil {
		t.Fatalf("QuaternionGuidedSearch failed: %v", err)
	}

	// Expect 5 samples × 3 steps = 15 structures
	expectedCount := config.NumSamples * config.SlerpSteps
	if len(ensemble) != expectedCount {
		t.Errorf("Expected %d structures, got %d", expectedCount, len(ensemble))
	}

	// Verify all structures have correct number of residues
	for i, protein := range ensemble {
		if len(protein.Residues) != 3 {
			t.Errorf("Structure %d: expected 3 residues, got %d", i, len(protein.Residues))
		}
	}

	// Verify structures are different (not identical to initial)
	angles0 := geometry.CalculateRamachandran(initial)
	angles1 := geometry.CalculateRamachandran(ensemble[0])

	if anglesEqual(angles0, angles1, 0.01) {
		t.Error("Sampled structure identical to initial (no exploration)")
	}
}

// TestFibonacciSphereUniformity checks that Fibonacci sphere sampling is well-distributed
func TestFibonacciSphereUniformity(t *testing.T) {
	initial := createTestProtein(1) // Single residue

	config := QuaternionSearchConfig{
		NumSamples:         20,
		SlerpSteps:         1,
		PerturbRadius:      0.5,
		Seed:               42,
		UseFibonacciSphere: true,
	}

	ensemble, err := QuaternionGuidedSearch(initial, config)
	if err != nil {
		t.Fatalf("QuaternionGuidedSearch failed: %v", err)
	}

	// Extract all sampled angles
	var phiValues, psiValues []float64
	for _, protein := range ensemble {
		angles := geometry.CalculateRamachandran(protein)
		if len(angles) > 0 {
			phiValues = append(phiValues, angles[0].Phi)
			psiValues = append(psiValues, angles[0].Psi)
		}
	}

	// Check coverage of Ramachandran space
	// Expect angles to span wide range (not clustered)
	phiRange := maxFloat(phiValues) - minFloat(phiValues)
	psiRange := maxFloat(psiValues) - minFloat(psiValues)

	if phiRange < 1.0 { // At least 1 radian coverage (≈57°)
		t.Errorf("Phi range too small: %.2f radians (expect >1.0)", phiRange)
	}
	if psiRange < 1.0 {
		t.Errorf("Psi range too small: %.2f radians (expect >1.0)", psiRange)
	}
}

// TestSlerpInterpolation verifies smooth interpolation between conformations
func TestSlerpInterpolation(t *testing.T) {
	initial := createTestProtein(2)

	config := QuaternionSearchConfig{
		NumSamples:         1,
		SlerpSteps:         10,
		PerturbRadius:      0.8,
		Seed:               42,
		UseFibonacciSphere: true,
	}

	ensemble, err := QuaternionGuidedSearch(initial, config)
	if err != nil {
		t.Fatalf("QuaternionGuidedSearch failed: %v", err)
	}

	// Verify smooth progression along slerp path
	// Angle changes should be approximately equal (constant angular velocity)
	prevAngles := geometry.CalculateRamachandran(initial)

	for i, protein := range ensemble {
		currentAngles := geometry.CalculateRamachandran(protein)

		// Calculate angle change for first residue
		if i > 0 && len(prevAngles) > 0 && len(currentAngles) > 0 {
			dPhi := math.Abs(currentAngles[0].Phi - prevAngles[0].Phi)
			dPsi := math.Abs(currentAngles[0].Psi - prevAngles[0].Psi)

			// Angle changes should be finite (no jumps)
			if dPhi > math.Pi || dPsi > math.Pi {
				t.Errorf("Step %d: discontinuous angle change (dPhi=%.2f, dPsi=%.2f)", i, dPhi, dPsi)
			}
		}

		prevAngles = currentAngles
	}
}

// TestRandomVsFibonacci compares random and Fibonacci sampling
func TestRandomVsFibonacci(t *testing.T) {
	initial := createTestProtein(2)

	// Fibonacci sampling
	configFib := QuaternionSearchConfig{
		NumSamples:         10,
		SlerpSteps:         1,
		PerturbRadius:      0.5,
		Seed:               42,
		UseFibonacciSphere: true,
	}

	ensembleFib, err := QuaternionGuidedSearch(initial, configFib)
	if err != nil {
		t.Fatalf("Fibonacci sampling failed: %v", err)
	}

	// Random sampling
	configRand := configFib
	configRand.UseFibonacciSphere = false

	ensembleRand, err := QuaternionGuidedSearch(initial, configRand)
	if err != nil {
		t.Fatalf("Random sampling failed: %v", err)
	}

	// Both should generate same number of structures
	if len(ensembleFib) != len(ensembleRand) {
		t.Errorf("Different ensemble sizes: Fibonacci=%d, Random=%d", len(ensembleFib), len(ensembleRand))
	}

	// Both should be different from initial
	initialAngles := geometry.CalculateRamachandran(initial)
	fibAngles := geometry.CalculateRamachandran(ensembleFib[0])
	randAngles := geometry.CalculateRamachandran(ensembleRand[0])

	if anglesEqual(initialAngles, fibAngles, 0.01) {
		t.Error("Fibonacci sampling produced identical structure")
	}
	if anglesEqual(initialAngles, randAngles, 0.01) {
		t.Error("Random sampling produced identical structure")
	}
}

// TestBuildStructureFromAngles verifies structure building from angles
func TestBuildStructureFromAngles(t *testing.T) {
	template := createTestProtein(3)

	// Define specific angles
	angles := []geometry.RamachandranAngles{
		{Phi: -60.0 * math.Pi / 180.0, Psi: -45.0 * math.Pi / 180.0}, // Alpha helix
		{Phi: -60.0 * math.Pi / 180.0, Psi: -45.0 * math.Pi / 180.0},
		{Phi: -60.0 * math.Pi / 180.0, Psi: -45.0 * math.Pi / 180.0},
	}

	protein, err := buildStructureFromAngles(template, angles)
	if err != nil {
		t.Fatalf("buildStructureFromAngles failed: %v", err)
	}

	// Verify correct number of residues
	if len(protein.Residues) != 3 {
		t.Errorf("Expected 3 residues, got %d", len(protein.Residues))
	}

	// Verify backbone atoms present
	for i, res := range protein.Residues {
		if res.N == nil {
			t.Errorf("Residue %d: missing N atom", i)
		}
		if res.CA == nil {
			t.Errorf("Residue %d: missing CA atom", i)
		}
		if res.C == nil {
			t.Errorf("Residue %d: missing C atom", i)
		}
		if res.O == nil {
			t.Errorf("Residue %d: missing O atom", i)
		}
	}

	// Verify bond lengths are reasonable (N-CA ≈ 1.46 Å)
	for i, res := range protein.Residues {
		if res.N != nil && res.CA != nil {
			dist := distance(res.N, res.CA)
			if dist < 1.0 || dist > 2.0 {
				t.Errorf("Residue %d: N-CA bond length %.2f Å (expect ≈1.46 Å)", i, dist)
			}
		}
	}
}

// TestPerturbRadiusEffect verifies that perturbRadius controls exploration distance
func TestPerturbRadiusEffect(t *testing.T) {
	initial := createTestProtein(2)
	initialAngles := geometry.CalculateRamachandran(initial)

	// Small perturbation
	configSmall := QuaternionSearchConfig{
		NumSamples:         5,
		SlerpSteps:         5,
		PerturbRadius:      0.1, // Small exploration
		Seed:               42,
		UseFibonacciSphere: true,
	}

	ensembleSmall, _ := QuaternionGuidedSearch(initial, configSmall)

	// Large perturbation
	configLarge := QuaternionSearchConfig{
		NumSamples:         5,
		SlerpSteps:         5,
		PerturbRadius:      0.9, // Large exploration
		Seed:               42,
		UseFibonacciSphere: true,
	}

	ensembleLarge, _ := QuaternionGuidedSearch(initial, configLarge)

	// Calculate average deviation from initial
	avgDeviationSmall := averageAngleDeviation(initialAngles, ensembleSmall)
	avgDeviationLarge := averageAngleDeviation(initialAngles, ensembleLarge)

	// Large perturbation should explore farther
	if avgDeviationLarge <= avgDeviationSmall {
		t.Errorf("Large perturbation (%.2f) should explore more than small (%.2f)",
			avgDeviationLarge, avgDeviationSmall)
	}
}

// Helper functions

func createTestProtein(numResidues int) *parser.Protein {
	protein := &parser.Protein{
		Name:     "test",
		Residues: make([]*parser.Residue, numResidues),
		Atoms:    make([]*parser.Atom, 0, numResidues*4),
	}

	atomSerial := 1
	x := 0.0

	for i := 0; i < numResidues; i++ {
		res := &parser.Residue{
			Name:    "ALA",
			SeqNum:  i + 1,
			ChainID: "A",
		}

		// N atom
		res.N = &parser.Atom{
			Serial:  atomSerial,
			Name:    "N",
			ResName: "ALA",
			ChainID: "A",
			ResSeq:  i + 1,
			X:       x,
			Y:       0,
			Z:       0,
			Element: "N",
		}
		atomSerial++
		protein.Atoms = append(protein.Atoms, res.N)
		x += 1.46

		// CA atom
		res.CA = &parser.Atom{
			Serial:  atomSerial,
			Name:    "CA",
			ResName: "ALA",
			ChainID: "A",
			ResSeq:  i + 1,
			X:       x,
			Y:       0,
			Z:       0,
			Element: "C",
		}
		atomSerial++
		protein.Atoms = append(protein.Atoms, res.CA)
		x += 1.52

		// C atom
		res.C = &parser.Atom{
			Serial:  atomSerial,
			Name:    "C",
			ResName: "ALA",
			ChainID: "A",
			ResSeq:  i + 1,
			X:       x,
			Y:       0,
			Z:       0,
			Element: "C",
		}
		atomSerial++
		protein.Atoms = append(protein.Atoms, res.C)

		// O atom
		res.O = &parser.Atom{
			Serial:  atomSerial,
			Name:    "O",
			ResName: "ALA",
			ChainID: "A",
			ResSeq:  i + 1,
			X:       x,
			Y:       1.23,
			Z:       0,
			Element: "O",
		}
		atomSerial++
		protein.Atoms = append(protein.Atoms, res.O)

		protein.Residues[i] = res
		x += 1.33
	}

	return protein
}

func anglesEqual(a1, a2 []geometry.RamachandranAngles, tolerance float64) bool {
	if len(a1) != len(a2) {
		return false
	}

	for i := range a1 {
		if math.Abs(a1[i].Phi-a2[i].Phi) > tolerance {
			return false
		}
		if math.Abs(a1[i].Psi-a2[i].Psi) > tolerance {
			return false
		}
	}

	return true
}

func distance(a1, a2 *parser.Atom) float64 {
	dx := a1.X - a2.X
	dy := a1.Y - a2.Y
	dz := a1.Z - a2.Z
	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}

func minFloat(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	min := values[0]
	for _, v := range values[1:] {
		if v < min {
			min = v
		}
	}
	return min
}

func maxFloat(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	max := values[0]
	for _, v := range values[1:] {
		if v > max {
			max = v
		}
	}
	return max
}

func averageAngleDeviation(initial []geometry.RamachandranAngles, ensemble []*parser.Protein) float64 {
	if len(ensemble) == 0 || len(initial) == 0 {
		return 0
	}

	totalDeviation := 0.0
	count := 0

	for _, protein := range ensemble {
		angles := geometry.CalculateRamachandran(protein)
		if len(angles) > 0 && len(initial) > 0 {
			dPhi := math.Abs(angles[0].Phi - initial[0].Phi)
			dPsi := math.Abs(angles[0].Psi - initial[0].Psi)
			totalDeviation += (dPhi + dPsi) / 2.0
			count++
		}
	}

	if count == 0 {
		return 0
	}

	return totalDeviation / float64(count)
}
