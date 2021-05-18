package frame

import "strings"

// trie tree

type node struct {
	pattern string  //
	part string 	// chidren value
	children []*node// tree children
	flag  bool // flag=ture if part include "*" or ":"
}

func (n *node) matchFirstChild(part string) *node {
	for _,child := range n.children{
		if child.part == part || child.flag {
			return child
		}
	}
	return nil
}

func (n *node) matchAllChild(part string)[]*node {
	var nodes []*node
	for _ , child := range n.children{
		if child.part == part || child.flag{
			nodes =append(nodes, child)
		}
	}
	return nodes
}

func (n *node) insert(pattern string,parts []string,height int){
	if len(parts) == height{
		n.pattern = pattern
		return
	}
	part := parts[height]
	child := n.matchFirstChild(part)
	if child == nil{
		child = &node{
			part:     part,
			flag:     part[0] == ':' || part[0] == '*',
		}
		n.children = append(n.children,child)
	}
	child.insert(pattern,parts,height+1)
}

func (n *node) search(parts []string , height int) *node {
	if len(parts) == height || strings.HasPrefix(n.part,"*"){
		if n.pattern == ""{
			return nil
		}
		return n
	}
	part := parts[height]
	children := n.matchAllChild(part)
	for _,children := range children{
		result := children.search(parts,height+1)
		if result != nil{
			return result
		}
	}
	return nil
}
