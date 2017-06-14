[![CircleCI](https://circleci.com/gh/handlename/mdqi.svg?style=svg)](https://circleci.com/gh/handlename/mdqi)

# mdqi

`mdqi` is interactive interface for mdq.
It supports:

- Display query results with table based layout
- Query history
- Remember tag

## Example

// TODO: gif animation

## Installation

Download binary from releases and put it on `$PATH`ed place.

Or

```
$ go get github.com/handlename/mdqi/cmd/mdqi
```

## Usage

```
$ mdqi
> select * from items;
```

## TODO

- [x] run mdq through mdqi
- [x] display results as table
- [x] convert mdq formatted JSON as table passed by stdin
- [x] receive command line input
- [x] query history
- [x] logger
- [x] ignore empty query
- [x] CI on CircleCI
- [x] slash command
  - [x] `/quit`
  - [x] `/tag (set|remove|show)`
  - [x] `/help [category [name]]`
  - [x] `/display (vertical|horizontal)`
  - [x] `/v` alias for toggle display
- [x] load configuration file
- [x] -version option
- [x] multi line input
- [x] deploy from CircleCI
- [ ] Ctrl-Z add shortcut to liner
  - pending.
    github.com/peterh/liner converts Ctrl-Z to beep,
    so SIGTSTP will not reach to signal handler.
  - https://github.com/peterh/liner/blob/88609521dc4b6c858fd4c98b628147da928ce4ac/line.go#L852-L856

## Licence

[MIT](https://github.com/handlename/mdqi/blob/master/LICENSE)

## Author

[handlename](https://github.com/handlename)
