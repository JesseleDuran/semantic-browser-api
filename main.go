package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB

func init() {
	//open a db connection
	var err error
	db, err = gorm.Open("mysql", "root:root@/semantic-browser?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("failed to connect database")
	}

	//Migrate the schema
	db.AutoMigrate(&favModel{})
}

func main() {

	router := gin.Default()

	v1 := router.Group("/api/v1/favs")
	{
		v1.POST("/", addFav)
		v1.GET("/:id", fetchSingleFav)
		v1.GET("/:id/all", fetchAllFavsFromUser)
		v1.PUT("/:id", updateFav)
		v1.DELETE("/:id", deleteFav)
	}
	router.Run()

}

type (
	// favModel describes a favModel type
	favModel struct {
		gorm.Model
		Link   string `json:"link"`
		IDUser string `json:"id-user"`
	}

	// transformedFav represents a formatted fav
	transformedFav struct {
		ID     uint   `json:"id"`
		Link   string `json:"link"`
		IDUser string `json:"id-user"`
	}
)

// addFav add a new fav
func addFav(c *gin.Context) {
	fav := favModel{IDUser: c.PostForm("id-user"), Link: c.PostForm("link")}
	db.Save(&fav)
	c.JSON(http.StatusCreated, gin.H{"status": http.StatusCreated, "message": "Fav item added successfully", "resourceId": fav})
}

// fetchAllFavs fetch all favs
func fetchAllFavsFromUser(c *gin.Context) {
	var favs []favModel
	var _favs []transformedFav
	userID := c.Param("id")

	db.Where("id_user = ?", userID).Find(&favs)

	if len(favs) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No favs found!"})
		return
	}

	//transforms the favs for building a good response
	for _, item := range favs {
		_favs = append(_favs, transformedFav{ID: item.ID, Link: item.Link, IDUser: item.IDUser})
	}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": _favs})
}

// fetchSingleFav fetch a single fav
func fetchSingleFav(c *gin.Context) {
	var fav favModel
	favID := c.Param("id")

	db.First(&fav, favID)

	if fav.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No fav found!"})
		return
	}

	_fav := transformedFav{ID: fav.ID, Link: fav.Link, IDUser: fav.IDUser}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": _fav})
}

// updateFav update a fav
func updateFav(c *gin.Context) {
	var fav favModel
	favID := c.Param("id")

	db.First(&fav, favID)

	if fav.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No fav found!"})
		return
	}
	fav.Link = c.PostForm("link")
	db.Save(&fav)
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Fav updated successfully!"})
}

// deleteFav remove a fav
func deleteFav(c *gin.Context) {
	var fav favModel
	favID := c.Param("id")

	db.First(&fav, favID)

	if fav.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No fav found!"})
		return
	}

	db.Delete(&fav)
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "fav deleted successfully!"})
}
