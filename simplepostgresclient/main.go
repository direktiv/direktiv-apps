package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strings"

	_ "github.com/lib/pq"
	"github.com/direktiv/direktiv-apps/pkg/direktivapps"
)

var debug bool

type operationDoerFunc func(tx *sql.Tx, table string) (interface{}, error)
type operationValidatorFunc func(map[string]interface{}) (operationDoerFunc, error)

var operations = map[string]operationValidatorFunc{
	"delete": deleteOpValidator,
	"insert": insertOpValidator,
	"select": selectOpValidator,
	"update": updateOpValidator,
}

func wheres(input map[string]interface{}, required bool) (map[string]interface{}, error) {

	x, ok := input["where"]
	if !ok {
		if !required {
			return map[string]interface{}{}, nil
		}
		return nil, fmt.Errorf("required parameter 'where' missing")
	}

	wheres, ok := x.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("'where' parameter must be a JSON object")
	}

	if len(wheres) == 0 {
		if !required {
			return map[string]interface{}{}, nil
		}
		return nil, fmt.Errorf("'where' parameter must specify at least one condition")
	}

	for k, v := range wheres {
		switch v.(type) {
		case map[string]interface{}:
			return nil, fmt.Errorf("'where' parameter has non-primitive parameter '%s' (it's a JSON object)", k)
		case []interface{}:
			return nil, fmt.Errorf("'where' parameter has non-primitive parameter '%s' (it's an array)", k)
		default:
		}
	}

	return wheres, nil

}

func wheresString(wheres map[string]interface{}) string {

	var conditions []string

	for k, v := range wheres {
		key := `"` + strings.ReplaceAll(fmt.Sprintf("%v", k), `"`, `""`) + `"`
		val := "'" + strings.ReplaceAll(fmt.Sprintf("%v", v), "'", "''") + "'"
		expr := fmt.Sprintf("%s=%s", key, val)
		if v == nil {
			expr = fmt.Sprintf("%s IS NULL", key)
		}
		conditions = append(conditions, expr)
	}

	return strings.Join(conditions, " AND ")

}

type deleteOp struct {
	wheres map[string]interface{}
	aid    string
}

func deleteOpValidator(input map[string]interface{}) (operationDoerFunc, error) {

	var err error
	op := new(deleteOp)

	expect := map[string]bool{
		"type":  true,
		"where": true,
	}

	for k := range input {
		_, expected := expect[k]
		if !expected {
			return nil, fmt.Errorf("unexpected parameter on 'delete' operation: %s", k)
		}
	}

	op.wheres, err = wheres(input, true)
	if err != nil {
		return nil, err
	}

	return op.do, nil

}

func (op *deleteOp) do(tx *sql.Tx, table string) (interface{}, error) {

	wheres := wheresString(op.wheres)
	query := fmt.Sprintf(`DELETE FROM %s WHERE %s`, table, wheres)
	result, err := tx.Exec(query)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to delete: %w", err)
	}

	k, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to lookup rows affected: %w", err)
	}

	return map[string]int64{
		"rowsAffected": k,
	}, nil

}

type insertOp struct {
	records []map[string]interface{}
	aid     string
}

func insertOpValidator(input map[string]interface{}) (operationDoerFunc, error) {

	op := new(insertOp)

	expect := map[string]bool{
		"type": true,
		"data": true,
	}

	for k := range input {
		_, expected := expect[k]
		if !expected {
			return nil, fmt.Errorf("unexpected parameter on 'insert' operation: %s", k)
		}
	}

	var records []map[string]interface{}

	x, ok := input["data"]
	if !ok {
		return nil, fmt.Errorf("required parameter 'data' missing")
	}

	type1, ok1 := x.([]map[string]interface{})
	if ok1 {
		records = type1
	}

	type2, ok2 := x.(map[string]interface{})
	if ok2 {
		records = append(records, type2)
	}

	if !ok1 && !ok2 {
		return nil, fmt.Errorf("'data' parameter must be a JSON object or an array of JSON objects")
	}

	if len(records) == 0 {
		return nil, fmt.Errorf("'data' parameter must specify at least one record to add to the database")
	}

	for i, rec := range records {
		for k, v := range rec {
			switch v.(type) {
			case map[string]interface{}:
				return nil, fmt.Errorf("'data' parameter has non-primitive parameter '%s' in record %d (it's a JSON object)", k, i)
			case []interface{}:
				return nil, fmt.Errorf("'data' parameter has non-primitive parameter '%s' in record %d (it's an array)", k, i)
			default:
			}
		}
	}

	op.records = records

	return op.do, nil

}

func (op *insertOp) do(tx *sql.Tx, table string) (interface{}, error) {

	for i, record := range op.records {

		var keys, vals, obscuredVals []string

		for k, v := range record {
			key := `"` + strings.ReplaceAll(fmt.Sprintf("%v", k), `"`, `""`) + `"`
			val := "'" + strings.ReplaceAll(fmt.Sprintf("%v", v), "'", "''") + "'"
			if v == nil {
				val = "NULL"
			}
			obscuredVal := "'****'"
			keys = append(keys, key)
			vals = append(vals, val)
			obscuredVals = append(obscuredVals, obscuredVal)
		}

		ks := strings.Join(keys, ", ")
		vs := strings.Join(vals, ", ")
		// obscuredVs := strings.Join(obscuredVals, ", ")

		query := fmt.Sprintf(`INSERT INTO %s (%s) VALUES (%s)`, table, ks, vs)
		// obscuredQuery := fmt.Sprintf(`INSERT INTO %s (%s) VALUES (%s)`, table, ks, obscuredVs)

		_, err := tx.Exec(query)
		if err != nil {
			return nil, fmt.Errorf("failed to insert record %d: %w", i, err)
		}

	}

	return map[string]interface{}{}, nil

}

type selectOp struct {
	fields   []string
	wildcard bool
	wheres   map[string]interface{}
	aid      string
}

func selectOpValidator(input map[string]interface{}) (operationDoerFunc, error) {

	var err error
	op := new(selectOp)

	expect := map[string]bool{
		"type":   true,
		"where":  true,
		"fields": true,
	}

	for k := range input {
		_, expected := expect[k]
		if !expected {
			return nil, fmt.Errorf("unexpected parameter on 'select' operation: %s", k)
		}
	}

	op.wheres, err = wheres(input, false)
	if err != nil {
		return nil, err
	}

	var fields []string

	x, ok := input["fields"]
	if !ok {
		return nil, fmt.Errorf("required parameter 'fields' missing")
	}

	s, ok1 := x.(string)
	y, ok2 := x.([]interface{})
	if (!ok1 && !ok2) || (ok1 && s != "*") {
		return nil, fmt.Errorf(`'fields' parameter must be "*" or a JSON array of strings`)
	}

	if ok1 {
		op.wildcard = true
		return op.do, nil
	}

	if len(y) == 0 {
		return nil, fmt.Errorf("'fields' parameter must request at least one column")
	}

	for i, z := range y {
		s, ok := z.(string)
		if !ok {
			return nil, fmt.Errorf("'fields' parameter has non-string element %d", i)
		}
		fields = append(fields, s)
	}

	op.fields = fields

	return op.do, nil

}

func (op *selectOp) do(tx *sql.Tx, table string) (interface{}, error) {

	var fieldsStr string
	if op.wildcard {
		fieldsStr = "*"
	} else {
		var fields []string
		for _, field := range op.fields {
			fieldName := `"` + strings.ReplaceAll(field, `"`, `""`) + `"`
			fields = append(fields, fieldName)
		}
		fieldsStr = strings.Join(fields, ", ")
	}

	var query string
	if len(op.wheres) == 0 {
		query = fmt.Sprintf(`SELECT %s FROM %s`, fieldsStr, table)
	} else {
		wheres := wheresString(op.wheres)
		query = fmt.Sprintf(`SELECT %s FROM %s WHERE %s`, fieldsStr, table, wheres)
	}

	rows, err := tx.Query(query)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to select: %v", err)
	}
	defer rows.Close()

	var outputRecords = make([]map[string]interface{}, 0)

	for rows.Next() {

		cols, err := rows.Columns()
		if err != nil {
			return nil, fmt.Errorf("failed to load columns from response: %v", err)
		}

		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i, _ := range cols {
			columnPointers[i] = &columns[i]
		}

		if err := rows.Scan(columnPointers...); err != nil {
			return nil, fmt.Errorf("failed to load results from query: %v", err)
		}

		m := make(map[string]interface{})
		for i, colName := range cols {
			val := columnPointers[i].(*interface{})
			m[colName] = *val
		}

		outputRecords = append(outputRecords, m)

	}

	return map[string][]map[string]interface{}{
		"rows": outputRecords,
	}, nil

}

type updateOp struct {
	set    map[string]interface{}
	wheres map[string]interface{}
	aid    string
}

func updateOpValidator(input map[string]interface{}) (operationDoerFunc, error) {

	var err error
	op := new(updateOp)

	expect := map[string]bool{
		"type":  true,
		"where": true,
		"set":   true,
	}

	for k := range input {
		_, expected := expect[k]
		if !expected {
			return nil, fmt.Errorf("unexpected parameter on 'update' operation: %s", k)
		}
	}

	op.wheres, err = wheres(input, true)
	if err != nil {
		return nil, err
	}

	x, ok := input["set"]
	if !ok {
		return nil, fmt.Errorf("required parameter 'set' missing")
	}

	y, ok := x.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("required parameter 'set' must be a JSON object")
	}

	if len(y) == 0 {
		return nil, fmt.Errorf("'set' parameter must specify at least one column to update on the database")
	}

	for k, v := range y {
		switch v.(type) {
		case map[string]interface{}:
			return nil, fmt.Errorf("'set' parameter has non-primitive parameter '%s' (it's a JSON object)", k)
		case []interface{}:
			return nil, fmt.Errorf("'set' parameter has non-primitive parameter '%s' (it's an array)", k)
		default:
		}
	}

	op.set = y

	return op.do, nil

}

func (op *updateOp) do(tx *sql.Tx, table string) (interface{}, error) {

	var changes, obscuredChanges []string

	for k, v := range op.set {
		key := `"` + strings.ReplaceAll(fmt.Sprintf("%v", k), `"`, `""`) + `"`
		val := "'" + strings.ReplaceAll(fmt.Sprintf("%v", v), "'", "''") + "'"
		if v == nil {
			val = "NULL"
		}
		changes = append(changes, fmt.Sprintf("%s=%s", key, val))
		obscuredChanges = append(obscuredChanges, fmt.Sprintf("%s='****'", key))
	}

	sets := strings.Join(changes, ", ")
	// obscuredSets := strings.Join(obscuredChanges, ", ")
	wheres := wheresString(op.wheres)

	// obscuredQuery := fmt.Sprintf(`UPDATE %s SET %s WHERE %s`, table, obscuredSets, wheres)
	query := fmt.Sprintf(`UPDATE %s SET %s WHERE %s`, table, sets, wheres)

	result, err := tx.Exec(query)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to update: %w", err)
	}

	k, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to lookup rows affected: %w", err)
	}

	return map[string]int64{
		"rowsAffected": k,
	}, nil

}

type Input struct {
	Conn        string                   `json:"conn"`
	Debug       bool                     `json:"debug"`
	Table       string                   `json:"table"`
	Transaction []map[string]interface{} `json:"transaction"`
}

func validateInput(input *Input) ([]operationDoerFunc, error) {

	log.Println("Validating input data...")

	if input.Conn == "" {
		return nil, errors.New("missing input parameter: 'conn'")
	}

	if input.Table == "" {
		return nil, errors.New("missing input parameter: 'table'")
	}

	input.Table = `"` + strings.ReplaceAll(input.Table, `"`, `""`) + `"`

	if input.Transaction == nil {
		return nil, errors.New("missing input parameter: 'transaction'")
	}

	if len(input.Transaction) == 0 {
		return nil, errors.New("input parameter 'transaction' is empty")
	}

	var doers []operationDoerFunc

	for i := range input.Transaction {

		op := input.Transaction[i]
		if op == nil {
			return nil, fmt.Errorf("input parameter 'transaction' element %d is null", i)
		}

		x, ok := op["type"]
		if !ok {
			return nil, fmt.Errorf("input parameter 'transaction' element %d is missing required parameter 'type'", i)
		}

		typ, ok := x.(string)
		if !ok {
			return nil, fmt.Errorf("input parameter 'transaction' element %d has bad 'type' parameter: must be a string", i)
		}

		validator, exists := operations[typ]
		if !exists {
			var types []string
			for k := range operations {
				types = append(types, k)
			}
			sort.Strings(types)
			return nil, fmt.Errorf("input parameter 'transaction' element %d has bad 'type' parameter: got '%s' but must be one of %v", i, typ, types)
		}

		doer, err := validator(op)
		if err != nil {
			return nil, fmt.Errorf("input parameter 'transaction' element %d failed to validate: %w", i, err)
		}

		doers = append(doers, doer)

	}

	return doers, nil

}

func begin(input *Input, aid string) (*sql.Tx, error) {

	direktivapps.LogDouble(aid, "Connecting to postgres database...")

	db, err := sql.Open("postgres", input.Conn)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	direktivapps.LogDouble(aid, "Beginning database transaction...")

	tx, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	return tx, nil

}

var code = "com.simplepostgres.error"

func SimplePostgresClient(w http.ResponseWriter, r *http.Request) {
	input := new(Input)
	output := make([]interface{}, 0)

	aid, err := direktivapps.Unmarshal(input, r)
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	steps, err := validateInput(input)
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	tx, err := begin(input, aid)
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}
	defer tx.Rollback()

	for i, step := range steps {
		out, err := step(tx, input.Table)
		if err != nil {
			direktivapps.RespondWithError(w, code, fmt.Sprintf("transaction step %d failed: %v", i, err))
			return
		}
		output = append(output, out)
	}

	data, err := json.Marshal(output)
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	direktivapps.LogDouble(aid, "Committing transaction to database...")
	err = tx.Commit()
	if err != nil {
		direktivapps.RespondWithError(w, code, err.Error())
		return
	}

	direktivapps.LogDouble(aid, "Transaction complete!")
	direktivapps.Respond(w, data)
}

func main() {
	direktivapps.StartServer(SimplePostgresClient)
}
