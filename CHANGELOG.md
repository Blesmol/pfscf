# Changelog

All notable changes to this project will be documented in this file.

## [Unreleased]

### Added
- Align wording: An 'ID' is no longer sometimes called 'Name'
- Allow to fill out chronicles with dummy/example values

## Changed
- Updated pdfcpu from v0.3.3 to v0.3.4
- Use global temp dir now for storing intermediate files
- The "default" section in yaml files was replaced with the section "default" in "presets"

### Fixed
- Now printing filename if an error occurs during reading a yaml file

## [0.3.2] - 2020-06-27

### Added
- Short aliases for cmd line commands, e.g. `f` for `fill`, `t` for `template`
- `verbose` flag and output for `template list` and `template describe`
- Provide example values in verbose output of `template describe`

## [0.3.1] - 2020-06-26

### Added
- First version of `template describe` command

## [0.3.0] - 2020-06-26

### Added
- Proper handling of template files
- First version of `template list` command
- Stubs for other `template` commands

### Fixed
- Allow to execute `pfsct` command when in a different directory

## [0.2.0] - 2020-06-24

### Added
- Check whether required template fields are present
- Mechanism for default values in templates

### Changed
- Configs are now named templates, and thus the 'config' subdir was renamed to 'template' as well
- Yaml unmarshalling now set to strict

## [0.1] - 2020-06-20

### Added
- First more or less working version that can fill out chronicles for PFS2
