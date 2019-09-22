package readers

import (
	"fmt"
	"testing"
)

const resourcesdir = "../tests_resources"

var files = []string{
	fmt.Sprintf("%s/1.log", resourcesdir),
	fmt.Sprintf("%s/2.log", resourcesdir),
	fmt.Sprintf("%s/3.log", resourcesdir),
}

func Test_ReadFilesSerially_isreading(t *testing.T) {
	ch := ReadFilesSerially(files)
	expected := [3]string{"one", "beer", "curry"}
	var v [3]string
	v[0] = <-ch
	v[1] = <-ch
	v[2] = <-ch
	if v != expected {
		t.Error(
			"For", files,
			"expected", expected,
			"got", v,
		)
	}
}

func TestReadFilesSerially_isterminating(t *testing.T) {
	expectedlines := 20
	numlines := 0
	for line := range ReadFilesSerially(files) {
		if line != "" {
			numlines++
		}
	}
	if numlines != 20 {
		t.Error("Expected to read ", expectedlines, "lines ",
			" but got ", numlines,
		)
	}
}
