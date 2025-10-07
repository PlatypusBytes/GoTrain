# Make a release

To make a release we need to create a tag.
The tag name should follow semantic versioning, vMAJOR.MINOR.PATCH, where:


 - MAJOR: Incremented for incompatible API changes.
 - MINOR: Incremented for added functionality in a backwards-compatible manner.
 - PATCH: Incremented for backwards-compatible bug fixes.

To create a new tag, run the following commands:

```bash
git tag v1.x.y
git push origin v1.x.y
```
The workflow [build-on-tag.yaml](.github/workflows/build-on-tag.yaml) will be triggered by this tag, and the binaries will be built and uploaded to the releases page.

# Build and format code
To build and format the code, there is a [Makefile](Makefile) with the following commands:

```bash
make fmt
make tidy
make build
```
- `make fmt`: Formats the code using `gofmt`.
- `make tidy`: Cleans up the `go.mod` and `go.sum` files.
- `make build`: Builds the project and creates the binary in the `bin/` directory.
- `make test`: Runs the tests.



# Add the package to pkg.go.dev

To trigger a new documentation build on pkg.go.dev, you can use the following command:

```bash
go get github.com/PlatypusBytes/GoTrain@v1.1.1
```

Make sure to replace `v1.1.1` with the desired version tag.