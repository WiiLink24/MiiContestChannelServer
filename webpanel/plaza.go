package webpanel

import (
	"encoding/base64"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	GetPlaza = `SELECT m.entry_id, m.artisan_id, m.initials, m.nickname, m.gender, m.country_id, m.wii_number, m.mii_id, m.likes, m.perm_likes, m.mii_data, a.name FROM miis m LEFT JOIN artisans a ON m.artisan_id = a.artisan_id ORDER BY entry_id`
	GetPages = `SELECT COUNT(*) FROM miis`
	DeleteMii = `DELETE FROM miis WHERE entry_id = $1`
	GetMiiDetails = `SELECT m.entry_id, m.artisan_id, m.initials, m.nickname, m.gender, m.country_id, m.wii_number, m.mii_id, m.likes, m.perm_likes, m.mii_data, a.name FROM miis m LEFT JOIN artisans a ON m.artisan_id = a.artisan_id WHERE entry_id = $1`
	SearchMiis = `SELECT entry_id, artisan_id, initials, nickname, gender, country_id, wii_number, mii_id, likes, perm_likes, mii_data FROM miis WHERE nickname ILIKE '%' || $1 || '%' ORDER BY entry_id`
	GetPlazaTop50 = `SELECT m.entry_id, m.artisan_id, m.initials, m.nickname, m.gender, m.country_id, m.wii_number, m.mii_id, m.likes, m.perm_likes, m.mii_data, a.name FROM miis m LEFT JOIN artisans a ON m.artisan_id = a.artisan_id ORDER BY perm_likes DESC LIMIT 50`
	GetPlazaNew = `SELECT m.entry_id, m.artisan_id, m.initials, m.nickname, m.gender, m.country_id, m.wii_number, m.mii_id, m.likes, m.perm_likes, m.mii_data, a.name FROM miis m LEFT JOIN artisans a ON m.artisan_id = a.artisan_id ORDER BY likes DESC`
)

type Plaza struct {
	EntryId int
	ArtisanId int
	Initials string
	Nickname string
	Gender int
	CountryId int
	WiiNumber int
	MiiId []byte
	Likes int
	PermLikes int
	MiiData []byte
	ArtisanName string
	MiiDataEncoded string
}

func (w *WebPanel) ViewPlaza(c *gin.Context) {
	pageStr := c.Param("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 0 {
		page = 1
	}

	const itemsPerPage = 150
	offset := (page - 1) * itemsPerPage

	query := GetPlaza + " LIMIT $1 OFFSET $2"

	rows, err := w.Pool.Query(w.Ctx, query, itemsPerPage, offset)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"Error": err.Error(),
		})
		return
	}

	//calculate number of pages
	var pages int
	err = w.Pool.QueryRow(w.Ctx, GetPages).Scan(&pages)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"Error": err.Error(),
		})
		return
	}
	pages = (pages + itemsPerPage - 1) / itemsPerPage

	var plaza []Plaza
	for rows.Next() {
		plazadata := Plaza{}
		err = rows.Scan(&plazadata.EntryId, &plazadata.ArtisanId, &plazadata.Initials, &plazadata.Nickname, &plazadata.Gender, &plazadata.CountryId, &plazadata.WiiNumber, &plazadata.MiiId, &plazadata.Likes, &plazadata.PermLikes, &plazadata.MiiData, &plazadata.ArtisanName)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "error.html", gin.H{
				"Error": err.Error(),
			})
			return
		}
		plazadata.MiiDataEncoded = base64.StdEncoding.EncodeToString(plazadata.MiiData)

		plaza = append(plaza, plazadata)
	}

	c.HTML(http.StatusOK, "view_plaza.html", gin.H{
		"numberOfMiis": len(plaza),
		"Plaza":        plaza,
		"Pages": 		pages,

	})
}


func (w *WebPanel) ViewPlazaTop50(c *gin.Context) {
	rows, err := w.Pool.Query(w.Ctx, GetPlazaTop50)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"Error": err.Error(),
		})
		return
	}

	var plaza []Plaza
	for rows.Next() {
		plazadata := Plaza{}
		err = rows.Scan(&plazadata.EntryId, &plazadata.ArtisanId, &plazadata.Initials, &plazadata.Nickname, &plazadata.Gender, &plazadata.CountryId, &plazadata.WiiNumber, &plazadata.MiiId, &plazadata.Likes, &plazadata.PermLikes, &plazadata.MiiData, &plazadata.ArtisanName)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "error.html", gin.H{
				"Error": err.Error(),
			})
			return
		}
		plazadata.MiiDataEncoded = base64.StdEncoding.EncodeToString(plazadata.MiiData)

		plaza = append(plaza, plazadata)
	}

	c.HTML(http.StatusOK, "view_plaza.html", gin.H{
		"numberOfMiis": len(plaza),
		"Plaza":         plaza,
	})
}

func (w *WebPanel) ViewPlazaNew(c *gin.Context) {
	pageStr := c.Param("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 0 {
		page = 1
	}

	const itemsPerPage = 150
	offset := (page - 1) * itemsPerPage

	query := GetPlazaNew + " LIMIT $1 OFFSET $2"

	rows, err := w.Pool.Query(w.Ctx, query, itemsPerPage, offset)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"Error": err.Error(),
		})
		return
	}

	//calculate number of pages
	var pages int
	err = w.Pool.QueryRow(w.Ctx, GetPages).Scan(&pages)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"Error": err.Error(),
		})
		return
	}
	pages = (pages + itemsPerPage - 1) / itemsPerPage

	var plaza []Plaza
	for rows.Next() {
		plazadata := Plaza{}
		err = rows.Scan(&plazadata.EntryId, &plazadata.ArtisanId, &plazadata.Initials, &plazadata.Nickname, &plazadata.Gender, &plazadata.CountryId, &plazadata.WiiNumber, &plazadata.MiiId, &plazadata.Likes, &plazadata.PermLikes, &plazadata.MiiData, &plazadata.ArtisanName)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "error.html", gin.H{
				"Error": err.Error(),
			})
			return
		}
		plazadata.MiiDataEncoded = base64.StdEncoding.EncodeToString(plazadata.MiiData)

		plaza = append(plaza, plazadata)
	}

	c.HTML(http.StatusOK, "view_plaza.html", gin.H{
		"numberOfMiis": len(plaza),
		"Plaza":        plaza,
		"Pages": 		pages,
	})
}

func (w *WebPanel) ViewMiiDetails(c *gin.Context) {
	entryId := c.Param("entry_id")
	row := w.Pool.QueryRow(w.Ctx, GetMiiDetails, entryId)

	MiiDetails := Plaza{}
	err := row.Scan(&MiiDetails.EntryId, &MiiDetails.ArtisanId, &MiiDetails.Initials, &MiiDetails.Nickname, &MiiDetails.Gender, &MiiDetails.CountryId, &MiiDetails.WiiNumber, &MiiDetails.MiiId, &MiiDetails.Likes, &MiiDetails.PermLikes, &MiiDetails.MiiData, &MiiDetails.ArtisanName)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"Error": err.Error(),
		})
		return
	}

	MiiDetails.MiiDataEncoded = base64.StdEncoding.EncodeToString(MiiDetails.MiiData)

	c.HTML(http.StatusOK, "view_mii.html", gin.H{
		"MiiDetails": MiiDetails,
	})
}

func (w *WebPanel) DeleteMii(c *gin.Context) {
	entryId := c.Param("entry_id")
	_, err := w.Pool.Exec(w.Ctx, DeleteMii, entryId)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"Error": err.Error(),
		})
		return
	}

	c.Redirect(http.StatusFound, "/panel/plaza#delete_success")
}

func (w *WebPanel) SearchPlaza(c *gin.Context) {
	search := c.PostForm("search")

	rows, err := w.Pool.Query(w.Ctx, SearchMiis, search)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"Error": err.Error(),
		})
		return
	}

	var SearchResults []Plaza
	for rows.Next() {
		search := Plaza{}
		err = rows.Scan(&search.EntryId, &search.ArtisanId, &search.Initials, &search.Nickname, &search.Gender, &search.CountryId, &search.WiiNumber, &search.MiiId, &search.Likes, &search.PermLikes, &search.MiiData)
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
		"SearchTerm":	 search,
		"SearchType":    "Plaza",
	})
}