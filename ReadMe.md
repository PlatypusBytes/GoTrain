![Tests](https://github.com/PlatypusBytes/GoTrain/actions/workflows/tests.yaml/badge.svg)
[![Go Reference](https://pkg.go.dev/badge/github.com/PlatypusBytes/GoTrain.svg)](https://pkg.go.dev/github.com/PlatypusBytes/GoTrain)
[![codecov](https://codecov.io/gh/PlatypusBytes/GoTrain/graph/badge.svg)](https://codecov.io/gh/PlatypusBytes/GoTrain)
[![Go Report Card](https://goreportcard.com/badge/github.com/PlatypusBytes/GoTrain)](https://goreportcard.com/report/github.com/PlatypusBytes/GoTrain)


# GoTrain

A Go library for analyzing critical speeds in railway systems, focusing on soil and track dispersion analysis.
This project is based on [TrainCritSpeed](https://github.com/PlatypusBytes/TrainCritSpeed), which is implemented in Python. GoTrain aims to provide similar functionality with improved performance and concurrency features.
The main difference with TrainCritSpeed is that GoTrain only computes the fundamental mode for the subsurface and
GoTrain does not make the dispersion field plot. If you wish to use these features, please use TrainCritSpeed.

The methodology for the computation of the critical train speed is based on the work of [Mezher et al. (2016)](https://www.sciencedirect.com/science/article/abs/pii/S2214391215000239).
The dispersion curve for the layered soil is based on the Fast Delta Matrix method proposed by [Buchen and Ben-Hador (1996)](https://academic.oup.com/gji/article-lookup/doi/10.1111/j.1365-246X.1996.tb05642.x).


## Building

To build the executables:

```bash
# Build the critical_speed command
go build -o bin/critical_speed ./cmd/critical_speed

# Build the batch runner command
go build -o bin/runner ./cmd/runner
```

## Commands

GoTrain provides two main commands:

### Critical Speed Calculator

The `critical_speed` command calculates dispersion curves and critical speeds for a single configuration:

```bash
# Run with a configuration file (required)
./bin/critical_speed -config /path/to/config.yaml
```

### Batch Runner

The `runner` command processes multiple configuration files in parallel:

```bash
# Run all YAML configs in a directory with 4 parallel workers
./bin/runner -dir /path/to/configs -workers 4
```

### Configuration File Format

The YAML configuration file defines track parameters and soil layers. The configuration file must specify:

- Track type (ballast or slab track)
- Frequency range for analysis
- Track parameters (specific to track type)
- Soil layer properties
- Output file location

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

See the `configs/sample_config.yaml` for the complete example.

### Results
The output will be a JSON file containing the computed dispersion curves for the specified track type and soil layers, as well as the critical speed. These are the keys in the JSON file:

```json
{
	"omega": [
    ...
    ],
  "track_phase_velocity": [
    ...
    ],
  "soil_phase_velocity": [
    ...
    ],
	"critical_omega": ...,
	"critical_velocity": ...
}
```

## Workflow

A typical workflow using GoTrain involves:

1. Create YAML configuration files for your track and soil conditions
2. For a single analysis:
   ```
   ./bin/critical_speed -config configs/your_config.yaml
   ```

3. For batch processing of multiple configurations:
   ```
   ./bin/runner -dir ./yamls/dir -workers 8
   ```

## License
This project is licensed under the BSD-3-Clause License. See the [LICENSE](LICENSE) file for details.