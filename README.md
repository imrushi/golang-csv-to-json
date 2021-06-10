# Convert CSV to JSON

## Package is used

[Cobra](https://github.com/spf13/cobra)

## To run program

- Clone the project
- Run `go get` command to get all packages
- To build project run `go build -o ./build/csv2json` (the build file extention will be depending on your system if you build on Windows it will generate `.exe`...)
- To build for specific OS `GOOS=linux GOARCH=amd64 go build` there are may support OS and ARCHITECTURES check this [link](https://gist.github.com/asukakenji/f15ba7e588ac42795f421b48b8aede63).
- To run build file `./build/csv2json` or `cd build` >> `./csv2json`

````$ ./build/csv2json.exe
csv2json: It is a CLI application for convert Csv file to Json file with multiple options.

Usage:
  csv2json [command]

Available Commands:
  convert     To convert file from CSV to JSON
  help        Help about any command

Flags:
  -h, --help               help for csv2json
  -p, --pretty             Generate pretty JSON
  -s, --separator string   Column Separator (default "comma")

Use "csv2json [command] --help" for more information about a command.```
````
