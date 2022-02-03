package user

import (
	"database/sql"
	"fmt"
	"golang-vscode-setup/repository/util/querybuilder"
	"golang-vscode-setup/service/user/defined"
	"golang-vscode-setup/util/logger"
	"net/url"
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

func NewUser(user defined.User) *User {
	return &User{
		Id:           user.Id,
		RoleId:       user.RoleId,
		Email:        user.Email,
		Password:     user.Password,
		PasswordSalt: user.PasswordSalt,
		PhoneNumber:  user.PhoneNumber,
		LoginAttempt: user.LoginAttempt,
		IsLogin:      user.IsLogin,
		IsActive:     user.IsActive,
		Created:      user.Created,
		Modified:     user.Modified,
	}
}

func (user *User) Map() defined.User {
	var data defined.User
	data.Id = user.Id
	data.RoleId = user.RoleId
	data.Email = user.Email
	data.Password = user.Password
	data.PasswordSalt = user.PasswordSalt
	data.PhoneNumber = user.PhoneNumber
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

func NewRepository(db *gorm.DB) (*Repository, error) {
	repository := Repository{db}
	return &repository, nil
}

type IRepository interface {
	List(filters url.Values, limit int, offset int) ([]*defined.User, error)
	Create(user defined.User) (*defined.User, error)
	FindById(id string) (*defined.User, error)
	UpdateById(id string, user defined.User) (*defined.User, error)
	DeleteById(id string) error
}

func (repository Repository) List(filters url.Values, limit int, offset int) ([]*defined.User, error) {
	users := []*User{}
	querybuilder.GormFilterBuilder(repository.db, filters, limit, offset).Find(&users)
	result := []*defined.User{}
	for _, data := range users {
		newData := data.Map()
		result = append(result, &newData)
	}
	return result, nil
}

func (repository Repository) Create(user defined.User) (*defined.User, error) {
	newUser := NewUser(user)
	inserted := repository.db.Create(&newUser)
	if inserted.RowsAffected == 0 {
		return nil, inserted.Error
	}
	result := newUser.Map()
	return &result, nil
}

func (repository Repository) FindById(id string) (*defined.User, error) {
	user := User{}
	data := repository.db.Find(&user, id)
	if data.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	result := user.Map()
	return &result, nil
}

func (repository Repository) UpdateById(id string, user defined.User) (*defined.User, error) {
	newUser := User{}
	data := repository.db.Find(&newUser, id)
	if data.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	newUser.RoleId = user.RoleId
	newUser.Email = user.Email
	newUser.PhoneNumber = user.PhoneNumber
	newUser.LoginAttempt = user.LoginAttempt
	newUser.IsLogin = user.IsLogin
	newUser.IsActive = user.IsActive
	err := repository.db.Save(&newUser)
	if err.Error != nil {
		return nil, err.Error
	}
	result := newUser.Map()
	return &result, nil
}

func (repository Repository) DeleteById(id string) error {
	newUser := User{}
	data := repository.db.Find(&newUser, id)
	if data.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	repository.db.Delete(&newUser)
	return nil
}
