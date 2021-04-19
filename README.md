# git-statuses

**git-statuses** finds local git repositories and show statuses of them.

[![PkgGoDev](https://pkg.go.dev/badge/kyoh86/git-statuses)](https://pkg.go.dev/kyoh86/git-statuses)
[![Go Report Card](https://goreportcard.com/badge/github.com/kyoh86/git-statuses)](https://goreportcard.com/report/github.com/kyoh86/git-statuses)
[![Coverage Status](https://img.shields.io/codecov/c/github/kyoh86/git-statuses.svg)](https://codecov.io/gh/kyoh86/git-statuses)
[![Release](https://github.com/kyoh86/git-statuses/workflows/Release/badge.svg)](https://github.com/kyoh86/git-statuses/releases)

## Install

```
go get github.com/kyoh86/git-statuses
```

## Usage

```sh
git statuses [--json] (pathspec)
```

or

```sh
cd (pathspec); git statuses [--json]
```

## Result

### Sample

```console
$ git statuses
-1 +2  U github.com/kyoh86/git-statuses
-1     U github.com/kyoh86/go-jsonrider
   +3 M  github.com/kyoh86/go-pcre
      M  github.com/kyoh86/gogh
      MU github.com/kyoh86/mogelo
```

### Format
```
(status) (repository path)
(status) (repository path)
:        :
```

|status|description                     |
|:----:|--------------------------------|
|`+n`  |Contains ahead commits          |
|`-n`  |Contains behind commits         |
|`M`   |Contains deleted/modified files |
|`U`   |Contains untracked files        |

# LICENSE

[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg)](http://www.opensource.org/licenses/MIT)

This is distributed under the [MIT License](http://www.opensource.org/licenses/MIT).
