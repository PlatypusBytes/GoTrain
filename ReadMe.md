![Tests](https://github.com/PlatypusBytes/GoTrain/actions/workflows/tests.yaml/badge.svg)
[![Go Reference](https://pkg.go.dev/badge/github.com/PlatypusBytes/GoTrain.svg)](https://pkg.go.dev/github.com/PlatypusBytes/GoTrain)
[![codecov](https://codecov.io/gh/PlatypusBytes/GoTrain/graph/badge.svg)](https://codecov.io/gh/PlatypusBytes/GoTrain)
[![Go Report Card](https://goreportcard.com/badge/github.com/PlatypusBytes/GoTrain)](https://goreportcard.com/report/github.com/PlatypusBytes/GoTrain)
[![License](https://img.shields.io/badge/License-BSD_3--Clause-blue.svg)](https://opensource.org/licenses/BSD-3-Clause)
[![DOI](https://zenodo.org/badge/993186160.svg)](https://doi.org/10.5281/zenodo.18744754)


# GoTrain

A high-performance Go library for analyzing critical speeds in railway systems, focusing on soil and track dispersion analysis.

## Overview

GoTrain computes the speed at which critical train speed occurs. The critical train speed is the speed at which the train speed matches the phase velocity of waves propagating through the track-soil system.

## Key Features

  - Critical speed calculation for railway track-soil systems
  - Support for both ballast and slab track configurations
  - Multi-layered soil profile modelling with elastic properties
  - High-performance parallel batch processing capabilities
  - JSON output format for integration with other tools

## Background

This project is based on [TrainCritSpeed](https://github.com/PlatypusBytes/TrainCritSpeed), originally implemented in Python. GoTrain reimplements the core functionality in Go, providing improved performance and native concurrency support.

**Main differences from TrainCritSpeed:**
- Computes only the fundamental mode for subsurface layers
- Does not generate dispersion field plots
- Significantly faster execution with Go's performance characteristics
- Built-in parallel processing for batch operations

For advanced features like higher-order modes and dispersion field visualization, please use the original [TrainCritSpeed](https://github.com/PlatypusBytes/TrainCritSpeed) Python implementation.

## Methodology

The critical speed computation is based on established scientific methods:

**Critical Speed Analysis**

Mezher, S. B., Connolly, D. P., Woodward, P. K., Laghrouche, O., Pombo, J., & Costa, P. A. (2016).
"Railway critical velocity - Analytical prediction and analysis"
*Transportation Geotechnics*, 6, 84–96.
[https://doi.org/10.1016/j.trgeo.2015.09.002](https://doi.org/10.1016/j.trgeo.2015.09.002)

The critical speed is identified at the intersection point of the track and soil dispersion curves, where the phase velocities match.

**Soil Dispersion Computation**

Buchen, P. W., & Ben-Hador, R. (1996).
"Free-mode surface-wave computations"
*Geophysical Journal International*, 124(3), 869–887.
[https://doi.org/10.1111/j.1365-246X.1996.tb05642.x](https://doi.org/10.1111/j.1365-246X.1996.tb05642.x)

The soil dispersion curves are computed using the Fast Delta Matrix method, which efficiently handles multi-layered soil profiles with varying elastic properties.

## Architecture

GoTrain is organized into several key components:

```
GoTrain/
├── cmd/
│   ├── critical_speed/     # Single configuration analyzer
│   └── runner/             # Batch processor
├── internal/
│   ├── critical_speed/     # Core critical speed analysis engine
│   ├── runner/             # Parallel batch processor
│   ├── soil_dispersion/    # Soil dispersion (Fast Delta Matrix)
│   └── track_dispersion/   # Track dispersion (ballast & slab)
├── pkg/
│   └── utils/              # Mathematical utilities (Brent's method, etc.)
├── configs/                # Sample configuration files
└── testdata/               # Test data and fixtures
```

**Component Descriptions:**
- `internal/critical_speed` - Core critical speed analysis engine
- `internal/runner` - Parallel batch processor for multiple configurations
- `internal/soil_dispersion` - Soil dispersion curve computation (Fast Delta Matrix)
- `internal/track_dispersion` - Track dispersion curve computation (ballast & slab tracks)
- `pkg/utils` - Mathematical utilities (Brent's method, linear interpolation, etc.)

## Installation

### Option 1: Download Pre-built Binaries (Recommended)

Download the latest release for your platform from the [GitHub Releases page](https://github.com/PlatypusBytes/GoTrain/releases).

You can download `critical_speed` (single configuration calculator) and `runner` (batch processor) directly.

**Available platforms:**
- Linux (amd64)
- Windows (amd64)

### Option 2: Build from Source

**Prerequisites:**
- Go 1.24.3 or later
- Git

**Build steps:**

```bash
# Clone the repository
git clone https://github.com/PlatypusBytes/GoTrain.git
cd GoTrain

# Build using Makefile
make build
```

This creates two executables in the `bin/` directory:
- `bin/critical_speed` - Single configuration calculator
- `bin/runner` - Batch processor for multiple configurations

## Commands

GoTrain provides two main command-line tools:

### 1. Critical Speed Calculator

Analyzes a single railway configuration and computes dispersion curves and critical speed.

**Usage:**
```bash
./critical_speed -config configs/sample_config.yaml
```

**What it does:**
- Loads configuration from YAML file
- Computes track dispersion curve (ballast or slab track)
- Computes soil dispersion curve (multi-layered profile)
- Identifies critical speed (intersection of dispersion curves)
- Outputs results to JSON file

**Output:** A JSON file containing omega values, track phase velocities, soil phase velocities, critical omega, and critical velocity.

**Command-line flags:**
- `-config` (required): Path to YAML configuration file

### 2. Batch Runner (`runner`)

Processes multiple YAML configuration files in parallel with configurable worker pools. Automatically discovers all `.yaml` files in a directory tree and processes them concurrently.

**Usage:**
```bash
./runner -dir testdata/batch -workers 4
```

**What it does:**
- Recursively scans directory for `.yaml` files
- Spawns worker goroutines for parallel processing
- Displays real-time progress bar
- Processes each configuration using `critical_speed` logic
- Maximizes throughput with concurrent execution

**Example output:**
```
Found 10 YAML files to process
[==========================                        ] 52.00% (5/10)

...

Completed processing 10 YAML files
```

**Command-line flags:**
- `-dir` (required): Directory containing YAML configuration files
- `-workers` (optional): Number of parallel workers (default: number of CPU cores)

## Configuration

Configuration files use YAML format and must specify:

- **Track type**: `"ballast"` or `"slabtrack"`
- **Frequency range**: min, max, and number of points
- **Track parameters**: rail, sleeper/slab, railpad properties
- **Soil layers**: multi-layer profile with elastic properties
- **Output**: JSON filename for results

### Example Configuration

An example configuration file is located at [`configs/sample_config.yaml`](configs/sample_config.yaml):

```yaml
# Track type: can be "ballast" or "slabtrack"
track_type: ballast

# Frequency range configuration
frequency:
  min: 1
  max: 314
  points: 100

# Ballast track parameters
ballast_track:
  EI_rail: 6.4e6         # Rail bending stiffness [N·m^2]
  m_rail: 60.21          # Rail mass per unit length [kg/m]
  k_rail_pad: 6e8        # Railpad stiffness [N/m]
  c_rail_pad: 2.5e5      # Railpad damping [N·s/m]
  m_sleeper: 238.5       # Sleeper (distributed) mass [kg/m]
  E_ballast: 100e6       # Young's modulus of ballast [Pa]
  h_ballast: 0.3         # Ballast (layer) thickness [m]
  width_sleeper: 1.25    # Half-track width [m]
  rho_ballast: 2000      # Ballast density [kg/m^3]
  soil_stiffness: 0.0    # Soil (spring) stiffness [N/m]

# Slab track parameters
slab_track:
  EI_rail: 1.29e7        # Rail bending stiffness [N·m^2]
  m_rail: 120            # Rail mass per unit length [kg/m]
  EI_slab: 6.40625e8     # Slab bending stiffness [N·m^2] (calculated from 30e9 * (1.25 * 0.35^3 / 12))
  m_slab: 1093.75        # Slab mass per unit length [kg/m] (calculated from 2500*1.25*0.35)
  k_rail_pad: 5e8        # Railpad stiffness [N/m]
  c_rail_pad: 2.5e5      # Railpad damping [N·s/m]
  soil_stiffness: 0.0    # Soil (spring) stiffness [N/m]

soil_layers:
  - thickness: 5          # Thickness of the soil layer [m]
    density: 1900         # Density of the soil layer [kg/m^3]
    young_modulus: 50e6   # Young  modulus of the soil layer [Pa]
    poisson_ratio: 0.3    # Poisson's ratio of the soil layer
  - thickness: 10         # Thickness of the second soil layer [m]
    density: 1900         # Density of the second soil layer [kg/m^3]
    young_modulus: 200e6  # Young  modulus of the second soil layer [Pa]
    poisson_ratio: 0.3    # Poisson's ratio of the second soil layer
  - thickness: 15         # Thickness of the third soil layer [m]
    density: 1900         # Density of the third soil layer [kg/m^3]
    young_modulus: 456e6  # Young  modulus of the third soil layer [Pa]
    poisson_ratio: 0.3    # Poisson's ratio of the third soil layer
  - thickness: .inf       # Thickness of the fourth soil layer [m]
    density: 1900         # Density of the fourth soil layer [kg/m^3]
    young_modulus: 810e6  # Young  modulus of the fourth soil layer [Pa]
    poisson_ratio: 0.33   # Poisson's ratio of the fourth soil layer

# Output file configuration
output:
  file_name: "dispersion_results.json"
```
## Output Format

Results are saved as JSON files with the following structure:

```json
{
  "omega": [1.0, 4.14, 7.28, ...],
  "track_phase_velocity": [245.3, 251.7, 258.1, ...],
  "soil_phase_velocity": [183.5, 185.2, 187.0, ...],
  "critical_omega": 125.66,
  "critical_velocity": 198.45
}
```

**Field descriptions:**
- `omega` - Angular frequencies [rad/s]
- `track_phase_velocity` - Phase velocities in track system [m/s]
- `soil_phase_velocity` - Phase velocities in soil layers [m/s]
- `critical_omega` - Critical angular frequency [rad/s]
- `critical_velocity` - Critical train speed [m/s]

## Examples: Typical Workflow

**Single Project Analysis:**

1. Create a YAML configuration file with your track and soil parameters
2. Run the analysis: `./critical_speed -config my_project.yaml`
3. Review the output JSON file
4. Adjust parameters if needed

**Parametric Studies:**

1. Create multiple YAML files with parameter variations
2. Organize them in a directory structure:
   ```
   parametric_study/
   ├── soft_soil/
   │   ├── config_1.yaml
   │   ├── config_2.yaml
   │   └── ...
   └── stiff_soil/
       ├── config_1.yaml
       └── ...
   ```
3. Run batch processing: `./runner -dir parametric_study -workers 8`
4. Compare results across configurations


## Contributing

Contributions are welcome! Please feel free to create a Fork and submit a Pull Request.

## Authors

See the [AUTHORS](AUTHORS) file for the list of contributors.


## License

This project is licensed under the BSD-3-Clause License. See the [LICENSE](LICENSE) file for details.
