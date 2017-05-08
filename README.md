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

- [ ] run mdq through mdqi
- [ ] display results as table
- [ ] load configuration file
- [ ] remember tag
- [ ] query history
- [ ] convert mdq formatted JSON as table passed by stdin

## Licence

[MIT](https://github.com/handlename/mdqi/blob/master/LICENSE)

## Author

[handlename](https://github.com/handlename)
