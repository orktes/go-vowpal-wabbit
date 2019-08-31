[![Build Status](https://travis-ci.org/orktes/go-vowpal-wabbit.svg?branch=master)](https://travis-ci.org/orktes/go-vowpal-wabbit)
[![GoDoc](https://godoc.org/github.com/orktes/go-vowpal-wabbit?status.svg)](http://godoc.org/github.com/orktes/go-vowpal-wabbit)


# go-vowpal-wabbit
Vowpal Wabbit bindings for Golang

```go
import "github.com/orktes/go-vowpal-wabbit"
```

## Usage

Library depends on the Vowpal Wabbit shared library. To install the shared lib please follow [VW installation instruction](https://github.com/VowpalWabbit/vowpal_wabbit/wiki/Getting-started). 

API is designed to closely resemble the C API of Vowpal Wabbit with minor changes done to make it more convinient to be used from go. Library is not thread safe and additional locking is required if libary is being called from multiple goroutines.

## Minimalistic example

```go
import (
    wabbit "github.com/orktes/go-vowpal-wabbit"
)

func main() {
	vw, _ := wabbit.New("-q st --noconstant --quiet")
	defer vw.Finish()

	// Learn an example
	ex, _ := vw.ReadExample("1 |s p^the_man w^the w^man |t p^un_homme w^un w^homme")
	vw.Learn(ex)

	// Predict with features
	ex, _ = vw.ReadExample("|s p^the_man w^the w^man |t p^un_homme w^un w^homme")
	res := vw.Predict(ex)

	fmt.Printf("Prediction: %f", res)
	// Output: Prediction: 0.855723
}

```

## TODO

- [ ] Handle exceptions from VW
- [ ] Add pooling for example creation (json & dsjson)

