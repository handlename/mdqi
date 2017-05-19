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
- [ ] slash command
    - [x] `/quit`
    - [ ] `/tag (set|remove|show)`
    - [ ] `/display (vertical|horizontal)`
- [ ] load configuration file
- [ ] CircleCI
- [ ] command help

## Licence

[MIT](https://github.com/handlename/mdqi/blob/master/LICENSE)

## Author

[handlename](https://github.com/handlename)
