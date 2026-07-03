package controllers

import (
	"net/http"
	"ppdb-be/models"
	"ppdb-be/services"
	"ppdb-be/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ConfigController interface {
	// SettingConfig
	CreateSettingConfig(c *gin.Context)
	GetSettingConfigs(c *gin.Context)
	GetSettingConfigByID(c *gin.Context)
	UpdateSettingConfig(c *gin.Context)
	DeleteSettingConfig(c *gin.Context)

	// Parameter
	CreateParameter(c *gin.Context)
	GetParameters(c *gin.Context)
	GetParameterByID(c *gin.Context)
	UpdateParameter(c *gin.Context)
	DeleteParameter(c *gin.Context)
}

type configController struct {
	configService services.ConfigService
}

func NewConfigController(configService services.ConfigService) ConfigController {
	return &configController{configService: configService}
}

// --- Setting Config ---

// @Summary Create a new Setting Config
// @Description Create a new Setting Config. Requires authorization.
// @Tags Setting Config
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body models.SettingConfigRequest true "Setting Config Request Payload"
// @Success 201 {object} utils.Response "Created"
// @Failure 400 {object} utils.ErrorResponse "Bad Request"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized"
// @Failure 500 {object} utils.ErrorResponse "Internal Server Error"
// @Router /settings [post]
func (ctrl *configController) CreateSettingConfig(c *gin.Context) {
	var req models.SettingConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorMsg(c, http.StatusBadRequest, utils.MsgErrBadRequest, err)
		return
	}

	userID := c.GetInt64("userID")
	if err := ctrl.configService.CreateSettingConfig(req, userID); err != nil {
		utils.SendErrorMsg(c, http.StatusInternalServerError, utils.MsgErrInternal, err)
		return
	}

	utils.SendSuccessMsg(c, http.StatusCreated, utils.MsgSuccessCreate, nil)
}

// @Summary Get all Setting Configs
// @Description Get a list of all Setting Configs.
// @Tags Setting Config
// @Produce json
// @Success 200 {object} utils.Response{data=[]models.SettingConfig} "Success"
// @Failure 500 {object} utils.ErrorResponse "Internal Server Error"
// @Router /settings [get]
func (ctrl *configController) GetSettingConfigs(c *gin.Context) {
	configs, err := ctrl.configService.GetSettingConfigs()
	if err != nil {
		utils.SendErrorMsg(c, http.StatusInternalServerError, utils.MsgErrInternal, err)
		return
	}

	utils.SendSuccessMsg(c, http.StatusOK, utils.MsgSuccessFetch, configs)
}

// @Summary Get Setting Config by ID
// @Description Get a Setting Config by its ID.
// @Tags Setting Config
// @Produce json
// @Param id path int true "Setting Config ID"
// @Success 200 {object} utils.Response{data=models.SettingConfig} "Success"
// @Failure 400 {object} utils.ErrorResponse "Bad Request"
// @Failure 404 {object} utils.ErrorResponse "Not Found"
// @Router /settings/{id} [get]
func (ctrl *configController) GetSettingConfigByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendErrorMsg(c, http.StatusBadRequest, utils.MsgErrInvalidID, err)
		return
	}

	config, err := ctrl.configService.GetSettingConfigByID(id)
	if err != nil {
		utils.SendErrorMsg(c, http.StatusNotFound, utils.MsgErrNotFound, err)
		return
	}

	utils.SendSuccessMsg(c, http.StatusOK, utils.MsgSuccessFetch, config)
}

// @Summary Update Setting Config
// @Description Update an existing Setting Config by its ID. Requires authorization.
// @Tags Setting Config
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Setting Config ID"
// @Param body body models.SettingConfigRequest true "Setting Config Request Payload"
// @Success 200 {object} utils.Response "Success"
// @Failure 400 {object} utils.ErrorResponse "Bad Request"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized"
// @Failure 500 {object} utils.ErrorResponse "Internal Server Error"
// @Router /settings/{id} [put]
func (ctrl *configController) UpdateSettingConfig(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendErrorMsg(c, http.StatusBadRequest, utils.MsgErrInvalidID, err)
		return
	}

	var req models.SettingConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorMsg(c, http.StatusBadRequest, utils.MsgErrBadRequest, err)
		return
	}

	userID := c.GetInt64("userID")
	if err := ctrl.configService.UpdateSettingConfig(id, req, userID); err != nil {
		utils.SendErrorMsg(c, http.StatusInternalServerError, utils.MsgErrInternal, err)
		return
	}

	utils.SendSuccessMsg(c, http.StatusOK, utils.MsgSuccessUpdate, nil)
}

// @Summary Delete Setting Config
// @Description Delete a Setting Config by its ID. Requires authorization.
// @Tags Setting Config
// @Produce json
// @Security BearerAuth
// @Param id path int true "Setting Config ID"
// @Success 200 {object} utils.Response "Success"
// @Failure 400 {object} utils.ErrorResponse "Bad Request"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized"
// @Failure 500 {object} utils.ErrorResponse "Internal Server Error"
// @Router /settings/{id} [delete]
func (ctrl *configController) DeleteSettingConfig(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendErrorMsg(c, http.StatusBadRequest, utils.MsgErrInvalidID, err)
		return
	}

	if err := ctrl.configService.DeleteSettingConfig(id); err != nil {
		utils.SendErrorMsg(c, http.StatusInternalServerError, utils.MsgErrInternal, err)
		return
	}

	utils.SendSuccessMsg(c, http.StatusOK, utils.MsgSuccessDelete, nil)
}

// --- Parameter ---

// @Summary Create a new Parameter
// @Description Create a new Parameter. Requires authorization.
// @Tags Parameters
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body models.ParameterRequest true "Parameter Request Payload"
// @Success 201 {object} utils.Response "Created"
// @Failure 400 {object} utils.ErrorResponse "Bad Request"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized"
// @Failure 500 {object} utils.ErrorResponse "Internal Server Error"
// @Router /parameters [post]
func (ctrl *configController) CreateParameter(c *gin.Context) {
	var req models.ParameterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorMsg(c, http.StatusBadRequest, utils.MsgErrBadRequest, err)
		return
	}

	userID := c.GetInt64("userID")
	if err := ctrl.configService.CreateParameter(req, userID); err != nil {
		utils.SendErrorMsg(c, http.StatusInternalServerError, utils.MsgErrInternal, err)
		return
	}

	utils.SendSuccessMsg(c, http.StatusCreated, utils.MsgSuccessCreate, nil)
}

// @Summary Get all Parameters
// @Description Get a list of all Parameters.
// @Tags Parameters
// @Produce json
// @Success 200 {object} utils.Response{data=[]models.Parameter} "Success"
// @Failure 500 {object} utils.ErrorResponse "Internal Server Error"
// @Router /parameters [get]
func (ctrl *configController) GetParameters(c *gin.Context) {
	params, err := ctrl.configService.GetParameters()
	if err != nil {
		utils.SendErrorMsg(c, http.StatusInternalServerError, utils.MsgErrInternal, err)
		return
	}

	utils.SendSuccessMsg(c, http.StatusOK, utils.MsgSuccessFetch, params)
}

// @Summary Get Parameter by ID
// @Description Get a Parameter by its ID.
// @Tags Parameters
// @Produce json
// @Param id path int true "Parameter ID"
// @Success 200 {object} utils.Response{data=models.Parameter} "Success"
// @Failure 400 {object} utils.ErrorResponse "Bad Request"
// @Failure 404 {object} utils.ErrorResponse "Not Found"
// @Router /parameters/{id} [get]
func (ctrl *configController) GetParameterByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendErrorMsg(c, http.StatusBadRequest, utils.MsgErrInvalidID, err)
		return
	}

	param, err := ctrl.configService.GetParameterByID(id)
	if err != nil {
		utils.SendErrorMsg(c, http.StatusNotFound, utils.MsgErrNotFound, err)
		return
	}

	utils.SendSuccessMsg(c, http.StatusOK, utils.MsgSuccessFetch, param)
}

// @Summary Update Parameter
// @Description Update an existing Parameter by its ID. Requires authorization.
// @Tags Parameters
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Parameter ID"
// @Param body body models.ParameterRequest true "Parameter Request Payload"
// @Success 200 {object} utils.Response "Success"
// @Failure 400 {object} utils.ErrorResponse "Bad Request"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized"
// @Failure 500 {object} utils.ErrorResponse "Internal Server Error"
// @Router /parameters/{id} [put]
func (ctrl *configController) UpdateParameter(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendErrorMsg(c, http.StatusBadRequest, utils.MsgErrInvalidID, err)
		return
	}

	var req models.ParameterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorMsg(c, http.StatusBadRequest, utils.MsgErrBadRequest, err)
		return
	}

	userID := c.GetInt64("userID")
	if err := ctrl.configService.UpdateParameter(id, req, userID); err != nil {
		utils.SendErrorMsg(c, http.StatusInternalServerError, utils.MsgErrInternal, err)
		return
	}

	utils.SendSuccessMsg(c, http.StatusOK, utils.MsgSuccessUpdate, nil)
}

// @Summary Delete Parameter
// @Description Delete a Parameter by its ID. Requires authorization.
// @Tags Parameters
// @Produce json
// @Security BearerAuth
// @Param id path int true "Parameter ID"
// @Success 200 {object} utils.Response "Success"
// @Failure 400 {object} utils.ErrorResponse "Bad Request"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized"
// @Failure 500 {object} utils.ErrorResponse "Internal Server Error"
// @Router /parameters/{id} [delete]
func (ctrl *configController) DeleteParameter(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendErrorMsg(c, http.StatusBadRequest, utils.MsgErrInvalidID, err)
		return
	}

	if err := ctrl.configService.DeleteParameter(id); err != nil {
		utils.SendErrorMsg(c, http.StatusInternalServerError, utils.MsgErrInternal, err)
		return
	}

	utils.SendSuccessMsg(c, http.StatusOK, utils.MsgSuccessDelete, nil)
}
