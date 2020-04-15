package main

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"errors"
	"github.com/NOVAPokemon/utils"
	"github.com/NOVAPokemon/utils/api"
	trainerdb "github.com/NOVAPokemon/utils/database/trainer"
	"github.com/NOVAPokemon/utils/experience"
	"github.com/NOVAPokemon/utils/items"
	"github.com/NOVAPokemon/utils/pokemons"
	"github.com/NOVAPokemon/utils/tokens"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"strings"
)

var decodeError = errors.New("an error occurred decoding the supplied resource")

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

func AddTrainer(w http.ResponseWriter, r *http.Request) {
	log.Infof("Request to add trainer")
	var trainer = utils.Trainer{}
	err := json.NewDecoder(r.Body).Decode(&trainer)
	if err != nil {
		handleError(decodeError, w)
		return
	}

	trainer.Items = generateStarterItems()

	log.Infof("Adding trainer: %+v", trainer)
	_, err = trainerdb.AddTrainer(trainer)
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
	trainerUsername := vars[api.UsernameVar]

	_, err := tokens.ExtractAndVerifyAuthToken(r.Header)
	if err != nil {
		handleError(err, w)
		return
	}

	var trainerStats = utils.TrainerStats{}
	err = json.NewDecoder(r.Body).Decode(&trainerStats)
	if err != nil {
		handleError(err, w)
		return
	}

	trainerStats.Level = experience.CalculateLevel(trainerStats.XP)
	updatedTrainerStats, err := trainerdb.UpdateTrainerStats(trainerUsername, trainerStats)
	if err != nil {
		handleError(err, w)
		return
	}

	toSend, err := json.Marshal(trainerStats)
	if err != nil {
		handleError(err, w)
		return
	}

	tokens.AddTrainerStatsToken(*updatedTrainerStats, w.Header())
	_, err = w.Write(toSend)

	if err != nil {
		panic(err)
	}
}

func AddPokemonToTrainer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	trainerUsername := vars[api.UsernameVar]

	var pokemon = pokemons.Pokemon{}

	err := json.NewDecoder(r.Body).Decode(&pokemon)

	pokemonId := primitive.NewObjectID()
	pokemon.Id = pokemonId

	if err != nil {
		handleError(err, w)
		return
	}

	updatedPokemons, err := trainerdb.AddPokemonToTrainer(trainerUsername, pokemon)
	if err != nil {
		handleError(err, w)
		return
	}
	tokens.AddPokemonsTokens(updatedPokemons, w.Header())

	toSend, err := json.Marshal(updatedPokemons[pokemonId.Hex()])
	if err != nil {
		handleError(err, w)
		return
	}

	_, err = w.Write(toSend)
	if err != nil {
		panic(err)
	}
}

func HandleUpdatePokemon(w http.ResponseWriter, r *http.Request) {
	log.Info("Update pokemon request")
	vars := mux.Vars(r)
	trainerUsername := vars[api.UsernameVar]
	pokemonId, err := primitive.ObjectIDFromHex(vars[api.PokemonIdVar])

	if err != nil {
		handleError(err, w)
		return
	}

	var pokemon = pokemons.Pokemon{}
	err = json.NewDecoder(r.Body).Decode(&pokemon)
	if err != nil {
		handleError(err, w)
		return
	}

	pokemon.Level = experience.CalculateLevel(pokemon.XP)
	updatedPokemons, err := trainerdb.UpdateTrainerPokemon(trainerUsername, pokemonId, pokemon)

	if err != nil {
		handleError(err, w)
		return
	}

	marshaledPokemons, err := json.Marshal(updatedPokemons[pokemonId.Hex()])

	if err != nil {
		handleError(err, w)
		return
	}

	tokens.AddPokemonsTokens(updatedPokemons, w.Header())
	_, err = w.Write(marshaledPokemons)

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

	oldTrainerPokemons, err := trainerdb.RemovePokemonFromTrainer(trainerUsername, pokemonId)

	if err != nil {
		handleError(err, w)
		return
	}

	removedPokemon, ok := oldTrainerPokemons[pokemonId.Hex()]

	if !ok {
		http.NotFound(w, r)
		return
	}

	delete(oldTrainerPokemons, pokemonId.Hex())
	tokens.AddPokemonsTokens(oldTrainerPokemons, w.Header())
	toSend, err := json.Marshal(removedPokemon)
	_, err = w.Write(toSend)

	if err != nil {
		panic(err)
	}

}

func AddItemsToTrainer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	_ = vars[api.UsernameVar]

	token, err := tokens.ExtractAndVerifyAuthToken(r.Header)
	if err != nil {
		log.Error(err)
		return
	}

	var itemsToAdd []items.Item
	err = json.NewDecoder(r.Body).Decode(&itemsToAdd)
	if err != nil {
		handleError(err, w)
		return
	}

	for _, item := range itemsToAdd {
		itemId := primitive.NewObjectID()
		item.Id = itemId
	}

	updatedItems, err := trainerdb.AddItemsToTrainer(token.Username, itemsToAdd)

	if err != nil {
		handleError(err, w)
		return
	}

	tokens.AddItemsToken(updatedItems, w.Header())
	toSend, err := json.Marshal(itemsToAdd)
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
	_ = vars[api.UsernameVar]

	token, err := tokens.ExtractAndVerifyAuthToken(r.Header)
	if err != nil {
		log.Error(err)
		return
	}

	var itemIds []primitive.ObjectID
	for _, itemIdStr := range strings.Split(vars[api.ItemIdVar], ",") {
		itemId, err := primitive.ObjectIDFromHex(itemIdStr)
		if err != nil {
			handleError(decodeError, w)
			return
		}
		itemIds = append(itemIds, itemId)
	}

	oldTrainerItems, err := trainerdb.RemoveItemsFromTrainer(token.Username, itemIds)
	if err != nil {
		handleError(err, w)
		return
	}

	removedItems := make(map[string]items.Item, len(itemIds))
	for i := 0; i < len(itemIds); i++ {
		item, ok := oldTrainerItems[itemIds[i].Hex()]
		if ok {
			removedItems[item.Id.Hex()] = item
			delete(oldTrainerItems, itemIds[i].Hex())
		}
	}

	toSend, err := json.Marshal(removedItems)
	if err != nil {
		handleError(err, w)
		return
	}
	tokens.AddItemsToken(oldTrainerItems, w.Header())
	_, err = w.Write(toSend)

	if err != nil {
		panic(err)
	}
}

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

func HandleVerifyTrainerStats(w http.ResponseWriter, r *http.Request) {
	log.Info("Verify Trainer Stats request")
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

func HandleVerifyTrainerItems(w http.ResponseWriter, r *http.Request) {
	log.Info("Verify items request")
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

	log.Info(token.Username)
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
	log.Info("Generate all tokens request")
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

func generateStarterItems() map[string]items.Item {
	pokeBallsAmount := 10
	masterBallsAmount := 1
	healAmount := 10
	reviveAmount := 1
	totalAmount := pokeBallsAmount + masterBallsAmount + healAmount + reviveAmount

	starterItems := make(map[string]items.Item, totalAmount)

	for i := 0; i < pokeBallsAmount; i++ {
		toAdd := items.PokeBallItem
		toAdd.Id = primitive.NewObjectID()
		starterItems[toAdd.Id.Hex()] = toAdd
	}

	for i := 0; i < masterBallsAmount; i++ {
		toAdd := items.MasterBallItem
		toAdd.Id = primitive.NewObjectID()
		starterItems[toAdd.Id.Hex()] = toAdd
	}

	for i := 0; i < healAmount; i++ {
		toAdd := items.HealItem
		toAdd.Id = primitive.NewObjectID()
		starterItems[toAdd.Id.Hex()] = toAdd
	}

	for i := 0; i < reviveAmount; i++ {
		toAdd := items.ReviveItem
		toAdd.Id = primitive.NewObjectID()
		starterItems[toAdd.Id.Hex()] = toAdd
	}

	return starterItems
}
