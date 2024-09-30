package di

import (
	"github.com/MuhammedAshifVnr/Gig_Space/Search_Svc/pkg/es"

	"github.com/MuhammedAshifVnr/Gig_Space/Search_Svc/pkg/internal/repo"
	"github.com/MuhammedAshifVnr/Gig_Space/Search_Svc/pkg/internal/service"
)

func InitElasticSearch() *service.SearchService {
	elastic := es.InitElasticSearch()
	redis := es.Redis()
	repo := repo.NewSearchRepository(elastic, redis)
	service := service.NewSearchService(repo)

	return service
}
