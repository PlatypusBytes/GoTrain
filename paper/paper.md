---
title: 'GoTrain: High-performance railway critical speed analysis'
tags:
  - Go
  - railways
  - soil dynamics
authors:
  - name: B. Zuada Coelho
    orcid: 0000-0002-9896-3248
    corresponding: true
    affiliation: 1
affiliations:
 - name: Deltares, The Netherlands
   index: 1
date: 18-11-2025
bibliography: references.bib
---


# Summary
GoTrain is an open-source Go library and command-line interface tool for computing the critical train speed of railway infrastructure.
The critical train speed is a key parameter in railway engineering, representing the speed at which dynamic amplification effects arise due to resonance between the train load and the supporting soil.
When a train reaches this speed, a resonance condition develops, leading to large track vibrations that may cause serviceability problems or long-term structural damage.

GoTrain implements a semi-analytical model to evaluate the critical train speed, as proposed by @mezher2016.
GoTrain provides a fast, open-source implementation of this analysis, enabling engineers and researchers to efficiently explore parameter variations and soil–track interactions.


# Statement of need

The estimation of critical train speed is essential for the assessment and design of railway infrastructure, particularly in soft soil conditions where the Rayleigh wave velocity of the soil layers are low.

Typically, these analyses involve numerical modelling by means of finite element or boundary element methods, which are very complex, time-consuming and computationally expensive [@Kacimi2013; @Galvin2010].
This makes these methods impractical to perform network analysis or extensive parameter studies or sensitivity analyses on the critical train speed.

The semi-analytical approach implemented in GoTrain provides a computationally efficient alternative that enables rapid evaluation of critical train speeds. Due to its efficiency, and parallelisation, it allows the performance of stochastic analyses to account for uncertainties in soil properties and track conditions. This leads to better risk assessment and design optimisation of railway track infrastructure.

GoTrain provides a fully open-source (BSD-3), cross-platform, high-performance implementation with a command-line interface and automatic parallelisation. GoTrain is intended for geotechnical engineers, railway infrastructure designers, academic researchers, and consultants working in soft-soil railway environments.


# Model formulation

GoTrain implements semi-analytical models for evaluating track and soil dispersion curves. The critical train speed is defined at the interception of the track and soil dispersion curves phase velocities.

## Track dispersion model

The railway track is modelled as a continuous, infinite Euler-Bernoulli beam subjected to vertical dynamic loading and supported by discrete elements such as sleepers, railpads, and ballast, which rest on a layered soil profile.

Under vertical harmonic excitation in the frequency–wavenumber domain, the governing equation can be written in matrix form as in \autoref{eq:track}:

\begin{equation}
  \mathbf{K}(k, \omega)\, \mathbf{u}(k, \omega) = \mathbf{P}(k, \omega),
  \label{eq:track}
\end{equation}


where $k$ is the wavenumber, $\omega$ is the angular frequency, $\mathbf{u}$ the vertical displacement, $\mathbf{P}$ the vertical force and $\mathbf{K}$ represents the dynamic stiffness matrix of the track system, which is available in either ballast or slab track configurations.
These configurations differ in the way the rail is supported: ballast track uses discrete railpads, sleepers, and a ballast layer, while in slab track the rail is directly supported by discrete railpads on a concrete slab.

The track dispersion relation is obtained by solving \autoref{eq:det_track} for all frequencies $\omega$ of interest.

\begin{equation}
  \det\left|\mathbf{K}(k,\omega)\right| = 0.
  \label{eq:det_track}
\end{equation}


## Soil dispersion model

Surface-wave dispersion in the soil is modelled assuming a horizontally layered,
elastic, isotropic material with constant density and P- and S-wave velocities in
each layer. GoTrain computes Rayleigh wave dispersion using the Fast Delta
Matrix Method [@buchen1996]. This method refines the classical Thomson–Haskell formulation by addressing numerical instabilities that occur at high frequencies, particularly when dealing with thin layers or stiff contrasts. The Rayleigh-wave dispersion function $D$ follows \autoref{eq:soil}:

\begin{equation}
   D(c, \omega) = \det\left|\mathbf{U^{\top}} \mathbf{T} \mathbf{V} \right| = 0,
  \label{eq:soil}
\end{equation}

where:

- $c = \omega / k$ is the Rayleigh-wave phase velocity;

- $\mathbf{T}$ is the full propagator matrix (product of individual layer transfer matrices);

- $\mathbf{U}$, $\mathbf{V}$ are boundary condition matrices (defined by free surface and radiation conditions).

The Fast Delta Matrix Method reformulates this using compound matrices (delta matrices), which eliminates numerical overflow and underflow issues by expressing the system in terms of second-order minors. This ensures stable and efficient computation of dispersion curves across a wide frequency range.


## Model assumptions

GoTrain implements analytical models subject to the following assumptions:

- the track and soil are modelled as uncoupled subsystems;

- soil and track components behave as linear elastic materials;

- soil layers are horizontally infinite and laterally homogeneous;

- nonlinear effects, such as soil plasticity or track degradation, are not considered.


# Functionality and features

GoTrain allows users to compute the critical train speed for railway tracks supported on layered soil profiles.
Key functionalities include:

- railway ballast and slab track configurations;

- layered soils;

- command-line interface for easy execution of analyses and integration with other tools;

- parallel processing capabilities to handle multiple analyses concurrently, ideal for stochastic studies;

- output results in JSON format.


# Usage
GoTrain can be installed from pre-built binaries available at the [GitHub releases page](https://github.com/PlatypusBytes/GoTrain/releases), or it can be built from source following the instructions in the [ReadMe file](https://github.com/PlatypusBytes/GoTrain/blob/main/ReadMe.md).

Two applications are available in GoTrain:

- `critical_speed`: runs one realisation of GoTrain;

- `runner`: runs GoTrain in parallel for multiple realisations (ideal for stochastic analysis).

To run GoTrain the user needs to prepare an input configuration file containing the information about the track and soil properties.
This is done by means of a [YAML format](https://github.com/PlatypusBytes/GoTrain/blob/main/configs/sample_config.yaml).
An example configuration file is shown below:

```YAML
# Track type: can be "ballast" or "slabtrack"
track_type: ballast

# Frequency range configuration
frequency:
  min: 1
  max: 400
  points: 100

# Ballast track parameters (used if track_type is "ballast")
ballast_track:
  EI_rail: 6.4e6          # Rail bending stiffness [N·m^2]
  m_rail: 60.21           # Rail mass per unit length [kg/m]
  k_rail_pad: 6e8         # Railpad stiffness [N/m]
  c_rail_pad: 2.5e5       # Railpad damping [N·s/m]
  m_sleeper: 238.5        # Sleeper (distributed) mass [kg/m]
  E_ballast: 100e6        # Young's modulus of ballast [Pa]
  h_ballast: 0.3          # Ballast (layer) thickness [m]
  width_sleeper: 1.25     # Half-track width [m]
  rho_ballast: 2000       # Ballast density [kg/m^3]
  soil_stiffness: 0.0     # Soil (spring) stiffness [N/m]

# Slab track parameters (used if track_type is "slabtrack")
slab_track:
  EI_rail: 1.29e7         # Rail bending stiffness [N·m^2]
  m_rail: 120             # Rail mass per unit length [kg/m]
  EI_slab: 600e6          # Slab bending stiffness [N·m^2]
  m_slab: 1093.75         # Slab mass per unit length [kg/m]
  k_rail_pad: 5e8         # Railpad stiffness [N/m]
  c_rail_pad: 2.5e5       # Railpad damping [N·s/m]
  soil_stiffness: 0.0     # Soil (spring) stiffness [N/m]

soil_layers:
  - thickness: 5          # Thickness of the soil layer [m]
    density: 1900         # Density of the soil layer [kg/m^3]
    young_modulus: 2.6e7  # Young  modulus of the soil layer [Pa]
    poisson_ratio: 0.30   # Poisson's ratio of the soil layer
  - thickness: 10         # Thickness of the second soil layer [m]
    density: 1900         # Density of the second soil layer [kg/m^3]
    young_modulus: 1.1e8  # Young  modulus of the second soil layer [Pa]
    poisson_ratio: 0.25   # Poisson's ratio of the second soil layer
  - thickness: .inf       # Thickness of the fourth soil layer [m]
    density: 1900         # Density of the fourth soil layer [kg/m^3]
    young_modulus: 4.7e8  # Young  modulus of the fourth soil layer [Pa]
    poisson_ratio: 0.30   # Poisson's ratio of the fourth soil layer

# Output file configuration
output:
  file_name: "dispersion_results.json"
```

The output of GoTrain is a JSON file containing the computed dispersion curves for soil and track and the critical speed.


## Example: `critical_speed`

Analyses a single railway configuration and computes dispersion curves and critical speed. The -config flag specifies the path to the YAML configuration file.

Usage:
```bash
./critical_speed -config <YAML configuration file>
```


## Example: `runner`

Processes multiple YAML configuration files in parallel with configurable worker pools. Automatically discovers all YAML files in the directory tree and processes them concurrently. The -dir flag specifies the directory to search for YAML files, and the -workers flag sets the number of concurrent workers.

Usage:
```bash
./runner -dir <directory tree> -workers <number of workers>
```

## Example output

GoTrain produces a JSON file containing the computed track and soil dispersion curves as well as the critical angular frequency and critical speed. A typical output file is follows:

```json
{
  "omega": [1.0, 2.0, 3.0, ..., 400.0],
  "track_phase_velocity": [8.84, 19.82, 26.60, ..., 160.14],
  "soil_phase_velocity": [280.8, 265.3, 248.09, ..., 67.74],
  "critical_omega": 63.98,
  "critical_speed": 70.54
}
```

# References
