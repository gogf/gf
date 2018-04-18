// gitissue #28

/*
(reference: http://goessner.net/articles/JsonPath/)
Let's practice JSONPath expressions by some more examples. We start with a simple JSON structure built after an XML example representing a bookstore (original XML file).

{ "store": {
    "book": [
      { "category": "reference",
        "author": "Nigel Rees",
        "title": "Sayings of the Century",
        "price": 8.95
      },
      { "category": "fiction",
        "author": "Evelyn Waugh",
        "title": "Sword of Honour",
        "price": 12.99
      },
      { "category": "fiction",
        "author": "Herman Melville",
        "title": "Moby Dick",
        "isbn": "0-553-21311-3",
        "price": 8.99
      },
      { "category": "fiction",
        "author": "J. R. R. Tolkien",
        "title": "The Lord of the Rings",
        "isbn": "0-395-19395-8",
        "price": 22.99
      }
    ],
    "bicycle": {
      "color": "red",
      "price": 19.95
    }
  }
}
XPath                JSONPath                 Result
/store/book/author   $.store.book[*].author   the authors of all books in the store
//author             $..author                all authors
/store/*             $.store.*                all things in store, which are some books and a red bicycle.
/store//price        $.store..price           the price of everything in the store.
//book[3]            $..book[2]               the third book
//book[last()]       $..book[(@.length-1)]
                     $..book[-1:]             the last book in order.
//book[position()<3] $..book[0,1]
                     $..book[:2]              the first two books
//book[isbn]         $..book[?(@.isbn)]       filter all books with isbn number
//book[price<10]     $..book[?(@.price<10)]   filter all books cheapier than 10
//*                  $..*                     all Elements in XML document. All members of JSON structure.

*/

package main

import (
	"fmt"

	"github.com/clbanning/mxj"
)

var data = []byte(`
{ "store": {
    "book": [ 
      { "category": "reference",
        "author": "Nigel Rees",
        "title": "Sayings of the Century",
        "price": 8.95
      },
      { "category": "fiction",
        "author": "Evelyn Waugh",
        "title": "Sword of Honour",
        "price": 12.99
      },
      { "category": "fiction",
        "author": "Herman Melville",
        "title": "Moby Dick",
        "isbn": "0-553-21311-3",
        "price": 8.99
      },
      { "category": "fiction",
        "author": "J. R. R. Tolkien",
        "title": "The Lord of the Rings",
        "isbn": "0-395-19395-8",
        "price": 22.99
      }
    ],
    "bicycle": {
      "color": "red",
      "price": 19.95
    }
  }
}`)

func main() {
	m, err := mxj.NewMapJson(data)
	if err != nil {
		fmt.Println("NewMapJson err:", err)
		return
	}

	// $.store.book[*].author   the authors of all books in the store
	list, err := m.ValuesForPath("store.book.author")
	if err != nil {
		fmt.Println("book author err:", err)
		return
	}
	fmt.Println("authors:", list)

	// $..author                all authors
	list, err = m.ValuesForKey("author")
	if err != nil {
		fmt.Println("author err:", err)
		return
	}
	fmt.Println("authors:", list)

	// $.store.*                all things in store, which are some books and a red bicycle.
	list, err = m.ValuesForKey("store")
	if err != nil {
		fmt.Println("store things err:", err)
	}
	fmt.Println("store things:", list)

	// /store//price        $.store..price           the price of everything in the store.
	list, err = m.ValuesForKey("price")
	if err != nil {
		fmt.Println("price of things err:", err)
	}
	fmt.Println("price of things:", list)

	// $..book[2]               the third book
	v, err := m.ValueForPath("store.book[2]")
	if err != nil {
		fmt.Println("price of things err:", err)
	}
	fmt.Println("3rd book:", v)

	// $..book[-1:]             the last book in order
	list, err = m.ValuesForPath("store.book")
	if err != nil {
		fmt.Println("list of books err:", err)
	}
	if len(list) <= 1 {
		fmt.Println("last book:", list)
	} else {
		fmt.Println("last book:", list[len(list)-1:])
	}

	// $..book[:2]              the first two books
	list, err = m.ValuesForPath("store.book")
	if err != nil {
		fmt.Println("list of books err:", err)
	}
	if len(list) <= 2 {
		fmt.Println("1st 2 books:", list)
	} else {
		fmt.Println("1st 2 books:", list[:2])
	}

	// $..book[?(@.isbn)]       filter all books with isbn number
	list, err = m.ValuesForPath("store.book", "isbn:*")
	if err != nil {
		fmt.Println("list of books err:", err)
	}
	fmt.Println("books with isbn:", list)

	// $..book[?(@.price<10)]   filter all books cheapier than 10
	list, err = m.ValuesForPath("store.book")
	if err != nil {
		fmt.Println("list of books err:", err)
	}
	var n int
	for _, v := range list {
		if v.(map[string]interface{})["price"].(float64) >= 10.0 {
			continue
		}
		list[n] = v
		n++
	}
	list = list[:n]
	fmt.Println("books with price < $10:", list)

	// $..*                     all Elements in XML document. All members of JSON structure.
	// 1st where values are not complex elements
	list = m.LeafValues()
	fmt.Println("list of leaf values:", list)

	// $..*                     all Elements in XML document. All members of JSON structure.
	// 2nd every value - even complex elements
	path := "*"
	list = make([]interface{}, 0)
	for {
		v, _ := m.ValuesForPath(path)
		if len(v) == 0 {
			break
		}
		list = append(list, v...)
		path = path + ".*"
	}
	fmt.Println("list of all values:", list)
}
