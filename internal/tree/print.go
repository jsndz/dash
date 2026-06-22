package tree

import (
	"fmt"
)

// PrintTreeColor prints the tree structure with clean colors and branch markers.
func (tree *Bptree) PrintTreeColor() {
	// ANSI escape colors
	const (
		Reset  = "\033[0m"
		Red    = "\033[31m" // Root
		Green  = "\033[32m" // Internal Nodes
		Yellow = "\033[33m" // Leaves
		Cyan   = "\033[36m" // Branch text labels
		Gray   = "\033[90m" // Decorative framework
	)

	fmt.Println(Gray + "\n================= COLORIZED B+ TREE =================\n" + Reset)
	if tree.Root == nil {
		fmt.Println("Empty Tree")
		return
	}

	var printNode func(node *Node, level int, prefix string, isLast bool)
	printNode = func(node *Node, level int, prefix string, isLast bool) {
		// Determine node type styling
		nodeType := "Internal"
		color := Green
		if node == tree.Root {
			nodeType = "Root"
			color = Red
		} else if node.IsLeaf {
			nodeType = "Leaf"
			color = Yellow
		}

		// Prepare the output line string
		nodeString := fmt.Sprintf("%s[%s]%s keys: %v", color, nodeType, Reset, node.Keys)

		// Render current line with proper branch structure symbols
		marker := "├── "
		if isLast {
			marker = "└── "
		}
		if level == 0 {
			marker = ""
		}

		fmt.Printf("%s%s%s\n", Gray+prefix, marker+Reset, nodeString)

		// Calculate trailing prefix spacing for children branches
		nextPrefix := prefix
		if level > 0 {
			if isLast {
				nextPrefix += "    "
			} else {
				nextPrefix += "│   "
			}
		} else {
			nextPrefix = ""
		}

		// Recurse down children branches if this is an internal node
		if !node.IsLeaf {
			for i, child := range node.Children {
				childIsLast := i == len(node.Children)-1
				label := fmt.Sprintf("%s(C%d)%s ", Cyan, i, Reset)
				printNode(child, level+1, nextPrefix+label, childIsLast)
			}
		}
	}

	printNode(tree.Root, 0, "", true)
	fmt.Println(Gray + "\n=====================================================" + Reset)
}
