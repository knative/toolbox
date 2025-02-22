# buoy

A tool to encode the Knative release process logic based on go module versions.

## Installation

`buoy` can be installed and upgraded by running:

```shell
go get knative.dev/toolbox/buoy
```

## Usage

```
Usage:
  buoy [command]

Available Commands:
  actions     Interact with GitHub Actions.
  check       Determine if this module has release branches or releases available from each dependency for a given release.
  float       Find latest versions of dependencies based on a release.
  help        Help about any command
  needs       Find dependencies based on a base import domain.
  exists      Determine if the release branch exists for a given module.
  repos       List the repos for a list of GitHub organizations.

Flags:
  -h, --help   help for buoy

Use "buoy [command] --help" for more information about a command.
```

### Actions

```
Interact with GitHub Actions.

Usage:
  buoy actions [command]

Available Commands:
  list        List GitHub Actions workflows for a given repository.
  run         Run a GitHub Actions workflow for a given repository.

Flags:
  -h, --help   help for actions

Use "buoy actions [command] --help" for more information about a command.
```

#### Actions List

```
List GitHub Actions workflows for a given repository.

Usage:
  buoy actions list org/repo [flags]

Flags:
  -h, --help                help for list
  -w, --only-workflow-id    Only output the workflow ID.
  -q, --query string        Search for a workflow by name.
      --short               Short output will only print the workflow url.
  -t, --token-path string   GitHub token file path.
```

#### Actions Run

```
Run a GitHub Actions workflow for a given repository.

Usage:
  buoy actions run org/repo --query OneResult [flags]

Flags:
  -h, --help                help for run
      --id int              Workflow ID.
      --inputs string       Workflow inputs.
  -q, --query string        Search for a workflow by name.
      --ref string          Ref to run workflow from. (default "master")
  -t, --token-path string   GitHub token file path.
```

### Check

```
The check command is used to evaluate if each dependency for the given module
meets the requirements for cutting a release branch. If the requirements are
met based on the ruleset selected, the command will exit with code 0, otherwise
an error message is generated and the with the failed dependencies and exit
code 1. Errors are written to stderr. Verbose output is written to stdout.

Rulesets,
  Release          check requires all dependencies to have tagged releases.
  Branch           check requires all dependencies to have a release branch.
  ReleaseOrBranch  check will use rule (Release || Branch).

Usage:
  buoy check go.mod [flags]

Flags:
  -d, --domain string    domain filter [required]
  -h, --help             help for check
  -r, --release string   release should be '<major>.<minor>' (i.e.: 1.23 or v1.23) [required]
      --ruleset string   The ruleset to evaluate the dependency refs. Rulesets: [Any, ReleaseOrBranch, Release, Branch] (default "ReleaseOrBranch")
  -v, --verbose          Print verbose output.
```

Example,

```
$ go run . check $HOME/go/src/knative.dev/eventing-github/go.mod --release 0.18 --verbose
[exit status 0]

$ go run . check $HOME/go/src/knative.dev/eventing-github/go.mod --release 0.19 --verbose
knative.dev/eventing-github not ready for release because of the following dependencies [knative.dev/eventing@master, knative.dev/pkg@master, knative.dev/serving@master, knative.dev/test-infra@master]
[exit status 1]
```

If you need to see a more verbose output, use `--verbose`:

```
$ buoy check $HOME/go/src/knative.dev/eventing-github/go.mod --domain knative.dev --release 0.18 --verbose
knative.dev/eventing-github
✔ knative.dev/eventing@v0.18.0
✔ knative.dev/pkg@release-0.18
✔ knative.dev/serving@v0.18.0
✔ knative.dev/test-infra@release-0.18
[exit status 0]

$ buoy check $HOME/go/src/knative.dev/eventing-github/go.mod --domain knative.dev  --release 0.19 --verbose
knative.dev/eventing-github
✘ knative.dev/eventing@master
✘ knative.dev/pkg@master
✘ knative.dev/serving@master
✘ knative.dev/test-infra@master
knative.dev/eventing-github not ready for release because of the following dependencies [knative.dev/eventing@master, knative.dev/pkg@master, knative.dev/serving@master, knative.dev/test-infra@master]
[exit status 1]
```

### Float

```
The goal of the float command is to find the best reference for a given release.
Float will select a ref for found dependencies, in this order (for the Any
ruleset, default):

1. A release tag with matching major and minor; choosing the one with the
   highest patch version, ex: "v0.1.2"
2. If no tags, choose the release branch, ex: "release-0.1"
3. Finally, the default branch, ex: "master"

The selection process for float can be modified by providing a ruleset.

Rulesets,
  Any              tagged releases, release branches, default branch
  Release          tagged releases
  Branch           release branches
  ReleaseOrBranch  tagged releases, release branch

For rulesets that that restrict the selection process, no ref is selected.

NOTE: For Go modules that don't live on the root of their repository, the tagged
releases are skipped. Golang requires the tags for sub-modules to be in form of
"sub/v1.99.0" to be recognized by "go get" command. It's cumbersome to have
both "v1.99.0" and "sub/v1.99.0", so float will skip looking for tags in this
case.

Usage:
  buoy float go.mod [flags]

Flags:
  -d, --domain string    domain filter [required]
  -h, --help             help for float
  -r, --release string   release should be '<major>.<minor>' (i.e.: 1.23 or v1.23) [required]
      --ruleset string   The ruleset to evaluate the dependency refs. Rulesets: [Any, ReleaseOrBranch, Release, Branch] (default "Any")
```

Example:

```
$ buoy float go.mod --domain knative.dev --release v0.15
knative.dev/eventing@v0.15.4
knative.dev/pkg@release-0.15
knative.dev/serving@v0.15.3
knative.dev/test-infra@release-0.15
```

Or set `domain` to and target release of that dependency:

```shell script
$ buoy float go.mod --release 0.18 --domain k8s.io
k8s.io/api@v0.18.10
k8s.io/apiextensions-apiserver@v0.18.10
k8s.io/apimachinery@v0.18.10
k8s.io/client-go@v0.18.10
k8s.io/code-generator@v0.18.10
k8s.io/gengo@master
k8s.io/klog@master
```

Note: the following are equivalent releases:

- `v0.1`
- `v0.1.0`
- `0.1`
- `0.1.0`

### Needs

```
Find dependencies based on a base import domain.

Usage:
  buoy needs go.mod [flags]

Flags:
  -d, --domain string   domain filter [required]
  -h, --help            help for needs
```

Example,

```
$ buoy needs $HOME/go/src/knative.dev/eventing-github/go.mod --domain knative.dev
knative.dev/eventing
knative.dev/pkg
knative.dev/serving
knative.dev/test-infra
```

Or `--domain` set to `k8s.io`:

```
$ buoy needs $HOME/go/src/knative.dev/eventing-github/go.mod --domain k8s.io
k8s.io/api
k8s.io/apimachinery
k8s.io/client-go
```

### Exists

```
Determine if the release branch exists for a given module.

Usage:
  buoy exists go.mod [flags]

Flags:
  -h, --help             help for next
  -r, --release string   release should be '<major>.<minor>' (i.e.: 1.23 or v1.23) [required]
  -t, --next             Print the next release tag (stdout)
  -v, --verbose          Print verbose output (stderr)
```

### Repos

```
List the repos for a list of GitHub organizations.

Usage:
  buoy repos org1 [org2 org3...] [flags]

Flags:
  -h, --help                help for repos
  -t, --token-path string   GitHub token file path.
```

## TODO:

- Support `go-import` with more than one import on a single page.
- Support release branch templates. For now, hardcoded to Knative style.
