This is a very stupid Promise implementation in golang.

## Installation
    $ go get github.com/march1993/gopromise

## Usage
Only one function called `Promise` is exported.
```go
package somepackage
import . "github.com/march1993/gopromise"
func () {

	// create a promise
	Promise(func(resolve func(interface{}), reject func(error)) {
		// call resolve or reject
	}).Then(func(value interface{}) interface{} {
		return something
	}, func(reason error) interface{} {
		return something
	}).Catch(func(reason error) interface{} {
		return something
	})

	// create a resolved promise
	Promise.Resolve(something).Then().Catch()

	// create a rejected promise
	Promise.Rejected(something).Then().Catch()

}
```

## Example
See [here](example/main.go).