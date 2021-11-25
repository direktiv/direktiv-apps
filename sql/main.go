package main

import (
	// "database/sql"

	"fmt"
	"net/http"
	"sync"

	ru "github.com/direktiv/direktiv-apps/pkg/reusable"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"

	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

type queries struct {
	Transaction bool   `json:"tx"`
	Query       string `json:"query"`
}

type input struct {
	Debug      string    `json:"debug"`
	DBType     string    `json:"db"`
	Connection string    `json:"connection"`
	Queries    []queries `json:"queries"`
	Fail       bool      `json:"fail"`
}

type queryResult struct {
	Success bool                     `json:"success"`
	Error   string                   `json:"error"`
	Result  []map[string]interface{} `json:"result"`
}

type queryResponse struct {
	Success bool           `json:"success"`
	Results []*queryResult `json:"results"`
}

const (
	dbError     = "io.direktiv.db"
	querryError = "io.direktiv.query"
)

var (
	dbConn     *sqlx.DB
	dbMutex    sync.Mutex
	connString string
)

func handleOpen(dbType, conn string) error {

	var err error

	if dbConn != nil {
		return nil
	}

	fmt.Printf("DB %v\n", dbType)

	// check db type
	switch dbType {
	case "mysql":
		fallthrough
	case "mssql":
		fallthrough
	case "postgres":
	default:
		return fmt.Errorf("database type not supported")
	}

	dbConn, err = sqlx.Open(dbType, conn)
	if err != nil {
		return err
	}

	err = dbConn.Ping()
	if err != nil {
		return err
	}

	connString = conn

	return nil

}

func reqHandler(w http.ResponseWriter, r *http.Request, ri *ru.RequestInfo) {

	var input input

	err := ru.Unmarshal(&input, true, r)
	if err != nil {
		ru.ReportError(w, ru.UnmarshallError, err)
		return
	}

	if input.Debug != "" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	// if connstring changed close it and reopen
	if dbConn != nil && connString != input.Connection {
		dbMutex.Lock()
		dbConn.Close()
		dbConn = nil
		dbMutex.Unlock()
	}

	// open db conn
	if dbConn == nil {
		dbMutex.Lock()
		defer dbMutex.Unlock()

		err := handleOpen(input.DBType, input.Connection)
		if err != nil {
			ru.ReportError(w, dbError, err)
			return
		}
		ri.Logger().Infof("database connection established")
	}

	resp := &queryResponse{
		Success: true,
	}

	for i := range input.Queries {
		q := input.Queries[i]

		ri.Logger().Infof("querying: %v", q.Query)

		var qr *queryResult
		// if transaction enabled
		if q.Transaction {
			qr, err = runTx(q.Query, ri)
		} else {
			qr, err = runNoTx(q.Query, ri)
		}

		// if fail set to true we retun with first unsuccessful query
		if err != nil && input.Fail {
			ru.ReportError(w, querryError, err)
			return
		}

		// append the result to queries
		// and mark as total as unsuccessful
		resp.Results = append(resp.Results, qr)
		if !qr.Success {
			resp.Success = false
		}

	}

	ru.ReportResult(w, resp)

}

func runNoTx(q string, ri *ru.RequestInfo) (*queryResult, error) {
	qr := &queryResult{
		Success: false,
	}

	rows, err := dbConn.Queryx(q)
	if err != nil {
		qr.Error = err.Error()
		return qr, err
	}

	// convert to json
	if rows != nil {
		tableData := make([]map[string]interface{}, 0)
		for rows.Next() {
			entry := make(map[string]interface{})
			err := rows.MapScan(entry)
			if err != nil {
				qr.Error = err.Error()
				return qr, err
			}

			tableData = append(tableData, entry)
		}
		qr.Result = tableData
	}

	// if no error query was ok
	qr.Success = true

	return qr, nil

}

func runTx(q string, ri *ru.RequestInfo) (*queryResult, error) {

	qr := &queryResult{
		Success: false,
	}

	tx, err := dbConn.Beginx()
	if err != nil {
		qr.Error = err.Error()
		return qr, err
	}
	defer tx.Commit()

	rows, err := tx.Queryx(q)
	if err != nil {
		qr.Error = err.Error()
		return qr, err
	}

	// convert to json
	if rows != nil {
		tableData := make([]map[string]interface{}, 0)
		for rows.Next() {
			entry := make(map[string]interface{})
			err := rows.MapScan(entry)
			if err != nil {
				qr.Error = err.Error()
				return qr, err
			}
			tableData = append(tableData, entry)
		}
		qr.Result = tableData
	}

	// if no error query was ok
	qr.Success = true

	return qr, nil
}

// closes db on close
func handleShutdown() {

	if dbConn != nil {
		fmt.Printf("closing database connection\n")
		defer dbConn.Close()
	}

}

func main() {
	handleShutdown()
	ru.StartServer(reqHandler, handleShutdown)
}
