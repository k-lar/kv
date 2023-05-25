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

## Installation

Building from source:
```console
git clone https://gitlab.com/k_lar/kv; cd kv/

# Build the program:
go build main.go

# Rename binary (optional)
mv main kv
```

## Usage

Basic usage:
```
# Initialize repository in the current directory
kv init

# Add files to staging area
kv add <path/to/file>

# Show what's in the staging area + modified and untracked files
kv status

# Commit changes to the repository
kv commit
```

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
filepath;modification date;hash;status
```

For example:
```
"src/main.c;2023-05-11 05:42:15;sa31sv35bc92js84jg33;created"
".kvignore;2023-05-16 05:11:04;hw81ks96ms24dl80sm12;created"
"README.md;2022-04-14 05:49:09;bl12sh56ka93tl22xc56;updated"
```

`final` is a folder that includes the result of merging all files with the specified
commit version. For example, `final/v2` includes a combination of `commit/v1` + `commit/v2`.

`status` is a file that keeps track of all files persistently. I clear the
contents of staging area after the successful commit.

Now to stop copying the inspiration page, I won't do this with the [cobra](https://github.com/spf13/cobra)
library. I'll try to do this with as few dependencies as I can. 

I'll keep track of deleted files

### Stages

```
+----------------+       +---------------+
|                |       |               |
|    Staging     |  -->  |   Commited    |
|                |       |               |
+----------------+       +---------------+
```

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

The first commit ever, gets put in this `v1` directory. All staged files will be copied
into the `v1` directory. Every file will be hashed with sha1, so that later comparisons can
be made.

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

Each file should be hashed so that it can be validated.
Each version (`v1`, `v2`, ...) will contain the following files:  
(`v1` directory can only have files that have been created, it would make no sense if it contained
deleted files or updated files.)
```
.kv
└── commit
  ├── v1
  │ ├── .kv_commit  <-- (Contains hash of the whole thing, time, and the commit messege + author)
  │ ├── file1
  │ ├── file2
  │ └── ...
  └── v2
    ├── .kv_commit  <-- (Contains hash of the whole thing + parent + above things)
    ├── file1
    ├── file2
    └── ...
```

## How commits know which files to take from previous versions

How I think it can work:

1. Check which files are in previous commit
2. Check which files are staged
3. Created files get put in the newest commit without trouble
4. Put all previously commited files in an array:
    - If staging doesn't contain filename of commited file, copy file from previous commit
       to newest commit
    - If staging contains commited filename (updated), copy it from repo to newest commit
    - If staging has a deleted file, don't copy it from anywhere

## Known issues

- Can not commit file if you're inside the `.kv/` directory
- Only absolute paths from repository root can  be added in staging

## Features I want to implement

- [X] Init the directories and files
- [X] Status, show what's in the staging area
- [X] Add files to the staging area
- [X] Commit files/changes
- [X] Implement SHA1 hashing of commits for integrity checks
- [ ] Show diff between files (make a builtin diff)
- [ ] View history
- [ ] Track deleted files
- [ ] Rollback to previous commits
- [ ] Decent UI/UX
