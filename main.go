/*
 * Copyright (C) 2021 Stuart Davies (stuartdd)
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */
package main

import (
	"fmt"
	"os"

	"github.com/stuartdd/jsonParserGo/parser"
)

var (
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
			"phoneNumbers": [{"type": "home"}, {"number": "7349282382"}, {"current":true, "loc":"UK"}]
		}
	 }`)
)

func main() {
	ExampleGet()
}

func ExampleGet() {
	n, err := parser.GetJsonParsed("http://n.n.n.n/files/temp.txt")
	if err != nil {
		panic(err.Error())
	}

	fmt.Println(n.JsonValueIndented(4))

	_, err = parser.PostJsonValueIndented(4, "http://n.n.n.n/files/temp2.txt", n)
	if err != nil {
		panic(err.Error())
	}
}

func ExampleDiagnostic() {
	root, err := parser.Parse(json)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(root.JsonValueIndented(4))
}

func ExampleWalkNodeTreeUntilTarget() {
	root, err := parser.Parse(json)
	if err != nil {
		panic(err.Error())
	}
	target, err := parser.Find(root, parser.NewDotPath("address.city"))
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

func LoadAfJsonFile() {
	if len(os.Args) < 2 {
		abortWithUsage("Missing file name")
	}
	filename := os.Args[1]
	dat, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Failed to read file %s. Error %s\n", filename, err.Error())
	}
	node, err := parser.Parse(dat)
	if err != nil {
		fmt.Printf("Failed to parse file %s. Error %s\n", filename, err.Error())
	}
	fmt.Println(parser.DiagnosticList(node))
}

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

func abortWithUsage(message string) {
	fmt.Printf(message+"\n  Usage: %s <filename>\n  Where: <filename> is a json file. E.g. config.json\n", os.Args[0])
	os.Exit(1)
}
