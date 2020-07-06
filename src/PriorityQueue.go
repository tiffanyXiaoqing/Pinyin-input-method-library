/*
功能：本程序是小根堆的实现，根据频率大小建立小根堆。小根堆保存汉字频率，
作者：xiaoqing_tiffany@foxmail.com
*/
package main

type Node struct {
	frequency int  //汉字频率
	Index int     //插入顺序
	content string  //汉字所在的文件名
}

type NodeHeap []*Node

func (h NodeHeap) Len() int { return len(h) }

func (h NodeHeap) Less(i, j int) bool { return h[i].frequency <= h[j].frequency }

func (h NodeHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

func (h *NodeHeap) Push(x interface{}) {
	*h = append(*h, x.(*Node))
}

func (h *NodeHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}