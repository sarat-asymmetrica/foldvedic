// PDB Downloader - fetches protein structures from RCSB PDB
//
// Usage: go run main.go 1UBQ 1CRN 2KXA
//
// Downloads structures to testdata/ directory
package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	pdbBaseURL = "https://files.rcsb.org/download/"
	outputDir  = "../../../testdata/"
)

func main() {
	// Default structures for Wave 2 validation
	pdbIDs := []string{"1UBQ", "1CRN", "2KXA", "1VII", "1L2Y"}

	if len(os.Args) > 1 {
		pdbIDs = os.Args[1:]
	}

	fmt.Println("FoldVedic PDB Downloader")
	fmt.Println(strings.Repeat("=", 80))
	fmt.Printf("Downloading %d PDB structures...\n\n", len(pdbIDs))

	// Ensure output directory exists
	os.MkdirAll(outputDir, 0755)

	successCount := 0
	for _, pdbID := range pdbIDs {
		pdbID = strings.ToUpper(strings.TrimSpace(pdbID))

		if downloadPDB(pdbID) {
			successCount++
		}

		// Be nice to RCSB servers
		time.Sleep(500 * time.Millisecond)
	}

	fmt.Println()
	fmt.Println(strings.Repeat("=", 80))
	fmt.Printf("Download complete: %d/%d structures successfully downloaded\n", successCount, len(pdbIDs))

	if successCount < len(pdbIDs) {
		fmt.Println("\nNote: Some downloads may have failed due to network issues.")
		fmt.Println("      These structures will be created as test stubs.")
		os.Exit(1)
	}
}

func downloadPDB(pdbID string) bool {
	url := fmt.Sprintf("%s%s.pdb", pdbBaseURL, pdbID)
	outputFile := fmt.Sprintf("%s%s.pdb", outputDir, strings.ToLower(pdbID))

	fmt.Printf("Downloading %s... ", pdbID)

	// Check if already exists
	if _, err := os.Stat(outputFile); err == nil {
		fmt.Println("✓ Already exists")
		return true
	}

	// Download
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("✗ Failed: %v\n", err)
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Printf("✗ Failed: HTTP %d\n", resp.StatusCode)
		return false
	}

	// Save to file
	out, err := os.Create(outputFile)
	if err != nil {
		fmt.Printf("✗ Failed to create file: %v\n", err)
		return false
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		fmt.Printf("✗ Failed to write file: %v\n", err)
		return false
	}

	// Get file size
	stat, _ := os.Stat(outputFile)
	size := stat.Size() / 1024 // KB

	fmt.Printf("✓ Downloaded (%d KB)\n", size)
	return true
}
