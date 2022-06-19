package main

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"time"
)

type spukmode struct {
	strbook         []rune        //加密文字參列表
	opera           [9]int        //位址用作大小,value 用作條件
	unbook          [2][10]string //主要記錄索引
	trial           [10]string    //僅記錄贏的一方，陣列索引值 將與unbook對照
	eat             [5]int        //亂數與寫入 變數 user
	neat            [5]int        //npc
	cheak           [5]int
	play            [5]int     //顯示與解密
	grave           []int      //墳牌
	GreatHolySpirit []int      //欺騙玩家的墳墓陣列  相關名詞是來自伊藤潤二漫畫
	sabetsu         [2]float64 //主要紀錄分數判斷
	puk             int        //抽卡容器
	pukct           string     //花色對照表
	event_loop      bool       //主要迴圈判斷
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
	ushepherd [2]int     //頭牌
	uETA      [5][5]int  //加密卡牌
	ub        [10]string // 手牌索引陣列 僅用作系統檢查 而非 載入與寫入
}
type npcmode struct {
	npoint    float64 //假分數
	nCR       int
	nrole     [2]int
	nETA      [5][5]int  //加密卡牌
	nshepherd [2]int     //頭牌
	nb        [10]string //[daabbcceed] 必需與開局陣列索引 比較 牌組判定使用
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
		perien = rand.Intn(99999999) + 1
		//for i := 0; i < index && b < index; {
		for i := 0; i < index; {
			if a[i] == perien {
				perien = rand.Intn(99999999) + 1
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

//howaito(&pc.play, &pc.eat, index 0 , m 01234)
func howaito(a *[5]int, b *[5]int, i int, j int) {
	switch j % 5 {
	case 0:
		a[j] = (b[i] * b[i]) //eat[][5] bebe[]  	a0  1
	case 1:
		a[j] = a[0] + (b[i] + 11) //				a1	13
	case 2:
		a[j] = a[1] + (b[i] * 4) //					a2	17
	case 3:
		a[j] = a[2] + (b[i] * 11) //				a3  28
	case 4:
		a[j] = a[3] + (b[i] + 22) //				a4  51
	}
}

/*
	//f[x]=x+ (y+22)     x-(y-22)
	x=  [0][0]   y*y												=1
		[0][1]	 y0/y0 - (y-11)										=(13-1)-11
		[0][2]	 (y0/y0 - (y1-11)) - (y/4)							=(17-13)/4
		[0][3]	 [(y0/y0 - (y1-11)) - (y2/4)]  -(y/11)				=(28-17)/11
		[0][4]	 {[(y0/y0 - (y1-11)) - (y2/4)]  -(y3/11)} - (y4-22) =(51-28)-22
		驗證公式為
		開根號a[0]
		(a[1]-a[0])-11
		(a[2]-a[1])/4
		(a[3]-a[2])/11
		(a[4]-a[3])-22
*/

// burakku(&pc.play, &UN.usermode.ub, m, pc.strbook)
func burakku(a *[5]int, x *[10]string, i int, f []rune) {
	var t rune //
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
	*/                        //-------------------
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

func rokkuokaijosuru(d *[5]int, i int, a *[5]int) {
	switch i % 5 {
	case 0:
		a[i] = int(math.Sqrt(float64(d[0]))) //1936 44
	case 1:
		a[i] = (d[1] - d[0]) - 11 //55   44
	case 2:
		a[i] = (d[2] - d[1]) / 4 //176	44
	case 3:
		a[i] = (d[3] - d[2]) / 11 //484  44
	case 4:
		a[i] = (d[4] - d[3]) - 22 //66	44   //f[x]=x+ (y+22)     x-(y-22)
	}
}

/*
	//f[x]=x+ (y+22)     x-(y-22)
	x=  [0][0]   y*y												=1				d0=1
		[0][1]	 y0/y0 - (y-11)										=(13-1)-11		d1=13
		[0][2]	 (y0/y0 - (y1-11)) - (y/4)							=(17-13)/4		d2=17
		[0][3]	 [(y0/y0 - (y1-11)) - (y2/4)]  -(y/11)				=(28-17)/11		d3=28
		[0][4]	 {[(y0/y0 - (y1-11)) - (y2/4)]  -(y3/11)} - (y4-22) =(51-28)-22		d4=51
		驗證公式為
		開根號a[0]
		(d1-a[0])-11
		(d2-a[1])/4
		(d3-a[2])/11
		(d4-a[3])-22
*/

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
func minpukid(a *[5]int, b int) int {
	var id, minmm int
	minmm = 65
	//if b != 4 {
	for iii := b; iii < len(a); iii++ {
		if (minmm-1)%13 > (a[iii]-1)%13 || ((minmm-1)%13 == (a[iii]-1)%13 && minmm > a[iii]) {
			//0~12  數值1[撲克2]~數值13[撲克A]  	      (相等	   &&  從遠始數值比大小)
			minmm = a[iii]
			id = iii
		} //
	}
	if minmm == 154 {
		return b
	}
	return id
	//} else {
	//return b
	//}  (minmm-1)%13 == (a[iii]-1)%13

	/*
		52 1 13 40 35
		1 52 13 40 35
		1 40 13 52 35
		1 40
	*/

}

func usersort(a *[5]int) {
	//var c, b, minad, Mn int
	var minad, c int
	// 0  =A   1 累計
	//c = 4294967296
	/*52 1 5 13 40
	(x-1)%13 for 找最小值
	/13
	//!!!!!!!!!!!!!!!!!按照牌面的大小順序排列 並且同花也由小到大排列!!!!!!!!!!!!!
	1	40	5	13	52


	1  40 5  13 52*/
	for index, _ := range a {
		minad = minpukid(a, index)
		//if minad < len(a) {
		c = a[index]
		a[index] = a[minad]
		a[minad] = c
		//fmt.Printf("!!!!!!!!!!!!!minad=%d,index=%d\n", minad, index)
		//fmt.Printf("a=%v\n", a)
		//}
	}

	/*
		for index, _ := range a {
			c = 4294967296
			for Mn = index; Mn < len(a); Mn++ { //40 27 14 1
				if c == 4294967296 || c%13 > a[Mn]%13 && !(a[Mn]%13 == 0) {
					//c會大於Mn , 因此c會儲存較小值
					//if c%13 == a[Mn]%13 && (c-1)/13 > (a[Mn]-1)/13 {
					c = a[Mn]
					minad = Mn
					//}
					fmt.Printf("\n搜尋c :%d ,user[%d]=%d ,最小值位址=%d", c, Mn, *a, minad)
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
				fmt.Printf("\nindex=%d 進入a撲克牌排序 :%d , len:%d", index, *a, len(a))
			} // [27]1  [13]0  47]8  35]9  5]5
			/*-------------------------------------------------
			for Mn = index; Mn < len(a); Mn++ {
				if c > a[Mn] {
					c = a[Mn]
					minad = Mn
					//fmt.Printf("\n排序c :%d ,user[%d]=%d ,最小值位址=%d", c, Mn, a, minad)
				} //-----min
			//}-------------------------------------------------
		}

		//a a a a 33
		//a A x 5 33     3  <=jzero[]    2
		//a a 1 5 33     2			     3
		//a 1 5 7 33     1  c mov =>     4
		//0 1 2 3 4
		//minad = jzero[4]
		c = 4294967296*/
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

//p_class_test(UN.usermode.uETA, &UN.usermode.ub, &UN.usermode.ushepherd, pc.opera, &UN.usermode.urole, pc.unbook, key, &pc.neat)
func p_class_test(a [5][5]int, b *[10]string, c *[2]int, d [9]int, e *[2]int, f [2][10]string, kkey string, z *[5]int) {
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
			z[j] = a[index][j]
		}
		rokkuokaijosuru(z, index, &koko)
	} // 卡牌解密 將其數值放在koko陣列  koko是卡牌陣列
	//----判斷陣列內部花色 並且累計hana
	for index, _ := range koko {
		z[index] = 0
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
		if (totol/gototal[4] == 5) && shadow == 4 { // 判斷同花值 與連續間格1 的次數

			if koko[4]%13 == 0 { //A
				e[0] = d[0]          //c[0]
				if koko[0]%13 == 1 { //2
					c[0] = koko[4] // = A
				} else if koko[0] == 9 { //10
					c[0] = koko[0] //
				}
			} else if koko[4]%13 == 5 { //23456
				e[0] = d[0]    // = 2
				c[0] = koko[0] //'2'3456
			} else if (koko[0]%13 <= 9) && (koko[0]%13 >= 2) {
				e[0] = d[0]
				c[0] = koko[0]
				//fmt.Printf("  1.koko:%v, c[0]=%d, e[0]=%d", koko, c[0], e[0])
			}

			if kkey == "3" {
				for i, _ := range koko {
					if koko[i] > 0 {
						b[i*2] = f[0][i*2]
						b[i*2+1] = f[0][i*2+1]
						z[i] = i + 1
						// 01 12 23 34 45
					}
				}
			} else if kkey == "" {
				for i, _ := range koko {
					if koko[i] > 0 {
						b[i*2] = f[1][i*2]
						b[i*2+1] = f[1][i*2+1]
					}
				}
			}
			//fmt.Printf("  2.koko:%v, c[0]=%d, e[0]=%d\n", koko, c[0], e[0])
			//fmt.Printf("  同花順, UN.usermode.ushepherd= %d \n", c[0])
			fmt.Printf("  同花順  \n")
		} else if (totol/gototal[4] == 5) && shadow != 4 {
			e[0] = d[3]
			c[0] = koko[4]
			if kkey == "3" {
				for i, v := range f {
					if i == 0 {
						for j, _ := range v {
							b[j] = f[i][j]
						}
					}
				}
				for i, _ := range z {
					z[i] = i + 1
				}
			} else if kkey == "" {
				for i, v := range f {
					if i == 1 {
						for j, _ := range v {
							b[j] = f[i][j]
						}
					}
				}
			}
			//fmt.Printf("  同花, UN.usermode.ushepherd= %d \n", c[0])
			fmt.Printf("  同花  \n")
		} else if shadow == 4 && !(totol/gototal[4] == 5) {
			//fmt.Printf("koko=%v\n", koko)
			if koko[4]%13 == 0 { //A
				e[0] = d[4]
				if koko[0] == 1 {
					c[0] = koko[4]
				} else if koko[0] == 9 {
					c[0] = koko[0]
				}
				//fmt.Printf("1. else if (koko[4])")
			} else if koko[4]%13 == 5 { //23456
				e[0] = d[4]
				c[0] = koko[0] //'2'3456
				//fmt.Printf("2.else if (koko[4])")
			} else if (koko[0]%13 <= 9) && (koko[0]%13 >= 2) {
				e[0] = d[4]
				c[0] = koko[0]
			}
			if kkey == "3" {
				for i, _ := range koko {
					if koko[i] > 0 {
						b[i*2] = f[0][i*2]
						b[i*2+1] = f[0][i*2+1]
						z[i] = i + 1
					}
				}
			} else if kkey == "" {
				for i, _ := range koko {
					if koko[i] > 0 {
						b[i*2] = f[1][i*2]
						b[i*2+1] = f[1][i*2+1]
					}
				}
			}
			//fmt.Printf("  蛇(順子), UN.usermode.ushepherd= %v \n", c[0])
			fmt.Printf("   蛇(順子)  \n")
		}
		totol = 0
		/*for index, _ := range koko {
			fmt.Printf("koko[%d]=%d\n", index, koko[index])
		}*/
		//if shadow < 4 {
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
				z[index] = index + 1 //0 12 1 123
				z[index+1] = index + 2
			} else if (index == len(koko)-1) && (koko[index]%13 == koko[index-(len(koko)-1)]%13) {
				if (c[0] == 0) || (c[0]%13 == koko[index]%13) { //4
					c[0] = koko[index]
					e[0]++
					//fmt.Printf("c[0]=%d\n", e[0])
				} else if c[1] == 0 || (c[1]%13 == koko[index]%13) {
					c[1] = koko[index]
					e[1]++
					//fmt.Printf("c[1]=%d\n", e[1])
				}
			}
		}
		//}
		//fmt.Printf("!!!!!c[0]=%d, e[0]=%d, c[1]=%d, e[1]=%d!!!!!!", c[0], e[0], c[1], e[1])
		//fmt.Printf("\n!手牌相同反應陣列k=%d\n", z)
		if (e[0] != d[0]) && (e[0] != d[3]) && (e[0] != d[4]) {
			switch {
			case e[0] == 3:
				{
					e[0] = d[1]
					//fmt.Printf("四條, 牌組= %v \n", e[0])
					fmt.Printf("  四條\n")
					if kkey == "3" {
						for i, _ := range z {
							if z[i] > 0 {
								b[i*2] = f[0][i*2]
								b[i*2+1] = f[0][i*2+1]
							}
						}
					} else if kkey == "" {
						for i, _ := range z {
							if z[i] > 0 {
								b[i*2] = f[1][i*2]
								b[i*2+1] = f[1][i*2+1]
							}
						}
					}
					//  0 1 2 3 4
					// k[1 2 3 4 5]
					// b[01 23 45 67 89]          3*2+1  4*2+1
				}
			case (e[0] == 2 && e[1] == 1) || (e[0] == 1 && e[1] == 2):
				{
					if e[0] > e[1] {
						c[1] = 0
					} else if e[0] < e[1] {
						c[0] = c[1]
					}
					e[0] = d[2]
					//fmt.Printf("夫佬, 牌組= %v \n", e[0])
					fmt.Printf("  夫佬\n")
					if kkey == "3" {
						for i, _ := range z {
							if z[i] > 0 {
								b[i*2] = f[0][i*2]
								b[i*2+1] = f[0][i*2+1]
							}
						}
					} else if kkey == "" {
						for i, _ := range z {
							if z[i] > 0 {
								b[i*2] = f[1][i*2]
								b[i*2+1] = f[1][i*2+1]
							}
						}
					}
				}
			case e[0] == 2 && e[1] == 0:
				{
					e[0] = d[5]
					e[1] = d[8]
					//fmt.Printf("三條, 牌組= %v \n", e[0])
					fmt.Printf("  三條\n")
					if kkey == "3" {
						for i, _ := range z {
							if z[i] > 0 {
								b[i*2] = f[0][i*2]
								b[i*2+1] = f[0][i*2+1]
							}
						}
					} else if kkey == "" {
						for i, _ := range z {
							if z[i] > 0 {
								b[i*2] = f[1][i*2]
								b[i*2+1] = f[1][i*2+1]
							}
						}
					}
				}
			case e[0] == 1 && e[1] == 1:
				{ //c 的正確性來自排序
					e[0] = d[6]
					e[1] = d[6]
					//fmt.Printf("兩對, 牌組= %v \n", e[0])
					fmt.Printf("  兩對\n")
					if kkey == "3" {
						for i, _ := range z {
							if z[i] > 0 {
								b[i*2] = f[0][i*2]
								b[i*2+1] = f[0][i*2+1]
							}
						}
					} else if kkey == "" {
						for i, _ := range z {
							if z[i] > 0 {
								b[i*2] = f[1][i*2]
								b[i*2+1] = f[1][i*2+1]
							}
						}
					}
				}
			case e[0] == 1 && e[1] == 0:
				{ //c 的正確性來自排序
					e[0] = d[7]
					e[1] = 0
					//fmt.Printf("對子, 牌組= %v \n", e[0])
					fmt.Printf("  對子\n")
					if kkey == "3" {
						for i, _ := range z {
							if z[i] > 0 {
								b[i*2] = f[0][i*2]
								b[i*2+1] = f[0][i*2+1]
							}
						}
					} else if kkey == "" {
						for i, _ := range z {
							if z[i] > 0 {
								b[i*2] = f[1][i*2]
								b[i*2+1] = f[1][i*2+1]
							}
						}
					}
				}
			default:
				{
					e[0] = d[8]
					e[1] = 0
					fmt.Printf("  散牌\n")
					//fmt.Printf("散牌, 牌組= %v \n", e[0])
				}
			}
		}
	} else if koko[0]%13 == 1 && koko[3]%13 == 12 && koko[4]%13 == 0 {
		e[0] = d[8]
		e[1] = 0
		fmt.Printf("  散牌\n")
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
	//fmt.Printf("ushepherd(頭牌):%v , urole:%v\n", c, e) //
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

func allclear(a *int, b *[2]int, c *[2]int) {
	*a = 999999
	*b = [2]int{0, 0}
	*c = [2]int{0, 0}
}
func yoki(a *float64, b int) {
	switch b {
	case 5:
		*a = 1.5
	case 4:
		*a = 1.4
	case 3:
		*a = 1.3
	case 2:
		*a = 1.2
	case 1:
		*a = 1.1
	}

}
func crlevel(a float64, b float64) float64 {
	switch a {
	case 0:
		a = 10
	case 1:
		a = 5
	case 2:
		a = 4.5
	case 3:
		a = 4
	case 4:
		a = 3
	case 5:
		a = 2.5
	case 6:
		a = 2
	case 7:
		a = 1.5
	case 8:
		a = 1
	}
	switch b {
	case 0:
		b = 10
	case 1:
		b = 5
	case 2:
		b = 4.5
	case 3:
		b = 4
	case 4:
		b = 3
	case 5:
		b = 2.5
	case 6:
		b = 2
	case 7:
		b = 1.5
	case 8:
		b = 1
	}
	return a * b

}
func maxid(a [5]int) int { //找到陣列最大值 位址
	var max, id int
	for pupu, _ := range a {
		if a[pupu] > max {
			max = a[pupu]
			id = pupu
		}
	}

	return id
}

/*
func maxran(a *[5]int) int { //找到陣列最大值並且回傳，並且串改為0
	var max, id int
	for pupu, _ := range a {
		if a[pupu] > max {
			max = a[pupu]
			id = pupu
		}
	}
	a[id] = 0
	return max
}*/

//xxxvsxxx(UN.usermode.uETA, UN.npcmode.nETA, pc.eat, pc.neat, UN.usermode.uCR "e",
// 	       UN.npcmode.nCR "f", UN.usermode.ushepherd "g", UN.npcmode.nshepherd "h",
//         &pc.sabetsu "i", UN.usermode.upoint "j", UN.npcmode.npoint "k", x "xxx")
func xxxvsxxx(a [5][5]int, b [5][5]int, c [5]int, d [5]int, e int, f int, g [2]int, h [2]int, i *[2]float64, j *float64, k *float64, xxx int) {
	var play [5]int
	var id [2]int //請使用id來判斷將由誰出牌  0 1
	var str string
	var ul, nl, zzz float64
	//var ul, nl float64
	// 紀錄ucr 與ncr  避免計算分數無法計算
	var execution, ground int = 5, 5
	for index, _ := range id {
		id[index] = rand.Intn(99999999) //{x,x}
	}
	for i := 0; i < 1; {
		if id[i] == id[i+1] {
			id[i] = rand.Intn(99999999)
		} else if id[i] != id[i+1] {
			i++
		}
	}
	for index, v := range a {
		for aaaa, _ := range v {
			c[aaaa] = a[index][aaaa]
		}
		rokkuokaijosuru(&c, index, &play)
	}
	c = play
	//fmt.Printf("c:%v; play:%v\n", c, play)
	for index, v := range b {
		for aaaa, _ := range v {
			d[aaaa] = b[index][aaaa]
		}
		rokkuokaijosuru(&d, index, &play)
	}
	d = play
	play = [5]int{0, 0, 0, 0, 0}
	//fmt.Printf("d:%v; play:%v\n", d, play)
	// 大老二比誰先脫手手中牌c是玩家 d是虛擬
	// 另 同花順*10>四條*8>夫佬*6>同花*5>順子*4.5>三條*3.5>兩對*2.4>對子*2>散牌+1.1 為x
	// 另 c 為對手手上抽到的散牌 pass回數 1回+1.1
	// 對手手上牌組  為 z
	// 贏家手上為牌組的計分方式 x*z+c			10*10+0		2.5*1.5+ 1.3+1.2+1.1
	//
	// 經過檢查失敗 要改成for 並且有switch 玩家 和莊家電腦 要輪流出牌 大的一方可以一直出牌
	// 而莊家永遠都有優先出牌的優勢
	//fmt.Printf("execution=%d, ground=%d, play[0]=%d,xxx=%d", execution, ground, play[0], xxx)
	for execution > 0 && ground > 0 {
		//fmt.Printf("execution=%d, ground=%d, play[0]=%d\n", execution, ground, play[0])
		//fmt.Printf("user:%v ,e=%d; npc:%v ,f=%d\n", c, e, d, f)
		fmt.Scanln(&str) //....debug
		//因為是0~4 無需分順序 直接攤牌 xxx不需要輪替
		if f < e && f < 5 { // f= 0~4  e=1~?  只需牌階 牌頭{'花''值'} 分數
			zzz = crlevel(float64(f), float64(e))
			// 0=10 1=5 2=4.5
			nl = nl + zzz
			ground = ground - 5
			xxx = 2
		} else if f > e && e < 5 {
			zzz = crlevel(float64(f), float64(e))
			// 0=10 1=5 2=4.5
			ul = ul + zzz
			execution = execution - 5
			xxx = 2
		} else if f == e && f == 0 || f == 4 { //相等階位 同花順與順子
			if (h[0]-1)/13 > (g[0]-1)/13 {
				zzz = crlevel(float64(f), float64(e))
				nl = nl + zzz
				ground = ground - 5
				xxx = 2
			} else if (h[0]-1)/13 < (g[0]-1)/13 { //花色對比
				zzz = crlevel(float64(f), float64(e))
				// 0=10 1=5 2=4.5
				ul = ul + zzz
				execution = execution - 5
				xxx = 2
			} else if (h[0])%13 == 0 && ((h[0]-1)/13 == (g[0]-1)/13) { //A
				zzz = crlevel(float64(f), float64(e))
				nl = nl + zzz
				ground = ground - 5
				xxx = 2
			} else if (g[0])%13 == 0 && ((h[0]-1)/13 == (g[0]-1)/13) {
				zzz = crlevel(float64(f), float64(e))
				// 0=10 1=5 2=4.5
				ul = ul + zzz
				execution = execution - 5
				xxx = 2
			} else if h[0]%13 == 1 && g[0]%13 != 0 { //2   !=A
				zzz = crlevel(float64(f), float64(e))
				nl = nl + zzz
				ground = ground - 5
				xxx = 2
			} else if g[0]%13 == 1 && h[0]%13 != 0 {
				zzz = crlevel(float64(f), float64(e))
				// 0=10 1=5 2=4.5
				ul = ul + zzz
				execution = execution - 5
				xxx = 2
			} else if h[0]%13 > g[0]%13 && (g[0]%13 != 0 && g[0]%13 != 1) { //
				zzz = crlevel(float64(f), float64(e))
				nl = nl + zzz
				ground = ground - 5
				xxx = 2
			} else if g[0]%13 > h[0]%13 && (h[0]%13 != 0 && h[0]%13 != 1) {
				zzz = crlevel(float64(f), float64(e))
				// 0=10 1=5 2=4.5
				ul = ul + zzz
				execution = execution - 5
				xxx = 2
			}
		} else if f == e && (f == 1 || f == 2) { //四條或夫老 同階 比牌頭的值即可
			if h[0]%13 > g[0]%13 { //
				zzz = crlevel(float64(f), float64(e))
				nl = nl + zzz
				ground = ground - 5
				xxx = 2
			} else if g[0]%13 > h[0]%13 {
				zzz = crlevel(float64(f), float64(e))
				// 0=10 1=5 2=4.5
				ul = ul + zzz
				execution = execution - 5
				xxx = 2
			}
		} else if f == e && (f == 3) { //同花 不需要順序 只要同花 比牌頭大小因此只比原始值
			if h[0] > g[0] {
				zzz = crlevel(float64(f), float64(e))
				nl = nl + zzz
				ground = ground - 5
				xxx = 2
			} else if g[0] > h[0] {
				zzz = crlevel(float64(f), float64(e))
				// 0=10 1=5 2=4.5
				ul = ul + zzz
				execution = execution - 5
				xxx = 2
			}
		}
		//if e == 8 && f == 8 { //散牌交鋒
		//因為牌階出牌不是全出，所以採用輪替選擇的方式一一出牌
		if xxx == 0 {
			switch f { //
			case 5: //   le0                 le1									le2
				if (play[0] == 0) || (play[0] != 0 && f < e && (e > 5 && e <= 8)) || (play[0] != 0 && f == e && h[0]%13 > play[0]%13) {
					zzz = crlevel(float64(f), float64(e))
					nl = nl + zzz       //該計算也只有在手上無牌 成立 ，只是先行計算 [隱]
					ground = ground - 3 // 條件達成可以放心出3張
					for iii, _ := range d {
						//由於場上無牌或牌階與牌頭的值勝過對方 因此可以放心將牌釋放
						if h[0]%13 == d[iii]%13 {
							play[0] = d[iii]
							d[iii] = 0

						}
					}
					play[0] = h[0]
					h[0] = 0
					if (e != 8 && g[0] == 0 && e != 6) || (e == 6 && g[1] == 0) {
						e = 8
					}
					//npc 手中的牌值去掉3個元素值並且將牌頭的值放在場上
					//牌階因為尚未比較因此先不將牌階轉換
					xxx = 1 //事後將其轉換出牌權 等待下一回
					//le1.....f=5 e=6~8 場上有牌
					//le2.....三條 同階 比牌頭 值
				} else {
					//zzz = crlevel(float64(f), float64(e))
					//ul = ul + zzz
					//execution = execution - 3
					e = 8
					xxx = 1
					//該牌階 與同階牌頭  都無法比，必需pass的動作 由於更高階的牌階判斷已經執行過
					//，因此對方牌階肯定是在5~8，然而是無法比因此視同階，將其轉換為散牌，等對方出牌
				}
			case 6: //成雙
				if (play[0] == 0) || ((play[0] != 0) && (e == 8)) || ((play[0] != 0) && ((e == 7) || (f == e)) && h[1]%13 > play[0]%13) {
					for iii, _ := range d {
						//由於場上無牌或牌階與牌頭的值勝過對方 因此可以放心將牌釋放
						if h[1]%13 == d[iii]%13 {
							play[0] = d[iii]
							d[iii] = 0
						}
					}
					if h[0] != 0 { //h0取代h1
						play[0] = h[1]
						h[1] = h[0] //
						h[0] = 0
						//h[0]==0  h[1]!=0
					} else if h[0] == 0 {
						zzz = crlevel(float64(f), float64(e))
						nl = nl + zzz //該計算也只有在手上無牌 成立 ，只是先行計算 [隱]
						play[0] = h[1]
						h[1] = 0 //h[h0]=0
					}
					ground = ground - 2 // 條件達成可以放心出3張
					//npc 手中的牌值去掉3個元素值並且將牌頭的值放在場上
					//牌階因為尚未比較因此先不將牌階轉換
					xxx = 1 //事後將其轉換出牌權 等待下一回
					//le1.....f=5 e=6~8 場上有牌
					//le2.....三條 同階 比牌頭 值
					if (e != 8 && g[0] == 0 && e != 6) || (e == 6 && g[1] == 0) {
						e = 8
					} //將對方牌面覆蓋 將其轉換牌階
				} else { //pass狀態
					if e == 6 && (g[1] != 0 && g[0] == 0 && play[0] != 0) {
						play[0] = 0
						xxx = 1
					} else if e == 6 && (g[1] == 0) {
						e = 8
						play[0] = 0
						xxx = 1
					} else if e == 7 || e == 5 {
						e = 8
						play[0] = 0
						xxx = 1
					}
					//zzz = crlevel(float64(f), float64(e))
					//ul = ul + zzz
					//execution = execution - 3
					//play[0]=g[0]

					//該牌階 與同階牌頭  都無法比，必需pass的動作 由於更高階的牌階判斷已經執行過
					//，因此對方牌階肯定是在5~8，然而是無法比因此視同階，將其轉換為散牌，等對方出牌
				}
			case 7: //對子
				if (play[0] == 0) || ((play[0] != 0) && (e == 8)) || ((play[0] != 0) && ((e == 6) || (f == e)) && h[0]%13 > play[0]%13) {
					for iii, _ := range d {
						//由於場上無牌或牌階與牌頭的值勝過對方 因此可以放心將牌釋放
						if h[0]%13 == d[iii]%13 {
							play[0] = d[iii]
							d[iii] = 0
						}
					}

					zzz = crlevel(float64(f), float64(e))
					nl = nl + zzz
					play[0] = h[0]
					h[0] = 0
					ground = ground - 2
					xxx = 1
					if (e != 8 && g[0] == 0 && e != 6) || (e == 6 && g[1] == 0) {
						e = 8
					}
				} else {
					if e == 6 && (g[1] != 0 && g[0] == 0 && play[0] != 0) {
						play[0] = 0
						xxx = 1
					} else if e == 6 && (g[1] == 0) {
						e = 8
						play[0] = 0
						xxx = 1
					} else if e == 7 || e == 5 {
						e = 8
						play[0] = 0
						xxx = 1
					}
				}
			case 8: //散牌
				if e == 8 || play[0] == 0 {
					if d[maxid(d)] != 0 && (play[0] == 0) || (play[0] != 0 && d[maxid(d)] > play[0]) {
						play[0] = d[maxid(d)]
						d[maxid(d)] = 0
						ground--
						xxx = 1
					} else if d[maxid(d)] < play[0] {
						xxx = 1
						ul = ul + 1.1
						play[0] = 0
					}
				} else if e != 8 {
					if e == 6 && (g[1] != 0 && g[0] == 0 && play[0] != 0) {
						play[0] = 0
						xxx = 1
					} else if e == 6 && (g[1] == 0) {
						e = 8
						play[0] = 0
						xxx = 1
					} else if e == 7 || e == 5 {
						e = 8
						play[0] = 0
						xxx = 1
					}
				}
			}
			//4.如有無法覆蓋的排則處以pass 轉換出牌權
			//  9  12 23 24 0
			//	6  11 22 33 0			44
			/*
				1.play=0 get
				2.get 換玩家比較 [牌階] [牌階相等比牌頭] [單張大小]
				3.玩家手禮可以蓋牌則出牌覆蓋 再次換出牌權
				4.如有無法覆蓋的排則處以pass 轉換出牌權   <============= 問題
				5.玩家再次出最大的牌
			*/
		} else if xxx == 1 {
			switch e {
			case 5:
				if (play[0] == 0) || (play[0] != 0 && f > e && (f > 5 && f <= 8)) || (play[0] != 0 && f == e && g[0]%13 > play[0]%13) {
					zzz = crlevel(float64(f), float64(e))
					ul = ul + zzz             //該計算也只有在手上無牌 成立 ，只是先行計算 [隱]
					execution = execution - 3 // 條件達成可以放心出3張
					for iii, _ := range c {
						//由於場上無牌或牌階與牌頭的值勝過對方 因此可以放心將牌釋放
						if g[0]%13 == c[iii]%13 {
							play[0] = c[iii]
							c[iii] = 0
						}
					}
					play[0] = g[0]
					g[0] = 0
					//npc 手中的牌值去掉3個元素值並且將牌頭的值放在場上
					//牌階因為尚未比較因此先不將牌階轉換
					xxx = 0
					if (f != 8 && (h[0] == 0) && f != 6) || f == 6 && h[1] == 0 {
						f = 8
					}
					//le1.....f=5 e=6~8 場上有牌
					//le2.....三條 同階 比牌頭 值
				} else {
					//zzz = crlevel(float64(f), float64(e))
					//ul = ul + zzz
					//execution = execution - 3
					f = 8
					xxx = 0
					//該牌階 與同階牌頭  都無法比，必需pass的動作 由於更高階的牌階判斷已經執行過
					//，因此對方牌階肯定是在5~8，然而是無法比因此不用設下條件視其為同階，將其轉換為散牌，等對方出牌
				}
			case 6:
				if (play[0] == 0) || (play[0] != 0 && (f == 8)) || (play[0] != 0 && ((f == 7) || (f == e)) && g[1]%13 > play[0]%13) {
					for iii, _ := range c {
						if g[1]%13 == c[iii]%13 {
							play[0] = c[iii]
							c[iii] = 0
						}
					}
					if g[0] != 0 {
						play[0] = g[1]
						g[1] = g[0]
						g[0] = 0
					} else if g[0] == 0 {
						zzz = crlevel(float64(f), float64(e))
						ul = ul + zzz
						play[0] = g[1]
						g[1] = 0
					}
					execution = execution - 2
					xxx = 0
					if (f != 8 && (h[0] == 0) && f != 6) || f == 6 && h[1] == 0 {
						f = 8
					}
				} else {
					if f == 6 && (h[1] != 0 && h[0] == 0 && play[0] != 0) {
						play[0] = 0
						xxx = 0
					} else if f == 6 && (h[1] == 0) {
						f = 8
						play[0] = 0
						xxx = 0
					} else if f == 7 || f == 5 {
						f = 8
						play[0] = 0
						xxx = 0
					}
				}
			case 7:
				if (play[0] == 0) || (play[0] != 0 && (f == 8)) || (play[0] != 0 && ((f == 6) || (f == e)) && g[0]%13 > play[0]%13) {
					for iii, _ := range d {
						//由於場上無牌或牌階與牌頭的值勝過對方 因此可以放心將牌釋放
						if g[0]%13 == c[iii]%13 {
							play[0] = c[iii]
							c[iii] = 0
						}
					}
					zzz = crlevel(float64(f), float64(e))
					ul = ul + zzz
					play[0] = g[0]
					g[0] = 0
					execution = execution - 2
					xxx = 0
					if (f != 8 && (h[0] == 0) && f != 6) || f == 6 && h[1] == 0 {
						f = 8
					}
				} else {
					if f == 6 && (h[1] != 0 && h[0] == 0 && play[0] != 0) {
						play[0] = 0
						xxx = 0
					} else if f == 6 && (h[1] == 0) {
						f = 8
						play[0] = 0
						xxx = 0
					} else if f == 7 || f == 5 {
						f = 8
						play[0] = 0
						xxx = 0
					}
				}
			case 8:
				if f == 8 || play[0] == 0 {
					if c[maxid(c)] != 0 && (play[0] == 0 || (play[0] != 0 && c[maxid(c)] > play[0])) {
						play[0] = c[maxid(c)]
						c[maxid(c)] = 0
						execution--
						xxx = 0
					} else if c[maxid(c)] < play[0] {
						xxx = 0
						nl = nl + 1.1
						play[0] = 0
					}
				} else if f != 8 {
					if f == 6 && (h[1] != 0 && h[0] == 0 && play[0] != 0) {
						play[0] = 0
						xxx = 0
					} else if f == 6 && (h[1] == 0) {
						f = 8
						play[0] = 0
						xxx = 0
					} else if f == 7 || f == 5 {
						f = 8
						play[0] = 0
						xxx = 0
					}
				}
			}
		}
		fmt.Printf("execution=%d, ground=%d, play[0]=%d, g[0]=%d\n", execution, ground, play[0], g[0])
		fmt.Printf("user:%v ,e=%d; npc:%v ,f=%d\n", c, e, d, f)
	}
	//fmt.Printf("\nexecution=%d; ground=%d\n", execution, ground)
	if execution == 0 {
		i[0] = i[0] + (ul * 2)
		*j = i[0] / 2
		//fmt.Printf("\n!!!!you win!!!!\n真分數=%f; 給玩家看的假分數=%f\n", i[0], j)
		fmt.Printf("\n!!!!you win!!!!\n 玩家的分數=%f\n", *j)
	} else if ground == 0 {
		i[1] = i[1] + (nl * 2)
		*k = i[1] / 2
		fmt.Printf("\n!!!!you lose!!!!\n NPC的分數=%f\n", *k)
		//fmt.Printf("\n!!!!you lose!!!!\n 真分數=%f; 給玩家看的假分數=%f\n", i[1], k)
	}

}

func main() {
	rand.Seed(time.Now().UnixNano())
	var key string
	var pc spukmode
	var UN npupukmode
	var hakari, x int
	bookcopy(&pc.strbook)
	//UN.usermode.ub[0]="155"
	pc.event_loop = true
	for pc.event_loop {
		fmt.Printf("\n y= 離開; 1=開局; 2.牌組與狀態; 3.發牌比輸贏:")
		fmt.Scanln(&key)
		switch key {
		case "y":
			{
				pc.event_loop = false
			}
		case "1":
			{
				x = rand.Intn(99999999) % 2 // 0=0 2 4 6 8    1=1 3 5 7 9
				if x == 0 {
					fmt.Printf("電腦先出牌\n")
				} else if x == 1 {
					fmt.Printf("玩家先出牌\n")
				}
				if len(pc.grave) >= 50 {
					pc.grave = nil
					pc.GreatHolySpirit = nil
					fmt.Printf("\n!!!! 牌墓 空間已經占滿, 因此重新洗牌!!!! \n")
				}
				whiteopera(&pc.opera)
				fmt.Printf("opera : %v\n", pc.opera)
				/*pc.eat[0] = 14
				pc.eat[1] = 41
				pc.eat[2] = 44
				pc.eat[3] = 31
				pc.eat[4] = 10*/
				for index, _ := range pc.eat {
					re_rand(&pc.puk, pc.grave)
					pc.eat[index] = pc.puk
					pc.puk = 0
					//fmt.Scanln(&pc.eat[index]) //......debug
					pc.grave = append(pc.grave, pc.eat[index])
					if index != 4 {
						pc.GreatHolySpirit = append(pc.GreatHolySpirit, pc.eat[index])
					}
				} //----------亂數[5]產生 與 墳場進入
				// 對照墳牌 重新搓牌
				usersort(&pc.eat)                      // 排序
				fmt.Printf("1.排序後的user :%v\n", pc.eat) //....debug
				//fmt.Printf("1.grave :%v\n", pc.grave)
				fmt.Printf("手牌: ")
				//---------------------------------------------
				for index, v := range UN.usermode.uETA {
					for m, _ := range v {
						howaito(&pc.play, &pc.eat, index, m) // 0  01234
						UN.usermode.uETA[index][m] = pc.play[m]
						//fmt.Printf("\nUN.usermode.uETA[%d][%d]=%d;  pc.eat[%d]=%d; pc.play[%d]=%d \n", index, m, UN.usermode.uETA[index][m], index, pc.eat[index], m, pc.play[m])
						pc.play[index] = pc.eat[index]
						burakku(&pc.play, &UN.usermode.ub, m, pc.strbook)
					} // 5 x{5 }
					pc.eat[index] = 0
					if index != 4 {
						suit(pc.pukct, pc.play[index])
					}
				} //------------------------------- 顯示花色 並且加密
				fmt.Printf("+++++ user :%d \n", pc.play)
				//fmt.Printf("\n")
				//fmt.Printf("+++++ ub :%s \n", UN.usermode.ub)
				for index, v := range pc.unbook {
					// 因此不對外執行 "SSS_rstring" =需要改成不對外的容器與函式=> "UN.usermode.ub"
					// ub 的參數會分成3類在三個資料庫
					for j, _ := range v {
						if index == 0 {
							//fmt.Printf("+++++  ub  :%s \n", UN.usermode.ub[j])
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
				fmt.Printf("+++++ uETA :%d \n", UN.usermode.uETA)
				//---------------------------------------------------------------------------------
				for fapai, _ := range pc.eat {
					re_rand(&pc.puk, pc.grave)
					pc.eat[fapai] = pc.puk
					pc.puk = 0
					pc.grave = append(pc.grave, pc.eat[fapai])
					if fapai != 4 {
						pc.GreatHolySpirit = append(pc.GreatHolySpirit, pc.eat[fapai])
					}
				} //npc cord-----in grave
				//yipuyaya(&pc.puk, pc.grave, &pc.eat, &UN.npcmode.nb, pc.strbook)
				usersort(&pc.eat) //npc sort
				fmt.Printf("\n____npc :%v\n", pc.eat)
				fmt.Printf("npc :")
				for index, v := range UN.npcmode.nETA {
					for m, _ := range v {
						howaito(&pc.play, &pc.eat, index, m) // 0  01234
						UN.npcmode.nETA[index][m] = pc.play[m]
						//fmt.Printf("\nUN.usermode.uETA[%d][%d]=%d;  pc.eat[%d]=%d; pc.play[%d]=%d \n", index, m, UN.usermode.uETA[index][m], index, pc.eat[index], m, pc.play[m])
						pc.play[index] = pc.eat[index]
						burakku(&pc.play, &UN.npcmode.nb, m, pc.strbook)
					} // 5 x{5 }
					pc.eat[index] = 0
					if index != 4 {
						suit(pc.pukct, pc.play[index])
					}
				} //------------------------------- 顯示花色 並且加密
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
				fmt.Printf("+++++ npc  :%d \n===============================\n", pc.play)
				//fmt.Printf("+++++ nETA :%d \n", UN.npcmode.nETA)
				//fmt.Printf("+++++ nb   :%s \n", UN.npcmode.nb)
				gravesort(pc.grave)
				//fmt.Printf("grave :%v\n", pc.grave)
				gravesort(pc.GreatHolySpirit)
				fmt.Printf("grave :%v\n", pc.GreatHolySpirit)
				key = ""
			}
		case "2": //解密顯示 玩家的牌組和狀態
			{

				if x == 0 {
					fmt.Printf("電腦先出牌\n")
				} else if x == 1 {
					fmt.Printf("玩家先出牌\n")
				}
				fmt.Printf("user: ")
				for index, v := range UN.usermode.uETA {
					for j, _ := range v {
						pc.eat[j] = UN.usermode.uETA[index][j]
						//fmt.Printf("UN.usermode.uETA[%d][%d]=%d ;pc.eat[%d]=%d\n", index, j, UN.usermode.uETA[index][j], j, pc.eat[j])
					}
					rokkuokaijosuru(&pc.eat, index, &pc.play)
					//fmt.Printf("UN.usermode.uETA[%d][0]=%d \n", index, UN.usermode.uETA[index][0])
					if index != 4 {
						suit(pc.pukct, pc.play[index])
					}
					pc.play[index] = 0
				} // 卡牌解密 將其數值放在koko陣列  koko是卡牌陣列
				fmt.Printf("\nnpc: ")
				for index, v := range UN.npcmode.nETA {
					for j, _ := range v {
						pc.eat[j] = UN.npcmode.nETA[index][j]
					}
					rokkuokaijosuru(&pc.eat, index, &pc.play)
					if index != 4 {
						suit(pc.pukct, pc.play[index])
					}
					pc.play[index] = 0
				}
				fmt.Printf("\n墳牌: ")
				for index, _ := range pc.GreatHolySpirit {
					suit(pc.pukct, pc.GreatHolySpirit[index]) //!!!!!墳牌是未加密狀態!!!!!
					//suit(pc.pukct, pc.grave[index]) //!!!!!墳牌是未加密狀態!!!!!
				}
				fmt.Printf("\n")
			}
		case "3":
			fmt.Printf("!!!!!!!!!!抽牌!!!!!!!!!\n")
			fmt.Printf("user: ")
			for index, v := range UN.usermode.uETA {
				for j, _ := range v {
					pc.eat[j] = UN.usermode.uETA[index][j]
				}
				rokkuokaijosuru(&pc.eat, index, &pc.play)
				suit(pc.pukct, pc.play[index])
				pc.play[index] = 0
				pc.eat[index] = 0
			} // 卡牌解密 將其數值放在koko陣列  koko是卡牌陣列
			p_class_test(UN.usermode.uETA, &UN.usermode.ub, &UN.usermode.ushepherd, pc.opera, &UN.usermode.urole, pc.unbook, key, &pc.cheak)
			key = ""
			for i, _ := range pc.opera {
				if (UN.usermode.urole[0] != 0) && (UN.usermode.urole[0] == pc.opera[i]) {
					UN.usermode.uCR = i
				} else if UN.usermode.urole[0] == 0 { //加密牌組值不可能為0 如果為零則顯示 yaku
					fmt.Printf("!!!!!!!!!!!!!!Yaku!!!!!!!!!!!!!!!\n")
				}
			}

			fmt.Printf("\nnpc: ")
			for index, v := range UN.npcmode.nETA {
				for j, _ := range v {
					pc.neat[j] = UN.npcmode.nETA[index][j]
				}
				rokkuokaijosuru(&pc.neat, index, &pc.play)
				suit(pc.pukct, pc.play[index])
				pc.play[index] = 0
				pc.neat[index] = 0
			}
			p_class_test(UN.npcmode.nETA, &UN.npcmode.nb, &UN.npcmode.nshepherd, pc.opera, &UN.npcmode.nrole, pc.unbook, key, &pc.play)
			for i, _ := range pc.opera {
				if (UN.npcmode.nrole[0] != 0) && (UN.npcmode.nrole[0] == pc.opera[i]) {
					UN.npcmode.nCR = i
				}
			}
			fmt.Printf("\n墳牌: ")
			for index, _ := range pc.grave {
				//suit(pc.pukct, pc.GreatHolySpirit[index]) //!!!!!墳牌是未加密狀態!!!!!
				suit(pc.pukct, pc.grave[index])
			}
			xxxvsxxx(UN.usermode.uETA, UN.npcmode.nETA, pc.eat, pc.neat, UN.usermode.uCR, UN.npcmode.nCR, UN.usermode.ushepherd, UN.npcmode.nshepherd, &pc.sabetsu, &UN.usermode.upoint, &UN.npcmode.npoint, x)
			/*------------------debug------------------------------
			for i,_ := range UN.usermode.ub{
				UN.usermode.ub[i]=""
			}
			*/ //------------------------------------------------
			//最後的檢查 目前停更  因為更新要刪減很多不必要的陣列 比方說eat 可以將neat 併用 ，或者將許多陣列併用成一個陣列
			if UN.usermode.uCR < 8 { //更新可以再加強 比方說 &UN.npcmode.nrole 可以用另外一種加密 並且作為2次確認檢查，檢查越多越準確，但也會增加執行消耗
				for i, _ := range pc.cheak { // 判定是牌組 內值 1 2 3 4 5
					if (pc.cheak[i] > 0) && (UN.usermode.ub[i*2] == pc.unbook[0][i*2]) && (UN.usermode.ub[i*2+1] == pc.unbook[0][i*2+1]) && (UN.usermode.ub[i*2] != "") && (UN.usermode.ub[i*2+1] != "") {
						//fmt.Printf("!!!!!!!!!!!!!!PASS!!!!!!!!!!!!!!!")
					} else {
						//fmt.Printf("!!!!!!!!!!!!!!NOGO!!!!!!!!!!!!!!!") // 只要UN.usermode.uCR < 8通過 NOGO五次就算作弊  因為牌已在事先發好加密好
						hakari++
					}
				}
			}
			if hakari == 5 { // 這個也是  也許多個sql 或者 新的?
				fmt.Printf("你作弊")
			}
			fmt.Printf("upoint=%f ; npoint=%f\n", UN.usermode.upoint, UN.npcmode.npoint)
			/*for in, v := range pc.unbook {....................debug
				for j, _ := range v {
					/*if in == 0 {
						//UN.usermode.ub[j] = pc.unbook[in][j]
						fmt.Printf("pc.unbook[%d][%d]=%s\n", in,j,pc.unbook[in][j])
					} else if in == 1 {
						//UN.npcmode.nb[j] = pc.unbook[in][j]
						fmt.Printf("UN.npcmode.nb=%s\n", UN.npcmode.nb[j])
					}
					//fmt.Printf("pc.unbook[%d][%d]=%s\n", in, j, pc.unbook[in][j])
				}
			}*/
			// 系統確認與判定 勝負結果
			/*
				uCR <=根據數值對應opera位址=urole  		[opera數值]=>opera[內碼......]
				ushepherd [0 1] == koko [0 1 2 3 4]    uETA=>eat=>  play[index]&&ub[index*2+1]==true 牌組位置索引
				UN.usermode.uCR int				牌階............UN.usermode.urole == pc.opera[]
				+UN.usermode.urole[2]int			牌組............UN.usermode.uETA >>> pc.opera[]
				+UN.usermode.ushepherd[2]int		頭牌............UN.usermode.uETA >>> pc.play >>> UN.usermode.ushepherd[2]
				+UN.usermode.ub[10]string		加密字串............UN.usermode.uETA >>> pc.play >>> pc.unbook >>> UN.usermode.ub
				+UN.usermode.uETA[5][5]int		加密數值............UN.usermode.uETA >>> pc.eat

				UN.usermode.upoint float64		............end
				pc.sabetsu[2] float64			............最好是在內部運算而成，或者在判別之前就有結果
			*/
			//fmt.Printf("uCR=%v, UN.usermode.urole=%v, ushepherd=%v", UN.usermode.uCR, UN.usermode.urole, UN.usermode.ushepherd) //ushepherd 牧羊人(牌頭)
			//fmt.Printf("ub: %s", UN.usermode.ub)
			//fmt.Printf("nb: %s", UN.npcmode.nb)
			//fmt.Printf("ub=%s\n un=%s\n", UN.usermode.ub, UN.npcmode.nb)
			//-----------------------------
			allclear(&UN.npcmode.nCR, &UN.npcmode.nrole, &UN.npcmode.nshepherd)
			allclear(&UN.usermode.uCR, &UN.usermode.urole, &UN.usermode.ushepherd)
			fmt.Printf("\n")
		}
		//fmt.Printf("grave :%v\n", pc.grave)

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
