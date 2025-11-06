package physics

import (
	"math"
	"testing"

	"github.com/sarat-asymmetrica/foldvedic/backend/internal/parser"
)

func TestCalculateBondEnergy(t *testing.T) {
	// Create two atoms at equilibrium distance
	atom1 := &parser.Atom{X: 0, Y: 0, Z: 0}
	atom2 := &parser.Atom{X: 1.335, Y: 0, Z: 0} // Peptide bond length

	params := BondParameters{
		K0: 490.0, // C-N force constant
		R0: 1.335, // Equilibrium length
	}

	energy := CalculateBondEnergy(atom1, atom2, params)

	// At equilibrium, energy should be zero
	if math.Abs(energy) > 0.001 {
		t.Errorf("Energy at equilibrium should be ~0, got %.6f", energy)
	}

	// Test stretched bond
	atom2.X = 1.5 // 0.165 Å stretch
	energy = CalculateBondEnergy(atom1, atom2, params)

	// E = k × (dr)² = 490 × (0.165)² ≈ 13.3 kcal/mol
	expectedEnergy := 490.0 * 0.165 * 0.165
	tolerance := 0.1

	if math.Abs(energy-expectedEnergy) > tolerance {
		t.Errorf("Stretched bond energy: expected %.2f, got %.2f", expectedEnergy, energy)
	}
}

func TestCalculateBondForce(t *testing.T) {
	// Two atoms with bond along X axis
	atom1 := &parser.Atom{X: 0, Y: 0, Z: 0}
	atom2 := &parser.Atom{X: 1.5, Y: 0, Z: 0} // 0.15 Å longer than equilibrium

	params := BondParameters{
		K0: 490.0,
		R0: 1.335,
	}

	force := CalculateBondForce(atom1, atom2, params)

	// Force should be along X axis
	if math.Abs(force.Y) > 0.001 || math.Abs(force.Z) > 0.001 {
		t.Error("Force should be along X axis only")
	}

	// Force magnitude: F = -2k(r-r0) = -2 × 490 × 0.165 ≈ -161.7
	dr := 1.5 - 1.335
	expectedForceMag := -2.0 * 490.0 * dr
	tolerance := 0.1

	if math.Abs(force.X-expectedForceMag) > tolerance {
		t.Errorf("Force magnitude: expected %.2f, got %.2f", expectedForceMag, force.X)
	}
}

func TestCalculateAngleEnergy(t *testing.T) {
	// Three atoms forming 90 degree angle
	atom1 := &parser.Atom{X: 1, Y: 0, Z: 0}
	atom2 := &parser.Atom{X: 0, Y: 0, Z: 0} // Central atom
	atom3 := &parser.Atom{X: 0, Y: 1, Z: 0}

	params := AngleParameters{
		K0:     63.0,                  // N-CA-C force constant
		Theta0: 110.1 * math.Pi / 180, // Equilibrium angle
	}

	energy := CalculateAngleEnergy(atom1, atom2, atom3, params)

	// Current angle is 90°, equilibrium is 110.1°
	// dθ = 90 - 110.1 = -20.1° = -0.351 rad
	// E = k × (dθ)² = 63 × (0.351)² ≈ 7.75 kcal/mol

	// Should be positive (deviation from equilibrium)
	if energy < 0 {
		t.Error("Angle energy should be positive (deviation from equilibrium)")
	}

	if energy < 5 || energy > 10 {
		t.Errorf("Angle energy out of expected range: %.2f kcal/mol", energy)
	}
}

func TestLennardJonesEnergy(t *testing.T) {
	// Two carbon atoms
	atom1 := &parser.Atom{X: 0, Y: 0, Z: 0, Element: "C"}
	atom2 := &parser.Atom{X: 4.0, Y: 0, Z: 0, Element: "C"} // 4 Å apart

	cutoff := 12.0 // Å

	energy := CalculateLennardJonesEnergy(atom1, atom2, cutoff)

	// At typical VdW distances (3-5 Å), LJ energy should be negative (attractive)
	if energy > 0 {
		t.Error("LJ energy should be negative (attractive) at this distance")
	}

	t.Logf("LJ energy at 4 Å: %.4f kcal/mol", energy)

	// Test at very short distance (strong repulsion)
	atom2.X = 1.0 // 1 Å apart (much less than σ)
	energyRepulsive := CalculateLennardJonesEnergy(atom1, atom2, cutoff)

	// Should be positive (repulsive)
	if energyRepulsive < 0 {
		t.Error("LJ energy should be positive (repulsive) at very short distance")
	}

	t.Logf("LJ energy at 1 Å: %.4f kcal/mol (repulsive)", energyRepulsive)
}

func TestElectrostaticEnergy(t *testing.T) {
	// Two oppositely charged atoms
	atom1 := &parser.Atom{X: 0, Y: 0, Z: 0}
	atom2 := &parser.Atom{X: 5.0, Y: 0, Z: 0} // 5 Å apart

	charge1 := 0.5973  // Carbonyl carbon (positive)
	charge2 := -0.5679 // Carbonyl oxygen (negative)
	cutoff := 15.0     // Å

	energy := CalculateElectrostaticEnergy(atom1, atom2, charge1, charge2, cutoff)

	// Opposite charges should give negative (attractive) energy
	if energy > 0 {
		t.Error("Electrostatic energy should be negative (attractive) for opposite charges")
	}

	t.Logf("Electrostatic energy at 5 Å: %.4f kcal/mol", energy)

	// Test same charges (repulsive)
	charge2Positive := 0.5973
	energyRepulsive := CalculateElectrostaticEnergy(atom1, atom2, charge1, charge2Positive, cutoff)

	if energyRepulsive < 0 {
		t.Error("Electrostatic energy should be positive (repulsive) for same charges")
	}

	t.Logf("Electrostatic energy (same charge) at 5 Å: %.4f kcal/mol (repulsive)", energyRepulsive)
}

func TestGetBondParams(t *testing.T) {
	// Test known bond type
	params := GetBondParams("C", "N")

	if params.R0 != 1.335 {
		t.Errorf("C-N equilibrium length: expected 1.335 Å, got %.3f Å", params.R0)
	}

	if params.K0 != 490.0 {
		t.Errorf("C-N force constant: expected 490.0, got %.1f", params.K0)
	}

	// Test reverse ordering
	paramsReverse := GetBondParams("N", "C")
	if paramsReverse.R0 != params.R0 {
		t.Error("Bond parameters should be symmetric")
	}

	// Test unknown bond type (should return default)
	paramsUnknown := GetBondParams("X", "Y")
	if paramsUnknown.R0 != 1.5 {
		t.Errorf("Unknown bond should return default R0=1.5, got %.3f", paramsUnknown.R0)
	}
}

func TestGetAngleParams(t *testing.T) {
	// Test known angle type
	params := GetAngleParams("N", "CA", "C")

	expectedTheta := 110.1 * math.Pi / 180 // radians
	tolerance := 0.01

	if math.Abs(params.Theta0-expectedTheta) > tolerance {
		t.Errorf("N-CA-C equilibrium angle: expected %.3f rad, got %.3f rad",
			expectedTheta, params.Theta0)
	}

	// Test reverse ordering
	paramsReverse := GetAngleParams("C", "CA", "N")
	if math.Abs(paramsReverse.Theta0-params.Theta0) > tolerance {
		t.Error("Angle parameters should be symmetric")
	}
}

func BenchmarkCalculateBondEnergy(b *testing.B) {
	atom1 := &parser.Atom{X: 0, Y: 0, Z: 0}
	atom2 := &parser.Atom{X: 1.335, Y: 0, Z: 0}
	params := BondParameters{K0: 490.0, R0: 1.335}

	for i := 0; i < b.N; i++ {
		_ = CalculateBondEnergy(atom1, atom2, params)
	}
}

func BenchmarkCalculateLennardJones(b *testing.B) {
	atom1 := &parser.Atom{X: 0, Y: 0, Z: 0, Element: "C"}
	atom2 := &parser.Atom{X: 4.0, Y: 0, Z: 0, Element: "C"}
	cutoff := 12.0

	for i := 0; i < b.N; i++ {
		_ = CalculateLennardJonesEnergy(atom1, atom2, cutoff)
	}
}
