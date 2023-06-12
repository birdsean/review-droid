package transformer

import (
	"strings"
	"testing"
)

func TestDiffTransformer_numberLines(t *testing.T) {
	type fields struct {
		rawDiff   string
		fileDiffs []string
		segments  []string
	}
	type args struct {
		segments []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []string
	}{
		{
			"numberLines",
			fields{
				"",
				[]string{},
				[]string{},
			},
			args{
				[]string{
					"line 1",
					"line 2",
					"line 3",
				},
			},
			[]string{
				"0 line 1",
				"1 line 2",
				"2 line 3",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dt := &DiffTransformer{
				rawDiff:   tt.fields.rawDiff,
				fileDiffs: tt.fields.fileDiffs,
				segments:  tt.fields.segments,
			}
			got := strings.Join(dt.numberLines(tt.args.segments), ",")
			want := strings.Join(tt.want, ",")
			if got != want {
				t.Errorf("DiffTransformer.numberLines() = %v, want %v", got, want)
			}
		})
	}
}
