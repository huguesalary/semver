package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/Masterminds/semver"
)

const (
	toolVersion string = "0.0.1-rc1"
)

func main() {
	major := 0
	minor := 0
	patch := 0
	prerelease := ""
	metadata := ""
	thisversion := false

	flag.IntVar(&major, "major", 0, "increase major by n")
	flag.IntVar(&minor, "minor", 0, "increase minor by n")
	flag.IntVar(&patch, "patch", 0, "increase patch by n")
	flag.StringVar(&prerelease, "prerelease", "", "set prerelease")
	flag.StringVar(&metadata, "metadata", "", "set metadata")
	flag.BoolVar(&thisversion, "version", false, "output this tool's version")
	flag.Parse()
	v := flag.Arg(0)  // Get first version
	v1 := flag.Arg(1) // Get second version

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), `
Usage: %s [flags] version [version]
	If [version] is provided, the command will output the whether version is lower (-1),
	equal (0), or greater (1) than [version] after the applying the major, minor, patch,
	prerelease and metadata transformations to version

`, os.Args[0])

		flag.PrintDefaults()
	}

	if thisversion {
		fmt.Print(toolVersion)
		os.Exit(0)
	}

	if v == "" {
		flag.Usage()
		os.Exit(1)
	}

	// Parse the version given as argument
	// Output and quit on error
	version, err := semver.NewVersion(v)
	if err != nil {
		log.Fatalf("error parsing semver for %s: %s", v, err)
	}

	// Parse the second version given as argument
	// Output and ignore on error
	var version1 *semver.Version
	if v1 != "" {
		version1, err = semver.NewVersion(v1)
		if err != nil {
			log.Printf("error parsing semver for secondary version %s: %s", v1, err)
		}
	}

	// If any of the najor, minor or patch are invalid, output the issue, but ignore it
	if major < 0 {
		log.Println("invalid major increment; only positive numbers allowed")
	}

	if minor < 0 {
		log.Println("invalid minor increment; only positive numbers allowed")
	}

	if patch < 0 {
		log.Println("invalid patch increment; only positive numbers allowed")
	}

	// Increase major, minor and patch
	for ; major > 0; major-- {
		*version = version.IncMajor()
	}

	for ; minor > 0; minor-- {
		*version = version.IncMinor()
	}

	for ; patch > 0; patch-- {
		*version = version.IncPatch()
	}

	// Parse prerelease and ignore errors
	if prerelease != "" {
		*version, err = version.SetPrerelease(prerelease)
		if err != nil && err == semver.ErrInvalidPrerelease {
			log.Printf("invalid prerelease format; your prerelease must follow this regex: %s", semver.ValidPrerelease)
		}
	}

	// Parse metadata and ignore errors
	if metadata != "" {
		*version, err = version.SetMetadata(metadata)
		if err != nil && err == semver.ErrInvalidMetadata {
			log.Printf("invalid metadata format; your metadata must follow this regex: %s", semver.ValidPrerelease)
		}
	}

	// If a second version was passed in
	// Add the result of the comparison to the output
	if version1 != nil {
		fmt.Print(version.Compare(version1))
	} else {
		// Construct the output
		fmt.Print(version)
	}
}
