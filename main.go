package main

import (
    "fmt"
    "os"
    "log"
    "time"
    "bufio"
    "strings"
    "crypto/sha1"
    "io"
    "encoding/hex"
    "io/ioutil"
    "strconv"
    "path/filepath"
    // "compress/zlib"
)

func getVersion() {
    version := "0.0.1"
    fmt.Println(version)
}

func help() {
    fmt.Println("Usage: kv [options]")
    fmt.Println()
    fmt.Println("Common uses:")
    fmt.Println("kv init                Create empty kv repository in the current directory")
    fmt.Println("kv add <filename>      Add <filename> to staging area")
    fmt.Println("kv status              Show what's inside the staging area")
    fmt.Println("kv commit              Commit changes to a new version")
}

func contentToBytes(filepath string) []byte {
   file, err := os.Open(filepath)
   if err != nil {
		fmt.Println(err)
      return nil
	}
   defer file.Close()

   // Get the file size
   stat, err := file.Stat()
   if err != nil {
      fmt.Println(err)
      return nil
   }

   // Read the file into a byte slice
   bs := make([]byte, stat.Size())
   _, err = bufio.NewReader(file).Read(bs)
   if err != nil && err != io.EOF {
      fmt.Println(err)
      return nil
   }

   return bs
}

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

func getHashSum(filepath string) string {
    sha1data := []byte(contentToBytes(filepath))
    sha1sum := sha1.Sum(sha1data)
    sha1Write := hex.EncodeToString(sha1sum[:])
    return sha1Write
}

func stageFile(fileToStage string) {
    f, err := os.OpenFile(".kv/staging-area.txt",
    	os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
    	log.Println(err)
    }
    defer f.Close()

    isDuplicate, lineNum := duplicateStageFile(fileToStage)
    if (isDuplicate) {
        // fmt.Printf("Found duplicate on line: %d\n", lineNum)
        deleteLine(".kv/staging-area.txt", lineNum)
    }

    if fileExists(fileToStage) {
        sha1Write := getHashSum(fileToStage)
        writeString := fileToStage + ";" + getCurrentTime() + ";" + sha1Write + ";created\n"
        if _, err := f.WriteString(writeString); err != nil {
        	log.Println(err)
        }

        fmt.Printf("Added %s to the repository.\n", fileToStage)
    } else {
        fmt.Printf("%s does not exist.\n", fileToStage)
    }
}

func deleteLine(file string, num int) {
    // Remove content from file if it exists
    f, _ := os.Open(file)

    // create and open a temporary file
    f_tmp, err := os.CreateTemp("", "tmpfile-*.txt")
    if err != nil {
        log.Fatal(err)
    }

    // Copy content from original to tmp
    _, err = io.Copy(f_tmp, f)
    if err != nil {
        log.Fatal(err)
    }

    f, _ = os.Create(file)
    tmpfile, _ := os.Open(f_tmp.Name())

    // Create new Scanner.
    scanner := bufio.NewScanner(tmpfile)
    line_num := 0
    // Use Scan.
    for scanner.Scan() {
        line_num++
        line := scanner.Text()
        if line_num != num {
            if _, err := f.Write([]byte(line + "\n")); err != nil {
                fmt.Println(err)
            }
        }
    }

    defer os.Remove(f_tmp.Name())
    defer f.Close()
}

func duplicateStageFile(filename string) (bool, int) {
    // check for duplicate stage files
    isDuplicate := false

    readFile, err := os.Open(".kv/staging-area.txt")
    if err != nil {
        log.Println(err)
    }

    fileScanner := bufio.NewScanner(readFile)

    fileScanner.Split(bufio.ScanLines)

    lineNum := 0
    for fileScanner.Scan() {
        lineNum = lineNum + 1
        line := fileScanner.Text()
        splitLine := strings.Split(line, ";")
        // splitLine[0] - filepath
        // splitLine[1] - modification date
        // splitLine[2] - sha1 hash
        // splitLine[3] - status (created/updated/deleted)

        if (splitLine[0] == filename) {
            isDuplicate = true
            // fmt.Printf("Found duplicate on line: %d", lineNum)
            return isDuplicate, lineNum
        }
    }
    return isDuplicate, lineNum
}

func stagingAreaLocation() string {
    var rootDir string
    if (len(getRootDir()) == 1) {
        rootDir = ".kv/"
    } else {
        rootDir = getRootDir() + "/.kv/"
    }

    return rootDir + "staging-area.txt"
}

func getStagingArea() {
    readFile, err := os.Open(stagingAreaLocation())
    if err != nil {
        log.Println(err)
    }

    fileScanner := bufio.NewScanner(readFile)
    fileScanner.Split(bufio.ScanLines)

    for fileScanner.Scan() {
        // lineNum = lineNum + 1
        line := fileScanner.Text()
        splitLine := strings.Split(line, ";")
        // splitLine[0] - filepath
        // splitLine[1] - modification date
        // splitLine[2] - sha1 hash
        // splitLine[3] - status (created/updated/deleted)

        fmt.Printf("%s: %s - %s\n", strings.ToUpper(splitLine[3]), splitLine[0], splitLine[1])
    }
}

func getRootDir() string {
    start := "."
    count := 0
    for count < 5 {
        check := filepath.Join(start, ".kv")
        fi, err := os.Stat(check)
        if err == nil && fi.IsDir() {
            return start
        }
        start = filepath.Join(start, "..")
        count++
    }
    return ""
}

// func dirHasUntrackedFiles() bool, string[] {
//     // Oh god this function is so ugly
//     // fix this with proper smaller functions
//
//     readFile, err := os.Open(".kv/staging-area.txt")
//     if err != nil {
//         log.Println(err)
//     }
//
//     fileScanner := bufio.NewScanner(readFile)
//     fileScanner.Split(bufio.ScanLines)
//
//     stagedFiles := string[] {}
//     for fileScanner.Scan() {
//         line := fileScanner.Text()
//         splitLine := strings.Split(line, ";")
//         // splitLine[0] - filepath
//         // splitLine[1] - modification date
//         // splitLine[2] - sha1 hash
//         // splitLine[3] - status (created/updated/deleted)
//         stagedFiles = append(stagedFiles, splitLine[0])
//     }   // This gets paths of all files in staging
//
//     newestCommitDir := ".kv/commit/v" + strconv.Itoa(commitNumber()) + "/"
//
//     readFile2, err := os.Open(newestCommitDir)
//     if err != nil {
//         log.Println(err)
//     }
//
//     fileScanner := bufio.NewScanner(readFile2)
//     fileScanner.Split(bufio.ScanLines)
//
//     commitedFiles := string[] {}
//     for fileScanner.Scan() {
//         line := fileScanner.Text()
//         splitLine := strings.Split(line, ";")
//         // splitLine[0] - filepath
//         // splitLine[1] - modification date
//         // splitLine[2] - sha1 hash
//         // splitLine[3] - status (created/updated/deleted)
//         commitedFiles = append(commitedFiles, splitLine[0])
//     }   // This gets paths of all files in the repo
//
//
//     repoFiles := string[] {}
//
//     err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
//         if err != nil {
//             log.Println(err)
//         }
//         // fmt.Printf("dir: %v: name: %s\n", info.IsDir(), path)
//         if (!info.IsDir()) {
//             repoFiles = append(repoFiles, path)
//         }
//         // return nil
//     })
//     if err != nil {
//         log.Println(err)
//     }
//
//     untrackedFiles := string[] {}
//     for _, a := range repoFiles {
//         exists := false
//         for _, b := range commitedFiles {
//             if a == b {
//                 exists = true
//                 break
//             }
//         }
//         if !exists {
//             untrackedFiles = append(untrackedFiles, a)
//         }
//     }
//
//     if (len(untrackedFiles) == 0) {
//         return false, untrackedFiles
//     }
//
//     return true, untrackedFiles
// }

func kvStatus() {
    if (!isStagingEmpty()) {
        fmt.Println("Staging:")
        fmt.Println("============================")
        getStagingArea()
        fmt.Println()
    }

    // if (dirHasUntrackedFiles()) {
    //     fmt.Println("Untracked files:")
    //     fmt.Println("============================")
    //     showUntrackedFiles()
    // }
}

func clearStagingArea() {
    if err := os.Truncate(stagingAreaLocation(), 0); err != nil {
        log.Printf("Failed to truncate: %v", err)
    }
}

func readLines(path string) ([]string, error) {
    file, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    var lines []string
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        lines = append(lines, scanner.Text())
    }
    return lines, scanner.Err()
}

func copyFile(src string, dst string) {
    // Read all content of src to data, may cause OOM for a large file.
    data, err := ioutil.ReadFile(src)
    if err != nil {
        log.Fatal(err)
    }

    // Write data to dst
    err = ioutil.WriteFile(dst, data, 0644)
    if err != nil {
        log.Fatal(err)
    }
}

func trimLeftChar(s string) string {
	for i := range s {
		if i > 0 {
			return s[i:]
		}
	}
	return s[:0]
}

func commitNumber() int {
    // returns the last commit number
    files, err := os.ReadDir(".kv/commit/")
    if err != nil {
        log.Fatal(err)
    }

    dirArray := []string {}
    for _, file := range files {
        dirArray = append(dirArray, file.Name())
    }

    return len(dirArray)
}

func isStagingEmpty() bool {
    var rootDir string
    if (len(getRootDir()) == 1) {
        rootDir = ".kv/"
    } else {
        rootDir = getRootDir() + "/.kv/"
    }
    stagingAreaLocation := rootDir + "staging-area.txt"



    stagingArea, err := os.Stat(stagingAreaLocation)
    if err != nil {
    	log.Println(err)
    }
    stagingSize := stagingArea.Size()

    if (stagingSize == 0) {
        return true
    }
    return false
}

func commitFiles() {
    // If first commit, make commitNum = 1 instead of 0
    commitNum := commitNumber() + 1

    commits, err := readLines(".kv/staging-area.txt")
    if err != nil {
        log.Println(err)
    }

    for i := 0; i < len(commits); i++ {
        singleCommit := strings.Split(commits[i], ";")
        // singleCommit[0] - filepath
        // singleCommit[1] - modification date
        // singleCommit[2] - sha1 hash
        // singleCommit[3] - status (created/updated/deleted)
        commitVersion := ".kv/commit/v" + strconv.Itoa(commitNum) + "/"
        commitedFile := commitVersion + singleCommit[0]

        os.MkdirAll(commitVersion, 0700) // Create your file

        f, err := os.Create(commitedFile)
        if err != nil {
            log.Fatal(err)
        }
        defer f.Close()

        copyFile(singleCommit[0], commitedFile)
    }
    clearStagingArea()
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

            case "status":
                kvStatus()
                os.Exit(0)

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

            case "commit":
                if (isStagingEmpty()) {
                    fmt.Println("Nothing to commit.")
                    os.Exit(0)
                }

                commitFiles()
                os.Exit(0)

            default:
                fmt.Printf("Unknown argument \"%s\"\n", os.Args[i])
        }
    }
}
