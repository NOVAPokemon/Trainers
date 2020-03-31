package main

import (
	"github.com/NOVAPokemon/utils"
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

const RemovePokemonName = "REMOVE_POKEMON"

// trainer bag
const AddItemToBagName = "ADD_TO_BAG"

const RemoveItemFromBagName = "REMOVE_FROM_BAG"

// tokens
const VerifyTrainerStats = "VERIFY_STATS"

const VerifyPokemons = "VERIFY_POKEMONS"

const VerifyItems = "VERIFY_ITEMS"

const GenerateAllTokens = "GENERATE_ALL_TOKENS"

const GenerateTrainerStatsToken = "GENERATE_TRAINER_STATS_TOKEN"

const GenerateItemsToken = "GENERATE_TRAINER_ITEMS_TOKEN"

const GeneratePokemonsToken = "GENERATE_TRAINER_POKEMONS_TOKEN"

var routes = utils.Routes{
	// TRAINERS

	utils.Route{
		Name:        AddTrainerName,
		Method:      GET,
		Pattern:     AddTrainerPath,
		HandlerFunc: AddTrainer,
	},


	utils.Route{
		Name:        UpdateTrainerInfo,
		Method:      PUT,
		Pattern:     UpdateTrainerInfoPath,
		HandlerFunc: HandleUpdateTrainerInfo,
	},

	utils.Route{
		Name:        GetTrainersName,
		Method:      GET,
		Pattern:     GetTrainersPath,
		HandlerFunc: GetAllTrainers,
	},

	utils.Route{
		Name:        GetTrainerByUsernameName,
		Method:      GET,
		Pattern:     GetTrainerByUsernamePath,
		HandlerFunc: GetTrainerByUsername,
	},

	// POKEMONS

	utils.Route{
		Name:        AddPokemonName,
		Method:      POST,
		Pattern:     AddPokemonPath,
		HandlerFunc: AddPokemonToTrainer,
	},

	utils.Route{
		Name:        RemovePokemonName,
		Method:      DELETE,
		Pattern:     RemovePokemonPath,
		HandlerFunc: RemovePokemonFromTrainer,
	},

	// ITEMS

	utils.Route{
		Name:        AddItemToBagName,
		Method:      POST,
		Pattern:     AddItemToBagPath,
		HandlerFunc: AddItemsToTrainer,
	},

	utils.Route{
		Name:        RemoveItemFromBagName,
		Method:      DELETE,
		Pattern:     RemoveItemFromBagPath,
		HandlerFunc: RemoveItemsFromTrainer,
	},

	// TOKENS

	utils.Route{
		Name:        VerifyTrainerStats,
		Method:      POST,
		Pattern:     VerifyTrainerStatsPath,
		HandlerFunc: HandleVerifyTrainerStats,
	},

	utils.Route{
		Name:        VerifyPokemons,
		Method:      POST,
		Pattern:     VerifyPokemonsPath,
		HandlerFunc: HandleVerifyTrainerPokemons,
	},

	utils.Route{
		Name:        GenerateAllTokens,
		Method:      GET,
		Pattern:     GenerateAllTokensPath,
		HandlerFunc: HandleGenerateAllTokens,
	},

	utils.Route{
		Name:        GenerateTrainerStatsToken,
		Method:      GET,
		Pattern:     GenerateTrainerStatsTokenPath,
		HandlerFunc: HandleGenerateTrainerStatsToken,
	},

	utils.Route{
		Name:        GeneratePokemonsToken,
		Method:      GET,
		Pattern:     GeneratePokemonsTokenPath,
		HandlerFunc: HandleGeneratePokemonsToken,
	},

	utils.Route{
		Name:        GenerateItemsToken,
		Method:      GET,
		Pattern:     GenerateItemsTokenPath,
		HandlerFunc: HandleGenerateItemsToken,
	},
	utils.Route{
		Name:        VerifyItems,
		Method:      POST,
		Pattern:     VerifyItemsPath,
		HandlerFunc: HandleVerifyTrainerItems,
	},

	utils.Route{
		Name:        VerifyTrainerStats,
		Method:      POST,
		Pattern:     VerifyTrainerStatsPath,
		HandlerFunc: HandleVerifyTrainerStats,
	},

	utils.Route{
		Name:        VerifyPokemons,
		Method:      POST,
		Pattern:     VerifyPokemonsPath,
		HandlerFunc: HandleVerifyTrainerPokemons,
	},

	utils.Route{
		Name:        VerifyItems,
		Method:      POST,
		Pattern:     VerifyItemsPath,
		HandlerFunc: HandleVerifyTrainerItems,
	},
}
