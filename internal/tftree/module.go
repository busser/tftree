package tftree

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	tfjson "github.com/hashicorp/terraform-json"
)

type Module struct {
	Name     string   `json:"name"`
	Source   string   `json:"source"`
	Children []Module `json:"children"`
}

type Alphabetical []Module

func (mods Alphabetical) Len() int {
	return len(mods)
}

func (mods Alphabetical) Swap(i, j int) {
	mods[i], mods[j] = mods[j], mods[i]
}

func (mods Alphabetical) Less(i, j int) bool {
	if mods[i].Name < mods[j].Name {
		return true
	}
	if mods[i].Name > mods[j].Name {
		return false
	}
	return mods[i].Source < mods[j].Source
}

func New(plan *tfjson.Plan, rootDir string) (Module, error) {
	c, err := childrenOf(plan.Config.RootModule)
	if err != nil {
		return Module{}, fmt.Errorf("could not compute module tree: %w", err)
	}

	root := Module{
		Name:     "root",
		Source:   rootDir,
		Children: c,
	}

	setFullPaths(&root)

	return root, nil
}

func childrenOf(mod *tfjson.ConfigModule) ([]Module, error) {
	var children []Module

	for name, call := range mod.ModuleCalls {
		c, err := childrenOf(call.Module)
		if err != nil {
			return nil, fmt.Errorf("module %q: %w", name, err)
		}

		m := Module{
			Name:     name,
			Source:   call.Source,
			Children: c,
		}

		children = append(children, m)
	}

	sort.Sort(Alphabetical(children))

	return children, nil
}

var moduleSourceLocalPrefixes = []string{
	"./",
	"../",
	".\\",
	"..\\",
}

func setFullPaths(root *Module) {
	var helper func(*Module, string)
	helper = func(mod *Module, parentSource string) {
		if mod.sourceIsLocal() {
			fullSource := "." + string(filepath.Separator) + filepath.Join(parentSource, mod.Source)
			mod.Source = fullSource
		}

		for i := range mod.Children {
			helper(&mod.Children[i], mod.Source)
		}
	}

	for i := range root.Children {
		helper(&root.Children[i], root.Source)
	}
}

func (mod Module) sourceIsLocal() bool {
	for _, prefix := range moduleSourceLocalPrefixes {
		if strings.HasPrefix(mod.Source, prefix) {
			return true
		}
	}
	return false
}
