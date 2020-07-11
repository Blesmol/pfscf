# Creating Template Files for Pfscf

After some time of using pfscf, you might be tempted to create your own templates or modify existing ones.
Or you are just curious on what is currently supported by pfscf.
Besides, it's always useful to provide proper documentation.
Even if perhaps only for the reason to show that the weird behavior for something might be actually intended!

## Template Structure on Disk

The template files are stored in a folder `templates` that is located in the same folder as the pfscf executable file.
Within this folder, each template is stored in a single file with file extension `.yml`.
It is also possible to create subfolder within the `templates` folder to organize the template files.
For example, this could be used to have separate folders for each supported game system, e.g. one folder for Pathfinder, one for Starfinder.

The filenames of the template files are irrelevant.
They should provide a hint to what is included, but besides from that they are not used anywhere.
This especially means that they are not used for the unique identifier which each template must have; this is done by the `id` field within the template file.
See also the `Template file format` section below on this.

## Template File Format

### YAML Format

Template files are stored in YAML format, one template per file.
The official spec of this format is maintained at [yaml.org](https://yaml.org/), and if you search the web you will find lots of examples, introductions and explanations to the format.
In this document, however, I will try to spare you the irrelevant parts and only include a very basic overview over what you might need.

A YAML file is a plain text file and can be opened with any text editor, if need be with the Windows Notepad application.
What is not possible is to use MS Word, LibreOffice and the like to modify such files.
Many modern text editors even bringt support for writing/modifying YAML files.
One such editor that I can recommend is [Notepad++](https://notepad-plus-plus.org) (not to be mixed up with the already mentioned MS Notepad).

YAML is intended to be a human-readable format.
Most lines in these files follow the format `<some identifier>: <some value>`, e.g. `description: This is a template for PFS`.
And until I have more time to write a proper YAML introduction, I would suggest you take a look at the existing files in the `templates` folder or search the web.

### File Layout

Template files have the following structure:
```yaml
id: <unique template identifier>
description: <short description on what this template contains>
inherit: <id of template that should be inherited>

presets:
  <presetId>:
    <content>
  <presetId>:
    <content>
  ...

content:
  <contentId>:
    <content>
  <contentId>:
    <content>
  ...
```

<details>
  <summary>Example</summary>

```yaml
id: myId
description: This is an example template
inherit: idOfMyOtherCoolTemplate

presets:
  topline:
    y:  100
    y2: 200
  bottomline:
    y:  400
    y2: 500

content:
  name:
    type: textCell
    presets: [ topline ]
    x: 50
    x2: 150
  date:
    type: textCell
    presets: [ bottomline ]
    x: 200
    x2: 230
```
</details>

The only mandatory top-level field in such a template is the `id`, all other top-level fields (`description`, `inherit`, `presets`, `content`) in this basic structure are optional.
Of course, a template that only consists of an `id`  does not make much sense.
But who am I to judge?

Each entry below `content` and `presets` must have an ID that is unique within that section.
That means, no two entries in the `content` section may have the same ID, and no two entries in the `presets` section may have the same ID.
The description on what may be included in a `content` or `presets` entry is described below.

### Content Types

The pfsct app supports different types of content entries.
Each content entry requires a mandatory field `type` where the type of the content entry must be added.

#### Type `textCell`

A `textCell` describes a rectangular cell on the PDF file where user-provided text is added.

It has a couple of mandatory and some optional fields:
| Field      | Description                                                                                             | Input type | Required? |
|:-----------|:--------------------------------------------------------------------------------------------------------|:----------:|:---------:|
| `desc`     | Gives a short description of what this is                                                               | Text       | TODO      |
| `x`        | First coordinate for the cell on the X axis                                                             | Number     | Mandatory |
| `y`        | First coordinate for the cell on the Y axis                                                             | Number     | Mandatory |
| `x2`       | Second coordinate for the cell on the X axis                                                            | Number     | Mandatory |
| `y2`       | Second coordinate for the cell on the Y axis                                                            | Number     | Mandatory |
| `font`     | Name of the font to use. See [the list of supported fonts](#fonts)                                      | Text       | Mandatory |
| `fontsize` | Fontsize in points                                                                                      | Number     | Mandatory |
| `align`    | Text alignment inside the rectangle. See [the list of supported alignments](#text-alignment).           | Text       | TODO      |
| `example`  | An example input value                                                                                  | Text       | Optional  |
| `presets`  | A list of preset IDs to apply to this content entry. See the [section about presets](#preset-mechanism) | List       | Optional  |


Regarding the coordinates: It does not matter whether `x` is smaller than `x2` or vice versa.
Same goes for `y` and `y2`.

<details>
  <summary>TextCell Example</summary>

```yaml
playername:
  type: textCell
  desc: Player name
  x:  40
  y:  125
  x2: 140
  y2: 110
  font: Helvetica
  fontsize: 14
  align: CB
  example: Bob
```
</details>

#### Type `societyId`

A `societyId` is a special type of content whose sole purpose is to bring a society ID to a chronicle with some special formatting.
A society ID should follow the pattern `<player_id>-<char_id>`, e.g. 123456-789.

So why not simply use one or multiple `textCell` entries for this?
Well, the `societyId` does a little bit more than that.
One special thing is that it blanks out (some of) the background underneath it.
For PFS2, there is already the "- 2" part of the society ID preprinted on the chronicles.
I can understand the motivation behind that, but for the automatically filled out chronicles I want to have everything using the same font and size.

Second reason for an own content type is the positioning.
With a text cell you always have to make assumptions on where exactly to put the text.
Should it be left-aligned?
But what if the player id is longer for some players as it has more digits?
Long story short: The `societyId` type takes one extra parameter `xpivot` that basically says on where exactly the dash should be located on the x axis.
Left of that we have the player id, right of that we have the char id.

| Field      | Description                                                                                             | Input type | Required? |
|:-----------|:--------------------------------------------------------------------------------------------------------|:----------:|:---------:|
| `desc`     | Gives a short description of what this is                                                               | Text       | TODO      |
| `x`        | First coordinate for the cell on the X axis                                                             | Number     | Mandatory |
| `y`        | First coordinate for the cell on the Y axis                                                             | Number     | Mandatory |
| `x2`       | Second coordinate for the cell on the X axis                                                            | Number     | Mandatory |
| `y2`       | Second coordinate for the cell on the Y axis                                                            | Number     | Mandatory |
| `xpivot`   | Location of the (center of the) dash on the x axis. Must lie between `x` and `x2`                       | Number     | Mandatory |
| `font`     | Name of the font to use. See [the list of supported fonts](#fonts)                                      | Text       | Mandatory |
| `fontsize` | Fontsize in points                                                                                      | Number     | Mandatory |
| `example`  | An example input value                                                                                  | Text       | Optional  |
| `presets`  | A list of preset IDs to apply to this content entry. See the [section about presets](#preset-mechanism) | List       | Optional  |

<details>
  <summary>SocietyId Example</summary>

```yaml
societyid:
  type: societyId
  desc: The players society id
  x:  40
  y:  125
  xpivot: 100
  x2: 140
  y2: 110
  font: Helvetica
  fontsize: 14
  example: 123456-789
```
</details>


### Other Formating Options

#### Text Alignment

Some content types, e.g. `textCell`, allow to choose an alignment.
This normally consists of a horizontal and a vertical alignment.
For example, selecting an alignment of `RT` for a `textCell` would indicate that the text should be aligned in the **top** **right** corner of the cell.
The possible values for horizontal and vertical alignment can be found below.
If you choose more than one horizontal or vertical alignment, no error will be thrown yet, but only one of the chosen alignments will be used.
The order of the alignment values does not matter, e.g. both `RT` and `TR` will have the same result.

##### Horizontal Alignment

* `L`: Left-bound
* `C`: Centered
* `R`: Right-bound

##### Vertical Alignment

* `T`: Top
* `M`: Middle
* `B`: Bottom
* `A`: Baseline

#### Fonts

The PDF format has builtin support for 14 fonts, the standard PDF fonts.
But of course normally you see PDFs out ther that use a plethora of other fonts.
So if someone uses some font for a PDF file that is not included in the standard fonts, the safest way to ensure that whoever opens a PDF file is not greeted with some nasty error message about missing fonts is to include that font in the resulting PDF.
However, that is currently not supported by `pfscf` and perhaps never will be.
So at the moment only the set of standard PDF fonts is officially supported.
If you use something else, you're doing it at your own risk.
It might work, or it might not.
At the moment there is no check included to filter out wrong or unsupported font names.

The list of fonts that can be used is as follows:
* Times-Roman
* Times-Bold
* Times-Italic
* Times-BoldItalic
* Courier
* Courier-Bold
* Courier-Oblique
* Courier-BoldOblique
* Helvetica
* Helvetica-Bold
* Helvetica-Oblique
* Helvetica-BoldOblique
* Symbol
* ZapfDingbats

### Preset Mechanism

## Template Inheritance

Templates can inherit other templates.
This means that all the content and presets from the original template will also be available in the new template.
The inheritance mechanism works tranistively.
So if template A inherits from template B, and template B inherits from template C, then all presets and content from both templates B and C will be available in template A.

In case of conflicting IDs, i.e. when an ID for a preset or content entry appears in both template, the entry from the inheriting template takes precedence.

Presets are only resolved after everything was inherited.
This means you can also change the appearance of inherited content by replacing presets that are used by this inherited content.

<details>
  <summary>Inheritance Example</summary>

File 1 (will be inherited by file 2 below):
```yaml
id: foo
presets:
  defaultfont:
    font: Helvetica
    fontsize: 14

  topline:
    y: 100
    y2: 200
    presets: [ defaultfont ]

content:
  charname:
    type: textCell
    x:  100
    x2: 200
    presets: [ topline ]

  xp:
    type: textCell
    x: 300
    y: 300
```

File 2 (inherits file 1 from above):
```yaml
id: foobar
inherit: foo
presets:
  defaultfont:
    font: Arial
    fontsize: 10

content:
  playername:
    type: textCell
    x:  300
    x2: 400
    presets: [ topline ]

  xp:
    type: textCell
    x: 450
    y: 450
```

Result after inheritance was resolved:
```yaml
id: foobar
presets:
  defaultfont:  # <= this comes from file 2
    font: Arial
    fontsize: 10

  topline:      # <= this comes from file 1
    y: 100
    y2: 200
    presets: [ defaultfont ]

content:
  charname:     # <= this comes from file 1
    type: textCell
    x:  100
    x2: 200
    presets: [ topline ]

  playername:   # <= this comes from file 2
    type: textCell
    x:  300
    x2: 400
    presets: [ topline ]

  xp:           # <= this comes from file 2
    x: 450
    y: 450
```

</details>

## Finding the correct coordinates

Let's be honest from the beginning: Finding the correct coordinates for adding own content is always fiddly and a lot of try-and-error.
However, there are a few ways to make life easier here.
There are currently three options for the `pfscf fill` command to support finding the correct coordinates:
* `--cellBorder` (short: `-c`)
* `--dummyValues` (short: `-d`)
* `--grid` (short: `-g`)

The `grid` option is probably the most useful.
It will draw a grid over the chronicle PDF file and print the x/y coordinates at the borders of the page.
The grid consists of major lines every 25 points and minor lines every 5 points.
This should allow to get at least the rough coordinates for whatever you want to add, although normally some fine-tuning is required afterwards.

Next is the `cellBorder` option.
When this is selected, `pfscf` will draw borders around the cells that it is printing on the page.
This allows to see the exact boundaries and locations of what you provided via coordinates.
At the moment, this is only done for content that is actually printed to the page.
So if you have some new content, but do not print it on the page, there won't be any borders displayed as well.

And finally the `dummyValues` option.
If you select this, normal input values like, e.g. `player=Bob` will be ignored.
Instead, every content that has an `example` value provided will be printed to the chronicle using exactly this value.

So from experience I would suggest to start with the `grid` option, get rough initial coordinates, and then switch to using both the `cellBorder` option and the `dummyValues` option to fine-tune everything.
