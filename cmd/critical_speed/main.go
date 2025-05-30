// Package main provides a command-line tool for calculating
// the critical speeds of railway tracks using dispersion analysis.
package main

import (
    "encoding/json"
    "fmt"
    "os"

    ballast_dispersion "github.com/PlatypusBytes/GoTrain/internal/track_dispersion"
    "github.com/PlatypusBytes/GoTrain/pkg/math_utils"
)

// main executes the critical speed analysis for a ballasted railway track.
func main() {
    // Initialize the parameters for the ballasted track
    ballast_parameters := ballast_dispersion.BallastTrackParameters{
        EIRail:        1.29e7,
        MRail:         120,
        KRailPad:      5e8,
        CRailPad:      2.5e5,
        MSleeper:      490,
        EBallast:      130e6,
        HBallast:      0.35,
        WidthSleeper:  1.25,
        RhoBallast:    1700,
        SoilStiffness: 0.0,
    }

    // Define the angular frequencies for the dispersion calculation
    omega := math_utils.Linspace(0.1, 250, 100) // 100 points from 0.1 to 250 rad/s

    // Calculate the phase velocity dispersion curve for the ballasted track
    track := ballast_dispersion.ComputeDispersion(ballast_parameters, omega)


    // Write the results into a json file
    type DispersionResults struct {
        Omega         []float64 `json:"omega"`
        PhaseVelocity []float64 `json:"phase_velocity"`
    }

    results := DispersionResults{
        Omega:         omega,
        PhaseVelocity: track,
    }

    jsonData, err := json.MarshalIndent(results, "", "\t")
    if err != nil {
        fmt.Println("Error marshaling to JSON:", err)
        return
    }

    // Write to file
    err = os.WriteFile("track_dispersion_results.json", jsonData, 0644)
    if err != nil {
        fmt.Println("Error writing JSON to file:", err)
        return
    }

    fmt.Println("Results written successfully to dispersion_results.json")
}
