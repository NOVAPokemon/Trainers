package main

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/NOVAPokemon/utils"
	trainerdb "github.com/NOVAPokemon/utils/database/trainer"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"reflect"
	"strings"
)

var TokenCookieName = "trainer_token"

type TrainerToken struct {
	TrainerUserName string

	TrainerHash   []byte
	PokemonHashes map[string][]byte
	ItemsHash     []byte

	jwt.StandardClaims
}

var decodeError = errors.New("error occurred decoding the supplied resource")

func AddTrainer(w http.ResponseWriter, r *http.Request) {

	var trainer = &utils.Trainer{}
	err := json.NewDecoder(r.Body).Decode(trainer)

	if err != nil {
		handleError(decodeError, w, r)
		return
	}

	_, err = trainerdb.AddTrainer(*trainer)

	if err != nil {
		handleError(err, w, r)
		return
	}

	toSend, err := json.Marshal(trainer)

	if err != nil {
		handleError(err, w, r)
		return
	}
	_, err = w.Write(toSend)

	if err != nil {
		panic(err)
	}

}

func UpdateTrainerInfo(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	trainerUsername := vars["username"]
	var trainer = &utils.Trainer{}

	err := json.NewDecoder(r.Body).Decode(trainer)

	if err != nil {
		handleError(err, w, r)
		return
	}

	trainer, err = trainerdb.UpdateTrainerStats(trainerUsername, *trainer)

	if err != nil {
		handleError(err, w, r)
		return
	}

	toSend, err := json.Marshal(trainer)

	if err != nil {
		handleError(err, w, r)
		return
	}
	_, err = w.Write(toSend)

	if err != nil {
		panic(err)
	}
}

func GetAllTrainers(w http.ResponseWriter, r *http.Request) {

	err, trainers := trainerdb.GetAllTrainers()

	if err != nil {
		http.Error(w, "An error occurred fetching trainers", http.StatusInternalServerError)
	}

	toSend, err := json.Marshal(trainers)

	if err != nil {
		handleError(err, w, r)
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
		handleError(err, w, r)
		return
	}

	toSend, err := json.Marshal(trainer)
	if err != nil {
		handleError(err, w, r)
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
		handleError(err, w, r)
		return
	}

	pokemon, err = trainerdb.AddPokemonToTrainer(trainerUsername, *pokemon)

	if err != nil {
		handleError(err, w, r)
		return
	}

	toSend, err := json.Marshal(*pokemon)
	if err != nil {
		handleError(err, w, r)
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
		handleError(err, w, r)
		return
	}

	err = trainerdb.RemovePokemonFromTrainer(trainerUsername, pokemonId)

	if err != nil {
		handleError(err, w, r)
		return
	}

}

func AddItemsToTrainer(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	trainerUsername := vars["username"]

	var item = &utils.Item{}

	err := json.NewDecoder(r.Body).Decode(item)

	if err != nil {
		handleError(err, w, r)
		return
	}

	addedItem, err := trainerdb.AddItemToTrainer(trainerUsername, *item)

	if err != nil {
		handleError(err, w, r)
		return
	}

	toSend, err := json.Marshal(*addedItem)
	if err != nil {
		handleError(err, w, r)
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
			handleError(decodeError, w, r)
			return
		}
		itemIds = append(itemIds, itemId)
	}

	removedItems, err := trainerdb.RemoveItemsFromTrainer(trainerUsername, itemIds)

	if err != nil {
		handleError(err, w, r)
		return
	}

	toSend, err := json.Marshal(removedItems)
	if err != nil {
		handleError(err, w, r)
		return
	}
	_, err = w.Write(toSend)

	if err != nil {
		panic(err)
	}
}

// receives a GET request with a cookie containing a trainer token
// returns a Json object with the outdated fields
func HandleVerifyTrainerToken(r *http.Request, w http.ResponseWriter) {

	c, err := r.Cookie(TokenCookieName)
	missingFields := utils.Trainer{}

	if err != nil {
		log.Error(err)
		http.Error(w, "Invalid cookie", http.StatusBadRequest)
		return
	}

	tknStr := c.Value
	token := &TrainerToken{}

	err = json.Unmarshal([]byte(tknStr), token)

	if err != nil {
		//TODO
		panic(err)
		return
	}

	trainer, err := trainerdb.GetTrainerByUsername(token.TrainerUserName)

	if err != nil {
		handleError(err, w, r)
		return
	}

	// verify trainer stats hash

	hash := md5.Sum([]byte(fmt.Sprintf("%s-%d-%d", trainer.Username, trainer.Level, trainer.Coins)))
	if !reflect.DeepEqual(hash, token.TrainerHash) {
		missingFields.Level = trainer.Level
		missingFields.Username = trainer.Username
		missingFields.Coins = trainer.Coins
		token.TrainerHash = hash[:]
	}

	// verify trainer items hash
	marshaled, _ := json.Marshal(trainer.Items)
	hash = md5.Sum(marshaled)
	if !reflect.DeepEqual(hash, token.ItemsHash) {
		missingFields.Items = trainer.Items
		token.ItemsHash = hash[:]
	}

	// verify pokemons hash

	for k, v := range trainer.Pokemons {
		marshaled, _ := json.Marshal(v)
		hash = md5.Sum(marshaled)
		if !reflect.DeepEqual(hash, token.PokemonHashes[k]) {
			missingFields.Pokemons[k] = v
			token.ItemsHash = hash[:]
		}
	}

	toSend, err := json.Marshal(missingFields)
	if err != nil {
		handleError(err, w, r)
		return
	}
	_, err = w.Write(toSend)

	if err != nil {
		panic(err)
	}

}

func handleError(err error, w http.ResponseWriter, r *http.Request) {

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
