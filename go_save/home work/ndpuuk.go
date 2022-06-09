package main

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"time"
)

type spukmode struct {
	strbook    []rune        //加密文字參列表
	opera      [9]int        //位址用作大小,value 用作條件
	unbook     [2][10]string //主要記錄索引
	trial      [10]string    //僅記錄贏的一方，陣列索引值 將與unbook對照
	eat        [5]int        //亂數與寫入 變數
	play       [5]int        //顯示與解密
	grave      []int         //墳牌
	sabetsu    [2]float64    //主要紀錄分數判斷
	puk        int           //抽卡容器
	pukct      string        //花色對照表
	event_loop bool          //主要迴圈判斷
	//牌組 對照表opera[0,1,2,3,4,5,6,7,8]int 主要記錄索引[ea ab bc cd de]string
	//dtd        time.Time
}

/*-------------------------------------------------
每次開局後主索引陣列 都必需判斷玩家的牌階{牌組陣列} 牌組花色陣列 頭牌陣列 索引陣列
比誰先將牌脫手...因此....
先比 牌階(role)>由加密值解給play顯示和eat?   索引值的必要?  花色>頭牌
再比 散牌的role>花色>頭牌
如何判斷輸贏? 贏的一方

顯示 牌組 花色 牌頭 大小
*/ //-----------------------------------------------
///*-------------------------------------------
type npupukmode struct {
	usermode
	npcmode
}
type usermode struct {
	upoint float64 //假分數顯示參考
	//同花順*10>四條*5>夫佬*4.5>同花*4>順子*3>三條*2.5>兩對*2>對子*1.5>散牌1.1~1.4
	uCR   int    // 對應opera(劇本)牌階 位址
	urole [2]int // 對應opera(劇本)牌組 value , main->role->uCR
	//由於花色其實可以用判斷就計算得出來，因此不列入變數
	ushepherd [2]int    //頭牌
	uETA      [5][5]int //加密卡牌
	/*		0	1	2	3	4x
	y0
	1
	2
	3
	4
	*/
	ub [10]string // 手牌索引陣列 僅用作系統檢查 而非 載入與寫入
}
type npcmode struct {
	npoint float64 //假分數
	nCR    int
	nrole  [2]int
	nETA   [5][5]int  //加密卡牌
	nb     [10]string //[daabbcceed] 必需與開局陣列索引 比較 牌組判定使用
}

//*/--------------------------------------------
func bookcopy(a *[]rune) { //0~61
	var letterRunes []rune = []rune("5ySKefrVz0FsT6vNDwLPtRWQ9x12Acd4HpqXU8uI3Jgh7BCiYOGabMklmjnoEZ") //26+26+10
	//fmt.Printf("\nletter: %c\n", letterRunes)
	*a = letterRunes
}
func whiteopera(a *[9]int) {
	var perien int //b 索引 perien亂數
	for index, _ := range a {
		perien = rand.Intn(99999999)
		//for i := 0; i < index && b < index; {
		for i := 0; i < index; {
			if a[i] == perien {
				perien = rand.Intn(99999999)
				i = 0
			} else if !(a[i] == perien) {
				i++
			}
		} //核對
		a[index] = perien
		//給值
	}
}

func re_rand(b *int, e []int) {
	//因為grave 不是指針式資料型態，所以re_rand函數"不能寫入"又每次執行都會回到原本的型態
	//因此有必要在re_rand函數外的main函數或者在main裡執行另外的函數，來去更新grave
	var btA_A, xxx int // 已知 函數判斷 永遠是以最初狀態而非動態
	//var nononn bool
	//btA_A = rand.Intn(51) + 1 //41+999999999 剛好一個亂數周期 每個數的機率會比較公平
	btA_A = 1
	if len(e) == 0 {
		btA_A = rand.Intn(51) + 1
		//fmt.Printf("grave=nil ,bt 抽到=%d AAAA\n", btA_A)
	} else if len(e) > 0 { //反之如果墳排有值
		for xxx < len(e) { //製造迴圈
			for btA_A == e[xxx] || btA_A == 0 { //若是 亂數 =ex
				btA_A = rand.Intn(51) + 1
				//fmt.Printf("AAAA發現相同e[%d]=%d ,因此RE!抽到=%d AAAA\n", xxx, e[xxx], btA_A)
				xxx = 0
			}
			if !(btA_A == e[xxx]) {
				xxx++
			}
		} //[40]
	}
	*b = btA_A
}
func SSS_rstring(a *[5]int, b *[5]int, x *[10]string, i int, f []rune) {
	var t rune
	var p int
	e := make([]rune, 30) //長度n  影響加密長度和迴圈次數

	//b[seppuku] = rand.Intn(51) //卡牌  因為是測試，事後組裝，所以才暫時套用執行，後續......
	//b[seppuku] = i
	//gg = rand.Intn(51)
	p = 1
	/*---------------dbug
	for i, _ := range a {
		a[i] = b[seppuku]
	}
	*/ //-------------------
	switch i % 5 {
	case 0:
		a[i] = (b[i] * b[i]) //eat[][5] bebe[]
	case 1:
		a[i] = (b[i] + 11)
	case 2:
		a[i] = (b[i] * 4) //
	case 3:
		a[i] = (b[i] * 11)
	case 4:
		a[i] = (b[i] + 22)
	}
	//-------------------
	for index, _ := range e { //    x^2+2x+1=49     => [2] 根號x +  +1  //len30
		//seppuku = time.Now()
		//pen = rand.Intn((seppuku.Nanosecond() + 47) % 62)
		//pen = int(math.Abs(float64(pen))) // 利用浮點數的abs轉乘int
		//pen = int(math.Abs(float64(rand.Intn((a.Nanosecond()+9)%62) + 1)))
		// 上面亂數會產生負數問題(原因是溢位),+9是為了公平地將其機率平均分配，可時間種子就是在亂數分配無論如何都會溢位
		//pen = rand.Intn((61)) //0~61 為y     y*z:{1~52}
		p = (p * (a[i])) % 62
		//   (x+1)^2(x+2)      (x+1)(x+2)
		//d[index] = p //						(pen)*(n+1)*(n + 2)		a  :56 1407	1551   =>  根號((a-3)/d)
		// 溢位的原因 rand.Intn()=0 在rand.go 154行 有乘上-n  至於為甚麼會有負數形成，目前沒頭緒(趕工)，
		// 目前解決辦法使用第73行可以穩定生成亂數避免負數形成
		e[index] = f[p]       //999999%62  0~61
		for i, _ := range f { ///len62
			if i < (len(f) - 1) { //  5*30*61 首次加總迴圈次數9150
				t = f[i]
				f[i] = f[i+1]
				f[i+1] = t
			} //	 1 2 3 4 0
		} // 	[0 1 2 3 4]  左移排序法
		//a[i] = p
		//fmt.Printf("b= %c\n", b)
		//b[i] = letterRunes[rand.Intn(len(letterRunes))]
		//len(letterRunes)=62=26+26+10y
	} //產生整數與字串的演算加密  eat[0~29]
	//fmt.Printf("d= %c\n", e)
	//x[i] = string(e)
	if i == 0 {
		x[i] = string(e) //x[i]={30s}
		x[len(x)-i-1] = x[i]
	} else if i == 1 {
		x[i] = string(e)
		x[i+1] = x[i]
	} else if i == 2 {
		x[i+1] = string(e)
		x[i+2] = x[i+1]
	} else if i == 3 {
		x[i+2] = string(e)
		x[i+3] = x[i+2]
	} else if i == 4 {
		x[i+3] = string(e)
		x[i+4] = x[i+3]
	}
}

func rokkuokaijosuru(d *[5]int, i int) {
	switch i % 5 {
	case 0:
		d[i] = int(math.Sqrt(float64(d[i]))) //1936 44
	case 1:
		d[i] = (d[i] - 11) //55   44
	case 2:
		d[i] = (d[i] / 4) //176	44
	case 3:
		d[i] = (d[i] / 11) //484  44
	case 4:
		d[i] = (d[i] - 22) //66	44
	}
}
func suit(b string, c int) { //顯示 數值 和花色
	if c > 0 && c < 53 {
		//symbol = puk / 13
		if c%13 == 0 {
			b = "A"
		} else if c%13 == 12 {
			b = "K"
		} else if c%13 == 11 {
			b = "Q"
		} else if c%13 == 10 {
			b = "J"
		} else {
			b = strconv.Itoa((c % 13) + 1)
		}
		if (c/13) == 0 || ((c/13) == 1 && c%13 == 0) {
			//if(puk<14)
			fmt.Printf("[\u2663 %s] ", b) //幸運草 0
		} else if (c/13) == 1 || ((c/13) == 2 && c%13 == 0) {
			fmt.Printf("[\u2666 %s] ", b) //方塊 1
		} else if (c/13) == 2 || ((c/13) == 3 && c%13 == 0) {
			fmt.Printf("[\u2665 %s] ", b)
		} else if (c/13) == 3 || ((c/13) == 4 && c%13 == 0) {
			fmt.Printf("[\u2660 %s] ", b)
		}
		/*--------------------
		fmt.Println("\u2660") //黑桃 3
		fmt.Println("\u2665") //紅心 2
		fmt.Println("\u2666") //方塊 1
		fmt.Println("\u2663") //幸運草 0
		*/
	}
}

func usersort(a *[5]int) {
	var c, b, minad, Mn int
	// 0  =A   1 累計
	//c = 4294967296
	for index, _ := range a {
		c = 4294967296
		for Mn = index; Mn < len(a); Mn++ {
			if (c == 4294967296 || c%13 > a[Mn]%13) && !(a[Mn]%13 == 0) {
				//c會大於Mn , 因此c會儲存較小值
				c = a[Mn]
				minad = Mn
				//fmt.Printf("\n搜尋c :%d ,user[%d]=%d ,最小值位址=%d", c, Mn, *a, minad)
			} else if a[Mn]%13 == 0 {
				//fmt.Printf("\nAAAAA搜尋c :%d ,user[%d]=%d ,最小值位址=%d", c, Mn, *a, minad)
			}
		}
		if c%13 == 0 {
		} else if !(c == 4294967296) {
			b = a[index]
			a[index] = c //  存入搜尋到的最小值
			a[minad] = b //  將原本的位置存到剛剛搬動的值
			c = 4294967296
			//fmt.Printf("\n 進入a撲克牌排序 :%d , len:%d", *a, len(a))
		} // [27]1  [13]0  47]8  35]9  5]5
		/*-------------------------------------------------
		for Mn = index; Mn < len(a); Mn++ {
			if c > a[Mn] {
				c = a[Mn]
				minad = Mn
				//fmt.Printf("\n排序c :%d ,user[%d]=%d ,最小值位址=%d", c, Mn, a, minad)
			} //-----min
		}*/
	}
	//a a a a 33
	//a A x 5 33     3  <=jzero[]    2
	//a a 1 5 33     2			     3
	//a 1 5 7 33     1  c mov =>     4
	//0 1 2 3 4
	//minad = jzero[4]
	c = 4294967296
}
func gravesort(a []int) {
	//var b = true
	var c, maxad, Gx int
	//var AT []int
	/*
		fmt.Printf("進入撲克牌排序 :%d", a)
		fmt.Printf("進入撲克牌排序 :%t", b)
		fmt.Printf("進入撲克牌排序 :%d", c)
	*/
	c = 4294967296
	//for b {
	//for index = 0; index < len(a); index++ {
	for index, _ := range a {
		for Gx = 0; Gx < len(a)-index; Gx++ {
			if c > a[Gx] {
				c = a[Gx]
				maxad = Gx
				//fmt.Printf("\n此時的c :%d ,a[%d]=%d ,最小值位址=%d", c, Gx, a, minad)
			} //-----max
		}
		if !(c == 4294967296) {
			a = append(a, 0)
			a[len(a)-1] = a[len(a)-index-2]
			a[len(a)-index-2] = c
			a[maxad] = a[len(a)-1]
			//a = append(a[:len(a)], a[len(a)+1:]...)
			a = a[:len(a)-1] //尾部刪除
			c = 4294967296
			//fmt.Printf("\n 進入a撲克牌排序 :%d , len:%d", a, len(a))
		}
	}
	//fmt.Printf("\n Final a撲克牌排序 :%d , len:%d", a, len(a))
	//b = false
	//}
}

//p_class_test(UN.usermode.uETA, &UN.usermode.ub, &UN.usermode.ushepherd, pc.opera, &UN.usermode.urole, pc.unbook, key)
func p_class_test(a [5][5]int, b *[10]string, c *[2]int, d [9]int, e *[2]int, f [2][10]string, kkey string) {
	/*	c散牌累計值 參考
		z花色值 1.1~1.4
		y同數累計值 1~4
		gototal[5]{y,z,y,z,c}-------*/
	var gototal [5]float64 //花色陣列  僅使用判斷"花與順"的數值
	var totol float64
	var koko [5]int //卡牌陣列
	var shadow int  //花色值

	for index, v := range a {
		for j, _ := range v {
			koko[j] = a[index][j]
			rokkuokaijosuru(&koko, j)
		}
	} // 卡牌解密 將其數值放在koko陣列  koko是卡牌陣列
	//----判斷陣列內部花色 並且累計hana
	for index, _ := range koko {
		if koko[index] <= 13 {
			gototal[index] = 1.0001
		} else if koko[index] <= 26 && koko[index] > 13 {
			gototal[index] = 1.001
		} else if koko[index] <= 39 && koko[index] > 26 {
			gototal[index] = 1.01
		} else if koko[index] <= 52 && koko[index] > 39 {
			gototal[index] = 1.1
		}
		//opera      [9]int
		if index == 0 {
		} else if gototal[index-1] == gototal[index] { //同花值累計
			totol = totol + gototal[index]
			//fmt.Printf("totol= %f ,index=%d\n", totol, index)
		}
		if index == 4 {
			//fmt.Printf("index=%d, len(gototal)-1=%d\n", index, len(gototal)-1)
			//fmt.Printf("gototal[index]=%f ,gototal[len(gototal)-1]=%f\n", gototal[index], gototal[len(gototal)-1])
			if gototal[len(gototal)-index] == gototal[len(gototal)-1] { //0 4
				totol = totol + gototal[index]
				//fmt.Printf("totol= %f ,index=%d\n", totol, index)
			}
		}
		/*if (totol/gototal[4] == 5) && (totol < 5.005) && (totol > 5.004) {
			totol = 5
			fmt.Printf("\n!!!!!!!!!!!!!!!!!!!totol= 5!!!!!!!!!!!!!!!!!!!!\n")
		}*/ //參考 IEEE 754 做浮點數運算總是會有誤差 因此判斷上需要注意

	} // *c 同花值累計尚未判別
	//---------------判斷同花順
	//0  1  2  3  4
	//2        K  A
	//ea ab bc cd de
	if !(koko[0]%13 == 1 && koko[3]%13 == 12 && koko[4]%13 == 0) { //KA2
		for index, _ := range a {
			if koko[0]%13 == 1 && !(index == 4) { // 2 X X x !2避免下一行執行錯誤
				if koko[index+1]%13-koko[index]%13 == 1 { // 確定間隔之間-1
					shadow++
					//e[0]++
					//fmt.Printf("bbbbbbbbbbbbbbb\n")
				} //23456 次二
			} else if koko[4]%13 == 0 && koko[3]%13 == 4 { //確認xxx5A "順"最大
				shadow++
				//e[0]++
				//fmt.Printf("cccccccccccccc\n")
			} // 2345A 23456
			if !(koko[0]%13 == 1) && !(index == 4) { //確認頭牌沒有2 [!2 x x x ]
				if (koko[index+1]%13)-(koko[index]%13) == 1 {
					shadow++
					//e[0]++
					//fmt.Printf("ddddddddddddddd\n")
				}
			} //10JQKA~34567
		}
		//fmt.Printf("upoint =%f", *c) // UN.usermode.ub[ea ab bc cd de]   pc.unbook[ea ab bc cd de]
		if (totol/gototal[4] == 5) && shadow == 4 {
			if koko[4]%13 == 0 { //A
				e[0] = d[0] //
			} else if koko[4]%13 == 5 { //23456
				e[0] = d[0]
			} else if koko[0] <= 9 && koko[0] >= 2 {
				e[0] = d[0]
			}
			fmt.Printf("同花順, UN.usermode.ushepherd= %d \n", e[0])
		} else if (totol/gototal[4] == 5) && shadow != 4 {
			e[0] = d[3]
			fmt.Printf("同花, UN.usermode.ushepherd= %d \n", e[0])
		} else if shadow == 4 && !(totol/gototal[4] == 5) {
			fmt.Printf("koko=%v\n", koko)
			if koko[4]%13 == 0 { //A
				e[0] = d[4]
				//fmt.Printf("1. else if (koko[4])")
			} else if koko[4]%13 == 5 { //23456
				e[0] = d[4]
				//fmt.Printf("2.else if (koko[4])")
			} else if (koko[0]%13 <= 9) && (koko[0]%13 >= 2) {
				e[0] = d[4]
			}
			fmt.Printf("蛇(順子), UN.usermode.ushepherd= %v \n", e[0])
		}
		totol = 0
		/*for index, _ := range koko {
			fmt.Printf("koko[%d]=%d\n", index, koko[index])
		}*/
		for index, _ := range koko {
			if (index != 4) && (koko[index]%13 == koko[index+1]%13) {
				if (c[0] == 0) || (c[0]%13 == koko[index]%13) {
					//!!!!!!!!已知 (c[0] == koko[index]%13)   18!=5
					c[0] = koko[index+1]
					e[0]++
					//fmt.Printf("e[0]=%d\n", e[0])
				} else if c[1] == 0 || (c[1]%13 == koko[index]%13) {
					c[1] = koko[index+1]
					e[1]++
					//fmt.Printf("e[1]=%d\n", e[1])
				}
			} else if (index == len(koko)-1) && (koko[index]%13 == koko[index-len(koko)+1]%13) {
				if (c[0] == 0) || (c[0]%13 == koko[index]%13) { //4
					c[0] = koko[index]
					e[0]++
					//.Printf("c[0]=%d\n", e[0])
				} else if c[1] == 0 || (c[1]%13 == koko[index]%13) {
					c[1] = koko[index]
					e[1]++
					//fmt.Printf("c[1]=%d\n", e[1])
				}
			}
		}
		if (e[0] != d[0]) && (e[0] != d[3]) && (e[0] != d[4]) {
			switch {
			case e[0] == 3:
				{
					e[0] = d[1]
					//fmt.Printf("四條, 牌組= %v \n", e[0])
					fmt.Printf("四條\n")
				}
			case (e[0] == 2 && e[1] == 1) || (e[0] == 1 && e[1] == 2):
				{
					e[0] = d[2]
					//fmt.Printf("夫佬, 牌組= %v \n", e[0])
					fmt.Printf("夫佬\n")
				}
			case e[0] == 2 && e[1] == 0:
				{
					e[0] = d[5]
					e[1] = d[8]
					//fmt.Printf("三條, 牌組= %v \n", e[0])
					fmt.Printf("三條\n")

				}
			case e[0] == 1 && e[1] == 1:
				{
					e[0] = d[6]
					e[1] = d[6]
					//fmt.Printf("兩對, 牌組= %v \n", e[0])
					fmt.Printf("兩對\n")
				}
			case e[0] == 1 && e[1] == 0:
				{
					e[0] = d[7]
					e[1] = d[8]
					//fmt.Printf("對子, 牌組= %v \n", e[0])
					fmt.Printf("對子\n")
				}
			default:
				{
					e[0] = d[8]
					e[1] = d[8]
					fmt.Printf("散牌\n")
					//fmt.Printf("散牌, 牌組= %v \n", e[0])
				}
			}
		}
	} else if koko[0]%13 == 1 && koko[3]%13 == 12 && koko[4]%13 == 0 {
		fmt.Printf("!!KA2不成順!! \n")
	}
	if (e[0] == d[0]) || (e[0] == d[3]) || (e[0] == d[4]) {
		c[0] = koko[4] //牌頭
		if kkey == "2" {
			for index, v := range f {
				if index == 0 { //
					for j, _ := range v {
						b[j] = f[0][j]
						//fmt.Printf("b[%d]=%s\n", j, b[j])
					}
				}
			}
		}
	}
	fmt.Printf("ushepherd(頭牌):%v , urole:%v\n", c, e) //
	/*---------------散牌
	for index, _ := range a { //
		if a[index] > b[index] {
			*c = *c + 1
		} else if a[index] < b[index] {
			*d = *d + 1
		}
	}*/
	//--------------user 與 npc 比較計分
}
func main() {
	rand.Seed(time.Now().UnixNano())
	var key string
	var pc spukmode
	var UN npupukmode
	bookcopy(&pc.strbook)
	//UN.usermode.ub[0]="155"
	pc.event_loop = true
	for pc.event_loop {
		if len(pc.grave) >= 50 {
			pc.grave = nil
			fmt.Printf("\n!!!! 牌墓 空間已經占滿, 因此重新洗牌!!!! \n")
		}
		fmt.Printf("\n y= 離開; 1=開局; 2牌組與狀態; 3=xxx發牌xxx; 4= xxx比輸贏xxx:")
		fmt.Scanln(&key)
		switch key {
		case "y":
			{
				pc.event_loop = false
			}
		case "1":
			{
				whiteopera(&pc.opera)
				fmt.Printf("opera : %v\n", pc.opera)
				pc.eat[0] = 1
				pc.eat[1] = 2
				pc.eat[2] = 3
				pc.eat[3] = 4
				pc.eat[4] = 13
				for index, _ := range pc.eat {
					//re_rand(&pc.puk, pc.grave)
					//pc.eat[index] = pc.puk
					pc.grave = append(pc.grave, pc.eat[index])
				} //----------亂數[5]產生 與 墳場進入
				// 對照墳牌 重新搓牌
				usersort(&pc.eat) // 排序
				//fmt.Printf("1.排序後的user :%v\n", pc.eat) //....debug
				//fmt.Printf("1.grave :%v\n", pc.grave)
				for index, _ := range pc.play {
					pc.play[index] = pc.eat[index]
					pc.eat[index] = 0
					suit(pc.pukct, pc.play[index])
				} //顯示花色
				for index, v := range UN.usermode.uETA {
					for m, _ := range v {
						SSS_rstring(&pc.eat, &pc.play, &UN.usermode.ub, m, pc.strbook)
						UN.usermode.uETA[index][m] = pc.eat[m]
					} // 5 x{5 }
				} // 加密
				for index, v := range pc.unbook {
					// 因此不對外執行 "SSS_rstring" =需要改成不對外的容器與函式=> "UN.usermode.ub"
					// ub 的參數會分成3類在三個資料庫
					for j, _ := range v {
						if index == 0 {
							pc.unbook[index][j] = UN.usermode.ub[j]
							UN.usermode.ub[j] = ""
						}
						//fmt.Printf("unbook[%d][%d]=%s\n", index, j, pc.unbook[index][j])
						//fmt.Printf("ub[%d]=%s\n", j, UN.usermode.ub[j])
						//理應sever=>sever data ,判定確立 sever=>user
						//data <=不對外=> sever <=對外=> user
					}
				} //將其陣列索引給予到pc主陣列
				// /*----------debug-------------------------------------
				fmt.Printf("+++++ user :%d \n", pc.play)
				//fmt.Printf("+++++ uETA :%d \n", UN.usermode.uETA)
				//fmt.Printf("+++++  ub  :%s \n", UN.usermode.ub)
				//---------------------------------------------------------------------------------
				for fapai, _ := range pc.eat {
					re_rand(&pc.puk, pc.grave)
					pc.eat[fapai] = pc.puk
					pc.grave = append(pc.grave, pc.eat[fapai])
				} //npc cord-----in grave
				//yipuyaya(&pc.puk, pc.grave, &pc.eat, &UN.npcmode.nb, pc.strbook)
				usersort(&pc.eat) //npc sort
				fmt.Printf("____npc :%v\n", pc.eat)
				for index, _ := range pc.play {
					pc.play[index] = pc.eat[index]
					pc.eat[index] = 0
					suit(pc.pukct, pc.play[index])
				} //顯示花色
				for fapai, v := range UN.npcmode.nETA { //nETA[fapai]: [v: [j]]
					for j, _ := range v {
						SSS_rstring(&pc.eat, &pc.play, &UN.npcmode.nb, j, pc.strbook)
						UN.npcmode.nETA[fapai][j] = pc.eat[j]
					}
				}
				for index, v := range pc.unbook {
					for j, _ := range v {
						if index == 1 {
							pc.unbook[index][j] = UN.npcmode.nb[j]
							UN.npcmode.nb[j] = ""
						}
						//fmt.Printf("unbook[%d][%d]=%s\n", index, j, pc.unbook[index][j])
						//fmt.Printf("nb[%d]=%s\n", j, UN.npcmode.nb[j])
					}
				}
				fmt.Printf("+++++ npc  :%d \n", pc.play)
				//fmt.Printf("+++++ nETA :%d \n", UN.npcmode.nETA)
				//fmt.Printf("+++++ nb   :%s \n", UN.npcmode.nb)
				gravesort(pc.grave)
				fmt.Printf("grave :%v\n", pc.grave)
				p_class_test(UN.usermode.uETA, &UN.usermode.ub, &UN.usermode.ushepherd, pc.opera, &UN.usermode.urole, pc.unbook, key)
				for i, _ := range pc.opera {
					if (UN.usermode.urole[0] != 0) && (UN.usermode.urole[0] == pc.opera[i]) {
						UN.usermode.uCR = i
					}
				}
				key = ""
			}
		case "2": //解密顯示 玩家的牌組和狀態
			{
				for index, v := range UN.usermode.uETA {
					for j, _ := range v {
						pc.play[j] = UN.usermode.uETA[index][j]
						rokkuokaijosuru(&pc.play, j)
					}
					suit(pc.pukct, pc.play[index])
					pc.play[index] = 0
				} // 卡牌解密 將其數值放在koko陣列  koko是卡牌陣列
				fmt.Println()
				for index, _ := range pc.grave {
					suit(pc.pukct, pc.grave[index]) //!!!!!墳牌是未加密狀態!!!!!
				}
				fmt.Println()
				switch UN.usermode.uCR {
				case 0:
					fmt.Printf("同花順\n")
				case 1:
					fmt.Printf("四條\n")
				case 2:
					fmt.Printf("夫佬\n")
				case 3:
					fmt.Printf("同花\n")
				case 4:
					fmt.Printf("順子\n")
				case 5:
					fmt.Printf("三條\n")
				case 6:
					fmt.Printf("兩對\n")
				case 7:
					fmt.Printf("對子\n")
				case 8:
					fmt.Printf("散牌\n")
				}
			}
		case "3":
			//fmt.Printf("grave :%v\n", pc.grave)
		}

		//for xxx, _ := range pukpc.grave {
		/*
			for index, v := range UN.usermode.uETA {
				for j, _ := range v {
					pc.play[j] = UN.usermode.uETA[index][j]
					rokkuokaijosuru(&pc.play, j)
				}
			}
			fmt.Printf("\n-----user :%v", pc.play)
			for index, v := range UN.npcmode.nETA {
				for j, _ := range v {
					pc.play[j] = UN.npcmode.nETA[index][j]
					rokkuokaijosuru(&pc.play, j)
				}
			}
			fmt.Printf("\n-----npc :%v\n---------grave :%v", pc.play, pc.grave)
		*/
	} //............for pc.event_loop{
}
