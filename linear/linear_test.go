package linear

import (
	"testing"
    "os"
    "fmt"

	"github.com/stretchr/testify/assert"
)

// abs(x) = |x|
func abs(x float64) float64 {
	if x < 0 {
		return -1 * x
	}
	return x
}

var flatX [][]float64
var flatY []float64

var increasingX [][]float64
var increasingY []float64

var threeDLineX [][]float64
var threeDLineY []float64

func init() {
    
    // create the /tmp/.goml/ dir for persistance testing
    // if it doesn't already exist!
    err := os.MkdirAll("/tmp/.goml", os.ModePerm)
    if err != nil {
        panic(fmt.Sprintf("You should be able to create the directory for goml model persistance testing.\n\tError returned: %v\n", err.Error()))
    }
    
	// the line y=3
	flatX = [][]float64{}
	flatY = []float64{}
	for i := -10; i < 10; i++ {
		for j := -10; j < 10; j++ {
			for k := -10; k < 10; k++ {
				flatX = append(flatX, []float64{float64(i), float64(j), float64(k)})
				flatY = append(flatY, 3.0)
			}
		}
	}

	// the line y=x
	increasingX = [][]float64{}
	increasingY = []float64{}
	for i := -10; i < 10; i++ {
		increasingX = append(increasingX, []float64{float64(i)})
		increasingY = append(increasingY, float64(i))
	}

	threeDLineX = [][]float64{}
	threeDLineY = []float64{}
	// the line z = 10 + (x/10) + (y/5)
	for i := -10; i < 10; i++ {
		for j := -10; j < 10; j++ {
			threeDLineX = append(threeDLineX, []float64{float64(i), float64(j)})
			threeDLineY = append(threeDLineY, 10+float64(i)/10+float64(j)/5)
		}
	}
}

// test y=3
func TestFlatLineShouldPass1(t *testing.T) {
	var err error

	model := NewLeastSquares(.000001, 800, flatX, flatY)

	err = model.Learn()
	assert.Nil(t, err, "Learning error should be nil")

	var guess []float64

	for i := -20; i < 20; i += 10 {
		for j := -20; j < 20; j += 10 {
			for k := -20; k < 20; k += 10 {
				guess, err = model.Predict([]float64{float64(i), float64(j), float64(k)})
				assert.Len(t, guess, 1, "Length of a LeastSquares model output from the hypothesis should always be a 1 dimensional vector. Never multidimensional.")
				assert.InDelta(t, 3, guess[0], 1e-2, "Guess should be really close to 3 (within 1e-2) for y=3")
				assert.Nil(t, err, "Prediction error should be nil")
			}
		}
	}
}

// test y=3 but don't have enough iterations
func TestFlatLineShouldFail1(t *testing.T) {
	var err error

	model := NewLeastSquares(.000001, 1, flatX, flatY)

	err = model.Learn()
	assert.Nil(t, err, "Learning error should be nil")

	var guess []float64
	var faliures int

	for i := -20; i < 20; i += 10 {
		for j := -20; j < 20; j += 10 {
			for k := -20; k < 20; k += 10 {
				guess, err = model.Predict([]float64{float64(i), float64(j), float64(k)})
				assert.Len(t, guess, 1, "Length of a LeastSquares model output from the hypothesis should always be a 1 dimensional vector. Never multidimensional.")
				if abs(3.0-guess[0]) > 1e-2 {
					faliures++
				}
				assert.Nil(t, err, "Prediction error should be nil")
			}
		}
	}

	assert.True(t, faliures > 40, "There should be more faliures than half of the training set")
}

// test y=3 but alpha is too large
func TestFlatLineShouldFail2(t *testing.T) {
	var err error

	model := NewLeastSquares(1, 800, flatX, flatY)

	err = model.Learn()
	assert.NotNil(t, err, "Learning error should not be nil")
}

// test y=3 but include an invalid data set
func TestFlatLineShouldFail3(t *testing.T) {
	var err error

	model := NewLeastSquares(1, 800, [][]float64{}, flatY)

	err = model.Learn()
	assert.NotNil(t, err, "Learning error should not be nil")
    
    model = NewLeastSquares(1, 800, [][]float64{[]float64{}, []float64{}}, flatY)

	err = model.Learn()
	assert.NotNil(t, err, "Learning error should not be nil")

	model = NewLeastSquares(1, 800, nil, flatY)

	err = model.Learn()
	assert.NotNil(t, err, "Learning error should not be nil")
}

// test y=3 but include an invalid data set
func TestFlatLineShouldFail4(t *testing.T) {
	var err error

	model := NewLeastSquares(1, 800, flatX, []float64{})

	err = model.Learn()
	assert.NotNil(t, err, "Learning error should not be nil")

	model = NewLeastSquares(1, 800, flatX, nil)

	err = model.Learn()
	assert.NotNil(t, err, "Learning error should not be nil")
}

// test y=x
func TestInclinedLineShouldPass1(t *testing.T) {
	var err error

	model := NewLeastSquares(.0001, 500, increasingX, increasingY)
	err = model.Learn()
	assert.Nil(t, err, "Learning error should be nil")

	var guess []float64

	for i := -20; i < 20; i++ {
		guess, err = model.Predict([]float64{float64(i)})
		assert.Len(t, guess, 1, "Length of a LeastSquares model output from the hypothesis should always be a 1 dimensional vector. Never multidimensional.")
		assert.InDelta(t, i, guess[0], 1e-2, "Guess should be really close to input (within 1e-2) for y=x")
		assert.Nil(t, err, "Prediction error should be nil")
	}
}

// test z = 10 + (x/10) + (y/5)
func TestThreeDimensionalLineShouldPass1(t *testing.T) {
	var err error

	model := NewLeastSquares(.0001, 1000, threeDLineX, threeDLineY)
	err = model.Learn()
	assert.Nil(t, err, "Learning error should be nil")

	var guess []float64

	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			guess, err = model.Predict([]float64{float64(i), float64(j)})
			assert.Len(t, guess, 1, "Length of a LeastSquares model output from the hypothesis should always be a 1 dimensional vector. Never multidimensional.")
			assert.InDelta(t, 10.0+float64(i)/10+float64(j)/5, guess[0], 1e-2, "Guess should be really close to i+x (within 1e-2) for line z=10 + (x+y)/10")
			assert.Nil(t, err, "Prediction error should be nil")
		}
	}
}

// test persisting y=x to file
func TestPersistLeastSquaresShouldPass1(t *testing.T) {
	var err error

	model := NewLeastSquares(.0001, 500, increasingX, increasingY)
	err = model.Learn()
	assert.Nil(t, err, "Learning error should be nil")

	var guess []float64

	for i := -20; i < 20; i++ {
		guess, err = model.Predict([]float64{float64(i)})
		assert.Len(t, guess, 1, "Length of a LeastSquares model output from the hypothesis should always be a 1 dimensional vector. Never multidimensional.")
		assert.InDelta(t, i, guess[0], 1e-2, "Guess should be really close to input (within 1e-2) for y=x")
		assert.Nil(t, err, "Prediction error should be nil")
	}

	// not that we know it works, try persisting to file,
	// then resetting the parameter vector theta, then
	// restoring it and testing that predictions are correct
	// again.

	err = model.PersistToFile("/tmp/.goml/LeastSquares")
	assert.Nil(t, err, "Persistance error should be nil")

	model.Parameters = make([]float64, len(model.Parameters))

	// make sure it WONT work now that we reset theta
	//
	// the result of Theta transpose * X should always
	// be 0 because theta is the zero vector right now.
	for i := -20; i < 20; i++ {
		guess, err = model.Predict([]float64{float64(i)})
		assert.Len(t, guess, 1, "Length of a LeastSquares model output from the hypothesis should always be a 1 dimensional vector. Never multidimensional.")
		assert.Equal(t, 0.0, guess[0], "Guess should be 0 when theta is the zero vector")
		assert.Nil(t, err, "Prediction error should be nil")
	}

	err = model.RestoreFromFile("/tmp/.goml/LeastSquares")
	assert.Nil(t, err, "Persistance error should be nil")

	for i := -20; i < 20; i++ {
		guess, err = model.Predict([]float64{float64(i)})
		assert.Len(t, guess, 1, "Length of a LeastSquares model output from the hypothesis should always be a 1 dimensional vector. Never multidimensional.")
		assert.InDelta(t, i, guess[0], 1e-2, "Guess should be really close to input (within 1e-2) for y=x")
		assert.Nil(t, err, "Prediction error should be nil")
	}
}
