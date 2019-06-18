package os_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	kit "github.com/ysmood/gokit"
	"github.com/ysmood/gokit/pkg/utils"
)

func TestThisDirPath(t *testing.T) {
	assert.Regexp(t, `.+/pkg/os/fs_test.go$`, kit.ThisFilePath())
	assert.Regexp(t, `.+/pkg/os$`, kit.ThisDirPath())
}

func TestOutputString(t *testing.T) {
	p := "tmp/" + kit.RandString(10)

	_ = kit.OutputFile(p, p, nil)

	c, err := kit.ReadString(p)

	if err != nil {
		panic(err)
	}

	assert.Equal(t, c, p)
}

func TestOutputBytes(t *testing.T) {
	p := "tmp/" + kit.RandString(10)

	_ = kit.OutputFile(p, []byte("test"), nil)

	c, err := kit.ReadString(p)

	if err != nil {
		panic(err)
	}

	assert.Equal(t, c, "test")
}

func TestOutputStringErr(t *testing.T) {
	err := kit.OutputFile(".", "", nil)

	assert.EqualError(t, err, "open .: is a directory")
}

func TestOutputStringErr2(t *testing.T) {
	err := kit.OutputFile("/a/a", "", nil)
	assert.EqualError(t, err, "mkdir /a: permission denied")
}

func TestOutputJSON(t *testing.T) {
	p := "tmp/deep/" + kit.RandString(10)

	data := map[string]interface{}{
		"A": p,
		"B": 10.0,
	}

	_ = kit.OutputFile(p, data, nil)

	var ret interface{}
	err := kit.ReadJSON(p, &ret)

	if err != nil {
		panic(err)
	}

	assert.Equal(t, ret, data)
}

func TestOutputJSONErr(t *testing.T) {
	p := "tmp/deep/" + kit.RandString(10)

	err := kit.OutputFile(p, make(chan kit.Nil), nil)

	assert.EqualError(t, err, "json: unsupported type: chan utils.Nil")
}

func TestReadJSONErr(t *testing.T) {
	p := "tmp/deep/" + kit.RandString(10)

	err := kit.ReadJSON(p, nil)

	assert.Regexp(t, "no such file or directory", err.Error())
}

func TestMkdir(t *testing.T) {
	p := "tmp/deep/a/b/c"
	_ = kit.Mkdir(p, nil)

	assert.Equal(t, true, kit.DirExists(p))
}

func TestGlob(t *testing.T) {
	_ = kit.OutputFile("tmp/glob/a/b", "", nil)
	_ = kit.OutputFile("tmp/glob/a/c", "", nil)

	l, err := kit.Walk("glob/**").Dir("tmp").List()
	kit.E(err)
	assert.Equal(t, 3, len(l))
}

func TestGlobGit(t *testing.T) {
	l, err := kit.Walk("**", kit.WalkGitIgnore).List()
	kit.E(err)
	fullPath, _ := filepath.Abs("fs.go")
	assert.Contains(t, l, fullPath)
}

func TestRemove(t *testing.T) {
	_ = kit.OutputFile("tmp/remove/a", "", nil)
	_ = kit.OutputFile("tmp/remove/b/c", "", nil)
	_ = kit.OutputFile("tmp/remove/b/d", "", nil)
	_ = kit.OutputFile("tmp/remove/e/f/g", "", nil)

	kit.E(kit.Remove("tmp/remove/**"))

	l, err := kit.Walk("tmp/remove/**").List()
	kit.E(err)
	assert.Equal(t, 0, len(l))
}

func TestRemoveDir(t *testing.T) {
	_ = kit.OutputFile("tmp/remove/a", "", nil)
	_ = kit.OutputFile("tmp/remove/b/c", "", nil)
	_ = kit.OutputFile("tmp/remove/b/d", "", nil)
	_ = kit.OutputFile("tmp/remove/e/f/g", "", nil)

	kit.E(kit.Remove("tmp/remove"))

	assert.False(t, kit.DirExists("tmp/remove"))
}

func TestRemoveDirErr(t *testing.T) {
	p := "tmp/" + kit.RandString(16)
	_ = kit.OutputFile(p+"/a", "", nil)
	_ = os.Chmod(p, 0400)
	defer func() { _ = os.Chmod(p, 0700) }()

	err := kit.Remove(p)

	assert.Regexp(t, "permission denied", err.Error())
}

func TestRemoveSingleFile(t *testing.T) {
	p := "tmp/remove-single/a"
	_ = kit.OutputFile(p, "", nil)

	assert.Equal(t, true, kit.FileExists(p))

	_ = kit.Remove(p)

	assert.Equal(t, false, kit.FileExists(p))
}

func TestMove(t *testing.T) {
	p := "tmp/" + kit.RandString(10)

	_ = kit.OutputFile(p+"/a/b", "", nil)
	_ = kit.OutputFile(p+"/a/c", "", nil)

	_ = kit.Move(p+"/a", p+"/d", nil)

	assert.True(t, kit.Exists(p+"/d/b"))
	assert.True(t, kit.DirExists(p+"/d"))
}

func TestMoveErr(t *testing.T) {
	p := "tmp/" + kit.RandString(10)

	_ = kit.OutputFile(p, nil, nil)

	err := kit.Move(p+"/a", p+"/b", nil)

	assert.Regexp(t, "not a directory", err.Error())
}

func TestGoPath(t *testing.T) {
	old := os.Getenv("GOPATH")
	utils.E(os.Setenv("GOPATH", ""))
	defer func() { utils.E(os.Setenv("GOPATH", old)) }()

	s := kit.GoPath()

	assert.True(t, kit.Exists(s))
}

func TestDirExists(t *testing.T) {
	p := "tmp/" + kit.RandString(10)

	assert.Equal(t, false, kit.DirExists(p))

	_ = kit.OutputFile(p, nil, nil)

	assert.Equal(t, false, kit.DirExists(p))
}

func TestFileExists(t *testing.T) {
	assert.Equal(t, false, kit.FileExists("."))
}
