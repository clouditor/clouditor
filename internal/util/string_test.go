package util

import (
	"reflect"
	"testing"

	"clouditor.io/clouditor/v2/api/ontology"
)

func TestCamelCaseToSnakeCase(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Input camelCase",
			args: args{input: "testCamelCaseString"},
			want: "test_camel_case_string",
		},
		{
			name: "Input snake_case",
			args: args{input: "test_camel_case_string"},
			want: "test_camel_case_string",
		},
		{
			name: "Input empty",
			args: args{input: ""},
			want: "",
		},
		{
			name: "Input variation 1",
			args: args{input: "TESTCamelCaseString"},
			want: "test_camel_case_string",
		},
		{
			name: "Input variation 2",
			args: args{input: "testCamelCaseSTRING"},
			want: "test_camel_case_string",
		},
		{
			name: "Input with digit 1",
			args: args{input: "3TestCamelCaseString"},
			want: "3_test_camel_case_string",
		},
		{
			name: "Input with digit 2",
			args: args{input: "3testCamelCaseString"},
			want: "3test_camel_case_string",
		},
		{
			name: "Input with digit 3",
			args: args{input: "test3CamelCaseString"},
			want: "test3_camel_case_string",
		},
		{
			name: "Input with digit 4",
			args: args{input: "t3CamelCaseString"},
			want: "t3_camel_case_string",
		},
		{
			name: "Input with digit 5",
			args: args{input: "T3CamelCaseString"},
			want: "t_3_camel_case_string",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CamelCaseToSnakeCase(tt.args.input); got != tt.want {
				t.Errorf("CamelCaseToSnakeCase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_marksNewWord(t *testing.T) {
	type args struct {
		i     int
		input []rune
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Correct values",
			args: args{
				i:     4,
				input: []rune("testCamelCaseString"),
			},
			want: true,
		},
		{
			name: "Index higher than length of string",
			args: args{
				i:     20,
				input: []rune("testCamelCaseString"),
			},
			want: false,
		},
		{
			name: "Index equals 0",
			args: args{
				i:     0,
				input: []rune("testCamelCaseString"),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := marksNewWord(tt.args.i, tt.args.input); got != tt.want {
				t.Errorf("marksNewWord() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestListResourceIDs(t *testing.T) {
	type args struct {
		r []ontology.IsResource
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "Empty input",
			args: args{},
			want: []string{},
		},
		{
			name: "Happy path",
			args: args{
				[]ontology.IsResource{
					&ontology.Account{Id: "test"},
					&ontology.Account{Id: "test2"},
				},
			},
			want: []string{"test", "test2"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ListResourceIDs(tt.args.r); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ListResourceIDs() = %v, want %v", got, tt.want)
			}
		})
	}
}
