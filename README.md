# Elo
is a dynamic scripting language targeted to be easily embeddable in host Go applications. It's overall design and syntax are inspired from Go itself, so there's almost no learning curve at all if you are already a Go developer.

## Status
List of things that are done and what needs to be.
- [x] Tokenization
- [x] Parsing
- [x] Syntax tree
- [x] Compilation to bytecode (Partially done, some statements are still not being compiled)
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
func factorial(n) => n <= 1 ? n : n * factorial(n - 1)
```

The `tests` folder contains some more examples, but they might be completely non-sense and change frequently since they're just tests. When the syntax is completely defined I'll create an `examples` folder with working code.

## Modules
Every .elo file is a module, which is just a top-level function that can contain any statement, not just declarations as in Go.
To import a module, one can use the `import` built-in function:
```go
// import the 'fmt' module
fmt := import('fmt')
```

Now, the `fmt` variable is an object containing all the public functions and values from the `fmt` module.
The `import` function returns whatever the requested module returns, for example:
```go
// file: fib.elo

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
