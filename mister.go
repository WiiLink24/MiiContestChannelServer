package main

import (
	"encoding/binary"
	"encoding/hex"
	"github.com/WiiLink24/nwc24"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"unicode/utf8"
)

const InsertArtisan = `INSERT INTO artisans 
    					(name, country_id, wii_number, mac_address, mii_data) 
						VALUES ($1, $2, $3, $4, $5)
						RETURNING artisan_id`

func mister(c *gin.Context) {
	name := c.PostForm("nickname")
	country := c.PostForm("country")
	miiData := []byte(c.PostForm("miidata"))
	macAddress := c.PostForm("macadr")
	strwiiNumber := c.PostForm("wiino")

	// Validate Wii Number
	wiiNumber, err := strconv.ParseUint(strwiiNumber, 10, 64)
	if err != nil {
		writeResult(c, 400)
		return
	}

	number := nwc24.LoadWiiNumber(wiiNumber)
	if !number.CheckWiiNumber() {
		writeResult(c, 310)
		return
	}

	// Validate Mii data
	if len(miiData) != 76 {
		writeResult(c, 305)
		return
	}

	// Validate nickname
	// We are required to use this function as non UTF-8 characters will go over the cap.
	strSize := utf8.RuneCountInString(name)
	if strSize < 1 || strSize > 10 {
		writeResult(c, 304)
		return
	}

	var artisanId uint32
	err = pool.QueryRow(ctx, InsertArtisan, name, country, wiiNumber, macAddress, miiData).Scan(&artisanId)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		writeResult(c, 500)
		return
	}

	data, _ := hex.DecodeString("4D5300000000000000000000000000000000000000000000FFFFFFFFFFFFFFFF454E000C00000001")
	data = binary.BigEndian.AppendUint32(data, artisanId)
	c.Data(http.StatusOK, "application/octet-stream", data)
}
