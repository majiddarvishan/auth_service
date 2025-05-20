package trie

type charNode struct {
	destination string

	parent *charNode
	depth  uint32
	index  int32

	hasValueFlag bool
	maxLength    uint32

	numberOfChildren uint32
	children         [256]*charNode
}

func NewCharNode(parent *charNode, depth uint32, index int32) *charNode {
	n := &charNode{
		parent:           parent,
		depth:            depth,
		index:            index,
		hasValueFlag:     false,
		maxLength:        0,
		numberOfChildren: 0,
	}
	return n
}

func (n *charNode) clean() {
	for i := 0; i < len(n.children); i++ {
		if n.children[i] != nil {
			n.children[i].clean()
			n.children[i] = nil
		}
	}

	n.numberOfChildren = 0
}

type Trie struct {
	tree_root *charNode
}

func NewTrie() *Trie {
	return &Trie{
		tree_root: NewCharNode(nil, 0, -1),
	}
}

func (t *Trie) AcceptAll() {
    t.tree_root.hasValueFlag = true
    t.tree_root.maxLength = 100
}

func (t *Trie) Add(prefix string, destination string, maxLength uint32 /*= 0xFFFFFFFF*/) bool {
	n := t.createSubTree(prefix, 0, uint32(len(prefix)), t.tree_root)
	if n != nil {
		n.destination = destination
		n.maxLength = maxLength
		n.hasValueFlag = true

		return true
	}

	// LOG_ERROR("Add Route Failed! A Route For {} Already Exists!", prefix);
	return false
}

func (t *Trie) Update(prefix string, destination string, maxLength uint32 /*= 0xFFFFFFFF*/) bool {
	n := t.findMatching(prefix, 0, uint32(len(prefix)), t.tree_root, nil)
	if (n != nil) && (n.depth == uint32(len(prefix))) {
		n.destination = destination
		n.maxLength = maxLength
		n.hasValueFlag = true

		return true
	}

	// LOG_ERROR("Update Route Failed! No Route For {} Exists!", prefix);
	return false
}

func (t *Trie) Remove(prefix string) bool {
	n := t.findMatching(prefix, 0, uint32(len(prefix)), t.tree_root, nil)
	if (n != nil) && (n.depth == uint32(len(prefix))) {
		n.hasValueFlag = false
		t.cleanUp(n)
		return true
	}

	// LOG_ERROR("Delete Route Failed! No Route For {} Exists!", prefix);
	return false
}

func (t *Trie) ClearAll() {
	t.tree_root.clean()
}

func (t *Trie) Find(prefix string) (bool, string) {
	n := t.findMatching(prefix, 0, uint32(len(prefix)), t.tree_root, nil)
	if n != nil {
		if n.maxLength >= uint32(len(prefix)) {
            return true, n.destination
		}
	}

	return false, ""
}

func (t *Trie) cleanUp(current *charNode) {
	if (!current.hasValueFlag) && (current.numberOfChildren == 0) {
		parent := current.parent
		if parent != nil {
			parent.numberOfChildren--
			parent.children[current.index] = nil
			t.cleanUp(parent)
		}
	}
}

func (t *Trie) createSubTree(prefix string, first uint32, last uint32, current *charNode) *charNode {
	if first < last {
		index := prefix[first]

		if current.children[index] == nil {
			current.children[index] = NewCharNode(current, current.depth+1, int32(index))
			current.numberOfChildren++
		}

		return t.createSubTree(prefix, first+1, last, current.children[index])
	}

	if current.hasValueFlag {
		return nil
	}

	return current
}

func (t *Trie) findMatching(prefix string, first uint32, last uint32, current *charNode, lastValidValue *charNode) *charNode {
	if current.hasValueFlag {
		lastValidValue = current
	}

	if first < last {
		index := prefix[first]

		if current.children[index] != nil {
			return t.findMatching(prefix, first+1, last, current.children[index], lastValidValue)
		}
	}

	return lastValidValue
}
