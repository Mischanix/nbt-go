package nbt

import (
	"fmt"
	"io"
	r "reflect"
)

type tagData tag

func (self *tagData) Type() TagType {
	return self.tagType
}

func (self *tagData) Name() string {
	return self.name
}

func (self *tagData) Parent() TagData {
	return self.parent
}

func loadTagData(rd io.Reader, t *tag) *tagData {
	d := (*tagData)(t)
	switch d.tagType {
	case Byte:
		d.data = readByte(rd)
	case Short:
		d.data = readShort(rd)
	case Int:
		d.data = readInt(rd)
	case Long:
		d.data = readLong(rd)
	case Float:
		d.data = readFloat(rd)
	case Double:
		d.data = readDouble(rd)
	case ByteArray:
		d.data = readByteArray(rd)
	case String:
		d.data = readString(rd)
	case List:
		loadTagList(rd, d)
	case Compound:
		loadTagCompound(rd, d)
	case IntArray:
		d.data = readIntArray(rd)
	}
	return d
}

func loadCompleteTag(rd io.Reader, parent *tagData) *tagData {
	return loadTagData(rd, loadNamedTag(rd, parent))
}

func saveTagData(w io.Writer, d TagData) (err error) {
	switch d.Type() {
	case Byte:
		err = write(w, d.Byte())
	case Short:
		err = write(w, d.Short())
	case Int:
		err = write(w, d.Int())
	case Long:
		err = write(w, d.Long())
	case Float:
		err = write(w, d.Float())
	case Double:
		err = write(w, d.Double())
	case ByteArray:
		err = writeByteArray(w, d.ByteArray())
	case String:
		err = writeString(w, d.StringData())
	case List:
		err = saveTagList(w, d.List())
	case Compound:
		err = saveTagCompound(w, d.Compound())
	case IntArray:
		err = writeIntArray(w, d.IntArray())
	}
	return err
}

func saveCompleteTag(w io.Writer, d TagData) (err error) {
	if err := saveNamedTag(w, d); err != nil {
		return err
	}
	return saveTagData(w, d)
}

func (self *tagData) asValue(tagType TagType) (result r.Value) {
	if self.tagType != tagType {
		tagName := [numTags]string{
			"End", "Byte", "Short", "Int", "Long", "Float", "Double", "ByteArray",
			"String", "List", "Compound", "IntArray",
		}[tagType]
		panic(fmt.Sprintf("%s() called on non-%s TagData",
			tagName, tagName,
		))
	}
	// Value's Kind may be wrong, but this is only indicative of a problem in
	// loadTagData, so a panic can wait until the type assertion.
	return r.ValueOf(self.data)
}

func (self *tagData) Byte() byte {
	return byte(self.asValue(Byte).Int())
}

func (self *tagData) Short() int16 {
	return int16(self.asValue(Short).Int())
}

func (self *tagData) Int() int32 {
	return int32(self.asValue(Int).Int())
}

func (self *tagData) Long() int64 {
	return int64(self.asValue(Long).Int())
}

func (self *tagData) Float() float32 {
	return float32(self.asValue(Float).Float())
}

func (self *tagData) Double() float64 {
	return float64(self.asValue(Double).Float())
}

func (self *tagData) ByteArray() []byte {
	return (self.asValue(ByteArray).Interface()).([]byte)
}

func (self *tagData) StringData() string {
	return (self.asValue(String).Interface()).(string)
}

func (self *tagData) List() TagList {
	return (self.asValue(List).Interface()).(*tagList)
}

func (self *tagData) Compound() TagCompound {
	return (self.asValue(Compound).Interface()).(*tagCompound)
}

func (self *tagData) IntArray() []int32 {
	return (self.asValue(IntArray).Interface()).([]int32)
}

var typeMap = map[r.Type]TagType{
	r.TypeOf(byte(0)):        Byte,
	r.TypeOf(int16(0)):       Short,
	r.TypeOf(int32(0)):       Int,
	r.TypeOf(int64(0)):       Long,
	r.TypeOf(float32(0)):     Float,
	r.TypeOf(float64(0)):     Double,
	r.TypeOf([]byte{}):       ByteArray,
	r.TypeOf(""):             String,
	r.TypeOf(&tagList{}):     List,
	r.TypeOf(&tagCompound{}): Compound,
	r.TypeOf([]int32{}):      IntArray,
}

func (self *tagData) Set(payload interface{}) error {
	if self.parent == nil {
		panic("TagData.Set() cannot be called on parentless tags")
	}
	if self.parent.tagType == List {
		panic("TagData.Set() must not be called directly on list items")
	}
	return self.set(payload)
}

func (self *tagData) set(payload interface{}) error {
	if tagType := payloadType(payload); tagType != invalid {
		self.tagType = tagType
	} else {
		return &payloadTypeError{"TagData.Set"}
	}
	self.data = payload
	return nil
}

func (self *tagData) String() string {
	return fmt.Sprintf("%v", self.data)
}

func payloadType(payload interface{}) TagType {
	if payloadTagType, ok := typeMap[r.TypeOf(payload)]; ok {
		return payloadTagType
	}
	return invalid
}

type payloadTypeError struct {
	Method string
}

func (e *payloadTypeError) Error() string {
	return e.Method + "() called with invalid payload type"
}
