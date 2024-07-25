package main

import (
	"github.com/WiiLink24/MiiContestChannel/common"
	"github.com/WiiLink24/MiiContestChannel/plaza"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

const (
	GetMiiName     = `SELECT miis.nickname FROM miis WHERE entry_id = $1`
	GetRelatedMiis = `SELECT miis.entry_id, miis.artisan_id, miis.initials, miis.perm_likes, miis.skill, miis.country_id, miis.mii_data, artisans.mii_data, artisans.is_master FROM miis, artisans WHERE artisans.artisan_id = miis.artisan_id AND miis.nickname ILIKE '%' || $1 || '%' LIMIT 50`
)

func nameSearch(c *gin.Context) {
	strEntryNumber := c.Query("entryno")

	entryId, err := strconv.Atoi(strEntryNumber)
	if err != nil {
		// TODO: Figure out invalid entry ID
		// It should never happen if this was sent by CMOC, but people will be people.
		writeResult(c, 108)
		return
	}

	// Next make sure the Mii exists.
	var exists bool
	err = pool.QueryRow(ctx, DoesMiiExist, entryId).Scan(&exists)
	if err != nil {
		writeResult(c, 500)
		return
	}

	if !exists {
		// Mii does not exist.
		c.Data(http.StatusOK, "application/octet-stream", plaza.MakeSearchList(common.NameSearch, nil, uint32(entryId)))
		return
	}

	var nickname string
	err = pool.QueryRow(ctx, GetMiiName, entryId).Scan(&nickname)
	if err != nil {
		writeResult(c, 500)
		return
	}

	rows, err := pool.Query(ctx, GetRelatedMiis, nickname)
	if err != nil {
		writeResult(c, 500)
		return
	}

	var miis []common.MiiWithArtisan

	defer rows.Close()
	for rows.Next() {
		mii := common.MiiWithArtisan{}

		var likes int
		var isMaster bool
		err = rows.Scan(&mii.EntryNumber, &mii.ArtisanId, &mii.Initials, &likes, &mii.Skill, &mii.CountryCode, &mii.MiiData, &mii.ArtisanMiiData, &isMaster)

		if isMaster {
			mii.IsMasterArtisan = 1
		}

		mii.Likes = uint8(likes)

		miis = append(miis, mii)
	}

	c.Data(http.StatusOK, "application/octet-stream", plaza.MakeSearchList(common.NameSearch, miis, uint32(entryId)))
}
