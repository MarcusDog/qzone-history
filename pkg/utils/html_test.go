package utils

import "testing"

func TestExtractStringBetween(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		source string
		start  string
		end    string
		want   string
	}{
		{name: "extracts exact range", source: "prefix html:'<li>内容</li>',opuin suffix", start: "html:'", end: "',opuin", want: "<li>内容</li>"},
		{name: "missing start", source: "',opuin", start: "html:'", end: "',opuin", want: ""},
		{name: "missing end", source: "html:'value", start: "html:'", end: "',opuin", want: ""},
		{name: "end before start", source: "',opuin html:'value", start: "html:'", end: "',opuin", want: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := extractStringBetween(tt.source, tt.start, tt.end); got != tt.want {
				t.Fatalf("extractStringBetween() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestProcessOldHTML(t *testing.T) {
	t.Parallel()

	input := `callback({html:'<li class="f-single">hello\x20world</li>',opuin:'123'})`
	want := `<li class="f-single">hello world</li>`
	if got := ProcessOldHTML(input); got != want {
		t.Fatalf("ProcessOldHTML() = %q, want %q", got, want)
	}
}

func TestProcessOldHTMLDecodesUTF8HexBytes(t *testing.T) {
	t.Parallel()

	input := `callback({html:'<li>\xe4\xb8\xad\xe6\x96\x87 C:\\backup</li>',opuin:'123'})`
	want := `<li>中文 C:\backup</li>`
	if got := ProcessOldHTML(input); got != want {
		t.Fatalf("ProcessOldHTML() = %q, want %q", got, want)
	}
}
