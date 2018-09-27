package data

type Data struct {
	repo map[string]*Repository
}

func NewData() *Data {
	return &Data{repo: make(map[string]*Repository)}
}
func (dat *Data) AddRepository(name string, repository *Repository) {
	if repository != nil {
		dat.repo[repository.Name] = repository
		dat.repo[name] = repository
	}
}
func (dat *Data) GetRepository(name string) *Repository {
	return dat.repo[name]
}
