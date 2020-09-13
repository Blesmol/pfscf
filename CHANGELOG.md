# Changelog

All notable changes to this project will be documented in this file.

## [Unreleased]

### Added

### Changed

### Removed

### Fixed

## v0.9.2 - 2020-09-13

### Added
- Release now includes a batch file that, when executed via clicking from windows explorer, will open a cmd prompt in the current directory. This should make it easier for windows users to use a command line application like pfscf.
- New subcommand `template search` to search for templates based on search terms.
- New content type `choice` for things like choosing the subtier
- New and improved documentation on usage

### Changed
- Releases for macOS now have "macOS" in their name instead of "Darwin" for clarification
 
## v0.9.1 - 2020-09-11

### Added
- Basic stubs for all PFS2 season 1 scenarios plus quests
- Templates now try to do some basic auto-guessing on possible page margins to reduce the number of cases where values on the produced sheet are misaligned.
- Text is now also automatically shrunk if a textcell is not high enough

### Changed
- Template listing (`pfscf template list`) now shows inheritance relations

## v0.9.0 - 2020-09-10

### Added
- Autoshrink for `textCell`: Automatically reduce font size if text is too wide
- If y2 coordinate is missing or 0, then the cell height is automatically determined via font size
- New template section `parameters`: No longer integrated into content entries. This allows for reusing parameters in a sheet, e.g. GM initials on a Starfinder chronicle
- Content type `rectangle`: Finally drawing colored boxes!
- Content type `trigger`: Can be used to conditionally print other content entries when a specific argument is provided
- Content type `canvas`: Allows to reduce the drawing canvas. Required for easily adapting PFS2 chronicle templates if the right sidebar has different coordinates again
- Chronicle template for Starfinder
- Several chronicle templates for Pathfinder 2 where the right sidebar had a different position in the released sheets

### Changed
- Switched template measurement unit from points to percent
- Empty/missing coordinates are now treated as 0
- Update golang dependencies
- Add missing error check during content generation
- Content entries do no longer have an ID
- Content type `textCell` now also supports static values not related to any passed arguments
- Basically replaced the complete internal data structures and data handling

### Removed
- Content type `societyid`

### Fixed
- Fixed wrong handling of command line arguments in batch mode
- Yaml files with empty sections like `presets:` left the underlying data structure uninitialized due to unexpected behavior in go-yaml module. Workaround was implemented.

## v0.8.0 - 2020-07-30

### Changed
- Heavily restructured internal coding regarding content handling

## v0.7.1 - 2020-07-18

### Changed
- Changed source code structure on disk and goreleaser config

## v0.7.0 - 2020-07-16

### Added
- Documentation on how to write templates
- Batch mode for filling out chronicles

### Changed
- Renamed `x1` and `y1` to `x` and `y`.
- Renamed content `societyid` to `societyId`
- Renamed content `code` to `eventcode`
- Renamed cmd line option `--dummyValues` to `--exampleValues`

### Fixed
- #33: Now correctly only reads files with .yml extension, not with .yml~ or .ymla
- Action `template describe -v` now works again

## v0.6.0 - 2020-07-08

### Added
- New content type `societyid`. This is specifically meant for printing a PFS society id following the pattern `<player_id>-<char_id>`, e.g. 123456-789. This is easier to use than providing both values separately, and also allows better formatting / placement.

### Changed
- Template `pfs2` now provides a `societyid` entry instead of separate `playerid` and `charid` entries. These were removed.

## v0.5.0 - 2020-07-07

### Added
- Template inheritance mechanism
- Mechanism for preset values

### Changed
- Template `pfs2` now uses presets instead of defaults
- Improved error texts

### Removed
- The `default` section is no longer supported / usable

## v0.4.0 - 2020-07-02

### Added
- Align wording: An  `ID` is no longer sometimes called `Name`
- Allow to fill out chronicles with dummy/example values

### Changed
- pfsct is now called pfscf
- Updated pdfcpu from v0.3.3 to v0.3.4
- Use global temp dir now for storing intermediate files
- The `default` section in yaml files was replaced with the section `default` in `presets`

### Fixed
- Now printing filename if an error occurs during reading a yaml file

## v0.3.2 - 2020-06-27

### Added
- Short aliases for cmd line commands, e.g. `f` for `fill`, `t` for `template`
- `verbose` flag and output for `template list` and `template describe`
- Provide example values in verbose output of `template describe`

## v0.3.1 - 2020-06-26

### Added
- First version of `template describe` command

## v0.3.0 - 2020-06-26

### Added
- Proper handling of template files
- First version of `template list` command
- Stubs for other `template` commands

### Fixed
- Allow to execute `pfsct` command when in a different directory

## v0.2.0 - 2020-06-24

### Added
- Check whether required template fields are present
- Mechanism for default values in templates

### Changed
- Configs are now named templates, and thus the `config` subdir was renamed to `template` as well
- Yaml unmarshalling now set to strict

## v0.1 - 2020-06-20

### Added
- First more or less working version that can fill out chronicles for PFS2
