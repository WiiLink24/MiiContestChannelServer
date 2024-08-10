package main

import (
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

const (
	IsContestJudging               = `SELECT EXISTS (SELECT 1 FROM contests WHERE contest_id = $1 AND status = 'judging')`
	DoesArtisanExistMacAddressNoId = `SELECT EXISTS(SELECT 1 FROM artisans WHERE mac_address = $1)`
	IsFirstVote                    = `SELECT COUNT(*) = 0 FROM contest_votes WHERE contest_id = $1 AND mac_address = $2`
	InsertContestVote              = `INSERT INTO contest_votes (contest_id, vote_1, vote_2, vote_3, mac_address) 
							VALUES ($1, $2, $3, $4, $5)`
	UpdateContestVote = `UPDATE contest_votes SET vote_1 = $1, vote_2 = $2, vote_3 = $3 WHERE contest_id = $4 AND mac_address = $5`
	GetContestVotes   = `SELECT vote_1, vote_2, vote_3 FROM contest_votes WHERE contest_id = $1 AND mac_address = $2`
)

func conVote(c *gin.Context) {
	contestId := c.PostForm("contestno")
	votes := strings.Split(c.PostForm("craftsno1"), ",")
	macAddress := c.PostForm("macadr")

	// Check if the contest is open for judging
	var isContestJudging bool
	err := pool.QueryRow(ctx, IsContestJudging, contestId).Scan(&isContestJudging)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		writeResult(c, 500)
		return
	}

	if !isContestJudging {
		data, _ := hex.DecodeString("435600000000000000000000000000000000000000000000FFFFFFFFFFFFFFFF")
		c.Data(http.StatusOK, "application/octet-stream", data)
		return
	}

	// Now check if this is a registered artisan
	var artisanExists bool
	err = pool.QueryRow(ctx, DoesArtisanExistMacAddressNoId, macAddress).Scan(&artisanExists)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		writeResult(c, 500)
		return
	}

	if !artisanExists {
		c.Status(http.StatusBadRequest)
		writeResult(c, 704)
		return
	}

	// Check if this is the artisans first time voting
	var isFirstVote bool
	err = pool.QueryRow(ctx, IsFirstVote, contestId, macAddress).Scan(&isFirstVote)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		writeResult(c, 500)
		return
	}

	if isFirstVote {
		_, err = pool.Exec(ctx, InsertContestVote, contestId, votes[0], votes[1], votes[2], macAddress)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			writeResult(c, 500)
			return
		}
	} else {
		previous := make([]int, 3)
		err = pool.QueryRow(ctx, GetContestVotes, contestId, macAddress).Scan(&previous[0], &previous[1], &previous[2])
		if err != nil {
			c.Status(http.StatusInternalServerError)
			writeResult(c, 500)
			return
		}

		for i := 0; i < 3; i++ {
			// Pad the votes
			if len(votes) < i+1 {
				votes = append(votes, strconv.Itoa(previous[i]))
			}
		}

		_, err = pool.Exec(ctx, UpdateContestVote, votes[0], votes[1], votes[2], contestId, macAddress)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			writeResult(c, 500)
			return
		}
	}

	data, _ := hex.DecodeString("435600000000000000000000000000000000000000000000FFFFFFFFFFFFFFFF4E4C001000000001FFFFFFFFFFFFFFFF")
	c.Data(http.StatusOK, "application/octet-stream", data)
}
