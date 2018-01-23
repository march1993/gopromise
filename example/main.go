package main

import (
	"errors"
	"fmt"
	. "github.com/march1993/gopromise"
	"math/rand"
	"time"
)

func rand1() bool {
	return rand.Float32() < 0.5
}

func main() {
	rand.Seed(time.Now().Unix())

	Promise(func(resolve func(interface{}), reject func(error)) {
		time.AfterFunc(time.Second, func() {
			resolve("<hello>")
		})
	}).Then(func(value interface{}) interface{} {
		fmt.Print(value.(string))
		return errors.New("<world>")
	}, nil).Catch(func(reason error) interface{} {
		fmt.Println(reason.Error())
		return Promise(func(resolve func(interface{}), reject func(error)) {

			if rand1() {
				time.AfterFunc(time.Second, func() {
					resolve("<nice>")
				})
			} else {
				time.AfterFunc(time.Second, func() {
					reject(errors.New("<good>"))
				})
			}
		})
	}).Then(func(value interface{}) interface{} {
		fmt.Println(value.(string))
		return nil
	}, nil).Catch(func(reason error) interface{} {
		fmt.Println(reason.Error())
		return nil
	})

	Promise.Resolve("resolved").Then(func(value interface{}) interface{} {
		fmt.Println(value.(string))
		return nil
	}, nil).Catch(func(reason error) interface{} {
		fmt.Println(reason.Error())
		return nil
	})

	Promise.Reject(errors.New("rejected")).Then(func(value interface{}) interface{} {
		fmt.Println(value.(string))
		return nil
	}, nil).Catch(func(reason error) interface{} {
		fmt.Println(reason.Error())
		return nil
	})

	time.Sleep(3 * time.Second)

}
