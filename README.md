# Went
is a dynamic scripting language targeted to be easily embeddable in host Go applications. It's overall design and syntax are inspired from Go itself, so there's almost no learning curve at all if you are already a Go developer.

This is a very early work in progress and a self-learning experience foremost, but pull requests are more than welcome.

## Status
List of things that are done and what needs to be.
- [x] Tokenization
- [x] Parsing
- [x] Syntax tree
- [x] Compilation to bytecode
- [ ] Register machine
- [ ] Go APIs
- [ ] Channels and goroutines
- [ ] Optimizations

## Syntax
The syntax is not well-defined yet, but what it's here should give you an idea of what'll be in the future.

Because it's heavily inspired in Go, you should almost feel no difference when switching between your compiled code and your script, minus the types.
```go
// Hello world in japanese (from google translator!)
println("こんにちは世界")

// factorial
func factorial(n) {
  return n <= 1 ? n : n * factorial(n - 1)
}
```

Note the addition of the ternary operator.
Also, this function could be one-lined with the short function syntax:
```go
func factorial(n) -> n <= 1 ? n : n * factorial(n - 1)
```

The `tests` folder contains some more examples, but they might be completely non-sense and change frequently since they're just tests. When the syntax is completely defined I'll create an `examples` folder with working code.

## Objects
Objects are the only hash-like structure in Went, they map strings (and only strings) to values. The literal syntax is almost identical to maps in Go:
```go
// map books names to their authors
authors := {
  "The Shining": "Stephen King",
  "Sherlock Holmes": "Sir Arthur Conan Doyle",
  "The Time Machine": "H.G. Wells",
}
```

To lookup a value from an object, or to set some key's value, you can use the following syntax: 
```go
timeMachineAuthor := authors["The Time Machine"]
authors["Harry Potter"] = "J.K. Rowling"
```

### Object inheritance
Besides being a hash map, objects are also useful for storing record-like data with meaningful properties, and actions (methods) that modify some state. Because of that, there is an alternative syntax to create and manipulate it's values:
```go
Book := {
  name: 'untitled',
  author: '',
  price: 0,
}

price := Book.price
Book.price = 57.50
```

With object inheritance, you can create an object from another object, they'll share the same properties, but each with it's own values. For that, there's the built-in function `inherit`:
```go
timeMachine := inherit(Book)
timeMachine.name = "The Time Machine"
timeMachine.author = "H.G. Wells"
timeMachine.price = 57.50
```

### Methods
Methods are just functions that take a receiver (the object it belongs to) as an implicit argument, through the special variable `this`:
```go
Book.comparePrice = func(otherBook) {
  return otherBook.price - this.price
}
```

These syntaxes are also valid:
```go
Book := {
  comparePrice: func(otherBook) {
    ...
  },
}

func Book.comparePrice(otherBook) {
  ...
}
```

Note that, unlike some other languages, methods are not functions bound to objects. Instead, the object is passed as a receiver to the function at call-time, and only when using the `dotted.syntax` so the following is invalid:
```go
comparePrice := someBook.comparePrice

// in both cases, 'this' will be nil
comparePrice(otherBook)
someBook["comparePrice"](otherBook)
```

If you need to call a method dynamically, there is also this syntax:
```go
// valid, 'this' will be set to 'someBook'
someBook.["comparePrice"](otherBook)
```

## Modules
Every .went file is a module, which is just a top-level function that can contain any statement, not just declarations as in Go.
To import a module, you can use the `import` built-in function:
```go
// import the 'fmt' module
fmt := import('fmt')
```

Now, the `fmt` variable is an object containing all the public functions and values from the `fmt` module.
"Public" is whatever the requested module returns, for example:
```go
// file: fib.went

// returns the [n]th number in the fibbonaci sequence
func nth(n) {
  if n <= 1 {
    return n
  }
  return nth(n - 1) + nth(n - 2)
}

// the import function will return this
return {
  nth: nth,
}
```

Now, in another module:
```go
fib := import('fib')
println(fib.nth(10)) //=> whatever number that is :)
```

## License
MIT
