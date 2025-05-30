![Tests](https://github.com/PlatypusBytes/GoTrain/actions/workflows/go.yaml/badge.svg)

# GoTrain

A Go library for analyzing critical speeds in railway systems, focusing on soil and track dispersion analysis.

The methodology for the computation of the critical train speed is based on the work of [Mezher et al. (2016)](https://www.sciencedirect.com/science/article/abs/pii/S2214391215000239).
The dispersion curve for the layered soil is based on the Fast Delta Matrix method proposed by [Buchen and Ben-Hador (1996)](https://academic.oup.com/gji/article-lookup/doi/10.1111/j.1365-246X.1996.tb05642.x).



## Building

To build the critical speed utility:

```bash
go build -o bin/criticalspeed ./cmd/critical_speed
```

Or use the Makefile:

```bash
make build
```

## Running

Execute the critical speed analysis tool:

```bash
./bin/criticalspeed
```

## License

[Specify your license here]

## Last Updated

May 30, 2025
