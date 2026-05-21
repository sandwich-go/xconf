package xconf

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// targetForIgnoreTest 用于 IgnoreFields 测试的最小目标 struct
type targetForIgnoreTest struct {
	Name string `xconf:"name"`
	Age  int    `xconf:"age"`
}

// nestedSubForIgnoreTest 嵌套子结构，用于精细路径粒度测试
type nestedSubForIgnoreTest struct {
	A string `xconf:"a"`
}

// nestedTargetForIgnoreTest 含嵌套子 struct 的目标
type nestedTargetForIgnoreTest struct {
	Name string                  `xconf:"name"`
	Sub  *nestedSubForIgnoreTest `xconf:"sub"`
}

// resetForTest 创建一个不依赖全局 FlagSet/Environ 的最小 XConf
func parseYAMLForTest(t *testing.T, yaml string, target interface{}, opts ...Option) error {
	t.Helper()
	allOpts := append([]Option{
		WithReaders(bytes.NewReader([]byte(yaml))),
		WithFlagSet(nil),
		WithEnviron(),
		WithErrorHandling(ContinueOnError),
	}, opts...)
	return New(allOpts...).Parse(target)
}

func TestMatchFieldPath(t *testing.T) {
	cases := []struct {
		name     string
		path     string
		patterns []string
		want     bool
	}{
		{"empty path", "", []string{"a"}, false},
		{"empty patterns", "a", nil, false},
		{"exact match", "plugin.foo", []string{"plugin.foo"}, true},
		{"prefix: pattern is path's ancestor (subtree)", "plugin.foo.bar", []string{"plugin.foo"}, true},
		{"prefix: pattern is path's ancestor deep", "plugin.foo.x.y.z", []string{"plugin.foo"}, true},
		{"prefix: path is pattern's ancestor (top-level shadow)", "plugin", []string{"plugin.foo"}, true},
		{"prefix: path is pattern's ancestor 2 levels", "plugin", []string{"plugin.foo.bar"}, true},
		{"non-prefix sibling", "plugin.foobar", []string{"plugin.foo"}, false},
		{"non-prefix sibling reversed", "pluginx", []string{"plugin"}, false},
		{"too short", "plugin.fo", []string{"plugin.foo"}, false},
		{"empty pattern in list ignored", "plugin.foo", []string{"", "plugin.foo"}, true},
		{"multiple patterns", "csharp.gen_csharp_proto_db", []string{"a.b", "csharp.gen_csharp_proto_db"}, true},
		{"no match in multiple", "a.b.c", []string{"x", "y.z"}, false},
		{"top level exact", "foo", []string{"foo"}, true},
		{"top level prefix subtree", "foo.bar", []string{"foo"}, true},
		{"top level non-match sibling", "foobar", []string{"foo"}, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := matchFieldPath(c.path, c.patterns)
			if got != c.want {
				t.Errorf("matchFieldPath(%q, %v) = %v, want %v", c.path, c.patterns, got, c.want)
			}
		})
	}
}

func TestToStringSliceLoose(t *testing.T) {
	cases := []struct {
		name string
		in   interface{}
		want []string
	}{
		{"nil", nil, nil},
		{"empty string", "", nil},
		{"single string", "a.b", []string{"a.b"}},
		{"[]string", []string{"a", "", "b"}, []string{"a", "b"}},
		{"[]interface{}", []interface{}{"a", "", "b", 123}, []string{"a", "b"}},
		{"unsupported type", 42, nil},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := toStringSliceLoose(c.in)
			if !slicesEqual(got, c.want) {
				t.Errorf("toStringSliceLoose(%v) = %v, want %v", c.in, got, c.want)
			}
		})
	}
}

func TestMergeIgnoreFields(t *testing.T) {
	cases := []struct {
		name     string
		base     []string
		newOnes  []string
		want     []string
	}{
		{"both nil", nil, nil, nil},
		{"base only", []string{"a"}, nil, []string{"a"}},
		{"new only", nil, []string{"a"}, []string{"a"}},
		{"merge dedupe", []string{"a", "b"}, []string{"b", "c"}, []string{"a", "b", "c"}},
		{"order preserved base first", []string{"b", "a"}, []string{"c"}, []string{"b", "a", "c"}},
		{"empty strings filtered", []string{"a", ""}, []string{"", "b"}, []string{"a", "b"}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := mergeIgnoreFields(c.base, c.newOnes)
			if !slicesEqual(got, c.want) {
				t.Errorf("mergeIgnoreFields(%v, %v) = %v, want %v", c.base, c.newOnes, got, c.want)
			}
		})
	}
}

func TestIgnoreFields_Basic(t *testing.T) {
	yaml := `
name: alice
age: 30
plugin_x: value_for_plugin
xconf_ignore_fields:
  - plugin_x
`
	var target targetForIgnoreTest
	err := parseYAMLForTest(t, yaml, &target, WithErrorUnused(true))
	if err != nil {
		t.Fatalf("expected no error with whitelist, got: %v", err)
	}
	if target.Name != "alice" || target.Age != 30 {
		t.Errorf("target not parsed correctly: %+v", target)
	}
}

func TestIgnoreFields_PrefixMatch_KnownParent(t *testing.T) {
	// 父字段 sub 在 struct 中存在，子字段不存在 → mapstructure 会下钻并报告精确路径 sub.unknown_child
	yaml := `
name: alice
sub:
  a: real_value
  unknown_child:
    nested:
      deep: 1
xconf_ignore_fields:
  - sub.unknown_child
`
	var target nestedTargetForIgnoreTest
	err := parseYAMLForTest(t, yaml, &target, WithErrorUnused(true))
	if err != nil {
		t.Fatalf("expected no error with prefix whitelist, got: %v", err)
	}
	if target.Sub == nil || target.Sub.A != "real_value" {
		t.Errorf("nested sub not parsed correctly: %+v", target)
	}
}

func TestIgnoreFields_PrefixMatch_UnknownTopLevel(t *testing.T) {
	// 顶层字段在 struct 中不存在 → mapstructure 仅报告顶层 key
	// 白名单写更细的路径（如 plugin.foo）也应放行（path 是 pattern 祖先的规则）
	yaml := `
name: alice
age: 30
plugin:
  foo:
    bar: 1
    baz:
      x: 2
xconf_ignore_fields:
  - plugin.foo
`
	var target targetForIgnoreTest
	err := parseYAMLForTest(t, yaml, &target, WithErrorUnused(true))
	if err != nil {
		t.Fatalf("expected no error with whitelist 'plugin.foo' covering top-level 'plugin', got: %v", err)
	}
}

func TestIgnoreFields_PrefixMatch_UnknownTopLevelExact(t *testing.T) {
	// 顶层字段在 struct 中不存在 → 白名单写顶层精确名
	yaml := `
name: alice
age: 30
plugin:
  foo:
    bar: 1
xconf_ignore_fields:
  - plugin
`
	var target targetForIgnoreTest
	err := parseYAMLForTest(t, yaml, &target, WithErrorUnused(true))
	if err != nil {
		t.Fatalf("expected no error with top-level whitelist, got: %v", err)
	}
}

func TestIgnoreFields_BoundarySafety_KnownParent(t *testing.T) {
	// 父字段 sub 已知，白名单 sub.foo 不应误命中 sub.foobar
	yaml := `
name: alice
sub:
  a: real
  foobar: oops
xconf_ignore_fields:
  - sub.foo
`
	var target nestedTargetForIgnoreTest
	err := parseYAMLForTest(t, yaml, &target, WithErrorUnused(true))
	if err == nil {
		t.Fatalf("expected error: 'sub.foobar' should NOT be matched by 'sub.foo'")
	}
	if !strings.Contains(err.Error(), "sub.foobar") {
		t.Errorf("expected error to mention sub.foobar, got: %v", err)
	}
}

func TestIgnoreFields_BoundarySafety_TopLevel(t *testing.T) {
	// 白名单 plugin 不应误命中 pluginx
	yaml := `
name: alice
age: 30
pluginx:
  oops: 1
xconf_ignore_fields:
  - plugin
`
	var target targetForIgnoreTest
	err := parseYAMLForTest(t, yaml, &target, WithErrorUnused(true))
	if err == nil {
		t.Fatalf("expected error: 'pluginx' should NOT be matched by 'plugin'")
	}
	if !strings.Contains(err.Error(), "pluginx") {
		t.Errorf("expected error to mention pluginx, got: %v", err)
	}
}

func TestIgnoreFields_DefaultBehaviorUnchanged(t *testing.T) {
	// 不写 xconf_ignore_fields 时，未识别字段应当报错（验证不破坏既有行为）
	yaml := `
name: alice
age: 30
unknown_field: oops
`
	var target targetForIgnoreTest
	err := parseYAMLForTest(t, yaml, &target, WithErrorUnused(true))
	if err == nil {
		t.Fatalf("expected error for unknown field without whitelist")
	}
	if !strings.Contains(err.Error(), "unknown_field") {
		t.Errorf("expected error to mention unknown_field, got: %v", err)
	}
}

func TestIgnoreFields_ErrorUnusedFalseStillSilent(t *testing.T) {
	// ErrorUnused=false 时无论是否配白名单都不报错（白名单不应改变这种语义）
	yaml := `
name: alice
unknown_field: oops
xconf_ignore_fields:
  - some_other
`
	var target targetForIgnoreTest
	err := parseYAMLForTest(t, yaml, &target, WithErrorUnused(false))
	if err != nil {
		t.Fatalf("expected no error when ErrorUnused=false, got: %v", err)
	}
}

func TestIgnoreFields_MultipleWhitelistEntries(t *testing.T) {
	yaml := `
name: alice
age: 30
csharp:
  gen_csharp_proto_db: foo
my_plugin:
  any:
    nested: value
xconf_ignore_fields:
  - csharp.gen_csharp_proto_db
  - my_plugin
`
	var target targetForIgnoreTest
	err := parseYAMLForTest(t, yaml, &target, WithErrorUnused(true))
	if err != nil {
		t.Fatalf("expected no error with multiple whitelist entries, got: %v", err)
	}
}

func TestIgnoreFields_CoexistWithFieldPathRemoved(t *testing.T) {
	// 白名单和 FieldPathRemoved 可以共存：白名单静默丢弃；FieldPathRemoved 告警丢弃
	yaml := `
name: alice
age: 30
old_field: keep_for_compat
plugin_x: ok
xconf_ignore_fields:
  - plugin_x
`
	var warnings []string
	logCapture := func(s string) { warnings = append(warnings, s) }
	var target targetForIgnoreTest
	err := parseYAMLForTest(t, yaml, &target,
		WithErrorUnused(true),
		WithFieldPathRemoved("old_field"),
		WithLogWarning(logCapture),
	)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	// 白名单字段不应触发告警；只有 FieldPathRemoved 命中的会触发告警
	hasDeprecatedWarning := false
	for _, w := range warnings {
		if strings.Contains(w, "DEPRECATED") && strings.Contains(w, "old_field") {
			hasDeprecatedWarning = true
		}
		if strings.Contains(w, "DEPRECATED") && strings.Contains(w, "plugin_x") {
			t.Errorf("plugin_x should be silenced by whitelist, but got DEPRECATED warning: %s", w)
		}
	}
	if !hasDeprecatedWarning {
		t.Errorf("expected DEPRECATED warning for old_field, got warnings: %v", warnings)
	}
}

func TestIgnoreFields_MultiFileInheritUnion(t *testing.T) {
	// 多文件继承场景：base.yaml 与 inherit.yaml 各自定义白名单 → 并集生效
	tmpDir, err := ioutil.TempDir("", "xconf_ignore_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	inheritFile := filepath.Join(tmpDir, "inherit.yaml")
	baseFile := filepath.Join(tmpDir, "base.yaml")

	inheritContent := []byte(`
inherit_only_field: from_inherit
xconf_ignore_fields:
  - inherit_only_field
`)
	baseContent := []byte(`
xconf_inherit_files:
  - inherit.yaml
name: alice
base_only_field: from_base
xconf_ignore_fields:
  - base_only_field
`)
	if err := ioutil.WriteFile(inheritFile, inheritContent, 0644); err != nil {
		t.Fatal(err)
	}
	if err := ioutil.WriteFile(baseFile, baseContent, 0644); err != nil {
		t.Fatal(err)
	}

	var target targetForIgnoreTest
	err = New(
		WithFiles(baseFile),
		WithFlagSet(nil),
		WithEnviron(),
		WithErrorHandling(ContinueOnError),
		WithErrorUnused(true),
	).Parse(&target)

	if err != nil {
		t.Fatalf("expected no error with merged whitelist, got: %v", err)
	}
	if target.Name != "alice" {
		t.Errorf("target not parsed correctly: %+v", target)
	}
}

// slicesEqual 简单切片相等比较（用于测试，避免引入新依赖）
func slicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
