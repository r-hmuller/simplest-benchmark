This is a benchmark to stress an URL.

There are 4 parameters that need to be passed, and they are
ordered:

- Number of threads to execute requests
- Number of requests _per thread_
- URL that will receive the requests
- Path to the file where the duration will be stored

Some requests will have a counter to check how long
it took to execute. Those times will be stored on the
file provided in the 4th argument.