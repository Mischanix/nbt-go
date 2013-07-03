package nbt

import (
	"fmt"
	"io"
)

type tagCompound struct {
	tag      *tagData
	children map[string]TagData
}

func (self *tagCompound) Type() TagType {
	return self.tag.tagType
}

func (self *tagCompound) Name() string {
	return self.tag.name
}

func (self *tagCompound) Parent() TagData {
	return self.tag.parent
}

func loadTagCompound(rd io.Reader, t *tagData) {
	c := &tagCompound{}
	c.tag = t
	c.children = make(map[string]TagData)
	for {
		child := loadCompleteTag(rd, t)
		if child.tagType == End {
			break
		}
		c.children[child.name] = child
	}
	t.data = c
}

func (self *tagCompound) At(name string) TagData {
	if child, ok := self.children[name]; ok {
		return child
	}
	return nil
}

func (self *tagCompound) Each(f func(name string, child TagData) bool) {
	for name, child := range self.children {
		if result := f(name, child); result == false {
			break
		}
	}
}

func (self *tagCompound) Path(pathString string) (result TagData) {
	self.Each(func(name string, child TagData) bool {
		if pathString == name {
			result = child
		} else if len(pathString) > len(name) && // Avoid slice size panics
			// pathString == {name}/...
			pathString[:len(name)] == name && pathString[len(name):][0] == '/' {

			if child.Type() != Compound {
				panic(&childNotACompoundError{name, pathString})
			} else {
				result = child.Compound().Path(pathString[len(name)+1:])
			}
		}
		return result == nil
	})
	if result != nil {
		return result
	} else {
		panic(&childNotFoundError{"Path", pathString})
	}
}

func (self *tagCompound) Set(name string, payload interface{}) error {
	if child := self.At(name); child != nil {
		child.Set(payload)
	} else {
		tagType := payloadType(payload)
		if tagType == invalid {
			panic(&payloadTypeError{"TagCompound.Set"})
		}
		self.children[name] = &tagData{name, tagType, self.tag, payload}
	}
	return nil
}

func (self *tagCompound) Remove(name string) error {
	if _, ok := self.children[name]; ok {
		delete(self.children, name)
		return nil
	} else {
		return &childNotFoundError{"TagCompound.Remove", name}
	}
}

func saveTagCompound(w io.Writer, c TagCompound) (err error) {
	c.Each(func(name string, child TagData) bool {
		if err = saveCompleteTag(w, child); err != nil {
			return false
		} else {
			return true
		}
	})
	if err != nil {
		return err
	} else {
		return write(w, []byte{0}) // End tag
	}
}

func (self *tagCompound) Save(w io.Writer) error {
	if err := saveNamedTag(w, self); err != nil {
		return err
	}
	return saveTagCompound(w, self)
}

func (self *tagCompound) String() string {
	return fmt.Sprintf("%s: %v", self.Name(), self.children)
}

type childNotFoundError struct {
	Method string
	Child  string
}

func (e *childNotFoundError) Error() string {
	return "'" + e.Child + "' not found by " + e.Method + "()"
}

type childNotACompoundError struct {
	Child string
	Path  string
}

func (e *childNotACompoundError) Error() string {
	return "Child '" + e.Child + "' was found by TagCompound.Path(\"" + e.Path +
		"\") but is not a TagCompound"
}
