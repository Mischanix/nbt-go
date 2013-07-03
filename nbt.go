// Package nbt provides an interface for reading, modifying, and writing
// Named Binary Tag (NBT) data atop Go streams.  Many methods in nbt will
// cause runtime panics if an invalid state occurs in favor of returning
// a meaningless nil and having that nil cause a runtime panic by being
// dereferenced elsewhere.
package nbt

import (
	"fmt"
	"io"
)

type TagType int8

const (
	// Marks the end of a Compound.
	End TagType = iota
	// int8
	Byte
	// int16
	Short
	// int32
	Int
	// int64
	Long
	// float32
	Float
	// float64
	Double
	// length-prefixed []byte
	ByteArray
	// length-prefixed string
	String
	// length-prefixed list of Tags
	List
	// End-delimited list of named tags
	Compound
	// length-prefixed []int32
	IntArray

	numTags int     = iota
	invalid TagType = -1
)

type Tag interface {
	Type() TagType
	// Name returns nil for unnamed tags.
	Name() string
	// Parent returns the compound or list for which this Tag is a child or
	// nil if it is a root tag
	Parent() TagData
}

// TagData represents a Tag of any type and provides accessors for the tag's
// payload data.  These will panic if the payload is not of the appropriate
// TagType.
type TagData interface {
	Tag
	Byte() int8
	Short() int16
	Int() int32
	Long() int64
	Float() float32
	Double() float64
	ByteArray() []int8
	StringData() string
	List() TagList
	Compound() TagCompound
	IntArray() []int32

	// Set attempts to modify the Tag's payload.  It will fail if the supplied
	// payload is of invalid type.  It will panic if the Tag's parent is a List.
	Set(payload interface{}) error

	// Used by TagList.Set
	set(payload interface{}) error
}

// TagList represents a List Tag and provides methods to add, access, modify,
// and remove list items.
type TagList interface {
	Tag
	ListType() TagType
	Length() int32
	// At returns the Tag at index.  It will panic on invalid index.
	At(index int32) TagData

	// Set replaces the tag at index with payload.  It will panic on invalid
	// payload type.  It will return a non-nil error on invalid index.
	Set(index int32, payload interface{}) error
	// Add appends payload to the list and resizes the list.  It will panic on
	// invalid payload type.
	Add(payload interface{}) error
	// Remove will remove the tag at index and resize the list.  It will return
	// a non-nil error on invalid index.
	Remove(index int32) error
}

// TagCompound represents a Compound Tag and provides methods to access, add,
// and modify the Tag's children.
type TagCompound interface {
	Tag
	// At returns the compound's child with Name name or nil.
	At(name string) TagData
	// Each iterates over the compound's children, calling the callback f on
	// each name,child pair.  If f returns false, Each will stop iterating.
	Each(f func(name string, child TagData) bool)
	// Path will iterate over the compound's children for names and do a simple
	// substring match against the pathString.  If it matches, and if there is
	// a '/' character directly after the match, it will follow.  If no suitable
	// children are found, it will panic.  This occurs recursively until the
	// pathString is emptied.
	Path(pathString string) TagData

	// Set will add a new tag or modify an existing tag on the compound to set
	// the child with Name name to the supplied payload.  If a payload with
	// invalid type is supplied, Set will panic.
	Set(name string, payload interface{}) error
	// If Remove cannot find the child, it will return a non-nil error.
	Remove(name string) error
	// Save attempts to write the compound to the writer.
	Save(w io.Writer) error
}

// Load deserializes a root tag from the supplied io.Reader.  Load will Read()
// the io.Reader until the end of the root tag or until an error occurs.
func Load(rd io.Reader) (root TagCompound, err error) {
	return loadCompleteTag(rd, nil).Compound(), nil
}

// Try wraps the supplied function in a deferred recover() that will put any
// panic result into the return value.  This allows closing interaction with
// nbt into a function while also not having to handle the panics directly.
func Try(f func()) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = &tryError{e}
		}
	}()
	f()
	return err
}

type tryError struct {
	Message interface{}
}

func (e *tryError) Error() string {
	return fmt.Sprintf("%v", e.Message)
}

// MakeRoot creates a new empty, parentless compound tag.
func MakeRoot(name string) (root TagCompound) {
	c := &tagCompound{nil, make(map[string]TagData)}
	tag := &tagData{name, Compound, nil, c}
	c.tag = tag
	return c
}
