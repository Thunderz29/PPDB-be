package controllers

import (
	"net/http"
	"ppdb-be/models"
	"ppdb-be/services"
	"ppdb-be/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MasterSchoolController interface {
	CreateSchool(c *gin.Context)
	GetSchools(c *gin.Context)
	GetSchoolByID(c *gin.Context)
	UpdateSchool(c *gin.Context)
	DeleteSchool(c *gin.Context)
}

type masterSchoolController struct {
	schoolService services.MasterSchoolService
}

func NewMasterSchoolController(schoolService services.MasterSchoolService) MasterSchoolController {
	return &masterSchoolController{schoolService: schoolService}
}

// @Summary Create a new Master School (Location)
// @Description Create a new Master School (Location). Requires authorization.
// @Tags School Locations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body models.MasterSchoolRequest true "Master School Request Payload"
// @Success 201 {object} utils.Response "Created"
// @Failure 400 {object} utils.ErrorResponse "Bad Request"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized"
// @Failure 500 {object} utils.ErrorResponse "Internal Server Error"
// @Router /school-locations [post]
func (ctrl *masterSchoolController) CreateSchool(c *gin.Context) {
	var req models.MasterSchoolRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorMsg(c, http.StatusBadRequest, utils.MsgErrBadRequest, err)
		return
	}

	userID := c.GetInt64("userID")
	if err := ctrl.schoolService.CreateSchool(req, userID); err != nil {
		utils.SendErrorMsg(c, http.StatusInternalServerError, utils.MsgErrInternal, err)
		return
	}

	utils.SendSuccessMsg(c, http.StatusCreated, utils.MsgSuccessCreate, nil)
}

// @Summary Get all Master Schools (Locations)
// @Description Get a list of all Master Schools (Locations).
// @Tags School Locations
// @Produce json
// @Success 200 {object} utils.Response{data=[]models.MasterSchool} "Success"
// @Failure 500 {object} utils.ErrorResponse "Internal Server Error"
// @Router /school-locations [get]
func (ctrl *masterSchoolController) GetSchools(c *gin.Context) {
	schools, err := ctrl.schoolService.GetSchools()
	if err != nil {
		utils.SendErrorMsg(c, http.StatusInternalServerError, utils.MsgErrInternal, err)
		return
	}

	utils.SendSuccessMsg(c, http.StatusOK, utils.MsgSuccessFetch, schools)
}

// @Summary Get Master School (Location) by ID
// @Description Get a Master School (Location) by its ID.
// @Tags School Locations
// @Produce json
// @Param id path int true "Master School ID"
// @Success 200 {object} utils.Response{data=models.MasterSchool} "Success"
// @Failure 400 {object} utils.ErrorResponse "Bad Request"
// @Failure 404 {object} utils.ErrorResponse "Not Found"
// @Router /school-locations/{id} [get]
func (ctrl *masterSchoolController) GetSchoolByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendErrorMsg(c, http.StatusBadRequest, utils.MsgErrInvalidID, err)
		return
	}

	school, err := ctrl.schoolService.GetSchoolByID(id)
	if err != nil {
		utils.SendErrorMsg(c, http.StatusNotFound, utils.MsgErrNotFound, err)
		return
	}

	utils.SendSuccessMsg(c, http.StatusOK, utils.MsgSuccessFetch, school)
}

// @Summary Update Master School (Location)
// @Description Update an existing Master School (Location) by its ID. Requires authorization.
// @Tags School Locations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Master School ID"
// @Param body body models.MasterSchoolRequest true "Master School Request Payload"
// @Success 200 {object} utils.Response "Success"
// @Failure 400 {object} utils.ErrorResponse "Bad Request"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized"
// @Failure 500 {object} utils.ErrorResponse "Internal Server Error"
// @Router /school-locations/{id} [put]
func (ctrl *masterSchoolController) UpdateSchool(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendErrorMsg(c, http.StatusBadRequest, utils.MsgErrInvalidID, err)
		return
	}

	var req models.MasterSchoolRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorMsg(c, http.StatusBadRequest, utils.MsgErrBadRequest, err)
		return
	}

	userID := c.GetInt64("userID")
	if err := ctrl.schoolService.UpdateSchool(id, req, userID); err != nil {
		utils.SendErrorMsg(c, http.StatusInternalServerError, utils.MsgErrInternal, err)
		return
	}

	utils.SendSuccessMsg(c, http.StatusOK, utils.MsgSuccessUpdate, nil)
}

// @Summary Delete Master School (Location)
// @Description Delete a Master School (Location) by its ID. Requires authorization.
// @Tags School Locations
// @Produce json
// @Security BearerAuth
// @Param id path int true "Master School ID"
// @Success 200 {object} utils.Response "Success"
// @Failure 400 {object} utils.ErrorResponse "Bad Request"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized"
// @Failure 500 {object} utils.ErrorResponse "Internal Server Error"
// @Router /school-locations/{id} [delete]
func (ctrl *masterSchoolController) DeleteSchool(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendErrorMsg(c, http.StatusBadRequest, utils.MsgErrInvalidID, err)
		return
	}

	if err := ctrl.schoolService.DeleteSchool(id); err != nil {
		utils.SendErrorMsg(c, http.StatusInternalServerError, utils.MsgErrInternal, err)
		return
	}

	utils.SendSuccessMsg(c, http.StatusOK, utils.MsgSuccessDelete, nil)
}
