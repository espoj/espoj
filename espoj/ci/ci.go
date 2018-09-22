package ci

import "github.com/samitc/espoj/espoj/data"

type ICIImporter interface {
	GetCI() data.CI
}

func CreateCI(ciImport ICIImporter) data.CI {
	return ciImport.GetCI()
}