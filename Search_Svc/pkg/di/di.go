package di

import (
	"github.com/MuhammedAshifVnr/Gig_Space/Search_Svc/pkg/es"
	"github.com/MuhammedAshifVnr/Gig_Space/Search_Svc/pkg/internal/repo"
	"github.com/MuhammedAshifVnr/Gig_Space/Search_Svc/pkg/internal/service"
)

func InitElasticSearch() *service.SearchService{
	es := es.InitElasticSearch()
	repo := repo.NewSearchRepository(es)
	service:=service.NewSearchService(repo)
	
	return service
}
