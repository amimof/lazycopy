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
  "path/filepath"
  "os"
  "time"
)

var (
	bar *pb.ProgressBar
	written int64
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

	var writtenTotal int64

	// Check if src is a dir
	if s.IsDir() {

		// Calculate total directory size before copying so that we can pass it to the progress bar
		dirSize, err := calculateSize(src)
		if err != nil {
			return 0, err
		}

		// Create progress bar and init some defaults
		bar = pb.New64(dirSize)
		initBar(src)

		// Start progress bar and start copying
		bar.Start()
		writtenTotal, err = copyDir(src, dst, ow)
		
	}

	// Check if src is a file
	if !s.IsDir() {

		// Create progress bar and init some defaults
		bar = pb.New64(s.Size())
		initBar(src)

		// Start progress bar and start  copying
		bar.Start()
		writtenTotal, err = copyFile(src, dst, ow)
		
	}

	bar.Finish()

	return writtenTotal, err

}

// Customize the bar. This function is mainly so that we don't have to write the same code twise.
// Don't like the fact that this function takes the argument src. Might need some work.
func initBar(src string) {
	bar.SetUnits(pb.U_BYTES)
	bar.SetRefreshRate(time.Millisecond*10)
	bar.Prefix(truncate(path.Base(src), 10, true)+": ")	
	bar.Format("<.- >")
	bar.ShowSpeed =  true
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
