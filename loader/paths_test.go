package loader

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"
)

type PathsSuite struct{}

func (s *PathsSuite) TestGetConfigPath(t sweet.T) {
	name := "api-test.yaml"
	ioutil.WriteFile(name, nil, os.ModePerm)
	defer os.RemoveAll(name)

	path, err := GetConfigPath("")
	Expect(err).To(BeNil())
	Expect(path).To(Equal(name))
}

func (s *PathsSuite) TestGetConfigPathMissing(t sweet.T) {
	_, err := GetConfigPath("")
	Expect(err).To(MatchError("could not infer config file"))
}

func (s *PathsSuite) TestGetConfigPathExplicit(t sweet.T) {
	path, err := GetConfigPath("foo.yaml")
	Expect(err).To(BeNil())
	Expect(path).To(Equal("foo.yaml"))
}

func (s *PathsSuite) TestGetOverridePath(t sweet.T) {
	name := "api-test.override.yaml"
	ioutil.WriteFile(name, nil, os.ModePerm)
	defer os.RemoveAll(name)

	path, err := GetOverridePath()
	Expect(err).To(BeNil())
	Expect(path).To(Equal(name))
}

func (s *PathsSuite) TestGetOverridePathMissing(t sweet.T) {
	path, err := GetOverridePath()
	Expect(err).To(BeNil())
	Expect(path).To(BeEmpty())
}

func (s *PathsSuite) TestFindFirstToExist(t sweet.T) {
	name := buildTempDir(buildEmptyFiles([]string{"b.txt", "c.txt"}))
	defer os.RemoveAll(name)

	path, err := findFirstToExist([]string{
		filepath.Join(name, "a.txt"),
		filepath.Join(name, "b.txt"),
		filepath.Join(name, "c.txt"),
	})

	Expect(err).To(BeNil())
	Expect(path).To(Equal(filepath.Join(name, "b.txt")))
}

func (s *PathsSuite) TestFindFirstToExistMissing(t sweet.T) {
	name := buildTempDir(buildEmptyFiles([]string{}))
	defer os.RemoveAll(name)

	path, err := findFirstToExist([]string{
		filepath.Join(name, "a.txt"),
		filepath.Join(name, "b.txt"),
		filepath.Join(name, "c.txt"),
	})

	Expect(err).To(BeNil())
	Expect(path).To(BeEmpty())
}

func (s *PathsSuite) TestFindFirstToExistDirectory(t sweet.T) {
	name := buildTempDir(buildEmptyFiles([]string{"sub/a.txt"}))
	defer os.RemoveAll(name)

	_, err := findFirstToExist([]string{
		filepath.Join(name, "sub"),
	})

	Expect(err).NotTo(BeNil())
	Expect(err.Error()).To(ContainSubstring("not a file"))
}

//
// Helpers

func buildTempDir(files map[string]string) string {
	name, _ := ioutil.TempDir("", "ij-test")

	for path, content := range files {
		path = filepath.Join(name, path)
		os.MkdirAll(filepath.Dir(path), os.ModePerm)
		ioutil.WriteFile(path, []byte(content), os.ModePerm)
	}

	return name
}

func buildEmptyFiles(keys []string) map[string]string {
	files := map[string]string{}
	for _, key := range keys {
		files[key] = ""
	}

	return files
}
