![Tests](https://github.com/PlatypusBytes/GoTrain/actions/workflows/go.yaml/badge.svg)
[![Go Reference](https://pkg.go.dev/badge/github.com/PlatypusBytes/GoTrain.svg)](https://pkg.go.dev/github.com/PlatypusBytes/GoTrain)
[![Go Report Card](https://goreportcard.com/badge/github.com/PlatypusBytes/GoTrain)](https://goreportcard.com/report/github.com/PlatypusBytes/GoTrain)


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
  min: 0.1
  max: 250
  points: 100

# Ballast track parameters
ballast_track:
  ei_rail: 1.29e7        # Rail bending stiffness [N·m^2]
  m_rail: 120            # Rail mass per unit length [kg/m]
  k_rail_pad: 5e8        # Railpad stiffness [N/m]
  # ... other parameters ...

# Slab track parameters
slab_track:
  ei_rail: 1.29e7        # Rail bending stiffness [N·m^2]
  m_rail: 120            # Rail mass per unit length [kg/m]
  # ... other parameters ...

# Output file configuration
output:
  file_name: "dispersion_results.json"
```

See the `configs/sample_config.yaml` for a complete example.


## License

This project is licensed under the BSD-3-Clause License. See the [LICENSE](LICENSE) file for details.