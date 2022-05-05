<h1 align="center">
    Gator
</h1>
<p align="center">Batch SQL Query runner for Postgres</p>

I deal with moving psql dump files and sql data insertion files around a lot.
It's mostly to seed data when a new developer joins the company and each of
them have a different tool/GUI for handling postgres and not all of them support
running each query separately. This creates an issue when the export of data by most
apps is done in an alphabetical manner and so foreign key issues are bound to exist.

Gator simplifies my work by separating each query into it's own execution context and
running each query about 10 times. This makes sure that all dependent queries are run and if there's still errors, you'll see them in the terminal.

**NOTE:** Gator, is a personal project which I wasn't build for the public but the source code is here cause there's nothing to "close source" here.

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
