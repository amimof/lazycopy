package cmd

/**
*
* Author: Amir Mofasser <amir.mofasser@gmail.com>
*	https://github.com/amimof
*
*/

import (
	"github.com/amimof/lazycopy/fileutils"
	"github.com/amimof/loglevel-go"
	"io/ioutil"
	"os"
	"regexp"
	"flag"
	"path"
	"fmt"
	"strings"
	"bufio"
)

var (
	prog string
	source string
	mroot string
	sroot string
	unit string
	overwrite bool
	confirm bool
	level int
	verify bool
)

type Movie struct {
	title string
	year string
	regexp string
	file os.FileInfo
}

type Serie struct {
	title string
	season string
	episode string
	regexp string
	file os.FileInfo
}

// List of file extensions. Otherwise we might get wierd matches when files contain numbers, such as log files.
// var extensions []string = []string{
// 	"mkv","MKV","mp4","MP4","m4p","M4P","m4v","M4V","mpg","MPG","mpeg","MPEG","mp2","MP2","mpe","MPE","mpv","MPV","3gp","3GP","nsv","NSV","f4v","F4V","f4p","F4P","f4a","F4A","f4b","F4P","vob","VOB","avi","AVI","mov","MOV","wmv","WMV","asd","ASD","flv","FLV","ogv","OGV","ogg","OGG","qt","QT","yuv","YUV","rm","RM","rmvb","RMVB",
// }
var extensions string = "\\.(mkv|MKV|mp4|MP4|m4p|M4P|m4v|M4V|mpg|MPG|mpeg|MPEG|mp2|MP2|mpe|MPE|mpv|MPV|3gp|3GP|nsv|NSV|f4v|F4V|f4p|F4P|f4a|F4A|f4b|F4P|vob|VOB|avi|AVI|mov|MOV|wmv|WMV|asd|ASD|flv|FLV|ogv|OGV|ogg|OGG|qt|QT|yuv|YUV|rm|RM|rmvb|RMVB)"
// Expressions to use when evaluating series
var spattern []string = []string{"(.*?)S(\\d{1,2})E(\\d{2})(.*)",
	"(.*?)s(\\d{1,2})e(\\d{2})(.*)",
	"(.*?)\\[?(\\d{1,2})x(\\d{2})\\]?(.*)",
	"(.*?)Season.?(\\d{1,2}).*?Episode.?(\\d{1,2})(.*)",
}
// Expressions to use when evaluating movies
var mpattern []string = []string{"(.*?)\\((17[0-9][0-9]|180[0-9]|181[0-9]|18[2-9]\\d|19\\d\\d|2\\d{3}|30[0-3]\\d|304[0-8])\\)(.*)",
	"(.*?)\\[(17[0-9][0-9]|180[0-9]|181[0-9]|18[2-9]\\d|19\\d\\d|2\\d{3}|30[0-3]\\d|304[0-8])\\](.*)",
	"(.*?)\\{(17[0-9][0-9]|180[0-9]|181[0-9]|18[2-9]\\d|19\\d\\d|2\\d{3}|30[0-3]\\d|304[0-8])\\}(.*)",
	"(.*?)(17[0-9][0-9]|180[0-9]|181[0-9]|18[2-9]\\d|19\\d\\d|2\\d{3}|30[0-3]\\d|304[0-8])(.*)",
	"(.*?)(\\d{3,4}p)(.*)",
}
var log *loglevel.Logger = loglevel.New()

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
func isMovie(file os.FileInfo, pattern []string) (*Movie, error) {
	for _, element := range pattern {
		if !file.IsDir() {
			element = element+extensions
		}
		r, err := regexp.Compile(element)
		if err == nil {
			match := r.MatchString(file.Name())
			if match {
				result := r.FindStringSubmatch(file.Name())
				movie := &Movie{strings.Replace(strings.Trim(strings.Trim(result[1], "."), " "), ".", " ", -1), strings.Trim(result[2], "."), element, file}
				return movie, nil
			}
		} else {
			return nil, err
		}
	}
	return nil, nil
}

func isSerie(file os.FileInfo, pattern []string) (*Serie, error) {
	for _, element := range pattern {
		if !file.IsDir() {
			element = element+extensions
		}
		r, err := regexp.Compile(element)
		if err == nil {
			match := r.MatchString(file.Name())
			if match {
				result := r.FindStringSubmatch(file.Name())
				serie := &Serie{strings.Replace(strings.Trim(strings.Trim(result[1], "."), " "), ".", " ", -1), strings.Trim(result[2], "."), strings.Trim(result[3], "."), element, file}
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

// Prompts the user for comfirmation by askig a yes/no question. The question can be 
// provided as msg. (yes/no) [default] will be appended to the msg. 
func confirmCopy(msg string, def bool) bool {

	var response string = ""
	var defaultChoice = "no"

	if def {
		defaultChoice = "yes"
	}
	fmt.Printf("%s (yes/no) [%s] ", msg, defaultChoice)

	scanner := bufio.NewScanner(os.Stdin)
	ok := scanner .Scan()

	if ok {
		
		response = strings.ToLower(strings.Trim(scanner.Text(), " "))

		if response == "y" || response == "yes" {
			return true
		} else if response == "n" || response == "no" {
			return false
		} else if response == "" && def {
			return true
		} else {
			fmt.Println("Please type (y)es or (n)o.")
			return confirmCopy(msg, def)
		}

	}

	return false

}

// Main loop
func Execute() {

	// Arguments and parameters to the program feels un-natural. Perhaps find a 
	// library that can handle the arguments for us. Also, thinking of renaming 'series' to 'tv-shows'

	// Read arguments
	flag.StringVar(&source, "S", ".", "Directories in which to look for media delimited by comma")
	flag.StringVar(&mroot, "m", ".", "Directory to your movies.")
	flag.StringVar(&sroot, "s", ".", "Directory to your series.")
	flag.BoolVar(&overwrite, "o", false, "Overwrite existing files/folders when copying")
	flag.BoolVar(&confirm, "c", false, "Prompt for confirm when overwriting existing files/folders")
	flag.IntVar(&level, "l", 1, "Log level. 3=DEBUG, 2=INFO, 1=WARN, 0=ERROR. (default \"0\")")
	flag.StringVar(&unit, "u", "g", "String representation of unit to use when calculating file sizes. Choices are k, m, g and t")
	flag.BoolVar(&verify, "v", false, "Verify, do not actually copy.")
	flag.Parse()

	// Sets the loglevel.
	// First we need to read from args and convert it to an int
	log.Level.SetLevel(level)
	log.PrintTime = false
	log.Debug("Log level is", level)

	// Check if movies root exists
	if !exists(mroot) {
		log.Errorf("Does not exist '%s'\n", mroot)
	}
	// Check if series root exists
	if !exists(sroot) {
		log.Errorf("Does not exist '%s'\n", sroot)
	}
	// Check if movies root is a directory
	if isFile(mroot) {
		log.Errorf("Is not a directory '%s'\n", mroot)
	}
	// Check if series root is a directory
	if isFile(sroot) {
		log.Errorf("Is not a directory '%s'\n", sroot)
	}

	sources := strings.Split(source, ",")

	for i, s := range sources {
		log.Debugf("[%b] - %s \n", i, s)
	
		// Check if source exists
		if !exists(s) {
			log.Error("Does not exist", s)
		}
		// Check if movies root is a directory
		if isFile(s) {
			log.Error("Is not a directory", s)
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
			movie, errM := isMovie(file, mpattern)
			if movie != nil {
				if errM == nil {
					log.Debugf("[%d] ==== MOVIE START ==== \n", j)
					log.Debugf("[%d] Movie. Title: '%s', Year: '%s', Filename: '%s'\n", j, movie.title, movie.year, movie.file.Name())
					log.Debugf("[%d] Movie matched regexp: '%s'\n", j, movie.regexp)
					mmatches = append(mmatches, file.Name())
					srcf := path.Join(src, file.Name())
					if confirm == true {
						overwrite = confirmCopy("Copy '"+file.Name()+"'?", true)
					}

					// Don't do anything if verify flag is true
					if !verify {
						var written int64
						dstf := path.Join(mroot, file.Name())
						log.Debugf("[%d] Dest is '%s' \n", j, dstf)

						// Start the copy
						written, err = fileutils.Copy(srcf, path.Join(mroot, file.Name()), overwrite)
						if err != nil {
							log.Errorf("[%b] Can't copy '%s'. %s \n", j, file.Name(), err)
						}

						totalWritten = totalWritten + written

					}
				}
				log.Debugf("[%d] ==== MOVIE END ==== \n", j)
			}

			// Check for series
			serie, errS := isSerie(file, spattern)
			if serie != nil {
				if errS == nil {
					log.Debugf("[%d] ==== SERIE START ====\n", j)
					log.Debugf("[%d] Serie. Title: '%s', Season: '%s', Episode: '%s', Filename: '%s'\n", j, serie.title, serie.season, serie.episode, serie.file.Name())
					log.Debugf("[%d] Serie matched regexp: '%s'\n", j, serie.regexp)
					var written int64
					smatches = append(smatches, file.Name())
					if confirm == true {
						overwrite = confirmCopy("Copy '"+file.Name()+"'?", true)
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
						written, err = fileutils.Copy(srcf, dstf, overwrite)
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

	fmt.Println("Movies matched:", len(mmatches))
	fmt.Println("Series matched:", len(smatches))
	fmt.Printf("Copied %.2f%s\n", convertFileSize(totalWritten, unit), unit)

}
