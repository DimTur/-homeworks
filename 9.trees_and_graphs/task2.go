// Реализуйте 2-3-дерево: вставка, поиск, удаление

package main

// Node interfaces
type Finder interface {
	Find(k int) bool
}

type Swapper interface {
	Swap(x, y *int)
}

type Sorter interface {
	Sort()
	Sort2(x, y *int)
	Sort3(x, y, z *int)
}

type Inserter interface {
	InsertToNode(k int)
}

type Remover interface {
	RemoveFromNode(k int)
}

type Transformer interface {
	BecomeNode2(k int, first, second *Node)
}

type LeafChecker interface {
	IsLeaf() bool
}

// Node struct
type Node struct {
	size   int
	key    [3]int
	first  *Node
	second *Node
	third  *Node
	fourth *Node
	parent *Node
}

// Implementing Node interface

// Find method return true if key `k` on the vertex, else `false`
func (n *Node) Find(k int) bool {
	for i := 0; i < n.size; i++ {
		if n.key[i] == k {
			return true
		}
	}
	return false
}

// Swap method
func (n *Node) Swap(x, y *int) {
	*x, *y = *y, *x
}

// Sort2 method
func (n *Node) Sort2(x, y *int) {
	if *x > *y {
		n.Swap(x, y)
	}
}

// Sort3 method
func (n *Node) Sort3(x, y, z *int) {
	if *x > *y {
		n.Swap(x, y)
	}
	if *x > *z {
		n.Swap(x, z)
	}
	if *y > *z {
		n.Swap(y, z)
	}
}

// Sort method
func (n *Node) Sort() {
	if n.size == 1 {
		return
	}
	if n.size == 2 {
		n.Sort2(&n.key[0], &n.key[1])
	}
	if n.size == 3 {
		n.Sort3(&n.key[0], &n.key[1], &n.key[2])
	}
}

// InsertToNode method insert key `k` to the vertex (not in tree)
func (n *Node) InsertToNode(k int) {
	n.key[n.size] = k
	n.size++
	n.Sort()
}

// RemoveFromNode method remove key `k` from the vertex (not from tree)
func (n *Node) RemoveFromNode(k int) {
	if n.size >= 1 && n.key[0] == k {
		n.key[0] = n.key[1]
		n.key[1] = n.key[2]
		n.size--
	} else if n.size == 2 && n.key[1] == k {
		n.key[1] = n.key[2]
		n.size--
	}
}

// BecomeNode2 method transform to second vertex
func (n *Node) BecomeNode2(k int, first, second *Node) {
	n.key[0] = k
	n.first = first
	n.second = second
	n.third = nil
	n.fourth = nil
	n.parent = nil
	n.size = 1
}

// IsLeaf method whether the node is a leaf; validation is used on insertion and deletion
func (n *Node) IsLeaf() bool {
	return n.first == nil && n.second == nil && n.third == nil
}

// Split func devide a node if it contains 3 keys
func Split(item *Node) *Node {
	if item.size < 3 {
		return item
	}

	x := &Node{
		size:   1,
		key:    [3]int{item.key[0]},
		first:  item.first,
		second: item.second,
		parent: item.parent,
	}
	y := &Node{
		size:   1,
		key:    [3]int{item.key[2]},
		first:  item.third,
		second: item.fourth,
		parent: item.parent,
	}

	switch {
	case x.first != nil:
		x.first.parent = x
	case x.second != nil:
		x.second.parent = x
	}
	switch {
	case y.first != nil:
		y.first.parent = y
	case y.second != nil:
		y.second.parent = y
	}

	if item.parent != nil {
		item.parent.InsertToNode(item.key[1])

		switch item {
		case item.parent.first:
			item.parent.first = nil
		case item.parent.second:
			item.parent.second = nil
		case item.parent.third:
			item.parent.third = nil
		}

		switch {
		case item.parent.first == nil:
			item.parent.fourth = item.parent.third
			item.parent.third = item.parent.second
			item.parent.second = y
			item.parent.first = x
		case item.parent.second == nil:
			item.parent.fourth = item.parent.third
			item.parent.third = y
			item.parent.second = x
		default:
			item.parent.fourth = y
			item.parent.third = x
		}

		parent := item.parent
		item = nil
		return Split(parent)
	} else {
		x.parent = item
		y.parent = item
		item.BecomeNode2(item.key[1], x, y)
		return item
	}
}

// Insert is func, that insert key `k` to the tree with root `p`
// Always reterns root of the tree, becouse it was change
func Insert(p *Node, k int) *Node {
	if p == nil {
		return &Node{size: 1, key: [3]int{k}}
	}

	switch {
	case p.IsLeaf():
		p.InsertToNode(k)
	case k <= p.key[0]:
		p.first = Insert(p.first, k)
	case p.size == 1 || (p.size == 2 && k <= p.key[1]):
		p.second = Insert(p.second, k)
	default:
		p.third = Insert(p.third, k)
	}

	return Split(p)
}

// Search func finding key `k` in a 2-3 tree with root `p`
func Search(p *Node, k int) *Node {
	if p == nil {
		return nil
	}

	switch {
	case p.Find(k):
		return p
	case k < p.key[0]:
		return Search(p.first, k)
	case (p.size == 2 && k < p.key[1]) || p.size == 1:
		return Search(p.second, k)
	case p.size == 2:
		return Search(p.third, k)
	default:
		return nil
	}
}

// SearchMin func finding minimal key `k` in a 2-3 tree with root `p`
func SearchMin(p *Node) *Node {
	if p == nil {
		return nil
	}
	if p.first == nil {
		return p
	}
	return SearchMin(p.first)
}

// The Fix function corrects the structure of 2-3 trees after deleting a key,
// checking and performing the necessary redistribution or merging of nodes
// to preserve the properties of the tree.
func Fix(leaf *Node) *Node {
	if leaf.size == 0 && leaf.parent == nil {
		return nil
	}
	if leaf.size != 0 {
		if leaf.parent != nil {
			return Fix(leaf.parent)
		}
		return leaf
	}

	parent := leaf.parent
	switch {
	case parent.first.size == 2 || parent.second.size == 2 || parent.size == 2:
		leaf = Redistribute(leaf)
	case parent.size == 2 && parent.third.size == 2:
		leaf = Redistribute(leaf)
	default:
		leaf = Merge(leaf)
	}

	return Fix(leaf)
}

// Redistribute func:
// We remove the key from the vertex and the vertex becomes empty.
// If at least one of the brothers has 2 keys,
// then we do a simple correct distribution and the work is finished.
// By correct distribution I mean that if you rotate keys between a parent and its children,
// you will also need to move the grandchildren of the parent.
// You can redistribute keys from any brother,
// but it is most convenient from the neighbor, which has 2 keys,
// while we cyclically shift all the keys
func Redistribute(leaf *Node) *Node {
	parent := leaf.parent
	first := parent.first
	second := parent.second
	third := parent.third

	switch {
	case parent.size == 2 && first.size < 2 && second.size < 2 && third.size < 2:
		switch {
		case first == leaf:
			parent.first = parent.second
			parent.second = parent.third
			parent.third = nil
			parent.first.InsertToNode(parent.key[0])
			parent.first.third = parent.first.second
			parent.first.second = parent.first.first

			if leaf.first != nil {
				parent.first.first = leaf.first
			} else if leaf.second != nil {
				parent.first.first = leaf.second
			}

			if parent.first.first != nil {
				parent.first.first.parent = parent.first
			}

			parent.RemoveFromNode(parent.key[0])
			return parent
		case second == leaf:
			first.InsertToNode(parent.key[0])
			parent.RemoveFromNode(parent.key[0])
			if leaf.first != nil {
				first.third = leaf.first
			} else if leaf.second != nil {
				first.third = leaf.second
			}

			if first.third != nil {
				first.third.parent = first
			}

			parent.second = parent.third
			parent.third = nil

			return parent
		case third == leaf:
			second.InsertToNode(parent.key[1])
			parent.RemoveFromNode(parent.key[1])
			if leaf.first != nil {
				second.third = leaf.first
			} else if leaf.second != nil {
				second.third = leaf.second
			}

			if second.third != nil {
				second.third.parent = second
			}

			parent.third = nil

			return parent
		}
	case parent.size == 2 && (first.size == 2 || second.size == 2 || third.size == 2):
		switch {
		case third == leaf:
			if leaf.first != nil {
				leaf.second = leaf.first
				leaf.first = nil
			}

			leaf.InsertToNode(parent.key[1])
			switch {
			case second.size == 2:
				parent.key[1] = second.key[1]
				second.RemoveFromNode(second.key[1])
				leaf.first = second.third
				second.third = nil
				if leaf.first != nil {
					leaf.first.parent = leaf
				}
			case first.size == 2:
				parent.key[1] = second.key[0]
				leaf.first = second.second
				second.second = second.first
				if leaf.first != nil {
					leaf.first.parent = leaf
				}

				second.key[0] = parent.key[0]
				parent.key[0] = first.key[1]
				first.RemoveFromNode(first.key[1])
				second.first = first.third
				if second.first != nil {
					second.first.parent = second
				}
				first.third = nil
			}
		case second == leaf:
			switch {
			case third.size == 2:
				if leaf.first == nil {
					leaf.first = leaf.second
					leaf.second = nil
				}
				second.InsertToNode(parent.key[1])
				parent.key[1] = third.key[0]
				third.RemoveFromNode(third.key[0])
				second.second = third.first
				if second.second != nil {
					second.second.parent = second
				}
				third.first = third.second
				third.second = third.third
				third.third = nil
			case first.size == 2:
				if leaf.second == nil {
					leaf.second = leaf.first
					leaf.first = nil
				}
				second.InsertToNode(parent.key[0])
				parent.key[0] = first.key[1]
				first.RemoveFromNode(first.key[1])
				second.first = first.third
				if second.first != nil {
					second.first.parent = second
				}
				first.third = nil
			}
		case first == leaf:
			if leaf.first == nil {
				leaf.first = leaf.second
				leaf.second = nil
			}
			first.InsertToNode(parent.key[0])
			switch {
			case second.size == 2:
				parent.key[0] = second.key[0]
				second.RemoveFromNode(second.key[0])
				first.second = second.first
				if first.second != nil {
					first.second.parent = first
				}
				second.first = second.second
				second.second = second.third
				second.third = nil
			case third.size == 2:
				parent.key[0] = second.key[0]
				second.key[0] = parent.key[1]
				parent.key[1] = third.key[0]
				third.RemoveFromNode(third.key[0])
				first.second = second.first
				if first.second != nil {
					first.second.parent = first
				}
				second.first = second.second
				second.second = third.first
				if second.second != nil {
					second.second.parent = second
				}
				third.first = third.second
				third.second = third.third
				third.third = nil
			}
		}
	case parent.size == 1:
		leaf.InsertToNode(parent.key[0])

		switch {
		case first == leaf && second.size == 2:
			parent.key[0] = second.key[0]
			second.RemoveFromNode(second.key[0])

			if leaf.first == nil {
				leaf.first = leaf.second
			}

			leaf.second = second.first
			second.first = second.second
			second.second = second.third
			second.third = nil
			if leaf.second != nil {
				leaf.second.parent = leaf
			}
		case second == leaf && first.size == 2:
			parent.key[0] = first.key[1]
			first.RemoveFromNode(first.key[1])

			if leaf.second == nil {
				leaf.second = leaf.first
			}

			leaf.first = first.third
			first.third = nil
			if leaf.first != nil {
				leaf.first.parent = leaf
			}
		}
	}

	return parent
}

// The Merge function merges a given leaf node with its neighbor node in a 2-3 tree,
// moving keys and children from the leaf to the neighbor node,
// and updates references to parent nodes to maintain a balanced tree structure.
func Merge(leaf *Node) *Node {
	parent := leaf.parent

	switch {
	case parent.first == leaf:
		parent.second.InsertToNode(parent.key[0])
		parent.second.third = parent.second.second
		parent.second.second = parent.second.first

		if leaf.first != nil {
			parent.second.first = leaf.first
		} else if leaf.second != nil {
			parent.second.first = leaf.second
		}

		if parent.second.first != nil {
			parent.second.first.parent = parent.second
		}

		parent.RemoveFromNode(parent.key[0])
		parent.first = nil

	case parent.second == leaf:
		parent.first.InsertToNode(parent.key[0])

		if leaf.first != nil {
			parent.first.third = leaf.first
		} else if leaf.second != nil {
			parent.first.third = leaf.second
		}

		if parent.first.third != nil {
			parent.first.third.parent = parent.first
		}

		parent.RemoveFromNode(parent.key[0])
		parent.second = nil
	}

	if parent.parent == nil {
		var tmp *Node
		if parent.first != nil {
			tmp = parent.first
		} else {
			tmp = parent.second
		}
		tmp.parent = nil
		return tmp
	}

	return parent
}

// The Remove function removes key k from a 2-3 tree rooted at p.
// It first finds the node containing the key, then, if necessary,
// looks for an equivalent replacement key, swaps them, and removes the key from the leaf.
// After this, the Fix function is called to restore the properties of the tree.
func Remove(p *Node, k int) *Node {
	item := Search(p, k)

	if item == nil {
		return p
	}

	var min *Node
	if item.key[0] == k {
		min = SearchMin(item.second)
	} else {
		min = SearchMin(item.third)
	}

	if min != nil {
		var z *int
		if k == item.key[0] {
			z = &item.key[0]
		} else {
			z = &item.key[1]
		}
		item.Swap(z, &min.key[0])
		item = min
	}

	item.RemoveFromNode(k)
	return Fix(item)
}
