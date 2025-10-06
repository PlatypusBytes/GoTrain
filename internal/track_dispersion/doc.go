// Package track_dispersion provides tools to compute phase velocity dispersion curves
// for railway track systems.
//
// The computation of the dispersion curve is based on the formulation:
// Mezher et al. (2016). "Railway critical velocity - Analytical prediction and analysis".
// Transportation Geotechnics, 6, 84–96.
// https://doi.org/10.1016/j.trgeo.2015.09.002
//
// # Supported Track Types
//
// The package supports two types of track systems:
//
//   - Ballast tracks: Modeled with rail, sleeper, railpad, ballast and soil
//   - Slab tracks: Modeled with rail, slab, railpad, and soil
//
// # Track Parameters
//
// The TrackParameters interface defines the contract that track parameter structs
// must implement. Two concrete implementations are provided:
//
//   - BallastTrackParameters: Holds parameters for ballast track models including
//     rail bending stiffness, rail mass, railpad properties, sleeper mass, ballast
//     properties, and soil stiffness.
//
//   - SlabTrackParameters: Holds parameters for slab track models including rail
//     bending stiffness, rail mass, railpad properties, slab properties, and soil
//     stiffness.
//
// # Dispersion Calculation
//
// The TrackDispersion function calculates the phase velocity dispersion curve for
// a railway track system using a numerical eigenvalue approach. It solves the
// dynamic equilibrium equations for the track-soil system at each frequency to
// determine the phase velocities.
//
// # Usage Example
//
//	params := track_dispersion.BallastTrackParameters{
//		EIRail:        6.4e6,
//		MRail:         60.21,
//		KRailPad:      6e8,
//		CRailPad:      2.5e5,
//		MSleeper:      238.5,
//		EBallast:      100e6,
//		HBallast:      0.3,
//		WidthSleeper:  1.25,
//		RhoBallast:    2000,
//		SoilStiffness: 0.0,
//	}
//	omega := math_utils.Linspace(1, 314, 100)
//	soilPhaseVelocity := []float64{...}
//	phaseVelocities := track_dispersion.TrackDispersion(params, omega, soilPhaseVelocity)
package track_dispersion
