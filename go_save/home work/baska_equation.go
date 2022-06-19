package main

import (
	"fmt"
)

type sedo struct {
	fx, a, b, x, y, n, i int
	it                   []int
}

func vrfp(a int, b int) int {
	var c int = 1
	var n = b
	for c = 1; n > 0; n-- {
		//fmt.Print("a=", a, ", n=", n)
		c = c * a
		//fmt.Println("\nfx=", c)

	}
	return c
}

func abn1(a int, b int, n int) {
	if a == 1 {
		fmt.Print("x")
	} else if a > 1 {
		fmt.Print(a, "x")
	}

	if a > 0 && b > 0 {
		fmt.Print("+")
	}

	if b == 1 {
		fmt.Print("y")
	} else if b > 1 {
		fmt.Print(b, "y")
	}

}

func abn(a int, b int, n int, fx int) {
	//var orz sedo =sedo{it[]:}

	if a == 0 && b == 0 {
		fmt.Println(0)

	} else if (!(a == 0) || !(b == 0)) && n == 0 {
		fmt.Println(1)
	}
	if n == 1 {
		fmt.Print("f(x)=")
		abn1(a, b, n)
		fmt.Print("\n")
	} else if n > 1 {
		fmt.Print("f(x)=(")
		abn1(a, b, n)
		fmt.Print(")^", n, "\n")
	}
}

func sisubaska(arr []int, a int, b int, i int) { arr[i] = arr[i] * a * b }

func baskasedo(a int, b int, n int, arr []int) { //<===  限制回傳一個整數參數
	var i int
	for i = 0; i <= n; i++ {
		//fmt.Println(arr[i] * vrfp(a, n-i) * vrfp(b, i))
		//fmt.Println("\na===", fx, x, y, i)
		//fx = fx * (a + b)
		//fmt.Println("b===", fx, x, y, i)
		/*
			fx=(5x+6y)(5x+6y)(5x+6y)     n=3
			   5^2[x^2] + 2*[5][6][xy] + 6^6*[y^2] (5x+6y)
			   5^3[x^3]	+ 3*[5^2][6][x^2][y] + 3*[5][6^2][x][y^2] +6^3[y^3]

			   i=0
			   it[i][a^n][b^i][x^n][y^i]+ it[i+1][a^n-1][b^i+1][x^n-1][y^i+1] + it[i+2][a^n-2][b^i+2][x^n-2][y^i+2] + it[i+3][a^n-3][b^i+3][x^n-3][y^i+3]

				a*a..../a  b/b

				vrfp(a *int, b *int)
		*/

		if n == 0 {
			fmt.Println(1)
		} else if (n == 1) && (i == 1) {
			fmt.Println(a, "x", "+", b, "y")
			sisubaska(arr, vrfp(a, n-i), vrfp(b, i), i)
			//arr[i] = arr[i] * vrfp(a, n-i) * vrfp(b, i)

			//fmt.Println("x", "+", "y")
		} else if n-i > 1 {
			if i == 0 {
				fmt.Print(arr[i]*vrfp(a, n-i)*vrfp(b, i), "[x^", n-i, "]+")
				sisubaska(arr, vrfp(a, n-i), vrfp(b, i), i)
				//fmt.Print("[x^", n-i, "]+")
			} else if i == 1 {
				fmt.Print(arr[i]*vrfp(a, n-i)*vrfp(b, i), "[x^", n-i, "][y]", "+")
				sisubaska(arr, vrfp(a, n-i), vrfp(b, i), i)
				//fmt.Print("[x^", n-i, "][y]", "+")
			} else {
				fmt.Print(arr[i]*vrfp(a, n-i)*vrfp(b, i), "[x^", n-i, "][y^", i, "]+")
				sisubaska(arr, vrfp(a, n-i), vrfp(b, i), i)
				//fmt.Print("[x^", n-i, "][y^", i, "]+")
			}
		} else if (n-i == 1) && n == 2 {
			fmt.Print(arr[i]*vrfp(a, n-i)*vrfp(b, i), "[x]", "[y]", "+")
			sisubaska(arr, vrfp(a, n-i), vrfp(b, i), i)
			//fmt.Print("[x]", "[y]", "+")
		} else if (n-i == 1) && !(n == 2) && !(n == 1) {
			fmt.Print(arr[i]*vrfp(a, n-i)*vrfp(b, i), "[x]", "[y^", i, "]+")
			sisubaska(arr, vrfp(a, n-i), vrfp(b, i), i)
			//fmt.Print("[x]", "[y^", i, "]+")
		} else if n-i == 0 {
			fmt.Print(arr[i]*vrfp(a, n-i)*vrfp(b, i), "[y^", i, "]")
			sisubaska(arr, vrfp(a, n-i), vrfp(b, i), i)
			//fmt.Print("[y^", i, "]")
		}

	}

}

func baska(n int, arr []int) []int {
	var b []int = []int{1}
	var bt, i int
	//var btt int = 1
	arr = append(arr, 1)
	//dandan = 0
	//var penso *int=&bt
	//fmt.Println("請決定次方數值為:")
	//fmt.Scanln(&n)
	fmt.Println("巴斯卡排列b陣列:", b)

	for i = 0; i < n; i++ { //^n
		arr = append(arr, 1)
		b = append(b, 1)
		for k, _ := range arr {
			//for dandan = 0; dandan < len(arr); dandan++ {
			if !(k == 0) && k < len(arr)-1 {
				bt = b[k-1] + b[k]
				arr[k] = bt
				//fmt.Println("此時的it 陣列", arr)
				//fmt.Println("此時的b 陣列", b)
			}
		}
		/*
			for dandan = 0; dandan < len(arr); dandan++ {
				bt = arr[dandan]
				b[dandan] = bt
			}
		*/
		//dandan=0
		for k, bt := range arr {
			bt = arr[k]
			b[k] = bt

			//fmt.Println("此時的it 陣列", arr)
			//fmt.Println("此時的dan", dandan)
			//fmt.Println("此時的bt", bt)
		}

		fmt.Println("巴斯卡排列b陣列:", b)
		//a11	a121	a1331
	}

	fmt.Println("最後it陣列的結果:", arr)
	return arr
}

func main() {
	var piricat sedo
	fmt.Print("f(x)=ax+by^n,目前不支援負數,請輸入a和b=")

	fmt.Scanln(&piricat.a, &piricat.b)
	fmt.Print("f(x)=ax+by^n,目前不支援負數,請輸入n次方=")
	fmt.Scanln(&piricat.n)
	abn(piricat.a, piricat.b, piricat.n, piricat.fx)

	//----------------組裝1....不是非常好算勉強，明明結構了變數，函數內的陣列卻沒有回傳給main
	//baska(piricat.n, piricat.it)
	piricat.it = baska(piricat.n, piricat.it)

	//---------------------組裝2.....ok
	fmt.Println("--------將係數帶入的式子--------------------------------------")
	baskasedo(piricat.a, piricat.b, piricat.n, piricat.it)
	fmt.Print("\n\na與b經過運算過的it陣列:", piricat.it, "\n ------------按下Enter結束程式執行-----------------")
	fmt.Scanln()
	//fmt.Println("\n", piricat.fx)
	// func函式  array陣列 for迴圈 *&指標  struct結構 append(a,1) a=a[:]
}
