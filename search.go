package main

import (
	"github.com/WiiLink24/MiiContestChannel/common"
	"github.com/WiiLink24/MiiContestChannel/plaza"
	"github.com/gin-gonic/gin"
	"net/http"
)

const GetMiiByEntryNumber = `SELECT miis.entry_id, miis.initials, miis.perm_likes, miis.skill, miis.country_id, miis.mii_data, 
       			artisans.mii_data, artisans.artisan_id, artisans.is_master 
				FROM miis, artisans WHERE miis.artisan_id = artisans.artisan_id AND miis.entry_id = $1
				ORDER BY miis.likes`

func search(c *gin.Context) {
	entryNumber := c.Query("entryno")

	var likes int
	var isMaster bool
	mii := common.MiiWithArtisan{}
	err := pool.QueryRow(ctx, GetMiiByEntryNumber, entryNumber).Scan(&mii.EntryNumber, &mii.Initials, &likes, &mii.Skill, &mii.CountryCode, &mii.MiiData, &mii.ArtisanMiiData, &mii.ArtisanId, &isMaster)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		writeResult(c, 500)
		return
	}

	if isMaster {
		mii.IsMasterArtisan = 1
	}

	mii.Likes = uint8(likes)

	c.Data(http.StatusOK, "application/octet-stream", plaza.MakeSearchList(common.Search, []common.MiiWithArtisan{mii}, mii.EntryNumber))
}
