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

var routes = utils.Routes{

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

	utils.Route{
		Name:        AddItemToBagName,
		Method:      POST,
		Pattern:     AddItemToBagPath,
		HandlerFunc: AddItemToTrainer,
	},

	utils.Route{
		Name:        RemoveItemFromBagName,
		Method:      DELETE,
		Pattern:     RemoveItemFromBagPath,
		HandlerFunc: RemoveItemToTrainer,
	},
}
