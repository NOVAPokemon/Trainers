package main

import (
	"fmt"
	"github.com/NOVAPokemon/utils"
	"github.com/NOVAPokemon/utils/api"
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
const VerifyPokemon = "VERIFY_POKEMONS"
const VerifyItems = "VERIFY_ITEMS"
const GenerateAllTokens = "GENERATE_ALL_TOKENS"
const GenerateTrainerStatsToken = "GENERATE_TRAINER_STATS_TOKEN"
const GenerateItemsToken = "GENERATE_TRAINER_ITEMS_TOKEN"
const GeneratePokemonsToken = "GENERATE_TRAINER_POKEMONS_TOKEN"

// ROUTES
const UsernameVar = "username"
const PokemonIdVar = "pokemonId"
const ItemIdVar = "itemId"

var GetTrainerByUsernameRoute = fmt.Sprintf(api.GetTrainerByUsernamePath, UsernameVar)
var UpdateTrainerStatsRoute = fmt.Sprintf(api.UpdateTrainerStatsPath, UsernameVar)

// trainer pokemons
var AddPokemonRoute = fmt.Sprintf(api.AddPokemonPath, UsernameVar)
var RemovePokemonRoute = fmt.Sprintf(api.RemovePokemonPath, UsernameVar, PokemonIdVar)

// trainer bag
var AddItemToBagRoute = fmt.Sprintf(api.AddItemToBagPath, UsernameVar)
var RemoveItemFromBagRoute = fmt.Sprintf(api.RemoveItemFromBagPath, UsernameVar, ItemIdVar)

// Tokens
var VerifyTrainerStatsRoute = fmt.Sprintf(api.VerifyTrainerStatsPath, UsernameVar)
var VerifyPokemonRoute = fmt.Sprintf(api.VerifyPokemonPath, UsernameVar)
var VerifyItemsRoute = fmt.Sprintf(api.VerifyItemsPath, UsernameVar)

var GenerateAllTokensRoute = fmt.Sprintf(api.GenerateAllTokensPath, UsernameVar)
var GenerateTrainerStatsTokenRoute = fmt.Sprintf(api.GenerateTrainerStatsTokenPath, UsernameVar)
var GenerateItemsTokenRoute = fmt.Sprintf(api.GenerateItemsTokenPath, UsernameVar)
var GeneratePokemonsTokenRoute = fmt.Sprintf(api.GeneratePokemonsTokenPath, UsernameVar)

var routes = utils.Routes{
	// TRAINERS

	utils.Route{
		Name:        AddTrainerName,
		Method:      GET,
		Pattern:     api.AddTrainerPath,
		HandlerFunc: AddTrainer,
	},

	utils.Route{
		Name:        UpdateTrainerInfo,
		Method:      PUT,
		Pattern:     UpdateTrainerStatsRoute,
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
		Pattern:     GetTrainerByUsernameRoute,
		HandlerFunc: GetTrainerByUsername,
	},

	// POKEMONS

	utils.Route{
		Name:        AddPokemonName,
		Method:      POST,
		Pattern:     AddPokemonRoute,
		HandlerFunc: AddPokemonToTrainer,
	},

	utils.Route{
		Name:        RemovePokemonName,
		Method:      DELETE,
		Pattern:     RemovePokemonRoute,
		HandlerFunc: RemovePokemonFromTrainer,
	},

	// ITEMS

	utils.Route{
		Name:        AddItemToBagName,
		Method:      POST,
		Pattern:     AddItemToBagRoute,
		HandlerFunc: AddItemsToTrainer,
	},

	utils.Route{
		Name:        RemoveItemFromBagName,
		Method:      DELETE,
		Pattern:     RemoveItemFromBagRoute,
		HandlerFunc: RemoveItemsFromTrainer,
	},

	// TOKENS

	utils.Route{
		Name:        GenerateAllTokens,
		Method:      GET,
		Pattern:     GenerateAllTokensRoute,
		HandlerFunc: HandleGenerateAllTokens,
	},

	utils.Route{
		Name:        GenerateTrainerStatsToken,
		Method:      GET,
		Pattern:     GenerateTrainerStatsTokenRoute,
		HandlerFunc: HandleGenerateTrainerStatsToken,
	},

	utils.Route{
		Name:        GeneratePokemonsToken,
		Method:      GET,
		Pattern:     GeneratePokemonsTokenRoute,
		HandlerFunc: HandleGeneratePokemonsToken,
	},

	utils.Route{
		Name:        GenerateItemsToken,
		Method:      GET,
		Pattern:     GenerateItemsTokenRoute,
		HandlerFunc: HandleGenerateItemsToken,
	},

	utils.Route{
		Name:        VerifyItems,
		Method:      POST,
		Pattern:     VerifyItemsRoute,
		HandlerFunc: HandleVerifyTrainerItems,
	},

	utils.Route{
		Name:        VerifyTrainerStats,
		Method:      POST,
		Pattern:     VerifyTrainerStatsRoute,
		HandlerFunc: HandleVerifyTrainerStats,
	},

	utils.Route{
		Name:        VerifyPokemon,
		Method:      POST,
		Pattern:     VerifyPokemonRoute,
		HandlerFunc: HandleVerifyTrainerPokemon,
	},
}
