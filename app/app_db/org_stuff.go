package app_db

import (
	"app/app_models"
	"gorm.io/gorm"
	"fmt"
)

// Org_GetAll returns all orgs with preloaded users, grouped and ordered by name.
func Org_GetAll() ([]app_models.Org, error) {
	var orgs []app_models.Org
	err := AppDB.
		Preload("Users").
		Find(&orgs).Error
	if err != nil {
		return nil, err
	}
	return orgs, nil
}

// Org_GetByUUID fetches an org by its UUID and preloads users.
func Org_GetByUUID(uuid string) (*app_models.Org, error) {
	var org app_models.Org
	err := AppDB.Preload("Users.AuthLevel").First(&org, "uuid = ?", uuid).Error
	if err != nil {
		// handle error
}
	return &org, nil
}

// Org_GetByUser returns all orgs associated with a given user ID.
func Org_GetByUser(userID uint) ([]app_models.Org, error) {
	var user app_models.Users
	if err := AppDB.Preload("Org").First(&user, userID).Error; err != nil {
		return nil, err
	}
	orgs := make([]app_models.Org, 0, len(user.Org))
	for _, o := range user.Org {
		var org app_models.Org
		if err := AppDB.First(&org, o).Error; err == nil {
			orgs = append(orgs, org)
		}
	}
	return orgs, nil
}

// Org_Delete deletes an org by UUID if it exists.
func Org_Delete(uuid string) error {
	var org app_models.Org
	if err := AppDB.Where("uuid = ?", uuid).First(&org).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		return err
	}
	return AppDB.Delete(&org).Error
}

// Org_AddUserto adds a user to an org (many2many relationship).
func Org_AddMember(orgID, userID int) error {
	var org app_models.Org
	if err := AppDB.First(&org, orgID).Error; err != nil {
		return err
	}
	var user app_models.Users
	if err := AppDB.First(&user, userID).Error; err != nil {
		return err
	}
	return AppDB.Model(&org).Association("Users").Append(&user)
}

func Org_RemoveMember(orgID, userID int) error {
	var org app_models.Org
	if err := AppDB.First(&org, orgID).Error; err != nil {
		return err
	}
	var user app_models.Users
	if err := AppDB.First(&user, userID).Error; err != nil {
		return err
	}
	return AppDB.Model(&org).Association("Users").Delete(&user)
}

// Org_AddUpd creates a new org if uuid == "0", otherwise updates the existing org.
func Org_AddUpd(uuid, name, note string) error {
	if uuid == "0" {
		// append user id 1 admin
		var users []*app_models.Users
		if err := AppDB.Find(&users, 1).Error; err != nil {
			return err
		}

		org := app_models.Org{
			Name:  name,
			Note:  note,
			Users: users,
		}
		return AppDB.Create(&org).Error
	}
	var exOrg app_models.Org
	if err := AppDB.Preload("Users").Where("uuid = ?", uuid).First(&exOrg).Error; err != nil {
		return err
	}
	exOrg.Name = name
	exOrg.Note = note
	fmt.Println("UUID", uuid, exOrg.Name, "update existing")
	return AppDB.Model(&exOrg).Updates(map[string]interface{}{
		"Name": exOrg.Name,
		"Note": exOrg.Note,
	}).Error
}