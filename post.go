package main

import (
	"encoding/binary"
	"encoding/hex"
	"github.com/WiiLink24/nwc24"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

const (
	DoesArtisanExist = `SELECT EXISTS(SELECT 1 FROM artisans WHERE artisan_id = $1 AND wii_number = $2 AND mac_address = $3)`

	InsertMii = `INSERT INTO miis 
    				(artisan_id, initials, skill, nickname, gender, country_id, wii_number, mii_data, mii_id)
					VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
					RETURNING entry_id`

	UpdateMii = `UPDATE miis 
					SET initials = $1, skill = $2, nickname = $3, gender = $4, country_id = $5, wii_number = $6, mii_data = $7, mii_id = $8 
            		WHERE artisan_id = $9 AND entry_id = $10`

	UpdatePostCount = `UPDATE artisans SET number_of_posts = number_of_posts + 1, last_post = $1 WHERE artisan_id = $2`
)

func post(c *gin.Context) {
	macAddress := c.PostForm("macadr")
	strEntryNumber := c.PostForm("entryno")
	strWiiNumber := c.PostForm("wiino")
	country := c.PostForm("country")
	gender := c.PostForm("sex")
	skill := c.PostForm("skill")
	name := c.PostForm("nickname")
	miiId := c.PostForm("miiid")
	artisanId := c.PostForm("craftsno")
	initials := c.PostForm("initial")
	miiData := []byte(c.PostForm("miidata"))

	// Validate Wii Number
	wiiNumber, err := strconv.ParseUint(strWiiNumber, 10, 64)
	if err != nil {
		writeResult(c, 400)
		return
	}

	number := nwc24.LoadWiiNumber(wiiNumber)
	if !number.CheckWiiNumber() {
		writeResult(c, 106)
		return
	}

	// Validate the artisan
	var exists bool
	err = pool.QueryRow(ctx, DoesArtisanExist, artisanId, wiiNumber, macAddress).Scan(&exists)
	if err != nil {
		writeResult(c, 500)
		return
	}

	if !exists {
		// Invalid artisan
		writeResult(c, 105)
		return
	}

	// Validate Mii data
	if len(miiData) != 76 {
		writeResult(c, 109)
		return
	}

	// Validate nickname
	if len(name) < 1 || len(name) > 10 {
		writeResult(c, 108)
		return
	}

	entryId, err := strconv.Atoi(strEntryNumber)
	if err != nil {
		// TODO: Figure out invalid entry ID
		// It should never happen if this was sent by MCC, but people will be people.
		writeResult(c, 108)
		return
	}

	// The Mii ID is a byte string which should be converted to bytes as it will save storage space by half.
	miiIdBytes, err := hex.DecodeString(miiId)
	if err != nil {
		// TODO: Figure out proper error for this
		writeResult(c, 108)
		return
	}

	if entryId != 0 {
		// Mii exists, update it.
		_, err = pool.Exec(ctx, UpdateMii, initials, skill, name, gender, country, wiiNumber, miiData, miiIdBytes, artisanId, entryId)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			writeResult(c, 500)
			return
		}
	} else {
		err = pool.QueryRow(ctx, InsertMii, artisanId, initials, skill, name, gender, country, wiiNumber, miiData, miiIdBytes).Scan(&entryId)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			writeResult(c, 500)
			return
		}
	}

	// Update the artisan's number of posts
	_, err = pool.Exec(ctx, UpdatePostCount, time.Now(), artisanId)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		writeResult(c, 500)
		return
	}

	data, _ := hex.DecodeString("505300000000000000000000000000000000000000000000FFFFFFFFFFFFFFFF454E000C00000001")
	data = binary.BigEndian.AppendUint32(data, uint32(entryId))
	c.Data(http.StatusOK, "application/octet-stream", data)
}
