package fx

import (
	"database/sql"
	"fmt"

	"github.com/moihn/whzycrmgo/internal/pkg/util"
)

type FxRate struct {
	FromCurrency Currency
	ToCurrency   Currency
	MultRate     float32
}

func getFxRate(tx *sql.Tx, fromCcyId int, toCcyId int) (*FxRate, *util.Status) {
	sqlQuery := `
		select c1.ISO_SYMBOL, er.MULT_RATE, c2.ISO_SYMBOL
		from EXCHANGE_RATE er, CURRENCY c1, CURRENCY c2
		where er.FROM_CCY_ID = :from_ccy_id
		  and er.TO_CCY_ID = :to_ccy_id
		  and c1.CURRENCY_ID = er.FROM_CCY_ID
		  and c2.CURRENCY_ID = er.TO_CCY_ID
		order by START_DATE desc
	`
	rows, err := tx.Query(sqlQuery,
		sql.Named("from_ccy_id", fromCcyId),
		sql.Named("to_ccy_id", toCcyId),
	)
	if err != nil {
		status := util.NewInternalServiceErrorStatus(
			fmt.Sprintf("failed to run exchange rate query: %v, from_ccy_id = %v, to_ccy_id = %v. Error: %v",
				sqlQuery, fromCcyId, toCcyId, err),
			"GET_FX_RATE_0")
		return nil, &status
	}
	defer rows.Close()

	var multRate float32
	var fromCcyIso string
	var toCcyIso string
	if rows.Next() {
		err := rows.Scan(&fromCcyIso, &multRate, &toCcyIso)
		if err != nil {
			status := util.NewInternalServiceErrorStatus(
				fmt.Sprintf("failed to extract data from exchange rate query: %v, from_ccy_id = %v, to_ccy_id = %v. Error: %v",
					sqlQuery, fromCcyId, toCcyId, err),
				"GET_FX_RATE_1")
			return nil, &status
		}
	} else {
		return nil, nil
	}
	return &FxRate{
		FromCurrency: Currency{
			Id:        fromCcyId,
			IsoSymbol: fromCcyIso,
		},
		ToCurrency: Currency{
			Id:        toCcyId,
			IsoSymbol: toCcyIso,
		},
		MultRate: multRate,
	}, nil
}

func GetFxRate(tx *sql.Tx, fromCcyId int, toCcyId int) (*float32, *util.Status) {
	multRate, status := getFxRate(tx, fromCcyId, toCcyId)
	if status != nil {
		return nil, status
	}
	if multRate != nil {
		if multRate.MultRate == 0 {
			status := util.NewInternalServiceErrorStatus(
				fmt.Sprintf("fx rate from currency %v to %v is zero",
					multRate.FromCurrency.IsoSymbol, multRate.ToCurrency.IsoSymbol),
				"GET_FX_RATE_2")
			return nil, &status
		}
		return &multRate.MultRate, nil
	}
	divRate, err := getFxRate(tx, toCcyId, fromCcyId)
	if err != nil {
		return nil, err
	}
	if divRate != nil {
		if divRate.MultRate == 0 {
			status := util.NewInternalServiceErrorStatus(
				fmt.Sprintf("fx rate from currency %v to %v is zero",
					divRate.FromCurrency.IsoSymbol, divRate.ToCurrency.IsoSymbol),
				"GET_FX_RATE_3")
			return nil, &status
		}
		multRate := 1.0 / (divRate.MultRate)
		return &multRate, nil
	}
	status2 := util.NewNotFoundStatus(
		fmt.Sprintf("no fx rate is found between currencies with id %v and %v", fromCcyId, toCcyId))
	return nil, &status2
}
