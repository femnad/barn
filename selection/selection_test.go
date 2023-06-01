package selection

import "testing"

func Test_getDisplayName(t *testing.T) {
	type args struct {
		entryName      string
		targetPath     string
		includeParents int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "No parents",
			args: args{
				entryName:      "foo",
				targetPath:     "/bar/baz",
				includeParents: 0,
			},
			want: "foo",
		},
		{
			name: "Non-zero parents",
			args: args{
				entryName:      "foo",
				targetPath:     "/bar/baz/foo",
				includeParents: 1,
			},
			want: "baz/foo",
		},
		{
			name: "Use all parents",
			args: args{
				entryName:      "foo",
				targetPath:     "/bar/baz/foo",
				includeParents: 2,
			},
			want: "/bar/baz/foo",
		},
		{
			name: "Too many parents",
			args: args{
				entryName:      "foo",
				targetPath:     "/bar/baz/foo",
				includeParents: 3,
			},
			want: "/bar/baz/foo",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getDisplayName(tt.args.entryName, tt.args.targetPath, tt.args.includeParents); got != tt.want {
				t.Errorf("getDisplayName() = %v, want %v", got, tt.want)
			}
		})
	}
}
