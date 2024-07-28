package main

import (
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

const (
	DoesArtisanExistMacAddress = `SELECT EXISTS(SELECT 1 FROM artisans WHERE mac_address = $1 AND artisan_id = $2)`
	DoesMiiExist               = `SELECT EXISTS(SELECT 1 FROM miis WHERE entry_id = $1)`
	DidAlreadyVote             = `SELECT EXISTS(SELECT 1 FROM votes WHERE artisan_id = $1 AND voted_mii = $2)`
	UpdateMiiVote              = `UPDATE miis SET likes = likes + 1, perm_likes = perm_likes + 1 WHERE entry_id = $1`
	UpdateArtisanVotes         = `UPDATE artisans SET total_likes = total_likes + 1 WHERE artisan_id = $1`
	InsertVote                 = `INSERT INTO votes (artisan_id, voted_mii) VALUES ($1, $2)`
)

func vote(c *gin.Context) {
	macAddress := c.PostForm("macadr")
	strEntryNumber := c.PostForm("entryno")
	artisanId := c.PostForm("craftsno")

	// Validate the artisan.
	var exists bool
	err := pool.QueryRow(ctx, DoesArtisanExistMacAddress, macAddress, artisanId).Scan(&exists)
	if err != nil {
		writeResult(c, 500)
		return
	}

	if !exists {
		// Invalid artisan.
		writeResult(c, 105)
		return
	}

	entryId, err := strconv.Atoi(strEntryNumber)
	if err != nil {
		// TODO: Figure out invalid entry ID
		// It should never happen if this was sent by MCC, but people will be people.
		writeResult(c, 108)
		return
	}

	// Next make sure the Mii exists.
	err = pool.QueryRow(ctx, DoesMiiExist, entryId).Scan(&exists)
	if err != nil {
		writeResult(c, 500)
		return
	}

	if !exists {
		// Mii does not exist.
		data, _ := hex.DecodeString("565400000000000000000000000000000000000000000000ffffffffffffffff454e002000000001000000010000000100000001000000010000000100000001")
		c.Data(http.StatusOK, "application/octet-stream", data)
		return
	}

	// Finally, we will make sure this artisan has not voted.
	err = pool.QueryRow(ctx, DidAlreadyVote, artisanId, entryId).Scan(&exists)
	if err != nil {
		writeResult(c, 500)
		return
	}

	if exists {
		// Duplicate vote
		data, _ := hex.DecodeString("565400000000000000000000000000000000000000000000ffffffffffffffff454e002000000001000000010000000100000001000000010000000100000001")
		c.Data(http.StatusOK, "application/octet-stream", data)
		return
	}

	// Now we can insert the vote
	_, err = pool.Exec(ctx, UpdateMiiVote, entryId)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		writeResult(c, 500)
		return
	}

	// TODO: Artisans have popularity, set it based on their likes
	_, err = pool.Exec(ctx, UpdateArtisanVotes, artisanId)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		writeResult(c, 500)
		return
	}

	_, err = pool.Exec(ctx, InsertVote, artisanId, entryId)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		writeResult(c, 500)
		return
	}

	data, _ := hex.DecodeString("565400000000000000000000000000000000000000000000ffffffffffffffff454e002000000001000000010000000100000001000000010000000100000001")
	c.Data(http.StatusOK, "application/octet-stream", data)
}
