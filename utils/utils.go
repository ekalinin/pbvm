package utils

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"

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
func DownloadVersion(app, version string, asset *github.ReleaseAsset, d func(ms ...interface{})) (bool, error) {
	d(" ... preparing home ...")
	if err := PrepareHomeDir(app); err != nil {
		return false, err
	}

	d(" ... checking if installed ...")
	installed, versionDir, err := IsInstalledVersion(app, version)
	if err != nil {
		return false, err
	}

	if !installed {
		d(" ... not installed :(")
		tmp, err := GetHomeTmpDir(app)
		if err != nil {
			return false, err
		}
		zipLocal := path.Join(tmp, *asset.Name)
		d(" ... checking if zip already downloaded ...")
		if _, err := os.Stat(zipLocal); os.IsNotExist(err) {
			d(" ... zip was not downloaded yet :( downloading ")
			if err := DownloadFile(*asset.BrowserDownloadURL, zipLocal); err != nil {
				return false, err
			}
		} else {
			d(" ... zip was downloaded already :) ")
		}

		d(" ... unzipping ... ")
		if _, err = Unzip(zipLocal, versionDir); err != nil {
			return false, err
		}

		d(" ... done. ")
		return true, nil
	} else {
		d(" ... installed :)")
	}

	return false, nil
}

// IsInstalledVersion returns true (first result) if version is installed
func IsInstalledVersion(app, version string) (bool, string, error) {
	versionDir, err := GetHomeVersionDir(app, version)
	if err != nil {
		return false, versionDir, err
	}

	if _, err := os.Stat(versionDir); os.IsNotExist(err) {
		return false, versionDir, nil
	}

	return true, versionDir, nil
}

// InstalledVersion describes installed version
type InstalledVersion struct {
	Version string
	Date    time.Time
	Active  bool
}

// ListInstalledVersions returns a slice of installed versions
func ListInstalledVersions(app string) ([]InstalledVersion, error) {
	versionsDir, err := GetHomeVersionsDir(app)
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(versionsDir); os.IsNotExist(err) {
		return nil, nil
	}

	files, err := ioutil.ReadDir(versionsDir)
	if err != nil {
		return nil, err
	}

	// Desc order by name
	sort.Slice(files, func(i, j int) bool { return files[i].Name() > files[j].Name() })
	res := []InstalledVersion{}
	for _, f := range files {
		if !f.IsDir() {
			continue
		}
		active, err := IsActiveVersion(app, f.Name())
		if err != nil {
			return nil, err
		}
		res = append(res, InstalledVersion{
			Version: f.Name(),
			Date:    f.ModTime(),
			Active:  active,
		})
	}

	return res, nil
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

// IsActiveVersion returns bool if version is active
func IsActiveVersion(app, version string) (bool, error) {
	activeVer, err := GetActiveVersion(app)
	if err != nil {
		return false, err
	}

	if activeVer == version {
		return true, nil
	}

	return false, nil
}

// GetActiveVersion returns active version
func GetActiveVersion(app string) (string, error) {
	activeDir, err := GetHomeActiveDir(app)
	if err != nil {
		return "", err
	}

	l, err := os.Readlink(path.Join(activeDir, "bin"))
	if err != nil {
		return "", err
	}
	// converting
	// "/home/user/.pbvm/versions/v3.12.3/bin" -> "v3.12.3"
	return path.Base(path.Dir(l)), nil
}

// FilterAsset finds an asset which is need to be downloaded
func FilterAsset(release *github.RepositoryRelease) *github.ReleaseAsset {
	arch := GetArch()
	for _, a := range release.Assets {
		if !IsSuitableAsset(*a.Name, arch, runtime.GOOS) {
			continue
		}
		return a
	}
	return nil
}

// DeleteVersion deletes version
func DeleteVersion(app, version string) error {
	versionDir, err := GetHomeVersionDir(app, version)
	if err != nil {
		return err
	}

	if err := os.RemoveAll(versionDir); err != nil {
		return err
	}

	return nil
}
