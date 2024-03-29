package product

import (
	"database/sql"

	"github.com/moihn/whzycrmgo/internal/pkg/dbmodel/dbobject"
)

type ClientProduct struct {
	ClientProductId int
	ClientId        int
	Reference       string
	Description     *string
	Narrative       *string
	Barcode         *string
}

type ClientProductModel struct {
	ClientProduct
	VendorProductIds []int
}

func AddClientProduct(tx *sql.Tx, clientProductModel ClientProductModel) (int, error) {
	newClientProductRow := dbobject.ClientProductRow{
		ClientId:    int64(clientProductModel.ClientId),
		Reference:   clientProductModel.Reference,
		Description: clientProductModel.Description,
		Narrative:   clientProductModel.Narrative,
		Barcode:     clientProductModel.Barcode,
	}
	err := newClientProductRow.Insert(tx)
	if err != nil {
		return 0, err
	}

	return int(newClientProductRow.ClientProductId), nil
}

func RemoveClientProduct(tx *sql.Tx, clntProdId int) error {
	if row, err := dbobject.ClientProductGetByClientProductId(tx, int64(clntProdId)); err != nil {
		return err
	} else {
		row.DeletedInd = "F"
		return row.Update(tx)
	}
}

func GetClientProductModels(tx *sql.Tx, clientId int, productRef string) ([]ClientProductModel, error) {
	clientProducts, err := dbobject.ClientProductPopulateByClientId(tx, int64(clientId))
	if err != nil {
		return nil, err
	}
	var result []ClientProductModel
	for _, clientProduct := range clientProducts {
		if len(productRef) == 0 || productRef == clientProduct.Reference {
			clientProductModel := ClientProductModel{
				ClientProduct{
					ClientProductId: int(clientProduct.ClientProductId),
					ClientId:        int(clientProduct.ClientId),
					Reference:       clientProduct.Reference,
				},
				[]int{},
			}
			clientProductModel.Description = clientProduct.Description
			clientProductModel.Narrative = clientProduct.Narrative
			clientProductModel.Barcode = clientProduct.Barcode
			clientProductItems, err := dbobject.ClientProductItemPopulateByClientProductId(tx, int64(clientProduct.ClientProductId))
			if err != nil {
				return nil, err
			}
			for _, clientProductItem := range clientProductItems {
				clientProductModel.VendorProductIds = append(clientProductModel.VendorProductIds, int(clientProductItem.VendorProductId))
			}
			result = append(result, clientProductModel)
		}
	}
	return result, nil
}

func GetClientProductModelById(tx *sql.Tx, clientProductId int) (*ClientProductModel, error) {
	clientProduct, err := dbobject.ClientProductGetByClientProductId(tx, int64(clientProductId))
	if err != nil {
		return nil, err
	}
	var result *ClientProductModel
	if clientProduct != nil {
		clientProductModel := ClientProductModel{
			ClientProduct{
				ClientProductId: int(clientProduct.ClientProductId),
				ClientId:        int(clientProduct.ClientId),
				Reference:       clientProduct.Reference,
			},
			[]int{},
		}
		clientProductModel.Description = clientProduct.Description
		clientProductModel.Narrative = clientProduct.Narrative
		clientProductModel.Barcode = clientProduct.Barcode
		clientProductItems, err := dbobject.ClientProductItemPopulateByClientProductId(tx, int64(clientProduct.ClientProductId))
		if err != nil {
			return nil, err
		}
		for _, clientProductItem := range clientProductItems {
			clientProductModel.VendorProductIds = append(clientProductModel.VendorProductIds, int(clientProductItem.VendorProductId))
		}
	}
	return result, nil
}
