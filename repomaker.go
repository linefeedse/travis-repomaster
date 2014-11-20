package main

import "encoding/json"
import "fmt"
import "io/ioutil"
import "strings"
import "net/url"
import "os"
import "github.com/jmcvetta/napping"

func check(e error) {
    if e != nil {
        panic(e)
    }
}

// Create a repo from json data (previously read from a file)

func createrepo(r []byte) {

    username := os.Getenv("GITHUBTOKEN")
    passwd   := "x-oauth-basic"

    res := struct {
        Id          int
        Name        string
        Fullname    string `json:"full_name"`
        Description string
        Private     bool
        Fork        bool
        Url         string
        Htmlurl     string `json:"html_url"`
        UpdatedAt   string `json:"updated_at"`
        CreatedAt   string `json:"created_at"`
        Owner       struct {
            Login string
            Id    int
        }
    }{}
    e := struct {
        Message string
        Errors  []struct {
            Resource string
            Field    string
            Code     string
        }
    }{}
    var payload map[string]interface{}
    err := json.Unmarshal(r, &payload)
    check(err)
    s := napping.Session{
        Userinfo: url.UserPassword(username, passwd),
    }
    url := "https://api.github.com/user/repos"
    resp, err := s.Post(url, &payload, &res, &e)
    check(err)
    if resp.Status() == 201 {
        fmt.Printf("Repo created at: %s\n\n", res.CreatedAt)
    } else if resp.Status() == 422 {
        fmt.Printf("Repo already exists:\n\n")
    } else {
        fmt.Println("Bad response status from Github server")
        fmt.Printf("\t Status: %v\n", resp.Status())
        fmt.Printf("\t Message: %v\n", e.Message)
        fmt.Printf("\t Errors: %v\n", e.Message)
    }
}

func main() {
    rdirs, err := ioutil.ReadDir(".")
    check(err)
    jrdirs := make([]string,0)
    for n := range rdirs {
        if rdirs[n].IsDir() {
            name := rdirs[n].Name()
            jrdirs = append(jrdirs,rdirs[n].Name())
            jsonname := strings.Join([]string{name, "/create.json"},"")
            rinfo , err := ioutil.ReadFile(jsonname)
            if err == nil {
                createrepo(rinfo)
            }
        }
    }
    jsonDirs, _ := json.Marshal(jrdirs)
    fmt.Println(string(jsonDirs))
}

// vim:et:sts=4:sw=4
