package util

import (
	"os"	
	"fmt"
	"path"
	"strings"
	"io/ioutil"
)

func folder_size(dir string, excludes_suffixes []string, size *AtomicInt64) error {
	stat, err := os.Stat(dir)
	if err != nil {
		return err
	}

	if stat.IsDir() {
		infos, err := ioutil.ReadDir(dir)
		if err != nil {
			return err
		}

		for _, info := range infos {
			if info.IsDir() {
				folder_size(path.Join(dir, info.Name()), excludes_suffixes, size)
				continue
			}

			hit := false
			for _,suffix := range excludes_suffixes {
				if HasSuffix(info.Name(), suffix) {
					hit = true
					break
				}				
			}
			if !hit {
				size.Caculate(info.Size())	
			} 					
		}
	} else {
		size.Caculate(stat.Size())
	}
	return nil
}

func FolderSize(folder string, excludes_suffixes []string) int64 {
	size := NewAtomicInt64(0)
	folder_size(folder, excludes_suffixes, size)
	return size.Value()
}

func ext(s string) string {
	r := strings.TrimPrefix(s, ".")
	return "." + r
}

func HasSuffix(fn, suffix string) bool {
	if strings.ToLower(path.Ext(fn)) == strings.ToLower(ext(suffix)) {
		return true
	}
	return false
}

func PathFolders(path string) []string {
	folders := []string{"/"}
	if path == "/" {
		return folders
	}
	
	path = strings.Trim(path, "/")
	paths := strings.Split(path, "/")
	if len(paths) == 1 && len(paths[0]) == 0 {
		return folders
	}
	return append(folders, paths...)
}

func NamespaceKey(uri string) (string, string) {
	if uri == "/" {
		return "/", ""
	}
	return path.Split(strings.TrimSuffix(uri, "/"))
}

func Namespace(uri string) (string, error) {
	if !strings.HasPrefix(uri, "/") {
		return "", fmt.Errorf("%s is not namespace format", uri)
	}

	if uri == "/" {
		return uri, nil
	}

	ns, _ := path.Split(strings.TrimSuffix(uri, "/"))
	return ns, nil
}

func Exist(path string) bool {
	if _, err := os.Stat(path); err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func QrcodeKey(uri string, size int) string {
	s := strings.TrimPrefix(uri, "/")
	s = strings.Replace(s, "/", "-", -1)
	s = strings.Replace(s, ".", "_", -1)
	return fmt.Sprintf("%sx%d",s, size)
}


