package main

import (
	"flag"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

var compareBranch string
var ignorecsv string
var ignore []string = []string{"master", "staging", "qa", "release", "uat", "production", "HEAD"}

func init() {
    ignoredDefault := ""
    for _, i := range ignore {
        ignoredDefault += i + " "
    }
    flag.StringVar(&compareBranch, "b", "", "Please provide the branch to compare (default develop)")
    flag.StringVar(&ignorecsv, "i", "", "Branches to ignore, csv separated. Default:" + ignoredDefault)
}

func main() {
	flag.Parse()

	ignoreMap := make(map[string]string)
	for _, i := range ignore {
		ignoreMap[i] = i
	}

    for _, i := range strings.Split(ignorecsv, ","){
        ignoreMap[i] = i
    }

	if compareBranch == "" {
		compareBranch = "develop"
	}

	cmd := exec.Command("git", "branch", "-a", "--merged", compareBranch)
	regex := regexp.MustCompile("(\\*)? +(remotes/)?(.*?)(\n| .*\n)")

	if b, err := cmd.Output(); err == nil {
		output := string(b)
		matches := regex.FindAllStringSubmatch(output, -1)
		for _, m := range matches {
			isCurrent := m[1] == "*"
			isRemote := m[2] == "remotes/"

			sp := strings.SplitN(m[3], "/", 2)

			remote := "origin"
			branch := ""

			if len(sp) == 2 {
				remote, branch = sp[0], sp[1]
			}

			if _, ignoreBranch := ignoreMap[branch]; len(branch) == 0 || ignoreBranch || isCurrent || branch == compareBranch {
				continue
			}

			delLocal := exec.Command("git", "branch", "-d", branch)
			delRemote := exec.Command("git", "push", "--delete", remote, branch)

			if isRemote {
				if remoteOutput, remerr := delRemote.Output(); remerr != nil {
					fmt.Println("command failed", delRemote, string(remoteOutput), remerr)
				}
			} else {
				if localOutput, locerr := delLocal.Output(); locerr != nil {
					fmt.Println("command failed", delLocal, string(localOutput), locerr)
				}
			}
		}
	} else {
		fmt.Println(err)
	}
}
