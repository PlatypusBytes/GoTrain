// Command critical_speed is a command-line tool for computing the critical train speed
// on railway tracks using dispersion analysis. It supports both ballast and slab track
// models and reads physical parameters from a YAML configuration file.
//
// The tool performs the following steps:
//   - Parses a YAML configuration file describing track and soil parameters.
//   - Computes the dispersion curves of the railway track.
//   - Computes the dispersion curves of the soil layers.
//   - Identifies the critical speed where the track and soil phase velocities intersect.
//   - Outputs the results (omega, phase velocities, critical values) to a structured JSON file.
//
// Usage:
//
//	go run cmd/critical_speed/main.go -config path/to/config.yaml
//
// Or using the compiled binary:
//
//	./bin/critical_speed -config path/to/config.yaml
//
// Required flags:
//
//	-config string
//	 	Path to the YAML configuration file defining model parameters.
//
// Configuration:
// The YAML file must specify the track type ("ballast" or "slabtrack"), the frequency range,
// track structure parameters, soil layer properties, and output file destination.
//
// For a complete example configuration file, see:
//
//	./configs/sample_config.yaml
package critical_speed
