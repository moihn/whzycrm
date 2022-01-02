package quotation

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/moihn/whzycrmgo/internal/pkg/dbmodel/fx"
	"github.com/moihn/whzycrmgo/internal/pkg/dbmodel/supplier/product"
	"github.com/moihn/whzycrmgo/internal/pkg/util"
)

type ItemInQuotation struct {
	VendorProductId int
	PriceAmount     float32
	MOQ             *int
	Narrative       string
}

type Quotation struct {
	QuotationId int
	Date        time.Time
	ClientId    int
	CurrencyId  int
	Sent        bool
	Items       []ItemInQuotation
}

func CreateClientQuotation(tx *sql.Tx, quotation Quotation) *util.Status {
	// insert parent row
	var quoteId int
	sqlQuery := `
		insert into CLIENT_QUOTATION
		(CLIENT_ID, CURRENCY_ID, UPDATED_DATE, SENT)
		values
		(:clientId, :ccyId, TRUNC(SYSDATE), 'F')
		returning QUOTATION_ID into :quoteId
	`
	if result, err := tx.Exec(sqlQuery,
		sql.Named("clientId", quotation.ClientId),
		sql.Named("ccyId", quotation.CurrencyId),
		sql.Named("quoteId", sql.Out{Dest: &quoteId})); err != nil {
		status := util.NewInternalServiceErrorStatus(
			fmt.Sprintf("Failed to execute query [%v]: %v",
				sqlQuery, err),
			"CreateClientQuotation_0")
		return &status
	} else {
		nRows, err := result.RowsAffected()
		if err != nil {
			status := util.NewInternalServiceErrorStatus(
				fmt.Sprintf("failed to get number of inserted row for query [%v]: %v", sqlQuery, err),
				"CreateClientQuotation_4")
			return &status
		}
		if nRows == 0 {
			status := util.NewInternalServiceErrorStatus(
				fmt.Sprintf("failed to insert client quotation [%v]: client_id=%v", sqlQuery, quotation.ClientId),
				"CreateClientQuotation_5")
			return &status
		}
	}

	for _, p := range quotation.Items {
		productInfo, status := product.GetVendorProductById(tx, p.VendorProductId)
		if status != nil {
			return status
		}
		if len(productInfo.PriceList) == 0 {
			status := util.NewNotFoundStatus(
				fmt.Sprintf("no valid price is found for product %v of vendor with ID %v",
					productInfo.Reference,
					productInfo.VendorId))
			return &status
		}

		moq := sql.NullInt64{}
		if len(productInfo.MoqList) > 0 {
			moq.Int64 = int64(productInfo.MoqList[0].Quantity)
			moq.Valid = true
		}
		fxRate, status := fx.GetFxRate(tx, productInfo.PriceList[0].CurrencyId, quotation.CurrencyId)
		if status != nil {
			return status
		}
		sqlQuery = `
			insert into CLIENT_QUOTATION_ITEM
			(QUOTATION_ID, VENDOR_PRODUCT_ID, PRICE, MOQ)
			select :quoteId, :vendorProductId, :rawPrice * (1 + pt.invoice_rate) * :fxRate, :moq
			from PRICE_TYPE pt
			where pt.PRICE_TYPE_ID = :rawPriceTypeId
		`
		if result, err := tx.Exec(sqlQuery,
			sql.Named("quoteId", quoteId),
			sql.Named("vendorProductId", p.VendorProductId),
			sql.Named("fxRate", fxRate),
			sql.Named("rawPrice", productInfo.PriceList[0].Amount),
			sql.Named("moq", moq),
			sql.Named("rawPriceTypeId", productInfo.PriceList[0].TypeId)); err != nil {
			status := util.NewInternalServiceErrorStatus(
				fmt.Sprintf("Failed to execute query [%v]: %v", sqlQuery, err),
				"CreateClientQuotation_1")
			return &status
		} else {
			nRows, err := result.RowsAffected()
			if err != nil {
				status := util.NewInternalServiceErrorStatus(
					fmt.Sprintf("failed to get number of inserted row for query [%v]: %v", sqlQuery, err),
					"CreateClientQuotation_2")
				return &status
			}
			if nRows == 0 {
				status := util.NewInternalServiceErrorStatus(
					fmt.Sprintf("failed to insert query [%v]: price_type_id=%v", sqlQuery, productInfo.PriceList[0].TypeId),
					"CreateClientQuotation_3")
				return &status
			}
		}
	}

	return nil
}

func FindClientQuotations(tx *sql.Tx, clientId *int, startDate time.Time) ([]Quotation, *util.Status) {
	var quotes []Quotation

	sqlQuery := `
		select QUOTATION_ID, CLIENT_ID, CURRENCY_ID, UPDATED_DATE, SENT
		from CLIENT_QUOTATION`
	sqlArgs := []interface{}{}
	if clientId != nil || !startDate.IsZero() {
		sqlQuery += ` where  1 = 1`
		if clientId != nil {
			sqlQuery += ` and CLIENT_ID = :clientId`
			sqlArgs = append(sqlArgs, sql.Named("clientId", *clientId))
		}
		if !startDate.IsZero() {
			sqlQuery += ` and UPDATED_DATE >= :startDate`
			sqlArgs = append(sqlArgs, sql.Named("startDate", startDate))
		}
	}

	rows, err := tx.Query(sqlQuery, sqlArgs...)
	if err != nil {
		status := util.NewInternalServiceErrorStatus(
			fmt.Sprintf("Failed to execute query [%v]: %v", sqlQuery, err),
			"GetClientQuotation_1")
		return nil, &status
	}
	defer rows.Close()

	for rows.Next() {
		quote := Quotation{}
		sent := "F"
		err = rows.Scan(&quote.QuotationId, &quote.ClientId, &quote.CurrencyId, &quote.Date, &sent)
		if err != nil {
			status := util.NewInternalServiceErrorStatus(
				fmt.Sprintf("Failed to extract result of query [%v]: %v", sqlQuery, err),
				"GetClientQuotation_2")
			return nil, &status
		}
		if sent == "T" {
			quote.Sent = true
		}
		quotes = append(quotes, quote)
	}

	// for each quotation, get items
	sqlQuery = `
		select VENDOR_PRODUCT_ID, PRICE, MOQ, NARRATIVE
		from CLIENT_QUOTATION_ITEM
		where QUOTATION_ID = :quoteId
	`
	for idx, _ := range quotes {
		quote := &quotes[idx]
		itemRows, err := tx.Query(sqlQuery, sql.Named("quoteId", quote.QuotationId))
		if err != nil {
			status := util.NewInternalServiceErrorStatus(
				fmt.Sprintf("Failed to execute query [%v]: %v", sqlQuery, err),
				"GetClientQuotation_3")
			return nil, &status
		}
		defer itemRows.Close()
		for itemRows.Next() {
			quoteItem := ItemInQuotation{}
			moq := sql.NullInt64{}
			err := itemRows.Scan(&quoteItem.VendorProductId, &quoteItem.PriceAmount, &moq, &quoteItem.Narrative)
			if err != nil {
				status := util.NewInternalServiceErrorStatus(
					fmt.Sprintf("Failed to extract result of query [%v]: %v", sqlQuery, err),
					"GetClientQuotation_4")
				return nil, &status
			}
			if moq.Valid {
				moq := int(moq.Int64)
				quoteItem.MOQ = &moq
			}
			quote.Items = append(quote.Items, quoteItem)
		}
	}

	return quotes, nil
}
