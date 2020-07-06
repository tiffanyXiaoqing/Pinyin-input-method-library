/*
功能：NewMyInputMethod函数和FindWords函数的单元测试
作者：xiaoqing_tiffany@foxmail.com
修改内容：expectedFilename需要修改对应的相对路径
*/
package main

import (
	"io/ioutil"
	"testing"
)


func TestFindWords(t *testing.T) {

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
	im := NewMyInputMethod(processFile)

	testCases := []struct{
		spell 	 		 string
		expectedWords 	 []string
	}{
		{
			spell: "a",
			expectedWords: []string{"阿", "啊", "呵", "嗄", "锕", "吖"},
		},
		{
			spell: "b",
			expectedWords: []string{"把", "百", "白", "报", "保", "被", "北", "本", "必", "比"},
		},
	}

	for _, c := range testCases {
		words := im.FindWords(c.spell)
		if len(words) != len(c.expectedWords){
			t.Fatalf("Expected message did not align with what was written:\n\texpected: %q\\n\\tactual: %q",
				c.expectedWords,words)
		}

		for i := 0; i < len(words); i++{
			if words[i] != c.expectedWords[i]{
				t.Fatalf("Expected message did not align with what was written:\n\texpected: %q\\n\\tactual: %q",
					c.expectedWords,words)
			}
		}
	}
}

func TestNewMyInputMethod(t *testing.T) {
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
	im := NewMyInputMethod(processFile)

	testCases := []struct{
		spell 	 		 	 byte
		expectedFilename 	 []string
	}{
		{
			spell: 'a',
			expectedFilename: []string{
				"./BuildData/WordLab/a.dat",
				"./BuildData/WordLab/ai.dat",
				"./BuildData/WordLab/an.dat",
				"./BuildData/WordLab/ang.dat",
				"./BuildData/WordLab/ao.dat",},
		},
		{
			spell: 'e',
			expectedFilename: []string{
				"./BuildData/WordLab/e.dat",
				"./BuildData/WordLab/ei.dat",
				"./BuildData/WordLab/en.dat",
				"./BuildData/WordLab/er.dat",},
		},
	}

	for _, c := range testCases {
		n := c.spell-'a'
		actulFilename := im.children[n].filename
		if len(actulFilename) != len(c.expectedFilename){
			t.Fatalf("Expected filename did not align with what was written:\n\texpected: %q\\n\\tactual: %q",
				c.expectedFilename,actulFilename)
		}

		for i := 0; i < len(actulFilename); i++{
			if actulFilename[i] != c.expectedFilename[i]{
				t.Fatalf("Expected message did not align with what was written:\n\texpected: %q\\n\\tactual: %q",
					c.expectedFilename,actulFilename)
			}
		}

	}
}