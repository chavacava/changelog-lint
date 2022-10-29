# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/)

## 1.9.0 - 1711-04-20

### Infrastructure
- Update actions/checkout to v3
- @actions/core `1.2.6` -> `1.10.6`
- @actions/github `4.0.0` -> `5.1.1`
- glob `7.1.6` -> `8.0.3`
- wildcard-match `5.1.0` -> `5.1.2`
- @zattoo/eslint-config `11.0.4` -> `14.0.0`
- eslint `7.14.0` -> `8.25.0`
- eslint-config-airbnb-base `14.2.0` -> `15.0.0`
- jest `26.6.3` -> `29.2.0`

## 1.8.0 - 0408-20-22
### Fixed
- [60](https://github.com/zattoo/changelog/issues/60) Supported validation of separator between version and date

### Added
- [65](https://github.com/zattoo/changelog/issues/65) Supported [Yanked](https://keepachangelog.com/en/1.0.0/#yanked) releases


## 1.7.0 - 2021-06-11

### Added
- [51](https://github.com/zattoo/changelog/issues/51) Ignore files
- [37](https://github.com/zattoo/changelog/issues/37) Supported alpha version order and `Unreleased` check as first heading
- [49](https://github.com/zattoo/changelog/issues/49) Limited release branch to ask for modified files only of its own project

### Infrastructure
- Added `v1` back-syncs
- Added tests for releases

## 1.6.0 - 2020-12-01

### Added 
- multiple changelogs support

## 1.6.1 - 2020-12-02

### Fixed 
- errors not throwing

## 1.5.0 - 2020-10-07

### Added
-  Support nested release branches

## 1.4.0 - 2020-09-15

### Added
- Validate breaking changes
- Support version checking for hotfix branch

## 1.3.1 - 2020-08-26

### Fixed
- Validate wrong version headings

## 1.3.0 - 2020-08-26

### Added
- Validate empty lines

## 1.2.0 - 2020-08-23

### Added
- Check if a version has repeated headings

### Removed
- Commit status

### Infrastructure
- Update packages

## 1.2.0

### Added
- Check if a version has repeated headings


## 1.1.0 - 2020-07-30


## 1.0.2 - 2020-07-09


## 1.0.1 - 2020-07-03


## 1.0.0 - 2020-06-29

### Added
- Initial Functionality
