package x2j

import (
	"fmt"
	"testing"
)

var doc1 = `
<doc>
	<books>
		<book seq="1">
			<author>William H. Gaddis</author>
			<title>The Recognitions</title>
			<review>One of the great seminal American novels of the 20th century.</review>
		</book>
		<book seq="2">
			<author>Austin Tappan Wright</author>
			<title>Islandia</title>
			<review>An example of earlier 20th century American utopian fiction.</review>
		</book>
		<book seq="3">
			<author>John Hawkes</author>
			<title>The Beetle Leg</title>
			<review>A lyrical novel about the construction of Ft. Peck Dam in Montana.</review>
		</book>
		<book seq="4">
			<author>
				<first_name>T.E.</first_name>
				<last_name>Porter</last_name>
			</author>
			<title>King's Day</title>
			<review>A magical novella.</review>
		</book>
	</books>
</doc>
`
var doc2 = `
<doc>
	<books>
		<author>
			<name>William H. Gaddis</name>
			<book seq="1">
				<title>The Recognitions</title>
				<review>One of the great seminal American novels of the 20th century.</review>
			</book>
			<book>
				<title>JR</title>
				<review>Won the National Book Award</review>
			</book>
		</author>
		<author>
			<name>John Hawkes</name>
			<books>
				<book>
					<title>The Beetle Leg</title>
				</book>
				<book>
					<title>The Blood Oranges</title>
				</book>
			</books>
		</author>
	</books>
</doc>
`
var msg1 = `
<msg>
	<pub>test</pub>
	<text>This is a long cold winter</text>
</msg>`

var msg2 = `
<msgs>
	<msg>
		<pub>test</pub>
		<text>This is a long cold winter</text>
	</msg>
	<msg>
		<pub>test2</pub>
		<text>I hope we have a cool summer, though</text>
	</msg>
</msgs>`

func TestValuesAtKeyPath(t *testing.T) {
	fmt.Println("\n============================ x2jat_test.go")
	fmt.Println("\n=============== TestValuesAtKeyPath ...")
	fmt.Println("\nValuesAtKeyPath ... doc1#author")
	m, _ := DocToMap(doc1)
	ss := PathsForKey(m,"author")
	fmt.Println("ss:", ss)
	for _,v := range ss {
		vv := ValuesAtKeyPath(m,v,true)
		fmt.Println("vv:", vv)
	}

	fmt.Println("\nValuesAtKeyPath ... doc1#first_name")
	// m, _ := DocToMap(doc1)
	ss = PathsForKey(m,"first_name")
	fmt.Println("ss:", ss)
	for _,v := range ss {
		vv := ValuesAtKeyPath(m,v,true)
		fmt.Println("vv:", vv)
	}

	fmt.Println("\nGetKeyPaths...doc2#book")
	m, _ = DocToMap(doc2)
	ss = PathsForKey(m,"book")
	fmt.Println("ss:", ss)
	for _,v := range ss {
		vv := ValuesAtKeyPath(m,v,true)
		fmt.Println("vv:", vv)
	}
	s := PathForKeyShortest(m,"book")
	vv := ValuesAtKeyPath(m,s)
	fmt.Println("vv,shortest_path:",vv)

	fmt.Println("\nValuesAtKeyPath ... msg1#pub")
	m, _ = DocToMap(msg1)
	ss = PathsForKey(m,"pub")
	fmt.Println("ss:", ss)
	for _,v := range ss {
		vv := ValuesAtKeyPath(m,v,true)
		fmt.Println("vv:", vv)
	}

	fmt.Println("\nValuesAtKeyPath ... msg2#pub")
	m, _ = DocToMap(msg2)
	ss = PathsForKey(m,"pub")
	fmt.Println("ss:", ss)
	for _,v := range ss {
		vv := ValuesAtKeyPath(m,v,true)
		fmt.Println("vv:", vv)
	}
}

func TestValuesAtTagPath(t *testing.T) {
	fmt.Println("\n=============== TestValuesAtTagPath ...")
	fmt.Println("\nValuesAtTagPath ... doc1#author")
	m, _ := DocToMap(doc1)
	ss := PathsForKey(m,"author")
	fmt.Println("ss:", ss)
	for _,v := range ss {
		vv,_ := ValuesAtTagPath(doc1,v,true)
		fmt.Println("vv:", vv)
	}

	fmt.Println("\nValuesAtTagPath ... doc1#first_name")
	// m, _ := DocToMap(doc1)
	ss = PathsForKey(m,"first_name")
	fmt.Println("ss:", ss)
	for _,v := range ss {
		vv,_ := ValuesAtTagPath(doc1,v,true)
		fmt.Println("vv:", vv)
	}

	fmt.Println("\nValuesAtTagPath...doc2#book")
	m, _ = DocToMap(doc2)
	ss = PathsForKey(m,"book")
	fmt.Println("ss:", ss)
	for _,v := range ss {
		vv,_ := ValuesAtTagPath(doc2,v,true)
		fmt.Println("vv:", vv)
	}
	s := PathForKeyShortest(m,"book")
	vv,_ := ValuesAtTagPath(doc2,s)
	fmt.Println("vv,shortest_path:",vv)
}

