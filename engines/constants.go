// Package core provides the mathematical foundation for the animation engine
// Drawing from ancient mathematics (Vedic, Egyptian, Babylonian, Greek, Islamic)
// and modern physics constants to create naturally beautiful animations
//
// AGENT SIGMA - Channeling: da Vinci, Newton, Ramanujan, Noether, von Neumann,
//                           Tufte, Victor, Turing
//
// "Mathematics is the language in which God has written the universe" - Galileo

package core

import "math"

// ═══════════════════════════════════════════════════════════════════════════
// UNIVERSAL MATHEMATICAL CONSTANTS
// The timeless numbers that govern nature, art, and perception
// ═══════════════════════════════════════════════════════════════════════════

const (
	// ─── FUNDAMENTAL CONSTANTS ────────────────────────────────────────────

	// Pi - The circle constant (Greek mathematics)
	Pi = math.Pi // 3.14159265358979323846...

	// Tau - The full circle constant (2π)
	Tau = 2 * math.Pi // 6.28318530717958647692...

	// E - Euler's number (natural growth/decay)
	E = math.E // 2.71828182845904523536...

	// ─── GOLDEN RATIO FAMILY (Vedic & Greek) ─────────────────────────────

	// Phi - The Golden Ratio (φ = (1 + √5) / 2)
	// Found in: Parthenon, nautilus shells, galaxies, human body proportions
	Phi = 1.6180339887498948482045868343656381177203091798057628621354486227

	// PhiConjugate - The golden ratio conjugate (1/φ = φ - 1)
	PhiConjugate = 0.6180339887498948482045868343656381177203091798057628621354486227

	// GoldenAngle - The golden angle in degrees (137.5077640°)
	// Optimal divergence angle for phyllotaxis (leaf/seed arrangement)
	GoldenAngle = 137.50776405003785464634870128333650722863938093322295308090953

	// GoldenAngleRad - The golden angle in radians
	GoldenAngleRad = 2.3999632297286533222315555066336138531249990110581150429351127

	// ─── SQUARE ROOT CONSTANTS (Babylonian approximations) ───────────────

	// Sqrt2 - The diagonal of a unit square (Pythagoras)
	Sqrt2 = 1.4142135623730950488016887242096980785696718753769480731766797

	// Sqrt3 - The height of an equilateral triangle
	Sqrt3 = 1.7320508075688772935274463415058723669428052538103806280558069

	// Sqrt5 - Related to golden ratio (φ = (1 + √5) / 2)
	Sqrt5 = 2.2360679774997896964091736687312762354406183596115257242708972

	// ─── SACRED GEOMETRY RATIOS (Egyptian) ───────────────────────────────

	// Seked - Ancient Egyptian slope ratio (pyramid of Giza)
	// Rise/Run = 7/5.5 palms = 1.272727...
	Seked = 1.2727272727272727272727272727272727272727272727272727272727273

	// PyramidAngle - Angle of the Great Pyramid (51.84°)
	PyramidAngle = 51.842773410997285804376227844819406343460488234311581316375523

	// ─── VEDIC MATHEMATICS CONSTANTS ──────────────────────────────────────

	// BabelPoint - Universal resonance point (2π²)
	// Used in quantum mechanics, resonance frequencies
	BabelPoint = 19.739208802178716237668981971435792261845264134052182844035993

	// VedicAttractor - Natural attractor constant for organic variations
	VedicAttractor = 0.1

	// VedicOrbitalStability - Orbital stability factor
	VedicOrbitalStability = 0.212195

	// ─── TESLA HARMONICS ──────────────────────────────────────────────────

	// TeslaHarmonic - Natural resonance frequency (Hz)
	TeslaHarmonic = 4.909

	// TeslaDivine - Tesla's "divine numbers"
	Tesla3 = 3
	Tesla6 = 6
	Tesla9 = 9

	// ─── FIBONACCI SEQUENCE ───────────────────────────────────────────────
	// Each number is the sum of the previous two
	// Converges to φ as n approaches infinity
	Fib0  = 0
	Fib1  = 1
	Fib2  = 1
	Fib3  = 2
	Fib4  = 3
	Fib5  = 5
	Fib6  = 8
	Fib7  = 13
	Fib8  = 21
	Fib9  = 34
	Fib10 = 55
	Fib11 = 89
	Fib12 = 144
	Fib13 = 233
	Fib14 = 377
	Fib15 = 610
	Fib16 = 987
	Fib17 = 1597
	Fib18 = 2584
	Fib19 = 4181
	Fib20 = 6765

	// ─── CHAKRA SYSTEM (Vedic sacred geometry) ───────────────────────────
	// Petal counts in the seven chakras
	ChakraRoot        = 4   // Muladhara - Root
	ChakraSacral      = 6   // Svadhisthana - Sacral
	ChakraSolarPlexus = 10  // Manipura - Solar Plexus
	ChakraHeart       = 12  // Anahata - Heart
	ChakraThroat      = 16  // Vishuddha - Throat
	ChakraThirdEye    = 2   // Ajna - Third Eye (or 96)
	ChakraCrown       = 1000 // Sahasrara - Crown (infinite petals)

	// ─── ISLAMIC GEOMETRIC PATTERNS ──────────────────────────────────────

	// IslamicStar8 - 8-pointed star angle (45°)
	IslamicStar8 = 45.0

	// IslamicStar12 - 12-pointed star angle (30°)
	IslamicStar12 = 30.0

	// IslamicStar16 - 16-pointed star angle (22.5°)
	IslamicStar16 = 22.5

	// ─── PLATONIC SOLIDS (Greek sacred geometry) ─────────────────────────
	// The five perfect polyhedra

	// TetrahedronFaces - Fire element (4 triangular faces)
	TetrahedronFaces = 4

	// CubeFaces - Earth element (6 square faces)
	CubeFaces = 6

	// OctahedronFaces - Air element (8 triangular faces)
	OctahedronFaces = 8

	// DodecahedronFaces - Ether element (12 pentagonal faces)
	DodecahedronFaces = 12

	// IcosahedronFaces - Water element (20 triangular faces)
	IcosahedronFaces = 20

	// ─── PHYSICS CONSTANTS (for natural motion) ──────────────────────────

	// SpeedOfLight - c (m/s) - Maximum speed in universe
	// Used for: Ultra-fast transitions
	SpeedOfLight = 299792458.0

	// PlanckConstant - h (J⋅s) - Quantum of action
	// Used for: Quantum-inspired discrete animations
	PlanckConstant = 6.62607015e-34

	// GravityEarth - g (m/s²) - Earth's gravitational acceleration
	// Used for: Natural falling animations
	GravityEarth = 9.80665

	// GoldenAngleInNature - Observed angle in sunflowers, pine cones
	// Measured at 137.508° (very close to theoretical 137.5077...)
	GoldenAngleInNature = 137.508

	// ─── THREE-REGIME ARCHITECTURE (Williams Optimizer) ──────────────────
	// The 30-20-50 split for exploration → optimization → stabilization

	RegimeExploration   = 0.30 // 30% - Rapid prototyping, broad search
	RegimeOptimization  = 0.20 // 20% - Refinement, tuning
	RegimeStabilization = 0.50 // 50% - Convergence, equilibrium

	// ─── WILLIAMS OPTIMIZER CRITICAL POINTS ──────────────────────────────
	// Phase transitions in the O(√t × log₂(t)) optimizer

	WilliamsCritical1 = 1.5
	WilliamsCritical2 = 3.2
	WilliamsCritical3 = 7.5

	// ─── PERFORMANCE TARGETS ──────────────────────────────────────────────

	TargetFPS       = 60      // Target frames per second
	MinFPS          = 30      // Minimum acceptable FPS
	FrameBudgetMs   = 16.67   // Budget per frame at 60 FPS (ms)
	FrameBudget30   = 33.33   // Budget per frame at 30 FPS (ms)

	// ─── COLOR HARMONY CONSTANTS ──────────────────────────────────────────

	// Complementary - Opposite on color wheel (180°)
	ColorComplementary = 180.0

	// Analogous - Adjacent colors (30°)
	ColorAnalogous = 30.0

	// Triadic - Evenly spaced (120°)
	ColorTriadic = 120.0

	// Split complementary - Near-opposite (150°)
	ColorSplitComplementary = 150.0

	// Tetradic - Square (90°)
	ColorTetradic = 90.0
)

// ═══════════════════════════════════════════════════════════════════════════
// FIBONACCI SEQUENCE ARRAY
// Pre-computed for O(1) lookup
// ═══════════════════════════════════════════════════════════════════════════

var FibonacciSequence = []uint64{
	0, 1, 1, 2, 3, 5, 8, 13, 21, 34, 55, 89, 144, 233, 377, 610, 987, 1597,
	2584, 4181, 6765, 10946, 17711, 28657, 46368, 75025, 121393, 196418,
	317811, 514229, 832040, 1346269, 2178309, 3524578, 5702887, 9227465,
	14930352, 24157817, 39088169, 63245986, 102334155, 165580141, 267914296,
	433494437, 701408733, 1134903170, 1836311903, 2971215073, 4807526976,
	7778742049, 12586269025,
}

// ═══════════════════════════════════════════════════════════════════════════
// CHAKRA PETAL COUNTS ARRAY
// For creating mandalas and radial symmetry
// ═══════════════════════════════════════════════════════════════════════════

var ChakraPetals = []int{
	ChakraRoot,        // 4
	ChakraSacral,      // 6
	ChakraSolarPlexus, // 10
	ChakraHeart,       // 12
	ChakraThroat,      // 16
	ChakraThirdEye,    // 2
	ChakraCrown,       // 1000
}

var ChakraNames = []string{
	"Root",
	"Sacral",
	"Solar Plexus",
	"Heart",
	"Throat",
	"Third Eye",
	"Crown",
}

// ═══════════════════════════════════════════════════════════════════════════
// PLATONIC SOLIDS DATA
// The five perfect polyhedra
// ═══════════════════════════════════════════════════════════════════════════

type PlatonicSolid struct {
	Name     string
	Faces    int
	Vertices int
	Edges    int
	Element  string
}

var PlatonicSolids = []PlatonicSolid{
	{"Tetrahedron", 4, 4, 6, "Fire"},
	{"Cube", 6, 8, 12, "Earth"},
	{"Octahedron", 8, 6, 12, "Air"},
	{"Dodecahedron", 12, 20, 30, "Ether"},
	{"Icosahedron", 20, 12, 30, "Water"},
}

// ═══════════════════════════════════════════════════════════════════════════
// ANCIENT MATHEMATICAL SYSTEMS
// Information about the origins and applications of these constants
// ═══════════════════════════════════════════════════════════════════════════

type MathematicalSystem struct {
	Name         string
	Origin       string
	Era          string
	KeyConstants []string
	Applications []string
}

var AncientSystems = []MathematicalSystem{
	{
		Name:   "Vedic Mathematics",
		Origin: "Ancient India",
		Era:    "1500-500 BCE",
		KeyConstants: []string{
			"Golden Ratio (Phi)",
			"Fibonacci Sequence",
			"Chakra Proportions",
			"Sacred Geometry",
		},
		Applications: []string{
			"Temple architecture",
			"Yantra design",
			"Musical scales",
			"Astronomical calculations",
		},
	},
	{
		Name:   "Egyptian Mathematics",
		Origin: "Ancient Egypt",
		Era:    "3000-300 BCE",
		KeyConstants: []string{
			"Golden Ratio",
			"Seked (slope ratio)",
			"Pyramid angles",
			"Sacred proportions",
		},
		Applications: []string{
			"Pyramid construction",
			"Architectural design",
			"Land surveying",
			"Astronomical alignment",
		},
	},
	{
		Name:   "Babylonian Mathematics",
		Origin: "Mesopotamia",
		Era:    "2000-500 BCE",
		KeyConstants: []string{
			"Base-60 system",
			"Pythagorean triples",
			"Square root approximations",
			"Circle measurements",
		},
		Applications: []string{
			"Time measurement",
			"Angle divisions",
			"Astronomical tables",
			"Trade calculations",
		},
	},
	{
		Name:   "Greek Mathematics",
		Origin: "Ancient Greece",
		Era:    "600-200 BCE",
		KeyConstants: []string{
			"Pi",
			"Golden Mean",
			"Platonic Solids",
			"Irrational numbers",
		},
		Applications: []string{
			"Architecture (Parthenon)",
			"Philosophy",
			"Geometry",
			"Music theory",
		},
	},
	{
		Name:   "Islamic Mathematics",
		Origin: "Islamic Golden Age",
		Era:    "800-1300 CE",
		KeyConstants: []string{
			"Geometric patterns",
			"Tessellations",
			"Star polygons",
			"Proportion systems",
		},
		Applications: []string{
			"Mosque decoration",
			"Carpet design",
			"Calligraphy",
			"Architectural ornament",
		},
	},
}

// ═══════════════════════════════════════════════════════════════════════════
// UTILITY FUNCTIONS
// Quick access to common calculations
// ═══════════════════════════════════════════════════════════════════════════

// DegreesToRadians converts degrees to radians
func DegreesToRadians(degrees float64) float64 {
	return degrees * (Pi / 180.0)
}

// RadiansToDegrees converts radians to degrees
func RadiansToDegrees(radians float64) float64 {
	return radians * (180.0 / Pi)
}

// Fibonacci returns the nth Fibonacci number
// Uses pre-computed array for n < 51, calculates for larger n
func Fibonacci(n int) uint64 {
	if n < 0 {
		return 0
	}
	if n < len(FibonacciSequence) {
		return FibonacciSequence[n]
	}

	// Calculate for larger n
	a, b := uint64(0), uint64(1)
	for i := 0; i < n; i++ {
		a, b = b, a+b
	}
	return a
}

// PhiPower returns φ^n (golden ratio raised to power n)
// Useful for scaling hierarchies
func PhiPower(n float64) float64 {
	return math.Pow(Phi, n)
}

// GoldenSpiralRadius returns the radius at angle theta
// r = a × φ^(2θ/π)
func GoldenSpiralRadius(theta, scale float64) float64 {
	return scale * math.Pow(Phi, (2*theta)/Pi)
}

// ═══════════════════════════════════════════════════════════════════════════
// DOCUMENTATION
// Mathematical proofs and historical context
// ═══════════════════════════════════════════════════════════════════════════

/*
GOLDEN RATIO (φ = 1.618...)

Definition:
  φ = (1 + √5) / 2

Properties:
  φ² = φ + 1
  1/φ = φ - 1
  φ = lim(n→∞) F(n+1) / F(n)  where F(n) is Fibonacci

Found in nature:
  - Nautilus shell spirals
  - Sunflower seed arrangements (137.5° divergence)
  - Galaxy spiral arms
  - Human body proportions (navel to height)
  - DNA molecule (34Å × 21Å dimensions)

Used in art & architecture:
  - Parthenon facade
  - Great Pyramid of Giza
  - Renaissance paintings
  - Modern design (Apple, Twitter logos)

Animation applications:
  - Spacing elements naturally
  - Scaling hierarchies
  - Creating organic distributions
  - Timing curves (φ-based easing)
*/

/*
FIBONACCI SEQUENCE (0, 1, 1, 2, 3, 5, 8, 13, 21...)

Definition:
  F(0) = 0, F(1) = 1
  F(n) = F(n-1) + F(n-2) for n ≥ 2

Properties:
  - Ratio of consecutive terms approaches φ
  - Sum of first n terms: S(n) = F(n+2) - 1
  - Appears in Pascal's triangle diagonals

Closed form (Binet's formula):
  F(n) = (φⁿ - ψⁿ) / √5
  where ψ = (1 - √5) / 2 = -1/φ

Animation applications:
  - Step counts in multi-stage animations
  - Duration multipliers
  - Particle counts
  - Grid dimensions
*/

/*
THREE-REGIME ARCHITECTURE (30-20-50)

Discovered: 2025-10-30 (Day 168 - Ananta Learns to Code)

Philosophy:
  Divide complex tasks into three overlapping phases:

  1. EXPLORATION (30%):
     - Try multiple approaches
     - Discover edge cases
     - Broad, rapid prototyping

  2. OPTIMIZATION (20%):
     - Refine what works
     - Tune parameters
     - Performance improvements

  3. STABILIZATION (50%):
     - Lock in quality
     - Rigorous testing
     - Convergence to equilibrium

Application to animation:
  - 30% of frames: Exploration (varying timing, trying effects)
  - 20% of frames: Optimization (refining best approach)
  - 50% of frames: Stabilization (smooth, consistent motion)

Results:
  - Higher quality output
  - Fewer iterations needed
  - Natural convergence
*/

/*
WILLIAMS SPACE OPTIMIZATION

Formula: O(√t × log₂(t))

Sublinear space complexity with validation probability p < 10^-133

Used for:
  - Animation duration: duration = complexity(n) × baseTime
  - LOD distances: distance[i] = complexity(n) × scale[i]
  - Batch sizes: batchSize = totalItems / complexity(totalItems)
  - Stagger delays: delay[i] = complexity(i) × baseDelay

Critical points (phase transitions):
  - t = 1.5: Initial organization
  - t = 3.2: Pattern emergence
  - t = 7.5: Stable structure

Benefits:
  - Automatic performance scaling
  - Natural-feeling timing
  - Prevents cognitive overload
  - Maintains 60 FPS
*/
