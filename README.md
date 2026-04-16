# gofilepath

Drop-in replacement for Go's `path/filepath` that works correctly with **cross-platform paths**.

## The Problem

Go's `filepath` package uses the **current OS separator** only. This breaks when you manipulate paths from a different OS:

```go
// Linux CLI handling a Windows remote path:
filepath.Base(`C:\Users\file.txt`)  // WRONG: returns "C:\\Users\\file.txt"
filepath.Join("/remote", `sub\dir`) // WRONG: returns "/remote/sub\\dir"
```

## What This Package Does

**Standard wrappers** — every `filepath` function is wrapped with `FromSlash` preprocessing, so `/` works everywhere:

```go
gofilepath.Base("/remote/file.txt")   // "file.txt" (same as filepath, but handles "/" on Windows)
gofilepath.Join("/a", "b")            // works correctly on Windows too
```

**Smart functions** — handle **both** `/` and `\` regardless of OS, output always uses `/`:

```go
// Works on ANY OS:
gofilepath.BaseSmart(`C:\Users\file.txt`)           // "file.txt"
gofilepath.DirSmart(`C:\Users\file.txt`)            // "C:/Users"
gofilepath.JoinSlash(`C:\Users`, `sub\dir`, "f.txt") // "C:/Users/sub/dir/f.txt"
gofilepath.RelSlash(`C:\a\b`, `C:\a\b\c\d`)         // "c/d"
```

## API

### Standard (OS-aware, wraps `filepath` with `FromSlash`)

| Function | Description |
|----------|-------------|
| `Base`, `Dir`, `Split`, `Join`, `Rel` | Same as `filepath.*` but accepts `/` on all platforms |
| `Clean`, `Ext`, `Abs`, `IsAbs` | Same behavior |
| `Walk`, `WalkDir` | Same behavior |
| `ToSlash`, `FromSlash` | Direct wrappers |
| `ToSlashSmart`, `FromSlashSmart` | Smart Windows drive letter handling (`/C/` <-> `C:/`) |

### Cross-Platform (OS-independent, both `/` and `\`)

| Function | Description |
|----------|-------------|
| `NormalizeSeparators(p)` | `\` -> `/` |
| `BaseSmart(p)` | Last path element |
| `DirSmart(p)` | Parent directory |
| `SplitSmart(p)` | Split into dir + file |
| `JoinSlash(elems...)` | Join with `/` |
| `RelSlash(base, targ)` | Relative path with `/` |
| `CleanSmart(p)` | Clean (`..`, `.`, `//`) |
| `ExtSmart(p)` | File extension |
| `BaseNoExtSmart(p)` | Filename without extension |
| `IsAbsSmart(p)` | Absolute check (incl. `C:\`) |

### Utilities

| Function | Description |
|----------|-------------|
| `JointSmart(fallback, elems...)` | Join preserving detected separator style |
| `RelSmart(base, targ)` | Rel preserving base's separator style |
| `ConvertPathSeparators(from, ref)` | Convert separators to match reference |
| `GetPathSeparator(p)` | Detect which separator a path uses |
| `PathIsExist`, `PathIsDir`, `PathIsFile` | Path type checks |
| `PathIsSymlink`, `PathIsSymlinkDir` | Symlink checks |
| `GetDrives()` | List drive letters (Windows) |
| `FindFilesMatch*` | Recursive file search with depth limit |

## Install

```
go get github.com/sonnt85/gofilepath@latest
```

## When to Use What

- **Local filesystem paths** -> standard functions (`Base`, `Join`, `Rel`)
- **Remote/network paths** (CLI -> server, different OS) -> cross-platform functions (`BaseSmart`, `JoinSlash`, `RelSlash`)
- **Preserving user's separator style** -> `JointSmart`, `RelSmart`

## Author

**sonnt85** — [thanhson.rf@gmail.com](mailto:thanhson.rf@gmail.com)

