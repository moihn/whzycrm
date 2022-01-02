package table

import (
	"database/sql"

	"github.com/sirupsen/logrus"
)

type ClientRow struct {
	Id        int
	Name      string
	CountryId int
}

func ClientPopulateByPk(tx *sql.Tx) []ClientRow {
	sqlQuery := `
		select CLIENT_ID, NAME, COUNTRY_ID
		from CLIENT`
	rows, err := tx.Query(sqlQuery)
	if err != nil {
		logrus.Fatalf("failed to run client table query: %v", err)
	}
	defer rows.Close()

	var result []ClientRow
	for rows.Next() {
		var clientId int
		var name string
		var countryId int
		err = rows.Scan(&clientId, &name, &countryId)
		if err != nil {
			logrus.Fatalf("failed to extract client row data: %v", err)
		}
		result = append(result, ClientRow{
			Id:        clientId,
			Name:      name,
			CountryId: countryId,
		})
	}
	return result
}
