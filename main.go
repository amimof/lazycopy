package main

/**
*
* Author: Amir Mofasser <amir.mofasser@gmail.com>
*	https://github.com/amimof
*
*/

import (
	"os"
	"io"
	"io/ioutil"
	"regexp"
	"flag"
	"path"
  "path/filepath"
	"strings"
	"bufio"
  "time"
	"github.com/amimof/loglevel-go"
	"github.com/cheggaaa/pb"
	"fmt"
)

var (
	bar *pb.ProgressBar
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

var extensions string = "\\.(mkv|MKV|mp4|MP4|m4p|M4P|m4v|M4V|mpg|MPG|mpeg|MPEG|mp2|MP2|mpe|MPE|mpv|MPV|3gp|3GP|nsv|NSV|f4v|F4V|f4p|F4P|f4a|F4A|f4b|F4P|vob|VOB|avi|AVI|mov|MOV|wmv|WMV|asd|ASD|flv|FLV|ogv|OGV|ogg|OGG|qt|QT|yuv|YUV|rm|RM|rmvb|RMVB)"
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
func isMovie(file os.FileInfo) (*Movie, error) {
	var pattern []string = []string{"(.*?)\\((17[0-9][0-9]|180[0-9]|181[0-9]|18[2-9]\\d|19\\d\\d|2\\d{3}|30[0-3]\\d|304[0-8])\\)(.*)",
		"(.*?)\\[(17[0-9][0-9]|180[0-9]|181[0-9]|18[2-9]\\d|19\\d\\d|2\\d{3}|30[0-3]\\d|304[0-8])\\](.*)",
		"(.*?)\\{(17[0-9][0-9]|180[0-9]|181[0-9]|18[2-9]\\d|19\\d\\d|2\\d{3}|30[0-3]\\d|304[0-8])\\}(.*)",
		"(.*?)(17[0-9][0-9]|180[0-9]|181[0-9]|18[2-9]\\d|19\\d\\d|2\\d{3}|30[0-3]\\d|304[0-8])(.*)",
		"(.*?)(\\d{3,4}p)(.*)",
	}
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

// Checks wether a given filename is considered to be a movie or series
// based on the specified regexp patterns
func isSerie(file os.FileInfo) (*Serie, error) {
	var pattern []string = []string{"(.*?)S(\\d{1,2})E(\\d{2})(.*)",
		"(.*?)s(\\d{1,2})e(\\d{2})(.*)",
		"(.*?)\\[?(\\d{1,2})x(\\d{2})\\]?(.*)",
		"(.*?)Season.?(\\d{1,2}).*?Episode.?(\\d{1,2})(.*)",
	}
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
func convertUnit(size int64) string {
	result := fmt.Sprintf("%d %s", size, "bytes")

	// Convert to kb
	if size >= 1024 {
		result = fmt.Sprintf("%d %s", size / 1024, "kB")
	}
	// Convert to mb
	if size >= (1024 * 1024) {
		result = fmt.Sprintf("%d %s", (size / 1024) / 1024, "MB")
	}
	// Convert to gb
	if size >= (1024 * 1024 * 1024) {
		result = fmt.Sprintf("%d %s", ((size / 1024) / 1024) / 1024, "GB")
	}
	// Convert to tb
	if size >= (1024 * 1024 * 1024 * 1024) {
		result = fmt.Sprintf("%d %s", (((size / 1024) / 1024) / 1024) / 1024, "TB")
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
	log.Printf("%s (yes/no) [%s] ", msg, defaultChoice)

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
			log.Println("Please type (y)es or (n)o.")
			return confirmCopy(msg, def)
		}

	}

	return false
}


/*
 * Copies a single file or a an entire directory with it's content from src to dst. 
 * Will overwrite any file that already exists if 'ow' is set to true
 *
 */
func copy(src, dst string, ow bool) (int64, error) {

	// Stat source
	s, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	var written int64
	var bar *pb.ProgressBar

	// Check if src is a dir
	if s.IsDir() {

		// Calculate total directory size before copying so that we can pass it to the progress bar
		dirSize, err := calculateSize(src)
		if err != nil {
			return 0, err
		}

		// Create progress bar and init some defaults
		bar = initBar(src, dirSize)

		// Start progress bar and start copying
		bar.Start()
		written, err = copyDir(src, dst, ow)
		
	}

	// Check if src is a file
	if !s.IsDir() {

		// Create progress bar and init some defaults
		bar = initBar(src, s.Size())

		// Start progress bar and start  copying
		bar.Start()
		written, err = copyFile(src, dst, ow)
		
	}

	bar.Finish()

	return written, err

}

// Customize the bar. This function is mainly so that we don't have to write the same code twise.
// Don't like the fact that this function takes the argument src. Might need some work.
func initBar(src string, size int64) *pb.ProgressBar {
	bar = pb.New64(size)
	bar.SetUnits(pb.U_BYTES)
	bar.SetRefreshRate(time.Millisecond*10)
	bar.Prefix(truncate(path.Base(src), 40, true)+": ")	
	bar.Format("[-> ]")
	bar.ShowSpeed =  true
	return bar
}

// Calculate total size of a directory
func calculateSize(src string) (int64, error) {
	var size int64
	err := filepath.Walk(src, func(_ string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	return size, err
}

// Copies a dir tree from src to dst. Overwrites content if ow is true
func copyDir(src, dst string, ow bool) (int64, error) {

	var written int64

	// Stat source
	s, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	// Ensure that src is a dir
	if !s.IsDir() {
		return 0, err
	}

	// Ensure that dst does not already exist
	d, err := os.Open(dst)
	if err != nil {
		if !os.IsNotExist(err) && ow == false {
			return 0, err	
		}
	}
	defer d.Close()

	// Create dest dir
	err = os.MkdirAll(dst, s.Mode())
	if err != nil {
		return 0, err
	}

	//var written int64
	entries, err := ioutil.ReadDir(src)

	if err != nil {
		return 0, err
	} else {

		for _, entry := range entries {

			sfp := path.Join(src, entry.Name())
			dfp := path.Join(dst, entry.Name())

			if entry.IsDir() {
				w, err := copyDir(sfp, dfp, ow)
				if err != nil {
					return 0, err
				}
				written += w
			} else {
				w, err := copyFile(sfp, dfp, ow)
				if err != nil {
					return 0, err
				}
				written += w
			}
		}
	}

	return written, nil

}

// Copies the content of src file to dst. Overwrites dst file if ow os true
func copyFile(src, dst string, ow bool) (int64, error) {

	// Create source
	s, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer s.Close()

	// Check if dst exists
	d, err := os.Open(dst)
	if !os.IsNotExist(err) {
		if ow == false {
		  return 0, err
		}
	}
	defer d.Close()

	// Create dest
	dest, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer dest.Close()

	// Copy
	writer := io.MultiWriter(dest, bar)
	written, err := io.Copy(writer, s)

	if err != nil {
		return 0, err
	}

	return written, nil
}


// Truncates string to the length specified with len.
// if bo is true, then truncated text will be represented by . chars
func truncate(str string, length int, bo bool) string {
	
	var cut int
	var newStr string
	var appendStr string

	// Dont cut more than the max number of chars in the string
	// Will throw 'slice bounds out of range' otherwise
	strLen := len([]rune(str))
	if length >= strLen {
		cut = strLen
	} else {
		cut = length
	}

	newStr = str[:cut]

	if bo == true {
		appendStr = str[(strLen-3):strLen]
	}

  returnStr := newStr + "..." + appendStr

  // If the new string is longer than the originial, the just return the originial string
  if returnStr >= str {
    returnStr = str
  }

	return returnStr
}

// Main loop
func main() {

	// Read arguments
	overwrite := flag.Bool("o", false, "Overwrite existing files/folders when copying")
	confirm := flag.Bool("c", false, "Prompt for confirm when overwriting existing files/folders")
	level := flag.Int("l", 1, "Log level. 3=DEBUG, 2=INFO, 1=WARN, 0=ERROR. (default \"0\")")
	version := flag.Bool("v", false, "Print version info")

	flag.Parse()

	// Print version info and exit
	if *version {
		fmt.Println("1.0.2")
		os.Exit(0)
	}

	// Sets the loglevel.
	// First we need to read from args and convert it to an int
	log.SetLevel(*level)
	log.PrintTime = true
	log.Debugf("Log level is '%d'", level)

	// Check command line arguments
	if len(os.Args) < 3 {
		log.Error("Missing arguments")
	}

	// Set up source and destination directories
	source := os.Args[len(os.Args)-2]
	if !exists(source) {
		log.Errorf("Does not exist '%s'", source)
	}
	destination := os.Args[len(os.Args)-1]
	if !exists(destination) {
		log.Errorf("Does not exist '%s'", destination)
	}
	// Check if movies root is a directory
	if isFile(source) {
		log.Errorf("Is not a directory '%s'", source)
	}
	// Check if series root is a directory
	if isFile(destination) {
		log.Errorf("Is not a directory '%s'", destination)
	}

	var totalWritten int64

	log.Debugf("Overwrite is set to '%t'", overwrite)

	files, err := ioutil.ReadDir(source)
	if err != nil {
		panic(err)
	}

	index := 0

	for j, file := range files {
		log.Debugf("Checking '%s'", file.Name())

		// Check for movies
		movie, err := isMovie(file)
		if err != nil {
			panic(err)
		}
		if movie != nil {
			log.Debugf("[%d] ==== MOVIE START ==== ", j)
			log.Debugf("[%d] Movie. Title: '%s', Year: '%s', Filename: '%s'", j, movie.title, movie.year, movie.file.Name())
			log.Debugf("[%d] Movie matched regexp: '%s'", j, movie.regexp)
			
			sourcef := path.Join(source, file.Name())
			
			if *confirm == true {
				*overwrite = confirmCopy("Overwrite '"+file.Name()+"'?", true)
			}

			var written int64
			dstFolder := path.Join(destination, "Movies")

			// If "Movies" sub folder doesn't exists, then create it.
			s, err := os.Stat(path.Dir(dstFolder))
			if err != nil {
				panic(err)
			}
			if !exists(dstFolder) {
				err = os.MkdirAll(dstFolder, s.Mode())	
			}

			log.Debugf("[%d] Destination folder is '%s'", j, dstFolder)

			// Start the copy
			written, err = copy(sourcef, path.Join(dstFolder, file.Name()), *overwrite)
			if err != nil {
				log.Errorf("[%b] Can't copy '%s'. %s", j, file.Name(), err)
			}

			totalWritten = totalWritten + written

			log.Debugf("[%d] ==== MOVIE END ====", j)
		}

		// Check for series
		serie, err := isSerie(file)
		if err != nil {
			panic(err)
		}
		if serie != nil {
			log.Debugf("[%d] ==== SERIE START ====", j)
			log.Debugf("[%d] Serie. Title: '%s', Season: '%s', Episode: '%s', Filename: '%s'", j, serie.title, serie.season, serie.episode, serie.file.Name())
			log.Debugf("[%d] Serie matched regexp: '%s'", j, serie.regexp)
			
			var written int64
			
			if *confirm == true {
				*overwrite = confirmCopy("Copy '"+file.Name()+"'?", true)
			}

			sourcef := path.Join(source, file.Name())

			// Stat source so that we can perserve permissions when creating the directories if necessary
			s, err := os.Stat(path.Dir(sourcef))
			if err != nil {
				log.Errorf("[%d] Couldn't stat. '%s'", j, err)
			}

			// Create serie folder and season folders resursively
			dstFolder := path.Join(destination, serie.title, "Season "+serie.season)
			if !exists(dstFolder) {
				log.Debugf("[%s] Dest does not exist, creating '%s'", j, dstFolder)
				err = os.MkdirAll(dstFolder, s.Mode())
				if err != nil {
					log.Errorf("[%d] Couldn't create '%s'. %s", j, dstFolder, err)
				}
			}

			// Start copying
			dstf := path.Join(dstFolder, file.Name())
			log.Debugf("[%d] Dest is '%s'", j, dstf)
			written, err = copy(source, dstf, *overwrite)
			if err != nil {
				log.Errorf("[%d] Can't copy '%s'. %s", j, file.Name(), err)
			}

			totalWritten = totalWritten + written

			log.Debugf("[%d] ==== SERIE END ====", j)
		}
		index++
	}
	
	log.Printf("Copied %s\n", convertUnit(totalWritten))

}
