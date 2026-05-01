# Prime CLI

## Overview

The Coinbase Prime command-line interface (CLI) to simplify programmatic interaction with [Coinbase Prime's](https://prime.coinbase.com/) [REST APIs](https://docs.cloud.coinbase.com/prime/reference).

## License

The Prime CLI is free and open source and released under the [Apache License, Version 2.0](LICENSE.txt).

The application and code are only available for demonstration purposes.

## Installation

Install the Prime CLI binary on your computer.

### MacOS

Ensure that you have [Homebrew](https://brew.sh/) installed.

Install the Coinbase Samples Homebrew tap:

```
brew tap coinbase-samples/homebrew-tap
```

Next, install the Prime CLI Homebrew formula:

```
brew install prime-cli
```

### Other Platforms

To install the Prime CLI on other platforms, clone this repository, build, and then add the binary to the [path](https://en.wikipedia.org/wiki/PATH_(variable)).

```
git clone git@github.com:coinbase-samples/prime-cli.git
cd prime-cli
go build -o primectl
```

## Configuration

Once you have the CLI installed, configure your environment.

Set an environment variable in your shell called `PRIME_CREDENTIALS` with your API and portfolio information.

Coinbase Prime API credentials can be created in the Prime web console under Settings -> APIs. Entity ID can be retrieved by calling [Get Portfolio](https://docs.cloud.coinbase.com/prime/reference/primerestapi_getportfolio). If you are not configured yet to call this endpoint, you may proceed without including the `entityId` key and value for the time being, but certain endpoints such as List Invoices and List Assets require it.

`PRIME_CREDENTIALS` should match the following format:
```
export PRIME_CREDENTIALS='{
"accessKey":"ACCESSKEY_HERE",
"passphrase":"PASSPHRASE_HERE",
"signingKey":"SIGNINGKEY_HERE",
"portfolioId":"PORTFOLIOID_HERE",
"svcAccountId":"SVCACCOUNTID_HERE",
"entityId":"ENTITYID_HERE"
}'
```

You may also pass an environment variable called `primeCliTimeout` which will override the default request timeout of 7 seconds. This value should be an integer in seconds.

## Usage

Build the application binary and specify an output name, e.g. `primectl`:

```
go build -o primectl
```

To ensure your project's dependencies are up-to-date, run:
```
go mod tidy
```

To verify that your application is installed correctly and accessible from any location, run the following command. It will include all available requests:

```
./primectl
```

Finally, to run commands for each endpoint, use the following format to test each endpoint. Please note that many endpoints require flags, which are detailed with the `--help` flag.

```
./primectl portfolios list
```

```
./primectl orders create --help
```

```
./primectl orders create-preview -b 0.001 -i ETH-USD -s BUY -t MARKET
```

As of v0.4.0, the CLI covers the full surface area of [prime-sdk-go](https://github.com/coinbase-samples/prime-sdk-go) v0.6.2, including the `advanced-transfers`, `futures`, and `positions` command groups.

## Releasing

The Prime CLI is distributed via the [`coinbase-samples/homebrew-tap`](https://github.com/coinbase-samples/homebrew-tap) Homebrew tap. Cutting a new release is a two-repo process: tag the source here, then bump the formula in the tap.

### 1. Bump the version in this repo

1. Update the version string in `cmd/version.go` (`primectlVersion`) to the new semver, e.g. `0.5.0`.
2. Add a new entry at the top of `CHANGELOG.md` following the existing `Added` / `Fixed` format.
3. If user-facing commands changed, refresh `README.md` and `COMMANDS.md`.
4. Open a PR, get it reviewed, and merge to `main`.

### 2. Tag and create a GitHub release

From `main` after the bump is merged:

```
git checkout main && git pull
git tag vX.Y.Z
git push origin vX.Y.Z

gh release create vX.Y.Z \
  --repo coinbase-samples/prime-cli \
  --title "vX.Y.Z" \
  --notes "See CHANGELOG.md for details."
```

Tags must be prefixed with `v` (e.g. `v0.5.0`) — the Homebrew formula's `url` resolves the tag by that exact name.

### 3. Compute the source tarball SHA256

Homebrew pins the GitHub source tarball by hash:

```
VERSION=X.Y.Z
curl -sL "https://github.com/coinbase-samples/prime-cli/archive/refs/tags/v${VERSION}.tar.gz" \
  | shasum -a 256
```

Copy the resulting 64-character hex digest.

### 4. Update the Homebrew formula

In a clone of [`coinbase-samples/homebrew-tap`](https://github.com/coinbase-samples/homebrew-tap), edit `Formula/prime-cli.rb` and update three lines:

```
url "https://github.com/coinbase-samples/prime-cli/archive/refs/tags/vX.Y.Z.tar.gz"
sha256 "<sha256 from the previous step>"
version "X.Y.Z"
```

Leave the `depends_on`, `install`, and `test` blocks untouched — the formula builds from source with `go build`.

### 5. Validate the formula locally

From the tap clone:

```
brew install --build-from-source ./Formula/prime-cli.rb
brew test prime-cli
brew audit --strict --online prime-cli
primectl version   # should print {"version":"X.Y.Z"}
```

`brew audit` will catch a mismatched `sha256`, a broken URL, or formula style issues before review.

### 6. Open the PR against the tap

```
gh pr create \
  --repo coinbase-samples/homebrew-tap \
  --title "prime-cli X.Y.Z" \
  --body "Bumps prime-cli formula to vX.Y.Z. See https://github.com/coinbase-samples/prime-cli/releases/tag/vX.Y.Z"
```

Once merged to `main` on the tap, end users pick up the new version with:

```
brew update
brew upgrade prime-cli
```

