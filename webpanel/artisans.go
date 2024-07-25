package webpanel

import (
	"encoding/base64"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	GetArtisans       = `SELECT artisan_id, name, country_id, wii_number, mac_address, number_of_posts, total_likes, is_master, last_post, mii_data FROM artisans ORDER BY artisan_id`
	GetArtisanDetails = `SELECT artisan_id, name, country_id, wii_number, mac_address, number_of_posts, total_likes, is_master, last_post, mii_data FROM artisans WHERE artisan_id = $1`
	GetArtisansMiis   = `SELECT entry_id, artisan_id, initials, nickname, gender, country_id, wii_number, mii_id, likes, perm_likes, mii_data FROM miis WHERE artisan_id = $1`
	SearchArtisans    = `SELECT artisan_id, name, country_id, wii_number, mac_address, number_of_posts, total_likes, is_master, last_post, mii_data FROM artisans WHERE name ILIKE '%' || $1 || '%' ORDER BY artisan_id`
	GetPagesArtisans  = `SELECT COUNT(*) FROM artisans`
)

type Artisans struct {
	ArtisanId      int
	Name           string
	CountryId      int
	WiiNumber      int
	MacAddress     string
	NumberOfPosts  int
	TotalLikes     int
	IsMaster       bool
	LastPost       time.Time
	MiiData        []byte
	MiiDataEncoded string
}

func (w *WebPanel) ViewArtisans(c *gin.Context) {
	pageStr := c.Param("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 0 {
		page = 1
	}

	const itemsPerPage = 150
	offset := (page - 1) * itemsPerPage

	query := GetArtisans + " LIMIT $1 OFFSET $2"

	rows, err := w.Pool.Query(w.Ctx, query, itemsPerPage, offset)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"Error": err.Error(),
		})
		return
	}

	//calculate number of pages
	var pages int
	err = w.Pool.QueryRow(w.Ctx, GetPagesArtisans).Scan(&pages)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"Error": err.Error(),
		})
		return
	}
	pages = (pages + itemsPerPage - 1) / itemsPerPage

	var artisans []Artisans
	for rows.Next() {
		artisan := Artisans{}
		err = rows.Scan(&artisan.ArtisanId, &artisan.Name, &artisan.CountryId, &artisan.WiiNumber, &artisan.MacAddress, &artisan.NumberOfPosts, &artisan.TotalLikes, &artisan.IsMaster, &artisan.LastPost, &artisan.MiiData)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "error.html", gin.H{
				"Error": err.Error(),
			})
			return
		}

		artisan.MiiDataEncoded = base64.StdEncoding.EncodeToString(artisan.MiiData)

		artisans = append(artisans, artisan)
	}

	c.HTML(http.StatusOK, "view_artisans.html", gin.H{
		"numberOfArtisans": len(artisans),
		"Artisans":         artisans,
		"Pages":            pages,
	})
}

func (w *WebPanel) ViewArtisanDetails(c *gin.Context) {
	artisanId := c.Param("artisan_id")
	row := w.Pool.QueryRow(w.Ctx, GetArtisanDetails, artisanId)

	artisanDetails := Artisans{}
	//fetch artisan details
	err := row.Scan(&artisanDetails.ArtisanId, &artisanDetails.Name, &artisanDetails.CountryId, &artisanDetails.WiiNumber, &artisanDetails.MacAddress, &artisanDetails.NumberOfPosts, &artisanDetails.TotalLikes, &artisanDetails.IsMaster, &artisanDetails.LastPost, &artisanDetails.MiiData)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"Error": err.Error(),
		})
		return
	}

	artisanDetails.MiiDataEncoded = base64.StdEncoding.EncodeToString(artisanDetails.MiiData)

	//fetch artisan miis
	rows, err := w.Pool.Query(w.Ctx, GetArtisansMiis, artisanId)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"Error": err.Error(),
		})
		return
	}

	var Miis []Plaza
	for rows.Next() {
		miidata := Plaza{}
		err = rows.Scan(&miidata.EntryId, &miidata.ArtisanId, &miidata.Initials, &miidata.Nickname, &miidata.Gender, &miidata.CountryId, &miidata.WiiNumber, &miidata.MiiId, &miidata.Likes, &miidata.PermLikes, &miidata.MiiData)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "error.html", gin.H{
				"Error": err.Error(),
			})
			return
		}
		miidata.MiiDataEncoded = base64.StdEncoding.EncodeToString(miidata.MiiData)

		Miis = append(Miis, miidata)
	}

	c.HTML(http.StatusOK, "artisan_details.html", gin.H{
		"ArtisanDetails": artisanDetails,
		"numberOfMiis":   len(Miis),
		"Miis":           Miis,
	})
}

func (w *WebPanel) SearchArtisans(c *gin.Context) {
	search := c.PostForm("search")

	rows, err := w.Pool.Query(w.Ctx, SearchArtisans, search)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"Error": err.Error(),
		})
		return
	}

	var SearchResults []Artisans
	for rows.Next() {
		search := Artisans{}
		err = rows.Scan(&search.ArtisanId, &search.Name, &search.CountryId, &search.WiiNumber, &search.MacAddress, &search.NumberOfPosts, &search.TotalLikes, &search.IsMaster, &search.LastPost, &search.MiiData)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "error.html", gin.H{
				"Error": err.Error(),
			})
			return
		}
		search.MiiDataEncoded = base64.StdEncoding.EncodeToString(search.MiiData)

		SearchResults = append(SearchResults, search)
	}

	c.HTML(http.StatusOK, "search_results.html", gin.H{
		"SearchResults": SearchResults,
		"SearchTerm":    search,
		"SearchType":    "Artisans",
	})
}
