package ballast_dispersion

import (
    "fmt"
    "math"
    "gonum.org/v1/gonum/mat"
    "GoTrain/pkg/math_utils"
)


// BallastTrackParameters holds the parameters for the ballast track dispersion model.
// These parameters are used to define the physical properties of the track,
type BallastTrackParameters struct {
    EIRail         float64 // Rail bending stiffness [N·m^2].
    MRail          float64 // Rail mass per unit length [kg/m].
    KRailPad       float64 // Railpad stiffness [N/m].
    CRailPad       float64 // Railpad damping [N·s/m].
    MSleeper       float64 // Sleeper (distributed) mass [kg/m].
    EBallast       float64 // Young's modulus of ballast [Pa].
    HBallast       float64 // Ballast (layer) thickness [m].
    WidthSleeper   float64 // Half-track width [m].
    RhoBallast     float64 // Ballast density [kg/m^3].
    SoilStiffness  float64 // Soil (spring) stiffness [N/m].
}


func ComputeDispersion(parameters BallastTrackParameters, omega []float64)  []float64 {

    // Convert angular frequency to Hz
    // frequency := make([]float64, len(omega))
    phase_velocity := make([]float64, len(omega))

    ini_wave_number := 0.001
    end_wave_number := 1000.0

    for i, omegaVal := range omega {
        // frequency[i] = omegaVal / (2.0 * math.Pi)
        // TrackStiffness(parameters, omegaVal, wave_number)

        // Define a function for the Brent method to find the wave number
        brentAuxiliar := func(wavenumber float64) float64 {
            return TrackStiffness(parameters, omegaVal, wavenumber)
        }

        wavenumber, err := math_utils.Brent(ini_wave_number, end_wave_number, 1e-9, brentAuxiliar)
        if err != nil {
            fmt.Println(err.Error())
        } else {
            // Calculate phase velocity from the found wave number
            phase_velocity[i] = omegaVal / wavenumber
        }
    }
    return phase_velocity

}



func TrackStiffness(parameters BallastTrackParameters, omega float64, wavenumber float64) float64 {

    // constant alpha
    alpha := 0.5
    // compression wave in ballast
    cp := math.Sqrt(parameters.EBallast / parameters.RhoBallast)

    // auxiliar values
    tan_value := math.Tan(omega * parameters.HBallast / cp) * cp
    sin_value := math.Sin(omega * parameters.HBallast / cp) * cp

    // railpad complex stiffness
    // rail_pad_complex_stiffness := complex(parameters.KRailPad, omega * parameters.CRailPad)
    rail_pad_complex_stiffness := parameters.KRailPad

    // stiffness matrix
    k11 := parameters.EIRail * math.Pow(wavenumber, 4) + rail_pad_complex_stiffness - math.Pow(omega, 2) * parameters.MRail
    k12 := -rail_pad_complex_stiffness
    k22 := rail_pad_complex_stiffness + (2 * omega * parameters.EBallast * parameters.WidthSleeper * alpha) / tan_value -
            math.Pow(omega, 2) * parameters.MSleeper
    k23 := -2 * omega * parameters.EBallast * parameters.WidthSleeper * alpha / sin_value
    k33 := 2 * omega * parameters.EBallast * parameters.WidthSleeper * alpha / tan_value + parameters.SoilStiffness

    stiffness := mat.NewDense(3, 3, []float64{
        k11, k12, 0,
        k12, k22, k23,
        0, k23, k33,
    })

    // Calculate the determinant of the stiffness matrix
    det := mat.Det(stiffness)

    return det

}
