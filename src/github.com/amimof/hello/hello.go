package main

import (
	"fmt"
	"os"
	"log"
	"regexp"
)

var prog string = os.Args[0]
var mroot string = "/Volumes/Store2/Media/Video/Movies"
var movie string = "C:\\temp\\Inception.2014.Blu-ray.mkv"
var spattern []string = []string{"(.*?)S(\\d{1,2})E(\\d{2})(.*)", "(.*?)s(\\d{1,2})e(\\d{2})(.*)", "(.*?)\\[?(\\d{1,2})x(\\d{2})\\]?(.*)", "(.*?)Season.?(\\d{1,2}).*?Episode.?(\\d{1,2})(.*)"}
var mpattern []string = []string{"(.*?)\\((17[0-9][0-9]|180[0-9]|181[0-9]|18[2-9]\\d|19\\d\\d|2\\d{3}|30[0-3]\\d|304[0-8])\\)(.*)", "(.*?)\\[(17[0-9][0-9]|180[0-9]|181[0-9]|18[2-9]\\d|19\\d\\d|2\\d{3}|30[0-3]\\d|304[0-8])\\](.*)", "(.*?)\\{(17[0-9][0-9]|180[0-9]|181[0-9]|18[2-9]\\d|19\\d\\d|2\\d{3}|30[0-3]\\d|304[0-8])\\}(.*)", "(.*?)(17[0-9][0-9]|180[0-9]|181[0-9]|18[2-9]\\d|19\\d\\d|2\\d{3}|30[0-3]\\d|304[0-8])(.*)", "(.*?)(\\d{3,4}p)(.*)"}

type Pattern struct {
	Data []string
}

func exists(path string) (bool) {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func isFile(path string) (bool) {
	file, err := os.Stat(path)
	if err == nil && file.IsDir() != true {
		return true
	} 
	return false
}

func usage() string {
	return prog + " <arguments>"
}

func isMovie(filename string) (bool, error) {
	log.Println("Checking if " + filename + " is a movie")
	for index, element := range mpattern {
		r, err := regexp.Compile(element)
		if err == nil {
			match := r.MatchString(filename)	
			if match {
				log.Println("Found match", index, match)
				return true, nil
			}
		} else {
			return false, err
		}		
	}
	log.Println("No matches found")
	return false, nil
}

func main() {	

	// Print usage to the user
	fmt.Println("Usage: " + usage())

	dirExists := exists(mroot)
	isFile := isFile(movie)

	ismovie, err := isMovie(movie)
	if ismovie {
		log.Println(ismovie, err)
	}

	if dirExists {
		fmt.Println("Dir exists")
	} else {
		fmt.Println("Dir does not exist")
	}

	if isFile {
		fmt.Println("This is a file")
	} else {
		fmt.Println("This is not a file")
	}

}
