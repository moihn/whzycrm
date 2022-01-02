package table

import "database/sql"

type ClientQuotationItem struct {
	QuotationId     int
	VendorProductId int
	Price           float32
	Narrative       string
	Moq             int
}

func ClientQuotationItemGetByPk(
	tx *sql.Tx,
	quotationId int,
	vendorProductId int,
) (*ClientQuotationItem, error) {
	sqlQuery := `
		select
			QUOTATION_ID,
			VENDOR_PRODUCT_ID,
			PRICE,
			NARRATIVE,
			MOQ
		from CLIENT_QUOTATION_ITEM
		where
		      QUOTATION_ID = :quotationId
		  and VENDOR_PRODUCT_ID = :vendorProductId
	`
	rows, err := tx.Query(sqlQuery,
		sql.Named("quotationId", quotationId),
		sql.Named("vendorProductId", vendorProductId),
	)
	if err != nil {
		return nil, err
	}
	if rows.Next() {
		result := ClientQuotationItem{}
		if err := rows.Scan(
			&result.QuotationId,
			&result.VendorProductId,
			&result.Price,
			&result.Moq,
		); err != nil {
			return nil, err
		}
		return &result, nil
	}
	return nil, nil
}

func ClientQuotationItemPopulateByFk1(tx *sql.Tx, quotationId int) ([]*ClientQuotationItem, error) {
	sqlQuery := `
		select
			QUOTATION_ID,
			VENDOR_PRODUCT_ID,
			PRICE,
			NARRATIVE,
			MOQ
		from CLIENT_QUOTATION_ITEM
		where
		      VENDOR_PRODUCT_ID = :vendorProductId
	`
	rows, err := tx.Query(sqlQuery,
		sql.Named("quotationId", quotationId),
	)
	if err != nil {
		return nil, err
	}

	results := []*ClientQuotationItem{}
	for rows.Next() {
		result := ClientQuotationItem{}
		if err := rows.Scan(
			&result.QuotationId,
			&result.VendorProductId,
			&result.Price,
			&result.Moq,
		); err != nil {
			return nil, err
		}
		results = append(results, &result)
	}
	return results, nil
}

func ClientQuotationItemPopulateByFk2(tx *sql.Tx, vendorProductId int) ([]*ClientQuotationItem, error) {
	sqlQuery := `
		select
			QUOTATION_ID,
			VENDOR_PRODUCT_ID,
			PRICE,
			NARRATIVE,
			MOQ
		from CLIENT_QUOTATION_ITEM
		where
		      QUOTATION_ID = :quotationId
	`
	rows, err := tx.Query(sqlQuery,
		sql.Named("vendorProductId", vendorProductId),
	)
	if err != nil {
		return nil, err
	}

	results := []*ClientQuotationItem{}
	for rows.Next() {
		result := ClientQuotationItem{}
		if err := rows.Scan(
			&result.QuotationId,
			&result.VendorProductId,
			&result.Price,
			&result.Moq,
		); err != nil {
			return nil, err
		}
		results = append(results, &result)
	}
	return results, nil
}
