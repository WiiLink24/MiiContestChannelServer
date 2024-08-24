package main

import (
	"errors"
	"github.com/WiiLink24/MiiContestChannel/common"
	"github.com/WiiLink24/MiiContestChannel/plaza"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"net/http"
)

const (
	GetArtisanInfo    = `SELECT artisan_id, is_master, total_likes FROM artisans WHERE artisan_id = $1`
	GetArtisanRanking = `WITH RankedCrafts AS (
    SELECT artisan_id, 
           ROW_NUMBER() OVER (ORDER BY total_likes DESC) AS ranking
    FROM artisans
    ORDER BY total_likes DESC
)
SELECT COALESCE(ranking, 0) AS ranking
FROM RankedCrafts
WHERE artisan_id = $1`

	GetMostPopularPost = `SELECT entry_id, mii_data, initials FROM miis WHERE artisan_id = $1 ORDER BY perm_likes DESC`
)

func info(c *gin.Context) {
	artisanId := c.Query("craftsno")

	var isMaster bool
	var popularity uint8
	var intArtisanId uint32
	err := pool.QueryRow(ctx, GetArtisanInfo, artisanId).Scan(&intArtisanId, &isMaster, &popularity)
	if errors.Is(err, pgx.ErrNoRows) {
		c.Data(http.StatusOK, "application/octet-stream", plaza.MakeArtisanInfo(intArtisanId, &plaza.MiiInfo{}, &plaza.ArtisanInfo{}))
	} else if err != nil {
		c.Status(http.StatusBadRequest)
		writeResult(c, 400)
		return
	}

	var ranking int
	err = pool.QueryRow(ctx, GetArtisanRanking, artisanId).Scan(&ranking)
	if errors.Is(err, pgx.ErrNoRows) {
		c.Data(http.StatusOK, "application/octet-stream", plaza.MakeArtisanInfo(intArtisanId, &plaza.MiiInfo{}, &plaza.ArtisanInfo{}))
	} else if err != nil {
		// TODO: If nil it is a different error
		c.Status(http.StatusBadRequest)
		writeResult(c, 400)
		return
	}

	var entryNumber uint32
	var miiData []byte
	var initials string
	err = pool.QueryRow(ctx, GetMostPopularPost, artisanId).Scan(&entryNumber, &miiData, &initials)
	if errors.Is(err, pgx.ErrNoRows) {
		// Artisan didn't post, send default Mii.
		miiData = []byte{128, 0, 0, 63, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 64, 64, 134, 57, 123, 203, 194, 38, 92, 40, 0, 4, 66, 64, 49, 189, 40, 162, 8, 140, 8, 64, 20, 73, 184, 141, 0, 138, 0, 138, 37, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 234, 41}
		entryNumber = 1
		initials = "00"
	} else if err != nil {
		c.Status(http.StatusBadRequest)
		writeResult(c, 400)
		return
	}

	miiInfo := &plaza.MiiInfo{
		Tag:         common.MiiInfo,
		TagSize:     96,
		Unknown:     1,
		EntryNumber: entryNumber,
		MiiData:     [76]byte{},
		Unknown2:    1,
		Initials:    [2]byte{},
	}

	copy(miiInfo.MiiData[:], miiData)
	copy(miiInfo.Initials[:], initials)

	artisanInfo := &plaza.ArtisanInfo{
		Tag:        common.ArtisanInfo,
		TagSize:    24,
		Unknown:    1,
		Unknown1:   1,
		Unknown2:   1,
		IsMaster:   0,
		Popularity: popularity,
		Ranking:    uint8(ranking),
		Unknown3:   1,
	}

	if isMaster {
		artisanInfo.IsMaster = 1
	}

	c.Data(http.StatusOK, "application/octet-stream", plaza.MakeArtisanInfo(intArtisanId, miiInfo, artisanInfo))
}
