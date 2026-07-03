package controllers

import (
	"net/http"
	"ppdb-be/models"
	"ppdb-be/services"
	"ppdb-be/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type SchoolProfileController interface {
	CreateSchoolProfile(c *gin.Context)
	GetSchoolProfiles(c *gin.Context)
	GetSchoolProfileByID(c *gin.Context)
	UpdateSchoolProfile(c *gin.Context)
	DeleteSchoolProfile(c *gin.Context)
}

type schoolProfileController struct {
	profileService services.SchoolProfileService
}

func NewSchoolProfileController(profileService services.SchoolProfileService) SchoolProfileController {
	return &schoolProfileController{profileService: profileService}
}

// @Summary Create a new School Profile
// @Description Create a new School Profile. Requires authorization.
// @Tags School Profiles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body models.SchoolProfileRequest true "School Profile Request Payload"
// @Success 201 {object} utils.Response "Created"
// @Failure 400 {object} utils.ErrorResponse "Bad Request"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized"
// @Failure 500 {object} utils.ErrorResponse "Internal Server Error"
// @Router /school-profiles [post]
func (ctrl *schoolProfileController) CreateSchoolProfile(c *gin.Context) {
	var req models.SchoolProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorMsg(c, http.StatusBadRequest, utils.MsgErrBadRequest, err)
		return
	}

	userID := c.GetInt64("userID")
	if err := ctrl.profileService.CreateSchoolProfile(req, userID); err != nil {
		utils.SendErrorMsg(c, http.StatusInternalServerError, utils.MsgErrInternal, err)
		return
	}

	utils.SendSuccessMsg(c, http.StatusCreated, utils.MsgSuccessCreate, nil)
}

// @Summary Get all School Profiles
// @Description Get a list of all School Profiles.
// @Tags School Profiles
// @Produce json
// @Success 200 {object} utils.Response{data=[]models.SchoolProfile} "Success"
// @Failure 500 {object} utils.ErrorResponse "Internal Server Error"
// @Router /school-profiles [get]
func (ctrl *schoolProfileController) GetSchoolProfiles(c *gin.Context) {
	profiles, err := ctrl.profileService.GetSchoolProfiles()
	if err != nil {
		utils.SendErrorMsg(c, http.StatusInternalServerError, utils.MsgErrInternal, err)
		return
	}

	utils.SendSuccessMsg(c, http.StatusOK, utils.MsgSuccessFetch, profiles)
}

// @Summary Get School Profile by ID
// @Description Get a School Profile by its ID.
// @Tags School Profiles
// @Produce json
// @Param id path int true "School Profile ID"
// @Success 200 {object} utils.Response{data=models.SchoolProfile} "Success"
// @Failure 400 {object} utils.ErrorResponse "Bad Request"
// @Failure 404 {object} utils.ErrorResponse "Not Found"
// @Router /school-profiles/{id} [get]
func (ctrl *schoolProfileController) GetSchoolProfileByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendErrorMsg(c, http.StatusBadRequest, utils.MsgErrInvalidID, err)
		return
	}

	profile, err := ctrl.profileService.GetSchoolProfileByID(id)
	if err != nil {
		utils.SendErrorMsg(c, http.StatusNotFound, utils.MsgErrNotFound, err)
		return
	}

	utils.SendSuccessMsg(c, http.StatusOK, utils.MsgSuccessFetch, profile)
}

// @Summary Update School Profile
// @Description Update an existing School Profile by its ID. Requires authorization.
// @Tags School Profiles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "School Profile ID"
// @Param body body models.SchoolProfileRequest true "School Profile Request Payload"
// @Success 200 {object} utils.Response "Success"
// @Failure 400 {object} utils.ErrorResponse "Bad Request"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized"
// @Failure 500 {object} utils.ErrorResponse "Internal Server Error"
// @Router /school-profiles/{id} [put]
func (ctrl *schoolProfileController) UpdateSchoolProfile(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendErrorMsg(c, http.StatusBadRequest, utils.MsgErrInvalidID, err)
		return
	}

	var req models.SchoolProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorMsg(c, http.StatusBadRequest, utils.MsgErrBadRequest, err)
		return
	}

	userID := c.GetInt64("userID")
	if err := ctrl.profileService.UpdateSchoolProfile(id, req, userID); err != nil {
		utils.SendErrorMsg(c, http.StatusInternalServerError, utils.MsgErrInternal, err)
		return
	}

	utils.SendSuccessMsg(c, http.StatusOK, utils.MsgSuccessUpdate, nil)
}

// @Summary Delete School Profile
// @Description Delete a School Profile by its ID. Requires authorization.
// @Tags School Profiles
// @Produce json
// @Security BearerAuth
// @Param id path int true "School Profile ID"
// @Success 200 {object} utils.Response "Success"
// @Failure 400 {object} utils.ErrorResponse "Bad Request"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized"
// @Failure 500 {object} utils.ErrorResponse "Internal Server Error"
// @Router /school-profiles/{id} [delete]
func (ctrl *schoolProfileController) DeleteSchoolProfile(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendErrorMsg(c, http.StatusBadRequest, utils.MsgErrInvalidID, err)
		return
	}

	if err := ctrl.profileService.DeleteSchoolProfile(id); err != nil {
		utils.SendErrorMsg(c, http.StatusInternalServerError, utils.MsgErrInternal, err)
		return
	}

	utils.SendSuccessMsg(c, http.StatusOK, utils.MsgSuccessDelete, nil)
}
