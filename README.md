# kv - Klear version control system

Seeing if I can make a semi-good basic version control system with go.

“As an analogy, imagine you are building a version control system like SVN or Git.
When a user commits a file for the first time, the system saves the whole file to disk.
Subsequent commits, reflecting changes to that file, might save only the delta — that is,
just the lines that were added, changed, or removed. Then, when the user checks out a certain
version, the system opens the version-0 file and applies all subsequent deltas, in order, to
derive the version the user asked for.”

Inspiration:
[here](https://levelup.gitconnected.com/how-was-i-build-a-version-control-system-vcs-using-pure-go-83ec8ec5d4f4)

## Simple breakdown

- Save a commited file to `.kv` folder

Tree structure of `.kv` folder should look something like this:
```
.kv
├── final
│ ├── v1
│ └── v2
├── commit
│ ├── v1
│ └── v2
├── staging-area.txt
└── status.txt
```

v1, v2, ... are commit versions.  
The `staging-area.txt` contains the files that are going to be in the next commit.
It should be a basic format (possibly csv) that shows this information:  
```
filepath;modification date;status
```

For example:
```
"src/main.c;2023-05-11 05:42:15;Created"
".kvignore;2023-05-16 05:11:04;Created"
"README.md;2022-04-14 05:49:09;Updated"
```

`final` is a folder that includes the result of merging all files with the specified
commit version. For example, `final/v2` includes a combination of `commit/v1` + `commit/v2`.

`status` is a file that keeps track of all files persistently. I clear the
contents of staging area after the successful commit.

Now to stop copying the inspiration page, I won't do this with the [cobra](https://github.com/spf13/cobra)
library. I'll try to do this with as few dependencies as I can. 

I'll keep track of deleted files

### How it starts

```
.kv
├── final
│ ├── v1
│ └── v2
└── commit
  ├── v1   <---- It starts here
  └── v2
```

The first commit ever, gets put in this `v1` directory. This is the only time that `kv`
will copy any file contents. All commited files will be copied into the `v1` directory.
Every other subsequent version/commit, will contain only deltas/diffs/changes done to the
`v1` directory. Every file will be hashed with sha1, so that later comparisons can be made.

### Hashing files

In each commit folder, for example v2:
```
.kv
├── final
│ ├── v1
│ └── v2
└── commit
  ├── v1
  └── v2  <---- This one
```

Each file that contains a delta/diff, will contain the hash of itself.

## Features I want to implement

- [X] Init the directories and files
- [X] Status, show what's in the staging area
- [X] Add files to the staging area
- [ ] Commit files/changes
- [ ] Implement SHA1 hashing of commits for integrity checks
- [ ] View history
- [ ] Track deleted files
- [ ] Rollback to previous commits
- [ ] Decent UI/UX
