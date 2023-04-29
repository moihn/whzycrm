package order

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	clientProduct "github.com/moihn/whzycrmgo/internal/pkg/dbmodel/client/product"
	"github.com/moihn/whzycrmgo/internal/pkg/dbmodel/table"
)

type OrderStatus int

const (
	Received OrderStatus = iota
	PI_SENT
	PI_SIGNED
	DEPOSIT_PAID
	IN_PRODUCTION
	SHIPPED
	BALANCE_PAID
	CLOSED
	CANCELLED
)

type PriceModel struct {
	PriceNumber float32
	PriceTypeId int
	CurrencyId  int
}

type ClientOrderItemModel struct {
	ClientOrder         *ClientOrderModel
	OrderItemId         int
	ClientProduct       clientProduct.ClientProductModel
	Quantity            int
	Price               PriceModel
	AddedDate           time.Time
	AlternativeShipDate *time.Time
}

type ClientOrderModel struct {
	OrderId              int
	OrderReference       string
	ClientId             int
	ClientOrderReference string
	CreationDate         time.Time
	ShipmentDate         time.Time
	Status               OrderStatus
	ClientOrderItems     []ClientOrderItemModel
}

func (model *ClientOrderItemModel) Upsert(tx *sql.Tx) (int, error) {
	if model.ClientOrder == nil {
		return 0, fmt.Errorf("ClientOrderItemModel without a ClientOrderModel: ClientProductId=%v, PriceNumber=%v, Quantity=%v",
			model.ClientProduct.ClientProductId, model.Price.PriceNumber, model.Quantity)
	}
	// check if the order item exists
	var row *table.ClientOrderItemRow
	var err error
	if model.OrderItemId > 0 {
		row, err = table.ClientOrderItem_GetByPk(tx, model.OrderItemId)
		if err != nil {
			return 0, err
		}
	}

	if row == nil {
		// insert order item rows
		model.OrderItemId, err = table.ClientOrderItem_Insert(
			tx,
			model.ClientOrder.OrderId,
			model.ClientProduct.ClientProductId,
			model.Quantity,
			model.Price.PriceNumber,
			model.Price.CurrencyId,
			model.AddedDate,
			model.AlternativeShipDate,
		)
		if err != nil {
			return 0, err
		}
		row, err = table.ClientOrderItem_GetByPk(tx, model.OrderItemId)
		if err != nil {
			return 0, err
		}
	} else {
		// update
		row.Update(tx)
	}
	return model.OrderItemId, nil
}

func (clientOrder *ClientOrderModel) Upsert(tx *sql.Tx) (int, error) {
	// check if the order reference exists
	var clientOrderRow *table.ClientOrderRow
	var err error
	if clientOrder.OrderId > 0 {
		clientOrderRow, err = table.ClientOrderGetByPk(tx, clientOrder.OrderId)
		if err != nil {
			return 0, err
		}
	}
	if clientOrderRow == nil && len(clientOrder.OrderReference) > 0 {
		clientOrderRow, err = table.ClientOrderGetByUi2(tx, clientOrder.OrderReference)
		if err != nil {
			return 0, err
		}
	}
	if clientOrderRow == nil && len(clientOrder.ClientOrderReference) > 0 {
		clientOrderRow, err = table.ClientOrderGetByUi1(tx, clientOrder.ClientId, clientOrder.ClientOrderReference)
		if err != nil {
			return 0, err
		}
	}

	if clientOrderRow == nil {
		// check and assign order reference
		if len(clientOrder.OrderReference) == 0 {
			count := 0
			orderRef := "HP" + time.Now().Format("060102") // YYMMDD
			for {
				if count != 0 {
					orderRef = orderRef + strconv.Itoa(count)
				}
				orderRow, err := table.ClientOrderGetByUi2(tx, clientOrder.OrderReference)
				if err != nil {
					return 0, err
				}
				if orderRow == nil {
					clientOrder.OrderReference = orderRef
					break
				}
				count++
			}
		}
		clientOrder.OrderId, err = table.ClientOrder_Insert(
			tx,
			clientOrder.OrderReference,
			clientOrder.ClientId,
			clientOrder.ClientOrderReference,
			clientOrder.CreationDate,
			clientOrder.ShipmentDate,
			int(clientOrder.Status),
		)
		if err != nil {
			return 0, err
		}
		clientOrderRow, err = table.ClientOrderGetByPk(tx, clientOrder.OrderId)
		if err != nil {
			return 0, err
		}
		if clientOrderRow == nil {
			return 0, fmt.Errorf("cannot find or create CLIENT_ORDER row, check system state")
		}
	} else {
		clientOrder.OrderId = clientOrderRow.OrderId
		clientOrderRow.OrderReference = clientOrder.OrderReference
		clientOrderRow.ClientId = clientOrder.ClientId
		clientOrderRow.ClientOrderReference = clientOrder.ClientOrderReference
		clientOrderRow.CreationDate = clientOrder.CreationDate
		clientOrderRow.ShipmentDate = clientOrder.ShipmentDate
		clientOrderRow.StatusId = int(clientOrder.Status)
		err := clientOrderRow.Update(tx)
		if err != nil {
			return 0, err
		}
	}
	// add client order items
	for _, orderItem := range clientOrder.ClientOrderItems {
		_, err := orderItem.Upsert(tx)
		if err != nil {
			return 0, err
		}
	}

	return clientOrder.OrderId, nil
}
