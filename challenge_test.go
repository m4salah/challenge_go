package challenge_test

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"testing"
	"unicode"

	"github.com/stretchr/testify/assert"
)

// Problem 01
// A palindromeRecurse is a word or sequence of characters which reads the same backward or
// forward. For example the words: level, anna, noon, rotator are all palindromes.
// Write a function palindrom that accepts a string as an argument and returns a boolean
// indicating whether the input is a palindromeRecurse or not, for example:
// palindromeRecurse("anna") # returns True
// palindromeRecurse("apple") # returns False
func palindromeRecurse(word string) bool {
	wordLength := len(word)
	if wordLength == 0 {
		return true
	}
	return word[0] == word[wordLength-1] && palindromeRecurse(word[1:wordLength-1])
}

func palindromeIter(word string) bool {
	wordLength := len(word)
	i, j := 0, wordLength-1
	for {
		if i >= j {
			return true
		}
		if word[i] == word[j] {
			i++
			j--
			continue
		}
		return false
	}
}

func TestPalindrome(t *testing.T) {
	testCases := []struct {
		input  string
		expect bool
	}{
		{input: "anna", expect: true},
		{input: "apple", expect: false},
		{input: "", expect: true},
	}

	for _, tc := range testCases {
		t.Run(
			fmt.Sprintf("palindromeRecurse should return %t for input=%s", tc.expect, tc.input),
			func(t *testing.T) {
				result := palindromeRecurse(tc.input)
				assert.Equal(t, tc.expect, result)
			},
		)
	}

	for _, tc := range testCases {
		t.Run(
			fmt.Sprintf("palindromeIter should return %t for input=%s", tc.expect, tc.input),
			func(t *testing.T) {
				result := palindromeIter(tc.input)
				assert.Equal(t, tc.expect, result)
			},
		)
	}
}

// Problem 02
// Write the Linux command needed to change a file name from original.txt to
// changed.txt
// > mv original.txt changed.txt
func changeFilename(originalName string, newName string) error {
	// using mv command we can change the file name
	cmd := exec.Command("mv", originalName, newName)
	_, err := cmd.Output()
	return err
}

func TestChangeFilename(t *testing.T) {
	t.Run("should return nil error", func(t *testing.T) {
		// Set up file paths
		originalFilename := "./tmp/original.txt"
		newFilename := "./tmp/changed.txt"

		// Ensure the ./tmp directory exists
		err := os.MkdirAll("./tmp", 0755)
		assert.NoError(t, err, "should create the tmp directory")

		// Create the original file
		file, err := os.Create(originalFilename)
		assert.NoError(t, err, "should create the original file")

		// Clean the files created by the test
		defer file.Close()
		defer os.RemoveAll("./tmp")

		err = changeFilename(originalFilename, newFilename)
		assert.NoError(t, err, "should rename the file")
	})
}

// Problem 03
// Given a string containing characters (a-z), implement a function runLengthEncode that
// compresses repeated ‘runs’ of the same character by storing the length of that run, and
// provide a function runLengthDecode to reverse the compression. The output can be
// anything, as long as you can recreate the input with it.
// For example:
// runLengthEncode("aaaaaaaaaabbbaxxxxyyyzyx") # returns "a10b3a1x4y3z1y1x1"
// runLengthDecode("a10b3a1x4y3z1y1x1") # returns "aaaaaaaaaabbbaxxxxyyyzyx"

func mustNoErr(err error) {
	if err != nil {
		panic(err)
	}
}
func runLengthEncode(input string) string {
	// Edge case: empty string
	if len(input) == 0 {
		return ""
	}
	result := strings.Builder{}
	repeats := 1
	for i := 1; i < len(input); i++ {
		if input[i] == input[i-1] {
			repeats++
		} else {
			err := result.WriteByte(input[i-1])
			mustNoErr(err)
			_, err = result.WriteString(fmt.Sprintf("%d", repeats))
			mustNoErr(err)
			repeats = 1
		}
	}

	result.WriteByte(input[len(input)-1])
	result.WriteString(fmt.Sprintf("%d", repeats))

	return result.String()
}

// extract first number encounter in a string, and how much step to skip this number
func extractFirstNumberFromString(input string) (int, int) {
	beginNumber, endNumber := 0, 0
	for i, c := range input {
		if unicode.IsDigit(c) {
			beginNumber = i
			endNumber = i + 1
			break
		}
	}

	for i, c := range input[beginNumber+1:] {
		if unicode.IsDigit(c) {
			endNumber += i + 1
		} else {
			break
		}
	}
	n, err := strconv.Atoi(input[beginNumber:endNumber])
	mustNoErr(err)
	return n, endNumber - beginNumber + 1
}

func runLengthDecode(input string) string {
	result := strings.Builder{}
	i := 0
	for i < len(input)-1 {
		repeats, skip := extractFirstNumberFromString(input[i+1:])
		result.WriteString(strings.Repeat(string(input[i]), repeats))
		i += skip
	}
	return result.String()
}

func TestLengthEncodeDecode(t *testing.T) {
	t.Run("check if the decoded encoded string equal to the original string", func(t *testing.T) {
		originalStrings := []string{
			"aaaaaaaaaabbbaxxxxyyyzyx",
			"a",
			"aa",
			"",
			"ab",
			"aab",
		}

		for _, originalString := range originalStrings {
			assert.Equal(
				t,
				originalString,
				runLengthDecode(runLengthEncode(originalString)),
			)

		}

	})
}

// Problem 04
// Let f and g be two one-argument functions. The composition f after g is defined to be the
// function . Define a function compose that implements composition. For
// example, if inc is a function that adds 1 to its argument, and square is a function that
// squares its argument, then:
// h = compose(square, inc) # create a new function h by composing inc and square
// h(6) # returns 49
func compose[A, B, C any](f1 func(B) C, f2 func(A) B) func(A) C {
	return func(input A) C {
		return f1(f2(input))
	}
}

func inc(n int) int {
	return n + 1
}

func square(n int) int {
	return n * n
}

func TestComposeFunction(t *testing.T) {
	t.Run("test componse function", func(t *testing.T) {
		h := compose(square, inc)
		assert.Equal(t, 49, h(6))

		h = compose(inc, square)
		assert.Equal(t, 37, h(6))

	})
}

// Problem 05
// Write a function unique that takes an array of strings as input and returns an array of
// the unique entries in the input, for example:
// unique(['apples', 'oranges', 'flowers', 'apples']) # returns ['oranges', 'flowers']
// unique(['apples', 'apples']) # returns []

func unique(words []string) []string {
	wordFreq := make(map[string]int)
	for _, word := range words {
		freq, ok := wordFreq[word]
		if ok {
			wordFreq[word] = freq + 1
		} else {
			wordFreq[word] = 1
		}
	}

	uniqueWords := []string{}
	for word, freq := range wordFreq {
		if freq == 1 {
			uniqueWords = append(uniqueWords, word)
		}
	}
	return uniqueWords
}

func TestUniqueFunction(t *testing.T) {
	testCases := []struct {
		input  []string
		expect []string
	}{
		{
			input:  []string{"apples", "oranges", "flowers", "apples"},
			expect: []string{"oranges", "flowers"},
		},
		{
			input:  []string{"apples", "apples"},
			expect: []string{},
		},
	}

	for _, tc := range testCases {
		t.Run(
			fmt.Sprintf("unique should return %v for input=%v", tc.expect, tc.input),
			func(t *testing.T) {
				result := unique(tc.input)
				assert.Equal(t, tc.expect, result)
			},
		)
	}
}

// Problem 06
// In linear algebra, the transpose of a matrix is another matrix created by writing the
// rows of as the columns of , for example:
// Write a function transpose that transposes a matrix, for example:
// transpose( [ [1,2], [3,4] ] ) # returns [ [1,3], [2,4] ]
// transpose( [ [1,2,3,4], [5,6,7,8] ] ) # returns [ [1,5], [2,6], [3,7], [4,8] ]

func transpose(matrix [][]int) [][]int {
	tansposedMatrix := [][]int{}

	numberOfRows := len(matrix)
	numberOfCols := len(matrix[0])
	for colIndex := 0; colIndex < numberOfCols; colIndex++ {
		transposedRow := []int{}
		for rowIndex := 0; rowIndex < numberOfRows; rowIndex++ {
			transposedRow = append(transposedRow, matrix[rowIndex][colIndex])
		}
		tansposedMatrix = append(tansposedMatrix, transposedRow)

	}
	return tansposedMatrix
}
func TestTransposeFunction(t *testing.T) {
	testCases := []struct {
		input  [][]int
		expect [][]int
	}{
		{
			input:  [][]int{[]int{1, 2}, []int{3, 4}},
			expect: [][]int{[]int{1, 3}, []int{2, 4}},
		},
		{
			input:  [][]int{[]int{1, 2, 3, 4}, []int{5, 6, 7, 8}},
			expect: [][]int{[]int{1, 5}, []int{2, 6}, []int{3, 7}, []int{4, 8}},
		},
	}

	for _, tc := range testCases {
		t.Run(
			fmt.Sprintf("transpose should return %v for input=%v", tc.expect, tc.input),
			func(t *testing.T) {
				result := transpose(tc.input)
				assert.Equal(t, tc.expect, result)
			},
		)
	}
}

// Problem 07
// You are given 2 containers: A and B. Container A can hold 5 litres of water, while
// container B can hold 3 litres. You are also given a source of water that you can use as
// you wish. Show how you can use the containers and the water source to put exactly 4
// litres of water in container A. No coding required, just write down the steps.
// 1. Fill container A          (A -> 5, B -> 0)
// 2. Pour water from A to B    (A -> 2, B -> 3) now B is full and A has 2 liter only.
// 3. Empty Container B         (A -> 2, B -> 0)
// 4. Pour Water from A to B    (A -> 0, B -> 2)
// 5. Fill A                    (A -> 5, B -> 2) now B can hold 1 liter only.
// 6. Pour water from A to B    (A -> 4, B -> 3)

// Problem 08
// Given an integer array of length n, find the index of the first duplicate element in the
// array and state the runtime and space complexity of your implementation, for example:
// # returns 3, assuming the index starts with 0
// index_of_first_duplicate( [ 5, 17, 3, 17, 4, -1 ] )

// Space complexity: O(n) -> the map size.
// Time complexity:  O(n) -> We make one iteration on the list, the map lookup is around O(1).
func indexOfFirstDuplicate(numbers []int) int {
	// construct a map that key is the number and value is it's index.
	// This construction is O(n) if setting the key value in the map is O(1)
	indexMap := map[int]int{}
	for index, number := range numbers {
		// check if that number appears previously in the map.
		if _, ok := indexMap[number]; ok {
			return index
		}
		indexMap[number] = index
	}
	// -1 indicate that we did not found any duplicate.
	return -1
}

func TestIndexOfFirstDuplicate(t *testing.T) {
	testCases := []struct {
		input  []int
		expect int
	}{
		{
			input:  []int{5, 17, 3, 17, 4, -1},
			expect: 3,
		},
		{
			input:  []int{1, 2, 3},
			expect: -1,
		},
	}

	for _, tc := range testCases {
		t.Run(
			fmt.Sprintf("indexOfFirstDuplicate should return %v for input=%v", tc.expect, tc.input),
			func(t *testing.T) {
				result := indexOfFirstDuplicate(tc.input)
				assert.Equal(t, tc.expect, result)
			},
		)
	}
}

// Problem 09
// Given the below tree structure, write a function sum that accepts a node and returns the
// sum of integers for the whole tree represented by the given node argument
type Node struct {
	value    int
	children []Node
}

func (n *Node) Sum() int {
	childrenSum := 0
	for _, node := range n.children {
		childrenSum += node.Sum()
	}
	return n.value + childrenSum
}

func TestNodeSum(t *testing.T) {
	testCases := []struct {
		input  Node
		expect int
	}{

		{
			input:  Node{},
			expect: 0,
		},
		{
			input:  Node{value: 5, children: []Node{}},
			expect: 5,
		},
		{
			input: Node{value: 5, children: []Node{
				Node{value: 6, children: []Node{}},
			}},
			expect: 11,
		},
		{
			input: Node{value: 5, children: []Node{
				Node{value: 6, children: []Node{
					Node{value: 6, children: []Node{}},
				}},
			}},
			expect: 17,
		},
	}

	for _, tc := range testCases {
		t.Run(
			fmt.Sprintf("node sum should return %v for input=%v", tc.expect, tc.input),
			func(t *testing.T) {
				result := tc.input.Sum()
				assert.Equal(t, tc.expect, result)
			},
		)
	}
}
