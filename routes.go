package main

import (
	"github.com/NOVAPokemon/utils"
)

const GET = "GET"
const POST = "POST"
const DELETE = "DELETE"
const PUT = "PUT"

// trainer
const AddTrainerPath = "/trainers/"
const AddTrainerName = "CREATE_TRAINER"

const GetTrainersPath = "/trainers/"
const GetTrainersName = "GET_TRAINERS"

const GetTrainerByUsernamePath = "/trainers/{username}"
const GetTrainerByUsernameName = "GET_TRAINER"

// trainer pokemons
const AddPokemonPath = "/trainers/{username}/pokemons/"
const AddPokemonName = "ADD_POKEMON"

const RemovePokemonPath = "/trainers/{username}/pokemons/{pokemonId}"
const RemovePokemonName = "REMOVE_POKEMON"

// trainer bag
const AddItemToBagPath = "/trainers/{username}/bag/"
const AddItemToBagName = "ADD_TO_BAG"

const RemoveItemFromBagPath = "/trainers/{username}/bag/{itemId}"
const RemoveItemFromBagName = "REMOVE_FROM_BAG"

// token

const VerifyTokenName = "/token"

var routes = utils.Routes{
	// TRAINERS

	utils.Route{
		Name:        AddTrainerName,
		Method:      GET,
		Pattern:     AddTrainerPath,
		HandlerFunc: AddTrainer,
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

	// TOKEN

	utils.Route{
		Name:        VerifyTokenName,
		Method:      GET,
		Pattern:     VerifyTokenPath,
		HandlerFunc: HandleVerifyTrainerToken,
	},
}
