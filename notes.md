# Release Notes

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

