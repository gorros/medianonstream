package main

import (
	"fmt"
	"github.com/gorros/medianonstream"
)

func main() {
	a := []float64{1,2,5.6,6,10,12,35,45,2,5,7,4.6,0}
	mos := medianonstream.NewMedianOnStream(30)
	for i := 0; i < len(a); i++  {
		mos.Insert(a[i])
		fmt.Printf("Added %6.3f, current median is %3.3f\n", a[i], mos.GetMedian())
	}
}
