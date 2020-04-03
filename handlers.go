package main

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"errors"
	"github.com/NOVAPokemon/utils"
	"github.com/NOVAPokemon/utils/api"
	trainerdb "github.com/NOVAPokemon/utils/database/trainer"
	"github.com/NOVAPokemon/utils/tokens"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"strings"
)

const serviceName = "Trainers"

var decodeError = errors.New("an error occurred decoding the supplied resource")

func AddTrainer(w http.ResponseWriter, r *http.Request) {

	log.Infof("Request to add trainer")
	var trainer = &utils.Trainer{}
	err := json.NewDecoder(r.Body).Decode(trainer)

	if err != nil {
		handleError(decodeError, w)
		return
	}

	log.Infof("Adding trainer: %+v", trainer)
	_, err = trainerdb.AddTrainer(*trainer)

	if err != nil {
		handleError(err, w)
		return
	}

	toSend, err := json.Marshal(*trainer)

	if err != nil {
		handleError(err, w)
		return
	}
	_, err = w.Write(toSend)

	if err != nil {
		panic(err)
	}

}

func HandleUpdateTrainerInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	trainerUsername := vars[api.UsernameVar]
	var trainerStats = &utils.TrainerStats{}

	err := json.NewDecoder(r.Body).Decode(trainerStats)

	if err != nil {
		handleError(err, w)
		return
	}

	trainerStats, err = trainerdb.UpdateTrainerStats(trainerUsername, *trainerStats)

	if err != nil {
		handleError(err, w)
		return
	}

	toSend, err := json.Marshal(trainerStats)

	if err != nil {
		handleError(err, w)
		return
	}
	_, err = w.Write(toSend)

	if err != nil {
		panic(err)
	}
}

func GetAllTrainers(w http.ResponseWriter, _ *http.Request) {
	trainers, err := trainerdb.GetAllTrainers()
	if err != nil {
		http.Error(w, "An error occurred fetching trainers", http.StatusInternalServerError)
	}

	toSend, err := json.Marshal(trainers)

	if err != nil {
		handleError(err, w)
		return
	}

	_, err = w.Write(toSend)

	if err != nil {
		panic(err)
	}

}

func GetTrainerByUsername(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	trainerUsername := vars[api.UsernameVar]

	trainer, err := trainerdb.GetTrainerByUsername(trainerUsername)

	if err != nil {
		handleError(err, w)
		return
	}

	toSend, err := json.Marshal(trainer)
	if err != nil {
		handleError(err, w)
		return
	}
	_, err = w.Write(toSend)

	if err != nil {
		panic(err)
	}
}

func AddPokemonToTrainer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	trainerUsername := vars[api.UsernameVar]

	var pokemon = &utils.Pokemon{}

	err := json.NewDecoder(r.Body).Decode(pokemon)

	if err != nil {
		handleError(err, w)
		return
	}

	pokemon, err = trainerdb.AddPokemonToTrainer(trainerUsername, *pokemon)

	if err != nil {
		handleError(err, w)
		return
	}

	toSend, err := json.Marshal(*pokemon)
	if err != nil {
		handleError(err, w)
		return
	}
	_, err = w.Write(toSend)

	if err != nil {
		panic(err)
	}
}

func RemovePokemonFromTrainer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	trainerUsername := vars[api.UsernameVar]
	pokemonId, err := primitive.ObjectIDFromHex(vars[api.PokemonIdVar])

	if err != nil {
		handleError(err, w)
		return
	}

	err = trainerdb.RemovePokemonFromTrainer(trainerUsername, pokemonId)

	if err != nil {
		handleError(err, w)
		return
	}

}

func AddItemsToTrainer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	trainerUsername := vars[api.UsernameVar]

	var item = &utils.Item{}

	err := json.NewDecoder(r.Body).Decode(item)

	if err != nil {
		handleError(err, w)
		return
	}

	addedItem, err := trainerdb.AddItemToTrainer(trainerUsername, *item)

	if err != nil {
		handleError(err, w)
		return
	}

	toSend, err := json.Marshal(*addedItem)
	if err != nil {
		handleError(err, w)
		return
	}
	_, err = w.Write(toSend)

	if err != nil {
		panic(err)
	}

}

func RemoveItemsFromTrainer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	trainerUsername := vars[api.UsernameVar]
	var itemIds []primitive.ObjectID

	for _, itemIdStr := range strings.Split(vars[api.ItemIdVar], ",") {
		itemId, err := primitive.ObjectIDFromHex(itemIdStr)

		if err != nil {
			handleError(decodeError, w)
			return
		}
		itemIds = append(itemIds, itemId)
	}

	removedItems, err := trainerdb.RemoveItemsFromTrainer(trainerUsername, itemIds)

	if err != nil {
		handleError(err, w)
		return
	}

	toSend, err := json.Marshal(removedItems)
	if err != nil {
		handleError(err, w)
		return
	}
	_, err = w.Write(toSend)

	if err != nil {
		panic(err)
	}
}

// receives a POST request with a hash of the pokemons stats
// returns true or false depending on if they are up to date
func HandleVerifyTrainerPokemons(w http.ResponseWriter, r *http.Request) {

	log.Info("Verify Pokemons request")

	token, err := tokens.ExtractAndVerifyAuthToken(r.Header)

	if err != nil {
		http.Error(w, "no auth token", http.StatusUnauthorized)
		return
	}

	var receivedHashes map[string][]byte
	err = json.NewDecoder(r.Body).Decode(&receivedHashes)

	if err != nil {
		handleError(decodeError, w)
		return
	}

	trainer, err := trainerdb.GetTrainerByUsername(token.Username)

	if err != nil {
		handleError(err, w)
		return
	}

	for pokemonId, currHash := range receivedHashes {
		pokemon, ok := trainer.Pokemons[pokemonId]

		if !ok {
			log.Info("Denied")
			w.WriteHeader(200)
			toSend, _ := json.Marshal(false)
			_, _ = w.Write(toSend)
		}

		pokemonBytes, _ := json.Marshal(pokemon)
		pokemonBytesTemp := md5.Sum(pokemonBytes)
		pokemonHash := pokemonBytesTemp[:]

		if !bytes.Equal(pokemonHash, currHash) {
			log.Info("Denied")
			w.WriteHeader(200)
			toSend, _ := json.Marshal(false)
			_, _ = w.Write(toSend)
			return
		}
	}
	log.Info("Accepted")
	w.WriteHeader(200)
	toSend, _ := json.Marshal(true)
	_, _ = w.Write(toSend)
}

// receives a POST request with a hash of the trainer stats
// returns true or false depending on if they are up to date
func HandleVerifyTrainerStats(w http.ResponseWriter, r *http.Request) {
	log.Info("Verify Trainer Stats requet")
	token, err := tokens.ExtractAndVerifyAuthToken(r.Header)
	if err != nil {
		handleError(err, w)
		return
	}

	var receivedHash []byte
	err = json.NewDecoder(r.Body).Decode(&receivedHash)
	if err != nil {
		log.Error(err)
		return
	}

	trainer, err := trainerdb.GetTrainerByUsername(token.Username)

	if err != nil {
		handleError(err, w)
		return
	}

	statsBytes, _ := json.Marshal(trainer.Stats)
	statsBytesTemp := md5.Sum(statsBytes)
	statsHash := statsBytesTemp[:]

	if bytes.Equal(statsHash, receivedHash) {
		w.WriteHeader(200)
		toSend, _ := json.Marshal(true)
		_, _ = w.Write(toSend)
	} else {
		w.WriteHeader(200)
		toSend, _ := json.Marshal(false)
		_, _ = w.Write(toSend)
	}

}

// receives a POST request with a hash of the trainer items
// returns true or false depending on if they are up to date
func HandleVerifyTrainerItems(w http.ResponseWriter, r *http.Request) {
	log.Info("Verify items requet")
	token, err := tokens.ExtractAndVerifyAuthToken(r.Header)
	if err != nil {
		handleError(err, w)
		return
	}

	var receivedHash []byte
	err = json.NewDecoder(r.Body).Decode(&receivedHash)
	if err != nil {
		log.Error(err)
		return
	}

	trainer, err := trainerdb.GetTrainerByUsername(token.Username)

	if err != nil {
		handleError(err, w)
		return
	}

	itemsBytes, _ := json.Marshal(trainer.Items)
	itemsHashTemp := md5.Sum(itemsBytes)
	itemsHash := itemsHashTemp[:]

	log.Info(itemsHash)
	log.Info(receivedHash)

	if bytes.Equal(itemsHash, receivedHash) {
		w.WriteHeader(200)
		toSend, _ := json.Marshal(true)
		_, _ = w.Write(toSend)
		log.Info("verify items: ", true)
	} else {

		w.WriteHeader(200)
		toSend, _ := json.Marshal(false)
		_, _ = w.Write(toSend)
		log.Info("verify items: ", false)
	}
}

func HandleGenerateAllTokens(w http.ResponseWriter, r *http.Request) {
	log.Info("Generate all tokens requet")
	token, err := tokens.ExtractAndVerifyAuthToken(r.Header)
	if err != nil {
		log.Error(err)
		handleError(err, w)
		return
	}

	trainer, err := trainerdb.GetTrainerByUsername(token.Username)

	if err != nil {
		log.Error(token.Username)
		handleError(err, w)
		return
	}

	tokens.AddItemsToken(trainer.Items, w.Header())
	tokens.AddPokemonsTokens(trainer.Pokemons, w.Header())
	tokens.AddTrainerStatsToken(trainer.Stats, w.Header())
}

func HandleGenerateTrainerStatsToken(w http.ResponseWriter, r *http.Request) {
	token, err := tokens.ExtractAndVerifyAuthToken(r.Header)
	if err != nil {
		handleError(err, w)
		return
	}

	trainer, err := trainerdb.GetTrainerByUsername(token.Username)
	if err != nil {
		handleError(err, w)
		return
	}
	tokens.AddTrainerStatsToken(trainer.Stats, w.Header())
}

func HandleGeneratePokemonsToken(w http.ResponseWriter, r *http.Request) {
	token, err := tokens.ExtractAndVerifyAuthToken(r.Header)
	if err != nil {
		handleError(err, w)
		return
	}

	trainer, err := trainerdb.GetTrainerByUsername(token.Username)
	if err != nil {
		handleError(err, w)
		return
	}
	tokens.AddPokemonsTokens(trainer.Pokemons, w.Header())
}

func HandleGenerateItemsToken(w http.ResponseWriter, r *http.Request) {
	token, err := tokens.ExtractAndVerifyAuthToken(r.Header)
	if err != nil {
		return
	}

	trainer, err := trainerdb.GetTrainerByUsername(token.Username)
	if err != nil {
		handleError(err, w)
		return
	}
	tokens.AddItemsToken(trainer.Items, w.Header())
}

func handleError(err error, w http.ResponseWriter) {
	log.Error(err)

	switch err {
	case trainerdb.ErrTrainerNotFound:
		http.Error(w, err.Error(), http.StatusNotFound)

	case trainerdb.ErrItemNotFound:
		http.Error(w, err.Error(), http.StatusNotFound)

	case trainerdb.ErrPokemonNotFound:
		http.Error(w, err.Error(), http.StatusNotFound)

	case trainerdb.ErrInvalidCoins:
		http.Error(w, err.Error(), http.StatusBadRequest)

	case trainerdb.ErrInvalidLevel:
		http.Error(w, err.Error(), http.StatusBadRequest)

	case decodeError:
		http.Error(w, err.Error(), http.StatusInternalServerError)

	default:
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
