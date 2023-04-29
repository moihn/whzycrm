//go:generate oapi-codegen -package oapi -generate types,chi-server,spec -o service.gen.go openapi.yaml

package oapi

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/deepmap/oapi-codegen/pkg/types"
	"github.com/moihn/whzycrmgo/internal/pkg/api"
	"github.com/moihn/whzycrmgo/internal/pkg/dbmodel/client/quotation"
	"github.com/moihn/whzycrmgo/internal/pkg/dbmodel/fx"
	"github.com/moihn/whzycrmgo/internal/pkg/dbmodel/supplier"
	vendorProduct "github.com/moihn/whzycrmgo/internal/pkg/dbmodel/supplier/product"
	"github.com/moihn/whzycrmgo/internal/pkg/util"
)

func sendResponse(w http.ResponseWriter, status *util.Status) {
	if status != nil {
		w.WriteHeader(status.Code)
		w.Write(status.ErrorBytes())
	}
}

func sendResponseWithBody(w http.ResponseWriter, response interface{}, errRef string) {
	w.WriteHeader(http.StatusOK)
	body, err := json.Marshal(response)
	if err != nil {
		status := util.NewInternalServiceErrorStatus(fmt.Sprintf("Failed to deserialize response body: %v", err), errRef)
		sendResponse(w, &status)
		return
	}
	w.Write(body)
}

type serviceHandlerImpl struct {
	timeout       int
	bodySizeLimit int64
	database      *sql.DB
}

func NewServiceHandler(timeout int, bodySizeLimit int64, db *sql.DB) ServerInterface {
	return &serviceHandlerImpl{
		timeout:       timeout,
		bodySizeLimit: bodySizeLimit,
		database:      db,
	}
}

func (handler *serviceHandlerImpl) getDbTransaction(w http.ResponseWriter) *sql.Tx {
	tx, err := handler.database.Begin()
	if err != nil {
		status := util.NewInternalServiceErrorStatus(fmt.Sprintf("Failed to get database transaction: %v", err), "GET_DB_TRANSACTION")
		sendResponse(w, &status)
		return nil
	}
	return tx
}

//
// Vendor
//

func (handler *serviceHandlerImpl) GetVendors(w http.ResponseWriter, r *http.Request) {
	/*userId := r.Header.Get("x-remote-user")
	userFormat := r.Header.Get("x-remote-user-format")
	if len(userFormat) == 0 {
		userFormat = "username"
	}
	*/
	tx := handler.getDbTransaction(w)
	if tx == nil {
		return
	}
	// this is a read-only operation, roll back any potential change
	defer tx.Rollback()

	// get data from model
	vendors, status := supplier.GetAll(tx)
	if status != nil {
		sendResponse(w, status)
		return
	}

	response := GetVendorsResponse{
		Vendors: &[]Vendor{},
	}
	for _, vendor := range vendors {
		*response.Vendors = append(*response.Vendors, Vendor{
			Name:     vendor.Name,
			VendorId: vendor.Id,
		})
	}

	sendResponseWithBody(w, response, "GET_VENDORS_1")
}

//
// Vendor Product
//

func (handler *serviceHandlerImpl) GetVendorProducts(w http.ResponseWriter, r *http.Request, vendorId int) {
	tx := handler.getDbTransaction(w)
	if tx == nil {
		return
	}
	// this is a read-only operation, roll back any potential change
	defer tx.Rollback()
	productSummaries, err := vendorProduct.GetVendorProducts(tx, vendorId)
	if err != nil {
		status := util.NewInternalServiceErrorStatus(err.Error(), "GET_VENDOR_PRODUCTS_0")
		sendResponse(w, &status)
		return
	}
	response := GetVendorProductsResponse{}
	for _, summary := range productSummaries {
		response.VendorProducts = append(response.VendorProducts, VendorProductSummary{
			VendorProductId: summary.VendorProductId,
			Reference:       summary.Reference,
			Description:     summary.Description,
			VendorId:        summary.VendorId,
		})
	}

	sendResponseWithBody(w, response, "GET_VENDOR_PRODUCTS_1")
}

func (handler *serviceHandlerImpl) GetVendorProduct(w http.ResponseWriter, r *http.Request, vendorId int, vendorProductReference string) {
	tx := handler.getDbTransaction(w)
	if tx == nil {
		return
	}
	// this is a read-only operation, roll back any potential change
	defer tx.Rollback()

	// get vendor product from model
	vp, status := vendorProduct.GetVendorProduct(tx, vendorId, vendorProductReference)
	if status != nil {
		sendResponse(w, status)
		return
	}

	if vp == nil {
		errMsg := fmt.Sprintf("Expecting valid vendor product being found for vendor %v reference %v", vendorId, vendorProductReference)
		status := util.NewNotFoundStatus(errMsg)
		sendResponse(w, &status)
		return
	}

	response := GetVendorProductResponse{
		VendorProduct: &VendorProduct{
			VendorId:        vendorId,
			VendorProductId: &vp.VendorProductId,
			Reference:       vendorProductReference,
			Description:     &vp.Description,
			MaterialTypeId:  vp.MaterialTypeId,
			Packing:         Packing{},
		},
	}

	if vp.UnitSize.Length != nil {
		response.VendorProduct.Packing.Length = &Measure{
			Value: vp.UnitSize.Length.Value,
			Unit:  vp.UnitSize.Length.Unit,
		}
	}

	if vp.UnitSize.Width != nil {
		response.VendorProduct.Packing.Width = &Measure{
			Value: vp.UnitSize.Width.Value,
			Unit:  vp.UnitSize.Width.Unit,
		}
	}

	if vp.UnitSize.Height != nil {
		response.VendorProduct.Packing.Height = &Measure{
			Value: vp.UnitSize.Height.Value,
			Unit:  vp.UnitSize.Height.Unit,
		}
	}

	if len(vp.PriceList) > 0 {
		for i := range vp.PriceList {
			response.VendorProduct.PriceHistory = append(response.VendorProduct.PriceHistory, PriceChange{
				StartDate: types.Date{Time: vp.PriceList[i].StartDate},
				Price: Price{
					CurrencyId:  vp.PriceList[i].CurrencyId,
					Value:       vp.PriceList[i].Amount,
					PriceTypeId: vp.PriceList[i].TypeId,
				},
			})
		}
	}

	// search quotation price

	// search order price
	sendResponseWithBody(w, response, "GET_VENDOR_PRODUCT_1")
}

func (handler *serviceHandlerImpl) AddVendorProduct(w http.ResponseWriter, r *http.Request) {
	var newProduct AddVendorProductJSONRequestBody
	status := util.DecodeJsonBodyAsObject(w, r, handler.bodySizeLimit, &newProduct)
	if status != nil {
		sendResponse(w, status)
		return
	}

	tx := handler.getDbTransaction(w)
	if tx == nil {
		return
	}

	description := ""
	if newProduct.Description != nil {
		description = *newProduct.Description
	}
	testPerformed := false
	if newProduct.TestPerformed != nil {
		testPerformed = *newProduct.TestPerformed
	}
	newVendorProductSpec := api.VendorNewProduct{
		VendorId:       newProduct.VendorId,
		Reference:      newProduct.ProductReference,
		Description:    description,
		MaterialTypeId: newProduct.MaterialTypeId,
		ProductTypeId:  newProduct.ProductTypeId,
		UnitTypeId:     newProduct.UnitTypeId,
		Length:         newProduct.Length,
		Width:          newProduct.Width,
		Height:         newProduct.Height,
		Weight:         newProduct.Weight,
		Price:          newProduct.Price,
		PriceCcyId:     newProduct.CcyId,
		PriceTypeId:    newProduct.PriceTypeId,
		Moq:            newProduct.Moq,
		TestPerformed:  testPerformed,
	}

	newVendorProductId, err := api.AddVendorProduct(tx, newVendorProductSpec)
	if err != nil {
		status := util.NewBadRequestStatus(err.Error())
		sendResponse(w, &status)
		return
	}
	response := AddVendorProductResponse{
		ProductId: newVendorProductId,
	}

	if err := tx.Commit(); err != nil {
		status := util.NewInternalServiceErrorStatus(fmt.Sprintf("failed to commit new vendor products: %v", err), "ADD_VENDOR_PRODUCT_1")
		sendResponse(w, &status)
	}

	sendResponseWithBody(w, response, "ADD_VENDOR_PRODUCT_2")
}

//
// Vendor Product Price
//

func (handler *serviceHandlerImpl) UpdateVendorPrice(w http.ResponseWriter, r *http.Request) {
	var newPriceUpdateBody UpdateVendorPriceJSONRequestBody
	status := util.DecodeJsonBodyAsObject(w, r, handler.bodySizeLimit, &newPriceUpdateBody)
	if status != nil {
		sendResponse(w, status)
		return
	}

	tx := handler.getDbTransaction(w)
	if tx == nil {
		return
	}
	// Roll back any potential uncommitted change
	defer tx.Rollback()
	err := vendorProduct.UpdatePrice(
		tx, newPriceUpdateBody.VendorProductId,
		newPriceUpdateBody.PriceUpdate.Price.Value,
		newPriceUpdateBody.PriceUpdate.Price.CurrencyId,
		newPriceUpdateBody.PriceUpdate.Price.PriceTypeId)

	if err != nil {
		status := util.NewBadRequestStatus(fmt.Sprintf("failed to update vendor product price: vendorProductId=%v, price=%v", newPriceUpdateBody.VendorProductId, newPriceUpdateBody.PriceUpdate.Price.Value))
		sendResponse(w, &status)
		return
	}
	if err := tx.Commit(); err != nil {
		status := util.NewInternalServiceErrorStatus(fmt.Sprintf("Failed to commit vendor price update: %v", err), "UPDATE_VENDOR_PRICE_0")
		sendResponse(w, &status)
	}
}

//
// Client Quotation
//

func (handler *serviceHandlerImpl) CreateClientQuotation(w http.ResponseWriter, r *http.Request) {
	var newQuotationBody CreateClientQuotationJSONRequestBody
	status := util.DecodeJsonBodyAsObject(w, r, handler.bodySizeLimit, &newQuotationBody)
	if status != nil {
		sendResponse(w, status)
		return
	}

	tx := handler.getDbTransaction(w)
	if tx == nil {
		return
	}
	// this is a read-only operation, roll back any potential change
	defer tx.Rollback()

	quoteCcy, status := fx.GetCurrencyByIsoSymbol(tx, "USD")
	if status != nil {
		sendResponse(w, status)
		return
	}

	if len(newQuotationBody.VendorProductIds) == 0 {
		// no product found in the list
		return
	}
	newQuote := quotation.Quotation{
		ClientId:   newQuotationBody.ClientId,
		CurrencyId: quoteCcy.Id,
		Date:       time.Now(),
	}

	for _, vendorProductId := range newQuotationBody.VendorProductIds {
		newQuote.Items = append(newQuote.Items, quotation.ItemInQuotation{
			VendorProductId: vendorProductId,
		})
	}

	if status = quotation.CreateClientQuotation(tx, newQuote); status != nil {
		sendResponse(w, status)
		return
	}

	if err := tx.Commit(); err != nil {
		status := util.NewInternalServiceErrorStatus(fmt.Sprintf("Failed to commit new client quotation: %v", err), "CREATE_CLIENT_QUOTATION_1")
		sendResponse(w, &status)
	}
}

func (handler *serviceHandlerImpl) UpdateClientQuotation(w http.ResponseWriter, r *http.Request) {

}

func (handler *serviceHandlerImpl) GetClientQuotations(w http.ResponseWriter, r *http.Request, params GetClientQuotationsParams) {
	tx := handler.getDbTransaction(w)
	if tx == nil {
		return
	}
	// this is a read-only operation, roll back any potential change
	defer tx.Rollback()

	var startDate time.Time
	if params.StartDate != nil {
		startDate = params.StartDate.Time
	}
	quotes, status := quotation.FindClientQuotations(tx, params.ClientId, startDate)
	if status != nil {
		sendResponse(w, status)
		return
	}

	response := GetClientQuotationsResponse{}
	if len(quotes) > 0 {
		response.Quotations = &[]ClientQuotation{}
	}
	for _, quote := range quotes {
		clientQuote := ClientQuotation{
			ClientId: quote.ClientId,
			UpdatedDate: types.Date{
				Time: quote.Date,
			},
			Items: []ItemInClientQuotation{},
		}
		for _, item := range quote.Items {
			clientQuote.Items = append(clientQuote.Items, ItemInClientQuotation{
				VendorProductId: item.VendorProductId,
				Price:           item.PriceAmount,
				Moq:             item.MOQ,
			})

		}
		*response.Quotations = append(*response.Quotations, clientQuote)
	}

	sendResponseWithBody(w, response, "GET_CLIENT_QUOTATIONS_0")
}

//
// Client Product
//

func (handler *serviceHandlerImpl) GetClientProduct(w http.ResponseWriter, r *http.Request, params GetClientProductParams) {
	tx := handler.getDbTransaction(w)
	if tx == nil {
		return
	}
	// this is a read-only operation, roll back any potential change
	defer tx.Rollback()

	clientProducts, err := api.GetClientProducts(tx, params.ClientId, params.ProductReference)
	if err != nil {
		status := util.NewInternalServiceErrorStatus(err.Error(), "GET_CLIENT_PRODUCT_1")
		sendResponse(w, &status)
		return
	}
	response := GetClientProductResponse{}
	for _, clientProduct := range clientProducts {
		clientProductResp := ClientProduct{
			ClientProductId:  clientProduct.Id,
			ProductReference: clientProduct.Reference,
			Description:      clientProduct.Description,
			Narrative:        clientProduct.Narrative,
			Barcode:          clientProduct.Barcode,
		}

		clientProductResp.VendorProductIds = append(clientProductResp.VendorProductIds, clientProduct.VendorProductIds...)
		response.ClientProducts = append(response.ClientProducts, clientProductResp)
	}
	sendResponseWithBody(w, response, "GET_CLIENT_PRODUCT_2")
}

func (handler *serviceHandlerImpl) AddClientProduct(w http.ResponseWriter, r *http.Request) {
	var newClientProduct NewClientProduct
	status := util.DecodeJsonBodyAsObject(w, r, handler.bodySizeLimit, &newClientProduct)
	if status != nil {
		sendResponse(w, status)
		return
	}
	tx := handler.getDbTransaction(w)
	if tx == nil {
		return
	}
	// Roll back any potential uncommitted change
	defer tx.Rollback()
	clientProductId, err := api.AddClientProduct(tx, api.ClientProduct{
		Description: newClientProduct.Description,
		ClientId:    newClientProduct.ClientId,
		Reference:   newClientProduct.ProductReference,
		Barcode:     newClientProduct.Barcode,
		Narrative:   newClientProduct.Narrative,
	})
	if err != nil {
		status := util.NewInternalServiceErrorStatus(err.Error(), "ADD_CLIENT_PRODUCT_1")
		sendResponse(w, &status)
		return
	}
	response := AddClientProductResponse{
		ClientProductId: clientProductId,
	}
	sendResponseWithBody(w, response, "ADD_CLIENT_PRODUCT_2")
}

func (handler *serviceHandlerImpl) RemoveClientProduct(w http.ResponseWriter, r *http.Request, clientProductId int) {
	tx := handler.getDbTransaction(w)
	if tx == nil {
		return
	}
	defer tx.Rollback()

	err := api.RemoveClientProduct(tx, clientProductId)
	if err != nil {
		status := util.NewInternalServiceErrorStatus(err.Error(), "REMOVE_CLIENT_PRODUCT_1")
		sendResponse(w, &status)
		return
	}
}

func (handler *serviceHandlerImpl) AddClientProductItem(w http.ResponseWriter, r *http.Request, clientProductId int) {
	var newClientProductItem NewClientProductItem
	status := util.DecodeJsonBodyAsObject(w, r, handler.bodySizeLimit, &newClientProductItem)
	if status != nil {
		sendResponse(w, status)
		return
	}
	tx := handler.getDbTransaction(w)
	if tx == nil {
		return
	}
	// Roll back any potential uncommitted change
	defer tx.Rollback()
	newClientProdItemIn := api.ClientProductItem{
		Description: newClientProductItem.Narrative,
	}
	api.AddClientProductItem(tx, clientProductId, newClientProductItem)
}

func (handler *serviceHandlerImpl) RemoveClientProductItem(w http.ResponseWriter, r *http.Request, clientProductId int, itemId int) {

}

//
// Client
//

func (handler *serviceHandlerImpl) GetClients(w http.ResponseWriter, r *http.Request) {
	tx := handler.getDbTransaction(w)
	if tx == nil {
		return
	}
	// this is a read-only operation, roll back any potential change
	defer tx.Rollback()

	clients := api.GetAllClients(tx)
	response := GetClientsResponse{}
	for _, client := range clients {
		response.Clients = append(response.Clients, Client{
			Id:        client.Id,
			Name:      client.Name,
			CountryId: client.CountryId,
		})
	}
	sendResponseWithBody(w, response, "GET_CLIENTS")
}

//
// Client Order
//

func (handler *serviceHandlerImpl) AddClientOrder(w http.ResponseWriter, r *http.Request) {
	var addClientOrderRequest AddClientOrderJSONRequestBody
	status := util.DecodeJsonBodyAsObject(w, r, handler.bodySizeLimit, &addClientOrderRequest)
	if status != nil {
		sendResponse(w, status)
		return
	}

	tx := handler.getDbTransaction(w)
	if tx == nil {
		return
	}

	orderReference := ""
	if addClientOrderRequest.OrderReference != nil {
		orderReference = *addClientOrderRequest.OrderReference
	}

	clientOrderReference := ""
	if addClientOrderRequest.ClientOrderReference != nil {
		clientOrderReference = *addClientOrderRequest.ClientOrderReference
	}

	newClientOrder := api.NewClientOrder{
		OrderReference:       orderReference,
		ClientId:             addClientOrderRequest.ClientId,
		ClientOrderReference: clientOrderReference,
		CreationDate:         addClientOrderRequest.OrderDate.Time,
		ShipmentDate:         addClientOrderRequest.DeliverDate.Time,
	}

	for _, orderItem := range addClientOrderRequest.ClientProducts {
		orderItemId := -1
		if orderItem.ClientOrderItemId != nil {
			orderItemId = *orderItem.ClientOrderItemId
		}
		var alternativeShipDate *time.Time
		if orderItem.AlternativeShipDate != nil {
			alternativeShipDate = &orderItem.AlternativeShipDate.Time
		}
		newClientOrder.ClientOrderItems = append(newClientOrder.ClientOrderItems, api.ClientOrderItem{
			OrderItemId:     orderItemId,
			ClientProductId: orderItem.ClientProductId,
			Quantity:        orderItem.Quantity,
			Price: api.Price{
				PriceNumber: orderItem.Price.Value,
				CurrencyId:  orderItem.Price.CurrencyId,
				PriceTypeId: orderItem.Price.PriceTypeId,
			},
			AddedDate:           orderItem.AddedDate.Time,
			AlternativeShipDate: alternativeShipDate,
		})
	}

	newOrderId, err := api.AddClientOrder(tx, newClientOrder)
	if err != nil {
		status := util.NewBadRequestStatus(err.Error())
		sendResponse(w, &status)
		return
	}

	response := AddClientOrderResponse{
		ClientOrderId: newOrderId,
	}

	if err := tx.Commit(); err != nil {
		status := util.NewInternalServiceErrorStatus(fmt.Sprintf("failed to commit new vendor products: %v", err), "ADD_VENDOR_PRODUCT_1")
		sendResponse(w, &status)
	}

	sendResponseWithBody(w, response, "ADD_VENDOR_PRODUCT_2")
}
