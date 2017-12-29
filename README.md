# inject - Dependency Injection and Lazy Construction in Go

[![Build Status](https://travis-ci.org/lisitsky/inject.svg?branch=master)](https://travis-ci.org/lisitsky/inject)
[![Coverage Status](https://coveralls.io/repos/github/lisitsky/inject/badge.svg?branch=docs)](https://coveralls.io/github/lisitsky/inject?branch=docs)

    import "github.com/lisitsky/inject"

Package `inject` provides utilities for total application 
components decoupling and lazy object initialization with automatic 
and manual dependencies injection resolving.

## Key advantages

 * Constructor arguments could be resolved automatically
 * Component can have standard public constructor for manuall calls 
 and private/public constructor for automatical injection
 * Total decoupling as between dependencies and dependency acceptors
 so as different interfaces implementation 
 * Main application has to "know" only **what** does it want to obtain, 
 but not **how** and **from where** it should be constructed.
 * No direct constructor invocations at all. Any package providing dependency 
 can be substituted at any moment to another one providing the same interface .
 * Easy **reassignment** of constructors **at late stages**. This helps to provide 
 one or two mocks and makes tests very easy without changes in another application pars.
       


## Simple example:
Suppose we have a package providing struct matching `fmt.Stringer` interface .
All we have to do is to add a constructor for it (here `NewStringer` function) 
then announce it in `init` function with `inject.Provide()` method. 

This declares: "Function `NewStringer` provides a way to construct `fmt.Stringer`. 
And anybody who can use it without further declarations." 


```go
package dependency
    
import (
   "fmt"
   
   "github.com/lisitsky/inject"
)
    
func init() {
   inject.Provide(NewStringer)
}
    
type stringer struct{}
    
func (s stringer) String() string {
    return "Hello, World"
}
    
func NewStringer() fmt.Stringer {
    return stringer{}
}
```

On a side accepting dependency (`main.go`):

```go
package main
 

import (
   "fmt"
	
   "github.com/lisitsky/inject"
   
   _ "github.com/lisitsky/inject/examples/simple/dependency"
)
 
var ( 
   str fmt.Stringer
)
 
func main() {
   inject.Construct(&str)
   fmt.Println("My Stringer is:", str)
}
```

Output:

    My Stringer is: Hello, World
    


Language Translations:
* [简体中文](translations/README_zh_cn.md)

## Usage

#### func  InterfaceOf

```go
func InterfaceOf(value interface{}) reflect.Type
```
InterfaceOf dereferences a pointer to an Interface type. It panics if value is
not an pointer to an interface.

#### type Applicator

```go
type Applicator interface {
	// Maps dependencies in the Type map to each field in the struct
	// that is tagged with 'inject'. Returns an error if the injection
	// fails.
	Apply(interface{}) error
}
```

Applicator represents an interface for mapping dependencies to a struct.

#### type Injector

```go
type Injector interface {
	Applicator
	Invoker
	TypeMapper
	// SetParent sets the parent of the injector. If the injector cannot find a
	// dependency in its Type map it will check its parent before returning an
	// error.
	SetParent(Injector)
}
```

Injector represents an interface for mapping and injecting dependencies into
structs and function arguments.

#### func  New

```go
func New() Injector
```
New returns a new Injector.

#### type Invoker

```go
type Invoker interface {
	// Invoke attempts to call the interface{} provided as a function,
	// providing dependencies for function arguments based on Type. Returns
	// a slice of reflect.Value representing the returned values of the function.
	// Returns an error if the injection fails.
	Invoke(interface{}) ([]reflect.Value, error)
}
```

Invoker represents an interface for calling functions via reflection.

#### type TypeMapper

```go
type TypeMapper interface {
	// Maps the interface{} value based on its immediate type from reflect.TypeOf.
	Map(interface{}) TypeMapper
	// Maps the interface{} value based on the pointer of an Interface provided.
	// This is really only useful for mapping a value as an interface, as interfaces
	// cannot at this time be referenced directly without a pointer.
	MapTo(interface{}, interface{}) TypeMapper
	// Provide the dynamic type of interface{} returns,
	Provide(interface{}) TypeMapper
	// Provides a possibility to directly insert a mapping based on type and value.
	// This makes it possible to directly map type arguments not possible to instantiate
	// with reflect like unidirectional channels.
	Set(reflect.Type, reflect.Value) TypeMapper
	// Returns the Value that is mapped to the current type. Returns a zeroed Value if
	// the Type has not been mapped.
	Get(reflect.Type) reflect.Value
}
```

TypeMapper represents an interface for mapping interface{} values based on type.
