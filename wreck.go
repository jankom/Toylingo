package main

import strings "strings"
//import fmt "fmt"
import bufio "bufio"
import os "os"
//import strconv "strconv"
import list "container/list"

// then make a Val interface then types Int, Float, String (with getType getVal and toString) funcs) (look at go how this is done)
// make the reader so that it pases whole block char by char and creates the Ents (entities)


const (
	t_none	= iota
	t_str = iota
	t_int	= iota
	t_native = iota
	t_word	= iota
);

type Val interface {
	gst() string
	giv() int
	gsv() string
};

const (
	u_int = iota
	u_str = iota
	u_word = iota
	u_setword = iota
	u_getword = iota
	u_block = iota
)

type Unit struct {
	t int
	s string
}

func NewUnit(t int, s string) *Unit {
	r := new(Unit)
	r.t = t
	r.s = s
	return r
}

type TNative struct {
	n	string
	pc	int
	f	func(s *list.List) Val
};

func (V *TNative) gst() string { return "Native" }
func (V *TNative) gsv() string { return V.n }
func (V *TNative) giv() int { return t_none }

type TWord struct {
	v string
}

func (V *TWord) gst() string { return "Word" }
func (V *TWord) gsv() string { return V.v }
func (V *TWord) giv() int { return t_none }

type TString struct {
	v string
}

func (V *TString) gst() string { return "String" }
func (V *TString) gsv() string { return V.v }
func (V *TString) giv() int { return t_none }

type TInt struct {
	v int
}

func (V *TInt) gst() string { return "Int" }
func (V *TInt) gsv() string { return "" }
func (V *TInt) giv() int { return V.v }

const (
	p_null = iota
	p_unit = iota
	p_str = iota
	p_block = iota
)

type Parser struct {
	state int
	units *list.List
	curr string
	currt int
	currd int
	blockstr bool
}

func NewParser () *Parser {
	var r = new(Parser)
	r.reset()
	return r
}

func (r *Parser) reset () {
	r.state = 0
	r.units = list.New()
	r.curr = ""
	r.currt = 0
	r.currd = 0
}

func (r *Parser) doBlock (code string, i int) {
	var ch string = string(code[i]);
	switch r.state {
		case p_null:
			switch ch {
				case " ", "\t":
					
				case "[":
					r.currt = u_block
					r.state = p_block
					r.currd += 1
				case "\"":
					r.currt = u_str
					r.state = p_str
				case "0","1","2","3","4","5","6","7","8","9":
					r.currt = u_int
					r.state = p_unit
					r.curr += ch
				case ":":
					r.currt = u_getword
					r.state = p_unit
				default: 
					r.currt = u_word
					r.state = p_unit
					r.curr += ch
			}
		case p_unit:
			switch ch {
				case " ", "\t":
					r.units.PushBack(NewUnit(r.currt, r.curr))
					r.curr = ""
					r.state = p_null
				default: 
					r.curr += ch
			}
		case p_block:
			switch ch {
				case "\"":
					r.blockstr = !r.blockstr
					r.curr += ch
				case "[":
					if (!r.blockstr) {
						r.currd += 1
					}
					r.curr += ch
				case "]":
					r.currd -= 1
					if r.currd == 0 {
						r.units.PushBack(NewUnit(r.currt, r.curr))
						r.curr = ""
						r.state = p_null
					} else {
						r.curr += ch
					}
				default: 
					r.curr += ch
			}
		case p_str:
			switch ch {
				case "\"":
					r.units.PushBack(NewUnit(r.currt, r.curr))
				default: 
					r.curr += ch
			}
	}
	if i < len(code) - 1 {
		r.doBlock(code, i + 1)		
	} else {
		switch r.state {
			case p_null:
			case p_unit:
						r.units.PushBack(NewUnit(r.currt, r.curr))
						r.curr = ""
						r.state = p_null
			case p_str:
		}
	}
}

func (r *Parser) print() {
	print("units:\n")
	for u := r.units.Front(); u != nil; u = u.Next() {
		u.Value.(*Unit).print()
	}
}

func (u *Unit) print() {
	var t = ""
	switch u.t { 
		case 0: t = "Int" 
		case 1: t = "Str" 
		case 2: t = "Word" 
		case 3: t = "SetWord" 
		case 4: t = "GetWord" 
		case 5: t = "Block" 
	}
	print(">> " + t + ": " + u.s + "\n")
}


type Context struct {
	words *map[string]Val
}


func main() {
	println("~ reck interpreter v0.0002 ~ ");
	var in *bufio.Reader = bufio.NewReader(os.Stdin); 
	var inp string;
	var parser = NewParser()
	for (strings.Trim(inp, "\n\r") != "exit!") {
		print(">> ");
		inp, _ := in.ReadBytes('\n')
		parser.doBlock(string(inp), 0);
		parser.print()
		parser.reset()
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
	return ToyVal{t: t_none}, idx + 1;
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
