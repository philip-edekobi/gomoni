package main

import (
	"test_proj/calc"
	"test_proj/calc/mul"
	"test_proj/out"
)

func main() {
	out.Prin(calc.Add(2, 3))
	out.Prin(mul.Mul(3, 4))
}
