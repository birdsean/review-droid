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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := generateComment(tt.args.rawComment, tt.args.originalCode, tt.args.filename); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("generateComment() = %+v, want %+v", got, tt.want)
			}
		})
	}
}
