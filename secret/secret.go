package secret

import (
	"io/ioutil"
	"path/filepath"
)

// Secret returns the value present in the secrets.txt file.
func Secret() string {
	// If using an absolute path just set it to it directly. For example:
	// path := "/tmp/dat"
	path, err := filepath.Abs("secret.txt")
	// This file isn't mandatory as the Mongo URL can also be added via flag. If the file isn't there, it just returns nothing.
	if err != nil {
		return ""
	}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return ""
	}
	return string(data)
}
