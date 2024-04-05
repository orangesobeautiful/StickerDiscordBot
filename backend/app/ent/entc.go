//go:build ignore

package main

import (
	"log"
	"strings"
	"text/template"

	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"golang.org/x/xerrors"
)

func main() {
	var err error
	defer func() {
		if err != nil {
			log.Fatalf("error occure:\n%+v", err)
		}
	}()

	err = entc.Generate("./schema",
		&gen.Config{
			Features: []gen.Feature{
				gen.FeatureUpsert,
			},
			Templates: []*gen.Template{
				gen.MustParse(gen.NewTemplate("static").
					Funcs(template.FuncMap{"title": strings.ToTitle}).
					ParseDir("./extemplates"),
				),
			},
		},
	)
	if err != nil {
		err = xerrors.Errorf("running ent codegen: %w", err)
		return
	}
}
