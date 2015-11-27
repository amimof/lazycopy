package main

import (
	// "github.com/cheggaaa/pb"
	"github.com/amimof/logger"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	// "time"
)

var prog string = os.Args[0]
var mroot string = os.Args[1]
var sroot string = os.Args[2]
var target string = os.Args[3]
var extensions string = "\\.(mkv|MKV|mp4|MP4|m4p|M4P|m4v|M4V|mpg|MPG|mpeg|MPEG|mp2|MP2|mpe|MPE|mpv|MPV|3gp|3GP|nsv|NSV|f4v|F4V|f4p|F4P|f4a|F4A|f4b|F4P|vob|VOB|avi|AVI|mov|MOV|wmv|WMV|asd|ASD|flv|FLV|ogv|OGV|ogg|OGG|qt|QT|yuv|YUV|rm|RM|rmvb|RMVB)"
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

// Prints program usage to the user
func usage() string {
	return *(&prog) + " <movies> <series> <target> <loglevel>"
}

// Checks wether a given filename is considered to be a movie
func isMovie(filename string) (bool, error) {
	log.Debug("Checking if " + filename + " is a movie")
	for index, element := range mpattern {
		r, err := regexp.Compile(element)
		if err == nil {
			match := r.MatchString(filename)
			if match {
				log.Debug(index, match)
				return true, nil
			}
		} else {
			return false, err
		}
	}
	log.Debug(filename, "did not match")
	return false, nil
}

func main() {

	// Sets the loglevel.
	// First we need to read from args and convert it to an int
	loglevel, err := strconv.Atoi(os.Args[4])
	if err != nil {
		log.Error("Unexpected arguments")
	}
	log.Level.SetLevel(loglevel)

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
	f, err := ioutil.ReadDir(target)
	log.Debug("Looking in", target)
	for i, file := range f {
		if isM, errM := isMovie(file.Name()); isM {
			if errM == nil {
				log.Debug("Is a movie", file.Name(), i, err)
				mmatches = append(mmatches, file.Name())
			} else {
				log.Error("Error", errM)
			}
		}
	}
	log.Info("Movies matched in dir is", len(mmatches))

}
