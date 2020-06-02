package main

import (
	"fmt"
	"github.com/NOVAPokemon/utils"
	"github.com/NOVAPokemon/utils/api"
	"strings"
)

const GET = "GET"
const POST = "POST"
const DELETE = "DELETE"
const PUT = "PUT"

// trainer
const AddTrainerName = "CREATE_TRAINER"
const GetTrainersName = "GET_TRAINERS"

const GetTrainerByUsernameName = "GET_TRAINER"
const UpdateTrainerInfo = "UPDATE_TRAINER_INFO"

// trainer pokemons
const AddPokemonName = "ADD_POKEMON"
const UpdatePokemonName = "UPDATE_POKEMON"
const RemovePokemonName = "REMOVE_POKEMON"

// trainer bag
const AddItemsName = "ADD_ITEMS"
const RemoveItemsName = "REMOVE_ITEMS"

// tokens
const VerifyTrainerStats = "VERIFY_STATS"
const VerifyPokemon = "VERIFY_POKEMONS"
const VerifyItems = "VERIFY_ITEMS"
const GenerateAllTokens = "GENERATE_ALL_TOKENS"
const GenerateTrainerStatsToken = "GENERATE_TRAINER_STATS_TOKEN"
const GenerateItemsToken = "GENERATE_TRAINER_ITEMS_TOKEN"
const GeneratePokemonsToken = "GENERATE_TRAINER_POKEMONS_TOKEN"

var routes = utils.Routes{
	api.GenStatusRoute(strings.ToLower(fmt.Sprintf("/%s", serviceName))),
	utils.Route{
		Name:        AddTrainerName,
		Method:      POST,
		Pattern:     api.AddTrainerPath,
		HandlerFunc: AddTrainer,
	},

	utils.Route{
		Name:        UpdateTrainerInfo,
		Method:      PUT,
		Pattern:     api.UpdateTrainerStatsRoute,
		HandlerFunc: HandleUpdateTrainerInfo,
	},

	utils.Route{
		Name:        GetTrainersName,
		Method:      GET,
		Pattern:     api.GetTrainersPath,
		HandlerFunc: GetAllTrainers,
	},

	utils.Route{
		Name:        GetTrainerByUsernameName,
		Method:      GET,
		Pattern:     api.GetTrainerByUsernameRoute,
		HandlerFunc: GetTrainerByUsername,
	},

	// POKEMONS

	utils.Route{
		Name:        AddPokemonName,
		Method:      POST,
		Pattern:     api.AddPokemonRoute,
		HandlerFunc: AddPokemonToTrainer,
	},

	utils.Route{
		Name:        UpdatePokemonName,
		Method:      PUT,
		Pattern:     api.UpdatePokemonRoute,
		HandlerFunc: HandleUpdatePokemon,
	},

	utils.Route{
		Name:        RemovePokemonName,
		Method:      DELETE,
		Pattern:     api.RemovePokemonRoute,
		HandlerFunc: RemovePokemonFromTrainer,
	},

	// ITEMS

	utils.Route{
		Name:        AddItemsName,
		Method:      POST,
		Pattern:     api.AddItemToBagRoute,
		HandlerFunc: AddItemsToTrainer,
	},

	utils.Route{
		Name:        RemoveItemsName,
		Method:      DELETE,
		Pattern:     api.RemoveItemFromBagRoute,
		HandlerFunc: RemoveItemsFromTrainer,
	},

	// TOKENS

	utils.Route{
		Name:        GenerateAllTokens,
		Method:      GET,
		Pattern:     api.GenerateAllTokensRoute,
		HandlerFunc: HandleGenerateAllTokens,
	},

	utils.Route{
		Name:        GenerateTrainerStatsToken,
		Method:      GET,
		Pattern:     api.GenerateTrainerStatsTokenRoute,
		HandlerFunc: HandleGenerateTrainerStatsToken,
	},

	utils.Route{
		Name:        GeneratePokemonsToken,
		Method:      GET,
		Pattern:     api.GeneratePokemonsTokenRoute,
		HandlerFunc: HandleGeneratePokemonsToken,
	},

	utils.Route{
		Name:        GenerateItemsToken,
		Method:      GET,
		Pattern:     api.GenerateItemsTokenRoute,
		HandlerFunc: HandleGenerateItemsToken,
	},

	utils.Route{
		Name:        VerifyItems,
		Method:      POST,
		Pattern:     api.VerifyItemsRoute,
		HandlerFunc: HandleVerifyTrainerItems,
	},

	utils.Route{
		Name:        VerifyTrainerStats,
		Method:      POST,
		Pattern:     api.VerifyTrainerStatsRoute,
		HandlerFunc: HandleVerifyTrainerStats,
	},

	utils.Route{
		Name:        VerifyPokemon,
		Method:      POST,
		Pattern:     api.VerifyPokemonRoute,
		HandlerFunc: HandleVerifyTrainerPokemons,
	},
}
