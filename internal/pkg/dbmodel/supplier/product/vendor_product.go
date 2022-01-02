package product

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/moihn/whzycrmgo/internal/pkg/dbmodel/common"
	"github.com/moihn/whzycrmgo/internal/pkg/dbmodel/table"
	"github.com/moihn/whzycrmgo/internal/pkg/util"
)

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

type VendorProduct struct {
	VendorId        int
	Reference       string
	VendorProductId int
	Description     string
	MaterialTypeId  *int
	UnitSize        common.Packing
	PriceList       []common.Price
	MoqList         []common.MoQ
	Packing         common.Packing
}

func getBasicVendorProductInfoById(tx *sql.Tx, vendorProductId int) (*VendorProduct, *util.Status) {
	// search for vendor_product
	sqlQuery := `
		select VENDOR_ID, REFERENCE, DESCRIPTION, MATERIAL_TYPE_ID, LENGTH, WIDTH, HEIGHT, WEIGHT
		from VENDOR_PRODUCT
		where VENDOR_PRODUCT_ID = :vendorProductId
	`
	rows, err := tx.Query(sqlQuery, sql.Named("vendorProductId", vendorProductId))
	if err != nil {
		status := util.NewInternalServiceErrorStatus(fmt.Sprintf("Failed to run query %v: %v", sqlQuery, err), "GET_BASIC_VENDOR_PRODUCT_BY_ID_0")
		return nil, &status
	}
	defer rows.Close()

	var vendorId int
	var vendorProductReference, description string
	var materialTypeId sql.NullInt64
	var length, width, height, weight sql.NullFloat64
	if rows.Next() {
		err := rows.Scan(&vendorId, &vendorProductReference, &description, &materialTypeId, &length, &width, &height, &weight)
		if err != nil {
			status := util.NewInternalServiceErrorStatus(fmt.Sprintf("failed to read result of query [%v]: %v", sqlQuery, err), "GET_BASIC_VENDOR_PRODUCT_BY_ID_1")
			return nil, &status
		}
	} else {
		status := util.NewNotFoundStatus(fmt.Sprintf("product %v from Vendor %v is not found", vendorProductReference, vendorId))
		return nil, &status
	}

	vendorProduct := &VendorProduct{
		VendorId:        vendorId,
		VendorProductId: vendorProductId,
		Reference:       vendorProductReference,
		Description:     description,
	}

	// try to get material information
	if materialTypeId.Valid {
		id := int(materialTypeId.Int64)
		vendorProduct.MaterialTypeId = &id
	}

	if length.Valid {
		vendorProduct.UnitSize.Length = &common.Measure{
			Value: float32(length.Float64),
			Unit:  "CM",
		}
	}
	if width.Valid {
		vendorProduct.UnitSize.Width = &common.Measure{
			Value: float32(width.Float64),
			Unit:  "CM",
		}
	}

	if height.Valid {
		vendorProduct.UnitSize.Height = &common.Measure{
			Value: float32(height.Float64),
			Unit:  "CM",
		}
	}

	if weight.Valid {
		vendorProduct.UnitSize.GrossWeight = &common.Measure{
			Value: float32(weight.Float64),
			Unit:  "G",
		}
	}

	return vendorProduct, nil
}

func getBasicVendorProductInfo(tx *sql.Tx, vendorId int, vendorProductReference string) (*VendorProduct, *util.Status) {
	// search for vendor_product
	sqlQuery := `
		select VENDOR_PRODUCT_ID, DESCRIPTION, MATERIAL_TYPE_ID, LENGTH, WIDTH, HEIGHT, WEIGHT
		from VENDOR_PRODUCT
		where VENDOR_ID = :vendorId
		  and REFERENCE = :reference
	`
	rows, err := tx.Query(sqlQuery, sql.Named("vendorId", vendorId), sql.Named("reference", vendorProductReference))
	if err != nil {
		status := util.NewInternalServiceErrorStatus(fmt.Sprintf("Failed to run query %v: %v", sqlQuery, err), "GET_VENDOR_PRODUCT_INFO_1")
		return nil, &status
	}
	defer rows.Close()

	var vendorProductId int
	var description string
	var materialTypeId sql.NullInt64
	var length, width, height, weight sql.NullFloat64
	if rows.Next() {
		err := rows.Scan(&vendorProductId, &description, &materialTypeId, &length, &width, &height, &weight)
		if err != nil {
			status := util.NewInternalServiceErrorStatus(fmt.Sprintf("Failed to read result of query [%v]: %v", sqlQuery, err), "GET_VENDOR_PRODUCT_INFO_2")
			return nil, &status
		}
	} else {
		status := util.NewNotFoundStatus(fmt.Sprintf("Product %v from Vendor %v is not found", vendorProductReference, vendorId))
		return nil, &status
	}

	vendorProduct := &VendorProduct{
		VendorId:        vendorId,
		VendorProductId: vendorProductId,
		Reference:       vendorProductReference,
		Description:     description,
	}

	// try to get material information
	if materialTypeId.Valid {
		id := int(materialTypeId.Int64)
		vendorProduct.MaterialTypeId = &id
	}

	if length.Valid {
		vendorProduct.UnitSize.Length = &common.Measure{
			Value: float32(length.Float64),
			Unit:  "CM",
		}
	}
	if width.Valid {
		vendorProduct.UnitSize.Width = &common.Measure{
			Value: float32(width.Float64),
			Unit:  "CM",
		}
	}

	if height.Valid {
		vendorProduct.UnitSize.Height = &common.Measure{
			Value: float32(height.Float64),
			Unit:  "CM",
		}
	}

	if weight.Valid {
		vendorProduct.UnitSize.GrossWeight = &common.Measure{
			Value: float32(weight.Float64),
			Unit:  "G",
		}
	}

	return vendorProduct, nil
}

func getVendorProductPriceHistory(tx *sql.Tx, vendorProductId int) ([]common.Price, *util.Status) {
	sqlQuery := `
		select PRICE, START_DATE, CURRENCY_ID, PRICE_TYPE_ID
		from VENDOR_PRODUCT_PRICE
		where VENDOR_PRODUCT_ID = :vendorProductId
		order by START_DATE desc
	`
	rows, err := tx.Query(sqlQuery, sql.Named("vendorProductId", vendorProductId))
	if err != nil {
		status := util.NewInternalServiceErrorStatus(fmt.Sprintf("failed to execute query [%v]: %v", sqlQuery, err), "GET_VENDOR_PRODUCT_PRICE_HIST_0")
		return nil, &status
	}
	defer rows.Close()

	priceList := []common.Price{}
	for rows.Next() {
		price := common.Price{}
		err := rows.Scan(&price.Amount, &price.StartDate, &price.CurrencyId, &price.TypeId)
		if err != nil {
			status := util.NewInternalServiceErrorStatus(fmt.Sprintf("Failed to parse result of query [%v]: %v", sqlQuery, err), "GET_VENDOR_PRODUCT_PRICE_1")
			return nil, &status
		}
		priceList = append(priceList, price)
	}
	return priceList, nil
}

func getVendorProductMoQHistory(tx *sql.Tx, vendorProductId int) ([]common.MoQ, *util.Status) {
	sqlQuery := `
		select QUANTITY, START_DATE
		from VENDOR_PRODUCT_MOQ
		where VENDOR_PRODUCT_ID = :vendorProductId
		order by START_DATE desc
	`
	rows, err := tx.Query(sqlQuery, sql.Named("vendorProductId", vendorProductId))
	if err != nil {
		status := util.NewInternalServiceErrorStatus(fmt.Sprintf("failed to execute query [%v]: %v", sqlQuery, err), "GET_VENDOR_PRODUCT_MOQ_HIST_0")
		return nil, &status
	}
	defer rows.Close()

	moqList := []common.MoQ{}
	for rows.Next() {
		moq := common.MoQ{}
		err := rows.Scan(&moq.Quantity, &moq.StartDate)
		if err != nil {
			status := util.NewInternalServiceErrorStatus(fmt.Sprintf("Failed to parse result of query [%v]: %v", sqlQuery, err), "GET_VENDOR_PRODUCT_PRICE_1")
			return nil, &status
		}
		moqList = append(moqList, moq)
	}
	return moqList, nil
}

func addNewVendorProduct(tx *sql.Tx, newProduct VendorNewProduct) (int, error) {
	return table.VendorProductInsert(
		tx, newProduct.VendorId, newProduct.Reference,
		newProduct.Description, newProduct.TestPerformed,
		newProduct.ProductTypeId, newProduct.MaterialTypeId,
		newProduct.UnitTypeId, newProduct.Length, newProduct.Width,
		newProduct.Height, newProduct.Weight)
}

func updateVendorProductPrice(
	tx *sql.Tx,
	vendorProductId int,
	priceDate time.Time,
	price float32,
	priceCcyId int,
	priceTypeId int) error {
	rows := table.VendorProductPricePopulateByVendorProductId(tx, vendorProductId)

	if len(rows) == 0 ||
		rows[0].Price != price ||
		rows[0].CurrencyId != priceCcyId ||
		rows[0].PriceTypeId != priceTypeId {
		// we need to insert a new price
		_, err := table.VendorProductPriceInsert(tx, vendorProductId, priceDate, price, priceCcyId, priceTypeId)
		if err != nil {
			return err
		}
	}
	return nil
}

func updateVendorProductMoq(
	tx *sql.Tx,
	vendorProductId int,
	startDate time.Time,
	quantity int) error {
	rows := table.VendorProductMoqPopulateByVendProductId(tx, vendorProductId)
	if len(rows) == 0 ||
		rows[0].Quantity != quantity {
		// we need to insert a new price
		err := table.VendorProductMoqInsert(tx, vendorProductId, quantity, startDate)
		if err != nil {
			return err
		}
	}
	return nil
}

//
// Public methods
//

func AddNewVendorProducts(tx *sql.Tx, newProducts []VendorNewProduct) ([]int, error) {
	result := []int{}
	for _, newProduct := range newProducts {
		vendorProductId, err := addNewVendorProduct(tx, newProduct)
		if err != nil {
			return nil, err
		}
		if newProduct.Price != nil {
			if newProduct.PriceCcyId == nil {
				return nil, fmt.Errorf("product %v of vendor %v has a price but missing currency information",
					newProduct.Reference,
					newProduct.VendorId)
			}
			if newProduct.PriceTypeId == nil {
				return nil, fmt.Errorf("product %v of vendor %v has a price but missing price type information",
					newProduct.Reference,
					newProduct.VendorId)
			}
			err = updateVendorProductPrice(tx, vendorProductId, time.Now(), *newProduct.Price, *newProduct.PriceCcyId, *newProduct.PriceTypeId)
			if err != nil {
				return nil, err
			}
		}
		if newProduct.Moq != nil {
			err = updateVendorProductMoq(tx, vendorProductId, time.Now(), *newProduct.Moq)
			if err != nil {
				return nil, err
			}
		}
		result = append(result, vendorProductId)
	}
	return result, nil
}

func GetVendorProduct(tx *sql.Tx, vendorId int, vendorProductReference string) (*VendorProduct, *util.Status) {
	vp, status := getBasicVendorProductInfo(tx, vendorId, vendorProductReference)
	if status != nil {
		return nil, status
	}

	priceList, err := getVendorProductPriceHistory(tx, vp.VendorProductId)
	if status != nil {
		return nil, err
	}

	if len(priceList) > 0 {
		vp.PriceList = priceList
	}

	return vp, nil
}

func GetVendorProductById(tx *sql.Tx, vendorProductId int) (*VendorProduct, *util.Status) {
	vp, status := getBasicVendorProductInfoById(tx, vendorProductId)
	if status != nil {
		return nil, status
	}

	priceList, status := getVendorProductPriceHistory(tx, vp.VendorProductId)
	if status != nil {
		return nil, status
	}

	moqList, status := getVendorProductMoQHistory(tx, vp.VendorProductId)
	if status != nil {
		return nil, status
	}

	if len(priceList) > 0 {
		vp.PriceList = priceList
	}

	if len(moqList) > 0 {
		vp.MoqList = moqList
	}

	return vp, nil
}
