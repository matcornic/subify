package notif

import (
	"io/ioutil"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
)

const (
	iconURL           = "https://github.com/matcornic/subify/raw/master/images/icon.png"
	subifyTempFolder  = ".subify"
	subifyIconsFolder = "icons"
	subifyIconName    = "icon.png"
)

func getCachePath(path ...string) (string, error) {
	usr, err := user.Current()

	if err != nil {
		return "", err
	}

	cachePath := filepath.Join(usr.HomeDir, subifyTempFolder)
	if len(path) > 0 {
		cachePath = filepath.Join(cachePath, filepath.Join(path...))
	}

	return cachePath, nil
}

func existsInCache(path string) bool {
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return true
	}
	return false
}

func requestToPath(url, path string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	parent := filepath.Join(path, "..")
	err = os.MkdirAll(parent, 0777)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path, data, 0777)
	if err != nil {
		return err
	}
	return nil
}

// downloadIcon downloads the icon if needed
func downloadIcon() string {
	path, _ := getCachePath(subifyIconsFolder, "icon.png")
	if !existsInCache(path) {
		err := requestToPath(iconURL, path)
		if err != nil {
			return ""
		}
		return path
	}
	return path
}
