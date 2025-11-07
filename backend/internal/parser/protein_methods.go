package parser

// Copy creates a deep copy of a Protein structure
// This is needed by Wave 4 optimization agents
func (p *Protein) Copy() *Protein {
	if p == nil {
		return nil
	}

	clone := &Protein{
		Name:     p.Name,
		Residues: make([]*Residue, len(p.Residues)),
		Atoms:    make([]*Atom, len(p.Atoms)),
	}

	// Clone atoms
	atomMap := make(map[*Atom]*Atom)
	for i, atom := range p.Atoms {
		clonedAtom := &Atom{
			Serial:    atom.Serial,
			Name:      atom.Name,
			AltLoc:    atom.AltLoc,
			ResName:   atom.ResName,
			ChainID:   atom.ChainID,
			ResSeq:    atom.ResSeq,
			ICode:     atom.ICode,
			X:         atom.X,
			Y:         atom.Y,
			Z:         atom.Z,
			Occupancy: atom.Occupancy,
			TempFacto: atom.TempFacto,
			Element:   atom.Element,
		}
		clone.Atoms[i] = clonedAtom
		atomMap[atom] = clonedAtom
	}

	// Clone residues with updated atom pointers
	for i, res := range p.Residues {
		clonedRes := &Residue{
			Name:    res.Name,
			SeqNum:  res.SeqNum,
			ChainID: res.ChainID,
		}
		if res.N != nil {
			clonedRes.N = atomMap[res.N]
		}
		if res.CA != nil {
			clonedRes.CA = atomMap[res.CA]
		}
		if res.C != nil {
			clonedRes.C = atomMap[res.C]
		}
		if res.O != nil {
			clonedRes.O = atomMap[res.O]
		}
		clone.Residues[i] = clonedRes
	}

	return clone
}

// Sequence returns the amino acid sequence as a string
func (p *Protein) Sequence() string {
	if p == nil || len(p.Residues) == 0 {
		return ""
	}

	sequence := make([]byte, len(p.Residues))
	for i, res := range p.Residues {
		// Convert three-letter code to one-letter
		// Simplified mapping (expand as needed)
		code := threeToOne(res.Name)
		sequence[i] = code
	}
	return string(sequence)
}

// threeToOne converts three-letter amino acid code to one-letter
func threeToOne(threeLetter string) byte {
	mapping := map[string]byte{
		"ALA": 'A', "CYS": 'C', "ASP": 'D', "GLU": 'E',
		"PHE": 'F', "GLY": 'G', "HIS": 'H', "ILE": 'I',
		"LYS": 'K', "LEU": 'L', "MET": 'M', "ASN": 'N',
		"PRO": 'P', "GLN": 'Q', "ARG": 'R', "SER": 'S',
		"THR": 'T', "VAL": 'V', "TRP": 'W', "TYR": 'Y',
	}

	if code, ok := mapping[threeLetter]; ok {
		return code
	}
	return 'X' // Unknown
}
