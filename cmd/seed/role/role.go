package role

import (
	"fmt"
	"simple-restful-api-echo-golang-port-adapter-archcore/domains/role"

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
		Id:   role.ADMINISTRATOR_ROLE_ID,
		Name: "Administrator",
		Desc: "Administrator",
	})

	data = append(data, role.Role{
		Id:   role.CUSTOMER_ROLE_ID,
		Name: "Customer",
		Desc: "Customer",
	})
	return data
}
