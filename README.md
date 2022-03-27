# Dualerr

## Why Dualerr?

Since go1.13, go supports error wrapping. It looks great at first glance.
But as I started to use the wrapping, it became painful.

When developing backend server, I expect two things on logs.

1. Trace the stack that an error occurred.
2. Get an error message for client response.

A vanilla go error is quite goot for the requirement number one. See the below
example.

```go
func one() error {
  err := fmt.Errorf("complex error from database")
  return fmt.Errorf("one: %w", err)
}

func two() error {
  return fmt.Errorf("two: %w", one())
}

func three() error {
  return fmt.Errorf("three: %w", two())
}

func main() {
  err := three()
  fmt.Println(err)
  // Output: three: two: one: complex error from database
}
```

By watching the error message, we can figure out a stack trace of the error
chain. It seems good enough, but keep in mind that we are developing a backend
server. Should we send the `err` as client response? Of course not, the message
of `err` is too detailed for a client error message.

See the solution by using `Dualerr`.

```go
func one() error {
  err := fmt.Errorf("complex error from database")
  return dualerr.New(err, "simple message for client[%s]", "you")
}

func two() error {
  return dualerr.Wrap(one())
}

func three() error {
  return dualerr.Wrap(two())
}

func main() {
  err := three()

  fmt.Println(err)
  // Output: main.three: main.two: main.one: complex error from database

  var dual dualerr.Error
  errors.As(err, &dual)

  fmt.Println(dual)
  // Output: main.three: main.two: main.one: complex error from database

  fmt.Println(dual.Simple())
  // Output: simple message for client[you]
}
```

You still can use `err` for detailed logging, but also you can extract simple
error message using `dual.Simple()` just for your client error response.

Also note that function names are automatically included into the error chain.
