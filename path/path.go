package path

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func CleanPathname(s ...string) string {
	res := []byte(strings.Join(s, "_"))
	for i, c := range res {
		switch c {
		case '-':
			// remove paths that can be opts, ex '-rf'
			if i == 0 {
				res[i] = '_'
			}
		case '/':
			res[i] = ','
		case '\\', '!', '*', '?', '\'', '"', '<', '>', '|':
			res[i] = '.'
		}
	}
	return string(res)
}

func FindUpwards(filepath string) (string, error) {
	if path.IsAbs(filepath) {
		return "", fmt.Errorf("path must be relative")
	}
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return findUpwardsRec(wd, filepath)
}

func findUpwardsRec(dir, p string) (string, error) {
	testpath := path.Join(dir, p)
	file, err := os.Open(testpath)
	if err == nil {
		file.Close()
		return testpath, nil
	} else {
		err := err.(*os.PathError)
		if err.Err.Error() != "no such file or directory" {
			return "", err
		}
	}
	if dir == "/" {
		// not found
		return "", nil
	}
	return findUpwardsRec(filepath.Dir(dir), p)
}
