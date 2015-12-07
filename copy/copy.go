package lazycopy

/**
*
* Author: Amir Mofasser <amir.mofasser@gmail.com>
*	https://github.com/amimof
*
*/

import (
  logger "github.com/amimof/lazycopy/logger"
  "github.com/cheggaaa/pb"
  "io/ioutil"
  "io"
  "path"
  "os"
  "time"
)

var log *logger.Logger = logger.SetupNew("COPY")

// Copies a dir tree from src to dst. Overwrites content if ow is true
func CopyDir(src, dst string, ow bool) (int64, error) {

	// Stat source
	s, err := os.Stat(src)
	if err != nil {
		log.Error("Couldn't stat", src, err)
		return 0, err
	}

	// Ensure that src is a dir
	if !s.IsDir() {
		log.Error("Source is not a directory")
		return 0, err
	}

	// Ensure that dst does not already exist
	d, err := os.Open(dst)
  if err != nil {
    log.Error("Couldn't open", dst, err)
  }
  if ow == false {
    if !os.IsNotExist(err) {
      return 0, err
    }
  }
	defer d.Close()

	// Create dest dir
	err = os.MkdirAll(dst, s.Mode())
	if err != nil {
		log.Error("Couldn't create", dst, err)
		return 0, err
	}

	var written int64
	entries, err := ioutil.ReadDir(src)
	if err != nil {
		log.Error("Couldn't read", src, err)
		return 0, err
	} else {
		for _, entry := range entries {

			sfp := path.Join(src, entry.Name())
			dfp := path.Join(dst, entry.Name())

			if entry.IsDir() {
				written, err = CopyDir(sfp, dfp, ow)
				if err != nil {
					log.Error("Couldn't copy", err)
					return 0, err
				}
			} else {
				written, err = CopyFile(sfp, dfp, ow)
				if err != nil {
					log.Error("Couldn't copy", err)
					return 0, nil
				}
			}
		}
	}

	return written, nil

}

// Copies the content of src file to dst. Overwrites dst file if ow os true
func CopyFile(src, dst string, ow bool) (int64, error) {

	// Create source
	var source io.Reader
	s, err := os.Open(src)
	if err != nil {
		log.Error("Couldn't open", src, err)
		return 0, err
	}
	defer s.Close()

	// Stat source
	srcStat, err := s.Stat()
	if err != nil {
		log.Error("Couldn't stat", src, err)
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
		log.Error("Couldn't create", dst, err)
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
		log.Error("Couldn't copy", err)
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
