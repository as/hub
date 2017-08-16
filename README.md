# Hub
Package hub provides a concurrent multiplexer on text editing operations
from multiple clients. To connect to a hub and make concurrent edits, see hub/client


# usage
This package is still in its experimental state. To see a demonstration of current functionality:

```
cd github.com/as/hub/example
example.exe localhost:8888 example.go &
cd ../client/example
example.exe localhost:8888 &
example.exe localhost:8888 &
example.exe localhost:8888 &
```
