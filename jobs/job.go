package jobs

import (
	"flag"
	"fmt"
	"log"
	"os"
	"stindex/readers"
	"stindex/trees"
	"time"
)

const version string = "0.1"
const myname string = "stindex"
const defaultRAMSize uint = 1000000

// usage returns the usage help string
func usage() string {
	return "NAME\n" +
		"\t" + myname + " -- search terms index utility\n\n" +
		"VERSION\n" +
		"\t" + version + "\n\n" +
		"SYNOPSIS\n" +
		"\t" + myname + " [-h] -ramsize -outfile [file ...]\n\n" +
		"\t" + myname + " is a small program that will process a series of\n" +
		"\tinput files containing search terms (one term per line)\n" +
		"\tcreating a final file where all the terms will be alphabetically\n" +
		"\tsorted\n\n" +
		"Mandatory parameters:\n\n" +
		"-outfile\tthe file where the output result will be\n" +
		"\t\twritten. The file must not exist.\n\n" +
		"Optional parameters:\n\n" +
		"-ramsize\tan integer greater than 1 signifying the number\n" +
		"\t\tof search terms that can be hosted in program's logical RAM\n" +
		"\t\t(defaults to " + fmt.Sprintf("%d", defaultRAMSize) + ")\n\n" +
		"-h\toutput this usage message and exit\n\n" +
		"EXIT STATUS\n" +
		"\tExits 0 on success and >0 if an error occurs\n\n" +
		"EXAMPLES\n" +
		"\tProcess all files under path resources/ with a virtual RAM of size 100\n" +
		"\tand save results in file out.log:\n" +
		"\t\t" + myname + " -ramsize=100 -outfile=out.log resources/*\n" +
		"AUTHORS\n" +
		"\tzerogvt (vas)\n" +
		"\t9 July 2019\n\n"
}

// Job struct hosts all info needed for a particular job plus
// some interim data structures and stats that get updated as the
// job goes along.
type Job struct {
	Log            *log.Logger // a logger object that will handle all logging for this Job
	RAMSize        uint        // logical simulation of a limited RAM as per problem statement
	Files          []string    // input files containg the lists transpired search terms
	TermsTree      *trees.BST  //interim data structure keeping processed terms
	TermsNum       uint64      // how many terms we processed
	Outfile        string      // output file
	StartProcessTS time.Time   // timestamp of when we start processing into interim data structs
	EndProcessTS   time.Time   // matching end timestamp
	StartUnzipTS   time.Time   //timestamp of when we start unzipping interim data structs to output file
	EndUnzipTS     time.Time   // matching end timestamp
}

// Run is driving the processing of search terms according to the current configuration.
func (job *Job) Run() {
	// start processing
	job.Log.Print("[INFO] RAM Size: ", job.RAMSize)
	job.Log.Print("[INFO] Input files: ", job.Files)
	job.Log.Print("[INFO] Starting processing.")
	job.StartProcessTS = time.Now()
	// "RAM" will be simulated by a slice of RAMSize capacity
	ram := make([]string, 0, job.RAMSize)
	// Get a channel that will provide the search terms from logs
	ch := readers.ReadFilesSerially(job.Files)
	// Get a BST which will be our interim data structure
	job.TermsTree = &trees.BST{}
	// While there is more search terms coming...
	for moretoread := true; moretoread; {
		// ...read as much terms as RAM allows in one go
		var i uint
		for i = 0; moretoread && i < job.RAMSize; i++ {
			term, ok := <-ch
			if ok {
				ram = append(ram, term)
				job.TermsNum++
			}
			moretoread = ok
		}
		// ...insert them in interim data structure
		for _, term := range ram {
			job.TermsTree.Insert(term)
		}
		// ...and clean up RAM
		ram = nil
		fmt.Print(".")
	}
	job.EndProcessTS = time.Now()
	fmt.Println()
	job.Log.Print("[INFO] ",
		"Processed ", job.TermsNum, " search terms.",
		" Elapsed time: ", job.EndProcessTS.Sub(job.StartProcessTS))
}

// PersistResults unzips the binary search tree into the output file
// defined in the current Job context
func (job *Job) PersistResults() {
	job.StartUnzipTS = time.Now()
	job.TermsTree.UnzipToFile(job.Outfile)
	job.EndUnzipTS = time.Now()
	job.Log.Print("[INFO] ",
		"Unzipped to destination file: ", job.Outfile,
		" Elapsed time: ", job.EndUnzipTS.Sub(job.StartUnzipTS))
}

// CreateJob reads and parses user input and if valid returns a Job
// object that describes the configuration of the job.
func CreateJob() *Job {
	help := false
	flag.BoolVar(&help, "help", false, "Usage")
	flag.BoolVar(&help, "h", false, "Usage")
	RAMSizeUsage := "RAM size must be a positive number."
	RAMSize := flag.Uint("ramsize", defaultRAMSize, RAMSizeUsage)

	outFileUsage := "You must provide an ouput destination file" +
		" via -outfile flag (e.g. -outfile=out.log)"
	outFile := flag.String("outfile", "", outFileUsage)
	flag.Parse()
	if help {
		fmt.Println(usage())
		os.Exit(0)
	}
	job := &Job{Log: log.New(os.Stdout, "", log.Ldate|log.Ltime),
		RAMSize: *RAMSize,
		Files:   make([]string, 0, len(flag.Args()))}
	if *RAMSize <= 0 {
		job.Log.Fatal("[ERROR] ", RAMSizeUsage, "It was set to: ", *RAMSize)
	}
	// set outfile if we pass validity tests
	if *outFile == "" {
		job.Log.Fatal("[ERROR] ", outFileUsage)
	}
	if _, err := os.Stat(*outFile); err == nil {
		job.Log.Fatal("[ERROR] Output file (", *outFile, ") exists.",
			" Either delete it or provide a new file name.")
	}
	job.Outfile = *outFile
	// set input files for the ones that pass validity tests
	for _, f := range flag.Args() {
		if _, err := os.Stat(f); err != nil {
			job.Log.Print("[WARN] File ", f, " does not exist. Will skip it.")
		} else {
			job.Files = append(job.Files, f)
		}
	}
	if len(job.Files) == 0 {
		job.Log.Fatal("[ERROR] No input files were given.")
	}
	return job
}
