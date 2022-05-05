package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	pg_query "github.com/pganalyze/pg_query_go/v2"
)

type AppState struct {
	Input    string
	Host     string
	Port     int64
	User     string
	Password string
	DBName   string
}

type GatorQuery struct {
	query string
	done  bool
	err   error
}

func bail(err error) {
	if err != nil {
		panic(err)
	}
}

func (app *AppState) connectionString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		app.User, app.Password, app.Host, app.Port, app.DBName)
}

func hasString(seq []string, key string) bool {
	for _, item := range seq {
		if item != key {
			continue
		}
		return true
	}
	return false
}

func (app *AppState) executeQueries(queries []GatorQuery) {
	log.Println(app.connectionString())
	db, err := sql.Open("postgres", appState.connectionString())
	bail(err)
	defer db.Close()

	fmt.Println("Checking connection to db...")

	err = db.Ping()
	bail(err)

	fmt.Println("Connected")

	fmt.Println("Executing Queries")

	var toExecQueries []GatorQuery = queries

	for i := 1; i <= 10; i++ {
		for queryIndex := len(toExecQueries) - 1; queryIndex >= 0; queryIndex-- {
			if toExecQueries[queryIndex].done {
				continue
			}

			_, err := db.Exec(toExecQueries[queryIndex].query)
			if err != nil {
				toExecQueries[queryIndex].err = err
			} else {
				toExecQueries[queryIndex].done = true
				toExecQueries[queryIndex].err = nil
			}
		}
	}

	var failedQueries = 0

	for _, query := range toExecQueries {
		if !query.done {
			failedQueries += 1
		}
	}

	if failedQueries > 0 {
		log.Println("There were errors")
		for _, query := range toExecQueries {
			if !query.done {
				log.Println(query.err)
			}
		}
	} else {
		fmt.Println("Executed Statements successfully")
	}

}

func (app *AppState) Gator() {

	var queriesToExec []GatorQuery

	fileData, err := os.ReadFile(app.Input)
	bail(err)

	result, err := pg_query.Parse(string(fileData))

	if err != nil {
		panic(err)
	}

	for _, stmt := range result.Stmts {
		var raw []*pg_query.RawStmt
		raw = append(raw, stmt)
		singleStatement := &pg_query.ParseResult{
			Stmts: raw,
		}
		deparsed, err := pg_query.Deparse(singleStatement)
		bail(err)
		queriesToExec = append(queriesToExec, GatorQuery{
			query: deparsed,
			done:  false,
			err:   nil,
		})
	}

	// log.Println(queriesToExec[0])

	appState.executeQueries(queriesToExec)
}

var appState *AppState

func cli() {
	appState = &AppState{}
	inputFile := flag.String("file", "", "sql file to run")
	host := flag.String("host", "localhost", "sql file to run")
	port := flag.Int64("port", 5432, "sql file to run")
	user := flag.String("user", "postgres", "sql file to run")
	password := flag.String("password", "", "sql file to run")
	dbName := flag.String("db", "postgres", "sql file to run")
	flag.Parse()

	appState.Input = *inputFile
	appState.Host = *host
	appState.Port = *port
	appState.User = *user
	appState.Password = *password
	appState.DBName = *dbName
}

func main() {
	cli()
	appState.Gator()
}
