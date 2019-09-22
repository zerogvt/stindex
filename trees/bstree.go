// Package trees implements Binary Search Trees for
// search terms (i.e. strings)
package trees

import (
	"fmt"
	"log"
	"os"
)

// SearchTerm is a representation of a search term datum
type SearchTerm struct {
	Term  string //search term
	Times int64  //how many times it appears
}

// ToStr returns a string representation of the search term
func (stA *SearchTerm) ToStr() string {
	return fmt.Sprintf("[%s, %d]", stA.Term, stA.Times)
}

// Compare compares the caller SearchTerm with the argument and
// returns -1, 0, 1 if the caller is correspondingly less, equal, greater
// than the argument.
func (stA *SearchTerm) Compare(stB SearchTerm) int {
	if stA.Term < stB.Term {
		return -1
	} else if stA.Term == stB.Term {
		return 0
	} else {
		return 1
	}
}

// Node is a representation of a tree node
type Node struct {
	left  *Node
	right *Node
	data  SearchTerm
}

// BST is a representation of a Binary Search Tree
type BST struct {
	root *Node
}

// Insert inserts a new data string into a Binary Search Tree
// If the data string has not been seen before it will be added in
// the proper position in BST.
// If it is already in the tree the insertion will just update the
// counter of the representation that is already present.
func (t *BST) Insert(term string) *BST {
	data := SearchTerm{Term: term, Times: 1}
	if t.root == nil {
		t.root = &Node{data: data, left: nil, right: nil}
	} else {
		t.root.insert(data)
	}
	return t
}

// insert (pkg private) a new node with SearchTerm data 'data' under
// node n following the BST rules recursivelly.
// If the SearchTerm term already exists then increase its 'times' value.
func (n *Node) insert(data SearchTerm) {
	if n == nil {
		return
	} else if data.Compare(n.data) < 0 {
		if n.left == nil {
			n.left = &Node{data: data, left: nil, right: nil}
		} else {
			n.left.insert(data)
		}
	} else if data.Compare(n.data) > 0 {
		if n.right == nil {
			n.right = &Node{data: data, left: nil, right: nil}
		} else {
			n.right.insert(data)
		}
	} else if data.Compare(n.data) == 0 {
		n.data.Times++
	}
}

// TraverseInOrder traverses binary search tree in order,
// i.e. left child, parent, right child,
// thus creating showing the data in the tree in order.
func (t *BST) TraverseInOrder() chan SearchTerm {
	ch := make(chan SearchTerm)
	go func() {
		t.root.traverseInOrder(ch)
		close(ch)
	}()
	return ch
}

// inorder traversal of the nodes under node n
func (n *Node) traverseInOrder(ch chan SearchTerm) {
	if n == nil {
		return
	}
	n.left.traverseInOrder(ch)
	ch <- n.data
	n.right.traverseInOrder(ch)
}

// UnzipToFile traverses a BST in order and unzips the found search terms.
// I.e. a Search Term {"beer", 2} will be expanded to 2 lines of "beer".
// It is used to write down the final output file.
// IMPORTANT NOTE: If the file exists the new data will be appended to it.
func (t *BST) UnzipToFile(fpath string) {
	// open file in append mode
	f, err := os.OpenFile(fpath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	// traverse tree and expand each search term
	for node := range t.TraverseInOrder() {
		var i int64
		for i = 0; i < node.Times; i++ {
			tmp := fmt.Sprintf("%s\n", node.Term)
			if _, err := f.Write([]byte(tmp)); err != nil {
				log.Fatal(err)
			}
		}
	}
	// close file
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}
