package role

import (
	"database/sql"
	"fmt"
	"golang-vscode-setup/repository/util/querybuilder"
	"golang-vscode-setup/service/role/defined"
	"golang-vscode-setup/util/logger"
	"net/url"
	"time"

	"gorm.io/gorm"
)

const LOG_IDENTIFIER = "REPOSITORY_ROLE"
const TABLE_NAME = "role"

var log = logger.SetIdentifierField(LOG_IDENTIFIER)

type Role struct {
	Id       uint           `gorm:"column:id;primaryKey"`
	Name     string         `gorm:"column:name;size:32;unique"`
	Desc     sql.NullString `gorm:"column:desc;size:255"`
	Created  time.Time      `gorm:"column:created;autoCreateTime"`
	Modified time.Time      `gorm:"column:modified;autoUpdateTime"`
}

type Tabler interface {
	TableName() string
}

func (Role) TableName() string {
	return TABLE_NAME
}

func Migrate(db gorm.DB) error {
	if !db.Migrator().HasTable(Role{}) {
		err := db.Migrator().CreateTable(Role{})
		if err != nil {
			log.Error(err)
		}
	}
	indexName := "idx_role_name"
	if !db.Migrator().HasIndex(Role{}, indexName) {
		err := db.Exec(fmt.Sprintf("CREATE INDEX %s ON %s (name)", indexName, TABLE_NAME))
		if err != nil {
			log.Error(err)
		}
	}
	return nil
}

func NewRole(role defined.Role) *Role {
	return &Role{
		Id:       role.Id,
		Name:     role.Name,
		Desc:     role.Desc,
		Created:  role.Created,
		Modified: role.Modified,
	}
}

func (role *Role) Map() defined.Role {
	var data defined.Role
	data.Id = role.Id
	data.Name = role.Name
	data.Desc = role.Desc
	data.Created = role.Created
	data.Modified = role.Modified
	return data
}

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) (*Repository, error) {
	repository := Repository{db}
	return &repository, nil
}

type IRepository interface {
	List(filters url.Values, limit int, offset int) ([]*defined.Role, error)
	Create(role defined.Role) (*defined.Role, error)
	FindById(id string) (*defined.Role, error)
	UpdateById(id string, role defined.Role) (*defined.Role, error)
	DeleteById(id string) error
}

func (repository Repository) List(filters url.Values, limit int, offset int) ([]*defined.Role, error) {
	roles := []*Role{}
	querybuilder.GormFilterBuilder(repository.db, filters, limit, offset).Find(&roles)
	result := []*defined.Role{}
	for _, data := range roles {
		newData := data.Map()
		result = append(result, &newData)
	}
	return result, nil
}

func (repository Repository) Create(role defined.Role) (*defined.Role, error) {
	newRole := NewRole(role)
	inserted := repository.db.Create(&newRole)
	if inserted.RowsAffected == 0 {
		return nil, inserted.Error
	}
	result := newRole.Map()
	return &result, nil
}

func (repository Repository) FindById(id string) (*defined.Role, error) {
	role := Role{}
	data := repository.db.Find(&role, id)
	if data.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	result := role.Map()
	return &result, nil
}

func (repository Repository) UpdateById(id string, role defined.Role) (*defined.Role, error) {
	newRole := Role{}
	data := repository.db.Find(&newRole, id)
	if data.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	newRole.Name = role.Name
	newRole.Desc = role.Desc
	err := repository.db.Save(&newRole)
	if err.Error != nil {
		return nil, err.Error
	}
	result := newRole.Map()
	return &result, nil
}

func (repository Repository) DeleteById(id string) error {
	newRole := Role{}
	data := repository.db.Find(&newRole, id)
	if data.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	repository.db.Delete(&newRole)
	return nil
}
