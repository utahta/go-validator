package validator

import "testing"

func Test_tagParse(t *testing.T) {
	testcases := []struct {
		rawTag string
		want   tagChunk
	}{
		{
			rawTag: "required",
			want: tagChunk{
				Tags: []Tag{{Name: "required"}},
				Next: nil,
			},
		},
		{
			rawTag: "len(3)",
			want: tagChunk{
				Tags: []Tag{{Name: "len", Params: []string{"3"}}},
				Next: nil,
			},
		},
		{
			rawTag: "len(1|3)",
			want: tagChunk{
				Tags: []Tag{{Name: "len", Params: []string{"1", "3"}}},
				Next: nil,
			},
		},
		{
			rawTag: "len(1|3),len(AAA|BBB|CCC)",
			want: tagChunk{
				Tags: []Tag{
					{Name: "len", Params: []string{"1", "3"}},
					{Name: "len", Params: []string{"AAA", "BBB", "CCC"}},
				},
				Next: nil,
			},
		},
		{
			rawTag: "required, len(3)",
			want: tagChunk{
				Tags: []Tag{
					{Name: "required"},
					{Name: "len", Params: []string{"3"}},
				},
				Next: nil,
			},
		},
		{
			rawTag: "required, len(3), len(1|3)",
			want: tagChunk{
				Tags: []Tag{
					{Name: "required"},
					{Name: "len", Params: []string{"3"}},
					{Name: "len", Params: []string{"1", "3"}},
				},
				Next: nil,
			},
		},

		// next separator
		{
			rawTag: "required;",
			want: tagChunk{
				Tags: []Tag{{Name: "required"}},
				Next: &tagChunk{},
			},
		},
		{
			rawTag: "required ; ",
			want: tagChunk{
				Tags: []Tag{{Name: "required"}},
				Next: &tagChunk{},
			},
		},
		{
			rawTag: "required; required",
			want: tagChunk{
				Tags: []Tag{{Name: "required"}},
				Next: &tagChunk{
					Tags: []Tag{{Name: "required"}},
				},
			},
		},
		{
			rawTag: "required ; len(3)",
			want: tagChunk{
				Tags: []Tag{{Name: "required"}},
				Next: &tagChunk{
					Tags: []Tag{{Name: "len", Params: []string{"3"}}},
				},
			},
		},
		{
			rawTag: "len(3); required",
			want: tagChunk{
				Tags: []Tag{{Name: "len", Params: []string{"3"}}},
				Next: &tagChunk{
					Tags: []Tag{{Name: "required"}},
				},
			},
		},
		{
			rawTag: "; len(3)",
			want: tagChunk{
				Tags: []Tag{},
				Next: &tagChunk{
					Tags: []Tag{{Name: "len", Params: []string{"3"}}},
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
				Tags:     []Tag{{Name: "required"}},
				Optional: true,
				Next:     nil,
			},
		},
		{
			rawTag: "required,optional",
			want: tagChunk{
				Tags:     []Tag{{Name: "required"}},
				Optional: true,
				Next:     nil,
			},
		},
		{
			rawTag: "len(3),optional,required",
			want: tagChunk{
				Tags: []Tag{
					{Name: "len", Params: []string{"3"}},
					{Name: "required"},
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
					{Name: "max", Params: []string{"3"}},
				},
				Optional: true,
				Next: &tagChunk{
					Tags: []Tag{
						{Name: "required"},
						{Name: "len", Params: []string{"3"}},
					},
				},
			},
		},
		{
			rawTag: "max(3),optional; required,len(3)",
			want: tagChunk{
				Tags: []Tag{
					{Name: "max", Params: []string{"3"}},
				},
				Optional: true,
				Next: &tagChunk{
					Tags: []Tag{
						{Name: "required"},
						{Name: "len", Params: []string{"3"}},
					},
				},
			},
		},
		{
			rawTag: "max(3); optional,required,len(3)",
			want: tagChunk{
				Tags: []Tag{
					{Name: "max", Params: []string{"3"}},
				},
				Next: &tagChunk{
					Tags: []Tag{
						{Name: "required"},
						{Name: "len", Params: []string{"3"}},
					},
					Optional: true,
				},
			},
		},
		{
			rawTag: "max(3); required,len(3),optional",
			want: tagChunk{
				Tags: []Tag{
					{Name: "max", Params: []string{"3"}},
				},
				Next: &tagChunk{
					Tags: []Tag{
						{Name: "required"},
						{Name: "len", Params: []string{"3"}},
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
					Tags: []Tag{{Name: "required"}},
				},
			},
		},

		// OR
		{
			rawTag: "alpha|numeric",
			want: tagChunk{
				Tags: []Tag{{Name: "or", Params: []string{"alpha", "numeric"}}},
				Next: nil,
			},
		},
		{
			rawTag: "alpha|numeric|len(1|10)",
			want: tagChunk{
				Tags: []Tag{{Name: "or", Params: []string{"alpha", "numeric", "len(1|10)"}}},
				Next: nil,
			},
		},
		{
			rawTag: "optional|alpha|numeric",
			want: tagChunk{
				Tags:     []Tag{{Name: "or", Params: []string{"alpha", "numeric"}}},
				Optional: true,
				Next:     nil,
			},
		},
		{
			rawTag: "or(req|numeric)",
			want: tagChunk{
				Tags: []Tag{{Name: "or", Params: []string{"req", "numeric"}}},
				Next: nil,
			},
		},

		// OR, AND
		{
			rawTag: "alpha|numeric,len(1|10)",
			want: tagChunk{
				Tags: []Tag{
					{Name: "or", Params: []string{"alpha", "numeric"}},
					{Name: "len", Params: []string{"1", "10"}},
				},
				Next: nil,
			},
		},
		{
			rawTag: "alpha|numeric,min(1),max(10)",
			want: tagChunk{
				Tags: []Tag{
					{Name: "or", Params: []string{"alpha", "numeric"}},
					{Name: "min", Params: []string{"1"}},
					{Name: "max", Params: []string{"10"}},
				},
				Next: nil,
			},
		},
		{
			rawTag: "alpha|numeric,min(1)|max(10)",
			want: tagChunk{
				Tags: []Tag{
					{Name: "or", Params: []string{"alpha", "numeric"}},
					{Name: "or", Params: []string{"min(1)", "max(10)"}},
				},
				Next: nil,
			},
		},

		// OR, next separator
		{
			rawTag: "alpha|numeric ;",
			want: tagChunk{
				Tags: []Tag{
					{Name: "or", Params: []string{"alpha", "numeric"}},
				},
				Next: &tagChunk{},
			},
		},
		{
			rawTag: "alpha|numeric ; min(1)|max(10)",
			want: tagChunk{
				Tags: []Tag{{Name: "or", Params: []string{"alpha", "numeric"}}},
				Next: &tagChunk{
					Tags: []Tag{{Name: "or", Params: []string{"min(1)", "max(10)"}}},
				},
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.rawTag, func(t *testing.T) {
			chunk, err := New().tagParse(tc.rawTag)
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
				if wantTag.Name != chunk.Tags[i].Name {
					t.Errorf("want tag name %v, but got %v", wantTag.Name, chunk.Tags[i].Name)
				}

				if len(wantTag.Params) != len(chunk.Tags[i].Params) {
					t.Errorf("want tag params len %v, but got %v", len(wantTag.Params), len(chunk.Tags[i].Params))
				}

				for j := range wantTag.Params {
					if wantTag.Params[j] != chunk.Tags[i].Params[j] {
						t.Errorf("want tag params[%d] %v, but got %v", j, wantTag.Params[j], chunk.Tags[i].Params[j])
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
				if nextTag.Name != chunk.Next.Tags[i].Name {
					t.Errorf("want tag name %v, but got %v", nextTag.Name, chunk.Next.Tags[i].Name)
				}

				if len(nextTag.Params) != len(chunk.Next.Tags[i].Params) {
					t.Errorf("want tag params len %v, but got %v", len(nextTag.Params), len(chunk.Next.Tags[i].Params))
				}

				for j := range nextTag.Params {
					if nextTag.Params[j] != chunk.Next.Tags[i].Params[j] {
						t.Errorf("want tag params[%d] %v, but got %v", j, nextTag.Params[j], chunk.Next.Tags[i].Params[j])
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
			rawTag:    "len(1,2)",
			wantError: "parse: failed to new tag",
		},
		{
			rawTag:    "unknown",
			wantError: "parse: tag unknown function not found",
		},
	}

	for _, tc := range testcases {
		_, err := New().tagParse(tc.rawTag)
		if err.Error() != tc.wantError {
			t.Errorf("want `%v`, got `%v`", tc.wantError, err.Error())
		}
	}
}

func Test_tagCache(t *testing.T) {
	const rawTag = "required,min(1),max(10)"
	v := New()
	want, err := v.tagParse(rawTag)
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
		if wantTag.Name != got.Tags[i].Name {
			t.Errorf("want tag name %v, but got %v", wantTag.Name, got.Tags[i].Name)
		}

		if len(wantTag.Params) != len(got.Tags[i].Params) {
			t.Errorf("want tag params len %v, but got %v", len(wantTag.Params), len(got.Tags[i].Params))
		}

		for j := range wantTag.Params {
			if wantTag.Params[j] != got.Tags[i].Params[j] {
				t.Errorf("want tag params[%d] %v, but got %v", j, wantTag.Params[j], got.Tags[i].Params[j])
			}
		}
	}

	if want.Next == nil {
		if got.Next != nil {
			t.Fatalf("want next nil, but got %v", got)
		}
	}
}
