# Median on stream
Calculates median on stream with limited memory.

Default buffer size is 1000. It should be at least twice large than 
standard deviation of the sequence. 

## Example
```golang
	a := []float64{1,2,5.6,6,10,12,35,45,2,5,7,4.6,0}
	mos := medianonstream.NewMedianOnStream(30)
	for i := 0; i < len(a); i++  {
		mos.Insert(a[i])
		fmt.Printf("Added %6.3f, current median is %3.3f\n", a[i], mos.GetMedian())
	}
```

Output
```
Added  1.000, current median is 1.000
Added  2.000, current median is 1.500
Added  5.600, current median is 2.000
Added  6.000, current median is 3.800
Added 10.000, current median is 5.600
Added 12.000, current median is 5.800
Added 35.000, current median is 6.000
Added 45.000, current median is 8.000
Added  2.000, current median is 6.000
Added  5.000, current median is 5.800
Added  7.000, current median is 6.000
Added  4.600, current median is 5.800
Added  0.000, current median is 5.600

```


###Notice
The original version was written in Java by @shunanya