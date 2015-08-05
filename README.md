git-statuses
=====

**git-statuses** finds local git repositories and show statuses of them.

## Installation

`go get github.com/kyoh86/git-statuses/...`

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

## License

Released under the MIT License.
