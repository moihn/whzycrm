package client

import (
	"database/sql"

	"github.com/moihn/whzycrmgo/internal/pkg/dbmodel/table"
)

type ClientModel struct {
	Id        int
	Name      string
	CountryId int
}

func GetAllClientModels(tx *sql.Tx) []ClientModel {
	var result []ClientModel
	for _, row := range table.ClientPopulateByPk(tx) {
		result = append(result, ClientModel{
			Id:        row.Id,
			Name:      row.Name,
			CountryId: row.CountryId,
		})
	}
	return result
}
