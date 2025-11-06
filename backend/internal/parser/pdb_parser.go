// Package parser provides PDB file parsing functionality for protein structures.
//
// BIOCHEMIST: Parses PDB format to extract backbone atoms (N, Cα, C) for Ramachandran analysis
// PHYSICIST: Handles 3D coordinates with precision for accurate force calculations
// MATHEMATICIAN: Provides clean data structures for quaternion mapping
// ETHICIST: Robust error handling for real-world PDB files (missing atoms, alternate conformations)
package parser

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Atom represents a single atom in 3D space
type Atom struct {
	Serial    int     // Atom serial number
	Name      string  // Atom name (e.g., "CA", "N", "C")
	AltLoc    string  // Alternate location indicator
	ResName   string  // Residue name (e.g., "ALA", "GLY")
	ChainID   string  // Chain identifier
	ResSeq    int     // Residue sequence number
	ICode     string  // Insertion code
	X, Y, Z   float64 // Atomic coordinates (Angstroms)
	Occupancy float64 // Occupancy
	TempFacto float64 // Temperature factor
	Element   string  // Element symbol
}

// Residue represents an amino acid residue with its backbone atoms
type Residue struct {
	Name    string  // Three-letter code (ALA, GLY, etc.)
	SeqNum  int     // Sequence number
	ChainID string  // Chain identifier
	N       *Atom   // Nitrogen (backbone)
	CA      *Atom   // Alpha carbon (backbone)
	C       *Atom   // Carbonyl carbon (backbone)
	O       *Atom   // Carbonyl oxygen (backbone)
}

// Protein represents a complete protein structure
type Protein struct {
	Name     string     // Protein name/PDB ID
	Residues []*Residue // All residues in sequence
	Atoms    []*Atom    // All atoms
}

// ParsePDB parses a PDB file and extracts protein structure
//
// Citation: PDB format specification from RCSB PDB (www.wwpdb.org)
// Handles ATOM and HETATM records, filters for protein backbone atoms
func ParsePDB(filename string) (*Protein, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open PDB file: %w", err)
	}
	defer file.Close()

	protein := &Protein{
		Name:     filename,
		Residues: make([]*Residue, 0),
		Atoms:    make([]*Atom, 0),
	}

	// Map to group atoms by residue (chainID:resSeq)
	residueMap := make(map[string]*Residue)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Parse ATOM and HETATM records
		if len(line) >= 6 && (line[0:4] == "ATOM" || line[0:6] == "HETATM") {
			atom, err := parseAtomLine(line)
			if err != nil {
				// Skip malformed lines but continue parsing
				continue
			}

			protein.Atoms = append(protein.Atoms, atom)

			// Only process backbone atoms for Ramachandran analysis
			if isBackboneAtom(atom.Name) {
				resKey := fmt.Sprintf("%s:%d", atom.ChainID, atom.ResSeq)

				// Get or create residue
				res, exists := residueMap[resKey]
				if !exists {
					res = &Residue{
						Name:    atom.ResName,
						SeqNum:  atom.ResSeq,
						ChainID: atom.ChainID,
					}
					residueMap[resKey] = res
					protein.Residues = append(protein.Residues, res)
				}

				// Assign atom to residue based on atom name
				switch atom.Name {
				case "N":
					res.N = atom
				case "CA":
					res.CA = atom
				case "C":
					res.C = atom
				case "O":
					res.O = atom
				}
			}
		}

		// Stop at END or ENDMDL
		if len(line) >= 3 && (line[0:3] == "END" || line[0:6] == "ENDMDL") {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading PDB file: %w", err)
	}

	return protein, nil
}

// parseAtomLine parses a single ATOM/HETATM line from PDB format
//
// PDB format (fixed-width columns):
// ATOM      1  N   ALA A   1      11.104   6.134  -6.504  1.00  0.00           N
// Cols: 1-6 (record), 7-11 (serial), 13-16 (name), 17 (altLoc), 18-20 (resName),
//       22 (chainID), 23-26 (resSeq), 31-38 (x), 39-46 (y), 47-54 (z), etc.
func parseAtomLine(line string) (*Atom, error) {
	// Ensure line is long enough
	if len(line) < 54 {
		return nil, fmt.Errorf("line too short: %d characters", len(line))
	}

	// Pad line to ensure we can safely access all columns
	for len(line) < 80 {
		line += " "
	}

	atom := &Atom{}

	// Serial number (columns 7-11)
	if serial, err := strconv.Atoi(strings.TrimSpace(line[6:11])); err == nil {
		atom.Serial = serial
	}

	// Atom name (columns 13-16)
	atom.Name = strings.TrimSpace(line[12:16])

	// Alternate location (column 17)
	atom.AltLoc = strings.TrimSpace(line[16:17])

	// Residue name (columns 18-20)
	atom.ResName = strings.TrimSpace(line[17:20])

	// Chain ID (column 22)
	atom.ChainID = strings.TrimSpace(line[21:22])

	// Residue sequence number (columns 23-26)
	if resSeq, err := strconv.Atoi(strings.TrimSpace(line[22:26])); err == nil {
		atom.ResSeq = resSeq
	}

	// Insertion code (column 27)
	atom.ICode = strings.TrimSpace(line[26:27])

	// Coordinates (columns 31-38, 39-46, 47-54)
	if x, err := strconv.ParseFloat(strings.TrimSpace(line[30:38]), 64); err == nil {
		atom.X = x
	}
	if y, err := strconv.ParseFloat(strings.TrimSpace(line[38:46]), 64); err == nil {
		atom.Y = y
	}
	if z, err := strconv.ParseFloat(strings.TrimSpace(line[46:54]), 64); err == nil {
		atom.Z = z
	}

	// Occupancy (columns 55-60)
	if len(line) >= 60 {
		if occ, err := strconv.ParseFloat(strings.TrimSpace(line[54:60]), 64); err == nil {
			atom.Occupancy = occ
		}
	}

	// Temperature factor (columns 61-66)
	if len(line) >= 66 {
		if temp, err := strconv.ParseFloat(strings.TrimSpace(line[60:66]), 64); err == nil {
			atom.TempFacto = temp
		}
	}

	// Element symbol (columns 77-78)
	if len(line) >= 78 {
		atom.Element = strings.TrimSpace(line[76:78])
	}

	return atom, nil
}

// isBackboneAtom checks if an atom name is a backbone atom
//
// BIOCHEMIST: Backbone atoms are N (nitrogen), CA (alpha carbon), C (carbonyl carbon), O (carbonyl oxygen)
// These define the protein backbone and are used for Ramachandran analysis
func isBackboneAtom(name string) bool {
	return name == "N" || name == "CA" || name == "C" || name == "O"
}

// HasCompleteBackbone checks if a residue has all required backbone atoms
//
// BIOCHEMIST: Complete backbone required for dihedral angle calculations
// Missing atoms can occur in real PDB files (crystal disorder, missing density)
func (r *Residue) HasCompleteBackbone() bool {
	return r.N != nil && r.CA != nil && r.C != nil
}

// NumCompleteResidues returns the count of residues with complete backbone
func (p *Protein) NumCompleteResidues() int {
	count := 0
	for _, res := range p.Residues {
		if res.HasCompleteBackbone() {
			count++
		}
	}
	return count
}
