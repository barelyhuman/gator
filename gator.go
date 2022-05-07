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
	Input         string
	Host          string
	Port          int64
	User          string
	Password      string
	DBName        string
	SyncSequences bool
	DBConnection  *sql.DB
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

func (app *AppState) connectDB() {
	fmt.Println("Connecting with the following Creds")
	fmt.Printf("Host:%s\n Port:%d\n Database:%s User:%s\n",
		appState.Host,
		appState.Port,
		appState.DBName,
		appState.User,
	)
	db, err := sql.Open("postgres", appState.connectionString())
	bail(err)
	fmt.Println("Checking connection to db...")
	err = db.Ping()
	bail(err)
	fmt.Println("Connected")
	app.DBConnection = db
}

func (app *AppState) syncSequences() {
	fmt.Println("Syncing Sequences")
	_, err := app.DBConnection.Exec("DO $$ DECLARE i TEXT; BEGIN FOR i IN (select table_name from information_schema.tables where table_catalog='YOUR_DATABASE_NAME' and table_schema='public') LOOP EXECUTE 'Select setval('''||i||'id_seq'', (SELECT max(id) as a FROM ' || i ||')+1);'; END LOOP; END$$;")
	bail(err)
	fmt.Println("Synced")
}

func (app *AppState) executeQueries(queries []GatorQuery) {
	fmt.Println("Executing Queries")

	var toExecQueries []GatorQuery = queries

	for i := 1; i <= 10; i++ {
		for queryIndex := len(toExecQueries) - 1; queryIndex >= 0; queryIndex-- {
			if toExecQueries[queryIndex].done {
				continue
			}

			_, err := app.DBConnection.Exec(toExecQueries[queryIndex].query)
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

func (app *AppState) runSeedFile() {
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

	app.executeQueries(queriesToExec)
}

func (app *AppState) Gator() {

	app.connectDB()

	if app.DBConnection != nil {
		defer app.DBConnection.Close()
	}

	if len(app.Input) > 0 {
		app.runSeedFile()
	}

	if app.SyncSequences {
		app.syncSequences()
	}

}

var appState *AppState

func cli() {
	appState = &AppState{}
	inputFile := flag.String("file", "", "sql file to run")
	syncSequences := flag.Bool("sync-sequences", false, "Sync Sequences")
	host := flag.String("host", "localhost", "host address")
	port := flag.Int64("port", 5432, "port to connect")
	user := flag.String("user", "postgres", "user for authentication")
	password := flag.String("password", "", "password for authentication")
	dbName := flag.String("db", "postgres", "database name to run the file against")
	flag.Parse()

	appState.Input = *inputFile
	appState.Host = *host
	appState.Port = *port
	appState.User = *user
	appState.SyncSequences = *syncSequences
	appState.Password = *password
	appState.DBName = *dbName
}

func main() {
	cli()
	appState.Gator()
}
