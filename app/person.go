package app

import (
	"fmt"
	"gomdb/internal/pkg/domain"
	"gomdb/internal/pkg/logging"
)

type personSvc struct {
	DB  domain.PersonDB
	API domain.PersonAPI
}

func NewPersonSvc(db domain.PersonDB, api domain.PersonAPI) domain.PersonSvc {
	return personSvc{
		DB:  db,
		API: api,
	}
}

func (ps personSvc) Get(id int) (*domain.Person, error) {
	person, err := ps.DB.Get(id, "people")
	return (*person).(*domain.Person), err
}

func (ps personSvc) List(query string) ([]*domain.Person, error) {
	person, err := ps.DB.List(query, "people")
	personSlice := make([]*domain.Person, len(person))
	for i, v := range person {
		person, ok := (*v).(*domain.Person)
		if !ok {
			logging.Error(fmt.Sprintf("Failed to convert element at index %d to *domain.Person\n", i))
			continue
		}
		personSlice[i] = person
	}
	return personSlice, err
}

func (ps personSvc) Upsert(person *domain.Person) error {
	var personDB interface{} = person
	return ps.DB.Upsert(&personDB, "people")
}

func (ps personSvc) GetFromAPI(person *domain.Person) error {
	var personAPI interface{} = person
	return ps.API.GetFromAPI(&personAPI, "people")
}
