# execgroup

ExecGroup uses a WaitGroup to allow multiple functions to be called concurrently.

````
var eg ExecGroup

// Run two functions concurrently.
eg.Do(func() {
  // Function 1
})

eg.Do(func() {
  // Function 2
})

// Wait for both functions to finish and handle errors.
if err := eg.Wait(); err != nil {
    // handle errors.
}
````