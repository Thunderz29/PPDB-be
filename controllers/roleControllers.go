package controllers

import (
	"net/http"
	"ppdb-be/models"
	"ppdb-be/services"
	"ppdb-be/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RoleController interface {
	CreateRole(c *gin.Context)
	GetRoles(c *gin.Context)
	GetRoleByID(c *gin.Context)
	UpdateRole(c *gin.Context)
	DeleteRole(c *gin.Context)
}

type roleController struct {
	roleService services.RoleService
}

func NewRoleController(roleService services.RoleService) RoleController {
	return &roleController{roleService: roleService}
}

// @Summary Create a new Role
// @Description Create a new Role. Requires authorization.
// @Tags Roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body models.RoleRequest true "Role Request Payload"
// @Success 201 {object} utils.Response "Created"
// @Failure 400 {object} utils.ErrorResponse "Bad Request"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized"
// @Failure 500 {object} utils.ErrorResponse "Internal Server Error"
// @Router /roles [post]
func (ctrl *roleController) CreateRole(c *gin.Context) {
	var req models.RoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorMsg(c, http.StatusBadRequest, utils.MsgErrBadRequest, err)
		return
	}

	userID := c.GetInt64("userID")
	if err := ctrl.roleService.CreateRole(req, userID); err != nil {
		utils.SendErrorMsg(c, http.StatusInternalServerError, utils.MsgErrInternal, err)
		return
	}

	utils.SendSuccessMsg(c, http.StatusCreated, utils.MsgSuccessCreate, nil)
}

// @Summary Get all Roles
// @Description Get a list of all Roles.
// @Tags Roles
// @Produce json
// @Success 200 {object} utils.Response{data=[]models.Role} "Success"
// @Failure 500 {object} utils.ErrorResponse "Internal Server Error"
// @Router /roles [get]
func (ctrl *roleController) GetRoles(c *gin.Context) {
	roles, err := ctrl.roleService.GetRoles()
	if err != nil {
		utils.SendErrorMsg(c, http.StatusInternalServerError, utils.MsgErrInternal, err)
		return
	}

	utils.SendSuccessMsg(c, http.StatusOK, utils.MsgSuccessFetch, roles)
}

// @Summary Get Role by ID
// @Description Get a Role by its ID.
// @Tags Roles
// @Produce json
// @Param id path int true "Role ID"
// @Success 200 {object} utils.Response{data=models.Role} "Success"
// @Failure 400 {object} utils.ErrorResponse "Bad Request"
// @Failure 404 {object} utils.ErrorResponse "Not Found"
// @Router /roles/{id} [get]
func (ctrl *roleController) GetRoleByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendErrorMsg(c, http.StatusBadRequest, utils.MsgErrInvalidID, err)
		return
	}

	role, err := ctrl.roleService.GetRoleByID(id)
	if err != nil {
		utils.SendErrorMsg(c, http.StatusNotFound, utils.MsgErrNotFound, err)
		return
	}

	utils.SendSuccessMsg(c, http.StatusOK, utils.MsgSuccessFetch, role)
}

// @Summary Update Role
// @Description Update an existing Role by its ID. Requires authorization.
// @Tags Roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Role ID"
// @Param body body models.RoleRequest true "Role Request Payload"
// @Success 200 {object} utils.Response "Success"
// @Failure 400 {object} utils.ErrorResponse "Bad Request"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized"
// @Failure 500 {object} utils.ErrorResponse "Internal Server Error"
// @Router /roles/{id} [put]
func (ctrl *roleController) UpdateRole(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendErrorMsg(c, http.StatusBadRequest, utils.MsgErrInvalidID, err)
		return
	}

	var req models.RoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorMsg(c, http.StatusBadRequest, utils.MsgErrBadRequest, err)
		return
	}

	userID := c.GetInt64("userID")
	if err := ctrl.roleService.UpdateRole(id, req, userID); err != nil {
		utils.SendErrorMsg(c, http.StatusInternalServerError, utils.MsgErrInternal, err)
		return
	}

	utils.SendSuccessMsg(c, http.StatusOK, utils.MsgSuccessUpdate, nil)
}

// @Summary Delete Role
// @Description Delete a Role by its ID. Requires authorization.
// @Tags Roles
// @Produce json
// @Security BearerAuth
// @Param id path int true "Role ID"
// @Success 200 {object} utils.Response "Success"
// @Failure 400 {object} utils.ErrorResponse "Bad Request"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized"
// @Failure 500 {object} utils.ErrorResponse "Internal Server Error"
// @Router /roles/{id} [delete]
func (ctrl *roleController) DeleteRole(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendErrorMsg(c, http.StatusBadRequest, utils.MsgErrInvalidID, err)
		return
	}

	if err := ctrl.roleService.DeleteRole(id); err != nil {
		utils.SendErrorMsg(c, http.StatusInternalServerError, utils.MsgErrInternal, err)
		return
	}

	utils.SendSuccessMsg(c, http.StatusOK, utils.MsgSuccessDelete, nil)
}
