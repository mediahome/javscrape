package rule

import (
	"fmt"
	"reflect"
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
								BeforeURL: "test url",
								FromValue: []string{"test from"},
								AfterURL:  "test uri",
								Relative:  true,
								Selector:  "",
								Success:   nil,
							},
							Through:   false,
							OnSuccess: "detail",
							OnFailure: "",
							Success: []Process{
								{
									Name:     "nexturl",
									Trim:     true,
									Type:     "put",
									Property: "attr",
								},
							},
						},
						{
							Type:  "",
							Name:  "detail",
							Index: 2,
							Web: Web{
								Method:    "GET",
								Header:    nil,
								FromValue: []string{"nexturl"},
								//URL:       "test url",
								//URI:       "test uri",
								Selector: "",
								Success:  nil,
							},
							Through:   false,
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
							Through:   true,
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

func TestLoadRuleFromFile(t *testing.T) {
	type args struct {
		file string
	}
	tests := []struct {
		name    string
		args    args
		want    *Rule
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "",
			args: args{
				file: "tmp.toml",
			},
			want: &Rule{
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
						},
						Through:   false,
						OnSuccess: "",
						OnFailure: "",
					},
					{
						Type:  "",
						Name:  "",
						Index: 2,
						Web: Web{
							Method: "GET",
							Header: map[string][]string{
								"cookie": {"2"},
							},
						},
						Through:   false,
						OnSuccess: "",
						OnFailure: "",
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LoadRuleFromFile(tt.args.file)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadRuleFromFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Printf("json:%+v", got)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadRuleFromFile() got = %v, want %v", got, tt.want)
			}
		})
	}
}
