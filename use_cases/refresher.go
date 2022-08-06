package use_cases

import (
	"context"
	"log"
	"strings"
	"sync"

	"GoConcurrency-Bootcamp-2022/models"
)

const MAX_WORKERS = 128

type reader interface {
	Read() ([]models.Pokemon, error)
}

type saver interface {
	Save(context.Context, []models.Pokemon) error
}

type fetcher interface {
	FetchAbility(string) (models.Ability, error)
}

type Refresher struct {
	reader
	saver
	fetcher
}

func NewRefresher(reader reader, saver saver, fetcher fetcher) Refresher {
	return Refresher{reader, saver, fetcher}
}

func (r Refresher) Refresh(ctx context.Context) error {
	var pokemons []models.Pokemon

	dataChan := readCSVFilePokemons(r)

	// Workers
	var workers = make(map[int]<-chan models.Pokemon)
	for i := 0; i < MAX_WORKERS; i++ {
		workers[i] = fetchAbilities(dataChan, r)
	}

	for pokemon := range fanInPokemons(workers) {
		pokemons = append(pokemons, pokemon)
	}

	return r.Save(ctx, pokemons)
}

func readCSVFilePokemons(r Refresher) chan models.Pokemon {
	outputPokemons := make(chan models.Pokemon)
	go func() {
		defer close(outputPokemons)

		pokemons, err := r.Read()
		if err != nil {
			log.Printf("ERROR: csv file could not read! %s", err.Error())
		}
		for _, pokemon := range pokemons {
			outputPokemons <- pokemon
		}
	}()

	return outputPokemons
}

func fetchAbilities(pokemons <-chan models.Pokemon, r Refresher) <-chan models.Pokemon {
	outputAbilities := make(chan models.Pokemon)

	go func() {
		defer close(outputAbilities)
		for pokemon := range pokemons {
			urls := strings.Split(pokemon.FlatAbilityURLs, "|")
			var abilities []string
			for _, url := range urls {
				ability, err := r.FetchAbility(url)
				if err != nil {
					log.Printf("ERROR: something went wrong fetching the abilities")
				}

				for _, ee := range ability.EffectEntries {
					abilities = append(abilities, ee.Effect)
				}

				pokemon.EffectEntries = abilities
			}
			outputAbilities <- pokemon
		}
	}()

	return outputAbilities
}

func fanInPokemons(chans map[int]<-chan models.Pokemon) <-chan models.Pokemon {
	outputPokemons := make(chan models.Pokemon)
	var wg sync.WaitGroup

	wg.Add(len(chans))
	for _, ch := range chans {
		go func(ch <-chan models.Pokemon) {
			defer wg.Done()
			for pokemon := range ch {
				outputPokemons <- pokemon
			}
		}(ch)
	}

	go func() {
		wg.Wait()
		close(outputPokemons)
	}()

	return outputPokemons
}
