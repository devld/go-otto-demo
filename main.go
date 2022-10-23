package main

import (
	"fmt"
	"os"
	"time"

	"github.com/robertkrimen/otto"
)

type Utils struct {
	A       string
	B       string
	C       func(string)
	D       []byte
	Sleep   func(t uint64)
	Display func(otto.Value)
}

func loadJsLib(vm *otto.Otto, file string) error {
	dayjsBytes, e := os.ReadFile(file)
	if e != nil {
		return e
	}

	_, e = vm.Run(dayjsBytes)
	if e != nil {
		return e
	}
	return nil
}

var libs = []string{"dayjs.js"}

func main() {

	vm := otto.New()

	vm.Set("hello", func(a int64, b float64) float64 {
		fmt.Println("called from js: ", a, b)
		return float64(a) + b
	})

	utils := Utils{
		A: "hello",
		B: "world",
		C: func(a string) {
			fmt.Println("嘿嘿", a)
		},
		D: []byte("这是 bytes"),
		Sleep: func(t uint64) {
			time.Sleep(time.Duration(t) * time.Millisecond)
		},
		Display: func(v otto.Value) {
			z, e := v.Object().Get("z")
			if e != nil {
				fmt.Println("failed to get z", e)
				return
			}
			fmt.Println("Display", z)
			z.Call(otto.TrueValue(), 123)
		},
	}

	vm.Set("utils", utils)

	for _, lib := range libs {
		e := loadJsLib(vm, lib)
		if e != nil {
			fmt.Println("load "+lib+" failed", e)
		}
	}

	r, e := vm.Run(`
		console.log(Object.keys(this))

		var a = 123
		var b = 12.3
		var c = hello(a, b)
		console.log('hello js', c)
		Object.keys(utils).map(function(e){console.log('utils', e, utils[e])})

		console.log('bytes test', typeof utils.D, Array.isArray(utils.D), utils.D.length, typeof utils.D[0])

		console.log(dayjs('2022-12-12 12:34:56'))
		console.log(new Date())
		utils.Sleep(1000)
		console.log(new Date())

		utils.C('1235')

		utils.Display({
			x: 123,
			y: [4, 5, 6],
			z: function(a){console.log('Im js function: ' + a, 'this is ', this)}
		})

		console.log(Array.prototype.reduce)
	`)

	fmt.Println(r, e)

}
