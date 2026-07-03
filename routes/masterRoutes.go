package routes

import (
	"ppdb-be/controllers"
	"ppdb-be/middleware"
	"ppdb-be/repositories"

	"github.com/gin-gonic/gin"
)

func MasterRoutes(
	router *gin.RouterGroup,
	roleCtrl controllers.RoleController,
	profileCtrl controllers.SchoolProfileController,
	configCtrl controllers.ConfigController,
	schoolCtrl controllers.MasterSchoolController,
	sessionRepo repositories.SessionRepository,
) {

	// --- Roles ---
	roles := router.Group("/roles")
	{
		roles.GET("", roleCtrl.GetRoles)
		roles.GET("/:id", roleCtrl.GetRoleByID)
	}
	protectedRoles := router.Group("/roles")
	protectedRoles.Use(middleware.AuthMiddleware(sessionRepo))
	{
		protectedRoles.POST("", roleCtrl.CreateRole)
		protectedRoles.PUT("/:id", roleCtrl.UpdateRole)
		protectedRoles.DELETE("/:id", roleCtrl.DeleteRole)
	}

	// --- School Profiles ---
	profiles := router.Group("/school-profiles")
	{
		profiles.GET("", profileCtrl.GetSchoolProfiles)
		profiles.GET("/:id", profileCtrl.GetSchoolProfileByID)
	}
	protectedProfiles := router.Group("/school-profiles")
	protectedProfiles.Use(middleware.AuthMiddleware(sessionRepo))
	{
		protectedProfiles.POST("", profileCtrl.CreateSchoolProfile)
		protectedProfiles.PUT("/:id", profileCtrl.UpdateSchoolProfile)
		protectedProfiles.DELETE("/:id", profileCtrl.DeleteSchoolProfile)
	}

	// --- Settings Config ---
	settings := router.Group("/settings")
	{
		settings.GET("", configCtrl.GetSettingConfigs)
		settings.GET("/:id", configCtrl.GetSettingConfigByID)
	}
	protectedSettings := router.Group("/settings")
	protectedSettings.Use(middleware.AuthMiddleware(sessionRepo))
	{
		protectedSettings.POST("", configCtrl.CreateSettingConfig)
		protectedSettings.PUT("/:id", configCtrl.UpdateSettingConfig)
		protectedSettings.DELETE("/:id", configCtrl.DeleteSettingConfig)
	}

	// --- Parameters ---
	params := router.Group("/parameters")
	{
		params.GET("", configCtrl.GetParameters)
		params.GET("/:id", configCtrl.GetParameterByID)
	}
	protectedParams := router.Group("/parameters")
	protectedParams.Use(middleware.AuthMiddleware(sessionRepo))
	{
		protectedParams.POST("", configCtrl.CreateParameter)
		protectedParams.PUT("/:id", configCtrl.UpdateParameter)
		protectedParams.DELETE("/:id", configCtrl.DeleteParameter)
	}

	// --- Master Schools (Locations) ---
	schools := router.Group("/school-locations")
	{
		schools.GET("", schoolCtrl.GetSchools)
		schools.GET("/:id", schoolCtrl.GetSchoolByID)
	}
	protectedSchools := router.Group("/school-locations")
	protectedSchools.Use(middleware.AuthMiddleware(sessionRepo))
	{
		protectedSchools.POST("", schoolCtrl.CreateSchool)
		protectedSchools.PUT("/:id", schoolCtrl.UpdateSchool)
		protectedSchools.DELETE("/:id", schoolCtrl.DeleteSchool)
	}
}
