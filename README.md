# git-statuses

**git-statuses** finds local git repositories and show statuses of them.

[![Go Report Card](https://goreportcard.com/badge/github.com/kyoh86/git-statuses)](https://goreportcard.com/report/github.com/kyoh86/git-statuses)
[![Coverage Status](https://img.shields.io/codecov/c/github/kyoh86/git-statuses.svg)](https://codecov.io/gh/kyoh86/git-statuses)

## Install

```
go get github.com/kyoh86/git-statuses
```

## Usage

```sh
git statuses [--detail] [--relative] (pathspec)
```

or

```sh
GIT_STATUSES_TARGET=(pathspec); git statuses [--detail] [--relative]
```

or

```sh
cd (pathspec); git statuses [--detail] [--relative]
```

## Result

### Sample
```
~ $ git statuses
 + github.com/kyoh86/git-statuses
 + github.com/kyoh86/go-jsonrider
M  github.com/kyoh86/go-pcre
M  github.com/kyoh86/gogh
M+ github.com/kyoh86/mogelo
```

### Format
```
(status) (repository path)
(status) (repository path)
:        :
```

|status|description                     |
|:----:|--------------------------------|
|`M`   |Contains deleted/modified files |
|`+`   |Contains untracked files        |

# LICENSE

[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg)](http://www.opensource.org/licenses/MIT)

This is distributed under the [MIT License](http://www.opensource.org/licenses/MIT).
