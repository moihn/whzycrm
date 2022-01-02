package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/bdwilliams/go-jsonify/jsonify"
	_ "github.com/godror/godror"
	"github.com/integrii/flaggy"
	"github.com/joho/sqltocsv"
	"github.com/sirupsen/logrus"
)

func describeResult(rows *sql.Rows) map[string]*sql.ColumnType {
	columnTypes, err := rows.ColumnTypes()
	if err != nil {
		log.Fatalf("Failed to get column types for result set : %v", err)
	}
	columnTypesMap := map[string]*sql.ColumnType{}
	for i, columnType := range columnTypes {
		columnTypesMap[columnType.Name()] = columnType
		logrus.Infof("Column %v (%v) has type %v", i, columnType.Name(), columnType.ScanType().Name())
	}
	return columnTypesMap
}

func main() {
	var queryString, outputFormat, connectionString string
	logLevelName := "warn"
	flaggy.SetName("oraquery")
	flaggy.SetDescription("Database SQL query tool for Oracle")
	flaggy.DefaultParser.ShowHelpOnUnexpected = false
	flaggy.String(&connectionString, "c", "conn-string", "Connection string.")
	flaggy.String(&queryString, "q", "query", "Query string. Get from stdin when omitted.")
	flaggy.String(&outputFormat, "f", "format", "Format of output. CSV by default")
	flaggy.String(&logLevelName, "", "log-level", "Log level")
	flaggy.Parse()

	// Only log the warning severity or above.
	if logLevel, err := logrus.ParseLevel(logLevelName); err != nil {
		logrus.Fatal("Fail to parse logging level: ", logLevel)
	} else {
		logrus.SetLevel(logLevel)
		logrus.SetOutput(os.Stderr)
	}

	if queryString == "" {
		queryStringBytes, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			logrus.Fatal("Failed to get query from stdin: ", err)
		} else {
			queryString = string(queryStringBytes)
		}
	}

	if outputFormat == "" {
		outputFormat = "CSV"
	}

	db, err := sql.Open("godror", connectionString)
	if err != nil {
		logrus.Fatal(err)
	}
	rows, err := db.Query(queryString)
	if err != nil {
		log.Fatal("Failed to execute query: ", err)
	}
	defer rows.Close()

	describeResult(rows)

	switch strings.ToUpper(outputFormat) {
	case "CSV":
		converter := sqltocsv.New(rows)
		converter.TimeFormat = time.RFC3339
		if err := converter.Write(os.Stdout); err != nil {
			log.Fatal("Failed to write to stdout: ", err)
		}
	case "JSON":
		fmt.Println(jsonify.Jsonify(rows))
	}
}
