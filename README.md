# 拼音输入法核心库 
实现了⼀个基本的输⼊法核心库。根据输⼊的词典创建了前缀树索引，并根据传⼊的拼写查询条件，找出最匹配的候选汉字列表。例如输⼊"de"，则返回 ["的", "得", "地", "德"]。 
# 实现步骤
## 1.创建本地词条
 首先创建后缀为 .dat的文件本地⽂件。每个⽂件中包含多⾏内容，每一行包含一个汉字和这个汉字的出现频次得分，频次分最⼩为1，最大为 10，均为整数，得分越大表示越是⾼频词。相同拼音的汉字保存在同一文件下。
 例如 de.dat ⽂件中的内容可能是这样：<br/>
汉字 | 频率 
 -: | :-: 
  的 |10
  得 |10
  德 |4
  地 |10 
### 创建本地词条步骤
将来自[汉字单字字频总表](http://lingua.mtsu.edu/chinese-computing/statistics/char/list.php?Which=TO)的原始数据做处理。
原始文件包括以下词条：<br/>
序号 | 汉字 | 频率 | 累计频率 | 拼音以及声调  | 英文翻译
:-: | :-: | :-: | :-: | :-: | :-:
1 | 的 | 8302698 | 3.2074998098725 | de/di2/di4 | (possessive particle)/of, really and truly, aim/clear
2 | 一 | 3728398 | 4.6478552071336 | yi1 | one/1/single/a(n)
<br/>
创建本地词条的过程主要是提取原始项中文字、频率、（除去音调的）拼音，接着把中文字和频率加入到对应的拼英文件当中。具体代码实现过程见BuildData文件夹，本地词条的创建结果保存在WordLab下。<br/>另外，有两个点做重点阐述：<br/>
1.频次分最小为1，最大为 10，而原始文件频率大到几百万，故需要对频率的范围做缩小处理。将频率的值/5000，得到得分，如果结果等于0的话，那么取1，如果大于10的话，那么取10。<br/>
2.提取拼音的过程略微复制。提取每一行第五列内容，再根据/分割字符串，返回一个list，里面包含de di2 di4。若最后一个是数字，则返回范围是[:-1]的内容<br/>

## 2.实现拼音查找汉字 
### 功能说明
 - 如果输入是⼀个完整的拼音（判断标准为有对应的词典⽂件），则返回该拼音下所有的汉字，按照频率分数从高到低，若两个字⾼频分数⼀样，则按照词典⽂件中的顺序。
- 如果输入不是⼀个完整的拼音（判断标准为没有对应的词典⽂件），则返回所有前缀与输⼊相同的拼音的汉字中最⾼频次的10个，具体排序为：频率越高排在越前，频率相同的情况下按照拼音字母序的排在前，频率相同拼音字母序也相同按照文件中的顺序排列。
### 拼音查找汉字的实现
根据输入条件：如果输入不是⼀个完整的拼音，则返回所有前缀与输入相同的拼音的汉字。利用这个条件，将本地文件路径的拼音构建出前缀树。前缀树的数据结构：
```go
type MyInputMethod struct {
	filename []string
	children [26]*MyInputMethod
}
```
前缀树示意图如下，只展示了部分内容：
![前缀树示意图](https://github.com/tiffanyXiaoqing/Pinyin-input-method-library/blob/master/image/TireTree.png)
<br/>这里的前缀树设计的特别之处在于设定字符串数组，用来保存每次经过该节点的文件名。例如，插入"ju.dat"文件，这个文件名会同时保存在j节点以及它的子节点u节点。这样，如果输入j不是⼀个完整的拼音，就可以直接查找它的字符串数组，得到关联的汉字。
创建好前缀树以后，接下来是查找过程。
- 第一步：首先看输入字符串是否为一个完整的拼音，即查找文件夹当中有没有拼音对应的本地文件。如果是完整拼音直接跳到第三步。
- 第二步：不是完整拼音的情况，根据输入拼音查找前缀树节点的字符串数组，它包含所有相关dat文件路径。
- 第三步：创建一个小根堆，将每个dat文件当中的内容加入到大小为topk的小根堆中。得到该拼音对应的频率前10的汉字字符，返回这些字符。

**对特殊情况的考虑** 
若输入拼音是特殊字符呢？若输入为空行呢？<br/>
首先，将有可能的大写字母全部转换为小写字母。接着，对于输入是特殊字符的情况，其acsii码会超出字母的范围，结果直接返回[输入有误，请重新输入]若输入为空行，可以检测到没有实际输入，返回[输入为空行]

**对排序顺序的设计**
 第一条：频次越高的排在越前⾯，这个是由小根堆实现的。重点在于第二三条比较难实现。即，在频率相同的情况下，保证汉字对应的拼音的字母序排列，和在原始文件中的顺序排列。
为此，我特别设计了我的Priority_queue节点结构：

```go
type Node struct {
	frequency int  //汉字频率
	Index int     //插入顺序
	content string  //汉字所在的文件名
}
```
在每次插入节点时，记录插入顺序。由于文件本身就是按照字典序排列的，这样插入顺序index变量既代表拼音序，又代表在原始文件中的顺序。
在得到的最小堆结果后，每次pop一个节点，节点有3个内容：内容字，插入顺序，字拼音长度。在频率相同的情况下，插入顺序的值大于前一个的值，则交换位子。最后反转结果，就能保证排序顺序的实现。

## 单元测试
在InputMethod_test.go文件对NewInputMethod 和 FindWords 两个方法做单元测试。
