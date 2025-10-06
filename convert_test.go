package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestParseCell(t *testing.T) {
	tests := []struct {
		in   string
		want any
	}{
		{"", nil},
		{"TRUE", true},
		{"false", false},
		{"42", int64(42)},
		{"7.25", 7.25},
		{"001", int64(1)},
		{"8.0", 8.0},
	}
	for _, tt := range tests {
		got := parseCell(tt.in)
		if got != tt.want {
			t.Fatalf("parseCell(%q)=%v (%T), want %v (%T)", tt.in, got, got, tt.want, tt.want)
		}
	}
}

func TestConvertCSV_OrderAndNulls(t *testing.T) {
	csv := `value,income,age,rooms,bedrooms,pop,hh
		452600,8.3252,41,880,129,322,126
		352100,7.2574,52,1467,190,496,`
	var out bytes.Buffer
	if err := encodeJSONL(strings.NewReader(csv), &out); err != nil {
		t.Fatalf("encodeJSONL error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 JSON lines, got %d", len(lines))
	}

	wantPrefix := `{"value":452600,"income":8.3252,"age":41,"rooms":880,"bedrooms":129,"pop":322,"hh":126}`
	if lines[0] != wantPrefix {
		t.Fatalf("line 1 mismatch:\n got: %s\nwant: %s", lines[0], wantPrefix)
	}

	if !strings.HasSuffix(lines[1], `"hh":null}`) {
		t.Fatalf("expected null hh, got: %s", lines[1])
	}
}

func TestEncodeJSONL_EmptyHeaderError(t *testing.T) {
	var out bytes.Buffer
	err := encodeJSONL(strings.NewReader(""), &out)
	if err == nil || !strings.Contains(err.Error(), "reading header") {
		t.Fatalf("expected header read error, got: %v", err)
	}
}
