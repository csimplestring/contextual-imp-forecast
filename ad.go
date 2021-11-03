package main

type Ad interface {
	Features() []Feature
	Weight(Feature) float32
	Bid() float32
}
