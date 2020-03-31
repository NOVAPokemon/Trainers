package main

import (
	"fmt"
	"github.com/NOVAPokemon/trainers/exported"
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

// ROUTES
const UsernameVar = "username"
const PokemonIdVar = "pokemonId"
const ItemIdVar = "itemId"

var GetTrainerByUsernameRoute = fmt.Sprintf("/trainers/{%s}", UsernameVar)
var UpdateTrainerInfoRoute = fmt.Sprintf("/trainers/{%s}", UsernameVar)

// trainer pokemons
var AddPokemonRoute = fmt.Sprintf("/trainers/{%s}/pokemons/", UsernameVar)
var RemovePokemonRoute = fmt.Sprintf("/trainers/{%s}/pokemons/{%s}", UsernameVar, PokemonIdVar)

// trainer bag
var AddItemToBagRoute = fmt.Sprintf("/trainers/{%s}/bag/", UsernameVar)
var RemoveItemFromBagRoute = fmt.Sprintf("/trainers/{%s}/bag/{%s}", UsernameVar, ItemIdVar)

// Tokens
var VerifyTrainerStatsRoute = fmt.Sprintf("/trainers/{%s}/stats/verify", UsernameVar)
var VerifyPokemonsRoute = fmt.Sprintf("/trainers/{%s}/pokemons/verify", UsernameVar)
var VerifyItemsRoute = fmt.Sprintf("/trainers/{%s}/bag/verify", UsernameVar)
var GenerateAllTokensRoute = fmt.Sprintf("/trainers/{%s}/tokens", UsernameVar)
var GenerateTrainerStatsTokenRoute = fmt.Sprintf("/trainers/{%s}/stats/token", UsernameVar)
var GenerateItemsTokenRoute = fmt.Sprintf("/trainers/{%s}/items/token", UsernameVar)
var GeneratePokemonsTokenRoute = fmt.Sprintf("/trainers/{%s}/pokemons/token", UsernameVar)

var routes = utils.Routes{
	// TRAINERS

	utils.Route{
		Name:        AddTrainerName,
		Method:      GET,
		Pattern:     exported.AddTrainerPath,
		HandlerFunc: AddTrainer,
	},

	utils.Route{
		Name:        UpdateTrainerInfo,
		Method:      PUT,
		Pattern:     UpdateTrainerInfoRoute,
		HandlerFunc: HandleUpdateTrainerInfo,
	},

	utils.Route{
		Name:        GetTrainersName,
		Method:      GET,
		Pattern:     exported.GetTrainersPath,
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
		Name:        VerifyTrainerStats,
		Method:      POST,
		Pattern:     VerifyTrainerStatsRoute,
		HandlerFunc: HandleVerifyTrainerStats,
	},

	utils.Route{
		Name:        VerifyPokemons,
		Method:      POST,
		Pattern:     VerifyPokemonsRoute,
		HandlerFunc: HandleVerifyTrainerPokemons,
	},

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
		Name:        VerifyPokemons,
		Method:      POST,
		Pattern:     VerifyPokemonsRoute,
		HandlerFunc: HandleVerifyTrainerPokemons,
	},

	utils.Route{
		Name:        VerifyItems,
		Method:      POST,
		Pattern:     VerifyItemsRoute,
		HandlerFunc: HandleVerifyTrainerItems,
	},
}
