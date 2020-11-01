package pterm

import (
	"strings"
)

// TreeNode is used as items in a Tree.
type TreeNode struct {
	Children []TreeNode
	Text     string
}

// LeveledList is a list, which contains multiple LeveledListItem.
type LeveledList []LeveledListItem

// LeveledListItem combines a text with a specific level.
// The level is the indent, which would normally be seen in a BulletList.
type LeveledListItem struct {
	Level int
	Text  string
}

// DefaultTree contains standards, which can be used to render a Tree.
var DefaultTree = Tree{
	TreeStyle:            &ThemeDefault.TreeStyle,
	TextStyle:            &ThemeDefault.TreeTextStyle,
	TopRightCornerString: "└",
	HorizontalString:     "─",
	TopRightDownString:   "├",
	VerticalString:       "│",
	RightDownLeftString:  "┬",
	Indent:               2,
}

// Tree is able to render a list.
type Tree struct {
	Root                 TreeNode
	TreeStyle            *Style
	TextStyle            *Style
	TopRightCornerString string
	TopRightDownString   string
	HorizontalString     string
	VerticalString       string
	RightDownLeftString  string
	Indent               int
}

// WithTreeStyle returns a new list with a specific tree style.
func (p Tree) WithTreeStyle(style *Style) *Tree {
	p.TreeStyle = style
	return &p
}

// WithTextStyle returns a new list with a specific text style.
func (p Tree) WithTextStyle(style *Style) *Tree {
	p.TextStyle = style
	return &p
}

// WithTopRightCornerString returns a new list with a specific TopRightCornerString.
func (p Tree) WithTopRightCornerString(s string) *Tree {
	p.TopRightCornerString = s
	return &p
}

// WithTopRightDownStringOngoing returns a new list with a specific TopRightDownString.
func (p Tree) WithTopRightDownStringOngoing(s string) *Tree {
	p.TopRightDownString = s
	return &p
}

// WithHorizontalString returns a new list with a specific HorizontalString.
func (p Tree) WithHorizontalString(s string) *Tree {
	p.HorizontalString = s
	return &p
}

// WithVerticalString returns a new list with a specific VerticalString.
func (p Tree) WithVerticalString(s string) *Tree {
	p.VerticalString = s
	return &p
}

// WithRoot returns a new list with a specific Root.
func (p Tree) WithRoot(root TreeNode) *Tree {
	p.Root = root
	return &p
}

// WithIndent returns a new list with a specific amount of spacing between the levels.
// Indent must be at least 1.
func (p Tree) WithIndent(indent int) *Tree {
	if indent < 1 {
		indent = 1
	}
	p.Indent = indent
	return &p
}

// Render prints the list to the terminal.
func (p Tree) Render() error {
	s, err := p.Srender()
	if err != nil {
		return err
	}
	Println(s)

	return nil
}

// Srender renders the list as a string.
func (p Tree) Srender() (string, error) {
	if p.TreeStyle == nil {
		p.TreeStyle = NewStyle()
	}
	if p.TextStyle == nil {
		p.TextStyle = NewStyle()
	}

	return walkOverTree(p.Root.Children, p, ""), nil
}

// walkOverTree is a recursive function,
// which analyzes a Tree and connects the items with specific characters.
// Returns Tree as string.
func walkOverTree(list []TreeNode, p Tree, prefix string) string {
	var ret string
	for i, item := range list {
		if len(list) > i+1 { // if not last in list
			if len(item.Children) == 0 { // if there are no children
				ret += prefix + p.TreeStyle.Sprint(p.TopRightDownString) + strings.Repeat(p.TreeStyle.Sprint(p.HorizontalString), p.Indent) +
					p.TextStyle.Sprint(item.Text) + "\n"
			} else { // if there are children
				ret += prefix + p.TreeStyle.Sprint(p.TopRightDownString) + strings.Repeat(p.TreeStyle.Sprint(p.HorizontalString), p.Indent-1) +
					p.TreeStyle.Sprint(p.RightDownLeftString) + p.TextStyle.Sprint(item.Text) + "\n"
				ret += walkOverTree(item.Children, p, prefix+p.TreeStyle.Sprint(p.VerticalString)+strings.Repeat(" ", p.Indent-1))
			}
		} else if len(list) == i+1 { // if last in list
			if len(item.Children) == 0 { // if there are no children
				ret += prefix + p.TreeStyle.Sprint(p.TopRightCornerString) + strings.Repeat(p.TreeStyle.Sprint(p.HorizontalString), p.Indent) +
					p.TextStyle.Sprint(item.Text) + "\n"
			} else { // if there are children
				ret += prefix + p.TreeStyle.Sprint(p.TopRightCornerString) + strings.Repeat(p.TreeStyle.Sprint(p.HorizontalString), p.Indent-1) +
					p.TreeStyle.Sprint(p.RightDownLeftString) + p.TextStyle.Sprint(item.Text) + "\n"
				ret += walkOverTree(item.Children, p, prefix+strings.Repeat(" ", p.Indent))
			}
		}
	}
	return ret
}

// NewTreeFromLeveledList converts a TreeItems list to a TreeNode and returns it.
func NewTreeFromLeveledList(leveledListItems LeveledList) TreeNode {
	if len(leveledListItems) == 0 {
		return TreeNode{}
	}

	root := &TreeNode{
		Children: []TreeNode{},
		Text:     leveledListItems[0].Text,
	}

	for i, record := range leveledListItems {
		last := root

		if record.Level < 0 {
			record.Level = 0
			leveledListItems[i].Level = 0
		}

		if len(leveledListItems)-1 != i {
			if leveledListItems[i+1].Level-1 > record.Level {
				leveledListItems[i+1].Level = record.Level + 1
			}
		}

		for i := 0; i < record.Level; i++ {
			lastIndex := len(last.Children) - 1
			last = &last.Children[lastIndex]
		}
		last.Children = append(last.Children, TreeNode{
			Children: []TreeNode{},
			Text:     record.Text,
		})
	}

	return *root
}
