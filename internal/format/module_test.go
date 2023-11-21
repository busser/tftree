package format

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/busser/tftree/internal/tftree"
)

func TestModule(t *testing.T) {
	tt := []struct {
		name string

		mod     tftree.Module
		noColor bool

		want string
	}{
		{
			name: "empty",
			mod: tftree.Module{
				Name:     "root",
				Source:   "./workspace",
				Children: nil,
			},
			noColor: false,
			want:    filepath.Join("testdata", "module", "empty.txt"),
		},
		{
			name: "empty no color",
			mod: tftree.Module{
				Name:     "root",
				Source:   "./workspace",
				Children: nil,
			},
			noColor: true,
			want:    filepath.Join("testdata", "module", "empty-no-color.txt"),
		},
		{
			name: "complex",
			mod: tftree.Module{
				Name:   "root",
				Source: "./workspace",
				Children: []tftree.Module{
					{
						Name:   "child_1",
						Source: "./modules/module-1",
						Children: []tftree.Module{
							{
								Name:     "grandchild_11",
								Source:   "./modules/module-2",
								Children: []tftree.Module{},
							},
							{
								Name:     "grandchild_12",
								Source:   "./modules/module-3",
								Children: []tftree.Module{},
							},
						},
					},
					{
						Name:     "child_2",
						Source:   "./modules/module-2",
						Children: nil,
					},
					{
						Name:   "child_3",
						Source: "./modules/module-4",
						Children: []tftree.Module{
							{
								Name:     "grandchild_31",
								Source:   "./modules/module-3",
								Children: []tftree.Module{},
							},
						},
					},
				},
			},
			noColor: false,
			want:    filepath.Join("testdata", "module", "complex.txt"),
		},
		{
			name: "complex no color",
			mod: tftree.Module{
				Name:   "root",
				Source: "./workspace",
				Children: []tftree.Module{
					{
						Name:   "child_1",
						Source: "./modules/module-1",
						Children: []tftree.Module{
							{
								Name:     "grandchild_11",
								Source:   "./modules/module-2",
								Children: []tftree.Module{},
							},
							{
								Name:     "grandchild_12",
								Source:   "./modules/module-3",
								Children: []tftree.Module{},
							},
						},
					},
					{
						Name:     "child_2",
						Source:   "./modules/module-2",
						Children: nil,
					},
					{
						Name:   "child_3",
						Source: "./modules/module-4",
						Children: []tftree.Module{
							{
								Name:     "grandchild_31",
								Source:   "./modules/module-3",
								Children: []tftree.Module{},
							},
						},
					},
				},
			},
			noColor: true,
			want:    filepath.Join("testdata", "module", "complex-no-color.txt"),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {

			// Set NoColor for the duration of the test.
			originalNoColor := NoColor
			NoColor = tc.noColor
			defer func() {
				NoColor = originalNoColor
			}()

			actual := Module(tc.mod)

			if *update {
				stringToFile(t, tc.want, actual)
			}

			want := stringFromFile(t, tc.want)

			const escapeSequence = "\x1b"
			if tc.noColor && strings.Contains(want, escapeSequence) {
				t.Errorf("Module() output contains escape sequence %q even though color is disabled:\n%q", escapeSequence, want)
			}

			if want != actual {
				t.Errorf("Module() mismatch\nWant:\n%s\nGot:\n%s", want, actual)
			}
		})
	}
}
