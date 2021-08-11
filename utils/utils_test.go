package utils

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"testing"
)

func TestHash(t *testing.T) {
	hash := "e005c1d727f7776a57a661d61a182816d8953c0432780beeae35e337830b1746"
	s := struct{ Test string }{Test: "test"}
	t.Run("Hash is always same", func(t *testing.T) {
		x := Hash(s)
		if x != hash {
			t.Errorf("Expected %s, got %s", hash, x)
		}
	})
	t.Run("Hash is hex encoded", func(t *testing.T) {
		x := Hash(s)
		_, err := hex.DecodeString(x)
		if err != nil {
			t.Error("Hash should be hex encoded")
		}
	})
}

func ExampleHash() {
	s := struct{ Test string }{Test: "test"}
	x := Hash(s)
	fmt.Println(x)
	// Output: e005c1d727f7776a57a661d61a182816d8953c0432780beeae35e337830b1746
}

func TestToBytes(t *testing.T) {
	s := "test"
	b := ToBytes(s)
	k := reflect.TypeOf(b).Kind()
	if k != reflect.Slice {
		t.Errorf("ToBytes should return a slice of bytes got %s", k)
	}
}

func TestSplitter(t *testing.T) {
	type test struct {
		input  string
		sep    string
		index  int
		output string
	}
	tests := []test{
		{input: "0:6:0", sep: ":", index: 1, output: "6"},
		{input: "0:6:0", sep: ":", index: 10, output: ""},
		{input: "0:6:0", sep: "/", index: 0, output: "0:6:0"},
	}
	for _, tc := range tests {
		got := Splitter(tc.input, tc.sep, tc.index)
		if got != tc.output {
			t.Errorf("Expected %s and got %s", tc.output, got)
		}
	}
}

func TestHadleErr(t *testing.T) {
	oldLogFn := logFn
	defer func() {
		logFn = oldLogFn
	}()
	called := false
	logFn = func(v ...interface{}) {
		called = true
	}
	err := errors.New("test")
	HandleErr(err)
	if !called {
		t.Error("HandleError should call in")
	}
}

func TestFromBytes(t *testing.T) {
	type testStruct struct {
		Test string
	}
	var restored testStruct
	ts := testStruct{"test"}
	b := ToBytes(ts)
	FromBytes(&restored, b)
	if !reflect.DeepEqual(ts, restored) {
		t.Error("FromBytes() should restore struct.")
	}
}

func TestToJSON(t *testing.T) {
	type testStruct struct{ Test string }
	s := testStruct{"test"}
	b := ToJSON(s)
	k := reflect.TypeOf(b).Kind()
	if k != reflect.Slice {
		t.Errorf("Expected %v and got %v", reflect.Slice, k)
	}
	var restored testStruct
	json.Unmarshal(b, &restored)
	if !reflect.DeepEqual(s, restored) {
		t.Error("ToJSON() should encode to JSON correctly.")
	}
}
