package main

import strings "strings"
import fmt "fmt"
import bufio "bufio"
import os "os"
import strconv "strconv"

var GLOB = make(map[string] Literal, 100);
var NATIVE = make(map[string] Native, 100);

const (
	tnone	= iota
	tstring = iota
	tint	= iota
	tfloat	= iota
	tnative = iota
);

type Literal struct {
	t int
	sv string
	iv int
	fv float
};

type Native struct {
	n	string
	pc	int
	f	func(a Literal, b Literal) Literal
};

func wordDefined(w string) bool {
	return w == "add";
}

func doTag(tags []string, idx int) Literal {
	var tag = strings.Trim(tags[idx], "\n\r");
	print("tag: " ); println(tag);
	var iv, ier = strconv.Atoi(tag);
	var fv, fer = strconv.Atof(tag);
	if strings.HasSuffix(tag, ":") {
		GLOB[tag] = doTag(tags, idx + 1);	
	} else if strings.HasPrefix(tag, "\"") && strings.HasSuffix(tag, "\"") {
		return Literal{t: tstring, sv: tag};
	} else if ier == nil {
		return Literal{t: tint, iv: iv};
	} else if  fer == nil {
		return Literal{t: tfloat, fv: fv};
	} else if wordDefined(tag) {
		//NATIVE[tag] != nil
		if tag == "add"  {
			return (NATIVE["add"].f)(doTag(tags, idx+1), doTag(tags, idx+2));
		}
		return doTag(tags, idx + 1);
	} else {
		println("Error: undefined")
	}
	return Literal{t: tnone};
}

func doTags(tags []string) {
	for idx := range tags {
		var l = doTag(tags, idx);
		println(l.t); 
		println(l.sv); 
	}
}

func printGLOB() {
	for idx := range GLOB {
		var v = GLOB[idx];
		fmt.Printf("GLOB %s: %d, %s, %d, %f\n", idx, v.t, v.sv, v.iv, v.fv )
	}
}

func regoEval (code string) {
	println(code)
	var tags = strings.Split(code, " ", 0);
	doTags(tags);
}

func regNative (name string, pc int, f func(a Literal, b Literal) Literal) {
	NATIVE[name] = Native{name, pc, f};
}

func natives() {
	regNative("add", 2, func (a Literal, b Literal) Literal { return Literal{iv: a.iv + b.iv}; });
}

func main() {
	println("REGO interpreter v0.0001");
	natives();
	var in *bufio.Reader = bufio.NewReader(os.Stdin); 
	var inp string;
	for (inp != "exit\n") {
		println(">>");
		inp, _ := in.ReadBytes('\n')
		regoEval(string(inp));
		printGLOB();
	}
}
