package llis

import (
	"encoding/gob"
	"io"
	"os"
)

type Tag struct {
	Name            string
	Parent          *Tag
	Children        map[string]*Tag
	AssociatedNodes map[string]*Node
}

type Node struct {
	Name string
	Tags map[string]*Tag
}

type Database struct {
	Name  string
	Tags  map[string]*Tag
	Nodes map[string]*Node

	f *os.File
}

/*
Exports to the old file, closes it, then changes to the new file and exports to that

new can be nil (will disable saving the database to a file)
*/
func (self *Database) SwitchFile(new *os.File) error {
	self.CloseFile()
	self.f = new
	err := ExportDatabase(*self, self.f)
	if err != nil {
		return err
	}
	return nil
}

func (self *Database) CloseFile() error {
	if self.f != nil {
		err := ExportDatabase(*self, self.f)
		if err != nil {
			return err
		}
		self.f.Close()
	}
	return nil
}

func (self *Database) NewNode(name string) (node *Node, alreadyExists bool) {
	if n, ok := self.Nodes[name]; ok {
		return n, ok
	}
	n := Node{
		Name: name,
		Tags: map[string]*Tag{},
	}
	self.Nodes[name] = &n
	return &n, false
}

// parent can be nil
func (self *Database) NewTag(name string, parent *Tag) (tag *Tag, alreadyExists bool) {
	if t, ok := self.Tags[name]; ok {
		return t, ok
	}
	t := Tag{
		Name:            name,
		Parent:          parent,
		Children:        map[string]*Tag{},
		AssociatedNodes: map[string]*Node{},
	}
	if parent != nil {
		parent.Children[name] = &t
	}
	self.Tags[name] = &t
	return &t, false
}

func (self *Database) AddTagToNode(tagName, nodeName string) (tagExists, nodeExists bool) {
	t, okt := self.Tags[tagName]
	n, okn := self.Nodes[nodeName]
	if okt && okn {
		self.Nodes[nodeName].Tags[tagName] = t
		t.AssociatedNodes[nodeName] = n
	}
	return okt, okn
}

func (self *Database) GetQuerier() Querier {
	return Querier{db: self}
}

var Databases = map[string]*Database{}

func NewDatabase(name string) *Database {
	db := Database{
		Name:  name,
		Tags:  map[string]*Tag{},
		Nodes: map[string]*Node{},
	}
	Databases[name] = &db
	return &db
}

type backupIntermediate_Tag struct {
	// tag name
	Parent string
}

type backupIntermediate_Node struct {
	// tag nams
	Tags []string
}

type backupIntermediate_Database struct {
	Name  string
	Tags  map[string]backupIntermediate_Tag
	Nodes map[string]backupIntermediate_Node
}

/*
Exports a database to dst

If you want a raw []byte, use a bytes.Buffer for dst
*/
func ExportDatabase(db Database, dst io.Writer) error {
	in := backupIntermediate_Database{
		Name:  db.Name,
		Tags:  map[string]backupIntermediate_Tag{},
		Nodes: map[string]backupIntermediate_Node{},
	}

	for name, t := range db.Tags {
		tag_in := backupIntermediate_Tag{}
		if t.Parent != nil {
			tag_in.Parent = t.Parent.Name
		}
		in.Tags[name] = tag_in
	}

	for name, n := range db.Nodes {
		node_in := backupIntermediate_Node{
			Tags: []string{},
		}
		for tname := range n.Tags {
			node_in.Tags = append(node_in.Tags, tname)
		}
		in.Nodes[name] = node_in
	}

	enc := gob.NewEncoder(dst)
	return enc.Encode(in)
}

/*
!!! WARNING: THIS FUNCTION IS NOT USABLE YET !!! (TODO)

Imports a database that was exported with llis.ExportDatabase()

If you are importing directly from a []byte, use bytes.NewReader(b) for src
*/
func ImportDatabase(src io.Reader) (*Database, error) {
	var in backupIntermediate_Database
	enc := gob.NewDecoder(src)
	err := enc.Decode(&in)
	if err != nil {
		return nil, err
	}
	db := Database{Name: in.Name}

	/*
		This will first look for tags without a parent and put them into
		orphans, while putting the others into look4parents.

	*/

	// TODO: everything below here

	look4parents := map[string]backupIntermediate_Tag{}
	orphans := map[string]*Tag{}
	for name, it := range in.Tags {
		t := Tag{Name: name, Parent: nil}
		if it.Parent == "" {
			orphans[name] = &t
		} else {
			look4parents[name] = it
		}
	}

	return &db, nil
}

// !!! WARNING: THIS FUNCTION IS NOT USABLE YET !!!
func LoadDatabase(fpath string) (*Database, error) {
	f, err := os.Open(fpath)
	if err != nil {
		return nil, err
	}
	db, err := ImportDatabase(f)
	return db, err
}

/////////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////// database functions end here ////////////////////////////////
/////////////////////////////////////////////////////////////////////////////////////////////

/*
The Querier makes queries to the database using tag/node names.

It shouldn't be hand-initialized, but rather made using Database.GetQuerier()
*/
type Querier struct {
	db *Database
}

func (self *Querier) Select_AllTagsOfNode(nodeName string) (tags []string, nodeExists bool) {
	n, nodeExists := self.db.Nodes[nodeName]
	if !nodeExists {
		return
	}
	tags = make([]string, len(n.Tags))
	i := 0
	for name := range n.Tags {
		tags[i] = name
		i++
	}
	return
}

// Not to be confused with Querier.Select_AllNodesWithTags
func (self *Querier) Select_AllNodesWithTag(tagName string) (nodes []string, tagExists bool) {
	t, tagExists := self.db.Tags[tagName]
	if !tagExists {
		return
	}
	nodes = make([]string, len(t.AssociatedNodes))
	i := 0
	for name := range t.AssociatedNodes {
		nodes[i] = name
		i++
	}
	return
}

/*
Not to be confused with Querier.Select_AllNodesWithTag.

If tagNames is empty, nodes and nonExistentTag will be empty.

If the query fails due to a tag not existing, it will return the tag name in nonExistentTag.
If all tags exist, nonExistentTag will be empty.

tagNames[0] should preferably be a tag with a low len(AssociatedNodes) for (potentially significantly) enhanced performance. (TODO: optimize)
*/
func (self *Querier) Select_AllNodesWithTags(tagNames []string) (nodes []string, nonExistentTag string) { // TODO
	if len(tagNames) == 0 {
		return []string{}, ""
	}

	t, ok := self.db.Tags[tagNames[0]]
	if !ok {
		return nil, tagNames[0]
	}
	possible := t.AssociatedNodes
	for i := 1; i < len(tagNames); i++ {
		t, ok = self.db.Tags[tagNames[i]]
		if !ok {
			return nil, tagNames[i]
		}
		for posblName := range possible {
			if _, ok := t.AssociatedNodes[posblName]; !ok {
				delete(possible, posblName)
			}
		}
	}

	nodes = make([]string, len(possible))
	i := 0
	for name := range possible {
		nodes[i] = name
		i++
	}
	return
}
