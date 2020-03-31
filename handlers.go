package main

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"github.com/NOVAPokemon/utils"
	"github.com/NOVAPokemon/utils/cookies"
	trainerdb "github.com/NOVAPokemon/utils/database/trainer"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"reflect"
	"strings"
)

var key = []byte("my_secret_key")

const serviceName = "Trainers"

var decodeError = errors.New("an error occurred decoding the supplied resource")

func AddTrainer(w http.ResponseWriter, r *http.Request) {

	var trainer = &utils.Trainer{}
	err := json.NewDecoder(r.Body).Decode(trainer)

	if err != nil {
		handleError(decodeError, w)
		return
	}

	_, err = trainerdb.AddTrainer(*trainer)

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

func HandleUpdateTrainerInfo(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	trainerUsername := vars["username"]
	var trainer = &utils.Trainer{}

	err := json.NewDecoder(r.Body).Decode(trainer)

	if err != nil {
		handleError(err, w)
		return
	}

	trainer, err = trainerdb.UpdateTrainerStats(trainerUsername, *trainer)

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
	trainerUsername := vars["username"]

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
	trainerUsername := vars["username"]

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
	trainerUsername := vars["username"]
	pokemonId, err := primitive.ObjectIDFromHex(vars["pokemonId"])

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
	trainerUsername := vars["username"]

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
	trainerUsername := vars["username"]
	var itemIds []primitive.ObjectID

	for _, itemIdStr := range strings.Split(vars["itemIds"], ",") {

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
	token, err := cookies.ExtractAndVerifyAuthToken(&w, r, serviceName)
	if err != nil {
		return
	}

	var pokemonHashes map[string][]byte
	err = json.NewDecoder(r.Body).Decode(pokemonHashes)

	if err != nil {
		handleError(decodeError, w)
		return
	}

	trainer, err := trainerdb.GetTrainerByUsername(token.Username)

	if err != nil {
		handleError(err, w)
		return
	}

	for pokemonId, pokemon := range trainer.Pokemons {
		marshaled, _ := json.Marshal(pokemon)
		currHash := md5.Sum(marshaled)
		if reflect.DeepEqual(currHash, pokemonHashes[pokemonId]) {
			w.WriteHeader(200)
			toSend, _ := json.Marshal(true)
			_, _ = w.Write(toSend)
		} else {
			w.WriteHeader(200)
			toSend, _ := json.Marshal(false)
			_, _ = w.Write(toSend)
		}
	}
}

// receives a POST request with a hash of the trainer stats
// returns true or false depending on if they are up to date
func HandleVerifyTrainerStats(w http.ResponseWriter, r *http.Request) {
	token, err := cookies.ExtractAndVerifyAuthToken(&w, r, serviceName)
	if err != nil {
		return
	}

	hash := r.Body
	trainer, err := trainerdb.GetTrainerByUsername(token.Username)

	if err != nil {
		handleError(err, w)
		return
	}

	marshaled, _ := json.Marshal(trainer.Stats)
	currHash := md5.Sum(marshaled)
	if reflect.DeepEqual(currHash, hash) {
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
	token, err := cookies.ExtractAndVerifyAuthToken(&w, r, serviceName)
	if err != nil {
		return
	}

	receivedHash := r.Body
	trainer, err := trainerdb.GetTrainerByUsername(token.Username)

	if err != nil {
		handleError(err, w)
		return
	}

	itemsBytes, _ := json.Marshal(trainer.Items)
	itemsHash := md5.Sum(itemsBytes)
	if reflect.DeepEqual(itemsHash, receivedHash) {
		w.WriteHeader(200)
		toSend, _ := json.Marshal(true)
		_, _ = w.Write(toSend)
	} else {
		w.WriteHeader(200)
		toSend, _ := json.Marshal(false)
		_, _ = w.Write(toSend)
	}
}

func HandleGenerateAllTokens(w http.ResponseWriter, r *http.Request) {
	token, err := cookies.ExtractAndVerifyAuthToken(&w, r, serviceName)
	if err != nil {
		return
	}
	
	trainer, err := trainerdb.GetTrainerByUsername(token.Username)

	if err != nil {
		log.Error(token.Username)
		handleError(err, w)
		return
	}

	cookies.SetItemsCookie(trainer.Items, w, key)
	cookies.SetPokemonsCookie(trainer.Pokemons, w, key)
	cookies.SetTrainerStatsCookie(trainer.Stats, w, key)
}

func HandleGenerateTrainerStatsToken(w http.ResponseWriter, r *http.Request) {
	token, err := cookies.ExtractAndVerifyAuthToken(&w, r, serviceName)
	if err != nil {
		return
	}

	trainer, err := trainerdb.GetTrainerByUsername(token.Username)
	if err != nil {
		handleError(err, w)
		return
	}
	cookies.SetTrainerStatsCookie(trainer.Stats, w, key)
}

func HandleGeneratePokemonsToken(w http.ResponseWriter, r *http.Request) {
	token, err := cookies.ExtractAndVerifyAuthToken(&w, r, serviceName)
	if err != nil {
		return
	}

	trainer, err := trainerdb.GetTrainerByUsername(token.Username)
	if err != nil {
		handleError(err, w)
		return
	}
	cookies.SetPokemonsCookie(trainer.Pokemons, w, key)
}

func HandleGenerateItemsToken(w http.ResponseWriter, r *http.Request) {
	token, err := cookies.ExtractAndVerifyAuthToken(&w, r, serviceName)
	if err != nil {
		return
	}

	trainer, err := trainerdb.GetTrainerByUsername(token.Username)
	if err != nil {
		handleError(err, w)
		return
	}
	cookies.SetItemsCookie(trainer.Items, w, key)
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
		http.Error(w, "An error has occurred", http.StatusInternalServerError)
	}
}
