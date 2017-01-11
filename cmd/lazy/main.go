package cmd

/**
*
* Author: Amir Mofasser <amir.mofasser@gmail.com>
*	https://github.com/amimof
*
*/

import (
	copy "github.com/amimof/lazycopy/copy"
	logger "github.com/amimof/lazycopy/logger"
	"io/ioutil"
	"os"
	"regexp"
	"flag"
	"path"
	"fmt"
	"strings"
)

var (
	prog string
	source string
	mroot string
	sroot string
	unit string
	overwrite bool
	confirm bool
	loglevel int
	verify bool
)

type Movie struct {
	title string
	year string
	filename string
}

type Serie struct {
	title string
	season string
	episode string
	filename string
}

// Extensions in regexp disabled atm since it prevents us from detecting folders
// var extensions string = "\\.(mkv|MKV|mp4|MP4|m4p|M4P|m4v|M4V|mpg|MPG|mpeg|MPEG|mp2|MP2|mpe|MPE|mpv|MPV|3gp|3GP|nsv|NSV|f4v|F4V|f4p|F4P|f4a|F4A|f4b|F4P|vob|VOB|avi|AVI|mov|MOV|wmv|WMV|asd|ASD|flv|FLV|ogv|OGV|ogg|OGG|qt|QT|yuv|YUV|rm|RM|rmvb|RMVB)"
var extensions string = ""
var spattern []string = []string{"(.*?)S(\\d{1,2})E(\\d{2})(.*)"+extensions,
	"(.*?)s(\\d{1,2})e(\\d{2})(.*)"+extensions,
	"(.*?)\\[?(\\d{1,2})x(\\d{2})\\]?(.*)"+extensions,
	"(.*?)Season.?(\\d{1,2}).*?Episode.?(\\d{1,2})(.*)"+extensions}
var mpattern []string = []string{"(.*?)\\((17[0-9][0-9]|180[0-9]|181[0-9]|18[2-9]\\d|19\\d\\d|2\\d{3}|30[0-3]\\d|304[0-8])\\)(.*)"+extensions,
	"(.*?)\\[(17[0-9][0-9]|180[0-9]|181[0-9]|18[2-9]\\d|19\\d\\d|2\\d{3}|30[0-3]\\d|304[0-8])\\](.*)"+extensions,
	"(.*?)\\{(17[0-9][0-9]|180[0-9]|181[0-9]|18[2-9]\\d|19\\d\\d|2\\d{3}|30[0-3]\\d|304[0-8])\\}(.*)"+extensions,
	"(.*?)(17[0-9][0-9]|180[0-9]|181[0-9]|18[2-9]\\d|19\\d\\d|2\\d{3}|30[0-3]\\d|304[0-8])(.*)"+extensions,
	"(.*?)(\\d{3,4}p)(.*)"}
var log *logger.Logger = logger.SetupNew("MAIN")

// Returns true if path exists
func exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return true
}

// Return true if path is a file
func isFile(path string) bool {
	file, err := os.Stat(path)
	if err == nil && file.IsDir() != true {
		return true
	}
	return false
}

// Checks wether a given filename is considered to be a movie or series
// based on the specified regexp patterns
func isMovie(filename string, pattern []string) (*Movie, error) {
	for index, element := range pattern {
		_ = index
		r, err := regexp.Compile(element)
		if err == nil {
			match := r.MatchString(filename)
			if match {
				result := r.FindStringSubmatch(filename)
				movie := &Movie{strings.Replace(strings.Trim(strings.Trim(result[1], "."), " "), ".", " ", -1), strings.Trim(result[2], "."), filename}
				return movie, nil
			}
		} else {
			return nil, err
		}
	}
	return nil, nil
}

func isSerie(filename string, pattern []string) (*Serie, error) {
	for index, element := range pattern {
		_ = index
		r, err := regexp.Compile(element)
		if err == nil {
			match := r.MatchString(filename)
			if match {
				result := r.FindStringSubmatch(filename)
				// Fix this so that dots are replaced with space in in the title.
				serie := &Serie{strings.Replace(strings.Trim(strings.Trim(result[1], "."), " "), ".", " ", -1), strings.Trim(result[2], "."), strings.Trim(result[3], "."), filename}
				return serie, nil				
			}
		} else {
			return nil, err
		}
	}
	return nil, nil
}

// Converts file size from bytes to kb, mb, gb or tb
func convertFileSize(size int64, unit string) float64 {
	var result float64 = 0
	switch unit {
		case "k":
			result = float64(size) / 1024
		case "m":
			result = (float64(size) / 1024) / 1024
		case "g":
			result = ((float64(size) / 1024) / 1024) / 1024
		case "t":
			result = (((float64(size) / 1024) / 1024) / 1024) / 1024
		default:
			result = ((float64(size) / 1024) / 1024) / 1024
	}
	return result
}

// Confirm prompt. Accept y/n from user. Returns true or false
func confirmCopy(msg string) bool {

	fmt.Println(msg)

	var response string
	_, err := fmt.Scanln(&response)
	if err != nil {
		log.Error(err)
	}

	// Lowercase response. Remove whitespace
	response = strings.ToLower(response)
	response = strings.Trim(response, " ")

	r := string(response[0])
	if r == "y" {
		return true
	} else if r == "n" {
		return false
	} else {
		fmt.Println("Please type yes or no")
		return confirmCopy(msg)
	}
}

// Main loop
func Execute() {

	// Read arguments
	flag.StringVar(&source, "S", ".", "Directories in which to look for media delimited by comma")
	flag.StringVar(&mroot, "m", ".", "Directory to your movies.")
	flag.StringVar(&sroot, "s", ".", "Directory to your series.")
	flag.BoolVar(&overwrite, "o", false, "Overwrite existing files/folders when copying")
	flag.BoolVar(&confirm, "c", false, "Prompt for confirm when overwriting existing files/folders")
	flag.IntVar(&loglevel, "v", 0, "Log level. 3=DEBUG, 2=WARN, 1=INFO, 0=DEBUG. (default \"0\")")
	flag.StringVar(&unit, "u", "g", "String representation of unit to use when calculating file sizes. Choices are k, m, g and t")
	flag.BoolVar(&verify, "V", false, "Verify, do not actually copy.")
	flag.Parse()

	// Sets the loglevel.
	// First we need to read from args and convert it to an int
	log.Level.SetLevel(loglevel)
	log.Debug("Log level is", loglevel)

	// Check if movies root exists
	if !exists(mroot) {
		log.Errorf("Does not exist '%s'\n", mroot)
		os.Exit(1)
	}
	// Check if series root exists
	if !exists(sroot) {
		log.Errorf("Does not exist '%s'\n", sroot)
		os.Exit(1)
	}
	// Check if movies root is a directory
	if isFile(mroot) {
		log.Errorf("Is not a directory '%s'\n", mroot)
		os.Exit(1)
	}
	// Check if series root is a directory
	if isFile(sroot) {
		log.Errorf("Is not a directory '%s'\n", sroot)
		os.Exit(1)
	}

	sources := strings.Split(source, ",")

	for i, s := range sources {
		log.Debugf("[%b] - %s \n", i, s)
	
		// Check if source exists
		if !exists(s) {
			log.Error("Does not exist", s)
			os.Exit(1)
		}
		// Check if movies root is a directory
		if isFile(s) {
			log.Error("Is not a directory", s)
			os.Exit(1)
		}
	}

	// Main
	var mmatches []string
	var smatches []string
	var totalWritten int64

	log.Debug("Overwrite is set to", overwrite)

	for i, src := range sources {

		f, err := ioutil.ReadDir(src)
		log.Debugf("[%d] Source is '%s'\n", i, src)
		var index int64 = 0

		for j, file := range f {
			log.Debugf("Checking '%s' \n", file.Name())

			// Check for movies
			movie, errM := isMovie(file.Name(), mpattern)
			if movie != nil {
				if errM == nil {
					log.Debugf("[%d] ==== MOVIE START ==== \n", j)
					log.Debugf("[%d] Movie. Title: '%s', Year: '%s', Filename: '%s', \n", j, movie.title, movie.year, movie.filename)
					mmatches = append(mmatches, file.Name())
					srcf := path.Join(src, file.Name())
					if confirm == true {
						overwrite = confirmCopy("Copy? (y/n) '" + file.Name() + "'")
					}

					// Don't do anything if verify flag is true
					if !verify {
						var written int64
						dstf := path.Join(mroot, file.Name())
						log.Debugf("[%d] Dest is '%s' \n", j, dstf)

						// Start the copy
						written, err = copy.Copy(srcf, path.Join(mroot, file.Name()), overwrite)
						if err != nil {
							log.Errorf("[%b] Can't copy '%s'. %s \n", j, file.Name(), err)
						}

						totalWritten = totalWritten + written

					}
				}
				log.Debugf("[%d] ==== MOVIE END ==== \n", j)
			}

			// Check for series
			serie, errS := isSerie(file.Name(), spattern)
			if serie != nil {
				if errS == nil {
					log.Debugf("[%d] ==== SERIE START ==== \n", j)
					log.Debugf("[%d] Serie. Title: '%s', Season: '%s', Episode: '%s', Filename: '%s' \n", j, serie.title, serie.season, serie.episode, serie.filename)
					var written int64
					smatches = append(smatches, file.Name())
					if confirm == true {
						overwrite = confirmCopy("Copy? (y/n) '" + file.Name() + "'")
					}

					// Don't do anything if verify flag is true
					if !verify {
						srcf := path.Join(src, file.Name())

						// Stat source so that we can perserve permissions when creating the directories if necessary
						s, err := os.Stat(path.Dir(srcf))
						if err != nil {
							log.Errorf("[%d] Couldn't stat. '%s' \n", j, err)
						}

						// Create serie folder and season folders resursively
						dstFolder := path.Join(sroot, serie.title, "Season "+serie.season)
						if !exists(dstFolder) {
							log.Debugf("[%s] Dest does not exist, creating '%s' \n", j, dstFolder)
							err = os.MkdirAll(dstFolder, s.Mode())
							if err != nil {
								log.Errorf("[%d] Couldn't create '%s'. %s \n", j, dstFolder, err)
							}
						}

						// Start copying
						dstf := path.Join(dstFolder, file.Name())
						log.Debugf("[%d] Dest is '%s' \n", j, dstf)
						written, err = copy.Copy(srcf, dstf, overwrite)
						if err != nil {
							log.Errorf("[%d] Can't copy '%s'. %s", j, file.Name(), err)
						}

						totalWritten = totalWritten + written

					}	
				}
				log.Debugf("[%d] ==== SERIE END ==== \n", j)
			}
			index++
		}
	}

	for _, arg := range flag.Args() {
		fmt.Println(arg)
	}

	fmt.Println(totalWritten)

	fmt.Println("Movies matched:", len(mmatches))
	fmt.Println("Series matched:", len(smatches))
	fmt.Printf("Copied %.2f%s\n", convertFileSize(totalWritten, unit), unit)

}
