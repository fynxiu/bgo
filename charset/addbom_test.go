package charset

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_addBomToUTF8(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "bgo-test-")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	type args struct {
		filepath string
		content  []byte
	}
	tests := []struct {
		name        string
		args        args
		wantErr     bool
		wantContent []byte
	}{
		{
			name: "file already has BOM",
			args: args{
				filepath: filepath.Join(tmpDir, "test-without-bom.txt"),
				content:  []byte("test without BOM"),
			},
			wantErr:     false,
			wantContent: append(bom, []byte("test without BOM")...),
		},
		{
			name: "file has no BOM",
			args: args{
				filepath: filepath.Join(tmpDir, "test-without-bom.txt"),
				content:  append(bom, []byte("test without BOM")...),
			},
			wantErr:     false,
			wantContent: append(bom, []byte("test without BOM")...),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.NoError(t, os.WriteFile(tt.args.filepath, tt.args.content, 0644))
			if err := addBomToUTF8(tt.args.filepath); (err != nil) != tt.wantErr {
				t.Errorf("addBomToUTF8() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				got, err := os.ReadFile(tt.args.filepath)
				require.NoError(t, err)
				assert.Equal(t, tt.wantContent, got)
			}
		})
	}
}
