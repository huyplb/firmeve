package firmeve

import (
	"github.com/firmeve/firmeve/utils"
	"github.com/kataras/iris/core/errors"
	"reflect"
	"strings"
	"sync"
)

var (
	firmeve *Firmeve
	once    sync.Once
	mutex   sync.Mutex
)

type prototypeFunc func(container Container, params ...interface{}) interface{}

type Container interface {
	Get(id string) interface{}
	Has(id string) bool
	//Bind(id interface{})
}

type binding struct {
	share     bool
	instance  interface{}
	prototype prototypeFunc
}

type Firmeve struct {
	Container
	bindings map[string]*binding
	//bindingOptions
	//resolveOptions
}

// Create a new firmeve container
func NewFirmeve() *Firmeve {
	if firmeve != nil {
		return firmeve
	}

	once.Do(func() {
		firmeve = &Firmeve{
			bindings: make(map[string]*binding),
		}
	})

	return firmeve
}

// 0. 每一个使用的都必须要注册，但可以不是单例，在init函数中是个好的选择
// 0.0.0 如果是非singleton的，必须是func，类型，这样才能得到每次的原型
// 0.0 注册时候的路径包是个问题，没有路径包或名称就没法在resolve的时候进行取值
// 0.1 结构体可以通过反射取名称，但如果是func肯定是需要自己手动加入name
// 1. 全部存储函数类型
// 2. 如果是singleton保存到instance实例
// 3. 先实现这样的bind和resolve，后面增加struct tag 注入
// 4. 名称怎样实现惟一呢？，完整路径太长(github.com/b/c) 不是完整路径可能会存在冲突

type bindingOption struct {
	name      string
	share     bool
	prototype interface{}
}

func (f *Firmeve) Resolve(abstract interface{}, params ...interface{}) interface{} {
	reflectType := reflect.TypeOf(abstract)
	kind := reflectType.Kind()
	//var value interface{}

	switch kind {
	case reflect.String:
		if v, ok := f.bindings[abstract.(string)]; ok {
			return v
		} else {
			panic(`不存在`)
		}
	case reflect.Func:
		// 反射参数
		//params := reflectType.NumIn()
		newParams := make([]reflect.Value,0)
		for i := 0; i < reflectType.NumIn(); i++ {
			reflectSubType := reflectType.In(0)
			name, err := f.parseName(reflectSubType, "")
			if err != nil {
				panic(err)
			}

			result := f.bindings[name].prototype(f)
			if reflectSubType.Kind() == reflect.Func {
				// 参数暂时为空
				result = reflect.ValueOf(result).Call([]reflect.Value{})
			} else {

			}

			newParams = append(newParams,reflect.ValueOf(result))
		}

		return reflect.ValueOf(abstract).Call(newParams)[0].Interface()

	default:
		panic(`暂不支持`)
	}

	return nil
}

func WithBindShare(share bool) utils.OptionFunc {
	return func(option interface{}) {
		option.(*bindingOption).share = share
	}
}

func WithBindName(name string) utils.OptionFunc {
	return func(option interface{}) {
		option.(*bindingOption).name = strings.ToLower(name)
	}
}

func WithBindInterface(object interface{}) utils.OptionFunc {
	return func(option interface{}) {
		option.(*bindingOption).prototype = object
	}
}

func (f *Firmeve) Bind(options ...utils.OptionFunc) { //, value interface{}
	// default bindingOption
	bindingOption := newBindingOption()

	// 参数解析
	utils.ApplyOption(bindingOption, options...)

	// 反射对象类型
	reflectType := reflect.TypeOf(bindingOption.prototype)

	name, err := f.parseName(reflectType, bindingOption.name)
	if err != nil {
		panic(err)
	}

	mutex.Lock()
	defer mutex.Unlock()
	f.bindings[name] = newBinding(f.setPrototype(bindingOption.prototype), bindingOption.share)
	//kind := reflectType.Kind()
	//
	//if kind == reflect.Ptr {
	//
	//}
	//
	//// 处理对象
	//object := bindingOption.object
	//
	//binding := newBinding()
	//binding.shared = bindingOption.share
	//
	//kind := reflectType.Kind()
	//
	//if kind == reflect.Func {
	//	if bindingOption.name == "" {
	//		panic("函数类型，Name必须存在")
	//	}
	//
	//	// 并且函数必须要有返回值
	//	//binding.
	//	//if reflectType.NumOut() == 0 {
	//	//	panic("函数类型，必须要有返回值")
	//	//}
	//
	//	binding.prototype = func(app Container, params ...interface{}) interface{} {
	//		return object
	//	}
	//
	//	f.bindings[bindingOption.name] = binding
	//
	//	func3 := f.bindings[bindingOption.name].prototype(f)
	//	result := reflect.ValueOf(func3).Call([]reflect.Value{})
	//	fmt.Printf("%#v", result[1])
	//} else if kind == reflect.Ptr {
	//	fmt.Println("============")
	//	// get Name
	//
	//	name := strings.Join([]string{reflectType.Elem().PkgPath(), reflectType.Elem().Name()}, `.`)
	//	fmt.Println(name)
	//	binding.prototype = func(container Container, params ...interface{}) interface{} {
	//		return object
	//	}
	//
	//	f.bindings[name] = binding
	//
	//	func3 := f.bindings[name].prototype(f)
	//	fmt.Printf("%#v", func3)
	//} else if kind == reflect.Struct {
	//	name := strings.Join([]string{reflectType.PkgPath(), reflectType.Name()}, `.`)
	//	fmt.Println(name, "==struct==")
	//	binding.prototype = func(container Container, params ...interface{}) interface{} {
	//		return object
	//	}
	//
	//	f.bindings[name] = binding
	//	func3 := f.bindings[name].prototype(f)
	//	fmt.Printf("%#v", func3)
	//} else if kind == reflect.Slice {
	//	if bindingOption.name == "" {
	//		panic("函数类型，Name必须存在")
	//	}
	//	binding.prototype = func(app Container, params ...interface{}) interface{} {
	//		return object
	//	}
	//	func3 := binding.prototype(f).([]string)
	//	func3 = append(func3, "d")
	//	fmt.Printf("%#v", func3)
	//}

	//fmt.Println(f.bindings[bindingOption.name].prototype(f.(*Container)))
	//fmt.Printf("%#v",f.bindings)
	//fmt.Println(reflect.TypeOf(object).Kind().String())

	//if bindDefaultOption.name == `` {
	//
	//}

	//for _, opt := range options {
	//	opt(bindDefaultOption)
	//}

	//fmt.Printf("%#v", bindDefaultOption)

	//// 取得name和对象
	//paramsLen := len(params)
	//if paramsLen == 1 { //name需要自己反射，必须是结构体,share:false
	//
	//} else if paramsLen == 2 {
	//	if isShare, ok := params[1].(bool); ok && isShare { // name需要自己反射，单例类型
	//
	//	} else {
	//
	//	}
	//}
	//
	//fmt.Println(path.Base("a/b/c"))
	////f.bindings[id] = newBinding(false, value)
	//reflectType := reflect.TypeOf(object)
	//fmt.Printf("%p", reflectType)
	//reflectType.Elem()
	//println(reflectType.Kind().String())
	//if reflectType.Kind() == reflect.Ptr {
	//	fullName := strings.Join([]string{reflectType.Elem().PkgPath(), reflectType.Elem().Name()}, ".")
	//	f.bindings[fullName] = newBinding(share, object)
	//}

	//var method func(f *Firmeve) interface{}
	//
	//if reflectType.Kind() != reflect.Func {
	//	method = func(f *Firmeve) interface{} {
	//		return object
	//	}
	//} else {
	//	method = func(f *Firmeve) interface{} {
	//		reflect.ValueOf(object).Elem().Call()
	//	}
	//}

	//if share {
	//
	//} else {
	//
	//}
}

func (f *Firmeve) parseName(reflectType reflect.Type, name string) (string, error) {

	if name == "" {
		kind := reflectType.Kind()
		switch kind {
		case reflect.Ptr:
			name = strings.Join([]string{reflectType.Elem().PkgPath(), reflectType.Elem().Name()}, `.`)
		case reflect.Struct:
			name = strings.Join([]string{reflectType.PkgPath(), reflectType.Name()}, `.`)
		default:
			return ``, errors.New(`不支持的类型`)
		}
	}

	return strings.ToLower(name), nil
}

func (f *Firmeve) setPrototype(prototype interface{}) prototypeFunc {
	return func(container Container, params ...interface{}) interface{} {
		return prototype
	}
}

func (f *Firmeve) Get(id string) interface{} {
	//panic("implement me")
	return "abc"
}

func (f *Firmeve) Has(id string) bool {
	//panic("implement me")
	return false

}

// ---------------------------- bindingOption ------------------------

func newBindingOption() *bindingOption {
	return &bindingOption{share: false,}
}

// ---------------------------- binding ------------------------

func newBinding(prototype prototypeFunc, share bool) *binding {
	return &binding{
		prototype: prototype,
		share:     share,
	}
}
