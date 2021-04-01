package main

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/NOVAPokemon/utils"
	"github.com/NOVAPokemon/utils/api"
	trainerdb "github.com/NOVAPokemon/utils/database/trainer"
	"github.com/NOVAPokemon/utils/experience"
	"github.com/NOVAPokemon/utils/items"
	"github.com/NOVAPokemon/utils/pokemons"
	"github.com/NOVAPokemon/utils/tokens"
	"github.com/NOVAPokemon/utils/websockets"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	serverName   string
	commsManager websockets.CommunicationManager
)

func init() {
	if aux, exists := os.LookupEnv(utils.HostnameEnvVar); exists {
		serverName = aux
	} else {
		log.Fatal("Could not load server name")
	}
}

func getAllTrainers(w http.ResponseWriter, _ *http.Request) {
	trainers, err := trainerdb.GetAllTrainers()
	if err != nil {
		utils.LogAndSendHTTPError(&w, wrapGetAllTrainersError(err), http.StatusInternalServerError)
		return
	}

	toSend, err := json.Marshal(trainers)
	if err != nil {
		utils.LogAndSendHTTPError(&w, wrapGetAllTrainersError(err), http.StatusInternalServerError)
		return
	}

	log.Info("getting all trainers")

	_, err = w.Write(toSend)
	if err != nil {
		utils.LogAndSendHTTPError(&w, wrapGetAllTrainersError(err), http.StatusInternalServerError)
	}
}

func getTrainerByUsername(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	trainerUsername := vars[api.UsernameVar]

	trainer, err := trainerdb.GetTrainerByUsername(trainerUsername)
	if err != nil {
		utils.LogAndSendHTTPError(&w, wrapGetTrainerError(err), http.StatusInternalServerError)
		return
	}

	toSend, err := json.Marshal(trainer)
	if err != nil {
		utils.LogAndSendHTTPError(&w, wrapGetTrainerError(err), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(toSend)
	if err != nil {
		utils.LogAndSendHTTPError(&w, wrapGetTrainerError(err), http.StatusInternalServerError)
	}
}

func addTrainer(w http.ResponseWriter, r *http.Request) {
	log.Infof("Request to add trainer")

	trainer := utils.Trainer{}
	err := json.NewDecoder(r.Body).Decode(&trainer)
	if err != nil {
		utils.LogAndSendHTTPError(&w, wrapAddTrainerError(err), http.StatusInternalServerError)
		return
	}

	trainer.Items = generateStarterItems()

	log.Infof("Adding trainer: %s", trainer.Username)
	_, err = trainerdb.AddTrainer(trainer)
	if err != nil {
		utils.LogAndSendHTTPError(&w, wrapAddTrainerError(err), http.StatusInternalServerError)
		return
	}

	toSend, err := json.Marshal(trainer)
	if err != nil {
		utils.LogAndSendHTTPError(&w, wrapAddTrainerError(err), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(toSend)
	if err != nil {
		utils.LogAndSendHTTPError(&w, wrapAddTrainerError(err), http.StatusInternalServerError)
	}
}

func handleUpdateTrainerInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	trainerUsername := vars[api.UsernameVar]

	_, err := tokens.ExtractAndVerifyAuthToken(r.Header)
	if err != nil {
		utils.LogAndSendHTTPError(&w, wrapUpdateTrainerError(err), http.StatusUnauthorized)
		return
	}

	trainerStats := utils.TrainerStats{}
	err = json.NewDecoder(r.Body).Decode(&trainerStats)
	if err != nil {
		utils.LogAndSendHTTPError(&w, wrapUpdateTrainerError(err), http.StatusInternalServerError)
		return
	}

	trainerStats.Level = experience.CalculateLevel(trainerStats.XP)
	updatedTrainerStats, err := trainerdb.UpdateTrainerStats(trainerUsername, trainerStats)
	if err != nil {
		utils.LogAndSendHTTPError(&w, wrapUpdateTrainerError(err), http.StatusInternalServerError)
		return
	}

	toSend, err := json.Marshal(trainerStats)
	if err != nil {
		utils.LogAndSendHTTPError(&w, wrapUpdateTrainerError(err), http.StatusInternalServerError)
		return
	}

	tokens.AddTrainerStatsToken(*updatedTrainerStats, w.Header())

	_, err = w.Write(toSend)
	if err != nil {
		utils.LogAndSendHTTPError(&w, wrapUpdateTrainerError(err), http.StatusInternalServerError)
	}
}

func addPokemonToTrainer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	trainerUsername := vars[api.UsernameVar]

	pokemon := pokemons.Pokemon{}
	err := json.NewDecoder(r.Body).Decode(&pokemon)
	if err != nil {
		utils.LogAndSendHTTPError(&w, wrapAddPokemonToTrainerError(err), http.StatusInternalServerError)
		return
	}

	pokemonId := primitive.NewObjectID()
	pokemon.Id = pokemonId.Hex()

	updatedPokemons, err := trainerdb.AddPokemonToTrainer(trainerUsername, pokemon)
	if err != nil {
		utils.LogAndSendHTTPError(&w, wrapAddPokemonToTrainerError(err), http.StatusInternalServerError)
		return
	}

	tokens.AddPokemonsTokens(updatedPokemons, w.Header())

	toSend, err := json.Marshal(updatedPokemons[pokemonId.Hex()])
	if err != nil {
		utils.LogAndSendHTTPError(&w, wrapAddPokemonToTrainerError(err), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(toSend)
	if err != nil {
		utils.LogAndSendHTTPError(&w, wrapAddPokemonToTrainerError(err), http.StatusInternalServerError)
	}
}

func handleUpdatePokemon(w http.ResponseWriter, r *http.Request) {
	log.Info("Update pokemon request")

	vars := mux.Vars(r)
	trainerUsername := vars[api.UsernameVar]

	pokemonId, err := primitive.ObjectIDFromHex(vars[api.PokemonIdVar])
	if err != nil {
		utils.LogAndSendHTTPError(&w, wrapUpdatePokemonError(err), http.StatusBadRequest)
		return
	}

	pokemon := pokemons.Pokemon{}
	err = json.NewDecoder(r.Body).Decode(&pokemon)
	if err != nil {
		utils.LogAndSendHTTPError(&w, wrapUpdatePokemonError(err), http.StatusInternalServerError)
		return
	}

	pokemon.Level = experience.CalculateLevel(pokemon.XP)

	updatedPokemons, err := trainerdb.UpdateTrainerPokemon(trainerUsername, pokemonId, pokemon)
	if err != nil {
		utils.LogAndSendHTTPError(&w, wrapUpdatePokemonError(err), http.StatusInternalServerError)
		return
	}

	marshaledPokemons, err := json.Marshal(updatedPokemons[pokemonId.Hex()])
	if err != nil {
		utils.LogAndSendHTTPError(&w, wrapUpdatePokemonError(err), http.StatusInternalServerError)
		return
	}

	tokens.AddPokemonsTokens(updatedPokemons, w.Header())

	_, err = w.Write(marshaledPokemons)
	if err != nil {
		utils.LogAndSendHTTPError(&w, wrapUpdatePokemonError(err), http.StatusInternalServerError)
	}
}

func removePokemonFromTrainer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	trainerUsername := vars[api.UsernameVar]
	pokemonId, err := primitive.ObjectIDFromHex(vars[api.PokemonIdVar])
	if err != nil {
		utils.LogAndSendHTTPError(&w, wrapRemovePokemonFromTrainerError(err), http.StatusBadRequest)
		return
	}

	oldTrainerPokemons, err := trainerdb.RemovePokemonFromTrainer(trainerUsername, pokemonId.Hex())
	if err != nil {
		utils.LogAndSendHTTPError(&w, wrapRemovePokemonFromTrainerError(err), http.StatusInternalServerError)
		return
	}

	removedPokemon, ok := oldTrainerPokemons[pokemonId.Hex()]
	if !ok {
		err = wrapRemovePokemonFromTrainerError(newPokemonNotFoundError(pokemonId.Hex()))
		utils.LogAndSendHTTPError(&w, err, http.StatusNotFound)
		return
	}

	delete(oldTrainerPokemons, pokemonId.Hex())
	tokens.AddPokemonsTokens(oldTrainerPokemons, w.Header())
	toSend, err := json.Marshal(removedPokemon)

	_, err = w.Write(toSend)
	if err != nil {
		utils.LogAndSendHTTPError(&w, wrapRemovePokemonFromTrainerError(err), http.StatusInternalServerError)
	}
}

func addItemsToTrainer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	_ = vars[api.UsernameVar]

	token, err := tokens.ExtractAndVerifyAuthToken(r.Header)
	if err != nil {
		utils.LogAndSendHTTPError(&w, wrapAddItemsToTrainerError(err), http.StatusUnauthorized)
		return
	}

	var itemsToAdd []items.Item
	err = json.NewDecoder(r.Body).Decode(&itemsToAdd)
	if err != nil {
		utils.LogAndSendHTTPError(&w, wrapAddItemsToTrainerError(err), http.StatusInternalServerError)
		return
	}

	addedItems := make(map[string]items.Item, len(itemsToAdd))
	for i, item := range itemsToAdd {
		itemId := primitive.NewObjectID()
		idHex := itemId.Hex()
		item.Id = idHex
		itemsToAdd[i].Id = idHex
		addedItems[idHex] = item
	}

	updatedItems, err := trainerdb.AddItemsToTrainer(token.Username, itemsToAdd)
	if err != nil {
		log.Info("token username: ", token.Username)
		utils.LogAndSendHTTPError(&w, wrapAddItemsToTrainerError(err), http.StatusInternalServerError)
		return
	}

	tokens.AddItemsToken(updatedItems, w.Header())
	toSend, err := json.Marshal(addedItems)
	if err != nil {
		utils.LogAndSendHTTPError(&w, wrapAddItemsToTrainerError(err), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(toSend)
	if err != nil {
		utils.LogAndSendHTTPError(&w, wrapAddItemsToTrainerError(err), http.StatusInternalServerError)
	}
}

func removeItemsFromTrainer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	_ = vars[api.UsernameVar]

	token, err := tokens.ExtractAndVerifyAuthToken(r.Header)
	if err != nil {
		utils.LogAndSendHTTPError(&w, wrapRemoveItemsFromTrainerError(err), http.StatusUnauthorized)
		return
	}

	var itemIds []primitive.ObjectID
	for _, itemIdStr := range strings.Split(vars[api.ItemIdVar], ",") {
		var itemId primitive.ObjectID
		itemId, err = primitive.ObjectIDFromHex(itemIdStr)
		if err != nil {
			utils.LogAndSendHTTPError(&w, wrapRemoveItemsFromTrainerError(err), http.StatusBadRequest)
			return
		}

		itemIds = append(itemIds, itemId)
	}

	oldTrainerItems, err := trainerdb.RemoveItemsFromTrainer(token.Username, itemIds)
	if err != nil {
		utils.LogAndSendHTTPError(&w, wrapRemoveItemsFromTrainerError(err), http.StatusInternalServerError)
		return
	}

	removedItems := make(map[string]items.Item, len(itemIds))
	for i := 0; i < len(itemIds); i++ {
		item, ok := oldTrainerItems[itemIds[i].Hex()]
		if ok {
			removedItems[item.Id] = item
			delete(oldTrainerItems, itemIds[i].Hex())
		}
	}

	toSend, err := json.Marshal(removedItems)
	if err != nil {
		utils.LogAndSendHTTPError(&w, wrapRemoveItemsFromTrainerError(err), http.StatusInternalServerError)
		return
	}

	tokens.AddItemsToken(oldTrainerItems, w.Header())

	_, err = w.Write(toSend)
	if err != nil {
		utils.LogAndSendHTTPError(&w, wrapRemoveItemsFromTrainerError(err), http.StatusInternalServerError)
	}
}

func handleVerifyTrainerPokemons(w http.ResponseWriter, r *http.Request) {
	log.Info("Verify Pokemons request")

	token, err := tokens.ExtractAndVerifyAuthToken(r.Header)
	if err != nil {
		utils.LogAndSendHTTPError(&w, wrapVerifyTrainerPokemonsError(err), http.StatusUnauthorized)
		return
	}

	var receivedHashes map[string]string
	err = json.NewDecoder(r.Body).Decode(&receivedHashes)
	if err != nil {
		utils.LogAndSendHTTPError(&w, wrapVerifyTrainerPokemonsError(err), http.StatusInternalServerError)
		return
	}

	trainer, err := trainerdb.GetTrainerByUsername(token.Username)
	if err != nil {
		utils.LogAndSendHTTPError(&w, wrapVerifyTrainerPokemonsError(err), http.StatusInternalServerError)
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

		if !(currHash == tokens.GenerateHash(pokemon)) {
			log.Info("Denied")
			w.WriteHeader(200)
			toSend, _ := json.Marshal(false)
			_, _ = w.Write(toSend)
			return
		}
	}
	log.Info("Accepted")

	toSend, err := json.Marshal(true)
	if err != nil {
		utils.LogAndSendHTTPError(&w, wrapVerifyTrainerPokemonsError(err), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(toSend)
	if err != nil {
		utils.LogAndSendHTTPError(&w, wrapVerifyTrainerPokemonsError(err), http.StatusInternalServerError)
	}
}

func handleVerifyTrainerStats(w http.ResponseWriter, r *http.Request) {
	log.Info("Verify Trainer Stats request")

	token, err := tokens.ExtractAndVerifyAuthToken(r.Header)
	if err != nil {
		utils.LogAndSendHTTPError(&w, wrapVerifyTrainerStatsError(err), http.StatusUnauthorized)
		return
	}

	var receivedHash string
	err = json.NewDecoder(r.Body).Decode(&receivedHash)
	if err != nil {
		utils.LogAndSendHTTPError(&w, wrapVerifyTrainerStatsError(err), http.StatusInternalServerError)
		return
	}

	trainer, err := trainerdb.GetTrainerByUsername(token.Username)
	if err != nil {
		utils.LogAndSendHTTPError(&w, wrapVerifyTrainerStatsError(err), http.StatusInternalServerError)
		return
	}

	equal := tokens.GenerateHash(trainer.Stats) == receivedHash
	toSend, err := json.Marshal(equal)
	if err != nil {
		utils.LogAndSendHTTPError(&w, wrapVerifyTrainerStatsError(err), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(toSend)
	if err != nil {
		utils.LogAndSendHTTPError(&w, wrapVerifyTrainerStatsError(err), http.StatusInternalServerError)
	}
}

func handleVerifyTrainerItems(w http.ResponseWriter, r *http.Request) {
	log.Info("Verify items request")

	token, err := tokens.ExtractAndVerifyAuthToken(r.Header)
	if err != nil {
		utils.LogAndSendHTTPError(&w, wrapVerifyTrainerItemsError(err), http.StatusUnauthorized)
		return
	}

	var receivedHash string
	err = json.NewDecoder(r.Body).Decode(&receivedHash)
	if err != nil {
		utils.LogAndSendHTTPError(&w, wrapVerifyTrainerItemsError(err), http.StatusInternalServerError)
		return
	}

	trainer, err := trainerdb.GetTrainerByUsername(token.Username)
	if err != nil {
		utils.LogAndSendHTTPError(&w, wrapVerifyTrainerItemsError(err), http.StatusInternalServerError)
		return
	}

	equal := tokens.GenerateHash(trainer.Items) == receivedHash
	log.Info("verify items: ", equal)

	toSend, err := json.Marshal(equal)
	if err != nil {
		utils.LogAndSendHTTPError(&w, wrapVerifyTrainerItemsError(err), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(toSend)
	if err != nil {
		utils.LogAndSendHTTPError(&w, wrapVerifyTrainerItemsError(err), http.StatusInternalServerError)
		return
	}
}

func handleGenerateAllTokens(w http.ResponseWriter, r *http.Request) {
	log.Info("Generate all tokens request")
	token, err := tokens.ExtractAndVerifyAuthToken(r.Header)
	if err != nil {
		utils.LogAndSendHTTPError(&w, wrapGenerateAllTokensError(err), http.StatusUnauthorized)
		return
	}

	trainer, err := trainerdb.GetTrainerByUsername(token.Username)
	if err != nil {
		utils.LogAndSendHTTPError(&w, wrapGenerateAllTokensError(err), http.StatusInternalServerError)
		return
	}

	tokens.AddItemsToken(trainer.Items, w.Header())
	tokens.AddPokemonsTokens(trainer.Pokemons, w.Header())
	tokens.AddTrainerStatsToken(trainer.Stats, w.Header())
}

func handleGenerateTrainerStatsToken(w http.ResponseWriter, r *http.Request) {
	token, err := tokens.ExtractAndVerifyAuthToken(r.Header)
	if err != nil {
		utils.LogAndSendHTTPError(&w, wrapGenerateStatsTokenError(err), http.StatusUnauthorized)
		return
	}

	trainer, err := trainerdb.GetTrainerByUsername(token.Username)
	if err != nil {
		utils.LogAndSendHTTPError(&w, wrapGenerateStatsTokenError(err), http.StatusInternalServerError)
		return
	}

	tokens.AddTrainerStatsToken(trainer.Stats, w.Header())
}

func handleGeneratePokemonsToken(w http.ResponseWriter, r *http.Request) {
	token, err := tokens.ExtractAndVerifyAuthToken(r.Header)
	if err != nil {
		utils.LogAndSendHTTPError(&w, wrapGeneratePokemonsTokenError(err), http.StatusUnauthorized)
		return
	}

	trainer, err := trainerdb.GetTrainerByUsername(token.Username)
	if err != nil {
		utils.LogAndSendHTTPError(&w, wrapGeneratePokemonsTokenError(err), http.StatusInternalServerError)
		return
	}

	tokens.AddPokemonsTokens(trainer.Pokemons, w.Header())
}

func handleGenerateItemsToken(w http.ResponseWriter, r *http.Request) {
	token, err := tokens.ExtractAndVerifyAuthToken(r.Header)
	if err != nil {
		utils.LogAndSendHTTPError(&w, wrapGenerateItemsTokenError(err), http.StatusUnauthorized)
		return
	}

	trainer, err := trainerdb.GetTrainerByUsername(token.Username)
	if err != nil {
		utils.LogAndSendHTTPError(&w, wrapGenerateItemsTokenError(err), http.StatusInternalServerError)
		return
	}

	tokens.AddItemsToken(trainer.Items, w.Header())
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
		toAdd.Id = primitive.NewObjectID().Hex()
		starterItems[toAdd.Id] = toAdd
	}

	for i := 0; i < masterBallsAmount; i++ {
		toAdd := items.MasterBallItem
		toAdd.Id = primitive.NewObjectID().Hex()
		starterItems[toAdd.Id] = toAdd
	}

	for i := 0; i < healAmount; i++ {
		toAdd := items.HealItem
		toAdd.Id = primitive.NewObjectID().Hex()
		starterItems[toAdd.Id] = toAdd
	}

	for i := 0; i < reviveAmount; i++ {
		toAdd := items.ReviveItem
		toAdd.Id = primitive.NewObjectID().Hex()
		starterItems[toAdd.Id] = toAdd
	}

	return starterItems
}
