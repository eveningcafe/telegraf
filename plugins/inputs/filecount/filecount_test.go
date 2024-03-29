package filecount

import (
	"os"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/influxdata/telegraf/internal"
	"github.com/influxdata/telegraf/testutil"
	"github.com/stretchr/testify/require"
)

func TestNoFilters(t *testing.T) {
	fc := getNoFilterFileCount()
	matches := []string{"foo", "bar", "baz", "qux",
		"subdir/", "subdir/quux", "subdir/quuz",
		"subdir/nested2", "subdir/nested2/qux"}
	fileCountEquals(t, fc, len(matches), 5096)
}

func TestNoFiltersOnChildDir(t *testing.T) {
	fc := getNoFilterFileCount()
	fc.Directories = []string{getTestdataDir() + "/*"}
	matches := []string{"subdir/quux", "subdir/quuz",
		"subdir/nested2/qux", "subdir/nested2"}

	tags := map[string]string{"directory": getTestdataDir() + "/subdir"}
	acc := testutil.Accumulator{}
	acc.GatherError(fc.Gather)
	require.True(t, acc.HasPoint("filecount", tags, "count", int64(len(matches))))
	require.True(t, acc.HasPoint("filecount", tags, "size_bytes", int64(600)))
}

func TestNoRecursiveButSuperMeta(t *testing.T) {
	fc := getNoFilterFileCount()
	fc.Recursive = false
	fc.Directories = []string{getTestdataDir() + "/**"}
	matches := []string{"subdir/quux", "subdir/quuz", "subdir/nested2"}

	tags := map[string]string{"directory": getTestdataDir() + "/subdir"}
	acc := testutil.Accumulator{}
	acc.GatherError(fc.Gather)

	require.True(t, acc.HasPoint("filecount", tags, "count", int64(len(matches))))
	require.True(t, acc.HasPoint("filecount", tags, "size_bytes", int64(200)))
}

func TestNameFilter(t *testing.T) {
	fc := getNoFilterFileCount()
	fc.Name = "ba*"
	matches := []string{"bar", "baz"}
	fileCountEquals(t, fc, len(matches), 0)
}

func TestNonRecursive(t *testing.T) {
	fc := getNoFilterFileCount()
	fc.Recursive = false
	matches := []string{"foo", "bar", "baz", "qux", "subdir"}

	fileCountEquals(t, fc, len(matches), 4496)
}

func TestDoubleAndSimpleStar(t *testing.T) {
	fc := getNoFilterFileCount()
	fc.Directories = []string{getTestdataDir() + "/**/*"}
	matches := []string{"qux"}

	tags := map[string]string{"directory": getTestdataDir() + "/subdir/nested2"}

	acc := testutil.Accumulator{}
	acc.GatherError(fc.Gather)

	require.True(t, acc.HasPoint("filecount", tags, "count", int64(len(matches))))
	require.True(t, acc.HasPoint("filecount", tags, "size_bytes", int64(400)))
}

func TestRegularOnlyFilter(t *testing.T) {
	fc := getNoFilterFileCount()
	fc.RegularOnly = true
	matches := []string{
		"foo", "bar", "baz", "qux", "subdir/quux", "subdir/quuz",
		"subdir/nested2/qux"}

	fileCountEquals(t, fc, len(matches), 800)
}

func TestSizeFilter(t *testing.T) {
	fc := getNoFilterFileCount()
	fc.Size = internal.Size{Size: -100}
	matches := []string{"foo", "bar", "baz",
		"subdir/quux", "subdir/quuz"}
	fileCountEquals(t, fc, len(matches), 0)

	fc.Size = internal.Size{Size: 100}
	matches = []string{"qux", "subdir/nested2//qux"}

	fileCountEquals(t, fc, len(matches), 800)
}

func TestMTimeFilter(t *testing.T) {

	mtime := time.Date(2011, time.December, 14, 18, 25, 5, 0, time.UTC)
	fileAge := time.Since(mtime) - (60 * time.Second)

	fc := getNoFilterFileCount()
	fc.MTime = internal.Duration{Duration: -fileAge}
	matches := []string{"foo", "bar", "qux",
		"subdir/", "subdir/quux", "subdir/quuz",
		"subdir/nested2", "subdir/nested2/qux"}

	fileCountEquals(t, fc, len(matches), 5096)

	fc.MTime = internal.Duration{Duration: fileAge}
	matches = []string{"baz"}
	fileCountEquals(t, fc, len(matches), 0)
}

func getNoFilterFileCount() FileCount {
	return FileCount{
		Directories: []string{getTestdataDir()},
		Name:        "*",
		Recursive:   true,
		RegularOnly: false,
		Size:        internal.Size{Size: 0},
		MTime:       internal.Duration{Duration: 0},
		fileFilters: nil,
		Fs:          getFakeFileSystem(getTestdataDir()),
	}
}

func getTestdataDir() string {
	dir, err := os.Getwd()
	if err != nil {
		// if we cannot even establish the test directory, further progress is meaningless
		panic(err)
	}

	var chunks []string
	var testDirectory string

	if runtime.GOOS == "windows" {
		chunks = strings.Split(dir, "\\")
		testDirectory = strings.Join(chunks[:], "\\") + "\\testdata"
	} else {
		chunks = strings.Split(dir, "/")
		testDirectory = strings.Join(chunks[:], "/") + "/testdata"
	}
	return testDirectory
}

func getFakeFileSystem(basePath string) fakeFileSystem {
	// create our desired "filesystem" object, complete with an internal map allowing our funcs to return meta data as requested

	mtime := time.Date(2015, time.December, 14, 18, 25, 5, 0, time.UTC)
	olderMtime := time.Date(2010, time.December, 14, 18, 25, 5, 0, time.UTC)

	// set file permisions
	var fmask uint32 = 0666
	var dmask uint32 = 0666

	// set directory bit
	dmask |= (1 << uint(32-1))

	// create a lookup map for getting "files" from the "filesystem"
	fileList := map[string]fakeFileInfo{
		basePath:                         {name: "testdata", size: int64(4096), filemode: uint32(dmask), modtime: mtime, isdir: true},
		basePath + "/foo":                {name: "foo", filemode: uint32(fmask), modtime: mtime},
		basePath + "/bar":                {name: "bar", filemode: uint32(fmask), modtime: mtime},
		basePath + "/baz":                {name: "baz", filemode: uint32(fmask), modtime: olderMtime},
		basePath + "/qux":                {name: "qux", size: int64(400), filemode: uint32(fmask), modtime: mtime},
		basePath + "/subdir":             {name: "subdir", size: int64(4096), filemode: uint32(dmask), modtime: mtime, isdir: true},
		basePath + "/subdir/quux":        {name: "quux", filemode: uint32(fmask), modtime: mtime},
		basePath + "/subdir/quuz":        {name: "quuz", filemode: uint32(fmask), modtime: mtime},
		basePath + "/subdir/nested2":     {name: "nested2", size: int64(200), filemode: uint32(dmask), modtime: mtime, isdir: true},
		basePath + "/subdir/nested2/qux": {name: "qux", filemode: uint32(fmask), modtime: mtime, size: int64(400)},
	}

	fs := fakeFileSystem{files: fileList}
	return fs

}

func fileCountEquals(t *testing.T, fc FileCount, expectedCount int, expectedSize int) {
	tags := map[string]string{"directory": getTestdataDir()}
	acc := testutil.Accumulator{}
	acc.GatherError(fc.Gather)
	require.True(t, acc.HasPoint("filecount", tags, "count", int64(expectedCount)))
	require.True(t, acc.HasPoint("filecount", tags, "size_bytes", int64(expectedSize)))
}
