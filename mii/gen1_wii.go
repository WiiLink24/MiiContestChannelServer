// Code generated by kaitai-struct-compiler from a .ksy source file. DO NOT EDIT.

package mii

import (
	"github.com/kaitai-io/kaitai_struct_go_runtime/kaitai"
	"golang.org/x/text/encoding/unicode"
)

type Gen1Wii struct {
	Unknown1           bool
	Gender             bool
	BirthMonth         uint64
	BirthDay           uint64
	FavoriteColor      uint64
	Favorite           bool
	MiiName            string
	BodyHeight         uint8
	BodyWeight         uint8
	AvatarId           []uint8
	ClientId           []uint8
	FaceType           uint64
	FaceColor          uint64
	FacialFeature      uint64
	Unknown2           uint64
	Mingle             bool
	Unknown3           bool
	Downloaded         bool
	HairType           uint64
	HairColor          uint64
	HairFlip           bool
	Unknown4           uint64
	EyebrowType        uint64
	Unknown5           bool
	EyebrowRotation    uint64
	Unknown6           uint64
	EyebrowColor       uint64
	EyebrowSize        uint64
	EyebrowVertical    uint64
	EyebrowHorizontal  uint64
	EyeType            uint64
	Unknown7           uint64
	EyeRotation        uint64
	EyeVertical        uint64
	EyeColor           uint64
	Unknown8           bool
	EyeSize            uint64
	EyeHorizontal      uint64
	Unknown9           uint64
	NoseType           uint64
	NoseSize           uint64
	NoseVertical       uint64
	Unknown10          uint64
	MouthType          uint64
	MouthColor         uint64
	MouthSize          uint64
	MouthVertical      uint64
	GlassesType        uint64
	GlassesColor       uint64
	Unknown11          bool
	GlassesSize        uint64
	GlassesVertical    uint64
	FacialHairMustache uint64
	FacialHairBeard    uint64
	FacialHairColor    uint64
	FacialHairSize     uint64
	FacialHairVertical uint64
	MoleEnable         bool
	MoleSize           uint64
	MoleVertical       uint64
	MoleHorizontal     uint64
	Unknown12          bool
	CreatorName        string
	_io                *kaitai.Stream
	_root              *Gen1Wii
	_parent            interface{}
}

func NewGen1Wii() *Gen1Wii {
	return &Gen1Wii{}
}

func (this *Gen1Wii) Read(io *kaitai.Stream, parent interface{}, root *Gen1Wii) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp1, err := this._io.ReadBitsIntBe(1)
	if err != nil {
		return err
	}
	this.Unknown1 = tmp1 != 0
	tmp2, err := this._io.ReadBitsIntBe(1)
	if err != nil {
		return err
	}
	this.Gender = tmp2 != 0
	tmp3, err := this._io.ReadBitsIntBe(4)
	if err != nil {
		return err
	}
	this.BirthMonth = tmp3
	tmp4, err := this._io.ReadBitsIntBe(5)
	if err != nil {
		return err
	}
	this.BirthDay = tmp4
	tmp5, err := this._io.ReadBitsIntBe(4)
	if err != nil {
		return err
	}
	this.FavoriteColor = tmp5
	tmp6, err := this._io.ReadBitsIntBe(1)
	if err != nil {
		return err
	}
	this.Favorite = tmp6 != 0
	this._io.AlignToByte()
	tmp7, err := this._io.ReadBytes(int(20))
	if err != nil {
		return err
	}
	tmp7 = tmp7
	tmp8, err := kaitai.BytesToStr(tmp7, unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM).NewDecoder())
	if err != nil {
		return err
	}
	this.MiiName = tmp8
	tmp9, err := this._io.ReadU1()
	if err != nil {
		return err
	}
	this.BodyHeight = tmp9
	tmp10, err := this._io.ReadU1()
	if err != nil {
		return err
	}
	this.BodyWeight = tmp10
	for i := 0; i < int(4); i++ {
		_ = i
		tmp11, err := this._io.ReadU1()
		if err != nil {
			return err
		}
		this.AvatarId = append(this.AvatarId, tmp11)
	}
	for i := 0; i < int(4); i++ {
		_ = i
		tmp12, err := this._io.ReadU1()
		if err != nil {
			return err
		}
		this.ClientId = append(this.ClientId, tmp12)
	}
	tmp13, err := this._io.ReadBitsIntBe(3)
	if err != nil {
		return err
	}
	this.FaceType = tmp13
	tmp14, err := this._io.ReadBitsIntBe(3)
	if err != nil {
		return err
	}
	this.FaceColor = tmp14
	tmp15, err := this._io.ReadBitsIntBe(4)
	if err != nil {
		return err
	}
	this.FacialFeature = tmp15
	tmp16, err := this._io.ReadBitsIntBe(3)
	if err != nil {
		return err
	}
	this.Unknown2 = tmp16
	tmp17, err := this._io.ReadBitsIntBe(1)
	if err != nil {
		return err
	}
	this.Mingle = tmp17 != 0
	tmp18, err := this._io.ReadBitsIntBe(1)
	if err != nil {
		return err
	}
	this.Unknown3 = tmp18 != 0
	tmp19, err := this._io.ReadBitsIntBe(1)
	if err != nil {
		return err
	}
	this.Downloaded = tmp19 != 0
	tmp20, err := this._io.ReadBitsIntBe(7)
	if err != nil {
		return err
	}
	this.HairType = tmp20
	tmp21, err := this._io.ReadBitsIntBe(3)
	if err != nil {
		return err
	}
	this.HairColor = tmp21
	tmp22, err := this._io.ReadBitsIntBe(1)
	if err != nil {
		return err
	}
	this.HairFlip = tmp22 != 0
	tmp23, err := this._io.ReadBitsIntBe(5)
	if err != nil {
		return err
	}
	this.Unknown4 = tmp23
	tmp24, err := this._io.ReadBitsIntBe(5)
	if err != nil {
		return err
	}
	this.EyebrowType = tmp24
	tmp25, err := this._io.ReadBitsIntBe(1)
	if err != nil {
		return err
	}
	this.Unknown5 = tmp25 != 0
	tmp26, err := this._io.ReadBitsIntBe(4)
	if err != nil {
		return err
	}
	this.EyebrowRotation = tmp26
	tmp27, err := this._io.ReadBitsIntBe(6)
	if err != nil {
		return err
	}
	this.Unknown6 = tmp27
	tmp28, err := this._io.ReadBitsIntBe(3)
	if err != nil {
		return err
	}
	this.EyebrowColor = tmp28
	tmp29, err := this._io.ReadBitsIntBe(4)
	if err != nil {
		return err
	}
	this.EyebrowSize = tmp29
	tmp30, err := this._io.ReadBitsIntBe(5)
	if err != nil {
		return err
	}
	this.EyebrowVertical = tmp30
	tmp31, err := this._io.ReadBitsIntBe(4)
	if err != nil {
		return err
	}
	this.EyebrowHorizontal = tmp31
	tmp32, err := this._io.ReadBitsIntBe(6)
	if err != nil {
		return err
	}
	this.EyeType = tmp32
	tmp33, err := this._io.ReadBitsIntBe(2)
	if err != nil {
		return err
	}
	this.Unknown7 = tmp33
	tmp34, err := this._io.ReadBitsIntBe(3)
	if err != nil {
		return err
	}
	this.EyeRotation = tmp34
	tmp35, err := this._io.ReadBitsIntBe(5)
	if err != nil {
		return err
	}
	this.EyeVertical = tmp35
	tmp36, err := this._io.ReadBitsIntBe(3)
	if err != nil {
		return err
	}
	this.EyeColor = tmp36
	tmp37, err := this._io.ReadBitsIntBe(1)
	if err != nil {
		return err
	}
	this.Unknown8 = tmp37 != 0
	tmp38, err := this._io.ReadBitsIntBe(3)
	if err != nil {
		return err
	}
	this.EyeSize = tmp38
	tmp39, err := this._io.ReadBitsIntBe(4)
	if err != nil {
		return err
	}
	this.EyeHorizontal = tmp39
	tmp40, err := this._io.ReadBitsIntBe(5)
	if err != nil {
		return err
	}
	this.Unknown9 = tmp40
	tmp41, err := this._io.ReadBitsIntBe(4)
	if err != nil {
		return err
	}
	this.NoseType = tmp41
	tmp42, err := this._io.ReadBitsIntBe(4)
	if err != nil {
		return err
	}
	this.NoseSize = tmp42
	tmp43, err := this._io.ReadBitsIntBe(5)
	if err != nil {
		return err
	}
	this.NoseVertical = tmp43
	tmp44, err := this._io.ReadBitsIntBe(3)
	if err != nil {
		return err
	}
	this.Unknown10 = tmp44
	tmp45, err := this._io.ReadBitsIntBe(5)
	if err != nil {
		return err
	}
	this.MouthType = tmp45
	tmp46, err := this._io.ReadBitsIntBe(2)
	if err != nil {
		return err
	}
	this.MouthColor = tmp46
	tmp47, err := this._io.ReadBitsIntBe(4)
	if err != nil {
		return err
	}
	this.MouthSize = tmp47
	tmp48, err := this._io.ReadBitsIntBe(5)
	if err != nil {
		return err
	}
	this.MouthVertical = tmp48
	tmp49, err := this._io.ReadBitsIntBe(4)
	if err != nil {
		return err
	}
	this.GlassesType = tmp49
	tmp50, err := this._io.ReadBitsIntBe(3)
	if err != nil {
		return err
	}
	this.GlassesColor = tmp50
	tmp51, err := this._io.ReadBitsIntBe(1)
	if err != nil {
		return err
	}
	this.Unknown11 = tmp51 != 0
	tmp52, err := this._io.ReadBitsIntBe(3)
	if err != nil {
		return err
	}
	this.GlassesSize = tmp52
	tmp53, err := this._io.ReadBitsIntBe(5)
	if err != nil {
		return err
	}
	this.GlassesVertical = tmp53
	tmp54, err := this._io.ReadBitsIntBe(2)
	if err != nil {
		return err
	}
	this.FacialHairMustache = tmp54
	tmp55, err := this._io.ReadBitsIntBe(2)
	if err != nil {
		return err
	}
	this.FacialHairBeard = tmp55
	tmp56, err := this._io.ReadBitsIntBe(3)
	if err != nil {
		return err
	}
	this.FacialHairColor = tmp56
	tmp57, err := this._io.ReadBitsIntBe(4)
	if err != nil {
		return err
	}
	this.FacialHairSize = tmp57
	tmp58, err := this._io.ReadBitsIntBe(5)
	if err != nil {
		return err
	}
	this.FacialHairVertical = tmp58
	tmp59, err := this._io.ReadBitsIntBe(1)
	if err != nil {
		return err
	}
	this.MoleEnable = tmp59 != 0
	tmp60, err := this._io.ReadBitsIntBe(4)
	if err != nil {
		return err
	}
	this.MoleSize = tmp60
	tmp61, err := this._io.ReadBitsIntBe(5)
	if err != nil {
		return err
	}
	this.MoleVertical = tmp61
	tmp62, err := this._io.ReadBitsIntBe(5)
	if err != nil {
		return err
	}
	this.MoleHorizontal = tmp62
	tmp63, err := this._io.ReadBitsIntBe(1)
	if err != nil {
		return err
	}
	this.Unknown12 = tmp63 != 0
	this._io.AlignToByte()
	tmp64, err := this._io.ReadBytes(int(20))
	if err != nil {
		return err
	}
	tmp64 = tmp64
	tmp65, err := kaitai.BytesToStr(tmp64, unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM).NewDecoder())
	if err != nil {
		return err
	}
	this.CreatorName = tmp65
	return err
}
