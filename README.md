# `changelog-lint`

A simple linter for changelog files.

## Installation

```bash
go install github.com/chavacava/changelog-lint@master
```

## Usage

```
changelog-lint
```
will lint the `CHANGELOG.md` file in the current directory

```
changelog-lint some/path/changes.md
```
will lint the `changes.md` file in the `some/path` directory

To get the full list of command line flags:
```
changelog-lint -h
```

### Configuration

If defaults do not fit your needs you can configure the linter by providing a configuration file with the command line flag `-config`

Via the configuration file you can:
- Disable rules (all enabled by default)
- Provide arguments to rules

The format of the configuration file is TOML.
Example:
```toml
[rule.version-empty]
    Disabled=true
[rule.subsection-namming]
    Arguments=["Performance", "Refactoring"]
```

### Error codes
Executing the linter returns one of the following error codes

| Code | Meaning | 
| -----:| :---- |
|`0`| no error|
|`1`| bad execution parameters/flags (e.g. bad changelog filename)|
|`2`| syntax error in the changelog file|
|`3`| the linting found a problem in the changelog|
 
## Details

The linter will apply a set of predefined rules (se below).

The expected global format of the file is Markdown where:
* `#` is used for the main title,
* `##` is used as header for versions, 
* `###` is used as header for subsections of versions,
* `*` or `-` is used as item markers for change details entries

By default the following patterns are expected

| Section | Pattern | 
| -----:| :---- |
| Title | `^# .+$` |
| Version | `^## (\d+\.\d+.\d+\|\[Unreleased\])( .*)*$` |
| Subsection | `^### ([A-Z]+[a-z]+)[ ]*$` |
| Entry | `^[*-] .+$` |


Check this [CHANGELOG.md](CHANGELOG.md) as example of the expected format.

# Rules

| Name | Description | Arguments |
| -----:| :---- | :---- |
| `subsection-empty`| warns on subsections without any entry ||
| `subsection-namming`| warns on unknown subsection names (`Added`, `Changed`, `Deprecated`, `Fixed`, `Removed`, `Security` are known) | list of accepted section names (replaces the default list) |
| `subsection-order`| warns on subsections not listed alphabetically in a version |
| `subsection-repetition`| warns on subsections appearing more than once under the same version |
| `version-empty`| warns on versions without any subsection |
| `version-order`| warns on versions not well ordered wrt their semver number |
| `version-retpetition`| warns on versions appearing more than once |

You can contribute new rules by implementing the `linting.Rule` interface:

```go
type Rule interface {
	Apply(model.Changelog, chan Failure, RuleArgs)
	Name() string
}
```

## TODO

* Accept parser configuration
* Relase mode:
    - Check last version corresponds to a given one, then
    - `[Unreleased]` not allowed
