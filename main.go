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

type Session struct {
	bar *pb.ProgressBar
	overwrite *bool
	confirm	*bool
	debug *bool	
	quiet *bool
	version string
	written int64
	logger *loglevel.Logger
}

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

/*
 * A wrapper method that will copy a single file or a an entire directory with it's content from src to dst. 
 */
 func (s *Session) copy(src, dst string) (int64, error) {

	var written int64

	// Stat source
	stat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	// Check if src is a dir
	if stat.IsDir() {

		// Calculate total directory size before copying so that we can pass it to the progress bar
		dirSize, err := calculateSize(src)
		if err != nil {
			return 0, err
		}

		// Create progress bar and init some defaults
		s.bar = initBar(src, dirSize)
		
		// Start progress bar and start copying
		if !*s.quiet {
			s.bar.Start()
		}
		written, err = s.copyDir(src, dst)
		
	}

	// Check if src is a file
	if !stat.IsDir() {

		// Create progress bar and init some defaults
		s.bar = initBar(src, stat.Size())
		
		// Start progress bar and start  copying
		if !*s.quiet {
			s.bar.Start()
		}
		written, err = s.copyFile(src, dst)
		
	}

	if !*s.quiet {
		s.bar.Finish()
		s.written += written
	}

	return written, err

}

/*
 * Copies a dir tree from src to dst.
 */ 
func (s *Session) copyDir(src, dst string) (int64, error) {

	var written int64

	// Stat source
	stat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	// Ensure that src is a dir
	if !stat.IsDir() {
		return 0, err
	}

	// Ensure that dst does not already exist
	d, err := os.Open(dst)
	if err != nil {
		if !os.IsNotExist(err) && *s.overwrite == false {
			return 0, err	
		}
	}
	defer d.Close()

	// Create dest dir
	err = os.MkdirAll(dst, stat.Mode())
	if err != nil {
		return 0, err
	}

	entries, err := ioutil.ReadDir(src)

	if err != nil {
		return 0, err
	} else {

		for _, entry := range entries {

			sfp := path.Join(src, entry.Name())
			dfp := path.Join(dst, entry.Name())

			if entry.IsDir() {
				w, err := s.copyDir(sfp, dfp)
				if err != nil {
					return 0, err
				}
				written += w
			} else {
				w, err := s.copyFile(sfp, dfp)
				if err != nil {
					return 0, err
				}
				written += w
			}
		}
	}

	return written, nil

}

/*
 * Copies a file from srcpath to dstpath
 */ 
func (s *Session) copyFile(srcpath, dstpath string) (int64, error) {

	// Create source
	src, err := os.Open(srcpath)
	if err != nil {
		return 0, err
	}
	defer src.Close()

	// Check if dst exists
	d, err := os.Open(dstpath)
	if !os.IsNotExist(err) {
		if *s.overwrite == false {
		  return 0, err
		}
	}
	defer d.Close()

	// Create dest
	dest, err := os.Create(dstpath)
	if err != nil {
		return 0, err
	}
	defer dest.Close()

	// Copy
	writer := io.MultiWriter(dest, s.bar)
	written, err := io.Copy(writer, src)

	if err != nil {
		return 0, err
	}

	return written, nil
}

/*
 * Returns true if path exists
 */
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

/*
 * Return true if path is a file
 */
func isFile(path string) bool {
	file, err := os.Stat(path)
	if err == nil && file.IsDir() != true {
		return true
	}
	return false
}


/*
 * Checks wether a given filename is a movie based on regexp patterns
 */
func isMovie(file os.FileInfo) (*Movie, error) {
	patterns := []string{"(.*?)\\((17[0-9][0-9]|180[0-9]|181[0-9]|18[2-9]\\d|19\\d\\d|2\\d{3}|30[0-3]\\d|304[0-8])\\)(.*)",
		"(.*?)\\[(17[0-9][0-9]|180[0-9]|181[0-9]|18[2-9]\\d|19\\d\\d|2\\d{3}|30[0-3]\\d|304[0-8])\\](.*)",
		"(.*?)\\{(17[0-9][0-9]|180[0-9]|181[0-9]|18[2-9]\\d|19\\d\\d|2\\d{3}|30[0-3]\\d|304[0-8])\\}(.*)",
		"(.*?)(17[0-9][0-9]|180[0-9]|181[0-9]|18[2-9]\\d|19\\d\\d|2\\d{3}|30[0-3]\\d|304[0-8])(.*)",
		"(.*?)(\\d{3,4}p)(.*)",
	}
	extensions := "\\.(mkv|MKV|mp4|MP4|m4p|M4P|m4v|M4V|mpg|MPG|mpeg|MPEG|mp2|MP2|mpe|MPE|mpv|MPV|3gp|3GP|nsv|NSV|f4v|F4V|f4p|F4P|f4a|F4A|f4b|F4P|vob|VOB|avi|AVI|mov|MOV|wmv|WMV|asd|ASD|flv|FLV|ogv|OGV|ogg|OGG|qt|QT|yuv|YUV|rm|RM|rmvb|RMVB)"
	
	for _, element := range patterns {
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

/*
 * Checks wether a given filename is a series based on regexp patterns
 */
func isSerie(file os.FileInfo) (*Serie, error) {
	patterns := []string{"(.*?)S(\\d{1,2})E(\\d{2})(.*)",
		"(.*?)s(\\d{1,2})e(\\d{2})(.*)",
		"(.*?)\\[?(\\d{1,2})x(\\d{2})\\]?(.*)",
		"(.*?)Season.?(\\d{1,2}).*?Episode.?(\\d{1,2})(.*)",
	}
	extensions := "\\.(mkv|MKV|mp4|MP4|m4p|M4P|m4v|M4V|mpg|MPG|mpeg|MPEG|mp2|MP2|mpe|MPE|mpv|MPV|3gp|3GP|nsv|NSV|f4v|F4V|f4p|F4P|f4a|F4A|f4b|F4P|vob|VOB|avi|AVI|mov|MOV|wmv|WMV|asd|ASD|flv|FLV|ogv|OGV|ogg|OGG|qt|QT|yuv|YUV|rm|RM|rmvb|RMVB)"

	for _, element := range patterns {
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

/*
 * Converts bytes to a more human readable format. 
 * For example 104857600 bytes is converted to '100 MB'
 */
func convertUnit(bytes int64) string {
	result := fmt.Sprintf("%d %s", bytes, "bytes")

	// Convert to kb
	if bytes >= 1024 {
		result = fmt.Sprintf("%d %s", bytes / 1024, "kB")
	}
	// Convert to mb
	if bytes >= (1024 * 1024) {
		result = fmt.Sprintf("%d %s", (bytes / 1024) / 1024, "MB")
	}
	// Convert to gb
	if bytes >= (1024 * 1024 * 1024) {
		result = fmt.Sprintf("%.1f %s", ((float64(bytes) / 1024) / 1024) / 1024, "GB")
	}
	// Convert to tb
	if bytes >= (1024 * 1024 * 1024 * 1024) {
		result = fmt.Sprintf("%.2f %s", (((float64(bytes) / 1024) / 1024) / 1024) / 1024, "TB")
	}

	return result
}

/*
 * Prompts the user for confirmation before continuing with any operation. The prompt message is 
 * provided as msg. 
 */
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



/*
 * Setup an instance of ProgressBar and apply some customization to it.
 */
func initBar(src string, size int64) *pb.ProgressBar {
	bar := pb.New64(size)
	bar.SetUnits(pb.U_BYTES)
	bar.SetRefreshRate(time.Millisecond*10)
	bar.Prefix(truncate(path.Base(src), 40)+": ")	
	bar.Format("[-> ]")
	bar.ShowSpeed =  true
	return bar
}

/* 
 * Calculate total size of a directory
 */
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

/*
 * Truncates a string to the length specified by len.
 */
func truncate(str string, length int) string {
	
	var cut int
	var newStr string

	// Dont cut more than the max number of chars in the string
	// Will throw 'slice bounds out of range' otherwise
	strLen := len([]rune(str))
	if length >= strLen {
		cut = strLen
	} else {
		cut = length
	}

	newStr = str[:cut]
  returnStr := newStr + "..." + str[(strLen-3):strLen]

  // If the new string is longer than the originial, the just return the originial string
  if returnStr >= str {
    returnStr = str
  }

	return returnStr
}


func (s *Session) infof(format string, msg ...interface{}) {
	if *s.quiet {
		return
	}
	s.logger.Printf(format, msg...)
}

func (s *Session) debugf(format string, msg ...interface{}) {
	if *s.quiet {
		return
	}
	if *s.debug {
		s.logger.Debugf(format, msg...)
	}
}

func (s *Session) errorf(format string, msg ...interface{}) {
	if *s.quiet {
		return
	}
	s.logger.Errorf(format, msg...)	
}

func main() {
	
	// Read arguments
	session := &Session{
		overwrite: flag.Bool("o", false, "Overwrite existing files/folders when copying"),
		confirm: flag.Bool("c", false, "Prompt for confirm before overwriting existing files/folders"),
		debug: flag.Bool("d", false, "Debug mode"),
		quiet: flag.Bool("q", false, "Supress all output"),
		logger: loglevel.New(),
		version: "1.0.2",
	}

	printVer := flag.Bool("v", false, "Print version info")

	flag.Parse()

	// Print version info and exit
	if *printVer {
		fmt.Println(session.version)
		os.Exit(0)
	}

	// Setup logging
	session.logger.SetLevel(3)

	// Check command line arguments
	if len(os.Args) < 3 {
		session.errorf("Not enough arguments")
	}

	// Set up source and destination directories
	source := os.Args[len(os.Args)-2]
	if !exists(source) {
		session.errorf("Does not exist '%s'", source)
	}
	destination := os.Args[len(os.Args)-1]
	if !exists(destination) {
		session.errorf("Does not exist '%s'", destination)
	}
	// Check if movies root is a directory
	if isFile(source) {
		session.errorf("Is not a directory '%s'", source)
	}
	// Check if series root is a directory
	if isFile(destination) {
		session.errorf("Is not a directory '%s'", destination)
	}

	session.debugf("Overwrite is set to '%t'", *session.overwrite)

	files, err := ioutil.ReadDir(source)
	if err != nil {
		panic(err)
	}

	index := 0

	for j, file := range files {
		session.debugf("Checking '%s'", file.Name())

		// Check for movies
		movie, err := isMovie(file)
		if err != nil {
			panic(err)
		}
		if movie != nil {
			session.debugf("[%d] ==== MOVIE START ==== ", j)
			session.debugf("[%d] Movie. Title: '%s', Year: '%s', Filename: '%s'", j, movie.title, movie.year, movie.file.Name())
			session.debugf("[%d] Movie matched regexp: '%s'", j, movie.regexp)
			
			sourcef := path.Join(source, file.Name())
			
			if *session.confirm == true {
				*session.overwrite = confirmCopy("Overwrite '"+file.Name()+"'?", true)
			}

			dstFolder := path.Join(destination, "Movies")

			// If "Movies" sub folder doesn't exists, then create it.
			s, err := os.Stat(path.Dir(dstFolder))
			if err != nil {
				panic(err)
			}
			if !exists(dstFolder) {
				err = os.MkdirAll(dstFolder, s.Mode())	
			}

			session.debugf("[%d] Destination folder is '%s'", j, dstFolder)

			// Start the copy
			_, err = session.copy(sourcef, path.Join(dstFolder, file.Name()))
			if err != nil {
				session.errorf("[%b] Can't copy '%s'. %s", j, file.Name(), err)
			}
			session.debugf("[%d] ==== MOVIE END ====", j)
		}

		// Check for series
		serie, err := isSerie(file)
		if err != nil {
			panic(err)
		}
		if serie != nil {
			session.debugf("[%d] ==== SERIE START ====", j)
			session.debugf("[%d] Serie. Title: '%s', Season: '%s', Episode: '%s', Filename: '%s'", j, serie.title, serie.season, serie.episode, serie.file.Name())
			session.debugf("[%d] Serie matched regexp: '%s'", j, serie.regexp)
			
			sourcef := path.Join(source, file.Name())

			if *session.confirm == true {
				*session.overwrite = confirmCopy("Copy '"+file.Name()+"'?", true)
			}

			// Stat source so that we can perserve permissions when creating the directories if necessary
			s, err := os.Stat(path.Dir(sourcef))
			if err != nil {
				session.errorf("[%d] Couldn't stat. '%s'", j, err)
			}

			// Create serie folder and season folders resursively
			dstFolder := path.Join(destination, "Series", serie.title, "Season "+serie.season)
			if !exists(dstFolder) {
				session.debugf("[%s] Dest does not exist, creating '%s'", j, dstFolder)
				err = os.MkdirAll(dstFolder, s.Mode())
				if err != nil {
					session.debugf("[%d] Couldn't create '%s'. %s", j, dstFolder, err)
				}
			}

			// Start copying
			dstf := path.Join(dstFolder, file.Name())
			session.debugf("[%d] Dest is '%s'", j, dstf)
			_, err = session.copy(sourcef, dstf)
			if err != nil {
				session.errorf("[%d] Can't copy '%s'. %s", j, file.Name(), err)
			}

			session.debugf("[%d] ==== SERIE END ====", j)
		}
		index++
	}
	
	session.infof("Copied %s\n", convertUnit(session.written))

}
