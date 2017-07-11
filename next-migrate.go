package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {

	var (
		real bool
		git  func(args ...string) error
	)

	flag.BoolVar(&real, "real", false, "do the real stuff")

	flag.Parse()

	if real {
		git = gitExec

	} else {
		git = func(args ...string) error {
			fmt.Println("git", strings.Join(args, " "))
			return nil
		}
	}

	// This is work in progress.
	fmt.Println("Migrate hugoDocs to the new Hugo docs concept.")

	toMove, err := readMoveStatements()
	if err != nil {
		log.Fatal(err)
	}

	for from, to := range toMove {

		if from == to {
			continue
		}

		if _, err := os.Stat(to); err == nil {
			fmt.Println("Skip existing", to)
			continue
		}

		// Create dir if not exist
		os.MkdirAll(filepath.Dir(to), 0755)
		// We could probably just as well do a OS move, but let us have git do it.
		if err := git("mv", from, to); err != nil {
			log.Fatal(err)
		}
	}

}

func gitExec(args ...string) error {
	cmd := exec.Command("git", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git failed: %q: %q (%q)", err, out, args)
	}
	return nil
}

func readMoveStatements() (map[string]string, error) {
	f, err := os.Open("next-move.csv")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	fromTo := make(map[string]string)

	lines, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return nil, err
	}

	for i, line := range lines {
		if i == 0 {
			// header
			continue
		}

		fromTo[prepPath(line[0])] = prepPath(line[1])

	}

	return fromTo, nil
}

func prepPath(s string) string {
	return strings.TrimPrefix(s, "/")
}
