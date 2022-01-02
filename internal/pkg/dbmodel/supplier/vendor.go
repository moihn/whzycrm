package supplier

import (
	"database/sql"
	"fmt"

	"github.com/moihn/whzycrmgo/internal/pkg/util"
)

type Vendor struct {
	Id   int
	Name string
}

func GetAll(tx *sql.Tx) ([]Vendor, *util.Status) {
	sqlQuery := `
		select ID, NAME
		from VENDOR
	`
	rows, err := tx.Query(sqlQuery)
	if err != nil {
		status := util.NewInternalServiceErrorStatus(fmt.Sprintf("Failed to run query %v: %v", sqlQuery, err), "GET_VENDORS")
		return nil, &status
	}
	defer rows.Close()

	var vendorId int
	var vendorName string

	var vendors []Vendor
	for rows.Next() {
		err := rows.Scan(&vendorId, &vendorName)
		if err != nil {
			status := util.NewInternalServiceErrorStatus(fmt.Sprintf("Failed to read result of query [%v]: %v", sqlQuery, err), "GET_VENDORS")
			return nil, &status
		}
		vendors = append(vendors, Vendor{
			Id:   vendorId,
			Name: vendorName,
		})
	}
	return vendors, nil
}
