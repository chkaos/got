package types

import (
	"bytes"
	"encoding/gob"
	"sort"
)

// TreeObject represents a directory.
type TreeObject struct {
	Files []DirectoryEntry
}

// DirectoryEntry represents one entry in a TreeObject
type DirectoryEntry struct {
	Name string
	Hash Hash
}

// ByFileName implements sort.Interface for []Person based on
// the Age field.
type ByFileName []DirectoryEntry

func (a ByFileName) Len() int           { return len(a) }
func (a ByFileName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByFileName) Less(i, j int) bool { return a[i].Name < a[j].Name }

// NewTreeObject creates a new TreeObject
func NewTreeObject() TreeObject {
	return TreeObject{
		Files: make([]DirectoryEntry, 0),
	}
}

// AddFile adds a file to the given TreeObject
func (tree TreeObject) AddFile(name string, hash Hash) {
	tree.Files = append(tree.Files, DirectoryEntry{
		Name: name,
		Hash: hash,
	})
	sort.Sort(ByFileName(tree.Files))
}

// HasFile returns true if the tree has a file of the given name
func (tree TreeObject) HasFile(name string) bool {
	for _, entry := range tree.Files {
		if entry.Name == name {
			return true
		}
	}
	return false
}

// Serialize creates a representation of the current TreeObject.
func (tree TreeObject) Serialize() []byte {
	if len(tree.Files) > 0 {
		buffer := new(bytes.Buffer)
		e := gob.NewEncoder(buffer)

		err := e.Encode(tree)
		if err != nil {
			panic(err)
		}
		return buffer.Bytes()
	}

	// Empty directory corresponds to empty bytes
	return []byte{}
}

// DeserializeTreeObject returns a commit object deserialized from the given
// byte stream
func DeserializeTreeObject(input []byte) TreeObject {
	buffer := bytes.NewBuffer(input)
	dec := gob.NewDecoder(buffer)

	var tree TreeObject
	err := dec.Decode(&tree)
	if err != nil {
		panic(err)
	}
	return tree
}