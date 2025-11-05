package hosteloffice

import (
	"net/http"

	"github.com/LambdaIITH/mess_registration/models"
	"github.com/LambdaIITH/mess_registration/utils"
	"github.com/gin-gonic/gin"
)

func (oc *OfficeController) AddNewUser(c *gin.Context) {
	var input struct {
		Name   string  `json:"name" binding:"required,min=2,max=100"`
		RollNo string  `json:"roll_no" binding:"required"`
		Email  string  `json:"email" binding:"required"`
		Type   int8    `json:"user_type"`
		Phone  *string `json:"phone"`                            // optional
		Mess   int8    `json:"mess" binding:"oneof=0 1 2 3 4 5"` // optional, but validated
	}

	// Bind JSON
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}
	canRegister := false
	if input.Type == 0 {
		canRegister = true
	}

	// Create user model
	user := models.User{
		Name:        input.Name,
		RollNo:      input.RollNo,
		Type:        input.Type,
		Phone:       input.Phone,
		Mess:        input.Mess,
		Email:       input.Email,
		NextMess:    0,           // default unassigned
		CanRegister: canRegister, // default
	}

	// Insert into DB
	if err := oc.DB.Create(&user).Error; err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Could not create user: "+err.Error())
		return
	}

	utils.RespondWithJSON(c, http.StatusOK, models.APIResponse{
		Message: "User added successfully",
		Data:    user,
	})
}
