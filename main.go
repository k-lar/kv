package main

import (
    "fmt"
    "os"
    "log"
    "time"
	// "crypto/sha1"
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

// func contentToBytes(string filepath) []byte {
//    file, err := os.Open(filepath)
//    if err != nil {
// 		fmt.Println(err)
//       return nil
// 	}
//    defer file.Close()
//
//    // Get the file size
//    stat, err := file.Stat()
//    if err != nil {
//       fmt.Println(err)
//       return nil
//    }
//
//    // Read the file into a byte slice
//    bs := make([]byte, stat.Size())
//    _, err = bufio.NewReader(file).Read(bs)
//    if err != nil && err != io.EOF {
//       fmt.Println(err)
//       return nil
//    }
//
//    return bs
// }

func getCurrentTime() string {
    currentTime := time.Now()
    timeString := currentTime.Format("2006-01-02 15:04:05")
    return timeString
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
    fmt.Printf("Initialized empty kv repository in %s/.kv/\n", currentDir)
}

func checkInitFiles() bool {
    result := false
    if _, err := os.Stat(".kv/"); !os.IsNotExist(err) {
    	// .kv/ directory does exist
        result = true
    }
    return result
}

func kvInit() {
    if (checkInitFiles() == true) {
        fmt.Println("Directory already initialized!")
        os.Exit(0)
    }
    createInitFiles()
}

func stageFile(file_to_stage string) {
    f, err := os.OpenFile(".kv/staging-area.txt",
    	os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
    	log.Println(err)
    }
    defer f.Close()

    if fileExists(file_to_stage) {
        writeString := file_to_stage + ";" + getCurrentTime() + ";created\n"
        if _, err := f.WriteString(writeString); err != nil {
        	log.Println(err)
        }
        fmt.Printf("Added %s to the repository.\n", file_to_stage)
    } else {
        fmt.Printf("%s does not exist.", file_to_stage)
    }
}

func duplicateStageFile() {
    // check for duplicate stage files
}

func fileExists(filepath string) bool {
    result := false
    if _, err := os.Stat(filepath); err == nil {
       result = true

    }
    return result
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

            case "add":
                if (len(os.Args) > i+1) {
                    for i := i+1; i < len(os.Args); i++ {
                        // fmt.Printf("Added %s\n", os.Args[i])
                        stageFile(os.Args[i])

                    }

                } else {
                    fmt.Println("Nothing to add.")
                    os.Exit(0)
                }
                os.Exit(0)

            case "time":
                fmt.Println(getCurrentTime())

            default:
                fmt.Printf("Unknown argument \"%s\"\n", os.Args[i])
        }
    }
}
