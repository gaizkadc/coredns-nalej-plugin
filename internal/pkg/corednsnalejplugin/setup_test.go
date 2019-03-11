package corednsnalejplugin

import (
	"reflect"
	"testing"

	"github.com/coredns/coredns/plugin/pkg/response"

	"github.com/mholt/caddy"
)

func TestLogParse(t *testing.T) {
	tests := []struct {
		inputLogRules    string
		shouldErr        bool
		expectedLogRules []Rule
	}{
		{`corednsnalejplugin`, false, []Rule{{
			NameScope: ".",
			Format:    DefaultLogFormat,
			Class:     map[response.Class]struct{}{response.All: struct{}{}},
		}}},
		{`corednsnalejplugin example.org`, false, []Rule{{
			NameScope: "example.org.",
			Format:    DefaultLogFormat,
			Class:     map[response.Class]struct{}{response.All: struct{}{}},
		}}},
		{`corednsnalejplugin example.org. {common}`, false, []Rule{{
			NameScope: "example.org.",
			Format:    CommonLogFormat,
			Class:     map[response.Class]struct{}{response.All: struct{}{}},
		}}},
		{`corednsnalejplugin example.org {combined}`, false, []Rule{{
			NameScope: "example.org.",
			Format:    CombinedLogFormat,
			Class:     map[response.Class]struct{}{response.All: struct{}{}},
		}}},
		{`corednsnalejplugin example.org.
		corednsnalejplugin example.net {combined}`, false, []Rule{{
			NameScope: "example.org.",
			Format:    DefaultLogFormat,
			Class:     map[response.Class]struct{}{response.All: struct{}{}},
		}, {
			NameScope: "example.net.",
			Format:    CombinedLogFormat,
			Class:     map[response.Class]struct{}{response.All: struct{}{}},
		}}},
		{`corednsnalejplugin example.org {host}
			  corednsnalejplugin example.org {when}`, false, []Rule{{
			NameScope: "example.org.",
			Format:    "{host}",
			Class:     map[response.Class]struct{}{response.All: struct{}{}},
		}, {
			NameScope: "example.org.",
			Format:    "{when}",
			Class:     map[response.Class]struct{}{response.All: struct{}{}},
		}}},
		{`corednsnalejplugin example.org example.net`, false, []Rule{{
			NameScope: "example.org.",
			Format:    DefaultLogFormat,
			Class:     map[response.Class]struct{}{response.All: struct{}{}},
		}, {
			NameScope: "example.net.",
			Format:    DefaultLogFormat,
			Class:     map[response.Class]struct{}{response.All: struct{}{}},
		}}},
		{`corednsnalejplugin example.org example.net {host}`, false, []Rule{{
			NameScope: "example.org.",
			Format:    "{host}",
			Class:     map[response.Class]struct{}{response.All: struct{}{}},
		}, {
			NameScope: "example.net.",
			Format:    "{host}",
			Class:     map[response.Class]struct{}{response.All: struct{}{}},
		}}},
		{`corednsnalejplugin example.org example.net {when} {
			class denial
		}`, false, []Rule{{
			NameScope: "example.org.",
			Format:    "{when}",
			Class:     map[response.Class]struct{}{response.Denial: struct{}{}},
		}, {
			NameScope: "example.net.",
			Format:    "{when}",
			Class:     map[response.Class]struct{}{response.Denial: struct{}{}},
		}}},

		{`corednsnalejplugin example.org {
				class all
			}`, false, []Rule{{
			NameScope: "example.org.",
			Format:    CommonLogFormat,
			Class:     map[response.Class]struct{}{response.All: struct{}{}},
		}}},
		{`corednsnalejplugin example.org {
			class denial
		}`, false, []Rule{{
			NameScope: "example.org.",
			Format:    CommonLogFormat,
			Class:     map[response.Class]struct{}{response.Denial: struct{}{}},
		}}},
		{`corednsnalejplugin {
			class denial
		}`, false, []Rule{{
			NameScope: ".",
			Format:    CommonLogFormat,
			Class:     map[response.Class]struct{}{response.Denial: struct{}{}},
		}}},
		{`corednsnalejplugin {
			class denial error
		}`, false, []Rule{{
			NameScope: ".",
			Format:    CommonLogFormat,
			Class:     map[response.Class]struct{}{response.Denial: struct{}{}, response.Error: struct{}{}},
		}}},
		{`corednsnalejplugin {
			class denial
			class error
		}`, false, []Rule{{
			NameScope: ".",
			Format:    CommonLogFormat,
			Class:     map[response.Class]struct{}{response.Denial: struct{}{}, response.Error: struct{}{}},
		}}},
		{`corednsnalejplugin {
			class abracadabra
		}`, true, []Rule{}},
		{`corednsnalejplugin {
			class
		}`, true, []Rule{}},
		{`corednsnalejplugin {
			unknown
		}`, true, []Rule{}},
	}
	for i, test := range tests {
		c := caddy.NewTestController("dns", test.inputLogRules)
		actualLogRules, err := logParse(c)

		if err == nil && test.shouldErr {
			t.Errorf("Test %d with input '%s' didn't error, but it should have", i, test.inputLogRules)
		} else if err != nil && !test.shouldErr {
			t.Errorf("Test %d with input '%s' errored, but it shouldn't have; got '%v'",
				i, test.inputLogRules, err)
		}
		if len(actualLogRules) != len(test.expectedLogRules) {
			t.Fatalf("Test %d expected %d no of Log rules, but got %d ",
				i, len(test.expectedLogRules), len(actualLogRules))
		}
		for j, actualLogRule := range actualLogRules {

			if actualLogRule.NameScope != test.expectedLogRules[j].NameScope {
				t.Errorf("Test %d expected %dth LogRule NameScope for '%s' to be  %s  , but got %s",
					i, j, test.inputLogRules, test.expectedLogRules[j].NameScope, actualLogRule.NameScope)
			}

			if actualLogRule.Format != test.expectedLogRules[j].Format {
				t.Errorf("Test %d expected %dth LogRule Format for '%s' to be  %s  , but got %s",
					i, j, test.inputLogRules, test.expectedLogRules[j].Format, actualLogRule.Format)
			}

			if !reflect.DeepEqual(actualLogRule.Class, test.expectedLogRules[j].Class) {
				t.Errorf("Test %d expected %dth LogRule Class to be  %v  , but got %v",
					i, j, test.expectedLogRules[j].Class, actualLogRule.Class)
			}
		}
	}

}
