package gtree_test

import (
	"fmt"
	"testing"

	"github.com/emirpasic/gods/trees/btree"
	"github.com/gogf/gf/v2/container/gtree"
	"github.com/gogf/gf/v2/util/gutil"
)

func TestBtreeMain(t *testing.T) {
	tree := btree.NewWithIntComparator(3) // empty (keys are of type int)

	tree.Put(1, "x") // 1->x
	tree.Put(2, "b") // 1->x, 2->b (in order)
	tree.Put(1, "a") // 1->a, 2->b (in order, replacement)
	tree.Put(3, "c") // 1->a, 2->b, 3->c (in order)
	tree.Put(4, "d") // 1->a, 2->b, 3->c, 4->d (in order)
	tree.Put(5, "e") // 1->a, 2->b, 3->c, 4->d, 5->e (in order)
	tree.Put(6, "f") // 1->a, 2->b, 3->c, 4->d, 5->e, 6->f (in order)
	tree.Put(7, "g") // 1->a, 2->b, 3->c, 4->d, 5->e, 6->f, 7->g (in order)

	json, err := tree.MarshalJSON()
	if err != nil {
		return
	}
	fmt.Println(string(json))

	b := tree.Iterator() // returns a stateful iterator whose elements are key/value pairs
	for b.Next() {
		index, value := b.Key(), b.Value()
		fmt.Println(index, value)
	}

}

func TestGtree(t *testing.T) {
	tree := gtree.NewBTree(3, gutil.ComparatorString)
	tree.Set(1, "x") // 1->x
	tree.Set(2, "b") // 1->x, 2->b (in order)
	tree.Set(1, "a") // 1->a, 2->b (in order, replacement)
	tree.Set(3, "c") // 1->a, 2->b, 3->c (in order)
	tree.Set(4, "d") // 1->a, 2->b, 3->c, 4->d (in order)
	tree.Set(5, "e") // 1->a, 2->b, 3->c, 4->d, 5->e (in order)
	tree.Set(6, "f") // 1->a, 2->b, 3->c, 4->d, 5->e, 6->f (in order)
	tree.Set(7, "g") // 1->a, 2->b, 3->c, 4->d, 5->e, 6->f, 7->g (in order)

	fmt.Println(tree)
}
