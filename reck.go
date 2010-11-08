package main

import strings "strings"
import fmt "fmt"
import bufio "bufio"
import os "os"
import strconv "strconv"
import list "container/list"

// repl doesn't work at all now, but crashes. so 1. fix that
// then make a Val interface then types Int, Float, String (with getType getVal and toString) funcs) (look at go how this is done)
// make the reader so that it pases whole block char by char and creates the Ents (entities)


const (
	tnone	= iota
	tstring = iota
	tint	= iota
	tnative = iota
	tword	= iota
);

type Val interface {
	gst() string
	giv() int
	gsv() string
};

type Context struct {
	words *map[TWord]Val
	parent Context
}

type Unit struct {
	t int
	s string
}

type TNative struct {
	n	string
	pc	int
	f	func(s *list.List) ToyVal
};

func (V *TNative) gst() string { return "Native" }
func (V *TNative) gsv() string { return V.n }
func (V *TNative) giv() int { return tnone }

type TWord struct {
	v string
}

func (V *TWord) gst() string { return "Word" }
func (V *TWord) gsv() string { return V.v }
func (V *TWord) giv() int { return tnone }

type TString struct {
	v string
}

func (V *TString) gst() string { return "String" }
func (V *TString) gsv() string { return V.v }
func (V *TString) giv() int { return tnone }

type TInt struct {
	v int
}

func (V *TInt) gst() string { return "Int" }
func (V *TInt) gsv() string { return "" }
func (V *TInt) giv() int { return V.v }

type Reader struct {
	state int
}

func (r *Reader) doBlock (code string, i int) {
	// switch cases and set states
	print code[i]
	if len(code) < i {
		doBlock (code, i+1)
	}
}



func main() {
	println("~ reck interpreter v0.0002 ~ ");
	var in *bufio.Reader = bufio.NewReader(os.Stdin); 
	var inp string;
	for (strings.Trim(inp, "\n\r") != "exit!") {
		print(">> ");
		inp, _ := in.ReadBytes('\n')
		doBlock(string(inp));
		//printGLOB();
	}
}














/*
var GLOB = make(map[string] ToyVal, 100);
var NATIVE = make(map[string] Native, 100);
var STACK = list.New();


func getTag(tags []string, idx int) string {
	var tag = strings.Trim(tags[idx], "\n\r");
	//fmt.Printf("(tag: %s)", tag );
	return tag;
}

func isSetWord(tag string) bool { 
	return strings.HasSuffix(tag, ":"); 
}

func isString(tag string) bool { 
	return strings.HasPrefix(tag, "\"") && strings.HasSuffix(tag, "\"");
}

func isInt(tag string) bool { 
	var _, ier = strconv.Atoi(tag);
	return ier == nil;
}

func isNative(tag string) bool { 
	_, ok := NATIVE[tag];
	return  ok;
}

func doTag(tags []string, idx int) (ToyVal, int) {
	var tag = getTag(tags, idx);

	if isSetWord(tag) {

		GLOB[tag], idx = doTag(tags, idx + 1);	
		return GLOB[tag], idx;

	} else if isString(tag) {

		return ToyVal{t: tstring, s: tag}, idx + 1;

	} else if isInt(tag) {

		var iv, _ = strconv.Atoi(tag);
		return ToyVal{t: tint, i: iv }, idx + 1;

	} else if isNative(tag) {
		var n = NATIVE[tag];
		var s = list.New();
		idx = idx + 1;
		for i := 0; i<n.pc ; i++ {
			var V ToyVal;
			V, idx = doTag(tags, idx)
			s.PushFront(V);
		}
		return (NATIVE[tag].f)(s), idx;

	} else {

		fmt.Printf("\n***Error, undefined: %s\n", tag)

	}
	return ToyVal{t: tnone}, idx + 1;
}

func doTags(tags []string) {
	doTag(tags, 0);
}

func printGLOB() {
	for idx := range GLOB {
		var v = GLOB[idx];
		fmt.Printf("GLOB %s: %d, %s, %d\n", idx, v.t, v.s, v.i)
	}
}

func toyEval (code string) {
	var tags = strings.Split(code, " ", 0);
	doTags(tags);
}

func regNative (name string, pc int, f func(s *list.List)ToyVal) {
	NATIVE[name] = Native{name, pc, f};
}

func natives() {
	regNative("inc", 1, 
		func (s *list.List) ToyVal { 
			var v1 = s.Front().Value.(ToyVal);
			return ToyVal{i: v1.i + 1 }; });
	regNative("add", 2, 
		func (s *list.List) ToyVal { 
			var e1 = s.Front();
			var v1 = e1.Value.(ToyVal);
			var v2 ToyVal = e1.Next().Value.(ToyVal);
			return ToyVal{i: v1.i + v2.i }; });
	regNative("join", 2, 
		func (s *list.List) ToyVal { 
			var e1 = s.Front();
			var v1 = e1.Value.(ToyVal);
			var v2 ToyVal = e1.Next().Value.(ToyVal);
			return ToyVal{s: v1.s + v2.s }; });
}

func main() {
	println("~ Toylingo interpreter v0.0001 ~ ");
	natives();
	var in *bufio.Reader = bufio.NewReader(os.Stdin); 
	var inp string;
	for (strings.Trim(inp, "\n\r") != "exit!") {
		print(">> ");
		inp, _ := in.ReadBytes('\n')
		toyEval(string(inp));
		printGLOB();
	}
}
*/
