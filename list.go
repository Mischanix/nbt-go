package nbt

import (
	"fmt"
	"io"
)

type tagList struct {
	tag      *tagData
	listType TagType
	list     []TagData
}

func (self *tagList) Type() TagType {
	return self.tag.tagType
}

func (self *tagList) Name() string {
	return self.tag.name
}

func (self *tagList) Parent() TagData {
	return self.tag.parent
}

func loadTagList(rd io.Reader, t *tagData) {
	l := &tagList{}
	l.tag = t
	l.listType = TagType(readByte(rd))
	length := readInt(rd)
	l.list = make([]TagData, length)
	for i := int32(0); i < length; i++ {
		t := &tag{"", l.listType, l.tag, nil}
		l.list[i] = loadTagData(rd, t)
	}
	t.data = l
}

func saveTagList(w io.Writer, l TagList) error {
	if err := write(w, int8(l.ListType())); err != nil {
		return err
	}
	if err := write(w, l.Length()); err != nil {
		return err
	}
	for i := int32(0); i < l.Length(); i++ {
		if err := saveTagData(w, l.At(i)); err != nil {
			return err
		}
	}
	return nil
}

func (self *tagList) ListType() TagType {
	return self.listType
}

func (self *tagList) Length() int32 {
	return int32(len(self.list))
}

func (self *tagList) At(index int32) TagData {
	return self.list[index]
}

func (self *tagList) Set(index int32, payload interface{}) error {
	if index >= self.Length() || index < 0 {
		return &indexError{"TagList.Set"}
	}
	if payloadType(payload) != self.listType {
		panic(&payloadTypeError{"TagList.Set"})
	}
	return self.list[index].set(payload)
}

func (self *tagList) Add(payload interface{}) error {
	if payloadType(payload) != self.listType {
		panic(&payloadTypeError{"TagList.Add"})
	}
	t := &tagData{"", self.listType, self.tag, payload}
	self.list = append(self.list, t)
	return nil
}

func (self *tagList) Remove(index int32) error {
	if index >= self.Length() || index < 0 {
		return &indexError{"TagList.Set"}
	}
	self.list = append(self.list[0:index], self.list[index+1:]...)
	return nil
}

func (self *tagList) String() string {
	return fmt.Sprintf("%v", self.list)
}

type indexError struct {
	Method string
}

func (e *indexError) Error() string {
	return e.Method + "() called with invalid index"
}
