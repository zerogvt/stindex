package jobs

import (
	"fmt"
	"log"
	"os"
	"stindex/trees"
	"testing"
)

const resourcesdir = "../tests_resources"

func Test_createJob(t *testing.T) {
	outfile := "test_out.log"
	os.Remove(outfile)
	var ramsize uint = 1
	files := []string{resourcesdir + "/1.log",
		resourcesdir + "/2.log",
		resourcesdir + "/3.log"}
	os.Args = []string{"myself",
		fmt.Sprintf("--ramsize=%d", ramsize),
		"-outfile=" + outfile,
		files[0], files[1], files[2]}
	job := CreateJob()
	if job.Outfile != outfile {
		t.Error("Expected outfile: ", outfile, "but got", job.Outfile)
	}
	if job.RAMSize != ramsize {
		t.Error("Expected ramsize: ", ramsize, "but got", job.RAMSize)
	}
	for i, f := range files {
		if job.Files[i] != f {
			t.Error("Expected file: ", f, "but got", job.Files[i])
		}
	}
}

func Test_processSearchTerms(t *testing.T) {
	infile := resourcesdir + "/1.log"
	var incount uint64 = 5
	job := &Job{Log: log.New(os.Stdout, "", log.Ldate|log.Ltime),
		RAMSize: 10,
		Files:   []string{infile}}
	job.Run()
	if job.TermsNum != incount {
		t.Error("Expected terms count of: ", incount, "but got", job.TermsNum)
	}
	exp := []trees.SearchTerm{
		trees.SearchTerm{Term: "aeon", Times: 1},
		trees.SearchTerm{Term: "beer", Times: 2},
		trees.SearchTerm{Term: "curry", Times: 1},
		trees.SearchTerm{Term: "one", Times: 1},
	}
	i := 0
	for n := range job.TermsTree.TraverseInOrder() {
		if n.Compare(exp[i]) != 0 {
			t.Error("Expected term: ", exp[i].ToStr(),
				"but got", n.ToStr())
		}
		i++
	}

}
