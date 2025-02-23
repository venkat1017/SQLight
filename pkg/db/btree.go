package db

import (
	"sqlight/pkg/interfaces"
)

const (
	NodeTypeLeaf = iota
	NodeTypeInternal
)

const (
	LeafNodeMaxRecords  = 3 // Small for testing
	InternalNodeMaxKeys = 3
)

// Node represents a B+ tree node
type Node struct {
	IsLeaf   bool
	Keys     []int
	Records  []*interfaces.Record
	Children []*Node
	Next     *Node
	Parent   *Node
}

// BTree represents a B+ tree
type BTree struct {
	Root *Node
}

// NewBTree creates a new B+ tree
func NewBTree() *BTree {
	return &BTree{
		Root: &Node{
			IsLeaf:   true,
			Keys:     make([]int, 0),
			Records:  make([]*interfaces.Record, 0),
			Children: nil,
		},
	}
}

// Insert adds a new record to the B+ tree
func (t *BTree) Insert(id int, record *interfaces.Record) error {
	if t.Root == nil {
		t.Root = &Node{
			IsLeaf:   true,
			Keys:     make([]int, 0),
			Records:  make([]*interfaces.Record, 0),
			Children: nil,
		}
	}

	node := t.Root
	// Find the leaf node where this record should be inserted
	for !node.IsLeaf {
		pos := node.findPosition(id)
		node = node.Children[pos]
	}

	// Insert into leaf node
	t.insertIntoLeaf(node, id, record)

	// Check if we need to split
	if len(node.Keys) > LeafNodeMaxRecords {
		t.splitLeaf(node)
	}

	return nil
}

// Helper methods
func (n *Node) findPosition(key int) int {
	for i, k := range n.Keys {
		if key <= k {
			return i
		}
	}
	return len(n.Keys)
}

func (t *BTree) insertIntoLeaf(node *Node, key int, record *interfaces.Record) {
	pos := node.findPosition(key)

	// Insert key
	node.Keys = append(node.Keys, 0)
	copy(node.Keys[pos+1:], node.Keys[pos:])
	node.Keys[pos] = key

	// Insert record
	node.Records = append(node.Records, nil)
	copy(node.Records[pos+1:], node.Records[pos:])
	node.Records[pos] = record
}

func (t *BTree) splitLeaf(node *Node) {
	// Create new leaf node
	newNode := &Node{
		IsLeaf:   true,
		Keys:     make([]int, 0),
		Records:  make([]*interfaces.Record, 0),
		Children: nil,
		Next:     node.Next,
	}
	node.Next = newNode

	// Find split point
	splitPoint := (len(node.Keys) + 1) / 2

	// Move half of the keys and records to the new node
	newNode.Keys = append(newNode.Keys, node.Keys[splitPoint:]...)
	newNode.Records = append(newNode.Records, node.Records[splitPoint:]...)
	node.Keys = node.Keys[:splitPoint]
	node.Records = node.Records[:splitPoint]

	// Update parent
	if node == t.Root {
		// Create new root
		newRoot := &Node{
			IsLeaf:   false,
			Keys:     []int{newNode.Keys[0]},
			Children: []*Node{node, newNode},
		}
		t.Root = newRoot
		node.Parent = newRoot
		newNode.Parent = newRoot
	} else {
		// Insert into parent
		newNode.Parent = node.Parent
		t.insertIntoParent(node, newNode.Keys[0], newNode)
	}
}

func (t *BTree) insertIntoParent(leftNode *Node, key int, rightNode *Node) {
	parent := leftNode.Parent
	pos := parent.findPosition(key)

	// Insert key
	parent.Keys = append(parent.Keys, 0)
	copy(parent.Keys[pos+1:], parent.Keys[pos:])
	parent.Keys[pos] = key

	// Insert child pointer
	parent.Children = append(parent.Children, nil)
	copy(parent.Children[pos+2:], parent.Children[pos+1:])
	parent.Children[pos+1] = rightNode

	// Check if we need to split the parent
	if len(parent.Keys) > InternalNodeMaxKeys {
		t.splitInternal(parent)
	}
}

func (t *BTree) splitInternal(node *Node) {
	// Create new internal node
	newNode := &Node{
		IsLeaf:   false,
		Keys:     make([]int, 0),
		Children: make([]*Node, 0),
	}

	// Find split point
	splitPoint := len(node.Keys) / 2
	promotedKey := node.Keys[splitPoint]

	// Move keys and children to new node
	newNode.Keys = append(newNode.Keys, node.Keys[splitPoint+1:]...)
	newNode.Children = append(newNode.Children, node.Children[splitPoint+1:]...)
	node.Keys = node.Keys[:splitPoint]
	node.Children = node.Children[:splitPoint+1]

	// Update children's parent pointers
	for _, child := range newNode.Children {
		child.Parent = newNode
	}

	if node == t.Root {
		// Create new root
		newRoot := &Node{
			IsLeaf:   false,
			Keys:     []int{promotedKey},
			Children: []*Node{node, newNode},
		}
		t.Root = newRoot
		node.Parent = newRoot
		newNode.Parent = newRoot
	} else {
		// Insert into parent
		newNode.Parent = node.Parent
		t.insertIntoParent(node, promotedKey, newNode)
	}
}

// Delete removes a record with the given key from the B-tree
func (t *BTree) Delete(key int) {
	if t.Root == nil {
		return
	}

	// Find the leaf node containing the key
	node := t.Root
	for !node.IsLeaf {
		pos := node.findPosition(key)
		if pos >= len(node.Children) {
			return
		}
		node = node.Children[pos]
	}

	// Find the position of the key in the leaf node
	pos := -1
	for i, k := range node.Keys {
		if k == key {
			pos = i
			break
		}
	}

	// If key not found, return
	if pos == -1 {
		return
	}

	// Remove the key and record
	node.Keys = append(node.Keys[:pos], node.Keys[pos+1:]...)
	node.Records = append(node.Records[:pos], node.Records[pos+1:]...)

	// If root is a leaf node, we're done
	if node == t.Root {
		return
	}

	// If node has enough keys, we're done
	if len(node.Keys) >= LeafNodeMaxRecords/2 {
		return
	}

	// Try to borrow from siblings
	if node.Next != nil && len(node.Next.Keys) > LeafNodeMaxRecords/2 {
		// Borrow from right sibling
		node.Keys = append(node.Keys, node.Next.Keys[0])
		node.Records = append(node.Records, node.Next.Records[0])
		node.Next.Keys = node.Next.Keys[1:]
		node.Next.Records = node.Next.Records[1:]
		return
	}

	// If we can't borrow, merge with next sibling if possible
	if node.Next != nil {
		// Merge with right sibling
		node.Keys = append(node.Keys, node.Next.Keys...)
		node.Records = append(node.Records, node.Next.Records...)
		node.Next = node.Next.Next
	}
}

// Scan retrieves all records from the B-tree
func (t *BTree) Scan() []*interfaces.Record {
	if t.Root == nil {
		return nil
	}

	var records []*interfaces.Record
	node := t.Root

	// Find leftmost leaf node
	for !node.IsLeaf {
		node = node.Children[0]
	}

	// Traverse through leaf nodes
	for node != nil {
		records = append(records, node.Records...)
		node = node.Next
	}

	return records
}

// Search finds a record by key
func (t *BTree) Search(key int) *interfaces.Record {
	if t.Root == nil {
		return nil
	}

	node := t.Root

	// Find leaf node
	for !node.IsLeaf {
		pos := node.findPosition(key)
		node = node.Children[pos]
	}

	// Search in leaf node
	pos := node.findPosition(key)
	if pos < len(node.Keys) && node.Keys[pos] == key {
		return node.Records[pos]
	}

	return nil
}
