package main

type MyInputMethod struct {
	filename []string
	children [26]*MyInputMethod
}

func Constructor() *MyInputMethod {
	return &MyInputMethod{}
}

func (this *MyInputMethod) Insert(word string) {
	cur := this
	for _, c := range word {
		n := c - 'a'

		if n < 0 || n > 25{ //对于不符合规范的文件忽略
			continue
		}
		if cur.children[n] == nil {
			cur.children[n] = &MyInputMethod{}
		}

		cur = cur.children[n]
		cur.filename = append(cur.filename, "./BuildData/WordLab/" + word +".dat")
	}
}

