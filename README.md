# Office Neighbours

A console application to find a person who is neighbors around the central default point for a radius defined as a console option.

A source for neighbors data is a text file which possible stored each line in JSON or XML format.

The application is automatically match a format of a source file.

## Requirements

Go 1.10+ is the only required to build the application.

No one third-party library was used in development.

Makefile was created to quickly build and test the application.

## Build

1. Using go:

```
go build -o office-neighbors
```

2. Using make:

```
make
```

Using a standard Go environment to build for another OS. Ex.:

```
GOOS=windows make
```

## Test

1. Using go:

```
go test ./...
```

2. Using make:

```
make testing
```

## Run

An application has two options:

1. **distance**: a distance in km which around neighbors finding.

2. **source file**: a path of a text file, contained neighbors list.

A default distance 100 km:
```
./office-neighbors ./customers.txt
```

or a different value of a distance:
```
./office-neighbors -distance 50 ./customers.txt
```

## Trade-off

1. Callback performance.

    Using callback for processing each record is one of the possible performance bottlenecks.
    The implementation was inspired by Rust iterators to show how possible to implemented iterators in Go.
    Go++ is gifting a hope of improving an iterator using, performance in the near future.

2. Automatically matching just by the first line.

    The application is automatically matching a format of a line using just a first not empty line, but maybe 
    be better trying to math more lines if the first has a broken format.