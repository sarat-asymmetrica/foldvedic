package geometry

import (
	"math"
	"testing"

	"github.com/sarat-asymmetrica/foldvedic/backend/internal/parser"
)

func TestCalculateDihedral(t *testing.T) {
	// Test dihedral angle calculation with known geometry
	// Create four points defining a dihedral angle of ~90 degrees

	// Points forming a right angle dihedral
	p1 := Vector3{X: 0, Y: 0, Z: 0}
	p2 := Vector3{X: 1, Y: 0, Z: 0}
	p3 := Vector3{X: 1, Y: 1, Z: 0}
	p4 := Vector3{X: 1, Y: 1, Z: 1}

	angle := calculateDihedral(p1, p2, p3, p4)

	// Expected angle is approximately ±90 degrees (±π/2 radians)
	// Sign depends on rotation direction, so check magnitude
	expectedMag := math.Pi / 2.0
	tolerance := 0.1 // radians

	if math.Abs(math.Abs(angle)-expectedMag) > tolerance {
		t.Errorf("Expected dihedral angle magnitude ~%.2f, got %.2f", expectedMag, angle)
	}
}

func TestCalculateDihedralPlanar(t *testing.T) {
	// Test dihedral for planar configuration (should be 0 or 180)
	p1 := Vector3{X: 0, Y: 0, Z: 0}
	p2 := Vector3{X: 1, Y: 0, Z: 0}
	p3 := Vector3{X: 2, Y: 0, Z: 0}
	p4 := Vector3{X: 3, Y: 0, Z: 0}

	angle := calculateDihedral(p1, p2, p3, p4)

	// All points are collinear - dihedral is undefined but should be near 0 or π
	// Due to numerical issues, we just check it's a valid number
	if math.IsNaN(angle) {
		t.Error("Dihedral angle should not be NaN for collinear points")
	}
}

func TestVector3Operations(t *testing.T) {
	v1 := Vector3{X: 1, Y: 2, Z: 3}
	v2 := Vector3{X: 4, Y: 5, Z: 6}

	// Test subtraction
	sub := v1.Sub(v2)
	if sub.X != -3 || sub.Y != -3 || sub.Z != -3 {
		t.Errorf("Subtraction failed: got %+v", sub)
	}

	// Test dot product
	dot := v1.Dot(v2)
	expected := 1*4 + 2*5 + 3*6 // = 32
	if dot != float64(expected) {
		t.Errorf("Dot product: expected %d, got %f", expected, dot)
	}

	// Test magnitude
	v := Vector3{X: 3, Y: 4, Z: 0}
	mag := v.Magnitude()
	if math.Abs(mag-5.0) > 0.001 {
		t.Errorf("Magnitude: expected 5.0, got %f", mag)
	}

	// Test cross product (right-hand rule)
	i := Vector3{X: 1, Y: 0, Z: 0}
	j := Vector3{X: 0, Y: 1, Z: 0}
	cross := i.Cross(j)
	// i × j = k
	if cross.X != 0 || cross.Y != 0 || cross.Z != 1 {
		t.Errorf("Cross product i×j: expected (0,0,1), got %+v", cross)
	}

	// Test normalize
	norm := v.Normalize()
	normMag := norm.Magnitude()
	if math.Abs(normMag-1.0) > 0.001 {
		t.Errorf("Normalized magnitude: expected 1.0, got %f", normMag)
	}
}

func TestCalculateRamachandran(t *testing.T) {
	// Load test protein
	protein, err := parser.ParsePDB("../../../testdata/test_peptide.pdb")
	if err != nil {
		t.Fatalf("Failed to parse test PDB: %v", err)
	}

	// Calculate Ramachandran angles
	angles := CalculateRamachandran(protein)

	// Should have 3 residues
	if len(angles) != 3 {
		t.Fatalf("Expected 3 angles, got %d", len(angles))
	}

	// First residue: phi should be NaN (N-terminal)
	if !math.IsNaN(angles[0].Phi) {
		t.Error("First residue phi should be NaN (N-terminal)")
	}

	// Last residue: psi should be NaN (C-terminal)
	if !math.IsNaN(angles[2].Psi) {
		t.Error("Last residue psi should be NaN (C-terminal)")
	}

	// Middle residue: both angles should be defined
	if math.IsNaN(angles[1].Phi) || math.IsNaN(angles[1].Psi) {
		t.Error("Middle residue should have both phi and psi defined")
	}

	// Check angles are in valid range [-π, π]
	for i, a := range angles {
		if !math.IsNaN(a.Phi) {
			if a.Phi < -math.Pi || a.Phi > math.Pi {
				t.Errorf("Residue %d phi out of range: %f", i, a.Phi)
			}
		}
		if !math.IsNaN(a.Psi) {
			if a.Psi < -math.Pi || a.Psi > math.Pi {
				t.Errorf("Residue %d psi out of range: %f", i, a.Psi)
			}
		}
	}
}

func TestRamachandranAngleConversion(t *testing.T) {
	// Test degree conversion
	angles := RamachandranAngles{
		Phi: -math.Pi / 3, // -60 degrees
		Psi: math.Pi / 4,  // 45 degrees
	}

	phiDeg := angles.ToDegressPhi()
	psiDeg := angles.ToDegressPsi()

	if math.Abs(phiDeg-(-60.0)) > 0.1 {
		t.Errorf("Phi conversion: expected -60°, got %.2f°", phiDeg)
	}
	if math.Abs(psiDeg-45.0) > 0.1 {
		t.Errorf("Psi conversion: expected 45°, got %.2f°", psiDeg)
	}
}

func TestIsInAllowedRegion(t *testing.T) {
	tests := []struct {
		name     string
		phi      float64 // degrees
		psi      float64 // degrees
		expected bool
	}{
		{
			name:     "Alpha helix",
			phi:      -60,
			psi:      -45,
			expected: true,
		},
		{
			name:     "Beta sheet",
			phi:      -120,
			psi:      +120,
			expected: true,
		},
		{
			name:     "Left-handed helix (rare)",
			phi:      +60,
			psi:      +45,
			expected: true,
		},
		{
			name:     "Forbidden region",
			phi:      +90,
			psi:      -90,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Convert degrees to radians
			angles := RamachandranAngles{
				Phi: tt.phi * math.Pi / 180,
				Psi: tt.psi * math.Pi / 180,
			}

			result := angles.IsInAllowedRegion()
			if result != tt.expected {
				t.Errorf("IsInAllowedRegion(φ=%.0f°, ψ=%.0f°) = %v, expected %v",
					tt.phi, tt.psi, result, tt.expected)
			}
		})
	}
}
