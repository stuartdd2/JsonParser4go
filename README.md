# jsonParser4Go

**Some major refactoring is being done to simplify finding nodes and navigating the model.**

**These are breaking changes so avoid using the new version.**

Fix trailing ',' not detected.

---

Json parser written in GO (Why not!)

It returns a tree structure of different node types using maps and lists.

Each node has a specific type. It can be (cast) to a specific type and read/updated.

All objects returned from the API are pointer objects. This enables you to update the node in place and prevents copy on call which can consume more memory.

All String() methods return the value of the node as a String. For example for JsonNumber it returns the number as a String. For JsonBool it will return "true" of "false".

``` go
 \* Copyright (C) 2021 Stuart Davies (stuartdd)
 \*
 \* This program is free software: you can redistribute it and/or modify
 \* it under the terms of the GNU General Public License as published by
 \* the Free Software Foundation, either version 3 of the License, or
 \* (at your option) any later version.
 \*
 \* This program is distributed in the hope that it will be useful,
 \* but WITHOUT ANY WARRANTY; without even the implied warranty of
 \* MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 \* GNU General Public License for more details.
 \*
 \* You should have received a copy of the GNU General Public License
 \* along with this program.  If not, see <http://www.gnu.org/licenses/>.
```

#### To-Do

- Parse from a stream

- Profile (see how fast it is compared to other parsers)



### Import to your project

From within the project (where your go.mod) file is.

```bash
go get github.com/stuartdd2/JsonParser4go/parser
```

You should see somthing like:

```bash
go get: added github.com/stuartdd2/JsonParser4go/parser v0.0.0-20210923131243-6710c9cdd57d
```

The go.mod file for you project should now have the parser in the require section. For example:

```
require (
	github.com/stuartdd2/JsonParser4go/parser v0.0.0-20210923131243-6710c9cdd57d // indirect
)
```

To remove it:

```
go get github.com/stuartdd2/JsonParser4go/parser@none
```

## Parsing

The parse method takes a []byte as text.

Returns a node and an error. If error is not nil then node will be nil

```go
s := "[\"literal\", 1234, true]\"]"
rootNode, err := parser.Parse([]byte(s))
```

### Parsing a List

```go
rootNode, err := parser.Parse([]byte(`["literal", 1234, true]`))
```

### Parsing a file

```go
dat, err := os.ReadFile(filename)
if err != nil {
	fmt.Printf("Failed to read file %s. Error %s\n", filename, err.Error())
}
rootNode, err = parser.Parse(dat)
if err != nil {
	fmt.Printf("Failed to parse file %s. Error %s\n", filename, err.Error())
}
```

### Parsing a url (http get/post)

```go
n, err := parser.GetJsonParsed("http://n.n.n.n/files/temp.txt")
if err != nil {
	panic(err.Error())
}

fmt.Println(n.JsonValueIndented(4))

_, err = parser.PostJsonValueIndented(4, "http://n.n.n.n/files/temp2.txt", n)
if err != nil {
	panic(err.Error())
}
```

## Node types

All nodes implement the NodeI (capital i) interface.

| Data Type | Struct Name | Contains                             | Example                              |
| --------- | ----------- | ------------------------------------ | ------------------------------------ |
| String    | JsonString  | String literals                      | "name":"value"                       |
| Bool      | JsonBool    | true or false                        | "is":true                            |
| Number    | JsonNumber  | float64/int64                        | "count":10, "x":11.25                |
| List      | JsonList    | A list of any other Data Type        | [true,"abc",123] {"list":[{"A":[]}]} |
| Objects   | JsonObjects | A named list of any other data types |                                      |
| Null      | JsonNull    | null                                 | "thisIs":null                        |

All container nodes implement the NodeC interface. There are two container nodes, JsonObject and JsonList.

### Node Type and Enums

To detirmine a nodes type use the `node.GetNodeType()` function to return the NodeType enum.

| Enum Name | Struct (Node type) | Construct                                             |
| --------- | ------------------ | ----------------------------------------------------- |
| NT_STRING | JsonString         | NewJsonString(name string, value string) *JsonString  |
| NT_BOOL   | JsonBool           | NewJsonBool(name string, value bool) *JsonBool        |
| NT_NUMBER | JsonNumber         | NewJsonNumber(name string, value float64) *JsonNumber |
| NT_LIST   | JsonList           | NewJsonList(name string) *JsonList                    |
| NT_OBJECT | JsonObject         | NewJsonObject(name string) *JsonObject                |
| NT_NULL   | JsonNull           | NewJsonNull(name string) *JsonNull                    |

## Container nodes

List and Object nodes are container nodes. These implement the NodeC interface.

The function `IsContainer()` can be called on ANY node. It returns true if the note type is NT_OBJECT or NT_LIST. If this is true you can safely be cast to the NodeC interface.

A root JsonList node will not have a name and is rendered as [] if empty. JsonList nodes that have a name will be rendered "name":[]

A root JsonObject node will not have a name and is rendered as {} if empty. JsonObject nodes that have a name will be rendered "name":{}

### List Nodes

List nodes can contain literals or objects and implement the NodeC interface. For example:

```["literal", 123.4, true, {"num":99.9}, {"t":true}]```

The first 3 nodes are Literal nodes, a JsonString, JsonNumber and a JsonBool. They are rendered as literals because they do not have a name (name == "").

The other nodes are Object nodes (they all have names).

Example: To create the above json list:

```go
root := parser.NewJsonList("") // Create the root list node (no name)
root.Add(parser.NewJsonString("", "literal")) // Add a string literal (no name)
root.Add(parser.NewJsonNumber("", 123.4)) // Add a number literal (no name)
root.Add(parser.NewJsonBool("", true)) // Add a boolean literal (no name)
root.Add(parser.NewJsonNumber("num", 99.9)) // Add an number Object with name 'obj'
root.Add(parser.NewJsonBool("t", true)) // Add a boolean Object with name 't'
s := root.JsonValue() // s should contain the above json
```

The JsonList **Add(node NodeI) error** method will never return an error: The error return value is there only to comply with the NodeC interface. The return value can be ignored.

### Object Nodes

Example: Creating object nodes and adding objects to them.

```go
root := parser.NewJsonObject("") // Create the root object
root.Add(parser.NewJsonString("name", "ABC")) // Add a named string
root.Add(parser.NewJsonBool("male", true)) // Add a named boolean
root.Add(parser.NewJsonNumber("age", 123)) // Add a named number
fmt.Println(root.JsonValue()) // Prints {"name": "ABC","male": true,"age": 123}
```

The JsonObject **Add(node NodeI) error** method returns an error if:

- The node being added does not have a name
- The node being added has the same name as an existing node

The JsonObject node implements the NodeC interface

### Determine the node type

Example: Accessing nodes and their specific data access functions

```go
	name := node.GetName() // Get the node name on ANY node
	switch node.GetNodeType() { // Get the node type and evaluate with a switch statement
	case parser.NT_OBJECT:
        objectsNode := (node.(*parser.JsonObjects)) // Cast a node to a JsonObject
        sortedKeys := objectsNode.GetSortedKeys() // Call specific function on the node
	case parser.NT_LIST:
        listNode := (node.(*parser.JsonList)) // Cast a node to a JsonList
        count = listNode.Len() // Call specific function on the object node
	case parser.NT_NUMBER:
        numberNode := (node.(*parser.JsonNumber)) // Cast a node to a JsonNumber
        floafValue := numberNode.GetValue() // Now we can access the specific functions
        intValue := numberNode.GetIntValue()
	case parser.NT_STRING:
        stringValue := (node.(*parser.JsonString)).GetValue() // Get the string value
	case parser.NT_BOOL:
        boolValue := (node.(*parser.JsonBool)).GetValue() // Get the boolean value
    }
```

Also the **IsContainer()** function will return true if the node is a JsonList or JsonObject node and complies with a NodeC interface.

### Node Functions (NodeI (capital i) interface)

All Node Types have the 'NodeI' interface. This interface defines the following functions:

| function signature            | Desc                                                                                                                                                                                                                                                                                    |
| ----------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| String() string               | The Stringer interface. Returns the string representation of the value of the node. For example for JsonNumber it returns the number as a String. For JsonBool it will return "true" of "false".  For JsonList and JsonObject this will return the same as JsonValue() below            |
| JsonValue() string            | Returns the node as compact JSON. For example the address node in the JSON example below will return `"address": {"business": true,"phoneNumbers": [{"type": "home"},{"number": "7349282382"},{"current": true,"loc": "UK"}],"streetAddress": "101","city": "San Diego","state": "CA"}` |
| JsonValueIndented(int) string | Returns the node as indented (formatted) JSON. An example is included below:                                                                                                                                                                                                            |
| GetNodeType() NodeType        | Returns an ENUM of the node type, of type NodeType (see below)                                                                                                                                                                                                                          |
| GetName() string              | Returns the name of the node or "" if the node has no name                                                                                                                                                                                                                              |
| IsContainer() bool            | Returns true if the node implements the NodeC interface. This is defined by the nodes node type (GetNodeType()). It must be NT_OBJECT or NT_LIST.                                                                                                                                       |

### Container Functions (NodeC interface)

| function signature                 | Desc                                                                                                                                                      |
| ---------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------- |
| GetValues() []NodeI                | Return a list of ALL nodes in the container                                                                                                               |
| Len() int                          | Returns the number of nodes in the container                                                                                                              |
| Add(node NodeI) error              | Adds a node to a container                                                                                                                                |
| GetNodeWithName(name string) NodeI | Returns a node with a given name in the container. Note that lists can contain nodes that do not have names. These cannot be returned from this function. |
| Remove(nodeRemove NodeI) error     | Remove the node from the container node                                                                                                                   |

#### Example output indented by 4 spaces: ```JsonValueIndented(4) ```

```json
{
    "firstName": "Joe",
    "lastName": "Jackson",
    "gender": "male",
    "age": 28,
    "address": {
        "phoneNumbers": [
            {"type": "home"},
            {"number": "7349282382"},
            {
                "current": true,
                "loc": "UK"
            }
        ],
        "streetAddress": "101",
        "city": "San Diego",
        "state": "CA",
        "business": true
    }
}
```

### Specific Node Functions

These functions are based on the Node Type

| Node Type  | Func                                                  | Desc                                                                                                                                                                                                                                                                                                                                         |
| ---------- | ----------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| JsonString | NewJsonString(name string, value string) *JsonString  | Constructor. Creates a JsonString node with a name and a value. Returns a pointer to the node.                                                                                                                                                                                                                                               |
|            | GetValue() string                                     | Returns the 'string' value on the node                                                                                                                                                                                                                                                                                                       |
|            | SetValue(newValue string)                             | Updates the string value                                                                                                                                                                                                                                                                                                                     |
| JsonNumber | NewJsonNumber(name string, value float64) *JsonNumber | Constructor. Creates a JsonNumber node with a name and a value. Returns a pointer to the node.                                                                                                                                                                                                                                               |
|            | SetValue(newValue float64)                            | Sets the value of the nodes float64.                                                                                                                                                                                                                                                                                                         |
|            | SetIntValue(newValue int64)                           | Sets the value using an int64. Note internally the JsonNumber object stores a float64. The String() function will render '123.000000' as '123' and '123.100000' as '123.1' removing redundant '0' and '.' characters.                                                                                                                        |
|            | GetValue() float64                                    | Returns the value as a float64                                                                                                                                                                                                                                                                                                               |
|            | GetIntValue() int64                                   | Returns the value as an int. floating point values will be truncated.                                                                                                                                                                                                                                                                        |
| JsonBool   | NewJsonBool(name string, value bool) *JsonBool        | Constructor. Creates a JsonBool node with a name and a value. Returns a pointer to the node.                                                                                                                                                                                                                                                 |
|            | GetValue() bool                                       | Returns true or false                                                                                                                                                                                                                                                                                                                        |
|            | SetValue(newValue bool)                               | Sets the value to true or false                                                                                                                                                                                                                                                                                                              |
| JsonNull   | NewJsonNull(name string) *JsonNull                    | Constructor. Creates a JsonNull node with a name. The value is always 'null'. Returns a pointer to the node.                                                                                                                                                                                                                                 |
|            | Does not have any additional functions                | Renders as "name":null                                                                                                                                                                                                                                                                                                                       |
| JsonList   | NewJsonList(name string) *JsonList                    | Constructor. Creates a JsonList node with a name. The value is always an empty list. Returns a pointer to the node.                                                                                                                                                                                                                          |
|            | GetNodeAt(i int) NodeI                                | Returns the node a index i in the list.                                                                                                                                                                                                                                                                                                      |
|            | GetNodeWithName(name string) (NodeI, error)           | Returns the node with the name. If the node you are looking for does not have a name then this will return a error. <br/>Lists can contain nodes without names. These cannot be returned via this function, use GetNodeAt(n) instead. <br/>The order of nodes in a list is constant once they have been added, unlike nodes in a JsonObject. |
|            | Add(node NodeI)                                       | Adds a node. If the node to be added has a name a wrapper (parent) object (jsonObject) without a name is created and added to the list. If the node to be added does not have a name (a literal) then it is simply added to the list                                                                                                         |
|            | GetValues() []NodeI                                   | Returns a list of ALL the values. Altering this list (add, remove) has NO effect on the underlying JsonList.                                                                                                                                                                                                                                 |
|            | Remove(nodeRemove NodeI) error                        | Removes a given node from the JsonList. This does not use the node name. You need to Find the node first. If the node in NOT in the map an error is returned.                                                                                                                                                                                |
|            | Len()                                                 | Returns the combined number of literal objects and wrapper objects in the list.                                                                                                                                                                                                                                                              |
| JsonObject | NewJsonObject(name string) *JsonObject                | Constructor. Creates a JsonObject node with a name. The value is always an empty map. Returns a pointer to the node.                                                                                                                                                                                                                         |
|            | GetNodeWithName(name string) NodeI                    | Returns the node with the given name. Internally a map[string]\*NodeI contains all of the nodes. This simple returns the value. If the value is not found a nil is returned.                                                                                                                                                                 |
|            | GetSortedKeys() []string                              | Returns a list of keys from the map sorted a to z by name. Altering this list has NO effect on the underlying JsonObject.                                                                                                                                                                                                                    |
|            | GetValuesSorted() []NodeI                             | Returns a list of all the values in the map sorted a to z by name.  Altering this list (add, remove) has NO effect on the underlying JsonObject.                                                                                                                                                                                             |
|            | GetValues() []NodeI                                   | Returns a list of all the values in the map. Altering this list (add, remove) has NO effect on the underlying JsonObject.                                                                                                                                                                                                                    |
|            | Add(node NodeI)                                       | Adds the node to the map. It uses the node name as the map key.                                                                                                                                                                                                                                                                              |
|            | Remove(nodeRemove NodeI) error                        | Removes a given node from the JsonObject. This does not use the node name. You need to Find the node first. If the node in NOT in the map an error is returned.                                                                                                                                                                              |
|            | Len()                                                 | Returns the combined number of objects in the list.                                                                                                                                                                                                                                                                                          |

## Static functions (utils)

### A Path to a Node

A path to ANY node in a node tree can be described by a 'path' to the node. The Path type is designed to represent that path,  The following functions (described below) require a Path to be specified.

- Find(node NodeI, path *Path) (NodeI, error)
- CreateAndReturnNodeAtPath(root NodeI, path *Path, nodeType NodeType) (NodeI, error)

A Path is can be defined as follows:

- path "a.b.c" delim "." - The path separator (or delimiter) is the '.' character.

This example path could be used to find a node 'c' in a container node 'b' in a root container node 'a'. The '.' separator for the paths is defined by the second parameter 'delim'.

- path "a|b|c" delim "|" - The path separator (or delimiter) is the '.' character.

This example path could be used to find a node 'c' in a container node 'b' in a root container node 'a'. The '.' separator for the paths is defined by the second parameter 'delim'.



| Function                          | description                                                                                                                                                                                                                                                                                                                                | Example                                                                      |
| --------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | ---------------------------------------------------------------------------- |
| NewPath(path, delim string) *Path | Returns a path defined by as string and a delimeter. The example path could be used to find a node 'c' in a container 'b' in a root container 'a'. Note the '.' separator for the paths is defined by the second parameter 'delim'. Both of the examples are equivalent unless there is a node with the name 'a.b'. See the third example. | p:=NewPath("a.b.c",".") p:=NewPath("a\|b\|c","\|") p:=NewPath("a.b\|c","\|") |
|                                   |                                                                                                                                                                                                                                                                                                                                            |                                                                              |
|                                   |                                                                                                                                                                                                                                                                                                                                            |                                                                              |



These functions are stand alone utilities:

These function do NOT include the Structure creation functions such as ```NewJsonString(name string, value string)```. These are already covered above.

| Function name                                                                                    | Desc                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                       |
| ------------------------------------------------------------------------------------------------ | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Parse(json []byte) (node NodeI, err error)                                                       | Parse a []byte array from a string or file. Returns a root node or an error if the parser fails. This node will be either a JsonList or a JsonObject node.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                 |
| Find(node NodeI, path *Path) (NodeI, error)                                                      | Find a node from a Path. Returns the node or an error if not found                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                         |
| FindParentNode(root, target NodeI) (NodeI, bool)                                                 | Find the parent of a node. This function uses the WalkNodeTree function to search the tree structure (based at root) for the specific target node. If found it returns the parent node and a boolean indicating sucess. If not found it returns nil, false.  See below for an example                                                                                                                                                                                                                                                                                                                                                                      |
| WalkNodeTree(root, target NodeI, onEachNode func(NodeI, NodeI, NodeI) bool) (NodeI, NodeI, bool) | Walk the node tree and visit each node. This function starts at the 'root' node and visits EVERY node under the 'root'. For each node it visits it will call the function 'onEachNode' passing in the current node (first parameter), the current nodes parent node (the second parameter) and the 'target' node (the third parameter). If the 'onEachNode' function returns 'true' the walk is terminated and the current node, it's parent and 'true' are returned from WalkNodeTree. If  'onEachNode'  never returns true untill all nodes have been visited then nil, nil, false will be returned from WalkNodeTree function. See below for an example |
| DiagnosticList(node NodeI) string                                                                | Returns a string that represents the structure of the node provided. This is really useful if you cannot see how the node objects are packed.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                              |
| GetNodeTypeName(tt NodeType) string                                                              | Given a NodeType this returns a string value for the type. NT_STRING returns "STRING"                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                      |
| Remove(root, node NodeI) error                                                                   | Removes the node from the root. This will search for the node in the nodes under and including the root node. It will then remove that node if found. Will return an error if not found or the root node is the node to be removed.                                                                                                                                                                                                                                                                                                                                                                                                                        |
| Rename(root, node NodeI, newName string) error                                                   | Renames the node in the root. This will search for the node in the nodes under the root node. It will then rename that node if found. Will return an error if not found or the rename would cause a duplicate in a container node.                                                                                                                                                                                                                                                                                                                                                                                                                         |
| CreateAndReturnNodeAtPath(root NodeI, path *Path, nodeType NodeType) (NodeI, error)              | Given a path from the root node, this method will ensure that all nodes in that path exist and that the last (leaf) node is of the correct type. It will create all of the required nodes. An error is returned if the leaf node is found but is not of the correct type. An error is returned if any existing intermediate nodes are not container nodes.                                                                                                                                                                                                                                                                                                 |

### Web functions

| Func                                                                       | Desc                                                                                                                                                                                             |
| -------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| GetJsonParsed(getUrl string) (NodeI, error)                                | Fetch and parse json data from a URL using the HTTP GET protocol                                                                                                                                 |
| func GetJson(getUrl string) ([]byte, error)                                | Fetch data from a URL using the HTTP GET protocol                                                                                                                                                |
| PostData(postUrl string, contentType string, data []byte) ([]byte, error)  | Send data to a URL using the HTTP POST protocol. Content type must be defined. For Json it should be "application/json".                                                                         |
| PostJsonBytes(postUrl string, data []byte) ([]byte, error)                 | Uses PostData with a content type "application/octet-stream". This is usefull if the data is encrypted.                                                                                          |
| PostJsonText(postUrl string, data []byte) ([]byte, error)                  | Uses PostData with a content type "application/json". This is usefull if the data is plain json.                                                                                                 |
| PostJsonValue(postUrl string, node NodeI) ([]byte, error)                  | Uses PostData with a content type "application/json". The data is derived from the JsonValue() function on the node provided. This send plain json without formatting.                           |
| PostJsonValueIndented(tab int, postUrl string, node NodeI) ([]byte, error) | Uses PostData with a content type "application/json". The data is derived from the JsonValueIndented(n) function on the node provided. This send json formatted with a indent (tab) of n spaces. |

### Searching the tree

The Find function will search the tree structure for a specific node. It uses a path string to define the search path:

Given the following json:

```go
json = []byte(`{
		"firstName": "Joe",
		"lastName": "Jackson",
		"gender": "male",
		"age": 28,
		"address": {
			"streetAddress": "101",
			"city": "San Diego",
			"state": "CA",
			"business": true,
			"phoneNumbers": [{"type": "home"}, {"number": "7349282382"} {""}]
		}
	 }`)
```

Example: Find the address phone number.

```go
func ExampleFindWithPath() {
	root, err := parser.Parse(json)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(parser.Find(root, parser.NewDotPath("address.phoneNumbers.number")))    // Will find single value named nodes
	fmt.Println(parser.Find(root, parser.NewDotPath("address.phoneNumbers.0")))         // The first node
	fmt.Println(parser.Find(root, parser.NewDotPath("address.phoneNumbers.1")))         // The 2nd node
	fmt.Println(parser.Find(root, parser.NewDotPath("address.phoneNumbers.2.current"))) // The third node 'current'

	fmt.Println(parser.Find(root, parser.NewDotPath("address.city")))
	fmt.Println(parser.Find(root, parser.NewDotPath("address")))
}
OUTPUT: (note the <nil> is printed bacause the error is nil)
7349282382 <nil>
home <nil>
7349282382 <nil>
true <nil>
```

**Note:** Although 'phoneNumbers' is a list it will search for named objects within the list.

**Note:** The index notation can be used for accessing elements in lists directly. It will return literals as well as objects.

Be warned that the order of the lists (and objects) cannot be assured so this is a last stop for finding elements in a list.

Example: Find the city node of the address

```go
fmt.Println(parser.Find(root, parser.NewDotPath("address.city"))
OUTPUT: (note the <nil> is printed bacause the error is nil)
San Diego <nil>
```

Example: Find the address node

```go
fmt.Println(parser.Find(node, "address"))
"address":{"business": true,"phoneNumbers":[{"type": "home"},{"number": "7349282382"},{"current": true,"loc": "UK"}],"streetAddress": "101","city": "San Diego","state": "CA"} <nil>
```

Note: The default String() method on JsonList and JsonObject returns JsonValue(). For all other node types it returns the String of the value of the node.

### Finding a parent of a Node

The nodes in the tree structure do not have a pointer to their parent. This makes the tree smaller and faster to create.

Given any node (target) currently in the tree this fuction will return it's parent node.

With this function I can:

- Confirm that a node is in the tree
- Find the parent node of any node in the tree
- This function uses the WalkNodeTree function (see below)

Given the above JSON (obj2)

```go
func ExampleFindParent() {
	root, _ := parser.Parse(json)
	target1 := parser.NewJsonString("state", "CA")                      // This node is NOT in the root tree
	target2, _ := parser.Find(root, parser.NewDotPath("address.state")) // This node in in the root tree. Check err!

	fmt.Println(target1.String()) // Will print "state": "CA". These look the same!
	fmt.Println(target2.String()) // Will print "state": "CA"

	p2, ok2 := parser.FindParentNode(root, target2)
	fmt.Printf("%t\n", ok2)          // will print true because target2 is in the node tree
	fmt.Printf("%s\n", p2.GetName()) // will print the name of the parent of target2 'address':

	_, ok1 := parser.FindParentNode(root, target1)
	fmt.Printf("%t\n", ok1) // will print false because target1 is outside the node tree.
	//	The parent is nil
}
```

### Walking the tree

Walking the tree enables your logic to be applied to EVERY node in the tree. For example to print all the node names:

```go
func ExampleWalkNodeTree() {
	root, err := parser.Parse(json)
	if err != nil {
		panic(err.Error())
	}
	parser.WalkNodeTree(root, nil, func(n, p, t parser.NodeI) bool {
		if n.GetName() != "" { // If the node has a name
			fmt.Printf("%s,", n.GetName()) // Print the node name
		}
		return false // Continue until the end!
	})
}
OUTPUT:
firstName,lastName,gender,age,address,streetAddress,city,state,business,phoneNumbers,type,number,current,loc,
```

If you require to return at a specific node then pass in the target and return true when the conditions are met.

```
func ExampleWalkNodeTreeUntilConditionMet() {
	root, err := parser.Parse(json)
	if err != nil {
		panic(err.Error())
	}
	n, p, ok := parser.WalkNodeTree(root, nil, func(n, p, t parser.NodeI) bool {
		return n.String() == "CA"
	})
	if ok {
		fmt.Printf("Node   Name: %s\n", n.GetName()) // Print the node with the value 'CA'
		fmt.Printf("Parent Name: %s\n", p.GetName()) // Print the parent of the node with the value 'CA'
	}
}
OUTPUT:
Node   Name:state
Parent Name:address
```

If you pass a target in you can use that to meet the criteria.

```go
func ExampleWalkNodeTreeUntilTarget() {
	root, err := parser.Parse(json)
	if err != nil {
		panic(err.Error())
	}
	target, err := parser.Find(root, "address.city")
	if err != nil {
		panic("Target not found")
	}
	n, p, ok := parser.WalkNodeTree(root, target, func(n, p, t parser.NodeI) bool {
		return n == t // If the node equals the target this returns true!
	})
	if ok {
		fmt.Printf("Node   Name: %s\n", n.GetName()) // Print the node with the value 'CA'
		fmt.Printf("Parent Name: %s\n", p.GetName()) // Print the parent of the node with the value 'CA'
	}
}
OUTPUT:
Node   Name: city
Parent Name: address
```

This is how FindParentNode is implemented. It just discards the node (n) above.

### Understanding the structure

Given the above JSON:

Using the following code will print the structure of the objects on the node tree.

```go
fmt.Println(parser.DiagnosticList(node))
Diag
OBJECT: N:''
  OBJECT: N:'address'
    BOOL: N:'business' V:'true'
    STRING: N:'city' V:'San Diego'
    LIST: N:'phoneNumbers'
      STRING: N:'type' V:'home'
      STRING: N:'number' V:'7349282382'
      OBJECT: N:''
        BOOL: N:'current' V:'true'
        STRING: N:'loc' V:'UK'
    STRING: N:'state' V:'CA'
    STRING: N:'streetAddress' V:'101'
  NUMBER: N:'age' V:'28.000000'
  STRING: N:'firstName' V:'Joe'
  STRING: N:'gender' V:'male'
  STRING: N:'lastName' V:'Jackson'
```

This shows what type of node each node in the tree is. It also shows its name (N) and its value (V).

