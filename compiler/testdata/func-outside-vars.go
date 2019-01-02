package main

import "external"

func ex(fn func()) {
	fn()
}

func main() {

	outsideEscapes := 100
	stackval := 200

	ex(func () {
		external.Printf("escapes: %d\n", outsideEscapes)
		outsideEscapes = outsideEscapes + 1
	})

	ex(func () {
		external.Printf("escapes: %d\n", outsideEscapes)
		outsideEscapes = outsideEscapes + 1
	})

	ex(func () {
		anonval := 400
		external.Printf("anonval: %d\n", anonval)
	})

	external.Printf("stack: %d\n", stackval)
}
