<h1 align="center">
    Gator
</h1>
<p align="center">Batch SQL Query runner for Postgres</p>

Having new developers setup their databases is a common thing for me to handle and this has a lot of friction due to the
existing tools that devs use, some have TablePlus, some have postico and then there's other web based solutions as well.

Most of these export the data in an alphabetical order or in a foreign key dependency order. The 2nd one normally works fine
with having to just run the errored out queries once more. The issue is that not all these tools offer a way to skip errors
on queries and just run the next query. Which, is necessary when you are dealing with alphabetical(table name) order based data exports.

Gator simplifies my work by separating each query into it's own execution context and
running each query about 10 times (if it fails). This makes sure that all dependent queries are run and if there's still errors, you'll see them in the terminal.

**NOTE:** Gator, is a personal project which wasn't build for the public but the source code is here cause there's nothing to "close source"

## Install

Like everything else that I've written in go, the [releases](https://github.com/barelyhuman/releases) page has binaries for common Unix systems.

For other \*nix systems, I'd recommend building from source since the library depends on
[pg_query_go](https://github.com/pganalyze/pg_query_go) and cross compiling for each
operating system isn't feasible right now. You are free to Raise PR's for adding build scripts for your particular system.

## Usage

```sh
Usage of gator:
  -db string
        database name to run the file against (default "postgres")
  -file string
        sql file to run
  -host string
        host address (default "localhost")
  -password string
        password for authentication
  -port int
        port to connect (default 5432)
  -user string
        user for authentication (default "postgres")
```

## License

[MIT](https://github.com/barelyhuman/gator/blob/dev/license) | [Reaper](https://github.com/barelyhuman)
