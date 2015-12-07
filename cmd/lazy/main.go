package main

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
	mroot string
	sroot string
	target string
	unit string
	overwrite bool
	confirm bool
	loglevel int
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
				movie := &Movie{strings.Trim(result[1], "."), strings.Trim(result[2], "."), filename}
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
				serie := &Serie{strings.Trim(result[1], "."), strings.Trim(result[2], "."), strings.Trim(result[3], "."), filename}
				return serie, nil
			}
		} else {
			return nil, err
		}
	}
	return nil, nil
}

// Converts file size from bytes to kb, mb, gb or tb
func convertFileSize(size int64, unit string) int64 {
	var result int64
	switch unit {
		case "k":
			result = size / 1024
		case "m":
			result = (size / 1024) / 1024
		case "g":
			result = ((size / 1024) / 1024) / 1024
		case "t":
			result = (((size / 1024) / 1024) / 1024) / 1024
		default:
			result = ((size / 1024) / 1024) / 1024
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
func main() {

	// Read arguments
	flag.StringVar(&mroot, "m", ".", "Directory to your movies.")
	flag.StringVar(&sroot, "s", ".", "Directory to your series.")
	flag.StringVar(&target, "t", ".", "Target directory. Typically your Downloads folder.")
	flag.BoolVar(&overwrite, "o", false, "Overwrite existing files/folders when copying")
	flag.BoolVar(&confirm, "c", false, "Prompt for confirm when overwriting existing files/folders")
	flag.IntVar(&loglevel, "v", 0, "Log level. 3=DEBUG, 2=WARN, 1=INFO, 0=DEBUG. (default \"0\")")
	flag.StringVar(&unit, "u", "g", "String representation of unit to use when calculating file sizes. Choices are k, m, g and t")
	flag.Parse()

	// Sets the loglevel.
	// First we need to read from args and convert it to an int
	log.Level.SetLevel(loglevel)
	log.Debug("Log level is", loglevel)

	// Check if movies root exists
	if !exists(mroot) {
		log.Error("Does not exist", mroot)
		os.Exit(1)
	}
	// Check if series root exists
	if !exists(sroot) {
		log.Error("Does not exist", sroot)
		os.Exit(1)
	}
	// Check if movies root is a directory
	if isFile(mroot) {
		log.Error("Is not a directory", mroot)
		os.Exit(1)
	}
	// Check if series root is a directory
	if isFile(sroot) {
		log.Error("Is not a directory", sroot)
		os.Exit(1)
	}
	// Check if target exists
	if !exists(target) {
		log.Error("Does not exist", target)
		os.Exit(1)
	}
	// Check if target is a directory
	if isFile(target) {
		log.Error("Is not a directory", target)
		os.Exit(1)
	}

	// Main
	var mmatches []string
	var smatches []string
	f, err := ioutil.ReadDir(target)
	log.Debug("Overwrite is set to", overwrite)
	log.Debug("Looking in", target)
	var totalWritten int64
	for i, file := range f {
		log.Debug("Checking", file.Name())

		// Check for movies
		movie, errM := isMovie(file.Name(), mpattern)
		if movie != nil {
			if errM == nil {
				log.Debug("==== MOVIE START ====")
				log.Debug("Movie", movie.title, movie.year, movie.filename, i, err)
				mmatches = append(mmatches, file.Name())
				srcf := path.Join(target, file.Name())
				if confirm == true {
					overwrite = confirmCopy("Overwrite? (y/n) " + file.Name())
				}
				var written int64
				if isFile(srcf) {
					dstf := path.Join(mroot, file.Name())
					log.Debug("Dest is", dstf)
					log.Debug("Is file", dstf)
					written, err = copy.CopyFile(srcf, path.Join(mroot, file.Name()), overwrite)
					if err != nil {
						log.Error("Can't copy", file.Name(), err)
					}
					totalWritten = totalWritten + written
				} else {
					dstf := path.Join(mroot, file.Name())
					log.Debug("Is dir", dstf)
					written, err = copy.CopyDir(srcf, path.Join(mroot, file.Name()), overwrite)
					if err != nil {
						log.Error("Can't copy", file.Name(), err)
					}
					totalWritten = totalWritten + written
				}
			}
			log.Debug("==== MOVIE END ====")
		}

		// Check for series
		serie, errS := isSerie(file.Name(), spattern)
		if serie != nil {
			if errS == nil {
				log.Debug("==== MOVIE START ====")
				log.Debug("Serie", serie.title, serie.season, serie.episode, serie.filename)
				var written int64
				smatches = append(smatches, file.Name())
				if confirm == true {
					overwrite = confirmCopy("Overwrite? (y/n) " + file.Name())
				}
				srcf := path.Join(target, file.Name())
				dstFolder := path.Join(sroot, serie.title, "Season"+serie.season)
				if !exists(dstFolder) {
						log.Debug("Dest does not exist, creating", dstFolder)
						err = os.MkdirAll(dstFolder, 0700)
						if err != nil {
							log.Error("Couldn't create", dstFolder, err)
						}
				}
				dstf := path.Join(dstFolder, file.Name())
				log.Debug("Dest is", dstf)
				if isFile(srcf) {
					log.Debug("Is file", srcf)
					written, err = copy.CopyFile(srcf, dstf, overwrite)
					if err != nil {
						log.Error("Can't copy", file.Name(), err)
					}
					totalWritten = totalWritten + written
				} else {
					log.Debug("Is dir", srcf)
					written, err = copy.CopyDir(srcf, dstf, overwrite)
					if err != nil {
						log.Error("Can't copy", file.Name(), err)
					}
					totalWritten = totalWritten + written
				}
			}
			log.Debug("==== SERIE END ====")
		}
	}

	log.Debug("Movies matched in dir is", len(mmatches))
	log.Debug("Series matched in dir is", len(smatches))
	fmt.Println("Copied", convertFileSize(totalWritten, unit), unit)

}
