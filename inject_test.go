package inject_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/lisitsky/inject"
)

type SpecialString interface {
}

type TestStruct struct {
	Dep1 string        `inject:"t" json:"-"`
	Dep2 SpecialString `inject`
	Dep3 string
	Dep4 string `inject:"" json:"-"`
}

type Greeter struct {
	Name string
}

func (g *Greeter) String() string {
	return "Hello, My name is " + g.Name
}

/* Test Helpers */
func expect(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Errorf("Expected %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

func refute(t *testing.T, a interface{}, b interface{}) {
	if a == b {
		t.Errorf("Did not expect %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

func Test_InjectorInvoke(t *testing.T) {
	injector := inject.New()
	expect(t, injector == nil, false)

	dep := "some dependency"
	injector.Map(dep)
	dep2 := "another dep"
	injector.MapTo(dep2, (*SpecialString)(nil))
	dep3 := make(chan *SpecialString)
	dep4 := make(chan *SpecialString)
	typRecv := reflect.ChanOf(reflect.RecvDir, reflect.TypeOf(dep3).Elem())
	typSend := reflect.ChanOf(reflect.SendDir, reflect.TypeOf(dep4).Elem())
	injector.Set(typRecv, reflect.ValueOf(dep3))
	injector.Set(typSend, reflect.ValueOf(dep4))

	_, err := injector.Invoke(func(d1 string, d2 SpecialString, d3 <-chan *SpecialString, d4 chan<- *SpecialString) {
		expect(t, d1, dep)
		expect(t, d2, dep2)
		expect(t, reflect.TypeOf(d3).Elem(), reflect.TypeOf(dep3).Elem())
		expect(t, reflect.TypeOf(d4).Elem(), reflect.TypeOf(dep4).Elem())
		expect(t, reflect.TypeOf(d3).ChanDir(), reflect.RecvDir)
		expect(t, reflect.TypeOf(d4).ChanDir(), reflect.SendDir)
	})

	expect(t, err, nil)
}

func Test_InjectorInvokeReturnValues(t *testing.T) {
	injector := inject.New()
	expect(t, injector == nil, false)

	dep := "some dependency"
	injector.Map(dep)
	dep2 := "another dep"
	injector.MapTo(dep2, (*SpecialString)(nil))

	result, err := injector.Invoke(func(d1 string, d2 SpecialString) string {
		expect(t, d1, dep)
		expect(t, d2, dep2)
		return "Hello world"
	})

	expect(t, result[0].String(), "Hello world")
	expect(t, err, nil)
}

func Test_InjectorInvokeInvalidValues(t *testing.T) {
	injector := inject.New()
	expect(t, injector == nil, false)

	dep := "some dependency"
	//injector.Map(dep)  -- dependency not provided
	result, err := injector.Invoke(func(d1 string) string {
		expect(t, d1, dep)
		return "Hi"
	})
	expect(t, len(result) == 0, true)
	expect(t, err == nil, false)
}

func Test_Injector_MapToInvalidTypeShouldPanic(t *testing.T) {
	injector := inject.New()
	expect(t, injector == nil, false)

	dep := "some dep"
	injector.Map(dep)
	dep2 := "another dep"
	// Expecting panic
	defer func() {
		rec := recover()
		refute(t, rec, nil)
	}()
	injector.MapTo(dep2, (*int)(nil))

	result, err := injector.Invoke(func(d1 string, d2 SpecialString) string {
		expect(t, d1, dep)
		expect(t, d2, dep2)
		return "Hello world"
	})

	expect(t, result[0].String(), "Hello world")
	expect(t, err, nil)
}

func Test_InjectorInvokeNotAFunction(t *testing.T) {
	injector := inject.New()
	defer func() {
		rec := recover()
		refute(t, rec == nil, true)
	}()
	_, _ = injector.Invoke(42)
}

func Test_InjectorApply(t *testing.T) {
	injector := inject.New()

	injector.Map("a dep").MapTo("another dep", (*SpecialString)(nil))

	s := TestStruct{}
	err := injector.Apply(&s)
	expect(t, err, nil)

	expect(t, s.Dep1, "a dep")
	expect(t, s.Dep2, "another dep")
	expect(t, s.Dep3, "")
	expect(t, s.Dep4, "a dep")
}

func Test_InjectorApplyAbsentValues(t *testing.T) {
	injector := inject.New()
	injector.Map(42)
	s := TestStruct{}
	err := injector.Apply(&s)
	expect(t, err != nil, true)
}

func Test_InjectorApplyNotStruct(t *testing.T) {
	injector := inject.New()
	injector.Map(42)
	i := 42
	err := injector.Apply(&i)
	expect(t, err == nil, true)
}

func Test_InterfaceOf(t *testing.T) {
	iType := inject.InterfaceOf((*SpecialString)(nil))
	expect(t, iType.Kind(), reflect.Interface)

	iType = inject.InterfaceOf((**SpecialString)(nil))
	expect(t, iType.Kind(), reflect.Interface)

	// Expecting nil
	defer func() {
		rec := recover()
		refute(t, rec, nil)
	}()
	iType = inject.InterfaceOf((*testing.T)(nil))
}

func Test_InjectorSet(t *testing.T) {
	injector := inject.New()
	typ := reflect.TypeOf("string")
	typSend := reflect.ChanOf(reflect.SendDir, typ)
	typRecv := reflect.ChanOf(reflect.RecvDir, typ)

	// instantiating unidirectional channels is not possible using reflect
	// http://golang.org/src/pkg/reflect/value.go?s=60463:60504#L2064
	chanRecv := reflect.MakeChan(reflect.ChanOf(reflect.BothDir, typ), 0)
	chanSend := reflect.MakeChan(reflect.ChanOf(reflect.BothDir, typ), 0)

	injector.Set(typSend, chanSend)
	injector.Set(typRecv, chanRecv)

	expect(t, injector.Get(typSend).IsValid(), true)
	expect(t, injector.Get(typRecv).IsValid(), true)
	expect(t, injector.Get(chanSend.Type()).IsValid(), false)
}

func Test_InjectorGet(t *testing.T) {
	injector := inject.New()

	injector.Map("some dependency")

	expect(t, injector.Get(reflect.TypeOf("string")).IsValid(), true)
	expect(t, injector.Get(reflect.TypeOf(11)).IsValid(), false)
}

func Test_GetWithInvalidInvokeShouldPanic(t *testing.T) {
	injector := inject.New()
	expect(t, injector == nil, false)

	provider := func(s string) int {
		return 42
	}
	var i int
	injector.Provide(provider)

	// dependency with type `string` is not provided - expecting panic
	defer func() {
		rec := recover()
		refute(t, rec, nil)
		e, ok := rec.(error)
		refute(t, ok, false)
		expect(t, e.Error() == "Value not found for type string", true)
	}()

	result := injector.Get(reflect.TypeOf(i))
	expect(t, result.IsValid(), false)
}

func Test_InjectorSetParent(t *testing.T) {
	injector := inject.New()
	injector.MapTo("another dep", (*SpecialString)(nil))

	injector2 := inject.New()
	injector2.SetParent(injector)

	expect(t, injector2.Get(inject.InterfaceOf((*SpecialString)(nil))).IsValid(), true)
}

func TestInjectImplementors(t *testing.T) {
	injector := inject.New()
	g := &Greeter{"Jeremy"}
	injector.Map(g)

	expect(t, injector.Get(inject.InterfaceOf((*fmt.Stringer)(nil))).IsValid(), true)
}

func Test_InjectorProvideStruct(t *testing.T) {
	injector := inject.New()

	expect(t, injector.Get(reflect.TypeOf(&TestStruct{})).IsValid(), false)

	injector.Provide(func() *TestStruct {
		return &TestStruct{
			Dep3: "test",
		}
	})

	injectedStruct := injector.Get(reflect.TypeOf(&TestStruct{}))
	expect(t, injectedStruct.IsValid(), true)
	if injectedStruct.IsValid() {
		expect(t, injectedStruct.Interface().(*TestStruct).Dep3, "test")
	}

	_, err := injector.Invoke(func(s1 *TestStruct) {
		expect(t, s1.Dep3, "test")
	})
	expect(t, err, nil)
}

func Test_InjectorProvideInterface(t *testing.T) {
	injector := inject.New()

	expect(t, injector.Get(inject.InterfaceOf((*fmt.Stringer)(nil))).IsValid(), false)

	injector.Provide(func() fmt.Stringer {
		return &Greeter{"Jeremy"}
	})

	expect(t, injector.Get(inject.InterfaceOf((*fmt.Stringer)(nil))).IsValid(), true)

	_, err := injector.Invoke(func(stringer fmt.Stringer) {
		expect(t, stringer.String(), "Hello, My name is Jeremy")
	})
	expect(t, err, nil)
}

func Test_InjectorConstruct(t *testing.T) {
	injector := inject.New()
	injector.Map("a dep").MapTo("another dep", (*SpecialString)(nil))

	type Constructable struct {
		S  string
		SS SpecialString
	}

	FConstr := func(s string, ss SpecialString) Constructable {
		return Constructable{S: s, SS: ss}
	}
	injector.Provide(FConstr)

	type IConstructable interface{}
	FIConstr := func(s string, ss SpecialString) IConstructable {
		return IConstructable(&Constructable{S: s, SS: ss})
	}

	injector.Provide(FIConstr)

	var ic IConstructable
	var err error
	// good interface
	err = injector.Construct(&ic)
	expect(t, err, nil)
	icc := ic.(*Constructable)
	expect(t, icc.S == "a dep", true)
	expect(t, icc.SS == "another dep", true)

	// good structure
	var c Constructable
	err = injector.Construct(&c)
	expect(t, err, nil)
	expect(t, c.S == "a dep", true)
	expect(t, c.SS == "another dep", true)

	// bad - no constructor for good structure
	var c2 *Constructable
	err = injector.Construct(&c2)
	refute(t, err, nil)
	refute(t, c2, nil)

	// bad not structure not ptr
	var bad int
	err = injector.Construct(bad)
	refute(t, err, nil)
}

func TestConstructLaterStruct(t *testing.T) {
	injector := inject.New()
	injector.Map("a dep").MapTo("another dep", (*SpecialString)(nil))

	type Constructable struct {
		S  string
		SS SpecialString
	}
	FConstr := func(s string, ss SpecialString) Constructable {
		return Constructable{S: s, SS: ss}
	}
	injector.Provide(FConstr)

	var constructable Constructable
	var err error
	err = injector.ConstructLater(&constructable)
	expect(t, err, nil)
	// object is not ready at the moment
	expect(t, constructable.S == "", true)
	err = injector.FinishConstruct()
	// no error on finalizing
	expect(t, err, nil)
	expect(t, constructable.S, "a dep")
	expect(t, constructable.SS, "another dep")
}

func TestConstructLaterStructWithAbsentDependency(t *testing.T) {
	injector := inject.New()
	injector.Map("a dep").MapTo("another dep", (*SpecialString)(nil))

	type Constructable struct {
		S  string
		SS SpecialString
	}
	//FConstr := func(s string, ss SpecialString) Constructable {
	//	return Constructable{S: s, SS: ss}
	//}
	//injector.Provide(FConstr)

	var constructable Constructable
	var err error
	err = injector.ConstructLater(&constructable)
	expect(t, err, nil)

	// object is not ready at the moment
	expect(t, constructable.S == "", true)
	err = injector.FinishConstruct()

	// error on finalizing, fields should not get their values
	refute(t, err, nil)
	expect(t, constructable.S == "a dep", false)
	expect(t, constructable.SS == "another dep", false)
}

func TestConstructLaterInterface(t *testing.T) {
	injector := inject.New()
	injector.Map("a dep").MapTo("another dep", (*SpecialString)(nil))

	type Constructable struct {
		S  string
		SS SpecialString
	}

	type IConstructable interface{}
	FIConstr := func(s string, ss SpecialString) IConstructable {
		return IConstructable(&Constructable{S: s, SS: ss})
	}

	injector.Provide(FIConstr)

	var iconstructable IConstructable
	var err error
	err = injector.ConstructLater(&iconstructable)
	expect(t, err, nil)

	// object is not ready at the moment
	icc, ok := iconstructable.(*Constructable)
	expect(t, ok, false)
	expect(t, icc == nil, true)
	err = injector.FinishConstruct()

	// no error on finalizing
	expect(t, err, nil)

	// object is ok with dependencies
	icc, ok = iconstructable.(*Constructable)
	expect(t, ok, true)
	refute(t, icc == nil, true)
	expect(t, icc.S, "a dep")
	expect(t, icc.SS, "another dep")
}
