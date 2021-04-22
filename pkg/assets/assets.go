package assets

import (
	"embed"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/bmatcuk/doublestar/v3"
	"github.com/brad-jones/dotfiles/pkg/utils"
	"github.com/brad-jones/goerr/v2"
	"github.com/brad-jones/goprefix/v2/pkg/colorchooser"
)

//go:embed files
var FS embed.FS

// FSItemType is used to filter results from the glob pattern so
// we can apply permissions to either files or folders; or both.
type FSItemType int

const (
	FILE FSItemType = iota
	FOLDER
	ANY
)

type Permission struct {
	// Glob is a "github.com/bmatcuk/doublestar/v3" glob
	Glob string

	// What filesystem items found by the glob will this permission apply to?
	PathType FSItemType

	// The permission to set against the filtered filesystem items
	Mode fs.FileMode
}

// GetPermission returns the FileMode for a given path based
// on the data contained in the Permissions array. see: permissions.go
func GetPermission(path string, pathType FSItemType) fs.FileMode {
	path = translatePathRev(path)
	for _, p := range Permissions {
		match, err := doublestar.Match(replaceHomeDir(p.Glob), path)
		goerr.Check(err, "glob pattern is not valid", p.Glob)
		if match && pathType == p.PathType {
			return p.Mode
		}
	}
	return 0755 // https://stackoverflow.com/questions/23842247
}

// ReadFile will return the contents of an embedded file given the path to it.
// Note that the path is passed through translatePath for you so you can supply
// the path as it would be normally.
func ReadFile(filename string) []byte {
	dat, err := FS.ReadFile(translatePath(filename))
	goerr.Check(err, "failed to read file from embedded assets")
	return dat
}

// WriteFile will write an embedded file, given by the src path to the real
// filesystem, given by the dst path. Parent folders will be created as required.
// Permissions will be sourced from the GetPermission function for each path
// that is created.
func WriteFile(src, dst string) {
	parent := ""

	for _, part := range strings.Split(replaceWinSlashes(filepath.Dir(dst)), "/") {
		if part == "" {
			parent = "/"
			continue
		}
		if parent == "" || parent == "/" {
			parent = parent + part
		} else {
			parent = fmt.Sprintf("%s/%s", parent, part)
		}
		translatedParent := translatePathRev(parent)
		if !utils.FolderExists(translatedParent) {
			perms := GetPermission(parent, FOLDER)
			goerr.Check(
				os.Mkdir(translatedParent, perms),
				"failed to create parent folder", dst,
			)
			fmt.Println(colorchooser.Sprint("created-folder"), "|", translatedParent, perms)
		}
	}

	content := ReadFile(src)
	translatedDst := translatePathRev(dst)
	if utils.FileExists(translatedDst) {
		if utils.Sha256HashContent(content) == utils.Sha256HashFile(translatedDst) {
			fmt.Println(colorchooser.Sprint("hashes-match"), "|", translatedDst)
			return
		}
	}

	perms := GetPermission(dst, FILE)
	goerr.Check(
		ioutil.WriteFile(translatedDst, content, perms),
		"failed to write embedded asset to filesystem",
	)
	fmt.Println(colorchooser.Sprint("written-file"), "|", translatedDst, perms)
}

// WriteFileToHome does the same thing as WriteFile but does not allow for
// customization fo the destination path, the destination is relative to the
// the users home directory.
func WriteFileToHome(filename string) {
	WriteFile(filename, filepath.Join(utils.HomeDir(), filename))
}

// WriteFolder will recursively write an entire folder of embedded files,
// given by the src path to the real filesystem relative the given dst path.
func WriteFolder(src, dst string) {
	items, err := FS.ReadDir(translatePath(src))
	goerr.Check(err, "failed to read dir from embedded assets")
	for _, item := range items {
		newSrc := filepath.Join(src, item.Name())
		newDst := filepath.Join(dst, item.Name())
		if item.IsDir() {
			WriteFolder(newSrc, newDst)
		} else {
			WriteFile(newSrc, newDst)
		}
	}
}

// WriteFolderToHome does the same thing as WriteFolder but does not allow for
// customization fo the destination path, the destination is relative to the
// the users home directory.
func WriteFolderToHome(dir string) {
	WriteFolder(dir, filepath.Join(utils.HomeDir(), dir))
}

// translatePath solves the following for us: If a pattern names a directory,
// all files in the subtree rooted at that directory are embedded (recursively),
// except that files with names beginning with ‘.’ or ‘_’ are excluded.
func translatePath(in string) string {
	out := "files"

	for _, part := range strings.Split(replaceWinSlashes(in), "/") {
		if strings.HasPrefix(part, ".") {
			part = strings.Replace(part, ".", "dot_", 1)
		}
		if strings.HasPrefix(part, "_") {
			part = strings.Replace(part, "_", "underscore_", 1)
		}
		out = fmt.Sprintf("%s/%s", out, part)
	}

	return out
}

// translatePathRev does the opposite of translatePath so we can
// write the correct file & folder names to the real filesystem.
func translatePathRev(in string) string {
	out := ""

	for _, part := range strings.Split(replaceWinSlashes(in), "/") {
		if part == "" {
			out = "/"
			continue
		}
		if strings.HasPrefix(part, "dot_") {
			part = strings.Replace(part, "dot_", ".", 1)
		}
		if strings.HasPrefix(part, "underscore_") {
			part = strings.Replace(part, "underscore_", "_", 1)
		}
		if out == "" || out == "/" {
			out = out + part
		} else {
			out = fmt.Sprintf("%s/%s", out, part)
		}
	}

	return out
}

// replaces the tilde with the actual home dir
func replaceHomeDir(in string) string {
	if strings.HasPrefix(in, "~") {
		return strings.Replace(in, "~", replaceWinSlashes(utils.HomeDir()), 1)
	}
	return in
}

// normalizes slashes
func replaceWinSlashes(in string) string {
	return strings.ReplaceAll(in, "\\", "/")
}
