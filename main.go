package main

import (
    "fmt"
    "flag"
    "os/exec"
    "strings"
)

var compareBranch string

func init(){
    flag.StringVar(&compareBranch, "b", "", "Please provide the branch to compare (default develop)")
}


func main(){
    flag.Parse()

    if compareBranch == "" {
        compareBranch = "develop"
    }
    cmd := exec.Command("git", "branch", "--merged", compareBranch)

    if b, err := cmd.Output(); err == nil {
        output := string(b)
        branches := strings.Split(output, "\n")
        for _, branch := range branches {
            branch = strings.Trim(branch, " ")
            if !strings.HasPrefix(branch, "* ") {
                delLocal := exec.Command("git", "branch", "-d", branch)
                delRemote := exec.Command("git", "push", "--delete", "origin", branch)
            }
        }
    } else {
        fmt.Println(err)
    }
}