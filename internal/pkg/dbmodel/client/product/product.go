package product

import (
	"database/sql"

	"github.com/moihn/whzycrmgo/internal/pkg/dbmodel/dbobject"
)

type ClientProductModel struct {
	ClientProductId  int
	ClientId         int
	Reference        string
	Description      string
	Narrative        string
	Barcode          string
	VendorProductIds []int
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
				ClientProductId: int(clientProduct.ClientProductId),
				ClientId:        int(clientProduct.ClientId),
				Reference:       clientProduct.Reference,
				Description:     clientProduct.Description,
				Narrative:       clientProduct.Narrative,
				Barcode:         clientProduct.Barcode,
			}

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
			ClientProductId: int(clientProduct.ClientProductId),
			ClientId:        int(clientProduct.ClientId),
			Reference:       clientProduct.Reference,
			Description:     clientProduct.Description,
			Narrative:       clientProduct.Narrative,
			Barcode:         clientProduct.Barcode,
		}

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
