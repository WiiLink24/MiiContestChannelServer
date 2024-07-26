package webpanel

import (
	"bytes"
	"fmt"
	"github.com/WiiLink24/MiiContestChannel/common"
	"github.com/WiiLink24/MiiContestChannel/contest"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	GetContests = `SELECT contest_id, english_name, status, has_souvenir, has_thumbnail, has_special_award, open_time, close_time FROM contests ORDER BY contest_id`
	AddContest  = `INSERT INTO contests (has_special_award, has_thumbnail, has_souvenir, english_name, status, open_time, close_time) 
					VALUES ($1, $2, $3, $4, 'waiting', $5, $6) RETURNING contest_id`
	GetOneContest        = `SELECT contest_id, english_name, status, has_souvenir, has_thumbnail, has_special_award, open_time, close_time FROM contests WHERE contest_id = $1`
	DeleteContest        = `DELETE FROM contests WHERE contest_id = $1`
	DeleteContestEntries = `DELETE FROM contest_miis WHERE contest_id = $1`
	UpdateContest        = `UPDATE contests SET english_name = $1, open_time = $2, close_time = $3, has_special_award = $4, has_thumbnail = $5, has_souvenir = $6 WHERE contest_id = $7`
)

type Contests struct {
	ContestId       int
	Name            string
	Status          string
	HasSouvenir     bool
	HasThumbnail    bool
	HasSpecialAward bool
	OpenTime        time.Time
	CloseTime       time.Time
}

func (w *WebPanel) ViewContests(c *gin.Context) {
	rows, err := w.Pool.Query(w.Ctx, GetContests)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"Error": err.Error(),
		})
		return
	}

	var contests []Contests
	var waitingContests []Contests
	var openContests []Contests
	var judgingContests []Contests
	var resultsContests []Contests
	var closedContests []Contests
	for rows.Next() {
		contest := Contests{}
		err = rows.Scan(&contest.ContestId, &contest.Name, &contest.Status, &contest.HasSouvenir, &contest.HasThumbnail, &contest.HasSpecialAward, &contest.OpenTime, &contest.CloseTime)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "error.html", gin.H{
				"Error": err.Error(),
			})
			return
		}

		contests = append(contests, contest)
		if contest.Status == "waiting" {
			waitingContests = append(waitingContests, contest)
		} else if contest.Status == "open" {
			openContests = append(openContests, contest)
		} else if contest.Status == "judging" {
			judgingContests = append(judgingContests, contest)
		} else if contest.Status == "results" {
			resultsContests = append(resultsContests, contest)
		} else if contest.Status == "closed" {
			closedContests = append(closedContests, contest)
		}
	}

	c.HTML(http.StatusOK, "view_contests.html", gin.H{
		"numberOfContests":        len(contests),
		"Contests":                contests,
		"numberOfWaitingContests":  len(waitingContests),
		"waitingContests":          waitingContests,
		"numberOfOpenContests": len(openContests),
		"openContests":         openContests,
		"numberOfJudgingContests": len(judgingContests),
		"judgingContests":      judgingContests,
		"numberOfResultsContests": len(resultsContests),		
		"resultsContests":      resultsContests,
		"numberOfClosedContests": len(closedContests),
		"closedContests":      closedContests,
	})
}

func (w *WebPanel) AddContest(c *gin.Context) {
	c.HTML(http.StatusOK, "add_contest.html", nil)
}

func (w *WebPanel) AddContestPOST(c *gin.Context) {
	name := c.PostForm("name")
	strSpecialAward := c.PostForm("special_award")
	strOpenTime := c.PostForm("open_time")
	thumbnail, _ := c.FormFile("thumbnail")
	souvenir, _ := c.FormFile("souvenir")

	specialAward := false
	hasThumbnail := false
	hasSouvenir := false
	if thumbnail != nil {
		hasThumbnail = true
	}

	if souvenir != nil {
		hasSouvenir = true
	}

	if strSpecialAward == "on" {
		specialAward = true
	}

	openTime, err := time.Parse("2006-01-02", strOpenTime)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "add_contest.html", gin.H{
			"Error": err.Error(),
		})
		return
	}

	var contestId uint32
	err = w.Pool.QueryRow(w.Ctx, AddContest, specialAward, hasThumbnail, hasSouvenir, name, openTime, openTime.AddDate(0, 0, 7)).Scan(&contestId)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "add_contest.html", gin.H{
			"Error": err.Error(),
		})
		return
	}

	// Write the images now that we have the contest ID.
	if hasThumbnail {
		f, err := thumbnail.Open()
		defer f.Close()
		if err != nil {
			c.HTML(http.StatusInternalServerError, "add_contest.html", gin.H{
				"Error": err.Error(),
			})
			return
		}

		buffer := new(bytes.Buffer)
		_, err = io.Copy(buffer, f)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "add_contest.html", gin.H{
				"Error": err.Error(),
			})
			return
		}

		resized, err := resize(buffer, 96, 96)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "add_contest.html", gin.H{
				"Error": err.Error(),
			})
			return
		}

		// Create encrypted thumbnail
		err = contest.MakePhoto(common.Thumbnail, resized, contestId)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "add_contest.html", gin.H{
				"Error": err.Error(),
			})
			return
		}

		// Write unencrypted thumbnail
		err = os.WriteFile(fmt.Sprintf("%s/contest/%d/thumbnail.jpg", w.Config.AssetsPath, contestId), resized, 0666)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "add_contest.html", gin.H{
				"Error": err.Error(),
			})
			return
		}
	}

	if hasSouvenir {
		f, err := souvenir.Open()
		defer f.Close()
		if err != nil {
			c.HTML(http.StatusInternalServerError, "add_contest.html", gin.H{
				"Error": err.Error(),
			})
			return
		}

		buffer := new(bytes.Buffer)
		_, err = io.Copy(buffer, f)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "add_contest.html", gin.H{
				"Error": err.Error(),
			})
			return
		}

		resized, err := resize(buffer, 512, 384)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "add_contest.html", gin.H{
				"Error": err.Error(),
			})
			return
		}

		// Create encrypted thumbnail
		err = contest.MakePhoto(common.Souvenir, resized, contestId)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "add_contest.html", gin.H{
				"Error": err.Error(),
			})
			return
		}

		// Write unencrypted thumbnail
		err = os.WriteFile(fmt.Sprintf("%s/contest/%d/souvenir.jpg", w.Config.AssetsPath, contestId), resized, 0666)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "add_contest.html", gin.H{
				"Error": err.Error(),
			})
			return
		}
	}

	if hasSouvenir {
		f, err := souvenir.Open()
		defer f.Close()
		if err != nil {
			c.HTML(http.StatusInternalServerError, "add_contest.html", gin.H{
				"Error": err.Error(),
			})
			return
		}

		buffer := new(bytes.Buffer)
		_, err = io.Copy(buffer, f)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "add_contest.html", gin.H{
				"Error": err.Error(),
			})
			return
		}

		resized, err := resize(buffer, 512, 384)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "add_contest.html", gin.H{
				"Error": err.Error(),
			})
			return
		}

		// Create encrypted thumbnail
		err = contest.MakePhoto(common.Souvenir, resized, contestId)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "add_contest.html", gin.H{
				"Error": err.Error(),
			})
			return
		}

		//Write unencrypted souvenir
		err = os.WriteFile(fmt.Sprintf("%s/contest/%d/souvenir.jpg", w.Config.AssetsPath, contestId), resized, 0666)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "add_contest.html", gin.H{
				"Error": err.Error(),
			})
			return
		}
	}

	c.Redirect(http.StatusPermanentRedirect, "/panel/contests#add_success")
}

func (w *WebPanel) EditContest(c *gin.Context) {
	contestId := c.Param("contest_id")
	//fetch the contest data
	row := w.Pool.QueryRow(w.Ctx, GetOneContest, contestId)

	var contestdata Contests
	err := row.Scan(&contestdata.ContestId, &contestdata.Name, &contestdata.Status, &contestdata.HasSouvenir, &contestdata.HasThumbnail, &contestdata.HasSpecialAward, &contestdata.OpenTime, &contestdata.CloseTime)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"Error": err.Error(),
		})
		return
	}

	c.HTML(http.StatusOK, "edit_contest.html", gin.H{
		"ContestInfo": contestdata,
	})
}

func (w *WebPanel) EditContestPOST(c *gin.Context) {
	//Fetch the form data
	name := c.PostForm("name")
	strSpecialAward := c.PostForm("special_award")
	strOpenTime := c.PostForm("open_time")
	thumbnail, _ := c.FormFile("thumbnail")
	souvenir, _ := c.FormFile("souvenir")

	strContestId := c.Param("contest_id")

	var contestId uint32
	contestIdInt, err := strconv.Atoi(strContestId)
	contestId = uint32(contestIdInt)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "edit_contest.html", gin.H{
			"Error": err.Error(),
		})
		return
	}

	//convert the award, thumbnail and souvenir to boolean
	specialAward := false
	hasThumbnail := false
	hasSouvenir := false
	if thumbnail != nil {
		hasThumbnail = true
	}

	if souvenir != nil {
		hasSouvenir = true
	}

	if strSpecialAward == "on" {
		specialAward = true
	}

	//parse the date
	openTime, err := time.Parse("2006-01-02", strOpenTime)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "edit_contest.html", gin.H{
			"Error": err.Error(),
		})
		return
	}

	//update the contest
	_, err = w.Pool.Exec(w.Ctx, UpdateContest, name, openTime, openTime.AddDate(0, 0, 7), specialAward, hasThumbnail, hasSouvenir, strContestId)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "edit_contest.html", gin.H{
			"Error": err.Error(),
		})
		return
	}

	//generate new thumbnail and souvenir
	if hasThumbnail {
		f, err := thumbnail.Open()
		defer f.Close()
		if err != nil {
			c.HTML(http.StatusInternalServerError, "edit_contest.html", gin.H{
				"Error": err.Error(),
			})
			return
		}

		buffer := new(bytes.Buffer)
		_, err = io.Copy(buffer, f)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "edit_contest.html", gin.H{
				"Error": err.Error(),
			})
			return
		}

		resized, err := resize(buffer, 96, 96)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "edit_contest.html", gin.H{
				"Error": err.Error(),
			})
			return
		}

		// Create encrypted thumbnail
		err = contest.MakePhoto(common.Souvenir, resized, contestId)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "edit_contest.html", gin.H{
				"Error": err.Error(),
			})
			return
		}

		// Write unencrypted thumbnail
		err = os.WriteFile(fmt.Sprintf("%s/contest/%d/thumbnail.jpg", w.Config.AssetsPath, contestIdInt), resized, 0666)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "edit_contest.html", gin.H{
				"Error": err.Error(),
			})
			return
		}
	}

	if hasSouvenir {
		f, err := souvenir.Open()
		if err != nil {
			c.HTML(http.StatusInternalServerError, "edit_contest.html", gin.H{
				"Error": err.Error(),
			})
			return
		}

		buffer := new(bytes.Buffer)
		_, err = io.Copy(buffer, f)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "edit_contest.html", gin.H{
				"Error": err.Error(),
			})
			return
		}

		resized, err := resize(buffer, 512, 384)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "edit_contest.html", gin.H{
				"Error": err.Error(),
			})
			return
		}

		//Write encrypted souvenir
		err = contest.MakePhoto(common.Souvenir, resized, contestId)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "edit_contest.html", gin.H{
				"Error": err.Error(),
			})
			return
		}

		//Write unencrypted souvenir
		err = os.WriteFile(fmt.Sprintf("%s/contest/%d/souvenir.jpg", w.Config.AssetsPath, contestId), resized, 0666)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "edit_contest.html", gin.H{
				"Error": err.Error(),
			})
			return
		}
	}

	c.Redirect(http.StatusPermanentRedirect, "/panel/contests/#edit_success")
}

func (w *WebPanel) DeleteContest(c *gin.Context) {
	contestId := c.Param("contest_id")
	//delete the contest entries first
	_, err := w.Pool.Exec(w.Ctx, DeleteContestEntries, contestId)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"Error": err.Error(),
		})
		return
	}
	//delete the contest when the entries are deleted
	_, err = w.Pool.Exec(w.Ctx, DeleteContest, contestId)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"Error": err.Error(),
		})
		return
	}

	c.Redirect(http.StatusFound, "/panel/contests#delete_success")

}
