# alog

ALog is the simplest form of the leveled logger based on the standard `log.Logger`.

## Usage

### New

Without creating an instance

```go
package main

import "github.com/orangenumber/alog"

func main() {
	alog.Print("hello1")     // hello1
	alog.Info("hello2")      // hello2
	alog.Infof("hello%d", 3) // hello3
}
```

With creating an instance

```go
package main

import "github.com/orangenumber/alog"

func main() {
	l := alog.New()
	l.Print(alog.C_INFO, "hello") // 17:54:13 hello
	l.Info("test") // 17:54:13 test
}
```

### Format 

```go
package main

import "github.com/orangenumber/alog"

func main() {
	// hello
	// hello my name is gon
	alog.Print("hello")
	alog.Printf("hello my name is %s", "gon")

	// sup
	// name is gon
	alog.SetPrefix("test1.")
	alog.Print("sup")
	alog.Printf("name is %s", "gon")

	// test2.hey
	// test2.call me gon
	alog.SetPrefix("test2.")
	alog.SetFormat(alog.F_PREFIX)
	alog.Print("hey")
	alog.Printf("call me %s", "gon")


	// 17:45:16.922084 test3.hey
	// 17:45:16.922195 test3.call me gon
	alog.SetPrefix("test3.")
	alog.SetFormat(alog.F_PREFIX, alog.F_MICROSECONDS)
	alog.Print("hey")
	alog.Printf("call me %s", "gon")
}
```


### Category Filtering

Many other loggers call it a `level`, but `alog` calls it a `category`. All the built-in categories are available
with a prefix `alog.C_`.  

```go
package main

import (
	"github.com/orangenumber/alog"
)

func main() {
	l := alog.New() // create instance
	l.SetFilter(alog.C_API, alog.C_BACK) // create filter
	l.Print(alog.C_API|alog.C_FRONT, "API+Front ex")
	l.Print(alog.C_WARN|alog.C_API, "Warn+API ex")
	l.Print(alog.C_WARN|alog.C_BACK, "Warn+Back ex")
	l.Printf(alog.C_WARN|alog.C_FRONT, "Warn+Front: [%s]", "warning code 123")
}
```

### Custom Category

Categories, unlike levels, support multiple, different types together. For instance, 
ErrSauce, ErrWing, ErrOrder are the 
error codes; BadCustom, GoodCustomer are categories; and 
HQ, Branch1, Branch2 are locations.
Create whichever you want with `uint64` just like the example below.

```go
package main

import "github.com/orangenumber/alog"

const (
	ErrSauce uint64 = 1 << iota
	ErrWing
	ErrOrder
	BadCustomer
	GoodCustomer
	HQ
	Branch1
	Branch2
)

func main() {
	l := alog.New()

	// ============================== EX 1: ALL
	// 17:59:57 order1
	// 17:59:57 order2
	// 17:59:57 err3
	l.Print(BadCustomer|HQ, "order1")
	l.Print(GoodCustomer|Branch1|ErrOrder, "order2")
	l.Print(HQ|Branch1|ErrWing, "err3")

	// ============================== EX 2: Only from HQ
	// 18:00:35 order1
	// 18:00:35 err3
	l.SetFilter(HQ)
	l.Print(BadCustomer|HQ, "order1")
	l.Print(GoodCustomer|Branch1|ErrOrder, "order2")
	l.Print(HQ|Branch1|ErrWing, "err3")

	// ============================== EX 3: All errors
	// 18:01:43 order2
	// 18:01:43 err3
	l.SetFilter(ErrWing, ErrOrder, ErrSauce) 
	l.Print(BadCustomer|HQ, "order1")
	l.Print(GoodCustomer|Branch1|ErrOrder, "order2")
	l.Print(HQ|Branch1|ErrWing, "err3")
}
```

