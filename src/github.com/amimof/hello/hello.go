package main

/**
*
* Author: Amir Mofasser <amir.mofasser@gmail.com>
*	https://github.com/amimof
*
*/

import (
	// "github.com/cheggaaa/pb"
	"github.com/amimof/logger"
	"io/ioutil"
	"os"
	"regexp"
	"flag"
	"path"
	// "time"
)

var (
	prog string
	mroot string
	sroot string
	target string
	loglevel int
)

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

// Checks wether a given filename is considered to be a movie or series
// based on the specified regexp patterns
func isMedia(filename string, pattern []string) (bool, error) {
	for index, element := range pattern {
		_ = index
		r, err := regexp.Compile(element)
		if err == nil {
			match := r.MatchString(filename)
			if match {
				return true, nil
			}
		} else {
			return false, err
		}
	}
	return false, nil
}

func main() {

	// Read arguments
	flag.StringVar(&mroot, "m", ".", "Directory to your movies. Current directory (.) by default")
	flag.StringVar(&sroot, "s", ".", "Directory to your series. Current directory (.) by default")
	flag.StringVar(&target, "t", ".", "Target directory. Typically your Downloads folder. Current directory (.) by default")
	flag.IntVar(&loglevel, "v", 3, "Log level. 3=DEBUG, 2=WARN, 1=INFO, 0=DEBUG")
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
	log.Debug("Looking in", target)
	for i, file := range f {
		log.Debug("Checking", file.Name())
		if isM, errM := isMedia(file.Name(), mpattern); isM {
			if errM == nil {
				log.Debug("Movie found", file.Name(), i, err)
				mmatches = append(mmatches, file.Name())
				srcf := path.Join(target, file.Name())
				log.Info("Copying", srcf)
			} else {
				log.Error("Error", errM)
			}
		}
		if isS, errS := isMedia(file.Name(), spattern); isS {
			if errS == nil {
				log.Debug("Serie found", file.Name(), i, err)
				smatches = append(smatches, file.Name())
			} else {
				log.Error("Error", errS)
			}
		}
	}

	log.Info("Movies matched in dir is", len(mmatches))
	log.Info("Series matched in dir is", len(smatches))

}
