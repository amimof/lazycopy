# lazycopy
A utility that copies movies and tv-shows from source, usually a download directory, to a target directory. Scans the source directory for video content and organizes them at a target directory. It will automatically detect movies and tv-shows based on their file name using complex regular expressions. When matching media is found, the file or folder is copied to the target diretory, preserving permissions and creating any sub-folders if needed.

![Jan-08-2018 15-59-39](https://gfycat.com/ifr/FilthyImpossibleBluefintuna)

## Installation

### From source
```
$ git clone https://github.com/amimof/lazycopy.git
$ cd lazycopy
$ make
```

### From built binaries
Go to https://github.com/amimof/lazycopy/releases and download the binary for your target platform

## Usage
```
$ lazycopy [options] <source dir> <target dir>

options:
  -c    Prompt for confirm before overwriting existing files/folders
  -d    Debug mode
  -o    Overwrite existing files/folders when copying
  -q    Supress all output
  -v    Print version info
```

## Example

The download directory `~/Downloads/` contains following content:

```
.
..
Some.Movie.1977.mkv
Ultra.Hi.Resolution.Film.2016.Blu-Ray.avi
Movie.Directory.Split.In.Two.CDS.1998
Typicall.TV-show.S01E01
Some.Sitcom.S02E01
```

The following command will copy all movies, wether it's a file or a folder, and overwrite existing ones to `/Volumes/USBDrive/Movies` and `/Volumes/USBDrive/Series`. Note that the folders `Movies` and `Series` will be created automatically in `/Volumes/USBDrive` if they don't already exist.
```
$ lazycopy ~/Downloads /Volumes/USBDrive
```


The destination directory will look like this when the copy operation finishes
```
.
..
/Volumes/USBDrive/Movies/Some.Movie.1977.mkv
/Volumes/USBDrive/Movies/Ultra.Hi.Resolution.Film.2016.Blu-Ray.avi
/Volumes/USBDrive/Movies/Movie.Directory.Split.In.Two.CDS.1998
/Volumes/USBDrive/Series/Typicall TV-Show/Season 01/Typicall.TV-show.S01E01
/Volumes/USBDrive/Series/Some Sitcom/Season 02/Some.Sitcom.S02E01
```