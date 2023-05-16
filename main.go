package main

import (
    "fmt"
    "os"
    "log"
    // "io"
    // "strconv"
    // "strings"
)

func getVersion() {
    version := "0.0.1"
    fmt.Println(version)
}

func help() {
    fmt.Println("Usage: kv [options]")
    fmt.Println()
    fmt.Println("Common uses:")
    fmt.Println("kv init        Create empty kv repository in the current directory")
}

func createInitFiles() {
    // Create folders
    if err := os.MkdirAll(".kv/final/", os.ModePerm); err != nil {
        log.Fatal(err)
    }
    if err := os.MkdirAll(".kv/commit/", os.ModePerm); err != nil {
        log.Fatal(err)
    }
    if err := os.MkdirAll(".kv/commit/", os.ModePerm); err != nil {
        log.Fatal(err)
    }

    // Create files
    f1, err := os.Create(".kv/staging-area.txt")
    if err != nil {
        log.Fatal(err)
    }
    defer f1.Close()

    f2, err := os.Create(".kv/status.txt")
    if err != nil {
        log.Fatal(err)
    }
    defer f2.Close()

    currentDir, err := os.Getwd()
    if err != nil {
    	log.Println(err)
    }
    fmt.Printf("Initialized empty kv repository in %s\n", currentDir)
}

func checkInitFiles() {
    if _, err := os.Stat(".kv/"); !os.IsNotExist(err) {
    	// .kv/ directory does exist
        fmt.Println("Directory already initialized!")
        os.Exit(0)
    }
}

func kvInit() {
    checkInitFiles()
    createInitFiles()
}

func main() {
    if (len(os.Args) == 1) {
        help()
        os.Exit(0)
    }

    for i := 1; i < len(os.Args); i++ {
        switch os.Args[i] {
            case "-h", "--help":
                help()

            case "-v", "--version":
                getVersion()

            case "init":
                kvInit()

            default:
                fmt.Printf("Unknown argument \"%s\"\n", os.Args[i])
        }
    }
}
