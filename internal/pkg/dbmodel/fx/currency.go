package fx

import (
	"database/sql"
	"fmt"

	"github.com/moihn/whzycrmgo/internal/pkg/util"
)

type Currency struct {
	Id        int
	IsoSymbol string
}

func GetCurrencyById(tx *sql.Tx, ccyId int) (*Currency, *util.Status) {
	sqlQuery := `
		select ISO_SYMBOL
		from CURRENCY
		where CURRENCY_ID = :currencyId
	`
	rows, err := tx.Query(sqlQuery,
		sql.Named("currencyId", ccyId))
	if err != nil {
		status := util.NewInternalServiceErrorStatus(
			fmt.Sprintf("failed to run currency look up query: %v. Error: %v", sqlQuery, err),
			"GET_CURRENCY_BY_ID_0")
		return nil, &status
	}
	defer rows.Close()
	if rows.Next() {
		var iso string
		err := rows.Scan(&iso)
		if err != nil {
			status := util.NewInternalServiceErrorStatus(
				fmt.Sprintf("failed to extract result of query: %v. Error: %v", sqlQuery, err),
				"GET_CURRENCY_BY_ID_1")
			return nil, &status
		} else {
			ccy := Currency{
				Id:        ccyId,
				IsoSymbol: iso,
			}
			return &ccy, nil
		}
	}
	return nil, nil
}

func GetCurrencyByIsoSymbol(tx *sql.Tx, isoSymbol string) (*Currency, *util.Status) {
	sqlQuery := `
		select CURRENCY_ID
		from CURRENCY
		where ISO_SYMBOL = :isoSymbol
	`
	rows, err := tx.Query(sqlQuery,
		sql.Named("isoSymbol", isoSymbol))
	if err != nil {
		status := util.NewInternalServiceErrorStatus(
			fmt.Sprintf("failed to run currency look up query: %v. Error: %v", sqlQuery, err),
			"GET_CURRENCY_BY_ID_0")
		return nil, &status
	}
	defer rows.Close()
	if rows.Next() {
		var ccyId int
		err := rows.Scan(&ccyId)
		if err != nil {
			status := util.NewInternalServiceErrorStatus(
				fmt.Sprintf("failed to extract result of query: %v. Error: %v", sqlQuery, err),
				"GET_CURRENCY_BY_ID_1")
			return nil, &status
		} else {
			ccy := Currency{
				Id:        ccyId,
				IsoSymbol: isoSymbol,
			}
			return &ccy, nil
		}
	}
	return nil, nil
}
