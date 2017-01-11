package lazycopy

/**
*
* Author: Amir Mofasser <amir.mofasser@gmail.com>
*	https://github.com/amimof
*
*/

import (
  "github.com/cheggaaa/pb"
  "io/ioutil"
  "io"
  "path"
  "os"
  "time"
)

/*
 * Copies a single file or a en entire directory with it's content from src to dst. 
 * Will overwrite any file that already exists if 'ow' is set to true
 *
 */
func Copy(src, dst string, ow bool) (int64, error) {

	// Stat source
	s, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	var written int64

	// Check if src is a dir
	if s.IsDir() {
		written, err = copyDir(src, dst, ow)
	}

	// Check if src is a file
	if !s.IsDir() {
		written, err = copyFile(src, dst, ow)
	}

	return written, err

}

// Copies a dir tree from src to dst. Overwrites content if ow is true
func copyDir(src, dst string, ow bool) (int64, error) {

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

	var written int64
	entries, err := ioutil.ReadDir(src)

	if err != nil {
		return 0, err
	} else {

		for _, entry := range entries {

			sfp := path.Join(src, entry.Name())
			dfp := path.Join(dst, entry.Name())

			if entry.IsDir() {
				written, err = copyDir(sfp, dfp, ow)
				if err != nil {
					return 0, err
				}
			} else {
				written, err = copyFile(sfp, dfp, ow)
				if err != nil {
					return 0, nil
				}
			}
		}
	}

	return written, nil

}

// Copies the content of src file to dst. Overwrites dst file if ow os true
func copyFile(src, dst string, ow bool) (int64, error) {

	// Create source
	var source io.Reader
	s, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer s.Close()

	// Stat source
	srcStat, err := s.Stat()
	if err != nil {
		return 0, err
	}
	sourceSize := srcStat.Size()
	source = s

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

	// Create the progress bar
	bar := pb.New(int(sourceSize)).SetUnits(pb.U_BYTES).SetRefreshRate(time.Millisecond*10).Prefix(truncate(path.Base(src), 10, true)+": ")
	bar.Format("<.->")
	bar.ShowSpeed = true
	bar.Start()

	// Copy
	writer := io.MultiWriter(dest, bar)
	written, err := io.Copy(writer, source)
	if err != nil {
		return 0, err
	}

	bar.Finish()
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
