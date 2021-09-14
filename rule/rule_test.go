package rule

import (
	"testing"
)

func TestSaveRuleToFile(t *testing.T) {
	type args struct {
		file string
		r    *Rule
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "",
			args: args{
				file: "tmp.toml",
				r: &Rule{
					Entrance: "",
					MainPage: "",
					Actions: []Action{
						{
							Type:  "",
							Name:  "",
							Index: 0,
							Web: Web{
								Method: "GET",
								Header: map[string][]string{
									"cookie": {"1"},
								},
								Skip: []SkipType{
									SkipTypeInput, SkipTypeMainPage,
								},
								Value:    []string{"test from"},
								Relative: true,
								Selector: "",
								Success:  nil,
							},
							OnSuccess: "detail",
							OnFailure: "",
							Success: []Process{
								{
									Name:     "nexturl",
									Selector: "",
									Cache: map[string]string{
										"selector": "abc",
									},
									Index:         0,
									Trim:          true,
									Type:          "put",
									Property:      "attr",
									PropertyIndex: 0,
									PropertyName:  "",
									Value:         "",
								},
							},
						},
						{
							Type:  "",
							Name:  "detail",
							Index: 2,
							Web: Web{
								Method: "GET",
								Header: nil,
								Value:  []string{"$nexturl"},
								//URL:       "test url",
								//URI:       "test uri",
								Selector: "",
								Success: []Process{
									{
										Name:     "",
										Selector: "",
										Cache: map[string]string{
											"selector": "cache",
										},
										Index:         0,
										Trim:          false,
										Type:          "",
										Property:      "",
										PropertyIndex: 0,
										PropertyName:  "",
										Value:         "",
									},
								},
							},
							OnSuccess: "",
							OnFailure: "",
						},
						{
							Type:  "",
							Name:  "getvalue",
							Index: 3,
							Web:   Web{
								//Method: "GET",
								//Header: map[string][]string{
								//	"cookie": {"1"},
								//},
								//URL: "test url",
								//URI: "test uri",
							},
							OnSuccess: "",
							OnFailure: "",
							Success:   []Process{},
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SaveRuleToFile(tt.args.file, tt.args.r); (err != nil) != tt.wantErr {
				t.Errorf("SaveRuleToFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
