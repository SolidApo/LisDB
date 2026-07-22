package llis_test

import (
	"slices"
	"testing"

	"codeberg.org/SolidApo/lisdb/llis"
)

// compare unordered string slices
func uStrSliceCmp(s1 []string, s2 []string) bool {
	slices.Sort(s1)
	slices.Sort(s2)
	return slices.Compare(s1, s2) == 0
}

func TestDatabase(t *testing.T) {
	db := llis.NewDatabase("testdb")

	db.NewTag("TA", nil)
	db.NewTag("TB", nil)
	db.NewTag("TC", nil)

	db.NewNode("N1")
	db.NewNode("N2")
	db.NewNode("N3")
	db.NewNode("N4")
	db.NewNode("N5")
	db.NewNode("N6")
	db.NewNode("N7")
	db.NewNode("N8")

	// 1 to 8
	db.AddTagToNode("TA", "N1")
	db.AddTagToNode("TA", "N2")
	db.AddTagToNode("TA", "N3")
	db.AddTagToNode("TA", "N4")
	db.AddTagToNode("TA", "N5")
	db.AddTagToNode("TA", "N6")
	db.AddTagToNode("TA", "N7")
	db.AddTagToNode("TA", "N8")

	// 1 to 4
	db.AddTagToNode("TB", "N1")
	db.AddTagToNode("TB", "N2")
	db.AddTagToNode("TB", "N3")
	db.AddTagToNode("TB", "N4")

	// primes
	db.AddTagToNode("TC", "N2")
	db.AddTagToNode("TC", "N3")
	db.AddTagToNode("TC", "N5")
	db.AddTagToNode("TC", "N7")

	q := db.GetQuerier()

	nodes, _ := q.Select_AllNodesWithTag("TB")
	if !(len(nodes) == 4) { // assume that its correct if theres 4 etnries
		t.Fatal("ERROR: `nodes, _ := q.Select_AllNodesWithTag(\"TB\")`\n\t\t-> len(nodes) != 4")
	}

	nodes, _ = q.Select_AllNodesWithTags([]string{"TB", "TC"})
	if !uStrSliceCmp(nodes, []string{"N2", "N3"}) {
		t.Fatal("ERROR: `nodes, _ = q.Select_AllNodesWithTags([]string{\"TB\", \"TC\"})`\n\t\t-> nodes != [N2 N3]")
	}
}
