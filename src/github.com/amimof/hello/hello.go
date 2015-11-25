package main

import (
	"fmt"
	"os"
	"log"
	"regexp"
	"io/ioutil"
)

var prog string = os.Args[0]
//var mroot string = "/Users/amir/Documents/Media/Movies"
//var sroot string = "/Users/amir/Documents/Media/Series"
//var target string = "/Users/amir/Transmission"
var mroot string = os.Args[1]
var sroot string = os.Args[2]
var target string = os.Args[3]
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
	return prog + " <movies> <series> <target>"
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

	// Check if movies root exists
	if !exists(mroot) {
		log.Fatalln("Does not exist", mroot)
		os.Exit(1)
	}
	// Check if series root exists
	if !exists(sroot) {
		log.Fatalln("Does not exist", sroot)
		os.Exit(1)
	}
	// Check if movies root is a directory
	if isFile(mroot) {
		log.Fatalln("Is not a directory", mroot)
		os.Exit(1)
	}
	// Check if series root is a directory
	if isFile(sroot) {
		log.Fatalln("Is not a directory", sroot)
		os.Exit(1)
	}
	// Check if target exists
	if !exists(target) {
		log.Fatalln("Does not exist", target)
		os.Exit(1)
	}
	// Check if target is a directory
	if isFile(target) {
		log.Fatalln("Is not a directory", target)
		os.Exit(1)
	}

	// Main
	f, err := ioutil.ReadDir(target)
	for i, file := range f {
		log.Println(file.Name(), err, i)
	}



	// ismovie, err := isMovie(movie)
	// if ismovie {
	// 	log.Println(ismovie, err)
	// }

}
