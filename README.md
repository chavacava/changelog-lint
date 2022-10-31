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

## Details

The linter will apply a set of predefined rules (se below).

The expected global format of the file is Markdown where:
* `#` is used for the main title,
* `##` is used as header for versions, 
* `###` is used as header for subsections of versions,
* `*` or `-` is used as item markers for change details entries

```
# Main title (e.g. Changelog)
## Version (e.g. 1.2.3)
### Subsection (e.g. Added)
* Change entry
```

Check this [CHANGELOG.md](CHANGELOG.md) as example of the expected format.

# Rules

| Name | Description | 
| -----| :----: |
| `subsection-empty`| warns on subsections without any entry |
| `subsection-namming`| warns on unknown subsection names (`Added`, `Changed`, `Deprecated`, `Fixed`, `Removed`, `Security` are known) |
| `subsection-order`| warns on subsections not listed alphabetically in a version |
| `subsection-repetition`| warns on subsections appearing more than once under the same version |
| `version-empty`| warns on versions without any subsection |
| `version-order`| warns on versions not well ordered wrt their semver number |
| `version-retpetition`| warns on versions appearing more than once |