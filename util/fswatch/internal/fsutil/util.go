package fsutil

import (
	"os"
	"path"
)

func Exists(filePath string) (exists bool, err error) {
	//判断文件是否存在
	_, err = os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			err = nil
			exists = false
		} else {
			exists = true
		}
	} else {
		exists = true
	}
	return
}

func TouchFile(filePath string) (err error) {
	//判断文件是否存在
	_, err = os.Stat(filePath)

	if err != nil && os.IsNotExist(err) {
		dir, _ := path.Split(filePath)
		if len(dir) > 0 {
			if err = os.MkdirAll(dir, os.ModePerm); err != nil {
				return
			}
		}
		var pf *os.File
		if pf, err = os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, os.ModePerm); err != nil {
			return
		} else {
			defer pf.Close()
		}
	}
	return
}
