# go-support
Go program for finding and supporting sponsorable awesome Go packages in your project's module dependencies (go.mod).
Let's support the Go developers and make the ecosystem a better place!

## For developers 
Please help us to update the list of sponsorable packages (your package!) by opening PR on the [`list.go`](https://github.com/joshuabezaleel/go-support/blob/master/list.go) file. 
Put (1) your name and your projects in the [`authorToProjectsList`](https://github.com/joshuabezaleel/go-support/blob/master/list.go#L3) variable and then (2) your name and your sponsor URLs in the [`authorToSponsorURLsList`](https://github.com/joshuabezaleel/go-support/blob/master/list.go#L25) variable by conforming to the previous examples.

## Usage
1. Make sure Go 1.12 is installed in your machine.
2. Get the program using `go get -u github.com/joshuabezaleel/go-support`.
3. Run `go-support` in the root of your project modules.
4. Open the URLs in the browser and make donations!

## Prior Art
This project is highly inspired by the kind efforts of [feross'](https://github.com/feross) [thanks](https://github.com/feross/thanks), GitHub sponsor, and [npm fund](https://github.com/npm/cli/pull/273). Thank you very much for taking the first steps. 