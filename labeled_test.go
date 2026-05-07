package main

import (
	"html/template"
	"testing"

	"github.com/Lord-lami/render-html"
)

func TestRenderLabeled(t *testing.T) {
	type record struct {
		name string
		data any
		want template.HTML
	}
	test := record{}
	tests := []record{}
	labeledRenderFuncs := []render.RenderFunc{}

	// Test 1 - map test
	test.name = "Chirstmas"
	test.data = labeledMap(map[string][]render.DateString{"Next 3": {"25-12-2026", "25-12-2027", "25-12-2028"}})
	test.want = "<strong>Chirstmas: </strong><table class=\"Chirstmas\">\n    <thead>\n        <tr class=\"keys\"><th scope=\"col\"><span class=\"0\">Next 3</span>\n</th></tr>\n    </thead>\n    <tbody>\n        <tr class=\"values\"><td><ul class=\"0\">\n    <li><time class=\"0\" datetime=\"2026-12-25\">Fri, 25 Dec 2026</time>\n</li><li><time class=\"1\" datetime=\"2027-12-25\">Sat, 25 Dec 2027</time>\n</li><li><time class=\"2\" datetime=\"2028-12-25\">Mon, 25 Dec 2028</time>\n</li>\n</ul>\n</td></tr>\n    </tbody>\n</table>"
	tests = append(tests, test)
	labeledRenderFuncs = append(labeledRenderFuncs, renderLabeledMap)

	// Test 2 - string slice test
	test.name = "Messages"
	test.data = labeledStringSlice{"Merry Christmas", "Ugly Easter"}
	test.want = "<strong>Messages:</strong><ul class=\"Messages\">\n    <li><span class=\"0\">Merry Christmas</span>\n</li><li><span class=\"1\">Ugly Easter</span>\n</li>\n</ul>\n"
	tests = append(tests, test)
	labeledRenderFuncs = append(labeledRenderFuncs, renderLabeledStringSlice)

	// Test 3 - int test
	test.name = "Code"
	test.data = labeledInt(200)
	test.want = "<strong>Code: </strong><span class=\"Code\">200</span>\n"
	tests = append(tests, test)
	labeledRenderFuncs = append(labeledRenderFuncs, renderLabeledInt)

	// Test 4 - DateString test
	test.name = "Death Day"
	test.data = labeledDateString("21-02-2098")
	test.want = "<strong>Death  Day: </strong><time class=\"Death Day\" datetime=\"2098-02-21\">Fri, 21 Feb 2098</time>\n"
	tests = append(tests, test)
	labeledRenderFuncs = append(labeledRenderFuncs, renderLabeledDateString)

	for i, test := range tests {
		result := labeledRenderFuncs[i](test.name, test.data)
		if result != test.want {
			t.Errorf("expected %q got %q", test.want, result)
		}
	}
}

func TestSpaceOutCamelCase(t *testing.T) {
	type record struct {
		str  string
		want string
	}
	test := record{}

	test.str = "CamelCase"
	test.want = "Camel Case"

	result := spaceOutCamelCase(test.str)
	if result != test.want {
		t.Errorf("expected %q got %q", test.want, result)
	}
}
