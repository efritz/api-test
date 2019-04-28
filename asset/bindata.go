// Code generated by go-bindata. (@generated) DO NOT EDIT.
// sources:
// schema/config.yaml
// schema/include.yaml
package asset

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)
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

var _schemaConfigYaml = []byte(`---

definitions:
  stringOrList:
    oneOf:
      - type: string
      - type: array
        items:
          type: string
  global-request:
    base-url:
      type: string
    auth:
      $ref: '#/definitions/auth'
    headers:
      $ref: '#/definitions/headers'
  scenario:
    type: object
    properties:
      name:
        type: string
      dependencies:
        $ref: '#/definitions/stringOrList'
      tests:
        type: array
        items:
          $ref: '#/definitions/test'
    additionalProperties: false
    required:
      - name
      - tests
  test:
    type: object
    properties:
      name:
        type: string
      request:
        $ref: '#/definitions/request'
      response:
        $ref: '#/definitions/response'
    additionalProperties: false
    required:
      - request
  request:
    type: object
    properties:
      uri:
        type: string
      method:
        type: string
        enum:
          - get
          - post
          - put
          - patch
          - delete
          - options
      headers:
        $ref: '#/definitions/headers'
      auth:
        $ref: '#/definitions/auth'
      body:
        type: string
      json-body: {}
    additionalProperties: false
    required:
      - uri
  response:
    type: object
    properties:
      status:
        anyOf:
          - type: string
          - type: number
      headers:
        $ref: '#/definitions/headers'
      body:
        type: string
      extract:
        type: string
    additionalProperties: false
  headers:
    type: object
    additionalProperties:
      $ref: '#/definitions/stringOrList'
  auth:
    type: object
    properties:
      username:
        type: string
      password:
        type: string
    additionalProperties: false
    required:
      - username
      - password

type: object
properties:
  global-request:
    $ref: '#/definitions/global-request'
  scenarios:
    type: array
    items:
      $ref: '#/definitions/scenario'
  include:
    $ref: '#/definitions/stringOrList'
additionalProperties: false
`)

func schemaConfigYamlBytes() ([]byte, error) {
	return _schemaConfigYaml, nil
}

func schemaConfigYaml() (*asset, error) {
	bytes, err := schemaConfigYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "schema/config.yaml", size: 2060, mode: os.FileMode(420), modTime: time.Unix(1556466026, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _schemaIncludeYaml = []byte(`---

definitions:
  stringOrList:
    oneOf:
      - type: string
      - type: array
        items:
          type: string
  scenario:
    type: object
    properties:
      name:
        type: string
      dependencies:
        $ref: '#/definitions/stringOrList'
      tests:
        type: array
        items:
          $ref: '#/definitions/test'
    additionalProperties: false
    required:
      - name
      - tests
  test:
    type: object
    properties:
      name:
        type: string
      request:
        $ref: '#/definitions/request'
      response:
        $ref: '#/definitions/response'
    additionalProperties: false
    required:
      - request
  request:
    type: object
    properties:
      uri:
        type: string
      method:
        type: string
        enum:
          - get
          - post
          - put
          - patch
          - delete
          - options
      auth:
        $ref: '#/definitions/auth'
      headers:
        $ref: '#/definitions/headers'
      body:
        type: string
      json-body: {}
    additionalProperties: false
    required:
      - uri
  response:
    type: object
    properties:
      status:
        anyOf:
          - type: string
          - type: number
      headers:
        $ref: '#/definitions/headers'
      body:
        type: string
      extract:
        type: string
    additionalProperties: false
  headers:
    type: object
    additionalProperties:
      $ref: '#/definitions/stringOrList'
  auth:
    type: object
    properties:
      username:
        type: string
      password:
        type: string
    additionalProperties: false
    required:
      - username
      - password

type: object
properties:
  scenarios:
    type: array
    items:
      $ref: '#/definitions/scenario'
  include:
    $ref: '#/definitions/stringOrList'
additionalProperties: false
`)

func schemaIncludeYamlBytes() ([]byte, error) {
	return _schemaIncludeYaml, nil
}

func schemaIncludeYaml() (*asset, error) {
	bytes, err := schemaIncludeYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "schema/include.yaml", size: 1858, mode: os.FileMode(420), modTime: time.Unix(1556466030, 0)}
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
	"schema/config.yaml":  schemaConfigYaml,
	"schema/include.yaml": schemaIncludeYaml,
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
	"schema": &bintree{nil, map[string]*bintree{
		"config.yaml":  &bintree{schemaConfigYaml, map[string]*bintree{}},
		"include.yaml": &bintree{schemaIncludeYaml, map[string]*bintree{}},
	}},
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
