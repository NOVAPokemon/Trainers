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
	return errors.Wrap(err, fmt.Sprintf(utils.ErrorInHandlerFormat, GetTrainersName))
}

func wrapGetTrainerError(err error) error {
	return errors.Wrap(err, fmt.Sprintf(utils.ErrorInHandlerFormat, GetTrainerByUsernameName))
}

func wrapAddTrainerError(err error) error {
	return errors.Wrap(err, fmt.Sprintf(utils.ErrorInHandlerFormat, AddTrainerName))
}

func wrapUpdateTrainerError(err error) error {
	return errors.Wrap(err, fmt.Sprintf(utils.ErrorInHandlerFormat, UpdateTrainerInfo))
}

func wrapAddPokemonToTrainerError(err error) error {
	return errors.Wrap(err, fmt.Sprintf(utils.ErrorInHandlerFormat, AddPokemonName))
}

func wrapUpdatePokemonError(err error) error {
	return errors.Wrap(err, fmt.Sprintf(utils.ErrorInHandlerFormat, UpdatePokemonName))
}

func wrapRemovePokemonFromTrainerError(err error) error {
	return errors.Wrap(err, fmt.Sprintf(utils.ErrorInHandlerFormat, RemovePokemonName))
}

func wrapAddItemsToTrainerError(err error) error {
	return errors.Wrap(err, fmt.Sprintf(utils.ErrorInHandlerFormat, AddItemsName))
}

func wrapRemoveItemsFromTrainerError(err error) error {
	return errors.Wrap(err, fmt.Sprintf(utils.ErrorInHandlerFormat, RemoveItemsName))
}

func wrapVerifyTrainerPokemonsError(err error) error {
	return errors.Wrap(err, fmt.Sprintf(utils.ErrorInHandlerFormat, VerifyPokemon))
}

func wrapVerifyTrainerStatsError(err error) error {
	return errors.Wrap(err, fmt.Sprintf(utils.ErrorInHandlerFormat, VerifyTrainerStats))
}

func wrapVerifyTrainerItemsError(err error) error {
	return errors.Wrap(err, fmt.Sprintf(utils.ErrorInHandlerFormat, VerifyItems))
}

func wrapGenerateAllTokensError(err error) error {
	return errors.Wrap(err, fmt.Sprintf(utils.ErrorInHandlerFormat, GenerateAllTokens))
}

func wrapGenerateStatsTokenError(err error) error {
	return errors.Wrap(err, fmt.Sprintf(utils.ErrorInHandlerFormat, GenerateTrainerStatsToken))
}

func wrapGeneratePokemonsTokenError(err error) error {
	return errors.Wrap(err, fmt.Sprintf(utils.ErrorInHandlerFormat, GeneratePokemonsToken))
}

func wrapGenerateItemsTokenError(err error) error {
	return errors.Wrap(err, fmt.Sprintf(utils.ErrorInHandlerFormat, GenerateItemsToken))
}

// Error builders
func newPokemonNotFoundError(pokemonId string) error {
	return errors.New(fmt.Sprintf(errorPokemonNotFoundFormat, pokemonId))
}
