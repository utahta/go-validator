package validator

import (
	"context"
	"testing"
)

func Test_tagParse(t *testing.T) {
	testcases := []struct {
		rawTag string
		want   tagChunk
	}{
		{
			rawTag: "required",
			want: tagChunk{
				Tags: []Tag{{name: "required"}},
				Next: nil,
			},
		},
		{
			rawTag: "len(3)",
			want: tagChunk{
				Tags: []Tag{{name: "len", params: []string{"3"}}},
				Next: nil,
			},
		},
		{
			rawTag: "len(1|3)",
			want: tagChunk{
				Tags: []Tag{{name: "len", params: []string{"1", "3"}}},
				Next: nil,
			},
		},
		{
			rawTag: "tmp(text/plain;charset=UTF-8)",
			want: tagChunk{
				Tags: []Tag{{name: "tmp", params: []string{"text/plain;charset=UTF-8"}}},
				Next: nil,
			},
		},
		{
			rawTag: "tmp(a|b|c)",
			want: tagChunk{
				Tags: []Tag{{name: "tmp", params: []string{"a", "b", "c"}}},
				Next: nil,
			},
		},
		{
			rawTag: "tmp(a\\|b\\|c)",
			want: tagChunk{
				Tags: []Tag{{name: "tmp", params: []string{"a|b|c"}}},
				Next: nil,
			},
		},
		{
			rawTag: "tmp(a\\/b\\/c)",
			want: tagChunk{
				Tags: []Tag{{name: "tmp", params: []string{"a\\/b\\/c"}}},
				Next: nil,
			},
		},
		{
			rawTag: "tmp(a,b,c)",
			want: tagChunk{
				Tags: []Tag{{name: "tmp", params: []string{"a,b,c"}}},
				Next: nil,
			},
		},
		{
			rawTag: "tmp(a b　c\t)",
			want: tagChunk{
				Tags: []Tag{{name: "tmp", params: []string{"a b　c\t"}}},
				Next: nil,
			},
		},
		{
			rawTag: "tmp(あいうえお)",
			want: tagChunk{
				Tags: []Tag{{name: "tmp", params: []string{"あいうえお"}}},
				Next: nil,
			},
		},
		{
			rawTag: "tmp((a,b,c))",
			want: tagChunk{
				Tags: []Tag{{name: "tmp", params: []string{"(a,b,c)"}}},
				Next: nil,
			},
		},
		{
			rawTag: "tmp(a\nb\nc)",
			want: tagChunk{
				Tags: []Tag{{name: "tmp", params: []string{"a\nb\nc"}}},
				Next: nil,
			},
		},
		{
			rawTag: "len(1|3),len(AAA|BBB|CCC)",
			want: tagChunk{
				Tags: []Tag{
					{name: "len", params: []string{"1", "3"}},
					{name: "len", params: []string{"AAA", "BBB", "CCC"}},
				},
				Next: nil,
			},
		},
		{
			rawTag: "required, len(3)",
			want: tagChunk{
				Tags: []Tag{
					{name: "required"},
					{name: "len", params: []string{"3"}},
				},
				Next: nil,
			},
		},
		{
			rawTag: "required, len(3), len(1|3)",
			want: tagChunk{
				Tags: []Tag{
					{name: "required"},
					{name: "len", params: []string{"3"}},
					{name: "len", params: []string{"1", "3"}},
				},
				Next: nil,
			},
		},

		// next separator
		{
			rawTag: "required;",
			want: tagChunk{
				Tags: []Tag{{name: "required"}},
				Next: &tagChunk{},
			},
		},
		{
			rawTag: "required ; ",
			want: tagChunk{
				Tags: []Tag{{name: "required"}},
				Next: &tagChunk{},
			},
		},
		{
			rawTag: "required; required",
			want: tagChunk{
				Tags: []Tag{{name: "required"}},
				Next: &tagChunk{
					Tags: []Tag{{name: "required"}},
				},
			},
		},
		{
			rawTag: "required ; len(3)",
			want: tagChunk{
				Tags: []Tag{{name: "required"}},
				Next: &tagChunk{
					Tags: []Tag{{name: "len", params: []string{"3"}}},
				},
			},
		},
		{
			rawTag: "len(3); required",
			want: tagChunk{
				Tags: []Tag{{name: "len", params: []string{"3"}}},
				Next: &tagChunk{
					Tags: []Tag{{name: "required"}},
				},
			},
		},
		{
			rawTag: "; len(3)",
			want: tagChunk{
				Tags: []Tag{},
				Next: &tagChunk{
					Tags: []Tag{{name: "len", params: []string{"3"}}},
				},
			},
		},

		// Optional
		{
			rawTag: "optional",
			want: tagChunk{
				Tags:     []Tag{},
				Optional: true,
				Next:     nil,
			},
		},
		{
			rawTag: "optional,required",
			want: tagChunk{
				Tags:     []Tag{{name: "required"}},
				Optional: true,
				Next:     nil,
			},
		},
		{
			rawTag: "required,optional",
			want: tagChunk{
				Tags:     []Tag{{name: "required"}},
				Optional: true,
				Next:     nil,
			},
		},
		{
			rawTag: "len(3),optional,required",
			want: tagChunk{
				Tags: []Tag{
					{name: "len", params: []string{"3"}},
					{name: "required"},
				},
				Optional: true,
				Next:     nil,
			},
		},

		// Optional, next separator
		{
			rawTag: "optional,max(3); required,len(3)",
			want: tagChunk{
				Tags: []Tag{
					{name: "max", params: []string{"3"}},
				},
				Optional: true,
				Next: &tagChunk{
					Tags: []Tag{
						{name: "required"},
						{name: "len", params: []string{"3"}},
					},
				},
			},
		},
		{
			rawTag: "max(3),optional; required,len(3)",
			want: tagChunk{
				Tags: []Tag{
					{name: "max", params: []string{"3"}},
				},
				Optional: true,
				Next: &tagChunk{
					Tags: []Tag{
						{name: "required"},
						{name: "len", params: []string{"3"}},
					},
				},
			},
		},
		{
			rawTag: "max(3); optional,required,len(3)",
			want: tagChunk{
				Tags: []Tag{
					{name: "max", params: []string{"3"}},
				},
				Next: &tagChunk{
					Tags: []Tag{
						{name: "required"},
						{name: "len", params: []string{"3"}},
					},
					Optional: true,
				},
			},
		},
		{
			rawTag: "max(3); required,len(3),optional",
			want: tagChunk{
				Tags: []Tag{
					{name: "max", params: []string{"3"}},
				},
				Next: &tagChunk{
					Tags: []Tag{
						{name: "required"},
						{name: "len", params: []string{"3"}},
					},
					Optional: true,
				},
			},
		},
		{
			rawTag: "optional; required",
			want: tagChunk{
				Tags:     []Tag{},
				Optional: true,
				Next: &tagChunk{
					Tags: []Tag{{name: "required"}},
				},
			},
		},

		// OR
		{
			rawTag: "alpha|numeric",
			want: tagChunk{
				Tags: []Tag{{name: "or", params: []string{"alpha", "numeric"}}},
				Next: nil,
			},
		},
		{
			rawTag: "alpha|numeric|len(1|10)",
			want: tagChunk{
				Tags: []Tag{{name: "or", params: []string{"alpha", "numeric", "len(1|10)"}}},
				Next: nil,
			},
		},
		{
			rawTag: "optional|alpha|numeric",
			want: tagChunk{
				Tags:     []Tag{{name: "or", params: []string{"alpha", "numeric"}}},
				Optional: true,
				Next:     nil,
			},
		},
		{
			rawTag: "or(req|numeric)",
			want: tagChunk{
				Tags: []Tag{{name: "or", params: []string{"req", "numeric"}}},
				Next: nil,
			},
		},

		// OR, AND
		{
			rawTag: "alpha|numeric,len(1|10)",
			want: tagChunk{
				Tags: []Tag{
					{name: "or", params: []string{"alpha", "numeric"}},
					{name: "len", params: []string{"1", "10"}},
				},
				Next: nil,
			},
		},
		{
			rawTag: "alpha|numeric,min(1),max(10)",
			want: tagChunk{
				Tags: []Tag{
					{name: "or", params: []string{"alpha", "numeric"}},
					{name: "min", params: []string{"1"}},
					{name: "max", params: []string{"10"}},
				},
				Next: nil,
			},
		},
		{
			rawTag: "alpha|numeric,min(1)|max(10)",
			want: tagChunk{
				Tags: []Tag{
					{name: "or", params: []string{"alpha", "numeric"}},
					{name: "or", params: []string{"min(1)", "max(10)"}},
				},
				Next: nil,
			},
		},

		// OR, next separator
		{
			rawTag: "alpha|numeric ;",
			want: tagChunk{
				Tags: []Tag{
					{name: "or", params: []string{"alpha", "numeric"}},
				},
				Next: &tagChunk{},
			},
		},
		{
			rawTag: "alpha|numeric ; min(1)|max(10)",
			want: tagChunk{
				Tags: []Tag{{name: "or", params: []string{"alpha", "numeric"}}},
				Next: &tagChunk{
					Tags: []Tag{{name: "or", params: []string{"min(1)", "max(10)"}}},
				},
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.rawTag, func(t *testing.T) {
			v := New()
			v.Apply(WithFunc("tmp", func(context.Context, Field, FuncOption) (bool, error) { return true, nil }))

			chunk, err := v.parseTag(tc.rawTag)
			if err != nil {
				t.Fatal(err)
			}

			if len(tc.want.Tags) != len(chunk.Tags) {
				t.Fatalf("want tag len %v, but got %v", len(tc.want.Tags), len(chunk.Tags))
			}
			if tc.want.Optional != chunk.Optional {
				t.Fatalf("want optional %v, but got %v", tc.want.Optional, chunk.Optional)
			}
			for i, wantTag := range tc.want.Tags {
				if wantTag.name != chunk.Tags[i].name {
					t.Errorf("want tag name %v, but got %v", wantTag.name, chunk.Tags[i].name)
				}

				if len(wantTag.params) != len(chunk.Tags[i].params) {
					t.Errorf("want tag params len %v, but got %v", len(wantTag.params), len(chunk.Tags[i].params))
				}

				for j := range wantTag.params {
					if wantTag.params[j] != chunk.Tags[i].params[j] {
						t.Errorf("want tag params[%d] %v, but got %v", j, wantTag.params[j], chunk.Tags[i].params[j])
					}
				}
			}

			if tc.want.Next == nil {
				if chunk.Next != nil {
					t.Fatalf("want next is nil, but got %v", chunk.Next)
				}
				return
			}
			if len(tc.want.Next.Tags) != len(chunk.Next.Tags) {
				t.Fatalf("want next tag len %v, but got %v", len(tc.want.Next.Tags), len(chunk.Next.Tags))
			}
			if tc.want.Next.Optional != chunk.Next.Optional {
				t.Fatalf("want next optional %v, but got %v", tc.want.Next.Optional, chunk.Next.Optional)
			}
			for i, nextTag := range tc.want.Next.Tags {
				if nextTag.name != chunk.Next.Tags[i].name {
					t.Errorf("want tag name %v, but got %v", nextTag.name, chunk.Next.Tags[i].name)
				}

				if len(nextTag.params) != len(chunk.Next.Tags[i].params) {
					t.Errorf("want tag params len %v, but got %v", len(nextTag.params), len(chunk.Next.Tags[i].params))
				}

				for j := range nextTag.params {
					if nextTag.params[j] != chunk.Next.Tags[i].params[j] {
						t.Errorf("want tag params[%d] %v, but got %v", j, nextTag.params[j], chunk.Next.Tags[i].params[j])
					}
				}
			}
		})
	}
}

func Test_tagParseInvalid(t *testing.T) {
	testcases := []struct {
		rawTag    string
		wantError string
	}{
		{
			rawTag:    "req,,",
			wantError: "parse: invalid literal in tag separator",
		},
		{
			rawTag:    "req||",
			wantError: "parse: invalid literal in or separator",
		},
		{
			rawTag:    "unknown",
			wantError: "parse: tag unknown function not found",
		},
		{
			rawTag:    "unknown,req",
			wantError: "parse: tag unknown function not found",
		},
		{
			rawTag:    "unknown;req",
			wantError: "parse: tag unknown function not found",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.rawTag, func(t *testing.T) {
			_, err := New().parseTag(tc.rawTag)
			if err == nil {
				t.Fatal("want error, but got nil")
			}
			if err.Error() != tc.wantError {
				t.Errorf("want `%v`, got `%v`", tc.wantError, err.Error())
			}
		})
	}
}

func Test_tagCache(t *testing.T) {
	const rawTag = "required,min(1),max(10)"
	v := New()
	want, err := v.parseTag(rawTag)
	if err != nil {
		t.Fatal(err)
	}

	got, ok := v.tagCache.Load(rawTag)
	if !ok {
		t.Fatal("want load true, got false")
	}

	if len(want.Tags) != len(got.Tags) {
		t.Fatalf("want tag len %v, but got %v", len(want.Tags), len(got.Tags))
	}
	if want.Optional != got.Optional {
		t.Fatalf("want optional %v, but got %v", want.Optional, got.Optional)
	}
	for i, wantTag := range want.Tags {
		if wantTag.name != got.Tags[i].name {
			t.Errorf("want tag name %v, but got %v", wantTag.name, got.Tags[i].name)
		}

		if len(wantTag.params) != len(got.Tags[i].params) {
			t.Errorf("want tag params len %v, but got %v", len(wantTag.params), len(got.Tags[i].params))
		}

		for j := range wantTag.params {
			if wantTag.params[j] != got.Tags[i].params[j] {
				t.Errorf("want tag params[%d] %v, but got %v", j, wantTag.params[j], got.Tags[i].params[j])
			}
		}
	}

	if want.Next == nil {
		if got.Next != nil {
			t.Fatalf("want next nil, but got %v", got)
		}
	}
}
