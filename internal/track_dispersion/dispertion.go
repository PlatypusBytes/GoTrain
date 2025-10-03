// Package track_dispersion provides tools to compute phase velocity dispersion curves
// for railway track systems.
//
// The computation of the dispersion curve is based on the formulation:
// Mezher et al. (2016). "Railway critical velocity - Analytical prediction and analysis".
// Transportation Geotechnics, 6, 84–96.  See also: https://doi.org/10.1016/j.trgeo.2015.09.002
//
// It supports two types of track systems:
//   - Ballast tracks, modeled with rail, sleeper, railpad, ballast and soil
//   - Slab tracks, modeled with rail, slab, railpad, and soil
package track_dispersion

import (
	"fmt"
	"math"

	"github.com/PlatypusBytes/GoTrain/pkg/utils"
	"gonum.org/v1/gonum/mat"
)

// TrackParameters defines the interface that track parameter structs must implement
type TrackParameters interface {
	CalculateStiffness(omega float64, wavenumber float64) float64
}

// BallastTrackParameters holds the parameters for the ballast track model.
// These parameters are used to define the physical properties of the railway track,
// including rail, sleeper, railpad, ballast, and soil.
type BallastTrackParameters struct {
	EIRail        float64 // Rail bending stiffness [N·m^2].
	MRail         float64 // Rail mass per unit length [kg/m].
	KRailPad      float64 // Railpad stiffness [N/m].
	CRailPad      float64 // Railpad damping [N·s/m].
	MSleeper      float64 // Sleeper (distributed) mass [kg/m].
	EBallast      float64 // Young's modulus of ballast [Pa].
	HBallast      float64 // Ballast (layer) thickness [m].
	WidthSleeper  float64 // Half-track width [m].
	RhoBallast    float64 // Ballast density [kg/m^3].
	SoilStiffness float64 // Soil (spring) stiffness [N/m].
}

// CalculateStiffness implements the TrackParameters interface for BallastTrackParameters
func (p BallastTrackParameters) CalculateStiffness(omega float64, wavenumber float64) float64 {
	return BallastTrackStiffness(p, omega, wavenumber)
}

// SlabTrackParameters holds the parameters for the slab track model.
// These parameters define the physical properties of a slab track system,
// including rail, slab, railpad, and soil.
type SlabTrackParameters struct {
	EIRail        float64 // Rail bending stiffness [N·m^2].
	MRail         float64 // Rail mass per unit length [kg/m].
	EISlab        float64 // Slab bending stiffness [N·m^2].
	MSlab         float64 // Slab mass per unit length [kg/m].
	KRailPad      float64 // Railpad stiffness [N/m].
	CRailPad      float64 // Railpad damping [N·s/m].
	SoilStiffness float64 // Soil (spring) stiffness [N/m].
}

// CalculateStiffness implements the TrackParameters interface for SlabTrackParameters
func (p SlabTrackParameters) CalculateStiffness(omega float64, wavenumber float64) float64 {
	return SlabTrackStiffness(p, omega, wavenumber)
}

// RailTrackDispersion calculates the phase velocity dispersion curve for a railway track.
//
// Parameters:
//   - parameters: Physical parameters of the track system (BallastTrackParameters or SlabTrackParameters)
//   - omega: Array of angular frequencies [rad/s] at which to compute phase velocities
//
// Returns:
//   - An array of phase velocities [m/s] corresponding to each input angular frequency
func RailTrackDispersion(parameters TrackParameters, omega []float64) []float64 {

	phase_velocity := make([]float64, len(omega))

	ini_wave_number := 0.001
	end_wave_number := 1000.0

	for i, omegaVal := range omega {
		// Define a function for the Brent method to find the wave number
		brentAuxiliar := func(wavenumber float64) float64 {
			return parameters.CalculateStiffness(omegaVal, wavenumber)
		}

		wavenumber, err := math_utils.Brent(brentAuxiliar, ini_wave_number, end_wave_number, 1e-12)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			// Calculate phase velocity from the found wave number
			phase_velocity[i] = omegaVal / wavenumber
		}
	}
	return phase_velocity
}

// BallastTrackStiffness computes the determinant of the track-soil system stiffness matrix
// for a given angular frequency and wavenumber. This function is used in dispersion analysis
// to identify combinations of frequency and wavenumber where the determinant is zero,
// which correspond to wave propagation modes in the track-soil system.
//
// Parameters:
//   - parameters: Physical parameters of the ballast track system
//   - omega: Angular frequency [rad/s]
//   - wavenumber: Spatial frequency [1/m]
//
// Returns:
//   - Determinant of the 3x3 stiffness matrix representing the track-soil system
func BallastTrackStiffness(parameters BallastTrackParameters, omega float64, wavenumber float64) float64 {

	// constant alpha
	alpha := 0.5
	// compression wave in ballast
	cp := math.Sqrt(parameters.EBallast / parameters.RhoBallast)

	// auxiliar values
	tan_value := math.Tan(omega*parameters.HBallast/cp) * cp
	sin_value := math.Sin(omega*parameters.HBallast/cp) * cp

	// railpad complex stiffness
	// rail_pad_complex_stiffness := complex(parameters.KRailPad, omega * parameters.CRailPad)
	rail_pad_complex_stiffness := parameters.KRailPad

	// stiffness matrix
	k11 := parameters.EIRail*math.Pow(wavenumber, 4) + rail_pad_complex_stiffness - math.Pow(omega, 2)*parameters.MRail
	k12 := -rail_pad_complex_stiffness
	k22 := rail_pad_complex_stiffness + (2*omega*parameters.EBallast*parameters.WidthSleeper*alpha)/tan_value -
		math.Pow(omega, 2)*parameters.MSleeper
	k23 := -2 * omega * parameters.EBallast * parameters.WidthSleeper * alpha / sin_value
	k33 := 2*omega*parameters.EBallast*parameters.WidthSleeper*alpha/tan_value + parameters.SoilStiffness

	stiffness := mat.NewDense(3, 3, []float64{
		k11, k12, 0,
		k12, k22, k23,
		0, k23, k33,
	})

	// Calculate the determinant of the stiffness matrix
	det := mat.Det(stiffness)

	return det
}

// SlabTrackStiffness computes the determinant of the track-soil system stiffness matrix
// for a given angular frequency and wavenumber for slab track systems.
//
// Parameters:
//   - parameters: Physical parameters of the slab track system
//   - omega: Angular frequency [rad/s]
//   - wavenumber: Spatial frequency [1/m]
//
// Returns:
//   - Determinant of the stiffness matrix representing the track-soil system
func SlabTrackStiffness(parameters SlabTrackParameters, omega float64, wavenumber float64) float64 {
	// rail_pad_complex_stiffness := complex(parameters.KRailPad, omega * parameters.CRailPad)
	rail_pad_complex_stiffness := parameters.KRailPad

	// stiffness matrix
	k11 := parameters.EIRail*math.Pow(wavenumber, 4) + rail_pad_complex_stiffness - math.Pow(omega, 2)*parameters.MRail
	k12 := -rail_pad_complex_stiffness
	k22 := rail_pad_complex_stiffness + parameters.EISlab*math.Pow(wavenumber, 4) - math.Pow(omega, 2)*parameters.MSlab + parameters.SoilStiffness

	stiffness := mat.NewDense(2, 2, []float64{
		k11, k12,
		k12, k22,
	})

	// Calculate the determinant of the stiffness matrix
	det := mat.Det(stiffness)

	return det
}
