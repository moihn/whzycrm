package table

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

type VendorProductMoqRow struct {
	VendorProductId int
	Quantity        int
	StartDate       time.Time
}

func VendorProductMoqInsert(tx *sql.Tx, vendorProductId int, quantity int, startDate time.Time) error {
	sqlQuery := `
		INSERT INTO VENDOR_PRODUCT_MOQ
		(VENDOR_PRODUCT_ID, QUANTITY, START_DATE)
		VALUES
		(:vendorProductId, :quantity, trunc(:startDate))
	`
	if result, err := tx.Exec(sqlQuery,
		sql.Named("vendorProductId", vendorProductId),
		sql.Named("startDate", startDate),
		sql.Named("quantity", quantity)); err != nil {
		return fmt.Errorf("failed to execute query [%v]: %v", sqlQuery, err)
	} else {
		nRows, err := result.RowsAffected()
		if err != nil {
			logrus.Fatalf("failed to get number of inserted row for query [%v]: %v", sqlQuery, err)
		}
		if nRows == 0 {
			logrus.Fatalf("failed to insert VENDOR_PRODUCT_MOQ row [%v]: vendorProductId=%v, startDate=%v, quantity=%v", sqlQuery,
				vendorProductId, startDate, quantity)
		}
	}
	return nil
}

func VendorProductMoqPopulateByVendProductId(tx *sql.Tx, vendorProductId int) []VendorProductMoqRow {
	sqlQuery := `
		SELECT
			QUANTITY, START_DATE
		FROM VENDOR_PRODUCT_MOQ
		WHERE
			VENDOR_PRODUCT_ID = :vendorProductId
		ORDER BY START_DATE DESC
	`
	rows, err := tx.Query(sqlQuery, sql.Named("vendorProductId", vendorProductId))
	if err != nil {
		logrus.Fatalf("failed to query VENDOR_PRODUCT_PRICE table. sql=%v, error=%v", sqlQuery, err)
	}
	defer rows.Close()

	result := []VendorProductMoqRow{}
	for rows.Next() {
		var quantity int
		var startDate time.Time
		err = rows.Scan(&quantity, &startDate)
		if err != nil {
			logrus.Fatalf("failed to extract values for query: %v. query=%v", err, sqlQuery)
		}
		result = append(result, VendorProductMoqRow{
			VendorProductId: vendorProductId,
			StartDate:       startDate,
			Quantity:        quantity,
		})
	}
	return result
}
