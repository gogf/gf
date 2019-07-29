package x2j

import (
	"fmt"
	"testing"
)

var doc01 = `
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
var doc02 = `
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

// the basic demo/test case - a small bibliography with mixed element types
func TestPathsForKey(t *testing.T) {
	fmt.Println("\n================================ x2jfindPath_test.go")
	fmt.Println("\n=============== TestPathsForKey ...")
	fmt.Println("\nPathsForKey... doc01#author")
	m, _ := DocToMap(doc01)
	ss := PathsForKey(m, "author")
	fmt.Println("ss:", ss)

	fmt.Println("\nPathsForKey... doc01#books")
	// m, _ := DocToMap(doc01)
	ss = PathsForKey(m, "books")
	fmt.Println("ss:", ss)

	fmt.Println("\nPathsForKey...doc02#book")
	m, _ = DocToMap(doc02)
	ss = PathsForKey(m, "book")
	fmt.Println("ss:", ss)

	fmt.Println("\nPathForKeyShortest...doc02#book")
	m, _ = DocToMap(doc02)
	s := PathForKeyShortest(m, "book")
	fmt.Println("s:", s)
}

// the basic demo/test case - a small bibliography with mixed element types
func TestPathsForTag(t *testing.T) {
	fmt.Println("\n=============== TestPathsForTag ...")
	fmt.Println("\nPathsForTag... doc01#author")
	ss, _ := PathsForTag(doc01, "author")
	fmt.Println("ss:", ss)

	fmt.Println("\nPathsForTag... doc01#books")
	ss, _ = PathsForTag(doc01, "books")
	fmt.Println("ss:", ss)

	fmt.Println("\nPathsForTag...doc02#book")
	ss, _ = PathsForTag(doc02, "book")
	fmt.Println("ss:", ss)

	fmt.Println("\nPathForTagShortest...doc02#book")
	s, _ := PathForTagShortest(doc02, "book")
	fmt.Println("s:", s)
}
