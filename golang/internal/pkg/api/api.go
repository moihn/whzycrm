package api

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/moihn/whzycrmgo/internal/pkg/dbmodel/client"
	clientOrder "github.com/moihn/whzycrmgo/internal/pkg/dbmodel/client/order"
	clientProduct "github.com/moihn/whzycrmgo/internal/pkg/dbmodel/client/product"
	vendorProduct "github.com/moihn/whzycrmgo/internal/pkg/dbmodel/supplier/product"
)

type Client struct {
	Id        int
	Name      string
	CountryId int
}

type ClientProduct struct {
	Id               int
	Description      *string
	ClientId         int
	Reference        string
	Narrative        *string
	Barcode          *string
	VendorProductIds []int
}

type ClientProductItem struct {
	Id              int
	Narrative       *string
	VendorProductId int
}

type VendorProduct struct {
	Id             int
	VendorId       int
	Reference      string
	Description    string
	MaterialTypeId *int
	ProductTypeId  int
	UnitTypeId     *int
	Length         *int
	Width          *int
	Height         *int
	Weight         *int
}

type VendorNewProduct struct {
	VendorId       int
	Reference      string
	Description    string
	MaterialTypeId *int
	TestPerformed  bool
	ProductTypeId  int
	UnitTypeId     *int
	Length         *float32
	Width          *float32
	Height         *float32
	Weight         *float32
	Price          *float32
	PriceCcyId     *int
	PriceTypeId    *int
	Moq            *int
}

type Price struct {
	CurrencyId  int
	PriceNumber float32
	PriceTypeId int
}

type ClientOrderItem struct {
	OrderItemId         int
	ClientProductId     int
	Quantity            int
	Price               Price
	AddedDate           time.Time
	AlternativeShipDate *time.Time
}

type NewClientOrder struct {
	OrderReference       string
	ClientId             int
	ClientOrderReference string
	CreationDate         time.Time
	ShipmentDate         time.Time
	ClientOrderItems     []ClientOrderItem
}

func AddVendorProduct(tx *sql.Tx, newProduct VendorNewProduct) (int, error) {
	newProductSpec := vendorProduct.VendorNewProduct{
		VendorId:       newProduct.VendorId,
		Reference:      newProduct.Reference,
		Description:    newProduct.Description,
		MaterialTypeId: newProduct.MaterialTypeId,
		TestPerformed:  newProduct.TestPerformed,
		ProductTypeId:  newProduct.ProductTypeId,
		UnitTypeId:     newProduct.UnitTypeId,
		Length:         newProduct.Length,
		Width:          newProduct.Width,
		Height:         newProduct.Height,
		Weight:         newProduct.Weight,
		Price:          newProduct.Price,
		PriceCcyId:     newProduct.PriceCcyId,
		PriceTypeId:    newProduct.PriceTypeId,
		Moq:            newProduct.Moq,
	}
	return vendorProduct.AddNewVendorProduct(tx, newProductSpec)
}

func GetAllClients(tx *sql.Tx) []Client {
	var result []Client
	clientModels := client.GetAllClientModels(tx)
	for _, clientModel := range clientModels {
		result = append(result, Client{
			Id:        clientModel.Id,
			Name:      clientModel.Name,
			CountryId: clientModel.CountryId,
		})
	}
	return result
}

func AddClientProduct(tx *sql.Tx, clntProduct ClientProduct) (int, error) {
	model := clientProduct.ClientProductModel{
		clientProduct.ClientProduct{
			Description: clntProduct.Description,
			ClientId:    clntProduct.ClientId,
			Reference:   clntProduct.Reference,
			Narrative:   clntProduct.Narrative,
			Barcode:     clntProduct.Barcode,
		},
		[]int{},
	}
	return clientProduct.AddClientProduct(tx, model)
}

func AddClientProductItem(tx *sql.Tx, clntProdId int, clntProdItem ClientProductItem) (int, error) {
	GetClientProductById(tx, clntProdId)
}

func RemoveClientProduct(tx *sql.Tx, clntProdId int) error {
	return clientProduct.RemoveClientProduct(tx, clntProdId)
}

func GetClientProducts(tx *sql.Tx, clientId int, productRef string) ([]ClientProduct, error) {
	var result []ClientProduct
	clientProductModels, err := clientProduct.GetClientProductModels(tx, clientId, productRef)
	if err != nil {
		return nil, err
	}
	for _, clientProductModel := range clientProductModels {
		result = append(result, ClientProduct{
			Id:               clientProductModel.ClientProductId,
			Description:      clientProductModel.Description,
			ClientId:         clientProductModel.ClientId,
			Reference:        clientProductModel.Reference,
			Narrative:        clientProductModel.Narrative,
			Barcode:          clientProductModel.Barcode,
			VendorProductIds: clientProductModel.VendorProductIds,
		})
	}
	return result, nil
}

func GetClientProductById(tx *sql.Tx, clientProductId int) (*ClientProduct, error) {
	var result *ClientProduct
	clientProductModel, err := clientProduct.GetClientProductModelById(tx, clientProductId)
	if err != nil {
		return nil, err
	}
	if clientProductModel != nil {
		result = &ClientProduct{
			Id:               clientProductModel.ClientProductId,
			Description:      clientProductModel.Description,
			ClientId:         clientProductModel.ClientId,
			Reference:        clientProductModel.Reference,
			Narrative:        clientProductModel.Narrative,
			Barcode:          clientProductModel.Barcode,
			VendorProductIds: clientProductModel.VendorProductIds,
		}
	}
	return result, nil
}

func AddClientOrder(tx *sql.Tx, newClientOrder NewClientOrder) (int, error) {
	newClientOrderModel := clientOrder.ClientOrderModel{
		OrderId:              -1,
		OrderReference:       newClientOrder.OrderReference,
		ClientId:             newClientOrder.ClientId,
		ClientOrderReference: newClientOrder.ClientOrderReference,
		CreationDate:         newClientOrder.CreationDate,
		ShipmentDate:         newClientOrder.ShipmentDate,
		Status:               0,
	}
	for _, clientOrderItem := range newClientOrder.ClientOrderItems {
		clientProductModel, err := clientProduct.GetClientProductModelById(tx, clientOrderItem.ClientProductId)
		if err != nil {
			return 0, err
		}
		if clientProductModel == nil {
			return -1, fmt.Errorf("client product with ID %v is not found in database", clientOrderItem.ClientProductId)
		}
		newClientOrderModel.ClientOrderItems = append(newClientOrderModel.ClientOrderItems, clientOrder.ClientOrderItemModel{
			ClientOrder:   &newClientOrderModel,
			OrderItemId:   clientOrderItem.OrderItemId,
			ClientProduct: *clientProductModel,
			Quantity:      clientOrderItem.Quantity,
			Price: clientOrder.PriceModel{
				PriceNumber: clientOrderItem.Price.PriceNumber,
				PriceTypeId: clientOrderItem.Price.PriceTypeId,
				CurrencyId:  clientOrderItem.Price.CurrencyId,
			},
			AddedDate:           clientOrderItem.AddedDate,
			AlternativeShipDate: clientOrderItem.AlternativeShipDate,
		})
	}

	return newClientOrderModel.Upsert(tx)
}
