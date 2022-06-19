# Go Tour guide

## Install

macOS:

```bash
brew install go
```

## Compile and run

Inside of the directory where the main package is located type and run

```bash
go run .
```

Or if not a go module(doesn't contain `go.mod`)

```bash
go run <filename>.go
```

## How to setup a module

Steps:

* To create the module: `go mod init <name>`
* Fetch module dependencies `go mod tidy`

More info on the [link](https://go.dev/doc/tutorial/getting-started)

## Packages, variables, and functions

[Link](./0.packages-vars-functions/packagesVarsFunctions.go) to first program

## Flow control statements: for, if, else, switch and defer

[Link](./1.control-flow-statements/controlFlow.go) to control flow

## More types: structs, slices, and maps

[Link](./2.types-structs-slices-maps/types.go) to control flow

## Methods and interfaces

[Link](./3.methods-interfaces/methodsAndInterfaces.go) to methods and interfaces

## Generics

[Link](./4.gemerics/generics.go) to generics