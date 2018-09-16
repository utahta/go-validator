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
				{Name: "required", Enable: true, dig: true},
			},
		},
		{
			"len(3)",
			[]Tag{
				{Name: "len", Params: []string{"3"}, Enable: true, dig: true},
			},
		},
		{
			"len(1|3)",
			[]Tag{
				{Name: "len", Params: []string{"1", "3"}, Enable: true, dig: true},
			},
		},
		{
			"len(1|3),len(AAA|BBB|CCC)",
			[]Tag{
				{Name: "len", Params: []string{"1", "3"}, Enable: true, dig: true},
				{Name: "len", Params: []string{"AAA", "BBB", "CCC"}, Enable: true, dig: true},
			},
		},
		{
			"required, len(3)",
			[]Tag{
				{Name: "required", Enable: true, dig: true},
				{Name: "len", Params: []string{"3"}, Enable: true, dig: true},
			},
		},
		{
			"required, len(3), len(1|3)",
			[]Tag{
				{Name: "required", Enable: true, dig: true},
				{Name: "len", Params: []string{"3"}, Enable: true, dig: true},
				{Name: "len", Params: []string{"1", "3"}, Enable: true, dig: true},
			},
		},
		{
			"required;",
			[]Tag{
				{Name: "required", Enable: true, dig: false},
			},
		},
		{
			"required ; ",
			[]Tag{
				{Name: "required", Enable: true, dig: false},
			},
		},
		{
			"required; required",
			[]Tag{
				{Name: "required", Enable: true, dig: false},
				{Name: "required", Enable: false, dig: true},
			},
		},
		{
			"required ; len(3)",
			[]Tag{
				{Name: "required", Enable: true, dig: false},
				{Name: "len", Params: []string{"3"}, Enable: false, dig: true},
			},
		},
		{
			"len(3); required",
			[]Tag{
				{Name: "len", Params: []string{"3"}, Enable: true, dig: false},
				{Name: "required", Enable: false, dig: true},
			},
		},
		{
			"; len(3)",
			[]Tag{
				{Name: "len", Params: []string{"3"}, Enable: false, dig: true},
			},
		},
		{
			"optional,required",
			[]Tag{
				{Name: "required", Enable: true, dig: true, Optional: true},
			},
		},
		{
			"len(3),optional,required",
			[]Tag{
				{Name: "len", Params: []string{"3"}, Enable: true, dig: true, Optional: true},
				{Name: "required", Enable: true, dig: true, Optional: true},
			},
		},
		{
			"alpha|numeric",
			[]Tag{
				{Name: "or", Params: []string{"alpha", "numeric"}, Enable: true, dig: true},
			},
		},
		{
			"alpha|numeric ;",
			[]Tag{
				{Name: "or", Params: []string{"alpha", "numeric"}, Enable: true, dig: false},
			},
		},
		{
			"alpha|numeric|len(1|10)",
			[]Tag{
				{Name: "or", Params: []string{"alpha", "numeric", "len(1|10)"}, Enable: true, dig: true},
			},
		},
		{
			"alpha|numeric,len(1|10)",
			[]Tag{
				{Name: "or", Params: []string{"alpha", "numeric"}, Enable: true, dig: true},
				{Name: "len", Params: []string{"1", "10"}, Enable: true, dig: true},
			},
		},
		{
			"alpha|numeric,min(1),max(10)",
			[]Tag{
				{Name: "or", Params: []string{"alpha", "numeric"}, Enable: true, dig: true},
				{Name: "min", Params: []string{"1"}, Enable: true, dig: true},
				{Name: "max", Params: []string{"10"}, Enable: true, dig: true},
			},
		},
		{
			"alpha|numeric,min(1)|max(10)",
			[]Tag{
				{Name: "or", Params: []string{"alpha", "numeric"}, Enable: true, dig: true},
				{Name: "or", Params: []string{"min(1)", "max(10)"}, Enable: true, dig: true},
			},
		},
		{
			"alpha|numeric ; min(1)|max(10)",
			[]Tag{
				{Name: "or", Params: []string{"alpha", "numeric"}, Enable: true, dig: false},
				{Name: "or", Params: []string{"min(1)", "max(10)"}, Enable: false, dig: true},
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
					t.Errorf("want valid %v, got %v", tc.want[i].Enable, tags[i].Enable)
				}
				if tc.want[i].dig != tags[i].dig {
					t.Errorf("want dig %v, got %v", tc.want[i].dig, tags[i].dig)
				}
			}
		})
	}
}

func Test_tagCache(t *testing.T) {
	const rawTag = "required,min(1),max(10)"
	want, err := New().tagParse(rawTag)
	if err != nil {
		t.Fatal(err)
	}

	v, ok := cache.Load(rawTag)
	if !ok {
		t.Fatal("want load true, got false")
	}
	got := v.([]Tag)

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
		if want[i].dig != got[i].dig {
			t.Errorf("want dig %v, got %v", want[i].dig, got[i].dig)
		}
	}
}
