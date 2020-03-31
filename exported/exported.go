package exported

import "fmt"

const UsernamePathVar = "username"
const PokemonIdPathVar = "%s, PokemonIdPathVar"
const ItemIdPathVar = "itemId"

var AddTrainerPath = "/trainers/"
var GetTrainersPath = "/trainers/"

var GetTrainerByUsernamePath = fmt.Sprintf("/trainers/{%s}", UsernamePathVar)
var UpdateTrainerInfoPath = fmt.Sprintf("/trainers/{%s}", UsernamePathVar)

var AddPokemonPath = fmt.Sprintf("/trainers/{%s}/pokemons/", UsernamePathVar)
var RemovePokemonPath = fmt.Sprintf("/trainers/{%s}/pokemons/{%s}", UsernamePathVar, PokemonIdPathVar)

var AddItemToBagPath = fmt.Sprintf("/trainers/{%s}/bag/", UsernamePathVar)
var RemoveItemFromBagPath = fmt.Sprintf("/trainers/{%s}/bag/{%s}", UsernamePathVar, ItemIdPathVar)

var VerifyTrainerStatsPath = fmt.Sprintf("/trainers/{%s}/stats/verify", UsernamePathVar)
var VerifyPokemonsPath = fmt.Sprintf("/trainers/{%s}/pokemons/verify", UsernamePathVar)
var VerifyItemsPath = fmt.Sprintf("/trainers/{%s}/bag/verify", UsernamePathVar)

var GenerateAllTokensPath = fmt.Sprintf("/trainers/{%s}/tokens", UsernamePathVar)
var GenerateTrainerStatsTokenPath = fmt.Sprintf("/trainers/{%s}/stats/token", UsernamePathVar)
var GenerateItemsTokenPath = fmt.Sprintf("/trainers/{%s}/items/token", UsernamePathVar)
var GeneratePokemonsTokenPath = fmt.Sprintf("/trainers/{%s}/pokemons/token", UsernamePathVar)
