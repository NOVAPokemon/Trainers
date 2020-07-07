package main

import (
	"fmt"
	"strings"

	"github.com/NOVAPokemon/utils"
	"github.com/NOVAPokemon/utils/api"
)

const (
	// trainer
	addTrainerName           = "CREATE_TRAINER"
	getTrainersName          = "GET_TRAINERS"
	getTrainerByUsernameName = "GET_TRAINER"
	updateTrainerInfo        = "UPDATE_TRAINER_INFO"

	// trainer pokemons
	addPokemonName    = "ADD_POKEMON"
	updatePokemonName = "UPDATE_POKEMON"
	removePokemonName = "REMOVE_POKEMON"

	// trainer bag
	addItemsName    = "ADD_ITEMS"
	removeItemsName = "REMOVE_ITEMS"

	// tokens
	verifyTrainerStats        = "VERIFY_STATS"
	verifyPokemon             = "VERIFY_POKEMONS"
	verifyItems               = "VERIFY_ITEMS"
	generateAllTokens         = "GENERATE_ALL_TOKENS"
	generateTrainerStatsToken = "GENERATE_TRAINER_STATS_TOKEN"
	generateItemsToken        = "GENERATE_TRAINER_ITEMS_TOKEN"
	generatePokemonsToken     = "GENERATE_TRAINER_POKEMONS_TOKEN"
)

const (
	get        = "GET"
	post       = "POST"
	deleteVerb = "DELETE"
	put        = "PUT"
)

var routes = utils.Routes{
	api.GenStatusRoute(strings.ToLower(fmt.Sprintf("%s", serviceName))),
	utils.Route{
		Name:        addTrainerName,
		Method:      post,
		Pattern:     api.AddTrainerPath,
		HandlerFunc: addTrainer,
	},

	utils.Route{
		Name:        updateTrainerInfo,
		Method:      put,
		Pattern:     api.UpdateTrainerStatsRoute,
		HandlerFunc: handleUpdateTrainerInfo,
	},

	utils.Route{
		Name:        getTrainersName,
		Method:      get,
		Pattern:     api.GetTrainersPath,
		HandlerFunc: getAllTrainers,
	},

	utils.Route{
		Name:        getTrainerByUsernameName,
		Method:      get,
		Pattern:     api.GetTrainerByUsernameRoute,
		HandlerFunc: getTrainerByUsername,
	},

	// POKEMONS

	utils.Route{
		Name:        addPokemonName,
		Method:      post,
		Pattern:     api.AddPokemonRoute,
		HandlerFunc: addPokemonToTrainer,
	},

	utils.Route{
		Name:        updatePokemonName,
		Method:      put,
		Pattern:     api.UpdatePokemonRoute,
		HandlerFunc: handleUpdatePokemon,
	},

	utils.Route{
		Name:        removePokemonName,
		Method:      deleteVerb,
		Pattern:     api.RemovePokemonRoute,
		HandlerFunc: removePokemonFromTrainer,
	},

	// ITEMS

	utils.Route{
		Name:        addItemsName,
		Method:      post,
		Pattern:     api.AddItemToBagRoute,
		HandlerFunc: addItemsToTrainer,
	},

	utils.Route{
		Name:        removeItemsName,
		Method:      deleteVerb,
		Pattern:     api.RemoveItemFromBagRoute,
		HandlerFunc: removeItemsFromTrainer,
	},

	// TOKENS

	utils.Route{
		Name:        generateAllTokens,
		Method:      get,
		Pattern:     api.GenerateAllTokensRoute,
		HandlerFunc: handleGenerateAllTokens,
	},

	utils.Route{
		Name:        generateTrainerStatsToken,
		Method:      get,
		Pattern:     api.GenerateTrainerStatsTokenRoute,
		HandlerFunc: handleGenerateTrainerStatsToken,
	},

	utils.Route{
		Name:        generatePokemonsToken,
		Method:      get,
		Pattern:     api.GeneratePokemonsTokenRoute,
		HandlerFunc: handleGeneratePokemonsToken,
	},

	utils.Route{
		Name:        generateItemsToken,
		Method:      get,
		Pattern:     api.GenerateItemsTokenRoute,
		HandlerFunc: handleGenerateItemsToken,
	},

	utils.Route{
		Name:        verifyItems,
		Method:      post,
		Pattern:     api.VerifyItemsRoute,
		HandlerFunc: handleVerifyTrainerItems,
	},

	utils.Route{
		Name:        verifyTrainerStats,
		Method:      post,
		Pattern:     api.VerifyTrainerStatsRoute,
		HandlerFunc: handleVerifyTrainerStats,
	},

	utils.Route{
		Name:        verifyPokemon,
		Method:      post,
		Pattern:     api.VerifyPokemonRoute,
		HandlerFunc: handleVerifyTrainerPokemons,
	},
}
