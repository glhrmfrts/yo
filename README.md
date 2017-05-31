# ヨ (Yo)
is a dynamic scripting language targeted to be easily embeddable in host Go applications. It's overall design and syntax are inspired from Go itself, so there's almost no learning curve at all if you are already a Go developer.

This is a very early work in progress and a self-learning experience foremost, but pull requests are more than welcome.

## Status
List of things that are done and what needs to be.
- [x] Tokenization
- [x] Parsing
- [x] Syntax tree
- [x] Compilation to bytecode
- [x] Register machine (WIP)
- [ ] Go APIs
- [ ] Channels and goroutines
- [ ] Optimizations

## Syntax
Because it's heavily inspired in Go, you should almost feel no difference when switching between your compiled code and your script, minus the types.
```go
// Hello world, from Yo!
println("こんにちは世界、ヨから")

// functions
func add(x, y) {
  return x + y
}
println(add(2, 5))

// compile-time constants
const numExamples = 7
println(numExamples + 3) // compiles "println(10)"

// error handling, multiple return values, short variable declaration (:=)
content, err := ioutil.readFile("data.txt")
if err {
  log.fatal(err)
}

// arrays
arr := [1, 2, 3]
append(arr, 4, 5, 6)
println(len(arr), " ", arr)

// objects/maps
Vector2 := {
  x: 0,
  y: 0
}

Vector3 := new(Vector2, {z: 0})

positions := {
  "user1": new(Vector3),
  "user2": new(Vector3, {x: 525.4, y: 320, z: 110.54}),
}

// methods
func Vector3.mul(multiple) {
  return new(Vector3, {
    x: this.x * multiple,
    y: this.y * multiple,
    z: this.z * multiple,
  })
}

// closures
func seq(start) {
  i := 0
  return func() {
    return start + i++
  }
}

inc := seq(1)
println(inc()) // 1
println(inc()) // 2
println(inc()) // 3
```

## License
MIT
