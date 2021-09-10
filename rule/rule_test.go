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
								URL: "test url",
								URI: "test uri",
							},
							Through:   false,
							OnSuccess: "",
							OnFailure: "",
							Success: Process{
								Name:     "nexturl",
								Trim:     true,
								Type:     "put",
								Property: "attr",
							},
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
								URL: "test url",
								URI: "test uri",
							},
							Through:   false,
							OnSuccess: "",
							OnFailure: "",
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
							URL: "test url",
							URI: "test uri",
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
							URL: "test url",
							URI: "test uri",
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
