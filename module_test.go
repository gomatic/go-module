package module

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// The fixtures use neutral identities that share no substring with this
// project's own identity, so the rename command never rewrites these tests when
// it renames the project that ships them.

func TestParse(t *testing.T) {
	t.Parallel()
	tests := []struct {
		wantErr error
		name    string
		remote  Remote
		want    Path
	}{
		{name: "ssh scp form", remote: "git@example.com:org/before.cli.git", want: "example.com/org/before.cli"},
		{name: "https url", remote: "https://example.com/org/before.cli.git", want: "example.com/org/before.cli"},
		{
			name:   "https url without .git",
			remote: "https://example.com/org/after.cli",
			want:   "example.com/org/after.cli",
		},
		{name: "ssh url with userinfo", remote: "ssh://git@example.com/org/tool", want: "example.com/org/tool"},
		{name: "surrounding whitespace", remote: "  git@example.com:org/tool.git\n", want: "example.com/org/tool"},
		{name: "nested group", remote: "git@gitlab.example:group/sub/repo.git", want: "gitlab.example/group/sub/repo"},
		{name: "empty", remote: "", wantErr: ErrInvalidRemote},
		{name: "too few segments", remote: "example.com/onlyone", wantErr: ErrInvalidRemote},
		{name: "empty segment", remote: "https://example.com//repo.git", wantErr: ErrInvalidRemote},
		{name: "contains space", remote: "git@example.com:org/wid get", wantErr: ErrInvalidRemote},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			want, must := assert.New(t), require.New(t)

			got, err := Parse(tt.remote)

			if tt.wantErr != nil {
				must.Error(err)
				want.ErrorIs(err, tt.wantErr)
				return
			}
			must.NoError(err)
			want.Equal(tt.want, got)
		})
	}
}

func TestParseDoesNotLeakCredentials(t *testing.T) {
	t.Parallel()
	want, must := assert.New(t), require.New(t)

	_, err := Parse("https://user:s3cr3t-token@example.com/bad")

	must.Error(err)
	want.ErrorIs(err, ErrInvalidRemote)
	want.NotContains(err.Error(), "s3cr3t-token")
	want.NotContains(err.Error(), "user")
}

func TestPathRepo(t *testing.T) {
	t.Parallel()
	want := assert.New(t)

	want.Equal(Name("before.cli"), Path("example.com/org/before.cli").Repo())
	want.Equal(Name("solo"), Path("solo").Repo())
}

func TestNameIdentifier(t *testing.T) {
	t.Parallel()
	want := assert.New(t)

	want.Equal(Identifier("beforecli"), Name("before.cli").Identifier())
	want.Equal(Identifier("myapp2"), Name("My-App2").Identifier())
}

func TestNameEnvPrefix(t *testing.T) {
	t.Parallel()
	want := assert.New(t)

	want.Equal(EnvPrefix("BEFORE_CLI"), Name("before.cli").EnvPrefix())
	want.Equal(EnvPrefix("MY_APP2"), Name("my-app2").EnvPrefix())
}
