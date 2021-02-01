package main

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Profiles []Profile `gorm:"many2many:user_profiles;foreignKey:Refer;joinForeignKey:UserReferID;References:UserRefer;JoinReferences:UserRefer"`
	Refer    uint      `gorm:"index:,unique"`
}

type Profile struct {
	gorm.Model
	Name      string
	UserRefer uint `gorm:"index:,unique"`
}

// Which creates join table: user_profiles
//   foreign key: user_refer_id, reference: users.refer
//   foreign key: profile_refer, reference: profiles.user_refer

type UserSystem struct {
	ID       int64     `gorm:"column:us_id;primaryKey;autoIncrement;"`
	Name     string    `gorm:"column:us_name;type:varchar(255);not null;"`
	Branches []*Branch `gorm:"many2many:ucb_user_client_branch;"`
	Clients  []*Client `gorm:"many2many:ucb_user_client_branch;"`
}

type Client struct {
	ID         int64         `gorm:"column:un_id;primaryKey;autoIncrement;"`
	Name       string        `gorm:"column:un_name;type:varchar(255);not null;"`
	UserSystem []*UserSystem `gorm:"many2many:ucb_user_client_branch;"`
	Branches   []Branch
}

type Branch struct {
	ID         int64         `gorm:"column:bn_id;primaryKey;autoIncrement;"`
	Name       string        `gorm:"column:bn_name;type:varchar(255);not null;"`
	UserSystem []*UserSystem `gorm:"many2many:ucb_user_client_branch;"`
	ClientID   uint          `gorm:"column:bn_un_id"`
}

// type UserClientBranch struct {
// 	UserSystemID int64 `gorm:"column:ucb_us_id;"`
// 	ClientID     int64 `gorm:"column:ucb_un_id;"`
// 	BranchID     int64 `gorm:"column:ucb_bn_id;"`
// }

// func (UserClientBranch) TableName() string {
// 	return "ucb_user_client_branch"
// }

func main() {
	db, err := gorm.Open(postgres.Open("host=localhost user=postgres password=postgres dbname=postgres port=5432"), &gorm.Config{})

	if err != nil {
		panic(err.Error())
	}

	// db.Debug().AutoMigrate(&UserClientBranch{})
	db.Debug().AutoMigrate(&UserSystem{})
	db.Debug().AutoMigrate(&Client{})
	db.Debug().AutoMigrate(&Branch{})

	AssociationSeparada(db)
	AssociationUnico(db)
}

func AssociationSeparada(db *gorm.DB) {
	user := UserSystem{}
	clients := []Client{}

	db.First(&user)

	db.Debug().Model(&user).Association("Clients").Find(&clients)
	fmt.Println(user, clients)
}

func AssociationUnico(db *gorm.DB) {
	user := []UserSystem{}

	db.Debug().Preload("Clients").Preload("Branches").Find(&user)

	fmt.Println(user)

	for _, value := range user {
		fmt.Print(value.Name)
		for _, client := range value.Clients {
			fmt.Println(client)
		}

		for _, branch := range value.Branches {
			fmt.Println(branch)
		}
	}
}
