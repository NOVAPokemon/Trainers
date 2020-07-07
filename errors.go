package main

import (
	"fmt"

	"github.com/NOVAPokemon/utils"
	"github.com/pkg/errors"
)

const (
	errorPokemonNotFoundFormat = "could not find pokemon %s"
)

// Handler wrappers
func wrapGetAllTrainersError(err error) error {
	return errors.Wrap(err, fmt.Sprintf(utils.ErrorInHandlerFormat, getTrainersName))
}

func wrapGetTrainerError(err error) error {
	return errors.Wrap(err, fmt.Sprintf(utils.ErrorInHandlerFormat, getTrainerByUsernameName))
}

func wrapAddTrainerError(err error) error {
	return errors.Wrap(err, fmt.Sprintf(utils.ErrorInHandlerFormat, addTrainerName))
}

func wrapUpdateTrainerError(err error) error {
	return errors.Wrap(err, fmt.Sprintf(utils.ErrorInHandlerFormat, updateTrainerInfo))
}

func wrapAddPokemonToTrainerError(err error) error {
	return errors.Wrap(err, fmt.Sprintf(utils.ErrorInHandlerFormat, addPokemonName))
}

func wrapUpdatePokemonError(err error) error {
	return errors.Wrap(err, fmt.Sprintf(utils.ErrorInHandlerFormat, updatePokemonName))
}

func wrapRemovePokemonFromTrainerError(err error) error {
	return errors.Wrap(err, fmt.Sprintf(utils.ErrorInHandlerFormat, removePokemonName))
}

func wrapAddItemsToTrainerError(err error) error {
	return errors.Wrap(err, fmt.Sprintf(utils.ErrorInHandlerFormat, addItemsName))
}

func wrapRemoveItemsFromTrainerError(err error) error {
	return errors.Wrap(err, fmt.Sprintf(utils.ErrorInHandlerFormat, removeItemsName))
}

func wrapVerifyTrainerPokemonsError(err error) error {
	return errors.Wrap(err, fmt.Sprintf(utils.ErrorInHandlerFormat, verifyPokemon))
}

func wrapVerifyTrainerStatsError(err error) error {
	return errors.Wrap(err, fmt.Sprintf(utils.ErrorInHandlerFormat, verifyTrainerStats))
}

func wrapVerifyTrainerItemsError(err error) error {
	return errors.Wrap(err, fmt.Sprintf(utils.ErrorInHandlerFormat, verifyItems))
}

func wrapGenerateAllTokensError(err error) error {
	return errors.Wrap(err, fmt.Sprintf(utils.ErrorInHandlerFormat, generateAllTokens))
}

func wrapGenerateStatsTokenError(err error) error {
	return errors.Wrap(err, fmt.Sprintf(utils.ErrorInHandlerFormat, generateTrainerStatsToken))
}

func wrapGeneratePokemonsTokenError(err error) error {
	return errors.Wrap(err, fmt.Sprintf(utils.ErrorInHandlerFormat, generatePokemonsToken))
}

func wrapGenerateItemsTokenError(err error) error {
	return errors.Wrap(err, fmt.Sprintf(utils.ErrorInHandlerFormat, generateItemsToken))
}

// Error builders
func newPokemonNotFoundError(pokemonId string) error {
	return errors.New(fmt.Sprintf(errorPokemonNotFoundFormat, pokemonId))
}
