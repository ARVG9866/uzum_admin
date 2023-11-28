package docs

import "github.com/mvrilo/go-redoc"

func Initialize() redoc.Redoc {
	return redoc.Redoc{
		Title:       "Documentation of adminSystem",
		Description: "Documentation describes working procedures of adminSystem like structs, handlers, etc.",
		SpecFile:    "./docs/api_admin_v1.swagger.json",
		SpecPath:    "/docs/api_admin_v1.swagger.json",
		DocsPath:    "/docs",
	}

}
