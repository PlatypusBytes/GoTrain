// Package gotrain is a Go library for analyzing critical speeds in railway systems,
// focusing on soil and track dispersion analysis.
//
// This project is based on TrainCritSpeed (https://github.com/PlatypusBytes/TrainCritSpeed),
// which is implemented in Python. GoTrain aims to provide similar functionality with
// improved performance and concurrency features.
//
// The main difference with TrainCritSpeed is that GoTrain only computes the fundamental
// mode for the subsurface and GoTrain does not make the dispersion field plot. If you
// wish to use these features, please use TrainCritSpeed.
//
// # Methodology
//
// The methodology for the computation of the critical train speed is based on the work
// of Mezher et al. (2016): "Railway critical velocity - Analytical prediction and analysis".
// Transportation Geotechnics, 6, 84–96.
// https://www.sciencedirect.com/science/article/abs/pii/S2214391215000239
//
// The dispersion curve for the layered soil is based on the Fast Delta Matrix method
// proposed by Buchen and Ben-Hador (1996): "Free-mode surface-wave computations".
// Geophysical Journal International, 124(3), 869–887.
// https://academic.oup.com/gji/article-lookup/doi/10.1111/j.1365-246X.1996.tb05642.x
//
// # Commands
//
// GoTrain provides two main commands:
//
// Critical Speed Calculator (cmd/critical_speed): Calculates dispersion curves and
// critical speeds for a single configuration file.
//
// Batch Runner (cmd/runner): Processes multiple configuration files in parallel with
// configurable worker pools.
//
// # Package Structure
//
// The library is organized into the following packages:
//
//   - internal/soil_dispersion: Computes phase velocity dispersion curves for soil profiles
//   - internal/track_dispersion: Computes phase velocity dispersion curves for railway track systems
//   - pkg/utils: Mathematical utility functions for numerical computations
//
// # Example Usage
//
// To use GoTrain as a library, import the necessary packages:
//
//	import (
//		"github.com/PlatypusBytes/GoTrain/internal/soil_dispersion"
//		"github.com/PlatypusBytes/GoTrain/internal/track_dispersion"
//		"github.com/PlatypusBytes/GoTrain/pkg/utils"
//	)
//
// For command-line usage, see the documentation for cmd/critical_speed and cmd/runner.
package gotrain
