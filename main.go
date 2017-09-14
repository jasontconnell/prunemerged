package main

import (
    "fmt"
    "flag"
    "os/exec"
)

var branch string

func init(){
    flag.StringVar(&branch, "b", "", "Please provide the branch to compare (default develop)")
}


func main(){
    flag.Parse()

    if branch == "" {
        //fmt.Println("Usage ::  prunemerged -b [branch-name]")
        branch = "develop"
    }
    cmd := exec.Command("git", "branch", "--merged", branch)

    if b, err := cmd.Output(); err == nil {
        fmt.Println("got read", string(b))
    } else {
        fmt.Println(err)
    }
}