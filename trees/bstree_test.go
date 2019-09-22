package trees

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func Test_BST_Insertions(t *testing.T) {
	tree := &BST{}
	tree.Insert("c").
		Insert("d").
		Insert("b").
		Insert("e").
		Insert("a")
	expected := []SearchTerm{
		SearchTerm{"a", 1},
		SearchTerm{"b", 1},
		SearchTerm{"c", 1},
		SearchTerm{"d", 1},
		SearchTerm{"e", 1},
	}
	i := 0
	for node := range tree.TraverseInOrder() {
		if node != expected[i] {
			t.Error(
				"Expected search term: ", expected[i].ToStr(),
				"but got", node.ToStr())
		}
		i++
	}
}

func Test_BST_MultipleSameInsertions(t *testing.T) {
	tree := &BST{}
	tree.Insert("b").
		Insert("a").
		Insert("e").
		Insert("e").
		Insert("e")
	expected := []SearchTerm{
		SearchTerm{"a", 1},
		SearchTerm{"b", 1},
		SearchTerm{"e", 3},
	}
	i := 0
	for node := range tree.TraverseInOrder() {
		if node != expected[i] {
			t.Error(
				"Expected search term: ", expected[i].ToStr(),
				"but got", node.ToStr(),
			)
		}
		i++
	}
}

func Test_BST_UnzipToFile(t *testing.T) {
	const testfile = "test_unzip.log"
	tree := &BST{}
	tree.Insert("b").
		Insert("a").
		Insert("e").
		Insert("e").
		Insert("e")
	os.Remove(testfile)
	tree.UnzipToFile(testfile)
	file, err := os.Open(testfile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	actual, err := ioutil.ReadAll(file)
	expected := []byte(fmt.Sprintf("a\nb\ne\ne\ne\n"))
	os.Remove(testfile)
	if !bytes.Equal(actual, expected) {
		t.Error(
			"Expected", expected,
			"but got", actual,
		)
	}
}
