# lazycopy
A utility written in Go that copies movies and tv-shows from source, usually a download directory, to a target directory. Scans the source directories for video content and organizes them at a target directory. It is capable of
detecting movies and tv-shows based on their file name using complex regular expressions. 

## Usage
Search the source (-S) directories for movies and tv-shows. Copies any movies found to (-m) and series to (-s). Below example will overwrite any already existing files and folders (-o)
```
./lazycopy -S ~/Downloads -m ~/Media/Movies -s ~/Media/TV-Shows -o
```

Other available command line options are
```
Usage of ./lazycopy:
  -S string
      Directories in which to look for media delimited by comma (default ".")
  -c  Prompt for confirm when overwriting existing files/folders
  -l int
      Log level. 3=DEBUG, 2=WARN, 1=INFO, 0=ERROR. (default "0")
  -m string
      Directory to your movies. (default ".")
  -o  Overwrite existing files/folders when copying
  -s string
      Directory to your series. (default ".")
  -u string
      String representation of unit to use when calculating file sizes. Choices are k, m, g and t (default "g")
  -v  Verify, do not actually copy.
```

## Example

The download directory `~/Downloads/`  contains following content:

```
.
..
Some.Movie.1977.mkv
Ultra.Hi.Resolution.Film.2016.Blu-Ray.avi
Movie.Directory.Split.In.Two.CDS.1998
Typicall.TV-show.S01E01
Some.Sitcom.S02E01
```

The following command:
```
./lazycopy -S ~/Downloads -m ~/Media/Movies -s ~/Media/Tv-Shows -o
```

Will copy all movies, wether it's a file or a folder, and overwrite existing ones to `~/Media/Movies`. In this example, movies will be copied to:
```
.
..
~/Media/Movies/Some.Movie.1977.mkv
~/Media/Movies/Ultra.Hi.Resolution.Film.2016.Blu-Ray.avi
~/Media/Movies/Movie.Directory.Split.In.Two.CDS.1998
```

TV-Show season directories will be created if necessary. In this example, the tv-shows will be copied to:

```
.
..
~/Media/Tv-Shows/Typicall TV-Show/Season 01/Typicall.TV-show.S01E01
~/Media/Tv-Shows/Some Sitcom/Season 02/Some.Sitcom.S02E01
```