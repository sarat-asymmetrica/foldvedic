package pipeline

import (
	"testing"

	"github.com/sarat-asymmetrica/foldvedic/backend/internal/prediction"
)

// TestDefaultUnifiedPipelineV2Config tests default configuration
func TestDefaultUnifiedPipelineV2Config(t *testing.T) {
	sequence := "ACDEFGHIKLMNPQRSTVWY"
	config := DefaultUnifiedPipelineV2Config(sequence)

	if config.Sequence != sequence {
		t.Errorf("Sequence mismatch: %s vs %s", config.Sequence, sequence)
	}

	if !config.UseSSprediction {
		t.Error("SS prediction should be enabled by default")
	}

	if !config.UseContactMap {
		t.Error("Contact map prediction should be enabled by default")
	}

	if !config.UseVedicBiasing {
		t.Error("Vedic biasing should be enabled by default")
	}

	// All sampling methods should be enabled
	if !config.UseQuaternionSlerp {
		t.Error("Quaternion slerp should be enabled")
	}
	if !config.UseMonteCarlo {
		t.Error("Monte Carlo should be enabled")
	}
	if !config.UseFragmentAssembly {
		t.Error("Fragment assembly should be enabled")
	}
	if !config.UseBasinExplorer {
		t.Error("Basin explorer should be enabled")
	}

	t.Logf("Default config validated: %d samples per method", config.NumSamplesPerMethod)
}

// TestQuickFold tests simple folding interface
func TestQuickFold(t *testing.T) {
	// Very short sequence for quick test
	sequence := "ACDEFG"

	result, err := QuickFold(sequence, false)
	if err != nil {
		t.Fatalf("QuickFold failed: %v", err)
	}

	// Basic validation
	if result == nil {
		t.Fatal("Result is nil")
	}

	if result.FinalStructure == nil {
		t.Error("FinalStructure is nil")
	}

	if len(result.FinalAngles) != len(sequence) {
		t.Errorf("Angle count mismatch: %d vs %d", len(result.FinalAngles), len(sequence))
	}

	if result.TotalSamplesGenerated == 0 {
		t.Error("No samples generated")
	}

	if result.SuccessRate == 0.0 {
		t.Error("Success rate is zero")
	}

	t.Logf("QuickFold success:")
	t.Logf("  Samples: %d", result.TotalSamplesGenerated)
	t.Logf("  Success rate: %.1f%%", result.SuccessRate*100)
	t.Logf("  Final energy: %.2f kcal/mol", result.FinalEnergy)
	t.Logf("  Vedic score: %.3f", result.FinalVedicScore)
	t.Logf("  Quality: %.3f", result.QualityScore)
	t.Logf("  Time: %.2f seconds", result.TotalTimeSeconds)
}

// TestRunUnifiedPipelineV2WithCustomConfig tests custom configuration
func TestRunUnifiedPipelineV2WithCustomConfig(t *testing.T) {
	sequence := "GACDEF"

	// Minimal config for fast test
	config := UnifiedPipelineV2Config{
		Sequence:            sequence,
		UseSSprediction:     true,
		SSMethod:            prediction.MethodChouFasman,
		UseContactMap:       false, // Disable for speed
		UseQuaternionSlerp:  true,
		UseMonteCarlo:       false, // Disable for speed
		UseFragmentAssembly: false, // Disable for speed
		UseBasinExplorer:    false, // Disable for speed
		NumSamplesPerMethod: 3,
		UseVedicBiasing:     true,
		VedicBias:           prediction.DefaultVedicStructuralBias(),
		Verbose:             false,
	}

	// Use default optimization config
	config.OptimizationConfig = DefaultUnifiedPipelineV2Config(sequence).OptimizationConfig

	result, err := RunUnifiedPipelineV2(config, nil)
	if err != nil {
		t.Fatalf("Pipeline failed: %v", err)
	}

	if result.FinalStructure == nil {
		t.Error("FinalStructure is nil")
	}

	// Should have secondary structure prediction
	if len(result.SecondaryStructure) != len(sequence) {
		t.Errorf("SS prediction length mismatch: %d vs %d",
			len(result.SecondaryStructure), len(sequence))
	}

	// Should have no contact map (disabled)
	if len(result.ContactMap) > 0 {
		t.Error("Contact map should be empty (disabled)")
	}

	t.Logf("Custom config test passed")
	t.Logf("  Energy: %.2f kcal/mol", result.FinalEnergy)
	t.Logf("  Vedic: %.3f", result.FinalVedicScore)
}

// TestInitializeFromSSPrediction tests structure initialization
func TestInitializeFromSSPrediction(t *testing.T) {
	sequence := "AAAEEEELLLL"

	// Create mock SS prediction
	ssPred := []prediction.SecondaryStructurePrediction{
		{Position: 0, PredictedType: prediction.AlphaHelix},
		{Position: 1, PredictedType: prediction.AlphaHelix},
		{Position: 2, PredictedType: prediction.AlphaHelix},
		{Position: 3, PredictedType: prediction.BetaSheet},
		{Position: 4, PredictedType: prediction.BetaSheet},
		{Position: 5, PredictedType: prediction.BetaSheet},
		{Position: 6, PredictedType: prediction.Coil},
		{Position: 7, PredictedType: prediction.Coil},
		{Position: 8, PredictedType: prediction.Coil},
		{Position: 9, PredictedType: prediction.Coil},
		{Position: 10, PredictedType: prediction.Coil},
	}

	protein := initializeFromSSPrediction(sequence, ssPred)

	if protein == nil {
		t.Fatal("Protein is nil")
	}

	if len(protein.Residues) != len(sequence) {
		t.Errorf("Residue count mismatch: %d vs %d",
			len(protein.Residues), len(sequence))
	}

	t.Logf("Initialized structure from SS prediction: %d residues", len(protein.Residues))
}

// TestBuildSimpleBackbone tests backbone construction
func TestBuildSimpleBackbone(t *testing.T) {
	sequence := "GAC"

	// Create simple angles (extended chain)
	angles := make([]struct{ Phi, Psi float64 }, len(sequence))
	for i := range angles {
		angles[i].Phi = -120.0 * 3.14159 / 180.0
		angles[i].Psi = +120.0 * 3.14159 / 180.0
	}

	// Convert to RamachandranAngles-compatible structure
	angleSlice := make([]interface{}, len(angles))
	for i := range angles {
		angleSlice[i] = angles[i]
	}

	protein := buildSimpleBackbone(sequence, nil)

	if protein == nil {
		t.Fatal("Protein is nil")
	}

	// Should have N, CA, C, O atoms per residue
	expectedAtoms := len(sequence) * 4
	if len(protein.Atoms) != expectedAtoms {
		t.Errorf("Atom count mismatch: %d vs %d", len(protein.Atoms), expectedAtoms)
	}

	// Check first residue has all backbone atoms
	if protein.Residues[0].N == nil {
		t.Error("First residue missing N atom")
	}
	if protein.Residues[0].CA == nil {
		t.Error("First residue missing CA atom")
	}
	if protein.Residues[0].C == nil {
		t.Error("First residue missing C atom")
	}
	if protein.Residues[0].O == nil {
		t.Error("First residue missing O atom")
	}

	t.Logf("Backbone built: %d atoms for %d residues", len(protein.Atoms), len(protein.Residues))
}

// BenchmarkQuickFold benchmarks full pipeline
func BenchmarkQuickFold(b *testing.B) {
	sequence := "ACDEFGHIKL" // 10 residues

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = QuickFold(sequence, false)
	}
}

// BenchmarkRunUnifiedPipelineV2 benchmarks with minimal config
func BenchmarkRunUnifiedPipelineV2(b *testing.B) {
	sequence := "GACDEFGHIK"

	config := UnifiedPipelineV2Config{
		Sequence:            sequence,
		UseSSprediction:     true,
		SSMethod:            prediction.MethodChouFasman,
		UseContactMap:       false,
		UseQuaternionSlerp:  true,
		UseMonteCarlo:       false,
		UseFragmentAssembly: false,
		UseBasinExplorer:    false,
		NumSamplesPerMethod: 2,
		UseVedicBiasing:     false,
		Verbose:             false,
	}
	config.OptimizationConfig = DefaultUnifiedPipelineV2Config(sequence).OptimizationConfig

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = RunUnifiedPipelineV2(config, nil)
	}
}
