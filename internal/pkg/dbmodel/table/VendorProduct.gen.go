package table

import (
	"database/sql"
	"fmt"

	"github.com/sirupsen/logrus"
)

type VendorProductRow struct {
	VendorProductId int
	Reference       string
	TestPerformed   bool
	VendorId        int
	Description     string
	MaterialTypeId  *int
	ProductTypeId   *int
	UnitTypeId      *int
	Length          *float32
	Width           *float32
	Height          *float32
	Weight          *float32
}

func VendorProductInsert(
	tx *sql.Tx,
	vendorId int,
	reference string,
	description string,
	testPerformed bool,
	productTypeId int,
	materialTypeId *int,
	unitTypeId *int,
	length *float32,
	width *float32,
	height *float32,
	weight *float32,
) (int, error) {
	sqlQuery := `
	INSERT INTO VENDOR_PRODUCT
	(VENDOR_ID, REFERENCE, DESCRIPTION, PRODUCT_TYPE_ID, MATERIAL_TYPE_ID, UNIT_TYPE_ID, LENGTH, WIDTH, HEIGHT, WEIGHT)
	VALUES
	(:vendorId, :reference, :description, :productTypeId, :materialTypeId, :unitTypeId, :length, :width, :height, :weight)
	RETURNING VENDOR_PRODUCT_ID into :vendorProductId
	`
	var vendorProductId int
	if result, err := tx.Exec(sqlQuery,
		sql.Named("vendorId", vendorId),
		sql.Named("reference", reference),
		sql.Named("description", description),
		sql.Named("productTypeId", productTypeId),
		sql.Named("materialTypeId", materialTypeId),
		sql.Named("unitTypeId", unitTypeId),
		sql.Named("length", length),
		sql.Named("width", width),
		sql.Named("height", height),
		sql.Named("weight", weight),
		sql.Named("vendorProductId", sql.Out{Dest: &vendorProductId})); err != nil {
		return 0, fmt.Errorf("failed to execute query [%v]: %v", sqlQuery, err)
	} else {
		nRows, err := result.RowsAffected()
		if err != nil {
			logrus.Fatalf("failed to get number of inserted row for query [%v]: %v", sqlQuery, err)
		}
		if nRows == 0 {
			logrus.Fatalf("failed to insert vendor_product row [%v]: vendorId=%v, reference=%v, description=%v", sqlQuery,
				vendorId, reference, description)
		}
	}
	return vendorProductId, nil
}
