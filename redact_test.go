package xconf

import (
	"bytes"
	"flag"
	"strings"
	"testing"
)

type redactTestConfig struct {
	Token       string `xconf:"token" usage:"api token"`
	AccessKeyID string `xconf:"access_key_id" usage:"access key"`
	GameKey     string `xconf:"game_key" usage:"game key"`
	Password    string `xconf:"password" usage:"password"`
	Normal      string `xconf:"normal" usage:"normal value"`
}

func newRedactTestXConf(t *testing.T, enabled bool) (*XConf, *redactTestConfig) {
	t.Helper()
	cc := &redactTestConfig{
		Token:       "token-value",
		AccessKeyID: "access-key-id-value",
		GameKey:     "game-key-value",
		Password:    "password-value",
		Normal:      "normal-value",
	}
	x := New(
		WithFlagSet(flag.NewFlagSet("redact", flag.ContinueOnError)),
		WithFlagArgs(),
		WithEnviron(),
		WithErrorHandling(ContinueOnError),
		WithSensitiveDataRedaction(enabled),
	)
	if err := x.Parse(cc); err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	return x, cc
}

func TestSensitiveDataRedactionForLogOutput(t *testing.T) {
	x, _ := newRedactTestXConf(t, true)

	var usage bytes.Buffer
	x.UsageToWriter(&usage)
	assertRedacted(t, usage.String())

	var yamlHelp bytes.Buffer
	x.UsageToWriter(&yamlHelp, "--help=yaml")
	assertRedacted(t, yamlHelp.String())
}

func TestExplicitRedactedWriter(t *testing.T) {
	x, cc := newRedactTestXConf(t, true)

	var raw bytes.Buffer
	if err := x.SaveVarToWriter(cc, ConfigTypeYAML, &raw); err != nil {
		t.Fatalf("SaveVarToWriter() error = %v", err)
	}
	assertSensitiveValuesVisible(t, raw.String())
	assertSensitiveValuesVisible(t, string(x.MustSaveToBytes(ConfigTypeYAML)))

	var redacted bytes.Buffer
	if err := x.SaveVarToWriterRedacted(cc, ConfigTypeYAML, &redacted); err != nil {
		t.Fatalf("SaveVarToWriterRedacted() error = %v", err)
	}
	assertRedacted(t, redacted.String())
}

func TestAutomaticLogRedactionCanBeDisabled(t *testing.T) {
	x, cc := newRedactTestXConf(t, false)

	var usage bytes.Buffer
	x.UsageToWriter(&usage)
	assertSensitiveValuesVisible(t, usage.String())

	var yamlHelp bytes.Buffer
	x.UsageToWriter(&yamlHelp, "--help=yaml")
	assertSensitiveValuesVisible(t, yamlHelp.String())

	var output bytes.Buffer
	if err := x.SaveVarToWriterRedacted(cc, ConfigTypeYAML, &output); err != nil {
		t.Fatalf("SaveVarToWriterRedacted() error = %v", err)
	}
	assertRedacted(t, output.String())
}

func TestSensitiveDataKeywordMatching(t *testing.T) {
	x := New()
	for _, fieldPath := range []string{
		"service.api_token",
		"service.client_secret",
		"database.password",
		"oss.access_key_id",
		"platform.game_key",
		"token.signing_keys",
	} {
		if !x.shouldRedactSensitiveData(fieldPath) {
			t.Errorf("shouldRedactSensitiveData(%q) = false, want true", fieldPath)
		}
	}
	for _, fieldPath := range []string{
		"service.timeout",
		"service.keyboard_layout",
		"feature.monkey",
		"token.expiration",
		"token.key_in_header",
	} {
		if x.shouldRedactSensitiveData(fieldPath) {
			t.Errorf("shouldRedactSensitiveData(%q) = true, want false", fieldPath)
		}
	}
}

func assertRedacted(t *testing.T, output string) {
	t.Helper()
	for _, value := range sensitiveTestValues() {
		if strings.Contains(output, value) {
			t.Fatalf("sensitive value %q was not redacted in:\n%s", value, output)
		}
	}
	if !strings.Contains(output, SensitiveDataRedactedValue) {
		t.Fatalf("redacted marker not found in:\n%s", output)
	}
	if !strings.Contains(output, "normal-value") {
		t.Fatalf("normal value should remain visible in:\n%s", output)
	}
}

func assertSensitiveValuesVisible(t *testing.T, output string) {
	t.Helper()
	for _, value := range sensitiveTestValues() {
		if !strings.Contains(output, value) {
			t.Fatalf("sensitive value %q should be visible, got:\n%s", value, output)
		}
	}
}

func sensitiveTestValues() []string {
	return []string{"token-value", "access-key-id-value", "game-key-value", "password-value"}
}
