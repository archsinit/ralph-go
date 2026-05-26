# Cross-platform build

Goal: Produce multi-platform binaries and a build script.

Context:
- Pure-Go deps keep cross-compilation simple.

Do:
1. Add a Makefile or build.sh building for linux/amd64, linux/arm64, darwin/amd64, darwin/arm64 via GOOS/GOARCH.
2. Output to dist/.
3. Document the command in README.

Done when:
- Build script produces all four binaries.
- Each binary runs --help on its platform (at least linux verified here).

Files: Makefile, README.md
