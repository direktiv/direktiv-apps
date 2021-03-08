package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"sort"
	"strings"

	_ "github.com/lib/pq"
)

type operationDoerFunc func(tx *sql.Tx, table string) (interface{}, error)
type operationValidatorFunc func(map[string]interface{}) (operationDoerFunc, error)

var operations = map[string]operationValidatorFunc{
	"delete": deleteOpValidator,
	"insert": insertOpValidator,
	"select": selectOpValidator,
	"update": updateOpValidator,
}

func wheres(input map[string]interface{}) (map[string]interface{}, error) {

	x, ok := input["where"]
	if !ok {
		return nil, fmt.Errorf("required parameter 'where' missing")
	}

	wheres, ok := x.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("'where' parameter must be a JSON object")
	}

	if len(wheres) == 0 {
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
		key := "'" + strings.ReplaceAll(k, "'", "''") + "'"
		val := "'" + strings.ReplaceAll(fmt.Sprintf("%v", v), "'", "''") + "'"
		conditions = append(conditions, fmt.Sprintf("%s=%s", key, val))
	}

	return strings.Join(conditions, " AND ")

}

type deleteOp struct {
	wheres map[string]interface{}
}

func deleteOpValidator(input map[string]interface{}) (operationDoerFunc, error) {

	var err error
	op := new(deleteOp)

	op.wheres, err = wheres(input)
	if err != nil {
		return nil, err
	}

	return op.do, nil

}

func (op *deleteOp) do(tx *sql.Tx, table string) (interface{}, error) {

	aerr.ErrorCode = "sql.db.delete"

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
}

func insertOpValidator(input map[string]interface{}) (operationDoerFunc, error) {

	op := new(insertOp)

	var records []map[string]interface{}

	x, ok := input["set"]
	if !ok {
		return nil, fmt.Errorf("required parameter 'set' missing")
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
		return nil, fmt.Errorf("'set' parameter must be a JSON object or an array of JSON objects")
	}

	if len(records) == 0 {
		return nil, fmt.Errorf("'set' parameter must specify at least one record to add to the database")
	}

	for i, rec := range records {
		for k, v := range rec {
			switch v.(type) {
			case map[string]interface{}:
				return nil, fmt.Errorf("'set' parameter has non-primitive parameter '%s' in record %d (it's a JSON object)", k, i)
			case []interface{}:
				return nil, fmt.Errorf("'set' parameter has non-primitive parameter '%s' in record %d (it's an array)", k, i)
			default:
			}
		}
	}

	op.records = records

	return op.do, nil

}

func (op *insertOp) do(tx *sql.Tx, table string) (interface{}, error) {

	aerr.ErrorCode = "sql.db.insert"

	var rowIds []int64

	for i, record := range op.records {

		var keys, vals []string

		for k, v := range record {
			key := "'" + strings.ReplaceAll(k, "'", "''") + "'"
			val := "'" + strings.ReplaceAll(fmt.Sprintf("%v", v), "'", "''") + "'"
			keys = append(keys, key)
			vals = append(vals, val)
		}

		ks := strings.Join(keys, ", ")
		vs := strings.Join(vals, ", ")
		query := fmt.Sprintf(`INSERT INTO %s (%s) VALUES (%s)`, table, ks, vs)
		result, err := tx.Exec(query)
		if err != nil {
			return nil, fmt.Errorf("failed to insert record %d: %w", i, err)
		}

		id, err := result.LastInsertId()
		if err != nil {
			return nil, fmt.Errorf("failed to lookup inserted row id: %w", err)
		}

		rowIds = append(rowIds, id)

	}

	return map[string][]int64{
		"rowIds": rowIds,
	}, nil

}

type selectOp struct {
	fields []string
	wheres map[string]interface{}
}

func selectOpValidator(input map[string]interface{}) (operationDoerFunc, error) {

	var err error
	op := new(selectOp)

	op.wheres, err = wheres(input)
	if err != nil {
		return nil, err
	}

	var fields []string

	x, ok := input["fields"]
	if !ok {
		return nil, fmt.Errorf("required parameter 'fields' missing")
	}

	y, ok := x.([]interface{})
	if !ok {
		return nil, fmt.Errorf("'fields' parameter must be a JSON array of strings")
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

	aerr.ErrorCode = "sql.db.select"

	var fields []string
	for _, field := range op.fields {
		fieldName := "'" + strings.ReplaceAll(field, "'", "''") + "'"
		fields = append(fields, fieldName)
	}

	wheres := wheresString(op.wheres)
	query := fmt.Sprintf(`SELECT %s FROM %s WHERE %s`, fields, table, wheres)
	rows, err := tx.Query(query)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to select: %v", err)
	}
	defer rows.Close()

	var outputRecords = make([]map[string]interface{}, 0)

	for rows.Next() {

		columns := make([]interface{}, len(op.fields))
		columnPointers := make([]interface{}, len(op.fields))
		for i, _ := range columns {
			columnPointers[i] = &columns[i]
		}

		if err := rows.Scan(columnPointers...); err != nil {
			return nil, fmt.Errorf("failed to load results from query: %v", err)
		}

		m := make(map[string]interface{})
		for i, colName := range op.fields {
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
}

func updateOpValidator(input map[string]interface{}) (operationDoerFunc, error) {

	var err error
	op := new(updateOp)

	op.wheres, err = wheres(input)
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

	aerr.ErrorCode = "sql.db.update"

	var changes []string

	for k, v := range op.set {
		key := "'" + strings.ReplaceAll(k, "'", "''") + "'"
		val := "'" + strings.ReplaceAll(fmt.Sprintf("%v", v), "'", "''") + "'"
		changes = append(changes, fmt.Sprintf("%s=%s", key, val))
	}

	sets := strings.Join(changes, ", ")

	wheres := wheresString(op.wheres)
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

type ActionError struct {
	ErrorCode    string `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
}

func writeError(g ActionError) {
	log.Printf("ERROR:\n  CODE: %v\n  TEXT: %v\n", g.ErrorCode, g.ErrorMessage)
	b, _ := json.Marshal(g)
	ioutil.WriteFile("/direktiv-data/error.json", b, 0755)
}

type Input struct {
	Conn        string                   `json:"conn"`
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

func getInput() (*Input, error) {

	log.Println("Reading input data...")
	aerr.ErrorCode = "error.input"
	input := new(Input)

	data, err := ioutil.ReadFile("/direktiv-data/data.in")
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, input)
	if err != nil {
		return nil, err
	}

	return input, nil

}

func saveOutput(output []interface{}) error {

	data, err := json.Marshal(output)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile("/direktiv-data/data.out", data, 0755)
	if err != nil {
		return err
	}

	return nil

}

func begin(input *Input) (*sql.Tx, error) {

	log.Println("Connecting to postgres database...")
	aerr.ErrorCode = "error.conn"

	db, err := sql.Open("postgres", input.Conn)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	log.Println("Beginning database transaction...")
	aerr.ErrorCode = "error.tx"

	tx, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	return tx, nil

}

var aerr ActionError

func main() {

	input := new(Input)
	output := make([]interface{}, 0)

	defer func() {

		r := recover()
		if r != nil {
			aerr.ErrorCode = "error.panic"
			aerr.ErrorMessage = fmt.Sprintf("%v", r)
		}

		if aerr.ErrorMessage != "" {
			writeError(aerr)
		}

	}()

	input, err := getInput()
	if err != nil {
		aerr.ErrorMessage = err.Error()
		return
	}

	steps, err := validateInput(input)
	if err != nil {
		aerr.ErrorMessage = err.Error()
		return
	}

	tx, err := begin(input)
	if err != nil {
		aerr.ErrorMessage = err.Error()
		return
	}
	defer tx.Rollback()

	aerr.ErrorCode = "error.db"

	for i, step := range steps {
		out, err := step(tx, input.Table)
		if err != nil {
			aerr.ErrorMessage = fmt.Sprintf("transaction step %d failed: %v", i, err)
			return
		}
		output = append(output, out)
	}

	err = saveOutput(output)
	if err != nil {
		aerr.ErrorCode = "error.output"
		aerr.ErrorMessage = err.Error()
		return
	}

	log.Println("Committing transaction to database...")
	err = tx.Commit()
	if err != nil {
		aerr.ErrorCode = "error.db.commit"
		aerr.ErrorMessage = err.Error()
		return
	}

	log.Println("Transaction complete!")
	aerr.ErrorCode = ""
	aerr.ErrorMessage = ""

}
