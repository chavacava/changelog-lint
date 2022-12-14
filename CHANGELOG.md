# Changelog

## 0.3.0 - 2022/11/04

### Added
* `release` mode command line flag to enable release-related checks. Example, when releasing version 1.2.3:
```bash
changelog-lint -release 1.2.3
```

### Fixed
* Rule disabling not working
* Typo in rule name: `subsection-namming`

## 0.2.0 - 2022/11/03

### Added
- Rules configuration support
- Parser configuration support

### Changed
- `version-empty` rule now accepts `Unreleased` version to be empty

### Fixed
- Multiline entries are not accepted by the default parser
- `version-order` rule fails to compare `Unreleased` version wrt other versions

## 0.1.0 - 2022/10/31

### Added
- Default changelog parser for `.md` format
- Rule `subsection-empty`
- Rule `subsection-naming`
- Rule `subsection-order`
- Rule `subsection-repetition`
- Rule `version-empty`
- Rule `version-order`
- Rule `version-repetition`
