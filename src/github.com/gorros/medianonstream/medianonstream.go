package medianonstream


import (
	"math"
)

const DefBufferSize = 5

var nan = -1.0

type MedianOnStream struct {
	medianIndex         float64
	currentMedian       float64
	circularBufferStart int
	circularBuffer      []float64
	circularBufferSize  int
}

func NewMedianOnStream(size int) *MedianOnStream {
	mos := &MedianOnStream{}
	mos.medianIndex = nan
	if size > DefBufferSize {
		mos.circularBufferSize = size
	} else {
		mos.circularBufferSize = DefBufferSize
	}
	mos.circularBuffer = make([]float64, mos.circularBufferSize)
	for i := range mos.circularBuffer {
		mos.circularBuffer[i] = nan
	}

	return mos
}

func (mos *MedianOnStream) denormalizeInt(index int) int {
	if index < mos.circularBufferStart {
		return index + mos.circularBufferSize
	} else {
		return index
	}
}

func (mos *MedianOnStream) denormalizeFloat64(index float64) float64 {
	if index < float64(mos.circularBufferStart) {
		return index + float64(mos.circularBufferSize)
	} else {
		return index
	}
}

func (mos *MedianOnStream) normalizeInt(index int) int {
	return (index + mos.circularBufferSize) % mos.circularBufferSize
}

func (mos *MedianOnStream) normalizeFloat64(index float64) float64 {
	return math.Mod(index+float64(mos.circularBufferSize), float64(mos.circularBufferSize))
}

func (mos *MedianOnStream) getMiddleIndexNormalized() float64 {
	return math.Mod(float64(mos.circularBufferStart)+float64(mos.circularBufferSize)/2, float64(mos.circularBufferSize))
}

func (mos *MedianOnStream) getMaxIndex() int {
	return mos.normalizeInt(mos.circularBufferStart + mos.circularBufferSize - 1)
}

func (mos *MedianOnStream) shiftLeftAndInsert(index int, value float64) {
	for i := mos.circularBufferStart; i < mos.denormalizeInt(index); i++ {
		mos.circularBuffer[mos.normalizeInt(i)] = mos.circularBuffer[mos.normalizeInt(i+1)]
	}
	mos.circularBuffer[mos.normalizeInt(index)] = value
}

func (mos *MedianOnStream) shiftRightAndInsert(index int, value float64) {
	for i := mos.denormalizeInt(mos.getMaxIndex()); i > mos.denormalizeInt(index); i-- {
		mos.circularBuffer[mos.normalizeInt(i+1)] = mos.circularBuffer[mos.normalizeInt(i)]
	}
	mos.circularBuffer[mos.normalizeInt(index)] = value
}

func (mos *MedianOnStream) compare(a, b float64) int {
	if mos.denormalizeFloat64(a) > mos.denormalizeFloat64(b) {
		return 1
	} else if mos.denormalizeFloat64(a) < mos.denormalizeFloat64(b) {
		return -1
	}
	return 0
}

func (mos *MedianOnStream) GetMedian() float64 {
	ret := nan
	if mos.medianIndex != nan {
		ind := int(mos.medianIndex)
		if mos.medianIndex != float64(int64(mos.medianIndex)) {
			ret = (mos.circularBuffer[mos.normalizeInt(ind)] + mos.circularBuffer[mos.normalizeInt(ind+1)]) / 2
		} else {
			ret = mos.circularBuffer[mos.normalizeInt(ind)]
		}
	}
	return ret
}

func (mos *MedianOnStream) Insert(value float64) {
	currentMedian := mos.GetMedian()
	if currentMedian == nan { // first step
		mos.medianIndex = mos.getMiddleIndexNormalized()
		mos.circularBuffer[int(mos.medianIndex)] = value
	} else {
		if value > currentMedian { // go right
			// calculate the position of value
			for ind := int(mos.denormalizeFloat64(mos.medianIndex + 1.1)); ind < mos.circularBufferStart + mos.circularBufferSize; ind++ {
				if mos.circularBuffer[mos.normalizeInt(ind)] == nan ||
					value < mos.circularBuffer[mos.normalizeInt(ind)] {
					mos.shiftRightAndInsert(ind, value)
					break
				}
			}
			mos.medianIndex = mos.normalizeFloat64(mos.medianIndex + 0.5)
			if mos.compare(mos.medianIndex, mos.getMiddleIndexNormalized()) > 0 {
				mos.circularBufferStart = mos.normalizeInt(mos.circularBufferStart + 1)
				mos.circularBuffer[mos.circularBufferStart] = nan
			}
		} else { // go left
			// calculate the position of value
			for ind := int(mos.denormalizeFloat64(mos.medianIndex)); ind >= mos.circularBufferStart; ind-- {
				if mos.circularBuffer[mos.normalizeInt(ind)] == nan ||
					value > mos.circularBuffer[mos.normalizeInt(ind)] {
					mos.shiftLeftAndInsert(ind, value)
					break
				}
			}
			mos.medianIndex = mos.normalizeFloat64(mos.medianIndex - 0.5)
			if mos.compare(mos.medianIndex, mos.getMiddleIndexNormalized()) < 0 {
				mos.circularBufferStart = mos.normalizeInt(mos.circularBufferStart - 1)
				mos.circularBuffer[mos.getMaxIndex()] = nan
			}

		}
	}
}

func (mos *MedianOnStream) getBuffer() []float64 {
	buf := make([]float64, mos.circularBufferSize)
	for i := 0; i < mos.circularBufferSize; i++ {
		buf[i] = mos.circularBuffer[mos.normalizeInt(i+mos.circularBufferStart)]
	}
	return buf
}