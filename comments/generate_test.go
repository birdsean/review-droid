package comments

import (
	"reflect"
	"testing"
)

func Test_generateComment(t *testing.T) {
	type args struct {
		rawComment   string
		originalCode string
		filename     string
	}
	tests := []struct {
		name string
		args args
		want *Comment
	}{
		{
			name: "clean response",
			args: args{
				"[- Line 8] Potential Bug: Looks like the entire test function has been removed. Is this intentional?",
				"8 - func TestDiffTransformer_numberLines(t *testing.T) {",
				"test.go",
			},
			want: &Comment{
				CodeLine:    8,
				CommentBody: "Potential Bug: Looks like the entire test function has been removed. Is this intentional?",
				FileAddress: "test.go",
				Side:        "LEFT",
			},
		},
		{
			name: "range of lines",
			args: args{
				"[- Line 8-14] Potential Bug: Looks like the entire test function has been removed. Is this intentional?",
				"8 - func TestDiffTransformer_numberLines(t *testing.T) {",
				"test.go",
			},
			want: &Comment{
				CodeLine:    8,
				CommentBody: "Potential Bug: Looks like the entire test function has been removed. Is this intentional?",
				FileAddress: "test.go",
				Side:        "LEFT",
			},
		},
		{
			name: "missing plus or minus",
			args: args{
				"[Line 8] Suggestion: Consider renaming `numberedRawDiff` to `numberedDiff` for simplicity",
				"8 - func TestDiffTransformer_numberLines(t *testing.T) {",
				"test.go",
			},
			want: &Comment{
				CodeLine:    8,
				CommentBody: "Suggestion: Consider renaming `numberedRawDiff` to `numberedDiff` for simplicity",
				FileAddress: "test.go",
				Side:        "LEFT",
			},
		},
		{
			name: "detects plus correctly",
			args: args{
				"[+ Line 8] Suggestion: Consider renaming `numberedRawDiff` to `numberedDiff` for simplicity",
				"8 + func TestDiffTransformer_numberLines(t *testing.T) {",
				"test.go",
			},
			want: &Comment{
				CodeLine:    8,
				CommentBody: "Suggestion: Consider renaming `numberedRawDiff` to `numberedDiff` for simplicity",
				FileAddress: "test.go",
				Side:        "RIGHT",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := generateComment(tt.args.rawComment, tt.args.originalCode, tt.args.filename, false); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("generateComment() = %+v, want %+v", got, tt.want)
			}
		})
	}
}
