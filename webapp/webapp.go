package main

import "github.com/gopherjs/gopherjs/js"

type Component interface {
	Render() string
}

func main() {
	idx := IndexComponent{}
	js.Global.Get("document").Call("write", idx.Render())

}
