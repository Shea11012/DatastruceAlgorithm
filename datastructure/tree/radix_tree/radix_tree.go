package radix_tree

import "strings"

/**
radix tree 与 trie区别：
如果某个父节点有且仅有一个子节点，并且不存在单词以这个父节点为结尾，
则此时radix tree会将这个父节点与子节点进行合并，并把父、子节点的相对路径组装在一起
生成一段复合路径，以此作为新节点。
*/

type Radix struct {
	root *node
}

type node struct {
	// 当前节点的相对路径
	path string
	// 完整路径
	fullPath string
	// 每个indice字符对应一个子节点的path首字母
	indices string
	// 后继节点
	children []*node
	// 是否有路径以当前节点为终点
	end bool
	// 记录有多少路径途径当前节点
	passCnt int
}

func (n *node) insert(word string) {
	fullword := word

	// 如果当前节点为root，且之前没有注册过子节点，则直接插入
	if n.path == "" && len(n.children) == 0 {
		n.insertWord(word, word)
		return
	}
walk:
	for {
		// 获取公共前缀长度
		i := commonPrefixLen(word, n.path)
		// > 0 则表示一定经过当前节点，需要累加 passCnt
		if i > 0 {
			n.passCnt++
		}

		// 公共前缀小于当前节点的相对路径，需要对节点进行分解
		if i < len(n.path) {
			child := &node{
				path:     n.path[i:],
				fullPath: n.fullPath,
				children: n.children,
				indices:  n.indices,
				end:      n.end,
				passCnt:  n.passCnt - 1,
			}

			// 续接上子节点
			n.indices = string(n.path[i])
			n.children = []*node{child}
			// 调整原节点fullPath
			n.fullPath = n.fullPath[:len(n.fullPath)-len(n.path)-i]
			// 调整原节点path
			n.path = n.path[:i]
			// 原节点是新拆分出来的，当前不会有单词以该节点结尾
			n.end = false
		}

		// 公共前缀小于word长度
		if i < len(word) {
			// 去除公共前缀部分
			word = word[i:]
			// 获取word剩余部分首字母
			c := word[0]

			for i := 0; i < len(n.indices); i++ {
				// 与后继节点还有公共前缀，则将n指向子节点，重复执行该流程
				if n.indices[i] == c {
					n = n.children[i]
					continue walk
				}
			}

			// word剩余部分与后继节点没有公共前缀
			// 构造新节点插入
			n.indices += string(c)
			child := &node{}
			child.insertWord(word, fullword)
			n.children = append(n.children, child)
			return
		}

		// 公共前缀刚好是path，将end置为true
		n.end = true
		return
	}

}

func (n *node) insertWord(path string, fullPath string) {
	n.path, n.fullPath = path, fullPath
	n.passCnt = 1
	n.end = true
}

func (n *node) search(word string) *node {
walk:
	for {
		prefix := n.path
		if len(word) > len(prefix) {
			// 没匹配上，直接返回nil
			if word[:len(prefix)] != prefix {
				return nil
			}

			// 截取公共前缀后的剩余部分
			word = word[len(prefix):]
			c := word[0]
			for i := 0; i < len(n.indices); i++ {
				if c == n.indices[i] {
					n = n.children[i]
					continue walk
				}
			}

			// 到这里，表示word还有剩余部分但与后继节点不匹配
			return nil
		}

		// 和当前节点匹配上
		if word == prefix {
			return n
		}

		// 表示 len(word) <= len(n.path) && word != n.path
		return n
	}
}

// 获取a,b两个字符串的公共前缀长度
func commonPrefixLen(a, b string) int {
	move := 0
	for move < len(a) && move < len(b) && a[move] == b[move] {
		move++
	}

	return move
}

func NewRadix() *Radix {
	return &Radix{root: &node{}}
}

func (r *Radix) Insert(word string) {
	if r.Exists(word) {
		return
	}

	r.root.insert(word)
}

func (r *Radix) Search(word string) bool {
	node := r.root.search(word)
	return node != nil && node.fullPath == word && node.end
}

func (r *Radix) StartWith(prefix string) bool {
	node := r.root.search(prefix)
	return node != nil && strings.HasPrefix(node.fullPath, prefix)
}

func (r *Radix) PassCnt(prefix string) int {
	node := r.root.search(prefix)
	if node == nil || !strings.HasPrefix(node.fullPath, prefix) {
		return 0
	}

	return node.passCnt
}

func (r *Radix) Erase(word string) bool {
	if !r.Exists(word) {
		return false
	}

	// 与root匹配
	if r.root.fullPath == word {
		// 没有子节点
		if len(r.root.indices) == 0 {
			r.root.path = ""
			r.root.fullPath = ""
			r.root.end = false
			r.root.passCnt = 0
			return true
		}

		// 只有一个子节点
		if len(r.root.indices) == 1 {
			r.root.children[0].path = r.root.path + r.root.children[0].path
			r.root = r.root.children[0]
			return true
		}

		// 多个子节点
		for i := 0; i < len(r.root.indices); i++ {
			r.root.children[i].path = r.root.path + r.root.children[0].path
		}

		newRoot := &node{
			indices:  r.root.indices,
			children: r.root.children,
			passCnt:  r.root.passCnt,
		}
		r.root = newRoot
		return true
	}

	move := r.root
walk:
	for {
		move.passCnt--
		prefix := move.path
		word = word[len(prefix):]
		c := word[0]
		for i := 0; i < len(move.indices); i++ {
			if move.indices[i] != c {
				continue
			}

			// 命中，但是任有后继节点
			if move.children[i].path == word && move.children[i].passCnt > 1 {
				move.children[i].end = false
				move.children[i].passCnt--
				return true
			}

			// 找到对应的child，且没有后继节点
			if move.children[i].passCnt > 1 {
				move = move.children[i]
				continue walk
			}

			move.children = append(move.children[:i], move.children[i+1:]...)
			move.indices = move.indices[:i] + move.indices[i+1:]

			// 最后一个节点，只有一个子节点且自身end为false，则需要进行合并
			if !move.end && len(move.indices) == 1 {
				move.path = move.children[0].path
				move.fullPath = move.children[0].path
				move.end = move.children[0].end
				move.indices = move.children[0].indices
				move.children = move.children[0].children
			}

			return true
		}
	}
}

func (r *Radix) Exists(word string) bool {
	node := r.root.search(word)
	return node != nil && node.end == true
}
