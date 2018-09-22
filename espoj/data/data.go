package data

type Data struct {
	repo map[string]*Repository
}

func NewData() *Data {
	return &Data{repo: make(map[string]*Repository)}
}
func (dat *Data) AddRepository(repository *Repository) {
	if repository != nil {
		dat.repo[repository.Name] = repository
	}
}
func (dat *Data) GetRepository(name string) *Repository {
	return dat.repo[name]
}
