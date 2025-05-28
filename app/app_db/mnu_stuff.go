package app_db

import (
	"app/app_models"

	"log"
	"sort"
)

func Upd_MenuTitle(mnu_id int, value string) error {
	// Update the menu title
	var menu app_models.Menu
	err := AppDB.Model(&menu).Where("uuid = ?", mnu_id).Update("title", value).Error
	if err != nil {
		log.Println("Error updating menu title:", err)
		return err
	}
	return nil

}

func Get_MenuTitles() []app_models.Menu {
	// Get all menu titles and theis submenu items. Order by submenu name ascending
	var menu []app_models.Menu
	err := AppDB.Preload("SubItems").Find(&menu).Error
	if err != nil {
		log.Println("Error getting menu titles:", err)
		return nil
	}

	return sortMenusTitles(menu)
}

func sortMenusTitles(menu []app_models.Menu) []app_models.Menu {
	// Sort each menu's submenu by name ascending
	for i := range menu {
		sort.Slice(menu[i].SubItems, func(i2, j int) bool {
			return menu[i].SubItems[i2].Name < menu[i].SubItems[j].Name
		})
	}
	return menu
}

func Upd_MenuSubItems(mnuid int, name string) error {
	// Update the menu item title
	menu := app_models.SubMenu{}
	err := AppDB.Model(&menu).Where("id = ?", mnuid).Update("name", name).Error
	if err != nil {
		log.Println("Error updating menu item title:", err)
		return err
	}
	return nil
}

func Mnu_GetSubItemIdByType(mnuid, subtype any) int {
	// Get the submenu ID by type to get defaults
	var subitem app_models.SubMenu
	err := AppDB.Where("menu_id = ? AND type = ?", mnuid, subtype).First(&subitem).Error
	if err != nil {
		log.Println("Error getting submenu ID by type:", err)
		return 0
	}
	return subitem.ID
}

func Mnu_GetMenuTitle(id any) string {
	var title app_models.Menu
	AppDB.Where("id = ?", id).First(&title)

	if title.ID == 0 {
		return "nil"
	}

	return title.Title
}

func Sub_GetName(id any) string {
	var sub app_models.SubMenu
	AppDB.Where("id = ?", id).First(&sub)

	if sub.ID == 0 {
		return "nil"
	}

	return sub.Name
}