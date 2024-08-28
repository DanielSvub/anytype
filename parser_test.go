package anytype_test

import (
	"os"
	"testing"

	"github.com/DanielSvub/anytype"
)

func TestParser(t *testing.T) {

	t.Run("object", func(t *testing.T) {
		o := Object(
			"object", Object(
				"\"tes\u0074\n\"", 0,
				"innerList", List(0, 1, Object(), "2"),
			),
			"list", List(nil, false, 42, 1.6e-8, "\"Žř@./'á?\u0041\n\\\"", Object("\"1\"", "\"1\""), List(0)),
			"string", "test",
			"bool", true,
			"int", 1,
			"float", 3.14,
			"nil", nil,
		)
		parsed, err := anytype.ParseObject(o.String())
		if err != nil {
			t.Error("JSON parser failed")
		}
		if !parsed.Equals(o) {
			t.Error("JSON parsing or exporting of object does not work properly")
		}
	})

	t.Run("list", func(t *testing.T) {
		l := List(nil, false, 42, 1.6, "\"Žř@./'á?\u0041\n\\\"", List(0, "1"), Object(
			"bool", true,
			"int", 1,
			"float", 3.14,
			"innerList", List(0, 1, Object(), "2"),
			"nil", nil,
			"string", "test",
		))
		parsed, err := anytype.ParseList(l.String())
		if err != nil {
			t.Error("JSON parser failed")
		}
		if !parsed.Equals(l) {
			t.Error("JSON parsing or exporting of list does not work properly")
		}
	})

	t.Run("indent", func(t *testing.T) {
		l := List(1, 2)
		if l.FormatString(2) != "[\n  1,\n  2\n]" {
			t.Error("JSON formatted export does not work properly")
		}
		o := Object(
			"key", "value",
		)
		if o.FormatString(2) != "{\n  \"key\": \"value\"\n}" {
			t.Error("JSON formatted export does not work properly")
		}
	})

	t.Run("errors", func(t *testing.T) {
		if _, err := anytype.ParseObject(""); err == nil {
			t.Error("parser did not return expected error")
		}
		if _, err := anytype.ParseObject("{\"test\":1"); err == nil {
			t.Error("parser did not return expected error")
		}
		if _, err := anytype.ParseObject("{\"test\"1"); err == nil {
			t.Error("parser did not return expected error")
		}
		if _, err := anytype.ParseObject("{1:2}"); err == nil {
			t.Error("parser did not return expected error")
		}
		if _, err := anytype.ParseObject("{\"test\":[]2}"); err == nil {
			t.Error("parser did not return expected error")
		}
		if _, err := anytype.ParseList(""); err == nil {
			t.Error("parser did not return expected error")
		}
		if _, err := anytype.ParseList("[1,2"); err == nil {
			t.Error("parser did not return expected error")
		}
		if _, err := anytype.ParseList("[test]"); err == nil {
			t.Error("parser did not return expected error")
		}
		if _, err := anytype.ParseObject("{\xa0"); err == nil {
			t.Error("parser did not return expected error")
		}
		if _, err := anytype.ParseList("[\xa0"); err == nil {
			t.Error("parser did not return expected error")
		}
		if _, err := anytype.ParseObject("{\"test\":a}"); err == nil {
			t.Error("parser did not return expected error")
		}
		if _, err := anytype.ParseList("[a]"); err == nil {
			t.Error("parser did not return expected error")
		}
		if _, err := anytype.ParseObject("{\"test\":{\"test\":a}}"); err == nil {
			t.Error("parser did not return expected error")
		}
		if _, err := anytype.ParseObject("{\"test\":[a]}"); err == nil {
			t.Error("parser did not return expected error")
		}
		if _, err := anytype.ParseList("[[a]]"); err == nil {
			t.Error("parser did not return expected error")
		}
		if _, err := anytype.ParseList("[{\"test\":a}]"); err == nil {
			t.Error("parser did not return expected error")
		}
	})

	t.Run("whitespaces", func(t *testing.T) {
		if _, err := anytype.ParseObject("{\"first\" : \"1\" , \n\"second\": \"2\" }"); err != nil {
			t.Error("parser did not handle extra whitespaces")
		}
		if _, err := anytype.ParseList("[ \"1\" , \n \"2\" ]"); err != nil {
			t.Error("parser did not handle extra whitespaces")
		}
		if _, err := anytype.ParseObject("{ \"first\" : {} }"); err != nil {
			t.Error("parser did not handle extra whitespaces")
		}
		if _, err := anytype.ParseList("[ {} ]"); err != nil {
			t.Error("parser did not handle extra whitespaces")
		}
	})

	t.Run("repairs", func(t *testing.T) {
		if _, err := anytype.ParseObject("{\"first\":{}\"second\":2}"); err != nil {
			t.Error("parser did not repair broken JSON")
		}
		if _, err := anytype.ParseObject("{\"first\":1 2}"); err != nil {
			t.Error("parser did not repair broken JSON")
		}
		if _, err := anytype.ParseList("[nu ll,[]2]"); err != nil {
			t.Error("parser did not repair broken JSON")
		}
	})

	t.Run("file", func(t *testing.T) {
		if _, err := anytype.ParseFile("test.json"); err == nil {
			t.Error("opened a file which should not exist")
		}
		if err := os.WriteFile("test.json", []byte("{\"first\":[],\"second\":2}"), 0644); err != nil {
			t.Fatal("unable to create the JSON file")
		}
		if _, err := anytype.ParseFile("test.json"); err != nil {
			t.Error("cannot parse JSON from file")
		}
		if err := os.Remove("test.json"); err != nil {
			t.Fatal("unable to delete the JSON file")
		}
	})

}
