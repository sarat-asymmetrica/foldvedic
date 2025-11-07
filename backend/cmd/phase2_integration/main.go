// Phase 2 Integration - Advanced Sampling Pipeline
//
// 🧬 FOLDVEDIC PHASE 2: ADVANCED SAMPLING
//
// Mission: Generate 100+ diverse protein structures using 4 advanced sampling methods
// Target: RMSD <15 Å (2× improvement from Phase 1: 26.45 Å)
// Quality: ≥0.92 (LEGENDARY tier)
//
// The 4 Agents:
// 1. Fibonacci Sphere Sampling (Golden angle uniform sampling on S³)
// 2. Vedic Monte Carlo (Metropolis-Hastings with digital root biasing)
// 3. Fragment Assembly (Rosetta-style with Vedic ranking)
// 4. Ramachandran Basin Explorer (Systematic exploration of allowed regions)
//
// Philosophy: Wright Brothers + Quaternion-first + Cross-domain fearlessness
package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"sort"
	"time"

	"github.com/sarat-asymmetrica/foldvedic/backend/internal/geometry"
	"github.com/sarat-asymmetrica/foldvedic/backend/internal/optimization"
	"github.com/sarat-asymmetrica/foldvedic/backend/internal/parser"
	"github.com/sarat-asymmetrica/foldvedic/backend/internal/sampling"
	"github.com/sarat-asymmetrica/foldvedic/backend/internal/validation"
	"github.com/sarat-asymmetrica/foldvedic/backend/internal/vedic"
)

// Phase2Result holds comprehensive Phase 2 results
type Phase2Result struct {
	// Metadata
	ProteinName    string    `json:"protein_name"`
	Sequence       string    `json:"sequence"`
	NumResidues    int       `json:"num_residues"`
	Timestamp      time.Time `json:"timestamp"`
	TotalDuration  float64   `json:"total_duration_seconds"`

	// Sampling statistics
	TotalStructures       int     `json:"total_structures"`
	FibonacciStructures   int     `json:"fibonacci_structures"`
	MonteCarloStructures  int     `json:"monte_carlo_structures"`
	FragmentStructures    int     `json:"fragment_structures"`
	BasinStructures       int     `json:"basin_structures"`

	// RMSD statistics (against experimental)
	BestRMSD      float64   `json:"best_rmsd_angstrom"`
	MedianRMSD    float64   `json:"median_rmsd_angstrom"`
	MeanRMSD      float64   `json:"mean_rmsd_angstrom"`
	WorstRMSD     float64   `json:"worst_rmsd_angstrom"`
	RMSDStdDev    float64   `json:"rmsd_std_dev"`
	RMSDImprovement float64 `json:"rmsd_improvement_vs_phase1"` // vs 26.45 Å

	// Energy statistics
	BestEnergy    float64 `json:"best_energy_kcal_mol"`
	MedianEnergy  float64 `json:"median_energy_kcal_mol"`
	MeanEnergy    float64 `json:"mean_energy_kcal_mol"`
	WorstEnergy   float64 `json:"worst_energy_kcal_mol"`

	// Vedic statistics
	BestVedic     float64 `json:"best_vedic_score"`
	MedianVedic   float64 `json:"median_vedic_score"`
	MeanVedic     float64 `json:"mean_vedic_score"`

	// Validation metrics (best structure)
	BestTMScore   float64 `json:"best_tm_score"`
	BestGDT_TS    float64 `json:"best_gdt_ts"`

	// Quality assessment
	QualityScore       float64 `json:"quality_score"`
	QualityTier        string  `json:"quality_tier"`
	MissionAccomplished bool   `json:"mission_accomplished"` // RMSD <15 Å && Quality ≥0.92

	// Sampling method performance
	BestMethod    string  `json:"best_sampling_method"`
	BestMethodRMSD float64 `json:"best_method_rmsd"`

	// Detailed structure metrics
	Structures []StructureMetric `json:"structures"`
}

// StructureMetric holds metrics for a single structure
type StructureMetric struct {
	ID              int     `json:"id"`
	SamplingMethod  string  `json:"sampling_method"`
	RMSD            float64 `json:"rmsd_angstrom"`
	Energy          float64 `json:"energy_kcal_mol"`
	VedicScore      float64 `json:"vedic_score"`
	TMScore         float64 `json:"tm_score"`
	GDT_TS          float64 `json:"gdt_ts"`
	OptimizationSteps int   `json:"optimization_steps"`
}

func main() {
	fmt.Println("╔══════════════════════════════════════════════════════════════════╗")
	fmt.Println("║  🧬 FOLDVEDIC PHASE 2: ADVANCED SAMPLING                        ║")
	fmt.Println("║                                                                  ║")
	fmt.Println("║  Mission: Generate 100+ diverse structures                      ║")
	fmt.Println("║  Target: RMSD <15 Å (2× improvement from 26.45 Å)              ║")
	fmt.Println("║  Quality: ≥0.92 (LEGENDARY tier)                                ║")
	fmt.Println("║                                                                  ║")
	fmt.Println("║  The 4 Agents:                                                   ║")
	fmt.Println("║  1. Fibonacci Sphere Sampling (Golden angle on S³)              ║")
	fmt.Println("║  2. Vedic Monte Carlo (Digital root biasing)                    ║")
	fmt.Println("║  3. Fragment Assembly (Rosetta-style)                           ║")
	fmt.Println("║  4. Ramachandran Basin Explorer (Systematic basins)             ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════════╝")
	fmt.Println()

	startTime := time.Now()

	// Test protein: Trp-cage (1L2Y) - 20 residues, well-characterized
	sequence := "NLYIQWLKDGGPSSGRPPPS"
	proteinName := "Trp-cage (1L2Y)"

	fmt.Printf("🎯 Target: %s\n", proteinName)
	fmt.Printf("📝 Sequence: %s\n", sequence)
	fmt.Printf("📏 Length: %d residues\n\n", len(sequence))

	// Load experimental structure for validation
	fmt.Println("Loading experimental structure...")
	experimental, err := parser.ParsePDB("/home/user/foldvedic/testdata/1L2Y.pdb")
	if err != nil {
		fmt.Printf("⚠️  Warning: Could not load experimental structure: %v\n", err)
		fmt.Println("Continuing with structure generation only...\n")
		experimental = nil
	} else {
		fmt.Printf("✅ Loaded: %d residues, %d atoms\n\n", len(experimental.Residues), len(experimental.Atoms))
	}

	// Initialize Phase 2 result
	result := &Phase2Result{
		ProteinName:  proteinName,
		Sequence:     sequence,
		NumResidues:  len(sequence),
		Timestamp:    time.Now(),
		Structures:   make([]StructureMetric, 0, 100),
	}

	// Build initial extended structure
	fmt.Println("Building initial extended structure...")
	initialAngles := make([]geometry.RamachandranAngles, len(sequence))
	for i := range initialAngles {
		initialAngles[i] = geometry.RamachandranAngles{
			Phi: -120.0 * math.Pi / 180.0,
			Psi: +120.0 * math.Pi / 180.0,
		}
	}
	initialStructure, err := geometry.BuildProteinFromAngles(sequence, initialAngles)
	if err != nil {
		fmt.Printf("❌ ERROR: Failed to build initial structure: %v\n", err)
		return
	}
	fmt.Println("✅ Initial structure ready\n")

	// ==================== AGENT 2.1: FIBONACCI SPHERE SAMPLING ====================
	fmt.Println("╔══════════════════════════════════════════════════════════════════╗")
	fmt.Println("║  AGENT 2.1: FIBONACCI SPHERE SAMPLING                            ║")
	fmt.Println("║  Golden angle (137.5°) uniform sampling on S³ hypersphere       ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════════╝")
	fmt.Println()

	fibConfig := sampling.DefaultQuaternionSearchConfig()
	fibConfig.NumSamples = 25 // 25 diverse samples
	fibConfig.SlerpSteps = 1  // No interpolation, just endpoints
	fibConfig.UseFibonacciSphere = true
	fibConfig.PerturbRadius = 0.3 // Moderate exploration

	fibStructures, err := sampling.QuaternionGuidedSearch(initialStructure, fibConfig)
	if err != nil {
		fmt.Printf("⚠️  Fibonacci sampling warning: %v\n", err)
		fibStructures = []*parser.Protein{}
	}

	fmt.Printf("✅ Generated %d structures via Fibonacci sphere sampling\n", len(fibStructures))
	result.FibonacciStructures = len(fibStructures)

	// Optimize and evaluate each Fibonacci structure
	for i, structure := range fibStructures {
		optimized, metrics := optimizeAndEvaluate(structure, experimental, fmt.Sprintf("fibonacci_%d", i))
		if metrics != nil {
			metrics.ID = len(result.Structures)
			metrics.SamplingMethod = "Fibonacci Sphere"
			result.Structures = append(result.Structures, *metrics)
		}
		_ = optimized
	}
	fmt.Println()

	// ==================== AGENT 2.2: VEDIC MONTE CARLO ====================
	fmt.Println("╔══════════════════════════════════════════════════════════════════╗")
	fmt.Println("║  AGENT 2.2: VEDIC MONTE CARLO                                    ║")
	fmt.Println("║  Metropolis-Hastings with digital root biasing                  ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════════╝")
	fmt.Println()

	mcConfig := sampling.DefaultMonteCarloConfig()
	mcConfig.NumSteps = 1000          // 1000 MC steps
	mcConfig.TemperatureInitial = 500.0
	mcConfig.TemperatureFinal = 10.0
	mcConfig.VedicWeight = 0.3        // 30% Vedic influence

	// Run 10 independent MC runs
	mcStructures := make([]*parser.Protein, 0, 10)
	for run := 0; run < 10; run++ {
		mcConfig.Seed = 42 + int64(run)
		mcResult, err := sampling.MonteCarloVedic(initialStructure, mcConfig)
		if err != nil {
			fmt.Printf("⚠️  MC run %d warning: %v\n", run, err)
			continue
		}
		mcStructures = append(mcStructures, mcResult.FinalStructure)
		fmt.Printf("  Run %d: Energy %.2f kcal/mol, Vedic %.3f, Acceptance %.2f%%\n",
			run+1, mcResult.FinalEnergy, mcResult.FinalVedicScore, mcResult.AcceptanceRate*100)
	}

	fmt.Printf("✅ Generated %d structures via Vedic Monte Carlo\n", len(mcStructures))
	result.MonteCarloStructures = len(mcStructures)

	// Optimize and evaluate each MC structure
	for i, structure := range mcStructures {
		optimized, metrics := optimizeAndEvaluate(structure, experimental, fmt.Sprintf("monte_carlo_%d", i))
		if metrics != nil {
			metrics.ID = len(result.Structures)
			metrics.SamplingMethod = "Vedic Monte Carlo"
			result.Structures = append(result.Structures, *metrics)
		}
		_ = optimized
	}
	fmt.Println()

	// ==================== AGENT 2.3: FRAGMENT ASSEMBLY ====================
	fmt.Println("╔══════════════════════════════════════════════════════════════════╗")
	fmt.Println("║  AGENT 2.3: FRAGMENT ASSEMBLY                                    ║")
	fmt.Println("║  Rosetta-style assembly with Vedic fragment ranking             ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════════╝")
	fmt.Println()

	fragLib := sampling.NewFragmentLibrary()
	fragConfig := sampling.DefaultFragmentAssemblyConfig()
	fragConfig.UseThreeMers = true
	fragConfig.UseNineMers = true
	fragConfig.VedicWeight = 0.3

	// Generate 25 diverse fragment assemblies
	fragStructures, err := sampling.GenerateFragmentEnsemble(sequence, fragLib, fragConfig, 25)
	if err != nil {
		fmt.Printf("⚠️  Fragment assembly warning: %v\n", err)
		fragStructures = []*parser.Protein{}
	}

	fmt.Printf("✅ Generated %d structures via fragment assembly\n", len(fragStructures))
	result.FragmentStructures = len(fragStructures)

	// Optimize and evaluate each fragment structure
	for i, structure := range fragStructures {
		optimized, metrics := optimizeAndEvaluate(structure, experimental, fmt.Sprintf("fragment_%d", i))
		if metrics != nil {
			metrics.ID = len(result.Structures)
			metrics.SamplingMethod = "Fragment Assembly"
			result.Structures = append(result.Structures, *metrics)
		}
		_ = optimized
	}
	fmt.Println()

	// ==================== AGENT 2.4: RAMACHANDRAN BASIN EXPLORER ====================
	fmt.Println("╔══════════════════════════════════════════════════════════════════╗")
	fmt.Println("║  AGENT 2.4: RAMACHANDRAN BASIN EXPLORER                          ║")
	fmt.Println("║  Systematic exploration of α-helix, β-sheet, PPII basins        ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════════╝")
	fmt.Println()

	basinConfig := sampling.DefaultBasinExplorerConfig()
	basinConfig.SamplesPerBasin = 5 // 5 samples per basin
	basinConfig.UseVedicBiasing = true

	// Mixed basin sampling for diversity
	basinStructures, err := sampling.MixedBasinSampling(sequence, basinConfig, 40)
	if err != nil {
		fmt.Printf("⚠️  Basin sampling warning: %v\n", err)
		basinStructures = []*parser.Protein{}
	}

	fmt.Printf("✅ Generated %d structures via basin exploration\n", len(basinStructures))
	result.BasinStructures = len(basinStructures)

	// Optimize and evaluate each basin structure
	for i, structure := range basinStructures {
		optimized, metrics := optimizeAndEvaluate(structure, experimental, fmt.Sprintf("basin_%d", i))
		if metrics != nil {
			metrics.ID = len(result.Structures)
			metrics.SamplingMethod = "Basin Explorer"
			result.Structures = append(result.Structures, *metrics)
		}
		_ = optimized
	}
	fmt.Println()

	// ==================== PHASE 2 ANALYSIS ====================
	fmt.Println("╔══════════════════════════════════════════════════════════════════╗")
	fmt.Println("║  PHASE 2 ANALYSIS & VALIDATION                                   ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════════╝")
	fmt.Println()

	result.TotalStructures = len(result.Structures)
	result.TotalDuration = time.Since(startTime).Seconds()

	if result.TotalStructures == 0 {
		fmt.Println("❌ ERROR: No structures generated!")
		return
	}

	// Calculate RMSD statistics
	if experimental != nil {
		rmsds := make([]float64, len(result.Structures))
		energies := make([]float64, len(result.Structures))
		vedics := make([]float64, len(result.Structures))

		for i, s := range result.Structures {
			rmsds[i] = s.RMSD
			energies[i] = s.Energy
			vedics[i] = s.VedicScore
		}

		sort.Float64s(rmsds)
		sort.Float64s(energies)
		sort.Float64s(vedics)

		result.BestRMSD = rmsds[0]
		result.WorstRMSD = rmsds[len(rmsds)-1]
		result.MedianRMSD = median(rmsds)
		result.MeanRMSD = mean(rmsds)
		result.RMSDStdDev = stdDev(rmsds)
		result.RMSDImprovement = (26.45 - result.BestRMSD) / 26.45 * 100 // vs Phase 1

		result.BestEnergy = energies[0]
		result.WorstEnergy = energies[len(energies)-1]
		result.MedianEnergy = median(energies)
		result.MeanEnergy = mean(energies)

		result.BestVedic = vedics[len(vedics)-1]
		result.MedianVedic = median(vedics)
		result.MeanVedic = mean(vedics)

		// Find best structure for TM-score and GDT_TS
		bestIdx := 0
		for i, s := range result.Structures {
			if s.RMSD < result.Structures[bestIdx].RMSD {
				bestIdx = i
			}
		}
		result.BestTMScore = result.Structures[bestIdx].TMScore
		result.BestGDT_TS = result.Structures[bestIdx].GDT_TS

		// Find best sampling method
		methodRMSDs := make(map[string][]float64)
		for _, s := range result.Structures {
			methodRMSDs[s.SamplingMethod] = append(methodRMSDs[s.SamplingMethod], s.RMSD)
		}

		bestMethodRMSD := math.Inf(1)
		for method, methodRmsds := range methodRMSDs {
			avgRMSD := mean(methodRmsds)
			if avgRMSD < bestMethodRMSD {
				bestMethodRMSD = avgRMSD
				result.BestMethod = method
				result.BestMethodRMSD = avgRMSD
			}
		}
	}

	// Calculate quality score
	result.QualityScore = calculateQualityScore(result)
	result.QualityTier = getQualityTier(result.QualityScore)
	result.MissionAccomplished = result.BestRMSD < 15.0 && result.QualityScore >= 0.92

	// Print results
	printPhase2Results(result)

	// Save results to JSON
	saveResults(result)

	// Generate markdown report
	generateReport(result)

	fmt.Println()
	fmt.Println("╔══════════════════════════════════════════════════════════════════╗")
	if result.MissionAccomplished {
		fmt.Println("║  🎉 MISSION ACCOMPLISHED! PHASE 2 SUCCESS!                      ║")
	} else {
		fmt.Println("║  ⚠️  MISSION INCOMPLETE - Further optimization needed           ║")
	}
	fmt.Println("╚══════════════════════════════════════════════════════════════════╝")
}

// optimizeAndEvaluate applies gentle relaxation and evaluates structure
func optimizeAndEvaluate(structure *parser.Protein, experimental *parser.Protein, name string) (*parser.Protein, *StructureMetric) {
	if structure == nil {
		return nil, nil
	}

	// Apply gentle relaxation (Phase 1 winner!)
	// Note: GentleRelax modifies protein in-place
	relaxConfig := optimization.DefaultGentleRelaxationConfig()
	relaxConfig.MaxSteps = 500 // Increased budget for Phase 2
	relaxConfig.StepSize = 0.01
	relaxConfig.VdWCutoff = 10.0
	relaxConfig.ElecCutoff = 12.0

	relaxResult, err := optimization.GentleRelax(structure, relaxConfig)
	if err != nil {
		return structure, nil
	}

	// The protein is modified in-place, so structure is now optimized
	optimized := structure

	// Calculate metrics
	metric := &StructureMetric{
		Energy:            relaxResult.FinalEnergy,
		OptimizationSteps: relaxResult.Steps,
	}

	// Calculate Vedic score
	angles := geometry.CalculateRamachandran(optimized)
	if len(angles) > 0 {
		vedicResult := vedic.CalculateVedicScore(optimized, angles)
		metric.VedicScore = vedicResult.TotalScore
	}

	// Calculate RMSD if experimental structure available
	if experimental != nil {
		rmsd, _ := validation.CalculateRMSD(optimized, experimental)
		metric.RMSD = rmsd

		tmScore := validation.CalculateTMScore(optimized, experimental, len(optimized.Residues))
		metric.TMScore = tmScore

		gdtTS := validation.CalculateGDT_TS(optimized, experimental)
		metric.GDT_TS = gdtTS
	}

	return optimized, metric
}

// calculateQualityScore computes overall Phase 2 quality score
func calculateQualityScore(result *Phase2Result) float64 {
	score := 0.0

	// Component 1: RMSD improvement (0-0.4 points)
	if result.BestRMSD < 15.0 {
		rmsdScore := 0.4 * (1.0 - result.BestRMSD/15.0)
		score += rmsdScore
	}

	// Component 2: Structure diversity (0-0.2 points)
	if result.TotalStructures >= 100 {
		score += 0.2
	} else {
		score += 0.2 * float64(result.TotalStructures) / 100.0
	}

	// Component 3: TM-score (0-0.2 points)
	if result.BestTMScore > 0.5 {
		score += 0.2 * result.BestTMScore
	}

	// Component 4: Vedic score (0-0.1 points)
	score += 0.1 * result.BestVedic

	// Component 5: Multi-method success (0-0.1 points)
	activeMethods := 0
	if result.FibonacciStructures > 0 {
		activeMethods++
	}
	if result.MonteCarloStructures > 0 {
		activeMethods++
	}
	if result.FragmentStructures > 0 {
		activeMethods++
	}
	if result.BasinStructures > 0 {
		activeMethods++
	}
	score += 0.1 * float64(activeMethods) / 4.0

	return score
}

// getQualityTier returns quality tier based on score
func getQualityTier(score float64) string {
	if score >= 0.95 {
		return "LEGENDARY++ (AlphaFold Competitor)"
	} else if score >= 0.92 {
		return "LEGENDARY (Mission Success)"
	} else if score >= 0.85 {
		return "EXCELLENT (Near Success)"
	} else if score >= 0.75 {
		return "GOOD (Significant Progress)"
	} else {
		return "DEVELOPING (Needs Work)"
	}
}

// Statistical helper functions
func mean(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}

func median(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	n := len(values)
	if n%2 == 0 {
		return (values[n/2-1] + values[n/2]) / 2.0
	}
	return values[n/2]
}

func stdDev(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	m := mean(values)
	sumSq := 0.0
	for _, v := range values {
		diff := v - m
		sumSq += diff * diff
	}
	return math.Sqrt(sumSq / float64(len(values)))
}

// printPhase2Results prints comprehensive results
func printPhase2Results(result *Phase2Result) {
	fmt.Println("📊 PHASE 2 RESULTS:")
	fmt.Println("═══════════════════════════════════════════════════════════════════")
	fmt.Println()

	fmt.Printf("🧬 Protein: %s\n", result.ProteinName)
	fmt.Printf("📝 Sequence: %s (%d residues)\n", result.Sequence, result.NumResidues)
	fmt.Printf("⏱️  Total Duration: %.2f seconds\n", result.TotalDuration)
	fmt.Println()

	fmt.Println("📈 SAMPLING STATISTICS:")
	fmt.Printf("  Total Structures Generated: %d\n", result.TotalStructures)
	fmt.Printf("  - Fibonacci Sphere: %d\n", result.FibonacciStructures)
	fmt.Printf("  - Vedic Monte Carlo: %d\n", result.MonteCarloStructures)
	fmt.Printf("  - Fragment Assembly: %d\n", result.FragmentStructures)
	fmt.Printf("  - Basin Explorer: %d\n", result.BasinStructures)
	fmt.Println()

	if result.BestRMSD > 0 {
		fmt.Println("🎯 RMSD STATISTICS (vs Experimental):")
		fmt.Printf("  Best RMSD: %.2f Å ⭐\n", result.BestRMSD)
		fmt.Printf("  Median RMSD: %.2f Å\n", result.MedianRMSD)
		fmt.Printf("  Mean RMSD: %.2f Å\n", result.MeanRMSD)
		fmt.Printf("  Worst RMSD: %.2f Å\n", result.WorstRMSD)
		fmt.Printf("  Std Dev: %.2f Å\n", result.RMSDStdDev)
		fmt.Printf("  Improvement vs Phase 1: %.1f%%\n", result.RMSDImprovement)
		fmt.Println()
	}

	fmt.Println("⚡ ENERGY STATISTICS:")
	fmt.Printf("  Best Energy: %.2f kcal/mol\n", result.BestEnergy)
	fmt.Printf("  Median Energy: %.2f kcal/mol\n", result.MedianEnergy)
	fmt.Printf("  Mean Energy: %.2f kcal/mol\n", result.MeanEnergy)
	fmt.Println()

	fmt.Println("🔮 VEDIC STATISTICS:")
	fmt.Printf("  Best Vedic Score: %.3f\n", result.BestVedic)
	fmt.Printf("  Median Vedic Score: %.3f\n", result.MedianVedic)
	fmt.Printf("  Mean Vedic Score: %.3f\n", result.MeanVedic)
	fmt.Println()

	if result.BestTMScore > 0 {
		fmt.Println("📐 VALIDATION METRICS (Best Structure):")
		fmt.Printf("  TM-score: %.3f\n", result.BestTMScore)
		fmt.Printf("  GDT_TS: %.3f\n", result.BestGDT_TS)
		fmt.Println()
	}

	fmt.Println("🏆 QUALITY ASSESSMENT:")
	fmt.Printf("  Quality Score: %.3f\n", result.QualityScore)
	fmt.Printf("  Quality Tier: %s\n", result.QualityTier)
	fmt.Printf("  Mission Accomplished: %v\n", result.MissionAccomplished)
	fmt.Println()

	fmt.Println("🥇 BEST SAMPLING METHOD:")
	fmt.Printf("  Method: %s\n", result.BestMethod)
	fmt.Printf("  Average RMSD: %.2f Å\n", result.BestMethodRMSD)
	fmt.Println()

	fmt.Println("═══════════════════════════════════════════════════════════════════")
}

// saveResults saves results to JSON file
func saveResults(result *Phase2Result) {
	filename := "PHASE_2_RESULTS.json"
	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("⚠️  Warning: Could not save results: %v\n", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(result); err != nil {
		fmt.Printf("⚠️  Warning: Could not encode results: %v\n", err)
		return
	}

	fmt.Printf("💾 Results saved to %s\n", filename)
}

// generateReport generates comprehensive markdown report
func generateReport(result *Phase2Result) {
	filename := "PHASE_2_REPORT.md"
	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("⚠️  Warning: Could not create report: %v\n", err)
		return
	}
	defer file.Close()

	fmt.Fprintf(file, "# FoldVedic Phase 2 Report\n\n")
	fmt.Fprintf(file, "**Advanced Sampling → <15 Å RMSD**\n\n")
	fmt.Fprintf(file, "**Date:** %s\n\n", result.Timestamp.Format("2006-01-02 15:04:05"))

	fmt.Fprintf(file, "## 🎯 Mission Statement\n\n")
	fmt.Fprintf(file, "Generate 100+ diverse protein structures using 4 advanced sampling methods to achieve:\n")
	fmt.Fprintf(file, "- **Target RMSD:** <15 Å (2× improvement from Phase 1: 26.45 Å)\n")
	fmt.Fprintf(file, "- **Target Quality:** ≥0.92 (LEGENDARY tier)\n")
	fmt.Fprintf(file, "- **Target Structures:** 100+ diverse conformations\n\n")

	fmt.Fprintf(file, "## 📊 Results Summary\n\n")
	fmt.Fprintf(file, "| Metric | Value |\n")
	fmt.Fprintf(file, "|--------|-------|\n")
	fmt.Fprintf(file, "| Protein | %s |\n", result.ProteinName)
	fmt.Fprintf(file, "| Residues | %d |\n", result.NumResidues)
	fmt.Fprintf(file, "| Total Structures | **%d** |\n", result.TotalStructures)
	fmt.Fprintf(file, "| Best RMSD | **%.2f Å** |\n", result.BestRMSD)
	fmt.Fprintf(file, "| Median RMSD | %.2f Å |\n", result.MedianRMSD)
	fmt.Fprintf(file, "| RMSD Improvement | **%.1f%%** vs Phase 1 |\n", result.RMSDImprovement)
	fmt.Fprintf(file, "| Quality Score | **%.3f** |\n", result.QualityScore)
	fmt.Fprintf(file, "| Quality Tier | **%s** |\n", result.QualityTier)
	fmt.Fprintf(file, "| Mission Success | **%v** |\n", result.MissionAccomplished)
	fmt.Fprintf(file, "| Duration | %.2f seconds |\n\n", result.TotalDuration)

	fmt.Fprintf(file, "## 🧪 The 4 Sampling Agents\n\n")
	fmt.Fprintf(file, "### Agent 2.1: Fibonacci Sphere Sampling\n")
	fmt.Fprintf(file, "- **Method:** Golden angle (137.5°) uniform sampling on S³ hypersphere\n")
	fmt.Fprintf(file, "- **Structures Generated:** %d\n\n", result.FibonacciStructures)

	fmt.Fprintf(file, "### Agent 2.2: Vedic Monte Carlo\n")
	fmt.Fprintf(file, "- **Method:** Metropolis-Hastings with digital root biasing\n")
	fmt.Fprintf(file, "- **Structures Generated:** %d\n\n", result.MonteCarloStructures)

	fmt.Fprintf(file, "### Agent 2.3: Fragment Assembly\n")
	fmt.Fprintf(file, "- **Method:** Rosetta-style assembly with Vedic fragment ranking\n")
	fmt.Fprintf(file, "- **Structures Generated:** %d\n\n", result.FragmentStructures)

	fmt.Fprintf(file, "### Agent 2.4: Ramachandran Basin Explorer\n")
	fmt.Fprintf(file, "- **Method:** Systematic exploration of α-helix, β-sheet, PPII basins\n")
	fmt.Fprintf(file, "- **Structures Generated:** %d\n\n", result.BasinStructures)

	fmt.Fprintf(file, "## 📈 RMSD Distribution\n\n")
	fmt.Fprintf(file, "```\n")
	fmt.Fprintf(file, "Best:    %.2f Å  ★\n", result.BestRMSD)
	fmt.Fprintf(file, "Median:  %.2f Å\n", result.MedianRMSD)
	fmt.Fprintf(file, "Mean:    %.2f Å\n", result.MeanRMSD)
	fmt.Fprintf(file, "Worst:   %.2f Å\n", result.WorstRMSD)
	fmt.Fprintf(file, "Std Dev: %.2f Å\n", result.RMSDStdDev)
	fmt.Fprintf(file, "```\n\n")

	fmt.Fprintf(file, "## 🏆 Best Sampling Method\n\n")
	fmt.Fprintf(file, "**Winner:** %s (Average RMSD: %.2f Å)\n\n", result.BestMethod, result.BestMethodRMSD)

	fmt.Fprintf(file, "## 🎓 Phase 2 Insights\n\n")
	if result.MissionAccomplished {
		fmt.Fprintf(file, "✅ **PHASE 2 SUCCESS!**\n\n")
		fmt.Fprintf(file, "We have successfully achieved:\n")
		fmt.Fprintf(file, "- RMSD <15 Å (Target: %.2f Å achieved!)\n", result.BestRMSD)
		fmt.Fprintf(file, "- Quality ≥0.92 (Achieved: %.3f)\n", result.QualityScore)
		fmt.Fprintf(file, "- 100+ diverse structures (Generated: %d)\n\n", result.TotalStructures)
	} else {
		fmt.Fprintf(file, "⚠️ **PHASE 2 INCOMPLETE**\n\n")
		fmt.Fprintf(file, "Further optimization needed:\n")
		if result.BestRMSD >= 15.0 {
			fmt.Fprintf(file, "- RMSD still %.2f Å (Target: <15 Å)\n", result.BestRMSD)
		}
		if result.QualityScore < 0.92 {
			fmt.Fprintf(file, "- Quality %.3f (Target: ≥0.92)\n", result.QualityScore)
		}
		if result.TotalStructures < 100 {
			fmt.Fprintf(file, "- Only %d structures (Target: 100+)\n", result.TotalStructures)
		}
		fmt.Fprintf(file, "\n")
	}

	fmt.Fprintf(file, "## 🔮 Philosophy & Innovation\n\n")
	fmt.Fprintf(file, "Phase 2 demonstrates the power of:\n\n")
	fmt.Fprintf(file, "1. **Wright Brothers Empiricism:** Gentle relaxation (Phase 1 winner) applied to all structures\n")
	fmt.Fprintf(file, "2. **Quaternion-First Thinking:** Fibonacci sphere on S³ for uniform coverage\n")
	fmt.Fprintf(file, "3. **Cross-Domain Fearlessness:** Computer graphics + Vedic mathematics + statistical mechanics\n")
	fmt.Fprintf(file, "4. **Multi-Method Exploration:** 4 diverse sampling strategies → comprehensive conformational search\n\n")

	fmt.Fprintf(file, "## 📝 Next Steps: Phase 3\n\n")
	fmt.Fprintf(file, "- Fix L-BFGS explosion via quaternion parameterization\n")
	fmt.Fprintf(file, "- Implement simulated annealing (robust alternative)\n")
	fmt.Fprintf(file, "- Increase minimization budget (1000-5000 steps)\n")
	fmt.Fprintf(file, "- Add structural constraints (contacts, secondary structure, membrane)\n")
	fmt.Fprintf(file, "- **Target:** <10 Å RMSD (Rosetta-competitive)\n\n")

	fmt.Fprintf(file, "---\n\n")
	fmt.Fprintf(file, "*Report generated by FoldVedic Phase 2 Integration Pipeline*\n")
	fmt.Fprintf(file, "*Built with full agency by Autonomous AI*\n")
	fmt.Fprintf(file, "*May this work benefit all of humanity* 🌍\n")

	fmt.Printf("📄 Report saved to %s\n", filename)
}
