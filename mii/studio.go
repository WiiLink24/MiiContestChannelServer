package mii

import (
	"bytes"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"github.com/kaitai-io/kaitai_struct_go_runtime/kaitai"
	"log"
	"net/http"
	"reflect"
)

const (
	Wii        = "wii"
	Gen2       = "gen2"       // gen2_wiiu_3ds_miitomo, CFLStoreData/FFLStoreData/AFLStoreData
	SwitchDB   = "switch"     // gen3_sdb, nn::mii::detail::StoreDataRaw
	SwitchGame = "switchgame" // gen3_swi, nn::mii::CharInfo
	Switch     = "switch"     // generic used for comparisons
	MiiStudio  = "studio"
)

var (
	MiiError = gin.H{
		"mii":    "000f165d65747981878c8e94a3a8b4acb3b7bec5cbd2dfe6f5f8fffc040b1c232b313f3a405a60697178808f9198a4",
		"studio": "Error!",
	}

	Wrinkles = map[uint8]uint8{4: 5, 5: 2, 6: 3, 7: 7, 8: 8, 10: 9, 11: 11}

	Makeup = map[uint8]uint8{1: 1, 2: 6, 3: 9, 9: 10}
)

/*
This is kept for compatibility reasons.
*/

// StudioResponse represents a response from studio.cgi.
type StudioResponse struct {
	// Certain fields are not used on Mii Studio and not included
	Mii           string `json:"mii"`
	MiiStudio     string `json:"miistudio"`
	Name          string `json:"name"`
	CreatorName   string `json:"creator_name"`
	//Birthday      string `json:"birthday"`
	FavoriteColor string `json:"favorite_color"`
	Height        uint8  `json:"height"`
	Build         uint8  `json:"build"`
	Gender        string `json:"gender"`
	//Mingle        string `json:"mingle"`
	//Copying       string `json:"copying"`
}

func Studio(c *gin.Context) {
	inputType := c.PostForm("platform")
	inputData, _ := c.FormFile("data")

	f, err := inputData.Open()
	defer f.Close()
	if err != nil {
		c.JSON(400, MiiError)
		return
	}

	var mii any

	switch inputType {
	case Wii:
		m := NewGen1Wii()
		err = m.Read(kaitai.NewStream(f), nil, m)
		if err != nil {
			c.JSON(400, MiiError)
			return
		}
		mii = m
	case Gen2:
		m := NewGen2Wiiu3dsMiitomo()
		err = m.Read(kaitai.NewStream(f), nil, m)
		if err != nil {
			c.JSON(400, MiiError)
			return
		}
		mii = m
	case SwitchDB:
		m := NewMiidataSdb()
		err = m.Read(kaitai.NewStream(f), nil, m)
		if err != nil {
			c.JSON(400, MiiError)
			return
		}
		// SwitchDB specific hack
		m.EyebrowVertical += 3
		mii = m
	case SwitchGame:
		m := NewMiidataSwi()
		err = m.Read(kaitai.NewStream(f), nil, m)
		if err != nil {
			c.JSON(400, MiiError)
			return
		}
		mii = m
	case MiiStudio:
		m := NewGen3Studio()
		err = m.Read(kaitai.NewStream(f), nil, m)
		if err != nil {
			c.JSON(400, MiiError)
			return
		}
		mii = m

	default:
		// that input type does not exist
		c.JSON(400, MiiError)
		return
	}

	// for comparison purposes
	if inputType == SwitchDB || inputType == SwitchGame || inputType == MiiStudio {
		inputType = Switch
	}

	response := CreateStudioMii(mii, inputType)

	c.JSON(http.StatusOK, response)
}

type ctx struct {
	n         int
	mii       any
	Buffer    *bytes.Buffer
	BufferRaw *bytes.Buffer
}

// used in response
// https://github.com/RiiConnect24/CMOC-Server/blob/master/miicontestp/cgi-bin/studio.cgi
var genderStringMap = []string{ "Male", "Female", }
var favoriteColorMap = []string{
    "Red",
    "Orange",
    "Yellow",
    "Lime Green",
    "Forest Green",
    "Royal Blue",
    "Sky Blue",
    "Pink",
    "Purple",
    "Brown",
    "White",
    "Black",
}
// for mingle and copying
//var yesNoNAMap = []string{ "Yes", "No", "N/A" }

// CreateStudioMii converts any Mii into the format the Nintendo Mii Studio expects.
func CreateStudioMii(mii any, miiType string) StudioResponse {
	c := &ctx{
		n:         256,
		Buffer:    new(bytes.Buffer),
		BufferRaw: new(bytes.Buffer), // for mii studio data
		mii:       mii,
	}

	// Init ctx
	c.Buffer.WriteByte(0)

	if miiType != Switch {
		if c.getField("FacialHairColor") == 0 {
			c.writeValue(8)
		} else {
			c.writeValue(c.getField("FacialHairColor"))
		}
	} else {
		c.writeValue(c.getField("FacialHairColor"))
	}

	c.writeValue(c.getField("FacialHairBeard"))
	bodyWeight := c.getField("BodyWeight")
	c.writeValue(bodyWeight)

	if miiType == Wii {
		c.writeValue(3)
	} else {
		c.writeValue(c.getField("EyeStretch"))
	}

	if miiType != Switch {
		c.writeValue(c.getField("EyeColor") + 8)
	} else {
		c.writeValue(c.getField("EyeColor"))
	}

	c.writeValue(c.getField("EyeRotation"))
	c.writeValue(c.getField("EyeSize"))
	c.writeValue(c.getField("EyeType"))
	c.writeValue(c.getField("EyeHorizontal"))
	c.writeValue(c.getField("EyeVertical"))

	if miiType == Wii {
		c.writeValue(3)
	} else {
		c.writeValue(c.getField("EyebrowStretch"))
	}

	if miiType != Switch {
		if c.getField("EyebrowColor") == 0 {
			c.writeValue(8)
		} else {
			c.writeValue(c.getField("EyebrowColor"))
		}
	} else {
		c.writeValue(c.getField("EyebrowColor"))
	}

	c.writeValue(c.getField("EyebrowRotation"))
	c.writeValue(c.getField("EyebrowSize"))
	c.writeValue(c.getField("EyebrowType"))
	c.writeValue(c.getField("EyebrowHorizontal"))

	c.writeValue(c.getField("EyebrowVertical"))

	c.writeValue(c.getField("FaceColor"))

	if miiType == Wii {
		if v, ok := Makeup[c.getField("FacialFeature")]; ok {
			c.writeValue(v)
		} else {
			c.writeValue(0)
		}
	} else {
		c.writeValue(c.getField("FaceMakeup"))
	}

	c.writeValue(c.getField("FaceType"))

	if miiType == Wii {
		if v, ok := Wrinkles[c.getField("FacialFeature")]; ok {
			c.writeValue(v)
		} else {
			c.writeValue(0)
		}
	} else {
		c.writeValue(c.getField("FaceWrinkles"))
	}

	favoriteColor := c.getField("FavoriteColor")
	c.writeValue(favoriteColor)
	gender := c.getField("Gender")
	c.writeValue(gender)

	if miiType != Switch {
		if c.getField("GlassesColor") == 0 {
			c.writeValue(8)
		} else if c.getField("GlassesColor") < 6 {
			c.writeValue(c.getField("GlassesColor") + 13)
		} else {
			c.writeValue(0)
		}
	} else {
		c.writeValue(c.getField("GlassesColor"))
	}

	c.writeValue(c.getField("GlassesSize"))
	c.writeValue(c.getField("GlassesType"))
	c.writeValue(c.getField("GlassesVertical"))

	if miiType != Switch {
		if c.getField("HairColor") == 0 {
			c.writeValue(8)
		} else {
			c.writeValue(c.getField("HairColor"))
		}
	} else {
		c.writeValue(c.getField("HairColor"))
	}

	c.writeValue(c.getField("HairFlip"))
	c.writeValue(c.getField("HairType"))
	bodyHeight := c.getField("BodyHeight")
	c.writeValue(bodyHeight)
	c.writeValue(c.getField("MoleSize"))
	c.writeValue(c.getField("MoleEnable"))
	c.writeValue(c.getField("MoleHorizontal"))
	c.writeValue(c.getField("MoleVertical"))

	if miiType == Wii {
		c.writeValue(3)
	} else {
		c.writeValue(c.getField("MouthStretch"))
	}

	if miiType != Switch {
		if c.getField("MouthColor") < 5 {
			c.writeValue(c.getField("MouthColor") + 19)
		} else {
			c.writeValue(0)
		}
	} else {
		c.writeValue(c.getField("MouthColor"))
	}

	c.writeValue(c.getField("MouthSize"))
	c.writeValue(c.getField("MouthType"))
	c.writeValue(c.getField("MouthVertical"))
	c.writeValue(c.getField("FacialHairSize"))
	c.writeValue(c.getField("FacialHairMustache"))
	c.writeValue(c.getField("FacialHairVertical"))
	c.writeValue(c.getField("NoseSize"))
	c.writeValue(c.getField("NoseType"))
	c.writeValue(c.getField("NoseVertical"))

	response := StudioResponse{
		Mii:           hex.EncodeToString(c.Buffer.Bytes()),
		MiiStudio:     hex.EncodeToString(c.BufferRaw.Bytes()),
		Height:        bodyHeight,
		Build:         bodyWeight,
		Gender:        genderStringMap[gender],
		FavoriteColor: favoriteColorMap[favoriteColor],
	}

	// set extra response parameters
	r := reflect.ValueOf(c.mii)
	i := reflect.Indirect(r)
	nameField := i.FieldByName("MiiName")
	if nameField.IsValid() {
		response.Name = nameField.Interface().(string)
	}
	creatorNameField := i.FieldByName("CreatorName")
	if creatorNameField.IsValid() {
		response.CreatorName = creatorNameField.Interface().(string)
	}

	return response
}

// getField is an unsafe way to dynamically get a field value from all Mii types.
// Although unsafe in theory, in practice all fields are guaranteed to exist.
// It also converts the value into the one Nintendo expects for the field.
func (c *ctx) getField(field string) uint8 {
	r := reflect.ValueOf(c.mii)
	i := reflect.Indirect(r)

	f := i.FieldByName(field)

	if !f.IsValid() {
		// this may be a method?
		m := r.MethodByName(field)
		if m.IsValid() {
			results := m.Call(nil)
			if len(results) > 0 {
				switch results[0].Kind() {
				case reflect.Bool:
					if results[0].Bool() {
						return 1
					} else {
						return 0
					}
				default:
					uint8Type := reflect.TypeOf((uint8)(0))
					convToUint8 := results[0].Convert(uint8Type)
					return convToUint8.Interface().(uint8)
				}
			}
		} else {
			log.Println("WARNING: Cannot fetch", field, "as field OR method from type", i.Type().Name())
			return 0
		}
	}

	if f.Kind() == reflect.Bool {
		if f.Bool() {
			return 1
		}

		return 0
	}

	return uint8(f.Uint())
}

func (c *ctx) writeValue(v uint8) {
	c.BufferRaw.WriteByte(byte(v))
	eo := (7 + (int(v) ^ c.n)) % 256
	c.n = eo
	c.Buffer.WriteByte(byte(eo))
}
