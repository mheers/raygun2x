# Raygun2x

## Overview
Raygun2x is a CLI tool that converts [raygun](https://github.com/mheers/opa-raygun) test suites into either an OpenAPI specification or a Postman Collection. It allows users to specify an input folder, choose the conversion mode, and decide whether to print the output to stdout or save it to a file.

## Installation
Ensure you have Go installed, then build the tool using:

```sh
make binary
```

## Usage
Run the command with the following options:

```sh
./raygun2x convert --input <folder> --mode <openapi|postman> --output <file> --stdout
```

### Options:
- `-i, --input`   (default: `presets/`) - The folder containing test suites.
- `-m, --mode`    (default: `openapi`) - The conversion mode: `openapi` or `postman`.
- `-o, --output`  (optional) - The file to save the output. If not specified, output is printed to stdout.
- `-s, --stdout`  (optional) - Force printing to stdout instead of writing to a file.

### Examples
Convert test suites to an OpenAPI spec and save to a file:
```sh
./raygun2x convert --input my-tests/ --mode openapi --output openapi.yaml
```

Convert test suites to a Postman Collection and print to stdout:
```sh
./raygun2x convert --input my-tests/ --mode postman --stdout
```

## License
This project is licensed under the MIT License.
