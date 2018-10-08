# Store ~ Local Object Storage made Simple

Store is a simple library to help store groups of objects in files.
	 
The following code snippets are how to get started using store quickly.

```go
import "github.com/rustyeddy/store"

func () {
  // Use the specified file path for storage.
  st, err  := UseStore("/tmp/somewhere")

  // Example 

  
  // StoreObject will take the Go object record, serialize it to JSON
  // (or other format if configured) then write data to disk.  NOTE
  // just like json.Unmarshal(because of json.Unmarshal) record has to
  // be the pointer to a structure.  That structure will be fill in.
  var record *ImportRecord
  obj, err := st.StoreObject("important", record)
  
  // FetchObject does the opposite of StoreObject. FetchObject 
  // will determine the filename based on object name, hence
  // figure out what Type(s)
  obj, err := st.FetchObject("important", record)
}
```	 


