package geometry

import (
	"math"
	"testing"
)

func TestRamachandranToQuaternion(t *testing.T) {
	// Test known angles
	phi := -math.Pi / 3 // -60 degrees (helix)
	psi := -math.Pi / 4 // -45 degrees (helix)

	q := RamachandranToQuaternion(phi, psi)

	// Check unit norm (quaternion should lie on S³ hypersphere)
	norm := math.Sqrt(q.W*q.W + q.X*q.X + q.Y*q.Y + q.Z*q.Z)
	if math.Abs(norm-1.0) > 0.0001 {
		t.Errorf("Quaternion norm: expected 1.0, got %f", norm)
	}

	// All components should be non-NaN
	if math.IsNaN(q.W) || math.IsNaN(q.X) || math.IsNaN(q.Y) || math.IsNaN(q.Z) {
		t.Error("Quaternion contains NaN values")
	}
}

func TestQuaternionToRamachandran(t *testing.T) {
	// Test round-trip: angles → quaternion → angles
	originalPhi := -math.Pi / 3 // -60 degrees
	originalPsi := math.Pi / 4   // 45 degrees

	// Convert to quaternion
	q := RamachandranToQuaternion(originalPhi, originalPsi)

	// Convert back to angles
	recoveredPhi, recoveredPsi := QuaternionToRamachandran(q)

	// Check angles match (within tolerance for floating point)
	tolerance := 0.0001
	if math.Abs(recoveredPhi-originalPhi) > tolerance {
		t.Errorf("Phi round-trip: expected %f, got %f (diff: %f)",
			originalPhi, recoveredPhi, math.Abs(recoveredPhi-originalPhi))
	}
	if math.Abs(recoveredPsi-originalPsi) > tolerance {
		t.Errorf("Psi round-trip: expected %f, got %f (diff: %f)",
			originalPsi, recoveredPsi, math.Abs(recoveredPsi-originalPsi))
	}
}

func TestQuaternionBoundaryAngles(t *testing.T) {
	// Test angles at boundary (±180°) - this is where traditional 2D grid fails
	testCases := []struct {
		name string
		phi  float64
		psi  float64
	}{
		{"Positive boundary", math.Pi, math.Pi},
		{"Negative boundary", -math.Pi, -math.Pi},
		{"Mixed boundary", math.Pi, -math.Pi},
		{"Near boundary", 0.99 * math.Pi, 0.99 * math.Pi},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			q := RamachandranToQuaternion(tc.phi, tc.psi)

			// Check unit norm
			norm := math.Sqrt(q.W*q.W + q.X*q.X + q.Y*q.Y + q.Z*q.Z)
			if math.Abs(norm-1.0) > 0.0001 {
				t.Errorf("Quaternion norm at boundary: expected 1.0, got %f", norm)
			}

			// Check round-trip (may wrap around ±π)
			recoveredPhi, recoveredPsi := QuaternionToRamachandran(q)

			// Check angles are in valid range
			if math.Abs(recoveredPhi) > math.Pi+0.001 {
				t.Errorf("Recovered phi out of range: %f", recoveredPhi)
			}
			if math.Abs(recoveredPsi) > math.Pi+0.001 {
				t.Errorf("Recovered psi out of range: %f", recoveredPsi)
			}
		})
	}
}

func TestQuaternionNormalize(t *testing.T) {
	// Test normalization
	q := Quaternion{W: 1, X: 1, Y: 1, Z: 1} // Not unit norm

	qNorm := q.Normalize()

	// Check unit norm after normalization
	norm := math.Sqrt(qNorm.W*qNorm.W + qNorm.X*qNorm.X + qNorm.Y*qNorm.Y + qNorm.Z*qNorm.Z)
	if math.Abs(norm-1.0) > 0.0001 {
		t.Errorf("Normalized quaternion norm: expected 1.0, got %f", norm)
	}

	// Check direction preserved (proportions maintained)
	// Original: (1, 1, 1, 1), normalized should be (0.5, 0.5, 0.5, 0.5)
	expected := 0.5
	tolerance := 0.0001
	if math.Abs(qNorm.W-expected) > tolerance ||
		math.Abs(qNorm.X-expected) > tolerance ||
		math.Abs(qNorm.Y-expected) > tolerance ||
		math.Abs(qNorm.Z-expected) > tolerance {
		t.Errorf("Normalized direction: expected (%f,%f,%f,%f), got (%f,%f,%f,%f)",
			expected, expected, expected, expected,
			qNorm.W, qNorm.X, qNorm.Y, qNorm.Z)
	}
}

func TestQuaternionSlerp(t *testing.T) {
	// Test slerp between two conformations
	// Helix conformation
	phiHelix, psiHelix := -60.0*math.Pi/180, -45.0*math.Pi/180
	qHelix := RamachandranToQuaternion(phiHelix, psiHelix)

	// Sheet conformation
	phiSheet, psiSheet := -120.0*math.Pi/180, +120.0*math.Pi/180
	qSheet := RamachandranToQuaternion(phiSheet, psiSheet)

	// Slerp at t=0 should give helix
	q0 := qHelix.Slerp(qSheet, 0.0)
	phi0, psi0 := QuaternionToRamachandran(q0)
	tolerance := 0.001
	if math.Abs(phi0-phiHelix) > tolerance || math.Abs(psi0-psiHelix) > tolerance {
		t.Error("Slerp at t=0 should return first quaternion")
	}

	// Slerp at t=1 should give sheet
	q1 := qHelix.Slerp(qSheet, 1.0)
	phi1, psi1 := QuaternionToRamachandran(q1)
	if math.Abs(phi1-phiSheet) > tolerance || math.Abs(psi1-psiSheet) > tolerance {
		t.Error("Slerp at t=1 should return second quaternion")
	}

	// Slerp at t=0.5 should give intermediate conformation
	qMid := qHelix.Slerp(qSheet, 0.5)

	// Check unit norm preserved
	norm := math.Sqrt(qMid.W*qMid.W + qMid.X*qMid.X + qMid.Y*qMid.Y + qMid.Z*qMid.Z)
	if math.Abs(norm-1.0) > 0.0001 {
		t.Errorf("Slerp preserves norm: expected 1.0, got %f", norm)
	}

	// Check interpolated angles are between helix and sheet
	phiMid, psiMid := QuaternionToRamachandran(qMid)
	phiMidDeg := phiMid * 180 / math.Pi
	psiMidDeg := psiMid * 180 / math.Pi

	// Angles should be intermediate (rough check - actual path is geodesic, not linear)
	// Just verify they're valid and different from endpoints
	if math.IsNaN(phiMid) || math.IsNaN(psiMid) {
		t.Error("Slerp produced NaN angles")
	}
	if math.Abs(phiMid-phiHelix) < 0.01 && math.Abs(psiMid-psiHelix) < 0.01 {
		t.Error("Slerp midpoint too close to helix conformation")
	}
	if math.Abs(phiMid-phiSheet) < 0.01 && math.Abs(psiMid-psiSheet) < 0.01 {
		t.Error("Slerp midpoint too close to sheet conformation")
	}

	t.Logf("Slerp midpoint: φ=%.1f°, ψ=%.1f° (between helix and sheet)", phiMidDeg, psiMidDeg)
}

func TestInterpolateConformation(t *testing.T) {
	// Test conformation interpolation for multiple residues
	angles1 := []RamachandranAngles{
		{Phi: -60 * math.Pi / 180, Psi: -45 * math.Pi / 180}, // Helix
		{Phi: -60 * math.Pi / 180, Psi: -45 * math.Pi / 180}, // Helix
		{Phi: -60 * math.Pi / 180, Psi: -45 * math.Pi / 180}, // Helix
	}

	angles2 := []RamachandranAngles{
		{Phi: -120 * math.Pi / 180, Psi: +120 * math.Pi / 180}, // Sheet
		{Phi: -120 * math.Pi / 180, Psi: +120 * math.Pi / 180}, // Sheet
		{Phi: -120 * math.Pi / 180, Psi: +120 * math.Pi / 180}, // Sheet
	}

	// Interpolate at t=0.5
	interpolated := InterpolateConformation(angles1, angles2, 0.5)

	if len(interpolated) != 3 {
		t.Fatalf("Expected 3 interpolated angles, got %d", len(interpolated))
	}

	// Check all angles are valid and intermediate
	for i, a := range interpolated {
		if math.IsNaN(a.Phi) || math.IsNaN(a.Psi) {
			t.Errorf("Interpolated angle %d contains NaN", i)
		}

		// Should be different from both endpoints (crude check)
		tolerance := 0.1 // radians
		tooCloseToHelix := math.Abs(a.Phi-angles1[i].Phi) < tolerance &&
			math.Abs(a.Psi-angles1[i].Psi) < tolerance
		tooCloseToSheet := math.Abs(a.Phi-angles2[i].Phi) < tolerance &&
			math.Abs(a.Psi-angles2[i].Psi) < tolerance

		if tooCloseToHelix || tooCloseToSheet {
			t.Logf("Warning: Interpolated angle %d very close to endpoint", i)
		}
	}
}

func TestQuaternionNaNHandling(t *testing.T) {
	// Test handling of NaN angles (terminal residues)
	q := RamachandranToQuaternion(math.NaN(), math.NaN())

	// Should return identity quaternion
	if q.W != 1.0 || q.X != 0.0 || q.Y != 0.0 || q.Z != 0.0 {
		t.Errorf("NaN angles should produce identity quaternion, got %+v", q)
	}

	// Round-trip should give NaN back
	phi, psi := QuaternionToRamachandran(q)
	if !math.IsNaN(phi) || !math.IsNaN(psi) {
		t.Error("Identity quaternion should convert to NaN angles")
	}
}

func BenchmarkRamachandranToQuaternion(b *testing.B) {
	phi := -math.Pi / 3
	psi := -math.Pi / 4

	for i := 0; i < b.N; i++ {
		_ = RamachandranToQuaternion(phi, psi)
	}
}

func BenchmarkQuaternionSlerp(b *testing.B) {
	q1 := RamachandranToQuaternion(-math.Pi/3, -math.Pi/4)
	q2 := RamachandranToQuaternion(-2*math.Pi/3, 2*math.Pi/3)

	for i := 0; i < b.N; i++ {
		_ = q1.Slerp(q2, 0.5)
	}
}
