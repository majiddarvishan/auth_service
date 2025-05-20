package trie

import (
	"fmt"
	"math"
	"strings"
)

func trimSuffix(s, suffix string) (string, bool) {
	isTrimed := false
	if strings.HasSuffix(s, suffix) {
		s = s[:len(s)-len(suffix)]
		isTrimed = true
	}
	return s, isTrimed
}

type TrieManager struct {
	trie *Trie
}

func NewTrieManager() *TrieManager {
	return &TrieManager{
		trie: NewTrie(),
	}
}

func (manager *TrieManager) Add(userName, prefix string) error {
	p, isTrimed := trimSuffix(prefix, "*")
	if isTrimed {
		if !manager.trie.Add(p, userName, math.MaxInt32) {
			return fmt.Errorf("add prefix is failed")
		}
	} else {
		if !manager.trie.Add(p, userName, uint32(len(p))) {
			return fmt.Errorf("add prefix is failed")
		}
	}

	return nil
}

func (manager *TrieManager) Remove(userName, prefix string) error {
    p, _ := trimSuffix(prefix, "*")

	if !manager.trie.Remove(p) {
		return fmt.Errorf("remove route is failed")
	}

	return nil
}

func (manager *TrieManager) Find(prefix string) (bool, string) {
	return manager.trie.Find(prefix)
}


var TrieManagerInstance *TrieManager