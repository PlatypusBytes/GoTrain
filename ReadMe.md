![Tests](https://github.com/PlatypusBytes/GoTrain/actions/workflows/go.yaml/badge.svg)
[![Go Reference](https://pkg.go.dev/badge/github.com/PlatypusBytes/GoTrain.svg)](https://pkg.go.dev/github.com/PlatypusBytes/GoTrain)
[![codecov](https://codecov.io/gh/PlatypusBytes/GoTrain/graph/badge.svg)](https://codecov.io/gh/PlatypusBytes/GoTrain)[![Go Report Card](https://goreportcard.com/badge/github.com/PlatypusBytes/GoTrain)](https://goreportcard.com/report/github.com/PlatypusBytes/GoTrain)


# GoTrain

A Go library for analyzing critical speeds in railway systems, focusing on soil and track dispersion analysis.

The methodology for the computation of the critical train speed is based on the work of [Mezher et al. (2016)](https://www.sciencedirect.com/science/article/abs/pii/S2214391215000239).
The dispersion curve for the layered soil is based on the Fast Delta Matrix method proposed by [Buchen and Ben-Hador (1996)](https://academic.oup.com/gji/article-lookup/doi/10.1111/j.1365-246X.1996.tb05642.x). (WIP)


## Building

To build the critical speed library:

```bash
go build -o bin/critical_speed ./cmd/critical_speed
```

Or use the Makefile:

```bash
make build
```

## Usage

The critical speed utility requires a YAML configuration file for specifying track parameters:

```bash
# Run with a configuration file (required)
./bin/critical_speed -config /path/to/config.yaml
```

### Configuration File Format

Create a YAML file with the following structure:

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
  EI_rail: 1.29e7        # Rail bending stiffness [N·m^2]
  m_rail: 120            # Rail mass per unit length [kg/m]
  k_rail_pad: 5e8        # Railpad stiffness [N/m]
  c_rail_pad: 2.5e5      # Railpad damping [N·s/m]
  m_sleeper: 490         # Sleeper (distributed) mass [kg/m]
  E_ballast: 130e6       # Young's modulus of ballast [Pa]
  h_ballast: 0.35        # Ballast (layer) thickness [m]
  width_sleeper: 1.25    # Half-track width [m]
  rho_ballast: 1700      # Ballast density [kg/m^3]
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
    young_modulus: 50.66666e6    # Young  modulus of the soil layer [Pa]
    poisson_ratio: 0.33333333    # Poisson's ratio of the soil layer
  - thickness: 10         # Thickness of the second soil layer [m]
    density: 1900         # Density of the second soil layer [kg/m^3]
    young_modulus: 202.6666e6  # Young  modulus of the second soil layer [Pa]
    poisson_ratio: 0.33333333    # Poisson's ratio of the second soil layer
  - thickness: 15         # Thickness of the third soil layer [m]
    density: 1900         # Density of the third soil layer [kg/m^3]
    young_modulus: 456e6    # Young  modulus of the third soil layer [Pa]
    poisson_ratio: 0.3    # Poisson's ratio of the third soil layer
  - thickness: .inf       # Thickness of the fourth soil layer [m]
    density: 1900         # Density of the fourth soil layer [kg/m^3]
    young_modulus: 810.6666e6  # Young  modulus of the fourth soil layer [Pa]
    poisson_ratio: 0.33333333    # Poisson's ratio of the fourth soil layer

# Output file configuration
output:
  file_name: "dispersion_results.json"
```

See the `configs/sample_config.yaml` for the complete example.

### Results
The output will be a JSON file containing the computed dispersion curves for the specified track type and soil layers, as well as the critical speed. These are the keys in the JSON file:

```
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

## License
This project is licensed under the BSD-3-Clause License. See the [LICENSE](LICENSE) file for details.