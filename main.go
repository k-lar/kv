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
    currentDir, err := os.Getwd()
    if err != nil {
    	log.Println(err)
    }
    currentDirAndFile := currentDir + "/" + fileToStage

    os.Chdir(getRootDir())
    rootDir, err := os.Getwd()
    if err != nil {
    	log.Println(err)
    }

    relPathToFile, err := filepath.Rel(rootDir, currentDirAndFile) //use the Rel function to get the relative path
    if err != nil {
       log.Println("Error:", err) //print error if no path is obtained
    }

    f, err := os.OpenFile(stagingAreaLocation(),
    	os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
    	log.Println(err)
    }
    defer f.Close()

    // TODO: Show if file is modified. (Don't always show created)
    isDuplicate, lineNum := duplicateStageFile(relPathToFile)
    if (isDuplicate) {
        // fmt.Printf("Found duplicate on line: %d\n", lineNum)
        deleteLine(stagingAreaLocation(), lineNum)
    }

    if fileExists(relPathToFile) {
        sha1Write := getHashSum(relPathToFile)
        writeString := relPathToFile + ";" + getCurrentTime() + ";" + sha1Write + ";created\n"
        if _, err := f.WriteString(writeString); err != nil {
        	log.Println(err)
        }

        fmt.Printf("Added %s to the repository.\n", relPathToFile)
    } else {
        fmt.Printf("%s does not exist.\n", relPathToFile)
    }

    os.Chdir(currentDir)
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

func getStagedFiles() [][]string {
    readFile, err := os.Open(stagingAreaLocation())
    if err != nil {
        log.Println(err)
    }

    fileScanner := bufio.NewScanner(readFile)
    fileScanner.Split(bufio.ScanLines)

    var result [][]string

    for fileScanner.Scan() {
        line := fileScanner.Text()
        splitLine := strings.Split(line, ";")
        // splitLine[0] - filepath
        // splitLine[1] - modification date
        // splitLine[2] - sha1 hash
        // splitLine[3] - status (created/updated/deleted)

        row := []string{splitLine[0], splitLine[2]}
        result = append(result, row)
    }
    return result
}

func getCommitedFiles() [][]string {
    var commitedFiles [][]string
    var commitDir string
    if (len(getRootDir()) == 1) {
        commitDir = ".kv/commit/v"
    } else {
        commitDir = getRootDir() + "/.kv/commit/v"
    }
    commitDir = commitDir + strconv.Itoa(commitNumber()) + "/"

    err := filepath.Walk(commitDir, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            log.Println(err)
        }
        if (!info.IsDir()) {
            hashOfFile := getHashSum(path)
            row := []string{path, hashOfFile}
            commitedFiles = append(commitedFiles, row)
        }
        return nil
    })
    if err != nil {
        log.Println(err)
    }

    return commitedFiles
}

func getAllFiles() [][]string {
    var filesInRootDir [][]string

    err := filepath.Walk(getRootDir(), func(path string, info os.FileInfo, err error) error {
        if err != nil {
            log.Println(err)
        }
        // This is ugly as hell. || operator doesnt work for some reason here so idk what else to do. Wtf
        if (!info.IsDir()) {
            if (!strings.HasPrefix(path, ".kv")) {
                // TODO: Do I keep this? I shouldn't, but it's convenient sometimes.
                // This ignores everything git related so it doesn't clutter up my testing.
                // Make this check for environment variable 'KV_GIT_IGNORE == true'
                if (!strings.HasPrefix(path, ".git")) {
                    hashOfFile := getHashSum(path)
                    row := []string{path, hashOfFile}
                    filesInRootDir = append(filesInRootDir, row)
                }
            }
        }
        return nil
    })
    if err != nil {
        log.Println(err)
    }

    return filesInRootDir
}

func getCommitedFilesShort() [][]string {
    commitedFiles := getCommitedFiles()
    shortCommitedFiles := [][]string {}
    for _, file := range commitedFiles {
        split := strings.Split(file[0], "/")
        split = split[3:]
        shortFilePath := ""
        for i := 0; i < len(split)-1; i++ {
            shortFilePath = shortFilePath + split[i] + "/"
        }
        shortFilePath = shortFilePath + split[len(split)-1]
        row := []string{shortFilePath, file[1]}
        shortCommitedFiles = append(shortCommitedFiles, row)
    }
    return shortCommitedFiles
}

func trackFiles() []string {
    // REALLY BAD CODE.

    // Go to rootDir first to be accurate
    oldDir, err := os.Getwd()
    if err != nil {
    	log.Println(err)
    }

    os.Chdir(getRootDir())

    allFiles := getAllFiles()
    stagedFiles := getStagedFiles()
    commitedFiles := getCommitedFilesShort()

    untrackedFiles := []string{}
    addedFiles := []string{}

    for _, fileInRepo := range allFiles {
        foundInStaged := false
        for _, fileStaged := range stagedFiles {
            if fileInRepo[0] == fileStaged[0] { // if filenames match
                if fileInRepo[1] != fileStaged[1] { // if hashes don't match
                    if !contains(addedFiles, fileInRepo[0]) {
                        // fileAndStatus := "updated;" + fileInRepo[0]
                        untrackedFiles = append(untrackedFiles, fileInRepo[0])
                        addedFiles = append(addedFiles, fileInRepo[0])
                    }
                    foundInStaged = true
                    break
                }
                foundInStaged = true
                if !contains(addedFiles, fileInRepo[0]) { // if hashes match and file not already added
                    addedFiles = append(addedFiles, fileInRepo[0])
                }
                break
            }
        }
        if foundInStaged {
            continue
        }

        foundInCommited := false
        for _, fileCommited := range commitedFiles {
            if fileInRepo[0] == fileCommited[0] { // if filenames match
                if fileInRepo[1] != fileCommited[1] { // if hashes don't match
                    if !contains(addedFiles, fileInRepo[0]) {
                        // fileAndStatus := "modified;" + fileInRepo[0]
                        untrackedFiles = append(untrackedFiles, fileInRepo[0])
                        addedFiles = append(addedFiles, fileInRepo[0])
                    }
                    break
                }
                foundInCommited = true
                break
            }
        }
        if foundInCommited {
            continue
        }
        if !foundInStaged {
            if !contains(addedFiles, fileInRepo[0]) {
                // fileAndStatus := "untracked;" + fileInRepo[0]
                untrackedFiles = append(untrackedFiles, fileInRepo[0])
                addedFiles = append(addedFiles, fileInRepo[0])
            }
        }
    }
    os.Chdir(oldDir)
    return untrackedFiles
}

func contains(arr []string, item string) bool {
    for _, arrItem := range arr {
        if arrItem == item {
            return true
        }
    }
    return false
}

func removeDuplicateStr(strSlice []string) []string {
    allKeys := make(map[string]bool)
    list := []string{}
    for _, item := range strSlice {
        if _, value := allKeys[item]; !value {
            allKeys[item] = true
            list = append(list, item)
        }
    }
    return list
}

func getDeletedFiles() []string {
    commitedFilesWithHashes := getCommitedFilesShort()
    commitedFiles := []string {}

    stagedFilesWithHashes := getStagedFiles()
    stagedFiles := []string {}

    allFilesWithHashes := getAllFiles()
    allFiles := []string {}

    for i := 0; i < len(commitedFilesWithHashes); i++ {
        commitedFiles = append(commitedFiles, commitedFilesWithHashes[i][0])
    }
    for i := 0; i < len(stagedFilesWithHashes); i++ {
        stagedFiles = append(stagedFiles, stagedFilesWithHashes[i][0])
    }
    for i := 0; i < len(allFilesWithHashes); i++ {
        allFiles = append(allFiles, allFilesWithHashes[i][0])
    }


    deletedFiles := []string {}

    for i := 0; i < len(commitedFiles); i++ {
        if (!contains(allFiles, commitedFiles[i])) {
            deletedFiles = append(deletedFiles, commitedFiles[i])
        }
    }
    for i := 0; i < len(stagedFiles); i++ {
        if (!contains(allFiles, stagedFiles[i])) {
            deletedFiles = append(deletedFiles, stagedFiles[i])
        }
    }

    return removeDuplicateStr(deletedFiles)
}

// TODO: Optimize kvStatus(), jesus christ is it slow.
// Pass variables to functions instead of always calculating everything again
func kvStatus() {
    if (!isStagingEmpty()) {
        fmt.Println("Staging:")
        fmt.Println("============================")
        getStagingArea()
        fmt.Println()
    }

    if (commitNumber() != 0) {
        untrackedFiles := trackFiles()
        if (len(untrackedFiles) != 0) {
            fmt.Println("Untracked files:")
            fmt.Println("============================")
            for i := 0; i < len(untrackedFiles); i++ {
                // splitUntrackedFiles := strings.Split(untrackedFiles[i], ";")

                // fmt.Printf("%s: %s\n", strings.ToUpper(splitUntrackedFiles[0]), splitUntrackedFiles[1])
                fmt.Printf("Untracked changes: %s\n", untrackedFiles[i])
            }
            fmt.Println()
        }
    }

    deletedFiles := getDeletedFiles()
    if (len(deletedFiles) > 0) {
        fmt.Println("Deleted files:")
        fmt.Println("============================")
        for i := 0; i < len(deletedFiles); i++ {
            // splitUntrackedFiles := strings.Split(untrackedFiles[i], ";")

            // fmt.Printf("%s: %s\n", strings.ToUpper(splitUntrackedFiles[0]), splitUntrackedFiles[1])
            fmt.Printf("Removed: %s\n", deletedFiles[i])
        }
    }
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
    var rootDir string
    if (len(getRootDir()) == 1) {
        rootDir = ".kv/"
    } else {
        rootDir = getRootDir() + "/.kv/"
    }
    commitDir := rootDir + "/commit/"

    files, err := os.ReadDir(commitDir)
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
    // FIXME: Can not commit files if you're inside the .kv/ directory
    stagingArea, err := os.Stat(stagingAreaLocation())
    if err != nil {
    	log.Println(err)
    }

    stagingSize := stagingArea.Size()
    if (stagingSize == 0) {
        return true
    }
    return false
}

func copyFromPrevCommit() {
    nextCommitNum := commitNumber() + 1
    nextCommitVersion := ".kv/commit/v" + strconv.Itoa(nextCommitNum) + "/"
    prevCommitFiles := getCommitedFiles()

    // If file not in root dir, create directories
    for i := 0; i < len(prevCommitFiles); i++ {
        pathToFile := ""
        commitedFile := prevCommitFiles[i][0]
        if (strings.Contains(prevCommitFiles[i][0], "/")) {
            dirToFile := strings.Split(prevCommitFiles[i][0], "/")
            dirToFile = append(dirToFile[3:]) // remove .kv/commit/vN from path
            for i := 0; i < len(dirToFile) - 1; i++ {
                pathToFile = pathToFile + dirToFile[i] + "/"
                // fmt.Println(pathToFile)
            }
            pathToFile = nextCommitVersion + pathToFile
            os.MkdirAll(pathToFile, 0700)
            commitedFile = pathToFile + dirToFile[len(dirToFile)-1]
        }

        f, err := os.Create(commitedFile)
        if err != nil {
            log.Fatal(err)
        }
        defer f.Close()

        copyFile(prevCommitFiles[i][0], commitedFile)
    }
}

func commitFiles() {
    // Don't know how I'm going to handle deleted files yet...

    // If first commit, make commitNum = 1 instead of 0
    commitNum := commitNumber() + 1

    commits, err := readLines(stagingAreaLocation())
    if err != nil {
        log.Println(err)
    }

    oldDir, err := os.Getwd()
    if err != nil {
    	log.Println(err)
    }

    os.Chdir(getRootDir())

    if (commitNum != 1) {
        copyFromPrevCommit()
    }

    for i := 0; i < len(commits); i++ {
        singleCommit := strings.Split(commits[i], ";")
        // singleCommit[0] - filepath
        // singleCommit[1] - modification date
        // singleCommit[2] - sha1 hash
        // singleCommit[3] - status (created/updated/deleted)
        commitVersion := ".kv/commit/v" + strconv.Itoa(commitNum) + "/"
        commitedFile := commitVersion + singleCommit[0]

        os.MkdirAll(commitVersion, 0700) // Create commit directory

        // If file not in root dir, create directories
        pathToFile := ""
        if (strings.Contains(singleCommit[0], "/")) {
            dirToFile := strings.Split(singleCommit[0], "/")
            for i := 0; i < len(dirToFile) - 1; i++ {
                pathToFile = pathToFile + dirToFile[i] + "/"
            }
            pathToFile = commitVersion + pathToFile
            os.MkdirAll(pathToFile, 0700)
            commitedFile = pathToFile + dirToFile[len(dirToFile)-1]
        }

        f, err := os.Create(commitedFile)
        if err != nil {
            log.Fatal(err)
        }
        defer f.Close()

        copyFile(singleCommit[0], commitedFile)
    }
    clearStagingArea()
    os.Chdir(oldDir)
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

            // Test argument, remove when implemented warning about untracked files
            case "staged":
                stagedFiles := getStagedFiles()
                for i := 0; i < len(stagedFiles); i++ {
                    fmt.Printf("name: %s\thash: %s\n", stagedFiles[i][0], stagedFiles[i][1])
                }
                os.Exit(0)

            // Test argument, remove when implemented warning about untracked files
            case "commited":
                commitedFiles := getCommitedFiles()
                for i := 0; i < len(commitedFiles); i++ {
                    fmt.Printf("name: %s\thash: %s\n", commitedFiles[i][0], commitedFiles[i][1])
                }
                os.Exit(0)

            // Test argument, remove when implemented warning about untracked files
            case "shortcommited":
                commitedFiles := getCommitedFilesShort()
                for i := 0; i < len(commitedFiles); i++ {
                    fmt.Printf("name: %s\thash: %s\n", commitedFiles[i][0], commitedFiles[i][1])
                }
                os.Exit(0)

            // Test argument, remove when implemented warning about untracked files
            case "all":
                allFiles := getAllFiles()
                for i := 0; i < len(allFiles); i++ {
                    fmt.Printf("name: %s\thash: %s\n", allFiles[i][0], allFiles[i][1])
                }
                os.Exit(0)

            // Test argument, remove when implemented warning about untracked files
            case "untracked":
                untrackedFiles := trackFiles()
                for i := 0; i < len(untrackedFiles); i++ {
                    fmt.Printf("%s\n", untrackedFiles[i])
                }
                os.Exit(0)

            case "add":
                if (len(os.Args) > i+1) {
                    for i := i+1; i < len(os.Args); i++ {
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
