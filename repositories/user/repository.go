package user

import (
	"database/sql"
	"fmt"
	"net/url"
	domain "simple-restful-api-echo-golang-port-adapter-archcore/domains/user"
	"simple-restful-api-echo-golang-port-adapter-archrepositories/util/querybuilder"
	"simple-restful-api-echo-golang-port-adapter-archutil/logger"
	"time"

	"gorm.io/gorm"
)

const LOG_IDENTIFIER = "REPOSITORY_USER"
const TABLE_NAME = "user"

var log = logger.SetIdentifierField(LOG_IDENTIFIER)

type User struct {
	Id           uint           `gorm:"column:id;primaryKey"`
	RoleId       uint           `gorm:"column:role_id"`
	Email        sql.NullString `gorm:"column:email;size:64;unique"`
	Password     string         `gorm:"column:password;size:255"`
	PasswordSalt string         `gorm:"column:password_salt"`
	PhoneNumber  sql.NullString `gorm:"column:phone_number"`
	LoginAttempt uint8          `gorm:"column:login_attempt;default:0"`
	IsLogin      bool           `gorm:"column:is_login;default:0"`
	IsActive     bool           `gorm:"column:is_active;default:0"`
	Created      time.Time      `gorm:"column:created;autoCreateTime"`
	Modified     time.Time      `gorm:"column:modified;autoUpdateTime"`
}

type Tabler interface {
	TableName() string
}

func (User) TableName() string {
	return TABLE_NAME
}

func Migrate(db gorm.DB) error {
	if !db.Migrator().HasTable(User{}) {
		err := db.Migrator().CreateTable(User{})
		if err != nil {
			log.Error(err)
		}
	}
	indexName := "idx_user_email"
	if !db.Migrator().HasIndex(User{}, indexName) {
		err := db.Exec(fmt.Sprintf("CREATE INDEX %s ON %s (email)", indexName, TABLE_NAME))
		if err != nil {
			log.Error(err)
		}
	}
	indexName = "idx_user_phone_number"
	if !db.Migrator().HasIndex(User{}, indexName) {
		err := db.Exec(fmt.Sprintf("CREATE INDEX %s ON %s (phone_number)", indexName, TABLE_NAME))
		if err != nil {
			log.Error(err)
		}
	}
	indexName = "idx_user_is_active"
	if !db.Migrator().HasIndex(User{}, indexName) {
		err := db.Exec(fmt.Sprintf("CREATE INDEX %s ON %s (is_active)", indexName, TABLE_NAME))
		if err != nil {
			log.Error(err)
		}
	}
	constraintName := "foreign_user_role_id"
	if !db.Migrator().HasConstraint(User{}, constraintName) {
		err := db.Exec(fmt.Sprintf("ALTER TABLE `%s` ADD CONSTRAINT `%s` FOREIGN KEY (`%s`) REFERENCES `%s` (`%s`) ON DELETE RESTRICT ON UPDATE RESTRICT", TABLE_NAME, constraintName, "role_id", "role", "id"))
		if err != nil {
			log.Error(err)
		}
	}
	return nil
}

func NewData(data domain.User) *User {
	return &User{
		Id:           data.Id,
		RoleId:       data.RoleId,
		Email:        sql.NullString{String: data.Email},
		Password:     data.Password,
		PasswordSalt: data.PasswordSalt,
		PhoneNumber:  sql.NullString{String: data.PhoneNumber},
		LoginAttempt: data.LoginAttempt,
		IsLogin:      data.IsLogin,
		IsActive:     data.IsActive,
		Created:      data.Created,
		Modified:     data.Modified,
	}
}

func (user *User) Map() domain.User {
	var data domain.User
	data.Id = user.Id
	data.RoleId = user.RoleId
	data.Email = user.Email.String
	data.Password = user.Password
	data.PasswordSalt = user.PasswordSalt
	data.PhoneNumber = user.PhoneNumber.String
	data.LoginAttempt = user.LoginAttempt
	data.IsLogin = user.IsLogin
	data.IsActive = user.IsActive
	data.Created = user.Created
	data.Modified = user.Modified
	return data
}

type Repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) (*Repository, error) {
	repository := Repository{db}
	return &repository, nil
}

type IRepository interface {
	List(filters url.Values, limit int, offset int) ([]*domain.User, error)
	Create(data domain.User) (*domain.User, error)
	FindById(id string) (*domain.User, error)
	UpdateById(id string, data domain.User) (*domain.User, error)
	DeleteById(id string) error
}

func (repository Repository) List(filters url.Values, limit int, offset int) ([]*domain.User, error) {
	users := []*User{}
	querybuilder.GormFilterBuilder(repository.db, filters, limit, offset).Find(&users)
	result := []*domain.User{}
	for _, data := range users {
		newData := data.Map()
		result = append(result, &newData)
	}
	return result, nil
}

func (repository Repository) Create(data domain.User) (*domain.User, error) {
	user := NewData(data)
	inserted := repository.db.Create(&user)
	if inserted.RowsAffected == 0 {
		return nil, inserted.Error
	}
	result := user.Map()
	return &result, nil
}

func (repository Repository) FindById(id string) (*domain.User, error) {
	user := User{}
	finded := repository.db.Find(&user, id)
	if finded.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	result := user.Map()
	return &result, nil
}

func (repository Repository) UpdateById(id string, data domain.User) (*domain.User, error) {
	user := User{}
	finded := repository.db.Find(&user, id)
	if finded.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	user.RoleId = data.RoleId
	user.Email = sql.NullString{String: data.Email}
	user.PhoneNumber = sql.NullString{String: data.PhoneNumber}
	user.LoginAttempt = data.LoginAttempt
	user.IsLogin = data.IsLogin
	user.IsActive = data.IsActive
	err := repository.db.Save(&user)
	if err.Error != nil {
		return nil, err.Error
	}
	result := user.Map()
	return &result, nil
}

func (repository Repository) DeleteById(id string) error {
	user := User{}
	finded := repository.db.Find(&user, id)
	if finded.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	repository.db.Delete(&user)
	return nil
}
