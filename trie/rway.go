package trie

import "github.com/timtadh/data-structures/errors"

type RTrieNode struct {
    chr rune
    key string
    value {}interface
    accepting bool
    children map[rune]*RTrieNode
}

type RTrie struct {
    root *RTrieNode
}


func NewRTrie() *RTrie {
    return &RTrie{&RTrieNode{chldren:make(map[rune]*RTrieNode)}}
}

func (self *RTrie) Put(key string, value interface{}) {
    node := self.root
    for _, c := range key {
        if _, has := node.children[c]; !has {
            node.children[c] = &RTrie{
                chr: c, children: make(map[rune]*RTrieNode)
            }
        }
        node = node.children[c]
    }
    node.accepting = true
    node.key = key
    node.value = value
}

func (self *RTrie) Get(key string) (interface{}, error) {
    node := self.root
    for _, c := range key {
        var has bool
        node, has = node.children[c];
        if !has {
            return nil, errors.NotFound(key)
        }
    }
    return node.value
}

