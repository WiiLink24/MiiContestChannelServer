package mii

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kaitai-io/kaitai_struct_go_runtime/kaitai"
	"io"
	"net/http"
	"net/url"
)

func Render(c *gin.Context) {
	inputData := c.Query("data")

	data, err := hex.DecodeString(inputData)
	if err != nil {
		c.JSON(400, MiiError)
		return
	}

	m := NewGen1Wii()
	err = m.Read(kaitai.NewStream(bytes.NewReader(data)), nil, m)
	if err != nil {
		c.JSON(400, MiiError)
		return
	}

	studioMii := CreateStudioMii(m, Wii)

	queryParams := url.Values{
		"data":          {studioMii},
		"type":          {"face"},
		"width":         {"512"},
		"instanceCount": {"1"},
	}

	resp, err := http.Get(fmt.Sprintf("https://studio.mii.nintendo.com/miis/image.png?%s", queryParams.Encode()))
	if err != nil || resp.StatusCode != 200 {
		c.JSON(400, MiiError)
		return
	}

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil || resp.StatusCode != 200 {
		c.JSON(400, MiiError)
		return
	}

	c.Data(http.StatusOK, "image/png", respBytes)
}
