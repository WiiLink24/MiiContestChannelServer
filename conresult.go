package main

import (
	"bytes"
	"encoding/binary"
	"github.com/gin-gonic/gin"
	"math"
	"net/http"
	"strconv"
)

const GetMiiRank = `SELECT rank FROM contest_miis WHERE artisan_id = $1 AND contest_id = $2`

/*
For some reason, conResult takes a binary blob akin to those found in the MiiContestChannel repository. Unlike them
however, they are unencrypted, uncompressed, and not signed.
*/

type ConResult struct {
	Header ConResultHeader
	Miis   []ConResultMii
}

type ConResultHeader struct {
	Tag       [2]byte
	_         uint16
	ContestId uint32
	_         uint32
	_         uint32
	_         [12]byte
	Padding   [4]byte
}

type ConResultMii struct {
	Tag       [2]byte
	TagSize   uint16
	MiiIndex  uint32
	ArtisanId uint32
	_         [82]byte
	Ranking   uint8
	_         uint8
}

func conResult(c *gin.Context) {
	strContestId := c.PostForm("contestno")
	contestId, err := strconv.Atoi(strContestId)
	if err != nil {
		c.Status(http.StatusBadRequest)
		writeResult(c, 400)
		return
	}

	// See if the user voted
	var artisans []string
	vote1 := c.PostForm("craftsno1")
	vote2 := c.PostForm("craftsno2")
	vote3 := c.PostForm("craftsno3")

	if vote2 == "" {
		artisans = append(artisans, vote1)
	} else {
		artisans = append(artisans, vote1, vote2, vote3)

		vote4 := c.PostForm("craftsno4")
		if vote4 != "" {
			artisans = append(artisans, vote4)
		}
	}

	result := ConResult{
		Header: ConResultHeader{
			Tag:       [2]byte{'C', 'R'},
			ContestId: uint32(contestId),
			Padding:   [4]byte{math.MaxUint8, math.MaxUint8, math.MaxUint8, math.MaxUint8},
		},
		Miis: make([]ConResultMii, len(artisans)),
	}

	for i, artisan := range artisans {
		var rank int
		err = pool.QueryRow(ctx, GetMiiRank, artisan, contestId).Scan(&rank)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			writeResult(c, 500)
			return
		}

		intArtisan, err := strconv.Atoi(artisan)
		if err != nil {
			c.Status(http.StatusBadRequest)
			writeResult(c, 400)
			return
		}

		result.Miis[i] = ConResultMii{
			Tag:       [2]byte{'C', 'C'},
			TagSize:   96,
			MiiIndex:  uint32(i + 1),
			ArtisanId: uint32(intArtisan),
			Ranking:   uint8(rank) * 10,
		}
	}

	// Now convert to bytes and send to the channel.
	data := new(bytes.Buffer)
	err = binary.Write(data, binary.BigEndian, result.Header)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		writeResult(c, 500)
		return
	}

	for _, mii := range result.Miis {
		err = binary.Write(data, binary.BigEndian, mii)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			writeResult(c, 500)
			return
		}
	}

	c.Data(http.StatusOK, "application/octet-stream", data.Bytes())
}
