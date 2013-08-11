package nbt

import (
	"io"
)

type tag struct {
	name    string
	tagType TagType
	parent  *tagData
	data    interface{}
}

func (self *tag) Type() TagType {
	return self.tagType
}

func (self *tag) Name() string {
	return self.name
}

func (self *tag) Parent() TagData {
	return self.parent
}

func createTag(name string, tagType TagType, parent *tagData) *tag {
	return &tag{name, tagType, parent, nil}
}

func loadNamedTag(rd io.Reader, parent *tagData) *tag {
	t := &tag{}
	t.parent = parent
	t.tagType = TagType(readByte(rd))
	if t.tagType != End {
		t.name = readString(rd)
	}
	return t
}

func saveNamedTag(w io.Writer, t Tag) error {
	if err := write(w, byte(t.Type())); err != nil {
		return err
	}
	return writeString(w, t.Name())
}
