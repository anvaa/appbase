package app_ctrl

import (
	"app/app_conf"
	"app/app_db"
	"app/app_models"
	"fmt"
	"strings"
	"strconv"

	"net/http"

	"github.com/gin-gonic/gin"
)

// Note_View displays a specific note
func Note_View(c *gin.Context) {
	// Implementation here

	_sort := app_conf.GetInt("note_sort_type")

	notes, err := app_db.GetAllNotes(c.Param("id"))
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": err.Error(),
		})
		return
	}

	// Render the note view template
	c.HTML(http.StatusOK, "note_view.html", gin.H{
		"js": "note.js",

		"projid": c.Param("id"),
		"typ":    app_db.Typ_GetAllTypsub(3),
		"sort":   _sort,
		"notes":  notes,
	})
}

// Note_Edit displays the edit page for a specific note
func Note_Edit(c *gin.Context) {
	// Implementation here
	note, err := app_db.GetNoteByUUID(c.Param("id"))
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": err.Error(),
		})
		return
	}

	// Render the note edit template
	c.HTML(http.StatusOK, "note_edit.html", gin.H{
		"js": "note.js",

		"note": note,
		"typ":  app_db.Typ_GetAllTypsub(3),
	})
}

// Note_AddUpd adds or updates a note
func Note_AddUpd(c *gin.Context) {
	var body struct {
		Projid  int    `json:"projid"`
		UUID    int    `json:"uuid"`
		Type    int    `json:"type"`
		Content string `json:"content"`
		Header  string `json:"header"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	_notes := app_models.Notes{
		TypsubID: body.Type,
		Content:  strings.TrimSpace(body.Content),
		Header:   strings.TrimSpace(body.Header),
	}

	url := fmt.Sprintf("/note/%d", body.Projid)

	// Implementation here
	if body.UUID > 0 {
		// Update existing note
		_notes.UpdatedbyID = CurUserID // Set the updater as well
		if err := app_db.AppDB.Model(&app_models.Notes{}).Where("uuid = ?", body.UUID).Updates(&_notes).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else {
		// Create new note
		_notes.CreatedbyID = CurUserID
		_notes.UpdatedbyID = CurUserID // Set the updater as well
		_notes.DeletedbyID = 1

		_notes.ProjectID = body.Projid

		if err := app_db.AppDB.Model(&app_models.Notes{}).Create(&_notes).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		} else {
			url = fmt.Sprintf("/note/edit/%d", _notes.UUID)
		}
	}

	fmt.Println("Note saved or added ", url)
	c.JSON(http.StatusOK, gin.H{
		"message":  "success",
		"redirect": url,
	})
}

// Note_Delete deletes a specific note
func Note_Delete(c *gin.Context) {
	// Implementation here
	noteUUID := c.Param("id")
	if noteUUID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Note UUID is required"})
		return
	}

	note, err := app_db.GetNoteByUUID(noteUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	projid := note.ProjectID
	note.DeletedbyID = CurUserID // Set the deleter as the current user

	if err := app_db.AppDB.Delete(&note).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	url := fmt.Sprintf("/note/%d", projid)
	c.JSON(http.StatusOK, gin.H{
		"message":  "success",
		"redirect": url,
	})

}

func Note_SaveSort(c *gin.Context) {

	val := c.Param("id")

	// set val to int
	_val, err := strconv.Atoi(val)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sort value"})
		return
	}

	app_conf.SetVal("note_sort_type", _val)

	c.JSON(http.StatusOK, gin.H{"message": "success"})
}