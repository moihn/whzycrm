package table

import (
	"database/sql"

	"github.com/sirupsen/logrus"
)

type ClientProductItemRow struct {
	ClientProductItemId int
	ClientProductId     int
	VendorProductId     int
}

func ClientProductItemGetByFk1(tx *sql.Tx, clientProductId int) []ClientProductItemRow {
	sqlQuery := `
	select CLIENT_PRODUCT_ITEM_ID, CLIENT_PRODUCT_ID, VENDOR_PRODUCT_ID
	from CLIENT_PRODUCT_ITEM
	where client_product_id = :clientProductId`
	rows, err := tx.Query(sqlQuery,
		sql.Named("clientProductId", &clientProductId),
	)
	if err != nil {
		logrus.Fatalf("failed to run CLIENT_PRODUCT_ITEM table by Fk1 query: %v", err)
	}
	defer rows.Close()

	var result []ClientProductItemRow
	for rows.Next() {
		var clientProductItemId, clientProductId, vendorProductId int
		err = rows.Scan(&clientProductItemId, &clientProductId, &vendorProductId)
		if err != nil {
			logrus.Fatalf("failed to extract CLIENT_PRODUCT row data: %v", err)
		}
		result = append(result, ClientProductItemRow{
			ClientProductId:     clientProductId,
			ClientProductItemId: clientProductItemId,
			VendorProductId:     vendorProductId,
		})
	}
	return result
}
