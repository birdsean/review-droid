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

const TEST_DIFF = `diff --git a/README.md b/README.md
index a4d3592..53adbf0 100644
--- a/README.md
+++ b/README.md
@@ -5,12 +5,14 @@ This is a simple app that will utilize AI to do an initial pass on your code rev
 ## Features
 
 ### MVP
-- [ ] GitHub OAuth: Implement the OAuth flow to allow users to authenticate with GitHub and authorize the plugin to access their repositories.
-- [ ] Webhook Listener: Implement the functionality to listen for pull request events from GitHub.
-- [ ] Code Retrieval: Implement the functionality to fetch the code changes from the pull request.
+- [ ] Code Retrieval: Implement the functionality to fetch the code changes from the pull request. Use env vars for permissions.
 - [ ] Language Model Integration: Integrate with OpenAI API.
 - [ ] Code Analysis: Feed the code changes retrieved from the pull request into the LLM. Extract comments or suggestions related to code quality, potential bugs, or best practices from the LLM's output.
-[] Comment Posting: Develop the functionality to post the extracted comments and suggestions as comments directly on the pull request within the GitHub interface.
+- [ ] Comment Posting: Develop the functionality to post the extracted comments and suggestions as comments directly on the pull request within the GitHub interface.
+
+### Backlog
+- [ ] GitHub App OAuth: Implement the OAuth flow to allow users to authenticate with GitHub and authorize the plugin to access their repositories. Implement as a GitHub App.
+- [ ] Webhook Listener: Implement the functionality to listen for pull request events from GitHub.
 
 ## Entities
 * TBD
\ No newline at end of file
diff --git a/go.mod b/go.mod
new file mode 100644
index 0000000..4fd5e32
--- /dev/null
+++ b/go.mod
@@ -0,0 +1,7 @@
+module github.com/birdsean/review-droid
+
+go 1.20
+
+require (
+       github.com/golang/protobuf v1.5.2 // indirect
+)
diff --git a/go.sum b/go.sum
new file mode 100644
index 0000000..e9c1376
--- /dev/null
+++ b/go.sum
@@ -0,0 +1,1 @@
+github.com/golang/protobuf v1.3.1/go.mod h1:6lQm79b+lXiMfvg/cZm0SGofjICqVBUtrP5yJMmIC1U=
diff --git a/main.go b/main.go
new file mode 100644
index 0000000..aab42e0
--- /dev/null
+++ b/main.go
@@ -0,0 +1,6 @@
+package main
+
+
+func main() {
+       // code here
+}`

func TestDiffTransformer_generateSegments(t *testing.T) {
	type fields struct {
		rawDiff   string
		fileDiffs []string
		segments  []string
	}
	tests := []struct {
		name   string
		fields fields
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dt := &DiffTransformer{
				rawDiff:   tt.fields.rawDiff,
				fileDiffs: tt.fields.fileDiffs,
				segments:  tt.fields.segments,
			}
			dt.generateSegments()
		})
	}
}

func TestDiffTransformer_splitIntoFiles(t *testing.T) {
	type fields struct {
		rawDiff   string
		fileDiffs []string
		segments  []string
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			"splitIntoFiles",
			fields{
				TEST_DIFF,
				[]string{},
				[]string{},
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
			dt.splitIntoFiles()
			if len(dt.fileDiffs) != 4 {
				t.Errorf("len(DiffTransformer.splitIntoFiles()) = %v, want %v", len(dt.fileDiffs), 4)
			}

			// make sure each file starts with "a/" or "b/"
			for _, file := range dt.fileDiffs {
				if !strings.HasPrefix(file, "a/") && !strings.HasPrefix(file, "b/") {
					firstTwo := file[:2]
					t.Errorf("DiffTransformer.splitIntoFiles() = %v, want %v", firstTwo, "a/ or b/")
				}
			}
		})
	}
}
