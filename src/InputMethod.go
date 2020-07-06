/*
功能：本程序是拼音输⼊法核⼼库的核心文件，能够根据用户的拼音输入查找最有可能的汉字排列
作者：xiaoqing_tiffany@foxmail.com
用户通过修改 skillfolder 修改词典的来源库
*/
package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

func NewMyInputMethod(dicts []string) *MyInputMethod {

	/*
	功能：根据传⼊的词典⽂件创建⼀个新的输⼊法实例，如果词典⽂件格式有误，忽略格式有误的⽂件
	将文件名字构建成一个前缀树
	如何建立呢？
	按照文件名字的字典顺序以及长度
	*/
	root := Constructor()
	for _, v := range dicts{
		/*
		left := strings.LastIndex(v, "/")
		right := strings.LastIndex(v, ".")
		if left == -1 || right == -1 {
			continue
		}
		word := v[left+1:right]
		*/
		right := strings.LastIndex(v, ".")
		if right == -1 {
			continue
		}
		word := v[0:right]
		root.Insert(word)
	}
	return root
}

//查找文件下是否有对应文件，有则返回true
func IfExit(spell string, folder string) bool {
	files, _ := ioutil.ReadDir(folder)
	for _,file := range files {
		if file.IsDir() {
			continue
		} else {
			if file.Name() == spell+".dat"{
				return true}
		}
	}
	return false
}

//对得到的文件中所有词条做词频统计。将词频和词加入到最小堆当中。依次输出最小堆当中的值
func smallHeap(dicname []string) (words []string) {
	h := &NodeHeap{}  //小根堆头节点
	topK := 10        //取频率最高的前10个汉字
	size := 0
	res := make([]string, 0)   //保存最后结果
	count := 0        //记录加入顺序，代表着字典顺序，也代表着在文件中的顺序

	for _, v := range dicname{
		file ,err := os.Open(v)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		//得到文件的长度。因为要对拼音的长度进行排序，文件越长，拼音也越长

		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanLines)   //一行一行读取

		for scanner.Scan() {  //是否有下一行
			text := scanner.Text()
			textRune := []rune(text)     //转换为rune类型，避免中文的乱码
			c := string(textRune[0:1])   //字符
			f := string(textRune[2:])   //频率
			i, err := strconv.Atoi(f)    //频率string to int
			if err == nil{
				if size < topK {
					heap.Push(h, &Node{
						frequency: i,
						Index: count,
						content: c,
					})
					size++
				} else {
					if i > (*h)[0].frequency {   //遍历到某一行时，如果该行频率大于小根堆的最小值，则替换这个节点
						heap.Pop(h)
						heap.Push(h, &Node{
							frequency: i,
							Index: count,
							content: c,
						})
					}
				}
				count = count + 1
			}


		}
	}

	res_frequency := make([]int, 0)   //保存topk10个节点的频率
	res_index := make([]int, 0)   	  //保存top10个节点的在文件中的顺序
	topK = min(10, h.Len())
	for i := 0; i < topK; i++ {
		popNode :=  heap.Pop(h).(*Node)

		res = append(res, popNode.content)
		res_frequency = append(res_frequency, popNode.frequency)
		res_index = append(res_index, popNode.Index)

		for j := i; j > 1 && res_frequency[j] == res_frequency[j-1] && res_index[j] > res_index[j-1];j-- {
				//在频率相同的情况下，将顺序大放前面
				res[j], res[j-1] = res[j-1], res[j]
				res_index[j], res_index[j-1] = res_index[j-1],res_index[j]
			}


	}

	//因为小根堆保存的频率是从小到大排列的，现在需要将结果反转
	for i := len(res)/2-1; i >= 0; i-- {
		opp := len(res)-1-i
		res[i], res[opp] = res[opp], res[i]
	}

	return res
}

func (mim *MyInputMethod)FindWords(spell string) (words []string) {
	/*
		查看spell有没有直接对应的文件：如果有，直接返回该文件
									否则，找到前缀树对应节点包含的所有路径。
		接着，将文件内容加入小根堆中
	*/
	spell = strings.ToLower(spell)
	cur := mim        //前缀树的根节点

	res := make([]string, 0)   //保存最后结果
	if len(spell) == 0{
		res = append(res, "输入为空行")
		return res
	}

	if IfExit(spell, "./BuildData/WordLab/"){
		dicname := make([]string,0)
		dicname = append(dicname, "./BuildData/WordLab/" + spell +".dat")
		return smallHeap(dicname)
	} else{
		for _, c := range spell {
			n := c - 'a'
			if n < 0 || n > 25{    //如果不在字母范围之内
				res = append(res, "输入有误，请重新输入")
				return res
			}
			if cur.children[n] == nil {  //如果找不到对应的拼音文件
				return res
			}
			cur = cur.children[n]
		}
		dicname := cur.filename   //找出所有备选的文件夹名字
		return smallHeap(dicname)
	}
}

type InputMethod interface {
	FindWords(string) []string
}

func loop(im InputMethod) {
	stdin := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("请输入拼音")
		spell, err := stdin.ReadString('\n')
		if err != nil {
			break
		}
		spell = strings.TrimRight(spell, "\n")
		words := im.FindWords(spell)
		fmt.Println(words)
	}
}

func main() {
	processFile := make([]string, 0)
		skillfolder := `E:\Go-work\code\BuildData\WordLab\`
		// 获取所有文件
		files, _ := ioutil.ReadDir(skillfolder)
		for _,file := range files {
			if file.IsDir() {
				continue
			} else {
				processFile = append(processFile, file.Name())
			}
	}
/*
	sort.Slice(processFile, func(i int, j int) bool { return processFile[i] < processFile[j] })
		fmt.Println(processFile)
*/


	im := NewMyInputMethod(processFile)
	loop(im)
/*可以选择运行.\BuildData\WordLab\文件夹下的部分文件
		im := NewMyInputMethod([]string{
		"e.dat",
		"ei.dat",
		"en.dat",
		"er.dat",
	})
	*/

}
