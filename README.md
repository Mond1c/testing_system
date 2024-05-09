# Simple testing system for contests

## Usage:

You can run the example using this command:

```
go run ./cmd -config examples/input.json -languages examples/languages.json
```

But you need to specify the path to the test folder for each problem.
!Important one problem = one folder.

In the folder you need to create files with test input and output, and the name of the files should satisfy this template:

```
in1.in or in1.out
1 - number of the test
```
