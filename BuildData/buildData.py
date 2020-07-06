# encoding:utf-8
# 本脚本是针对汉字库的数据清洗文件
# 作者 hkjoe8191@gamail.com
# 源数据来自 https://lingua.mtsu.edu/chinese-computing/statistics/char/list.php?Which=TO

import os

class generateDatFile():
    def __init__(self,path='origin.dat'):
        assert os.path.exists(path)
        #校验文件路径的合法性
        self.path=path
        self.run()

    def run(self):
        self.ReadDatFile()
        self.ExtractWord()
        self.ExtractFreq()
        self.ExtractPinyinForEachWord()
        self.MapPinyinToWord()
        self.Write2DatFile()

    def ReadDatFile(self):
        #读取dat文件当中的每一行内容，建议直接从附件当中读取
        with open(self.path, 'r', encoding='utf-8') as f:
            self.LinesOfDatFile = f.readlines() # dat文件当中的所有行

    def ExtractWord(self):
        #提取汉字
        Word=[]
        for line in self.LinesOfDatFile:
            word=line.split('\t')[1] #通过split风格每一行，以第一行为例：'1	的	8302698 3.2074998098725	de/di2/di4	(possessive particle)/of, really and truly, aim/clear
            #分割的结果为list：word=['1', '的', '8302698', '3.2074998098725', 'de/di2/di4', '(possessive particle)/of, really and truly, aim/clear\n']
            Word.append(word)
        self.Word=Word

    def ExtractFreq(self):
        #提取频率值
        self.Word2Freq={}
        self.Freq=[]
        for index,line in enumerate(self.LinesOfDatFile):
            freq=line.split('\t')[2]
            freq = int(freq)/20000   #将频率缩小处理
            if freq > 10:
                freq = 10
            elif freq < 1:
                freq = 1
            freq = round(freq)
            self.Freq.append(freq)
            self.Word2Freq[self.Word[index]]=freq

    def ExtractPinyinForEachWord(self):
        self.Word2Pinyin={} #为每一个汉字建立一个字典，键为汉字，值为拼音(去除声调)
        self.AllPinyin=[] #所有的拼音
        for indexLine,line in enumerate(self.LinesOfDatFile):
            Pinyin=line.split('\t')[4] #以第一行为例，Pinyin='de/di2/di4'
            if Pinyin.__len__()==0:
                continue #如果这一行没有提取到拼音，那么跳过
            Pinyin=Pinyin.split('/') #得到的结果为Pinyin=['de', 'di2', 'di4']
            Tone=['1','2','3','4'] #声调
            for index,py in enumerate(Pinyin): #去除声调
                py=py.strip() #去除拼音末尾的\n 空格等
                try:
                    if py[-1] in Tone:
                        Pinyin[index]=py[:-1] #如果包含声调那么去除声调，例如py=di2,那么去除之后为di
                except:
                    print(self.LinesOfDatFile[indexLine])
            PinyinSet=list(set(Pinyin)) #考虑到拼音去音调之后会有重复的，那么此处去重
            self.AllPinyin.extend(PinyinSet)
            self.Word2Pinyin[self.Word[indexLine]]=PinyinSet #为每一个汉字对应上拼音

    def MapPinyinToWord(self):
        #同音字归为一个拼音下面
        self.AllPinyin=list(set(self.AllPinyin))
        self.Pinyin2Word={}
        for pinyin in self.AllPinyin:
            self.Pinyin2Word[pinyin]=[]#为么每一个拼音建立一个空的list值
        #self.Pinyin2Word=dict(self.AllPinyin,[[]]*self.AllPinyin.__len__())
        #例如{'yi':[]},{'liu':[]}
        for word,pinyin in self.Word2Pinyin.items():
            for py in pinyin:
                self.Pinyin2Word[py].append(word) #按照字-拼音从前往后遍历保证了同一拼音下的字按照频次从大到小排列

    def Write2DatFile(self,FileName='WordLab'):
        # if os.path.exists(FileName):
        #     os.remove(FileName) #如果存在同名文件夹，那么先全部删除在创建
        #os.mkdir(FileName)

        for pinyin,word in self.Pinyin2Word.items():
            Lines=''
            for wd in word:
                Lines+=wd+' '+str(self.Word2Freq[wd])+'\n'
            #循环之后写入的内容writeLines变为'的 310025\n得 300152...'方便一次性直接写入
            try:
                #需要先创建WordLab文件夹
                with open('WordLab\%s.dat'%pinyin.lower(),'w',encoding='utf-8') as f: #lower防止有大写出现
                    f.writelines(Lines)
            except:
                print(pinyin,word)

if __name__=='__main__':
    GDF=generateDatFile()