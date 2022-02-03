package role

import (
	"database/sql"
	"fmt"
	"golang-vscode-setup/repository/role"
	"golang-vscode-setup/repository/role/defined"

	"gorm.io/gorm"
)

func Seed(db *gorm.DB) {
	data := data()
	for _, value := range data {
		result := db.Create(&value)
		if result.RowsAffected > 0 {
			message := fmt.Sprintf("inserted role with id %d", value.Id)
			fmt.Println(message)
		}
	}
}

func data() []role.Role {
	var data []role.Role

	data = append(data, role.Role{
		Id:   defined.ADMINISTRATOR_ROLE_ID,
		Name: "Administrator",
		Desc: sql.NullString{String: "Administrator"},
	})

	data = append(data, role.Role{
		Id:   defined.CUSTOMER_ROLE_ID,
		Name: "Customer",
		Desc: sql.NullString{String: "Customer"},
	})
	return data
}
