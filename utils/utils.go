package utils

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/google/go-github/v32/github"
)

var archs = map[string]string{
	"386":   "x86_32",
	"amd64": "x86_64",
}

// GetArch returns current arch
func GetArch() string {
	arch, ok := archs[runtime.GOARCH]
	if ok {
		return arch
	}
	return runtime.GOARCH
}

// IsSuitableAsset returns true if asset is suitable for download for current system
func IsSuitableAsset(assetName, arch, os string) bool {
	return strings.HasPrefix(assetName, "protoc") &&
		strings.Contains(assetName, fmt.Sprintf("-%s-", os)) &&
		strings.Contains(assetName, fmt.Sprintf("-%s.", arch))
}

// GetHomeDir returns home dir for app
func GetHomeDir(app string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return path.Join(home, "."+app), nil
}

// GetHomeVersionsDir returns home's versions dir for app
func GetHomeVersionsDir(app string) (string, error) {
	home, err := GetHomeDir(app)
	if err != nil {
		return "", err
	}

	return path.Join(home, "versions"), nil
}

// GetHomeVersionDir returns home's versions dir for a certain app version
func GetHomeVersionDir(app, versoin string) (string, error) {
	home, err := GetHomeVersionsDir(app)
	if err != nil {
		return "", err
	}

	return path.Join(home, versoin), nil
}

// GetHomeTmpDir returns home's tmp dir for app
func GetHomeTmpDir(app string) (string, error) {
	home, err := GetHomeDir(app)
	if err != nil {
		return "", err
	}

	return path.Join(home, "tmp"), nil
}

// GetHomeActiveDir returns home's active dir for app
func GetHomeActiveDir(app string) (string, error) {
	home, err := GetHomeDir(app)
	if err != nil {
		return "", err
	}

	return path.Join(home, "active"), nil
}

// PrepareHomeDir prepares home dir
func PrepareHomeDir(app string) error {
	fs := []func(string) (string, error){
		GetHomeDir,
		GetHomeTmpDir,
		GetHomeVersionsDir,
		GetHomeActiveDir,
	}
	for _, f := range fs {
		d, err := f(app)
		if err != nil {
			return err
		}
		if err = os.MkdirAll(d, 0755); err != nil {
			return err
		}
	}
	return nil
}

// DownloadFile will download a url to a local file. (it will
// write as it downloads and not load the whole file into memory)
func DownloadFile(url, filepath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

// Unzip will decompress a zip archive, moving all files and folders
// within the zip file (parameter 1) to an output directory (parameter 2).
func Unzip(src string, dest string) ([]string, error) {

	var filenames []string

	r, err := zip.OpenReader(src)
	if err != nil {
		return filenames, err
	}
	defer r.Close()

	for _, f := range r.File {

		// Store filename/path for returning and using later on
		fpath := filepath.Join(dest, f.Name)

		// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return filenames, fmt.Errorf("%s: illegal file path", fpath)
		}

		filenames = append(filenames, fpath)

		if f.FileInfo().IsDir() {
			// Make Folder
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		// Make File
		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return filenames, err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return filenames, err
		}

		rc, err := f.Open()
		if err != nil {
			return filenames, err
		}

		_, err = io.Copy(outFile, rc)

		// Close the file without defer to close before next iteration of loop
		outFile.Close()
		rc.Close()

		if err != nil {
			return filenames, err
		}
	}
	return filenames, nil
}

// DownloadVersion download a version if needed. Returns (false, nil) if version
// is already downloaded
func DownloadVersion(app, version string, asset *github.ReleaseAsset) (bool, error) {
	if err := PrepareHomeDir(app); err != nil {
		return false, err
	}

	versionDir, err := GetHomeVersionDir(app, version)
	if err != nil {
		return false, err
	}

	if _, err := os.Stat(versionDir); os.IsNotExist(err) {
		tmp, err := GetHomeTmpDir(app)
		if err != nil {
			return false, err
		}
		zipLocal := path.Join(tmp, *asset.Name)
		if err := DownloadFile(*asset.BrowserDownloadURL, zipLocal); err != nil {
			return false, err
		}

		if _, err = Unzip(zipLocal, versionDir); err != nil {
			return false, err
		}

		return true, nil
	}

	return false, nil
}

// ActivateVersion activates certain version
func ActivateVersion(app, version string) error {

	versionDir, err := GetHomeVersionDir(app, version)
	if err != nil {
		return err
	}

	activeDir, err := GetHomeActiveDir(app)
	if err != nil {
		return err
	}

	dirs := []string{"bin", "include"}
	for _, d := range dirs {
		v := path.Join(versionDir, d)
		l := path.Join(activeDir, d)
		if _, err := os.Lstat(l); err == nil {
			os.Remove(l)
		}
		if err := os.Symlink(v, l); err != nil {
			return err
		}
	}
	return nil
}
