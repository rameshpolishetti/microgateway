// Code generated by go-bindata. DO NOT EDIT. @generated
// sources:
// schema.json
package schema

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _schemaJson = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xec\x58\xcd\x8e\x9b\x30\x10\xbe\xef\x53\x20\x6f\x4f\x55\x77\x69\xa5\x9e\xf2\x06\x3d\x54\x5d\xb5\xc7\x6a\x0f\x0e\x0c\xc4\x2b\xb0\xbd\xe3\xa1\x55\x54\xe5\xdd\x2b\xfe\x1c\x08\x60\x87\x40\xfa\x23\x85\x43\x0e\xfe\x3e\x3c\xf6\xfc\x7c\x33\xe4\xd7\x5d\x10\x04\x01\x7b\x63\xa2\x1d\xe4\x9c\x6d\x02\xb6\x23\xd2\x9b\x30\x7c\x31\x4a\x3e\xd4\xab\x8f\x0a\xd3\x30\x46\x9e\xd0\xc3\xfb\x8f\x61\xbd\x76\xcf\xde\x35\x6f\x22\x24\xe5\x6b\xf7\x61\x0c\x89\x90\x82\x84\x92\x26\xfc\x2c\x22\x54\x29\x27\xf8\xc9\xf7\x2d\xb3\x83\xb3\x4d\x50\x1b\xae\x80\x1e\xb9\x8b\x54\x28\xc2\x6b\x21\x10\x62\xb6\x09\xbe\xf7\x90\x0a\x95\x3c\x87\xc6\x40\x6f\xdd\x10\x68\xc3\x7a\xeb\xcf\x7d\x1a\xd3\xa8\x34\x20\x09\x30\x03\xa3\xc7\xad\xc7\x90\x0a\xa5\xbd\x2e\x51\x66\x08\x85\x4c\xd9\x80\x74\x18\x39\x14\x82\xd1\x4a\x9a\x09\x83\x15\x45\x10\xe4\xd3\x70\xb0\x28\x56\xe3\xbb\x8d\xc6\xef\x6b\x73\xd2\xe1\xbd\x26\xee\xd6\xf3\x09\x47\xe4\xfb\xf3\x5c\x62\x00\x7f\x88\xe8\x3f\xf0\xc8\xb7\xfa\xa0\xf3\x1c\x52\x48\xf1\x5a\xc0\xa7\xe6\x06\x84\x05\xac\xe8\xb9\x2a\xc3\xff\x79\xb7\x11\xe8\x79\x3e\xcb\x85\x6c\x1d\xf6\xe1\x42\x6f\xdd\x39\xec\x30\x1e\xc7\xd5\xd9\x78\xf6\xd4\x15\x80\x84\x67\xe6\x24\x3a\xd6\x8e\xda\xbe\x40\x44\x47\x43\x9d\x2d\xd9\x97\x82\x74\x41\x73\x65\x2b\xe6\xc4\x97\xa9\x53\xa4\x62\x87\x3a\x4d\x5c\xd2\x9f\x81\xc3\xb3\x1e\xf7\xac\xbc\xed\xc8\x80\xad\x52\x19\x70\xe9\xa2\x08\x49\x90\x02\xba\x28\xb2\xc8\xb7\x3e\x46\x96\xb9\xf0\x26\x5c\x0e\xc6\x94\x68\x97\xcf\xf3\x59\xd5\x57\x05\xf0\xe6\xfe\x31\x7c\x75\xf7\xff\xc9\x72\xb6\x8d\x6f\x66\x41\x03\xa2\xc2\x65\x15\x5d\x6f\xe1\x1d\x38\xda\x30\x9f\x95\xa6\x22\x59\x77\x82\x51\xe3\x6a\x67\xf1\xd5\x5a\xca\x44\x3b\x69\xd4\xf6\xef\xe6\x48\x3b\x0a\xac\x34\xaa\x96\x17\x5d\x94\x38\x31\x98\x08\x85\x2e\x6f\xb8\x6e\xb4\xaf\x31\x01\xaf\x9c\x90\x06\x88\x84\x4c\x1d\x83\x90\xe6\x44\x80\xf2\xc9\xed\x44\x4b\x7f\x7c\xeb\xc4\x83\x8b\xe4\xfd\xf4\x8a\xd3\x32\x7f\xb4\xe1\x91\x7b\x4b\xf4\xcb\xbe\xa5\xfa\xe5\xdf\x52\xbd\x6d\xa0\xc3\x74\xb6\x03\xcb\xf3\xb6\x05\xcb\x74\xb5\x87\xf6\x19\xb6\x89\xf6\x39\x5c\xf4\xc9\x72\x5a\xef\xe3\xbb\x5d\x5b\x58\xca\x61\x79\xa6\xaa\x98\xb1\xef\x92\xb9\x0a\xb2\xe3\x99\x43\xd4\x2f\xa9\xcb\xb5\x1b\x8f\x90\xce\xbe\x73\x2b\x72\x4b\xbd\x15\xf9\xc0\x96\xb7\xc8\xa7\xff\x98\x58\x90\xc5\x57\x57\x8e\xbb\xfa\xf7\xf0\x3b\x00\x00\xff\xff\x9b\xe3\x98\x29\xc4\x13\x00\x00")

func schemaJsonBytes() ([]byte, error) {
	return bindataRead(
		_schemaJson,
		"schema.json",
	)
}

func schemaJson() (*asset, error) {
	bytes, err := schemaJsonBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "schema.json", size: 5060, mode: os.FileMode(420), modTime: time.Unix(1540500290, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"schema.json": schemaJson,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}
var _bintree = &bintree{nil, map[string]*bintree{
	"schema.json": &bintree{schemaJson, map[string]*bintree{}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}

