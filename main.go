package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

var ignore []string = []string{"master", "develop", "staging", "qa", "release", "uat", "production", "stage", "prod", "HEAD"}

func main() {
	compareBranch := flag.String("b", "develop", "branch to compare (default develop)")
	ignorecsv := flag.String("i", "", "branches to ignore, default: "+strings.Join(ignore, ", "))
	dryrun := flag.Bool("dry", false, "only output commands to run")
	//nopush := flag.Bool("nopush", false, "don't push changes to remote")
	flag.Parse()

	ignoreMap := make(map[string]string)
	for _, i := range ignore {
		ignoreMap[i] = i
	}

	for _, i := range strings.Split(*ignorecsv, ",") {
		ignoreMap[i] = i
	}

	cmd := exec.Command("git", "branch", "-a", "--merged", *compareBranch)
	regex := regexp.MustCompile("(\\*)? +(remotes/)?(.*?)(\n| .*\n)")

	b, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	output := string(b)
	matches := regex.FindAllStringSubmatch(output, -1)
	for _, m := range matches {
		isCurrent := m[1] == "*"
		isRemote := m[2] == "remotes/"

		sp := strings.SplitN(m[3], "/", 2)

		remote := "origin"
		branch := ""

		if isRemote && len(sp) == 2 {
			remote, branch = sp[0], sp[1]
		} else if !isRemote {
			remote, branch = "origin", m[3]
		}

		if _, ignoreBranch := ignoreMap[branch]; len(branch) == 0 || ignoreBranch || isCurrent || branch == *compareBranch {
			continue
		}

		delLocal := exec.Command("git", "branch", "-d", branch)
		delRemote := exec.Command("git", "push", "--delete", remote, branch)
		remotes := []string{}

		shouldPrune := false

		if *dryrun {
			fmt.Println(remote, branch)
		} else {
			if isRemote {
				fmt.Println("git", "push", "--delete", remote, branch)
				if remoteOutput, remerr := delRemote.CombinedOutput(); remerr != nil {
					fmt.Println("command failed, removing remote", string(remoteOutput), remerr)
				}

				shouldPrune = true
				remotes = append(remotes, remote)
			} else {
				fmt.Println("git", "branch", "-d", branch)
				if localOutput, locerr := delLocal.CombinedOutput(); locerr != nil {
					fmt.Println("command failed, removing local", string(localOutput), locerr)
				}
			}
		}

		if shouldPrune {
			for _, rem := range remotes {
				prune := exec.Command("git", "remote", "prune", rem)
				fmt.Println("git", "remote", "prune", rem)
				if pruneOutput, pruerr := prune.CombinedOutput(); pruerr != nil {
					fmt.Println("command failed, pruning remote", rem, string(pruneOutput), pruerr)
				}
			}
		}
	}
}
