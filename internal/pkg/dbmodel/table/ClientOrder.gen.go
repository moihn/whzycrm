package table

import (
	"database/sql"
	"fmt"
	"time"
)

type ClientOrderRow struct {
	OrderId              int
	OrderReference       string
	ClientId             int
	ClientOrderReference string
	CreationDate         time.Time
	ShipmentDate         time.Time
	StatusId             int
}

func ClientOrder_Insert(
	tx *sql.Tx,
	orderReference string,
	clientId int,
	clientOrderReference string,
	creationDate time.Time,
	shipmentDate time.Time,
	statusId int,
) (int, error) {
	sqlQuery := `
	INSERT INTO CLIENT_ORDER
	(ORDER_REFERENCE, CLIENT_ID, CLIENT_ORDER_REFERENCE, CREATION_DATE, SHIPMENT_DATE, STATUS_ID)
	VALUES
	(:orderReference, :clientId, :clientOrderRefernce, :creationDate, :shipmentDate, :statusId)
	RETURNING ORDER_ID into :orderId
	`
	var orderId int
	if result, err := tx.Exec(sqlQuery,
		sql.Named("orderReference", orderReference),
		sql.Named("clientId", clientId),
		sql.Named("clientOrderReference", clientOrderReference),
		sql.Named("creationDate", creationDate),
		sql.Named("shipmentDate", shipmentDate),
		sql.Named("statusId", statusId),
		sql.Named("orderId", sql.Out{Dest: &orderId})); err != nil {
		return 0, fmt.Errorf("failed to execute query [%v]: %v", sqlQuery, err)
	} else {
		nRows, err := result.RowsAffected()
		if err != nil {
			return 0, err
		}
		if nRows == 0 {
			fmt.Errorf("failed to insert CLIENT_ORDER row [%v]: orderReference=%v, clientId=%v, clientOrderReference=%v", sqlQuery,
				orderReference, clientId, clientOrderReference)
		}
	}
	return orderId, nil
}

func (row *ClientOrderRow) Update(tx *sql.Tx) error {
	sqlQuery := `
		UPDATE CLIENT_ORDER
		SET ORDER_REFERENCE = :orderReference,
			CLIENT_ID = :clientId,
			CLIENT_ORDER_REFERENCE = :clientOrderReference,
			CREATION_DATE = :creationDate,
			SHIPMENT_DATE = :shipmentDate,
			STATUS_ID = :statusId)
		WHERE
			CLIENT_ORDER_ID = :orderId`
	if result, err := tx.Exec(sqlQuery,
		sql.Named("orderReference", row.OrderReference),
		sql.Named("clientId", row.ClientId),
		sql.Named("clientOrderReference", row.ClientOrderReference),
		sql.Named("creationDate", row.CreationDate),
		sql.Named("shipmentDate", row.ShipmentDate),
		sql.Named("statusId", row.StatusId),
		sql.Named("orderId", row.OrderId)); err != nil {
		return fmt.Errorf("failed to execute query [%v]: %v", sqlQuery, err)
	} else {
		nRows, err := result.RowsAffected()
		if err != nil {
			return err
		}
		if nRows == 0 {
			return fmt.Errorf("failed to update CLIENT_ORDER row [%v]: orderReference=%v, clientId=%v, clientOrderReference=%v", sqlQuery,
				row.OrderReference, row.ClientId, row.ClientOrderReference)
		}
	}
	return nil
}

func ClientOrderGetByPk(tx *sql.Tx, clientOrderId int) (*ClientOrderRow, error) {
	sqlQuery := `
		select
			ORDER_ID,
			ORDER_REFERENCE,
			CLIENT_ID,
			CLIENT_ORDER_REFERENCE,
			CREATION_DATE,
			SHIPMENT_DATE,
			STATUS_ID
		from CLIENT_ORDER
		where CLIENT_ID = :clientId
		  and CLIENT_ORDER_REFERENCE = :clientOrderReference
	`
	rows, err := tx.Query(sqlQuery,
		sql.Named("clientOrderId", clientOrderId),
	)
	if err != nil {
		return nil, err
	}
	if rows.Next() {
		result := ClientOrderRow{}
		if err := rows.Scan(
			&result.OrderId,
			&result.OrderReference,
			&result.ClientId,
			&result.ClientOrderReference,
			&result.CreationDate,
			&result.ShipmentDate,
			&result.StatusId,
		); err != nil {
			return nil, err
		}
		return &result, nil
	}
	return nil, nil
}

func ClientOrderGetByUi1(tx *sql.Tx, clientId int, clientOrderReference string) (*ClientOrderRow, error) {
	sqlQuery := `
		select
			ORDER_ID,
			ORDER_REFERENCE,
			CLIENT_ID,
			CLIENT_ORDER_REFERENCE,
			CREATION_DATE,
			SHIPMENT_DATE,
			STATUS_ID
		from CLIENT_ORDER
		where CLIENT_ID = :clientId
		  and CLIENT_ORDER_REFERENCE = :clientOrderReference
	`
	rows, err := tx.Query(sqlQuery,
		sql.Named("clientId", clientId),
		sql.Named("clientOrderReference", clientOrderReference),
	)
	if err != nil {
		return nil, err
	}
	if rows.Next() {
		result := ClientOrderRow{}
		if err := rows.Scan(
			&result.OrderId,
			&result.OrderReference,
			&result.ClientId,
			&result.ClientOrderReference,
			&result.CreationDate,
			&result.ShipmentDate,
			&result.StatusId,
		); err != nil {
			return nil, err
		}
		return &result, nil
	}
	return nil, nil
}

func ClientOrderGetByUi2(tx *sql.Tx, orderReference string) (*ClientOrderRow, error) {
	sqlQuery := `
		select
			ORDER_ID,
			ORDER_REFERENCE,
			CLIENT_ID,
			CLIENT_ORDER_REFERENCE,
			CREATION_DATE,
			SHIPMENT_DATE,
			STATUS_ID
		from CLIENT_ORDER
		where ORDER_REFERENCE = :orderReference
	`
	rows, err := tx.Query(sqlQuery,
		sql.Named("orderReference", orderReference),
	)
	if err != nil {
		return nil, err
	}
	if rows.Next() {
		result := ClientOrderRow{}
		if err := rows.Scan(
			&result.OrderId,
			&result.OrderReference,
			&result.ClientId,
			&result.ClientOrderReference,
			&result.CreationDate,
			&result.ShipmentDate,
			&result.StatusId,
		); err != nil {
			return nil, err
		}
		return &result, nil
	}
	return nil, nil
}
