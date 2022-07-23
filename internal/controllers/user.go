package controllers

import (
	"net/http"
	"strconv"
	"truckx/internal/models"
	"truckx/internal/repository"

	"github.com/gin-gonic/gin"
)

type userCtrl struct {
	repo repository.UserRepo
}

func NewUserCtrl(repo repository.UserRepo) *userCtrl {
	return &userCtrl{repo}
}

func (ctrl *userCtrl) Get(c *gin.Context) {
	var ctxUsr models.User
	if usr, ok := c.Get(string(models.UserCtxKey)); !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not valid"})
		c.Abort()
		return
	} else {
		ctxUsr = usr.(models.User)
	}

	userIdStr, ok := c.Params.Get("id")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user id missing"})
		c.Abort()
		return
	}

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		c.Abort()
		return
	}

	if ctxUsr.RoleId != 1 && userId != int(ctxUsr.Id) {
		c.JSON(http.StatusForbidden, gin.H{"error": http.StatusText(http.StatusForbidden)})
		c.Abort()
		return
	}

	movies, err := ctrl.repo.Get(c.Request.Context(), userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, gin.H{"movies": movies})
}

func (ctrl *userCtrl) Delete(c *gin.Context) {
	var ctxUsr models.User
	if usr, ok := c.Get(string(models.UserCtxKey)); !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not valid"})
		c.Abort()
		return
	} else {
		ctxUsr = usr.(models.User)
	}

	userIdStr, ok := c.Params.Get("id")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user id missing"})
		c.Abort()
		return
	}

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		c.Abort()
		return
	}

	if ctxUsr.RoleId != 1 && userId != int(ctxUsr.Id) {
		c.JSON(http.StatusForbidden, gin.H{"error": http.StatusText(http.StatusForbidden)})
		c.Abort()
		return
	}

	movies, err := ctrl.repo.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, gin.H{"movies": movies})
}
