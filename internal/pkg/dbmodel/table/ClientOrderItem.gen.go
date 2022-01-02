package table

import (
	"database/sql"
	"fmt"
	"time"
)

type ClientOrderItemRow struct {
	OrderItemId         int
	OrderId             int
	ClientProductId     int
	Quantity            int
	Price               float32
	CurrencyId          int
	AddedDate           time.Time
	AlternativeShipDate *time.Time
}

func ClientOrderItem_Insert(
	tx *sql.Tx,
	orderId int,
	clientProductId int,
	quantity int,
	price float32,
	currencyId int,
	addedDate time.Time,
	alternativeShipmentDate *time.Time,
) (int, error) {
	sqlQuery := `
		INSERT INTO CLIENT_ORDER_ITEM
		(ORDER_ID, CLIENT_PRODUCT_ID, QUANTITY, PRICE, CURRENCY_ID, ADDED_DATE, ALTERNATIVE_SHIP_DATE)
		VALUES
		(:orderId, :clientProductId, :quantity, :price, :currencyId, :addedDate, :alternativeShipmentDate)
		RETURNING ORDER_ITEM_ID into :orderItemId
	`
	var orderItemId int
	if result, err := tx.Exec(sqlQuery,
		sql.Named("orderId", orderId),
		sql.Named("clientProductId", clientProductId),
		sql.Named("quantity", quantity),
		sql.Named("price", price),
		sql.Named("currencyId", currencyId),
		sql.Named("addedDate", addedDate),
		sql.Named("alternativeShipmentDate", alternativeShipmentDate),
		sql.Named("orderItemId", sql.Out{Dest: &orderItemId})); err != nil {
		return 0, fmt.Errorf("failed to execute query [%v]: %v", sqlQuery, err)
	} else {
		nRows, err := result.RowsAffected()
		if err != nil {
			return 0, err
		}
		if nRows == 0 {
			return 0, fmt.Errorf("failed to insert CLIENT_ORDER_ITEM row [%v]: orderId=%v, clientProductId=%v", sqlQuery,
				orderId, clientProductId)
		}
	}
	return orderId, nil
}

func ClientOrderItem_GetByPk(tx *sql.Tx, orderItemId int) (*ClientOrderItemRow, error) {
	sqlQuery := `
		SELECT
			ORDER_ITEM_ID,
			ORDER_ID,
			CLIENT_PRODUCT_ID,
			QUANTITY,
			PRICE,
			CURRENCY_ID,
			ADDED_DATE,
			ALTERNATIVE_SHIP_DATE
		FROM
			CLIENT_ORDER_ITEM
		WHERE
			ORDER_ITEM_ID = :orderItemId`
	rows, err := tx.Query(sqlQuery,
		sql.Named("orderItemId", orderItemId))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result *ClientOrderItemRow
	if rows.Next() {
		err = rows.Scan(&result.OrderItemId, &result.OrderId, &result.ClientProductId, &result.Quantity, &result.Price, &result.CurrencyId,
			&result.AddedDate, &result.AlternativeShipDate)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func (row *ClientOrderItemRow) Update(tx *sql.Tx) error {
	sqlQuery := `
		UPDATE CLIENT_ORDER_ITEM
		SET
			ORDER_ID = :orderId,
			CLIENT_PRODUCT_ID = :clientProductId,
			QUANTITY = :quantity,
			PRICE = :price,
			CURRENCY_ID = :currencyId,
			ADDED_DATE = :addedDate,
			ALTERNATIVE_SHIP_DATE = :alternativeShipmentDate
		WHERE
			ORDER_ITEM_ID = :orderItemId`
	if result, err := tx.Exec(sqlQuery,
		sql.Named("orderId", row.OrderId),
		sql.Named("clientProductId", row.ClientProductId),
		sql.Named("quantity", row.Quantity),
		sql.Named("price", row.Price),
		sql.Named("currencyId", row.CurrencyId),
		sql.Named("addedDate", row.AddedDate),
		sql.Named("alternativeShipmentDate", row.AlternativeShipDate),
		sql.Named("orderItemId", row.OrderItemId)); err != nil {
		return fmt.Errorf("failed to execute query [%v]: %v", sqlQuery, err)
	} else {
		nRows, err := result.RowsAffected()
		if err != nil {
			return err
		}
		if nRows == 0 {
			return fmt.Errorf("failed to insert CLIENT_ORDER_ITEM row [%v]: orderId=%v, clientProductId=%v", sqlQuery,
				row.OrderId, row.ClientProductId)
		}
	}
	return nil
}
