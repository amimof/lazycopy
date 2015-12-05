# video-scanner-go
A utility written in Go that copies your movies and series from source to target directory.
Scans a directory for video content and organizes them at a target directory. It is capable of
detecting movies and series based on their file name.

# Usage
Search the target (-t) directory ~/Downloads for media. Copies any movies (-m) to ~/Media/Movies and (-s) series to ~/Media/Series.
Below example will overwrite any already existing files and folders (-o)
```
./scanner -m ~/Media/Movies -s ~/Media/Series -t ~/Downloads -o
```

Other available command line options are
```
./scanner -h
Usage of ./scanner:
  -c    Prompt for confirm when overwriting existing files/folders
  -m string
        Directory to your movies. Current directory (.) by default (default ".")
  -o    Set to true to overwrite existing files/folders when copying
  -s string
        Directory to your series. Current directory (.) by default (default ".")
  -t string
        Target directory. Typically your Downloads folder. Current directory (.) by default (default ".")
  -v int
        Log level. 3=DEBUG, 2=WARN, 1=INFO, 0=DEBUG
```
