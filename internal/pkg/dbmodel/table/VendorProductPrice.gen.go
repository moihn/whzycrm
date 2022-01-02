package table

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

type VendorProductPriceRow struct {
	PriceId         int
	VendorProductId int
	StartDate       time.Time
	Price           float32
	CurrencyId      int
	PriceTypeId     int
}

func VendorProductPricePopulateByVendorProductId(tx *sql.Tx, vendorProductId int) []VendorProductPriceRow {
	sqlQuery := `
		SELECT
			PRICE_ID, START_DATE, PRICE, CURRENCY_ID, PRICE_TYPE_ID
		FROM VENDOR_PRODUCT_PRICE
		WHERE VENDOR_PRODUCT_ID = :vendorProductId
		ORDER BY START_DATE DESC
	`
	rows, err := tx.Query(sqlQuery, sql.Named("vendorProductId", vendorProductId))
	if err != nil {
		logrus.Fatalf("failed to query VENDOR_PRODUCT_PRICE table. sql=%v, error=%v", sqlQuery, err)
	}
	defer rows.Close()

	result := []VendorProductPriceRow{}
	for rows.Next() {
		var priceId, currencyId, priceTypeId int
		var startDate time.Time
		var price float32
		err = rows.Scan(&priceId, &startDate, &price, &currencyId, &priceTypeId)
		if err != nil {
			logrus.Fatalf("failed to extract values for query: %v. query=%v", err, sqlQuery)
		}
		result = append(result, VendorProductPriceRow{
			PriceId:         priceId,
			VendorProductId: vendorProductId,
			StartDate:       startDate,
			Price:           price,
			CurrencyId:      currencyId,
			PriceTypeId:     priceTypeId,
		})
	}
	return result
}

func VendorProductPriceInsert(
	tx *sql.Tx,
	vendorProductId int,
	startDate time.Time,
	price float32,
	currencyId int,
	priceTypeId int) (int, error) {
	sqlQuery := `
		INSERT INTO VENDOR_PRODUCT_PRICE
		(VENDOR_PRODUCT_ID, START_DATE, PRICE, CURRENCY_ID, PRICE_TYPE_ID)
		VALUES
		(:vendorProductId, trunc(:startDate), :price, :ccyId, :priceTypeId)
		RETURNING PRICE_ID INTO :priceId
	`
	var priceId int
	if result, err := tx.Exec(sqlQuery,
		sql.Named("vendorProductId", vendorProductId),
		sql.Named("startDate", startDate),
		sql.Named("price", price),
		sql.Named("ccyId", currencyId),
		sql.Named("priceTypeId", priceTypeId),
		sql.Named("priceId", sql.Out{Dest: &priceId})); err != nil {
		return 0, fmt.Errorf("failed to execute query [%v]: %v", sqlQuery, err)
	} else {
		nRows, err := result.RowsAffected()
		if err != nil {
			logrus.Fatalf("failed to get number of inserted row for query [%v]: %v", sqlQuery, err)
		}
		if nRows == 0 {
			logrus.Fatalf("failed to insert VENDOR_PRODUCT_PRICE row [%v]: vendorProductId=%v, startDate=%v, price=%v", sqlQuery,
				vendorProductId, startDate, price)
		}
	}
	return priceId, nil
}
