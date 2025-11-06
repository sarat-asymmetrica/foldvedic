// Benchmark validation framework for FoldVedic.ai
//
// WAVE 6: Large-scale validation against experimental structures
// Tests prediction accuracy across diverse protein folds
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/sarat-asymmetrica/foldvedic/backend/internal/folding"
	"github.com/sarat-asymmetrica/foldvedic/backend/internal/parser"
)

// BenchmarkProtein represents a test case with metadata
type BenchmarkProtein struct {
	PDBCode     string  `json:"pdb_code"`
	Name        string  `json:"name"`
	Length      int     `json:"length"`
	Description string  `json:"description"`
	FoldClass   string  `json:"fold_class"` // alpha, beta, alpha+beta, irregular
}

// BenchmarkResult stores prediction results for one protein
type BenchmarkResult struct {
	PDBCode     string    `json:"pdb_code"`
	Name        string    `json:"name"`
	Length      int       `json:"length"`
	FoldClass   string    `json:"fold_class"`

	// Validation metrics
	RMSD        float64   `json:"rmsd"`
	TMScore     float64   `json:"tm_score"`
	GDT_TS      float64   `json:"gdt_ts"`
	VedicScore  float64   `json:"vedic_score"`
	QualityScore float64  `json:"quality_score"`

	// Performance metrics
	TimeElapsed float64   `json:"time_elapsed_seconds"`
	Converged   bool      `json:"converged"`
	NumSteps    int       `json:"num_steps"`

	// Error tracking
	Success     bool      `json:"success"`
	ErrorMsg    string    `json:"error_msg,omitempty"`
}

// BenchmarkSummary holds aggregate statistics
type BenchmarkSummary struct {
	TotalProteins    int       `json:"total_proteins"`
	SuccessfulPreds  int       `json:"successful_predictions"`
	FailedPreds      int       `json:"failed_predictions"`

	// Mean metrics
	MeanRMSD         float64   `json:"mean_rmsd"`
	MedianRMSD       float64   `json:"median_rmsd"`
	MeanTMScore      float64   `json:"mean_tm_score"`
	MedianTMScore    float64   `json:"median_tm_score"`
	MeanGDT_TS       float64   `json:"mean_gdt_ts"`
	MeanQuality      float64   `json:"mean_quality_score"`

	// Performance
	TotalTime        float64   `json:"total_time_seconds"`
	MeanTime         float64   `json:"mean_time_seconds"`

	// Quality thresholds
	ExcellentPreds   int       `json:"excellent_predictions"` // RMSD < 2Å, TM > 0.6
	GoodPreds        int       `json:"good_predictions"`      // RMSD < 3.5Å, TM > 0.5
	AcceptablePreds  int       `json:"acceptable_predictions"` // RMSD < 5Å

	Results          []BenchmarkResult `json:"results"`
}

// Curated benchmark set covering diverse fold classes
// Selection criteria:
// - High-resolution structures (< 2.0 Å)
// - No missing residues
// - Diverse SCOP classifications
// - Range of sizes (20-200 residues)
var benchmarkSet = []BenchmarkProtein{
	// ALPHA-HELICAL PROTEINS
	{PDBCode: "1UBQ", Name: "Ubiquitin", Length: 76, FoldClass: "alpha+beta", Description: "Small regulatory protein"},
	{PDBCode: "1VII", Name: "Villin headpiece", Length: 36, FoldClass: "alpha", Description: "Ultra-fast folder"},
	{PDBCode: "2LVG", Name: "Myoglobin", Length: 154, FoldClass: "alpha", Description: "Classic globin fold"},
	{PDBCode: "1MBN", Name: "Myoglobin mini", Length: 153, FoldClass: "alpha", Description: "Oxygen carrier"},
	{PDBCode: "256B", Name: "Cytochrome b562", Length: 106, FoldClass: "alpha", Description: "Four-helix bundle"},

	// BETA PROTEINS
	{PDBCode: "1CRN", Name: "Crambin", Length: 46, FoldClass: "beta", Description: "Plant seed protein"},
	{PDBCode: "1PIN", Name: "Protein G B1", Length: 56, FoldClass: "alpha+beta", Description: "Immunoglobulin binding"},
	{PDBCode: "1SHG", Name: "SH3 domain", Length: 62, FoldClass: "beta", Description: "Beta barrel"},
	{PDBCode: "1TEN", Name: "Tenascin", Length: 90, FoldClass: "beta", Description: "Fibronectin type III"},

	// ALPHA+BETA PROTEINS
	{PDBCode: "1L2Y", Name: "Trp-cage miniprotein", Length: 20, FoldClass: "alpha", Description: "Smallest natural protein"},
	{PDBCode: "2KXA", Name: "WW domain", Length: 34, FoldClass: "beta", Description: "Fast-folding triple-stranded sheet"},
	{PDBCode: "1RIS", Name: "Ras protein", Length: 166, FoldClass: "alpha+beta", Description: "GTPase"},
	{PDBCode: "1ENH", Name: "Engrailed homeodomain", Length: 54, FoldClass: "alpha", Description: "DNA-binding helix-turn-helix"},
	{PDBCode: "1PGB", Name: "Protein G", Length: 56, FoldClass: "alpha+beta", Description: "Four-stranded beta sheet"},

	// IRREGULAR / DISORDERED
	{PDBCode: "1BDD", Name: "Hirudin", Length: 65, FoldClass: "irregular", Description: "Thrombin inhibitor"},

	// Additional diverse structures
	{PDBCode: "1YRF", Name: "BBA5", Length: 46, FoldClass: "alpha", Description: "Three-helix bundle"},
	{PDBCode: "1PRB", Name: "Protein B", Length: 53, FoldClass: "alpha", Description: "Four-helix bundle"},
	{PDBCode: "1IGD", Name: "Immunoglobulin", Length: 61, FoldClass: "beta", Description: "Beta sandwich"},
	{PDBCode: "1UTG", Name: "Uteroglobin", Length: 70, FoldClass: "alpha", Description: "Anti-inflammatory protein"},
	{PDBCode: "2PTN", Name: "Trypsin inhibitor", Length: 58, FoldClass: "alpha+beta", Description: "Disulfide-rich"},
}

func main() {
	fmt.Println("=== FoldVedic.ai Wave 6: Large-scale Benchmark Validation ===\n")

	// Create data directory
	dataDir := "testdata/benchmark"
	os.MkdirAll(dataDir, 0755)

	// Download benchmark structures (with progress)
	fmt.Printf("Downloading %d benchmark structures...\n", len(benchmarkSet))
	downloadBenchmarkSet(dataDir)

	// Run predictions in parallel
	fmt.Println("\nRunning predictions on benchmark set...")
	results := runBenchmark(dataDir)

	// Calculate statistics
	fmt.Println("\nCalculating statistics...")
	summary := calculateSummary(results)

	// Generate report
	fmt.Println("\nGenerating validation report...")
	generateReport(summary)

	// Save JSON results
	saveResults(summary, "WAVE_6_BENCHMARK_RESULTS.json")

	fmt.Println("\n=== Wave 6 Benchmark Validation Complete ===")
	fmt.Printf("Success rate: %.1f%% (%d/%d)\n",
		float64(summary.SuccessfulPreds)/float64(summary.TotalProteins)*100,
		summary.SuccessfulPreds, summary.TotalProteins)
	fmt.Printf("Mean RMSD: %.2f Å\n", summary.MeanRMSD)
	fmt.Printf("Mean TM-score: %.3f\n", summary.MeanTMScore)
	fmt.Printf("Quality score: %.3f\n", summary.MeanQuality)
}

func downloadBenchmarkSet(dataDir string) {
	var wg sync.WaitGroup
	sem := make(chan struct{}, 5) // Limit to 5 concurrent downloads

	for i, protein := range benchmarkSet {
		wg.Add(1)
		go func(idx int, prot BenchmarkProtein) {
			defer wg.Done()
			sem <- struct{}{} // Acquire semaphore
			defer func() { <-sem }() // Release semaphore

			filename := filepath.Join(dataDir, prot.PDBCode+".pdb")

			// Skip if already downloaded
			if _, err := os.Stat(filename); err == nil {
				fmt.Printf("[%d/%d] %s already downloaded\n", idx+1, len(benchmarkSet), prot.PDBCode)
				return
			}

			url := fmt.Sprintf("https://files.rcsb.org/download/%s.pdb", prot.PDBCode)

			// Download with retries
			var resp *http.Response
			var err error
			for retry := 0; retry < 3; retry++ {
				resp, err = http.Get(url)
				if err == nil && resp.StatusCode == 200 {
					break
				}
				time.Sleep(time.Second * time.Duration(1<<retry)) // Exponential backoff
			}

			if err != nil || resp.StatusCode != 200 {
				fmt.Printf("[%d/%d] Failed to download %s: %v\n", idx+1, len(benchmarkSet), prot.PDBCode, err)
				return
			}
			defer resp.Body.Close()

			// Save to file
			out, err := os.Create(filename)
			if err != nil {
				fmt.Printf("[%d/%d] Failed to create %s: %v\n", idx+1, len(benchmarkSet), prot.PDBCode, err)
				return
			}
			defer out.Close()

			_, err = io.Copy(out, resp.Body)
			if err != nil {
				fmt.Printf("[%d/%d] Failed to save %s: %v\n", idx+1, len(benchmarkSet), prot.PDBCode, err)
				return
			}

			fmt.Printf("[%d/%d] Downloaded %s (%s)\n", idx+1, len(benchmarkSet), prot.PDBCode, prot.Name)
			time.Sleep(500 * time.Millisecond) // Polite rate limiting
		}(i, protein)
	}

	wg.Wait()
	fmt.Println("Download complete!")
}

func runBenchmark(dataDir string) []BenchmarkResult {
	results := make([]BenchmarkResult, 0, len(benchmarkSet))
	var mu sync.Mutex
	var wg sync.WaitGroup
	sem := make(chan struct{}, 4) // Run 4 predictions in parallel

	for i, protein := range benchmarkSet {
		wg.Add(1)
		go func(idx int, prot BenchmarkProtein) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			result := runSinglePrediction(dataDir, prot, idx+1, len(benchmarkSet))

			mu.Lock()
			results = append(results, result)
			mu.Unlock()
		}(i, protein)
	}

	wg.Wait()
	return results
}

func runSinglePrediction(dataDir string, prot BenchmarkProtein, idx, total int) BenchmarkResult {
	result := BenchmarkResult{
		PDBCode:   prot.PDBCode,
		Name:      prot.Name,
		Length:    prot.Length,
		FoldClass: prot.FoldClass,
	}

	startTime := time.Now()

	// Load experimental structure
	pdbFile := filepath.Join(dataDir, prot.PDBCode+".pdb")
	experimental, err := parser.ParsePDB(pdbFile)
	if err != nil {
		result.Success = false
		result.ErrorMsg = fmt.Sprintf("Failed to parse PDB: %v", err)
		fmt.Printf("[%d/%d] %s FAILED (parse error)\n", idx, total, prot.PDBCode)
		return result
	}

	// Extract sequence
	sequence := extractSequence(experimental)
	if len(sequence) == 0 {
		result.Success = false
		result.ErrorMsg = "Empty sequence"
		fmt.Printf("[%d/%d] %s FAILED (empty sequence)\n", idx, total, prot.PDBCode)
		return result
	}

	// Run prediction
	config := folding.DefaultPredictionConfig(sequence)
	config.NumSamples = 5 // Use 5 samples for faster benchmarking
	config.MinimizerConfig.MaxSteps = 100 // Limit iterations

	predResult, err := folding.PredictStructure(config, experimental)
	if err != nil {
		result.Success = false
		result.ErrorMsg = fmt.Sprintf("Prediction failed: %v", err)
		fmt.Printf("[%d/%d] %s FAILED (prediction error)\n", idx, total, prot.PDBCode)
		return result
	}

	elapsed := time.Since(startTime).Seconds()

	// Extract metrics
	result.Success = true
	result.TimeElapsed = elapsed
	result.Converged = predResult.Converged
	result.NumSteps = predResult.NumSteps
	result.VedicScore = predResult.VedicScore.TotalScore
	result.QualityScore = predResult.QualityScore

	if predResult.Comparison != nil {
		result.RMSD = predResult.Comparison.RMSD
		result.TMScore = predResult.Comparison.TMScore
		result.GDT_TS = predResult.Comparison.GDT_TS
	}

	// Quality assessment
	quality := "POOR"
	if result.RMSD < 2.0 && result.TMScore > 0.6 {
		quality = "EXCELLENT"
	} else if result.RMSD < 3.5 && result.TMScore > 0.5 {
		quality = "GOOD"
	} else if result.RMSD < 5.0 {
		quality = "ACCEPTABLE"
	}

	fmt.Printf("[%d/%d] %s: RMSD=%.2fÅ TM=%.3f Quality=%s (%.1fs)\n",
		idx, total, prot.PDBCode, result.RMSD, result.TMScore, quality, elapsed)

	return result
}

func extractSequence(protein *parser.Protein) string {
	sequence := ""
	for _, res := range protein.Residues {
		// Convert three-letter to one-letter code (simplified)
		aa := threeToOne(res.Name)
		if aa != "" {
			sequence += aa
		}
	}
	return sequence
}

func threeToOne(three string) string {
	codes := map[string]string{
		"ALA": "A", "ARG": "R", "ASN": "N", "ASP": "D", "CYS": "C",
		"GLN": "Q", "GLU": "E", "GLY": "G", "HIS": "H", "ILE": "I",
		"LEU": "L", "LYS": "K", "MET": "M", "PHE": "F", "PRO": "P",
		"SER": "S", "THR": "T", "TRP": "W", "TYR": "Y", "VAL": "V",
	}
	return codes[three]
}

func calculateSummary(results []BenchmarkResult) BenchmarkSummary {
	summary := BenchmarkSummary{
		TotalProteins: len(results),
		Results:      results,
	}

	// Separate successful predictions
	var successResults []BenchmarkResult
	for _, r := range results {
		if r.Success {
			summary.SuccessfulPreds++
			successResults = append(successResults, r)
			summary.TotalTime += r.TimeElapsed
		} else {
			summary.FailedPreds++
		}
	}

	if len(successResults) == 0 {
		return summary
	}

	// Calculate means
	sumRMSD := 0.0
	sumTM := 0.0
	sumGDT := 0.0
	sumQuality := 0.0

	for _, r := range successResults {
		sumRMSD += r.RMSD
		sumTM += r.TMScore
		sumGDT += r.GDT_TS
		sumQuality += r.QualityScore

		// Count by quality threshold
		if r.RMSD < 2.0 && r.TMScore > 0.6 {
			summary.ExcellentPreds++
		} else if r.RMSD < 3.5 && r.TMScore > 0.5 {
			summary.GoodPreds++
		} else if r.RMSD < 5.0 {
			summary.AcceptablePreds++
		}
	}

	n := float64(len(successResults))
	summary.MeanRMSD = sumRMSD / n
	summary.MeanTMScore = sumTM / n
	summary.MeanGDT_TS = sumGDT / n
	summary.MeanQuality = sumQuality / n
	summary.MeanTime = summary.TotalTime / n

	// Calculate medians (simplified - just sort and take middle)
	rmsdValues := make([]float64, len(successResults))
	tmValues := make([]float64, len(successResults))
	for i, r := range successResults {
		rmsdValues[i] = r.RMSD
		tmValues[i] = r.TMScore
	}
	summary.MedianRMSD = median(rmsdValues)
	summary.MedianTMScore = median(tmValues)

	return summary
}

func median(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}

	// Simple bubble sort for small arrays
	sorted := make([]float64, len(values))
	copy(sorted, values)
	for i := 0; i < len(sorted); i++ {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[i] > sorted[j] {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}

	mid := len(sorted) / 2
	if len(sorted)%2 == 0 {
		return (sorted[mid-1] + sorted[mid]) / 2
	}
	return sorted[mid]
}

func generateReport(summary BenchmarkSummary) {
	report := fmt.Sprintf(`# Wave 6 Benchmark Validation Report

## Summary Statistics

**Dataset:** %d diverse protein structures (20-200 residues)
**Success Rate:** %.1f%% (%d/%d predictions completed)
**Total Time:** %.1f seconds
**Mean Time per Protein:** %.2f seconds

## Accuracy Metrics

| Metric | Mean | Median | Interpretation |
|--------|------|--------|----------------|
| **RMSD** | %.2f Å | %.2f Å | %s |
| **TM-score** | %.3f | %.3f | %s |
| **GDT_TS** | %.3f | - | %s |
| **Quality Score** | %.3f | - | %s |

## Quality Distribution

- **Excellent** (RMSD < 2Å, TM > 0.6): %d (%.1f%%)
- **Good** (RMSD < 3.5Å, TM > 0.5): %d (%.1f%%)
- **Acceptable** (RMSD < 5Å): %d (%.1f%%)

## Breakdown by Fold Class

`,
		summary.TotalProteins,
		float64(summary.SuccessfulPreds)/float64(summary.TotalProteins)*100,
		summary.SuccessfulPreds, summary.TotalProteins,
		summary.TotalTime,
		summary.MeanTime,
		summary.MeanRMSD, summary.MedianRMSD, interpretRMSD(summary.MeanRMSD),
		summary.MeanTMScore, summary.MedianTMScore, interpretTM(summary.MeanTMScore),
		summary.MeanGDT_TS, interpretGDT(summary.MeanGDT_TS),
		summary.MeanQuality, interpretQuality(summary.MeanQuality),
		summary.ExcellentPreds, float64(summary.ExcellentPreds)/float64(summary.SuccessfulPreds)*100,
		summary.GoodPreds, float64(summary.GoodPreds)/float64(summary.SuccessfulPreds)*100,
		summary.AcceptablePreds, float64(summary.AcceptablePreds)/float64(summary.SuccessfulPreds)*100,
	)

	// Per-fold-class analysis
	foldClasses := map[string][]BenchmarkResult{
		"alpha": {}, "beta": {}, "alpha+beta": {}, "irregular": {},
	}
	for _, r := range summary.Results {
		if r.Success {
			foldClasses[r.FoldClass] = append(foldClasses[r.FoldClass], r)
		}
	}

	for class, results := range foldClasses {
		if len(results) == 0 {
			continue
		}
		meanRMSD := 0.0
		meanTM := 0.0
		for _, r := range results {
			meanRMSD += r.RMSD
			meanTM += r.TMScore
		}
		meanRMSD /= float64(len(results))
		meanTM /= float64(len(results))

		report += fmt.Sprintf("**%s:** %d proteins, RMSD=%.2fÅ, TM=%.3f\n",
			class, len(results), meanRMSD, meanTM)
	}

	report += "\n## Individual Results\n\n"
	report += "| PDB | Name | Length | RMSD (Å) | TM-score | Quality | Time (s) |\n"
	report += "|-----|------|--------|----------|----------|---------|----------|\n"

	for _, r := range summary.Results {
		if !r.Success {
			report += fmt.Sprintf("| %s | %s | %d | FAILED | - | - | - |\n",
				r.PDBCode, r.Name, r.Length)
		} else {
			quality := "Poor"
			if r.RMSD < 2.0 && r.TMScore > 0.6 {
				quality = "Excellent"
			} else if r.RMSD < 3.5 && r.TMScore > 0.5 {
				quality = "Good"
			} else if r.RMSD < 5.0 {
				quality = "Acceptable"
			}

			report += fmt.Sprintf("| %s | %s | %d | %.2f | %.3f | %s | %.1f |\n",
				r.PDBCode, r.Name, r.Length, r.RMSD, r.TMScore, quality, r.TimeElapsed)
		}
	}

	report += "\n## Methodology\n\n"
	report += "- **Algorithm:** FoldVedic.ai (Vedic mathematics + quaternion geometry + AMBER ff14SB)\n"
	report += "- **Conformational Sampling:** 5 samples per protein\n"
	report += "- **Energy Minimization:** Steepest descent (max 100 iterations)\n"
	report += "- **Validation Metrics:** RMSD (Kabsch alignment), TM-score, GDT_TS\n"
	report += "- **Hardware:** CPU-only (no GPU required)\n"
	report += "\n---\n**Generated:** " + time.Now().Format("2006-01-02 15:04:05") + "\n"

	// Write report
	os.WriteFile("WAVE_6_VALIDATION_REPORT.md", []byte(report), 0644)
	fmt.Println("Report saved to WAVE_6_VALIDATION_REPORT.md")
}

func interpretRMSD(rmsd float64) string {
	if rmsd < 2.0 {
		return "Excellent (near-identical)"
	} else if rmsd < 3.5 {
		return "Good (same fold)"
	} else if rmsd < 5.0 {
		return "Acceptable"
	}
	return "Poor"
}

func interpretTM(tm float64) string {
	if tm > 0.6 {
		return "High confidence same fold"
	} else if tm > 0.5 {
		return "Same fold"
	} else if tm > 0.3 {
		return "Similar topology"
	}
	return "Different folds"
}

func interpretGDT(gdt float64) string {
	if gdt > 0.7 {
		return "Excellent"
	} else if gdt > 0.5 {
		return "Good"
	} else if gdt > 0.3 {
		return "Acceptable"
	}
	return "Poor"
}

func interpretQuality(quality float64) string {
	if quality > 0.7 {
		return "Excellent"
	} else if quality > 0.5 {
		return "Good"
	} else if quality > 0.3 {
		return "Acceptable"
	}
	return "Poor"
}

func saveResults(summary BenchmarkSummary, filename string) {
	data, err := json.MarshalIndent(summary, "", "  ")
	if err != nil {
		fmt.Printf("Failed to marshal results: %v\n", err)
		return
	}

	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		fmt.Printf("Failed to write results: %v\n", err)
		return
	}

	fmt.Printf("Results saved to %s\n", filename)
}
