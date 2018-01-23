package promise

import "sync"

type status uint

const (
	_PENDING status = iota
	_FULFILLED
	_REJECTED
)

type promise struct {
	value  interface{}
	reason error
	status
	onFulfilled func(interface{}) interface{}
	onReject    func(error) interface{}
	mutex       sync.Mutex
	next        *promise
}

func (p *promise) Then(onFulfilled func(interface{}) interface{}, onReject func(error) interface{}) *promise {

	genNext := func(result interface{}) *promise {
		switch v := result.(type) {
		case *promise:
			return v
		case error:
			return Promise.Reject(v)
		default:
			return Promise.Resolve(v)
		}
	}

	p.mutex.Lock()
	if p.status == _FULFILLED {
		p.mutex.Unlock()

		if onFulfilled == nil {
			return p
		}

		p.onFulfilled = onFulfilled
		ret := onFulfilled(p.value)

		p.next = genNext(ret)
		return p.next

	} else if p.status == _REJECTED {
		p.mutex.Unlock()

		if onReject == nil {
			return p
		}

		p.onReject = onReject
		ret := onReject(p.reason)

		p.next = genNext(ret)
		return p.next

	} else {
		// _PENDING
		if onFulfilled != nil {
			p.onFulfilled = onFulfilled
		}
		if onReject != nil {
			p.onReject = onReject
		}

		p.next = Promise(func(resolve func(interface{}), reject func(error)) {})

		p.mutex.Unlock()

		return p.next
	}
}
func (p *promise) Catch(onReject func(error) interface{}) *promise {
	return p.Then(nil, onReject)
}

func (p *promise) resolve(value interface{}) {
	p.mutex.Lock()
	p.value = value
	if p.status != _PENDING {
		p.mutex.Unlock()
		panic("Already resolved or rejected")
	}
	p.status = _FULFILLED
	p.mutex.Unlock()

	if p.onFulfilled != nil {
		ret := p.onFulfilled(p.value)
		if p2, ok := ret.(*promise); ok {
			p2.Then(func(v interface{}) interface{} {
				p.next.resolve(value)
				return p.next
			}, func(reason error) interface{} {
				p.next.reject(reason)
				return p.next
			})
		} else if p2, ok := ret.(error); ok {
			p.next.reject(p2)
		} else {
			p.next.resolve(ret)
		}
	} else if p.next != nil {
		p.next.resolve(p.value)
	}
}

func (p *promise) reject(reason error) {
	p.mutex.Lock()
	p.reason = reason
	if p.status != _PENDING {
		p.mutex.Unlock()
		panic("Already resolved or rejected")
	}
	p.status = _REJECTED
	p.mutex.Unlock()

	if p.onReject != nil {
		ret := p.onReject(p.reason)
		if p2, ok := ret.(*promise); ok {
			p2.Then(func(value interface{}) interface{} {
				p.next.resolve(value)
				return p.next
			}, func(reason error) interface{} {
				p.next.reject(reason)
				return p.next
			})
		} else if p2, ok := ret.(error); ok {
			p.next.reject(p2)
		} else {
			p.next.resolve(ret)
		}
	} else if p.next != nil {
		p.next.reject(p.reason)
	}
}

type executor func(func(interface{}), func(error))

type promiseWrapper func(executor) *promise

var Promise promiseWrapper

func (_ *promiseWrapper) Resolve(value interface{}) *promise {
	return Promise(func(resolve func(interface{}), reject func(error)) {
		resolve(value)
	})
}

func (_ *promiseWrapper) Reject(reason error) *promise {
	return Promise(func(resolve func(interface{}), reject func(error)) {
		reject(reason)
	})
}

func init() {
	Promise = func(executor executor) *promise {
		p := new(promise)

		resolve := func(value interface{}) {
			p.resolve(value)
		}

		reject := func(reason error) {
			p.reject(reason)
		}

		executor(resolve, reject)

		return p
	}
}
