package validator

import "testing"

func Test_tagParse(t *testing.T) {
	testcases := []struct {
		rawTag string
		want   []Tag
	}{
		{
			"required",
			[]Tag{
				{Name: "required", Enable: true, isDig: true},
			},
		},
		{
			"len(3)",
			[]Tag{
				{Name: "len", Params: []string{"3"}, Enable: true, isDig: true},
			},
		},
		{
			"len(1|3)",
			[]Tag{
				{Name: "len", Params: []string{"1", "3"}, Enable: true, isDig: true},
			},
		},
		{
			"len(1|3),len(AAA|BBB|CCC)",
			[]Tag{
				{Name: "len", Params: []string{"1", "3"}, Enable: true, isDig: true},
				{Name: "len", Params: []string{"AAA", "BBB", "CCC"}, Enable: true, isDig: true},
			},
		},
		{
			"required, len(3)",
			[]Tag{
				{Name: "required", Enable: true, isDig: true},
				{Name: "len", Params: []string{"3"}, Enable: true, isDig: true},
			},
		},
		{
			"required, len(3), len(1|3)",
			[]Tag{
				{Name: "required", Enable: true, isDig: true},
				{Name: "len", Params: []string{"3"}, Enable: true, isDig: true},
				{Name: "len", Params: []string{"1", "3"}, Enable: true, isDig: true},
			},
		},

		// Dig
		{
			"required;",
			[]Tag{
				{Name: "required", Enable: true, isDig: false},
			},
		},
		{
			"required ; ",
			[]Tag{
				{Name: "required", Enable: true, isDig: false},
			},
		},
		{
			"required; required",
			[]Tag{
				{Name: "required", Enable: true, isDig: false},
				{Name: "required", Enable: false, isDig: true},
			},
		},
		{
			"required ; len(3)",
			[]Tag{
				{Name: "required", Enable: true, isDig: false},
				{Name: "len", Params: []string{"3"}, Enable: false, isDig: true},
			},
		},
		{
			"len(3); required",
			[]Tag{
				{Name: "len", Params: []string{"3"}, Enable: true, isDig: false},
				{Name: "required", Enable: false, isDig: true},
			},
		},
		{
			"; len(3)",
			[]Tag{
				{Name: "len", Params: []string{"3"}, Enable: false, isDig: true},
			},
		},

		// Optional
		{
			"optional,required",
			[]Tag{
				{Name: "required", Enable: true, isDig: true, Optional: true},
			},
		},
		{
			"required,optional",
			[]Tag{
				{Name: "required", Enable: true, isDig: true, Optional: true},
			},
		},
		{
			"len(3),optional,required",
			[]Tag{
				{Name: "len", Params: []string{"3"}, Enable: true, isDig: true, Optional: true},
				{Name: "required", Enable: true, isDig: true, Optional: true},
			},
		},

		// Optional, Dig
		{
			"optional,max(3); required,len(3)",
			[]Tag{
				{Name: "max", Params: []string{"3"}, Enable: true, isDig: false, Optional: true},
				{Name: "required", Enable: false, isDig: true, Optional: false},
				{Name: "len", Params: []string{"3"}, Enable: false, isDig: true, Optional: false},
			},
		},
		{
			"max(3),optional; required,len(3)",
			[]Tag{
				{Name: "max", Params: []string{"3"}, Enable: true, isDig: false, Optional: true},
				{Name: "required", Enable: false, isDig: true, Optional: false},
				{Name: "len", Params: []string{"3"}, Enable: false, isDig: true, Optional: false},
			},
		},
		{
			"max(3); optional,required,len(3)",
			[]Tag{
				{Name: "max", Params: []string{"3"}, Enable: true, isDig: false, Optional: false},
				{Name: "required", Enable: false, isDig: true, Optional: true},
				{Name: "len", Params: []string{"3"}, Enable: false, isDig: true, Optional: true},
			},
		},
		{
			"max(3); required,len(3),optional",
			[]Tag{
				{Name: "max", Params: []string{"3"}, Enable: true, isDig: false, Optional: false},
				{Name: "required", Enable: false, isDig: true, Optional: true},
				{Name: "len", Params: []string{"3"}, Enable: false, isDig: true, Optional: true},
			},
		},
		{
			"optional; required",
			[]Tag{
				{Name: "required", Enable: false, isDig: true, Optional: false},
			},
		},

		// OR
		{
			"alpha|numeric",
			[]Tag{
				{Name: "or", Params: []string{"alpha", "numeric"}, Enable: true, isDig: true},
			},
		},
		{
			"alpha|numeric|len(1|10)",
			[]Tag{
				{Name: "or", Params: []string{"alpha", "numeric", "len(1|10)"}, Enable: true, isDig: true},
			},
		},
		{
			"optional|alpha|numeric",
			[]Tag{
				{Name: "or", Params: []string{"alpha", "numeric"}, Enable: true, isDig: true, Optional: true},
			},
		},

		// OR, AND
		{
			"alpha|numeric,len(1|10)",
			[]Tag{
				{Name: "or", Params: []string{"alpha", "numeric"}, Enable: true, isDig: true},
				{Name: "len", Params: []string{"1", "10"}, Enable: true, isDig: true},
			},
		},
		{
			"alpha|numeric,min(1),max(10)",
			[]Tag{
				{Name: "or", Params: []string{"alpha", "numeric"}, Enable: true, isDig: true},
				{Name: "min", Params: []string{"1"}, Enable: true, isDig: true},
				{Name: "max", Params: []string{"10"}, Enable: true, isDig: true},
			},
		},
		{
			"alpha|numeric,min(1)|max(10)",
			[]Tag{
				{Name: "or", Params: []string{"alpha", "numeric"}, Enable: true, isDig: true},
				{Name: "or", Params: []string{"min(1)", "max(10)"}, Enable: true, isDig: true},
			},
		},

		// OR, Dig
		{
			"alpha|numeric ;",
			[]Tag{
				{Name: "or", Params: []string{"alpha", "numeric"}, Enable: true, isDig: false},
			},
		},
		{
			"alpha|numeric ; min(1)|max(10)",
			[]Tag{
				{Name: "or", Params: []string{"alpha", "numeric"}, Enable: true, isDig: false},
				{Name: "or", Params: []string{"min(1)", "max(10)"}, Enable: false, isDig: true},
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.rawTag, func(t *testing.T) {
			tags, err := New().tagParse(tc.rawTag)
			if err != nil {
				t.Fatal(err)
			}

			if len(tc.want) != len(tags) {
				t.Fatalf("want length %v, got %v", len(tc.want), len(tags))
			}

			for i := range tc.want {
				if tc.want[i].Name != tags[i].Name {
					t.Errorf("want name %v, got %v", tc.want[i].Name, tags[i].Name)
				}
				if len(tc.want[i].Params) != len(tags[i].Params) {
					t.Fatalf("want params length %v, got %v", len(tc.want[i].Params), len(tags[i].Params))
				}
				for j := range tc.want[i].Params {
					if tc.want[i].Params[j] != tags[i].Params[j] {
						t.Errorf("want params %v, got %v", tc.want[i].Params[j], tags[i].Params[j])
					}
				}
				if tc.want[i].Enable != tags[i].Enable {
					t.Errorf("want enable %v, got %v", tc.want[i].Enable, tags[i].Enable)
				}
				if tc.want[i].isDig != tags[i].isDig {
					t.Errorf("want isDig %v, got %v", tc.want[i].isDig, tags[i].isDig)
				}
				if tc.want[i].Optional != tags[i].Optional {
					t.Errorf("want optional %v, got %v", tc.want[i].Optional, tags[i].Optional)
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

	if len(want) != len(got) {
		t.Fatalf("want len %v, got %v", len(want), len(got))
	}

	for i := range want {
		if want[i].Name != got[i].Name {
			t.Errorf("want name %v, got %v", want[i].Name, got[i].Name)
		}
		if len(want[i].Params) != len(got[i].Params) {
			t.Fatalf("want params length %v, got %v", len(want[i].Params), len(got[i].Params))
		}
		for j := range want[i].Params {
			if want[i].Params[j] != got[i].Params[j] {
				t.Errorf("want params %v, got %v", want[i].Params[j], got[i].Params[j])
			}
		}
		if want[i].Enable != got[i].Enable {
			t.Errorf("want valid %v, got %v", want[i].Enable, got[i].Enable)
		}
		if want[i].isDig != got[i].isDig {
			t.Errorf("want isDig %v, got %v", want[i].isDig, got[i].isDig)
		}
		if want[i].Optional != got[i].Optional {
			t.Errorf("want optional %v, got %v", want[i].Optional, got[i].Optional)
		}
	}
}
