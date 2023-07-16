package trie

type Trie struct {
	root *node
}

type node struct {
	// 限定26个小写字母a-z
	nexts [26]*node
	// 从根节点到当前节点的单词数量
	passCnt int
	// 是否存在以当前节点结尾的单词
	end bool
}

func NewTrie() *Trie {
	return &Trie{root: &node{}}
}

func (t *Trie) Exists(word string) bool {
	node := t.search(word)
	return node != nil && node.end
}

func (t *Trie) StartWith(prefix string) bool {
	return t.search(prefix) != nil
}

// PassCnt 统计以prefix为前缀的单词数量
func (t *Trie) PassCnt(prefix string) int {
	node := t.search(prefix)
	if node == nil {
		return 0
	}

	return node.passCnt
}

func (t *Trie) Insert(word string) {
	if t.Exists(word) {
		return
	}

	move := t.root

	for _, ch := range word {
		idx := ch - 'a'
		if move.nexts[idx] == nil {
			move.nexts[idx] = &node{}
		}
		move.nexts[idx].passCnt++
		move = move.nexts[idx]
	}

	move.end = true
}

// Erase 删除word
func (t *Trie) Erase(word string) bool {
	if !t.Exists(word) {
		return false
	}

	move := t.root
	for _, ch := range word {
		idx := ch - 'a'
		move.nexts[idx].passCnt--
		// passCnt 为0则直接舍弃子节点，返回
		if move.nexts[idx].passCnt == 0 {
			move.nexts[idx] = nil
			return true
		}

		move = move.nexts[idx]
	}

	// 遍历到单词末尾，将end置为false
	move.end = false
	return true
}

func (t *Trie) search(word string) *node {
	move := t.root

	for _, ch := range word {
		// 如果有一个字符不存在对应的节点，则说明没有插入过
		if move.nexts[ch-'a'] == nil {
			return nil
		}

		move = move.nexts[ch-'a']
	}

	return move
}
