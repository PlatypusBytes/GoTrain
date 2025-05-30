package main

import (
	"github.com/PlatypusBytes/GoTrain/internal/track_dispersion"
	"github.com/PlatypusBytes/GoTrain/pkg/math_utils"
	"fmt"
)

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
	// Create an array similar to numpy's linspace: np.linspace(0.1, 250, 100)
	omega := math_utils.Linspace(0.1, 250, 100) // 100 points from 0.1 to 250 rad/s
	track := ballast_dispersion.ComputeDispersion(ballast_parameters, omega)
	// Print track dispersion results (show first 5 and last 5 values)
	fmt.Println("Track Dispersion Results (first 5 and last 5 values):")
	fmt.Println("First 5 values:")
	for i := 0; i < 5 && i < len(track); i++ {
		fmt.Printf("Frequency %.2f rad/s: %v\n", omega[i], track[i])
	}

	if len(track) > 10 {
		fmt.Println("\nLast 5 values:")
		for i := len(track) - 5; i < len(track); i++ {
			fmt.Printf("Frequency %.2f rad/s: %v\n", omega[i], track[i])
		}
	}
}
