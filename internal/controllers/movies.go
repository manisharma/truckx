package controllers

import (
	"net/http"
	"strconv"
	"strings"
	"time"
	"truckx/internal/models"
	"truckx/internal/repository"
	"truckx/internal/services/util"

	"github.com/gin-gonic/gin"
)

type movieCtrl struct {
	repo repository.MovieRepo
}

func NewMovieCtrl(repo repository.MovieRepo) *movieCtrl {
	return &movieCtrl{repo}
}

func (ctrl *movieCtrl) Create(c *gin.Context) {
	var user *models.User
	if usr, ok := c.Get(string(models.UserCtxKey)); !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not valid"})
		c.Abort()
		return
	} else {
		user = usr.(*models.User)
	}
	if user.RoleId != 1 {
		c.JSON(http.StatusForbidden, gin.H{"error": http.StatusText(http.StatusForbidden)})
		c.Abort()
		return
	}

	movie := &models.Movie{}
	err := c.ShouldBindJSON(movie)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	movie, err = ctrl.repo.Create(c.Request.Context(), *movie)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, gin.H{"movie": movie})
}

func (ctrl *movieCtrl) List(c *gin.Context) {
	movies, err := ctrl.repo.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, gin.H{"movies": movies})
}

func (ctrl *movieCtrl) GetAvailability(c *gin.Context) {
	movieIdStr, ok := c.Params.Get("id")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user id missing"})
		c.Abort()
		return
	}

	movieId, err := strconv.Atoi(movieIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		c.Abort()
		return
	}

	state := c.Request.URL.Query().Get("state")
	if strings.EqualFold(state, "") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "location/state param missing"})
		c.Abort()
		return
	}

	var date time.Time
	dateStr := c.Request.URL.Query().Get("date")
	if strings.EqualFold(dateStr, "") { // default to today
		date = time.Now()
	} else {
		date, err = time.Parse(util.DateFormat, dateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date param, need YYYY-MM-DD"})
			c.Abort()
			return
		}
	}

	availability, err := ctrl.repo.GetAvailability(c.Request.Context(), movieId, state, date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, gin.H{"availability": availability})
}
