package instanceid

import (
	"strconv"
	"strings"

	"github.com/xlab/treeprint"
)

// This package supports the handling of instance-identification fields
// as proposed in Theos paper

// An instance is represented through it MIID
// MIID := <sN> "/" <vN> ["/" <vA>] "%" <t>s
// Example:
//		msA/1.1/feature-branch-2345abcd%222s
// The complete call-graph including it's own MIID
// is represented by:
// CIID := MIID [ "(" UIDs+ ")"]
// UIDs := CIID [ "+" CIID ]+

// This package provides some helpers to work with this
// type of instance-identification

// CIID := MIID [ "(" UIDs+ ")"]
// UIDs := CIID [ "+" CIID ]+
// MIID := <sN> "/" <vN> ["/" <vA>] "%" <t>s

// Ciid represents the complete call-graph as instance-id
type Ciid struct {
	Miid  Miid
	Ciids []Ciid
}

// Miid represents the instance only by it's name, version, additional information
// and epoch time
type Miid struct {
	Sn string
	Vn string
	Va string
	T  int
}

// PrintCiid prints a tree representation of the complete call-graph of a
// Ciid
func PrintCiid(ciid Ciid) string {
	tree := treeprint.New()
	tree = ciid.visitCiid(tree)
	return tree.String()
}

// NewCiid creates a new Ciid from a string in the form of
// Sn1/Vn1/Va1%t1s(Sn2/Vn2/Va2%t2s+Sn3/Vn3/Va3%t3s(Sn4/Vn4/Va4%t4s))
func NewCiid(id string) (ciid Ciid) {
	return parseCiid(id)
}

// NewMiid creates a new Miid from a string in the of
// Sn1/Vn1/Va1%t1s
// in case a Ciid is being provided the Miid part is only
// returned
// If there are syntax errors an empty Miid will be returned
func NewMiid(id string) (miid Miid) {
	return parseMIID(id)
}

func (c *Ciid) String() string {
	sB := strings.Builder{}
	sB.WriteString(c.Miid.String())
	if len(c.Ciids) > 0 {
		sB.WriteString("(")
		for i, a := range c.Ciids {
			sB.WriteString(a.String())
			if i+1 < len(c.Ciids) {
				sB.WriteString("+")
			}
		}
		sB.WriteString(")")
	}

	return sB.String()
}

func (m *Miid) String() string {
	sB := strings.Builder{}

	sB.WriteString(m.Sn)
	if m.Vn != "" {
		sB.WriteString("/" + m.Vn)
	}
	if m.Va != "" {
		sB.WriteString("/" + m.Va)
	}
	if m.T != 0 {
		sB.WriteString("%" + strconv.Itoa(m.T) + "s")
	}
	return sB.String()
}

func (m *Miid) metadata() string {
	sB := strings.Builder{}

	if m.Vn != "" {
		sB.WriteString(m.Vn)
	}
	if m.Va != "" {
		sB.WriteString("/" + m.Va)
	}
	if m.T != 0 {
		sB.WriteString("%" + strconv.Itoa(m.T) + "s")
	}
	return sB.String()
}

func parseCiid(id string) (ciid Ciid) {
	name, arg := seperateFNameFromArg(id)

	if arg == "" {
		return Ciid{Miid: parseMIID(name)}
	}
	me := Ciid{Miid: parseMIID(name)}
	me.Ciids = parseArguments(arg)
	return me
}

func parseArguments(arg string) (ciids []Ciid) {
	ss := splitOnPlus(arg)

	for _, a := range ss {
		ciids = append(ciids, parseCiid(a))
	}
	return ciids
}

func (c Ciid) visitCiid(t treeprint.Tree) treeprint.Tree {
	x := t.AddBranch(c.Miid.Sn + "/" + c.Miid.Vn)
	if c.Miid.metadata() != "" {
		x.SetMetaValue(strconv.Itoa(c.Miid.T) + "s")
	}
	for _, s := range c.Ciids {
		s.visitCiid(x)
	}
	return t
}

func parseMIID(id string) (miid Miid) {
	s := strings.SplitN(id, "/", -1)
	l := len(s)
	var r Miid
	if l >= 1 {
		r.Sn = s[0]
	}
	if l == 2 {
		e := strings.Split(s[1], "%")
		r.Vn = e[0]
		if len(e) > 1 {
			t, _ := strconv.Atoi(strings.Split(e[1], "s")[0])
			r.T = t
		}
	} else if l >= 2 {
		r.Vn = s[1]
	}

	if l == 3 {
		e := strings.Split(s[2], "%")
		r.Va = e[0]
		t, _ := strconv.Atoi(strings.Split(e[1], "s")[0])
		r.T = t
	}

	return r
}

func seperateFNameFromArg(signature string) (name, arg string) {
	n := strings.Builder{}
	a := strings.Builder{}
	var inArgs bool
	var count int
	for _, s := range signature {
		if s == '(' {
			count++
			inArgs = true
		}
		if s == ')' {
			count--
		}

		if !inArgs {
			n.WriteRune(s)
		} else if count == 1 && s != '(' {
			a.WriteRune(s)
		} else if count > 1 {
			a.WriteRune(s)
		}
	}

	return n.String(), a.String()
}

func deleteEmpty(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}

func splitOnPlus(s string) (ss []string) {
	var openClose int
	var splitPos []int

	// first find parenthis pairs
	for pos, char := range s {
		if char == '(' {
			openClose++
		} else if char == ')' {
			openClose--
		} else if char == '+' {
			if openClose == 0 {
				splitPos = append(splitPos, pos)
			}
		}
	}

	//   split arguments
	if len(splitPos) > 0 {
		oldPos := -1
		for _, s2 := range splitPos {
			ss = append(ss, s[oldPos+1:s2])
			oldPos = s2
		}
		ss = append(ss, s[oldPos+1:])
	} else {
		ss = append(ss, s)
	}
	return deleteEmpty(ss)
}
