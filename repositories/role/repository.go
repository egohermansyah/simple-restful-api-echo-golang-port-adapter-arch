package role

import (
	"database/sql"
	"fmt"
	"net/url"
	domain "simple-restful-api-echo-golang-port-adapter-archcore/domains/role"
	"simple-restful-api-echo-golang-port-adapter-archrepositories/util/querybuilder"
	"simple-restful-api-echo-golang-port-adapter-archutil/logger"
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

func NewData(data domain.Role) *Role {
	return &Role{
		Id:       data.Id,
		Name:     data.Name,
		Desc:     sql.NullString{String: data.Desc},
		Created:  data.Created,
		Modified: data.Modified,
	}
}

func (data *Role) Map() domain.Role {
	var result domain.Role
	result.Id = data.Id
	result.Name = data.Name
	result.Desc = data.Desc.String
	result.Created = data.Created
	result.Modified = data.Modified
	return result
}

type Repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) (*Repository, error) {
	repository := Repository{db}
	return &repository, nil
}

func (repository Repository) List(filters url.Values, limit int, offset int) ([]*domain.Role, error) {
	roles := []*Role{}
	querybuilder.GormFilterBuilder(repository.db, filters, limit, offset).Find(&roles)
	result := []*domain.Role{}
	for _, data := range roles {
		newData := data.Map()
		result = append(result, &newData)
	}
	return result, nil
}

func (repository Repository) Create(data domain.Role) (*domain.Role, error) {
	role := NewData(data)
	inserted := repository.db.Create(&role)
	if inserted.RowsAffected == 0 {
		return nil, inserted.Error
	}
	result := role.Map()
	return &result, nil
}

func (repository Repository) FindById(id string) (*domain.Role, error) {
	role := Role{}
	finded := repository.db.Find(&role, id)
	if finded.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	result := role.Map()
	return &result, nil
}

func (repository Repository) UpdateById(id string, data domain.Role) (*domain.Role, error) {
	role := Role{}
	finded := repository.db.Find(&role, id)
	if finded.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	role.Name = data.Name
	role.Desc = sql.NullString{String: data.Desc}
	err := repository.db.Save(&role)
	if err.Error != nil {
		return nil, err.Error
	}
	result := role.Map()
	return &result, nil
}

func (repository Repository) DeleteById(id string) error {
	role := Role{}
	finded := repository.db.Find(&role, id)
	if finded.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	repository.db.Delete(&role)
	return nil
}
