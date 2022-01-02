package table

import (
	"database/sql"
)

type ClientProductRow struct {
	ClientProductId int
	ClientId        int
	Reference       string
	Description     string
	Narrative       string
	Barcode         string
}

func ClientProductGetByFk1(tx *sql.Tx, clientId int) ([]ClientProductRow, error) {
	sqlQuery := `
		select CLIENT_PRODUCT_ID, REFERENCE, DESCRIPTION, NARRATIVE, BARCODE
		from CLIENT_PRODUCT
		where client_id = :clientId`
	rows, err := tx.Query(sqlQuery,
		sql.Named("clientId", &clientId),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []ClientProductRow
	for rows.Next() {
		var clientProductId int
		var reference, description, narrative, barcode string
		err = rows.Scan(&clientProductId, &reference, &description, &narrative, &barcode)
		if err != nil {
			return nil, err
		}
		result = append(result, ClientProductRow{
			ClientProductId: clientProductId,
			ClientId:        clientId,
			Reference:       reference,
			Description:     description,
			Narrative:       narrative,
			Barcode:         barcode,
		})
	}
	return result, nil
}

func ClientProductGetByPk(tx *sql.Tx, clientProductId int) (*ClientProductRow, error) {
	sqlQuery := `
		select CLIENT_ID, REFERENCE, DESCRIPTION, NARRATIVE, BARCODE
		from CLIENT_PRODUCT
		where CLIENT_PRODUCT_ID = :clientProductId`
	rows, err := tx.Query(sqlQuery,
		sql.Named("clientProductId", &clientProductId),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result *ClientProductRow
	if rows.Next() {
		var clientId int
		var reference, description, narrative, barcode string
		err = rows.Scan(&clientId, &reference, &description, &narrative, &barcode)
		if err != nil {
			return nil, err
		}
		result = &ClientProductRow{
			ClientProductId: clientProductId,
			ClientId:        clientId,
			Reference:       reference,
			Description:     description,
			Narrative:       narrative,
			Barcode:         barcode,
		}
	}
	return result, nil
}
