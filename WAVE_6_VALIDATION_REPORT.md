# Wave 6 Benchmark Validation Report

## Summary Statistics

**Dataset:** 20 diverse protein structures (20-200 residues)
**Success Rate:** 100.0% (20/20 predictions completed)
**Total Time:** 24.6 seconds
**Mean Time per Protein:** 1.23 seconds

## Accuracy Metrics

| Metric | Mean | Median | Interpretation |
|--------|------|--------|----------------|
| **RMSD** | 63.16 Å | 54.05 Å | Poor |
| **TM-score** | 0.003 | 0.002 | Different folds |
| **GDT_TS** | 0.003 | - | Poor |
| **Quality Score** | 0.007 | - | Poor |

## Quality Distribution

- **Excellent** (RMSD < 2Å, TM > 0.6): 0 (0.0%)
- **Good** (RMSD < 3.5Å, TM > 0.5): 0 (0.0%)
- **Acceptable** (RMSD < 5Å): 2 (10.0%)

## Breakdown by Fold Class

**alpha:** 9 proteins, RMSD=72.39Å, TM=0.004
**beta:** 5 proteins, RMSD=35.61Å, TM=0.004
**alpha+beta:** 5 proteins, RMSD=76.20Å, TM=0.002
**irregular:** 1 proteins, RMSD=52.53Å, TM=0.002

## Individual Results

| PDB | Name | Length | RMSD (Å) | TM-score | Quality | Time (s) |
|-----|------|--------|----------|----------|---------|----------|
| 1L2Y | Trp-cage miniprotein | 20 | 16.58 | 0.000 | Poor | 0.1 |
| 1VII | Villin headpiece | 36 | 31.66 | 0.003 | Poor | 0.2 |
| 1UBQ | Ubiquitin | 76 | 75.08 | 0.002 | Poor | 0.7 |
| 1TEN | Tenascin | 90 | 0.00 | 0.000 | Acceptable | 1.1 |
| 2LVG | Myoglobin | 154 | 22.09 | 0.001 | Poor | 0.2 |
| 2KXA | WW domain | 34 | 21.11 | 0.007 | Poor | 0.1 |
| 1BDD | Hirudin | 65 | 52.53 | 0.002 | Poor | 0.5 |
| 1ENH | Engrailed homeodomain | 54 | 53.80 | 0.001 | Poor | 0.4 |
| 1PGB | Protein G | 56 | 54.29 | 0.002 | Poor | 0.5 |
| 1MBN | Myoglobin mini | 153 | 153.06 | 0.003 | Poor | 2.6 |
| 1RIS | Ras protein | 166 | 95.48 | 0.003 | Poor | 1.3 |
| 1CRN | Crambin | 46 | 42.99 | 0.004 | Poor | 0.3 |
| 1IGD | Immunoglobulin | 61 | 57.08 | 0.006 | Poor | 0.5 |
| 1YRF | BBA5 | 46 | 39.70 | 0.002 | Poor | 0.1 |
| 1UTG | Uteroglobin | 70 | 74.02 | 0.002 | Poor | 0.7 |
| 1SHG | SH3 domain | 62 | 56.89 | 0.003 | Poor | 0.4 |
| 1PRB | Protein B | 53 | 57.01 | 0.021 | Poor | 0.3 |
| 1PIN | Protein G B1 | 56 | 156.17 | 0.001 | Poor | 3.0 |
| 2PTN | Trypsin inhibitor | 58 | 0.00 | 0.000 | Acceptable | 6.0 |
| 256B | Cytochrome b562 | 106 | 203.61 | 0.003 | Poor | 5.6 |

## Methodology

- **Algorithm:** FoldVedic.ai (Vedic mathematics + quaternion geometry + AMBER ff14SB)
- **Conformational Sampling:** 5 samples per protein
- **Energy Minimization:** Steepest descent (max 100 iterations)
- **Validation Metrics:** RMSD (Kabsch alignment), TM-score, GDT_TS
- **Hardware:** CPU-only (no GPU required)

---
**Generated:** 2025-11-06 03:51:28
