package exported


// PATHS

// trainer
const AddTrainerPath = "/trainers/"
const GetTrainersPath = "/trainers/"

const GetTrainerByUsernamePath = "/trainers/%s"
const UpdateTrainerInfoPath = "/trainers/%s"

// trainer pokemons
const AddPokemonPath = "/trainers/%s/pokemons/"
const RemovePokemonPath = "/trainers/%s/pokemons/%s"

// trainer bag
const AddItemToBagPath = "/trainers/%s/bag/"
const RemoveItemFromBagPath = "/trainers/%s/bag/%s"

// Tokens
const VerifyTrainerStatsPath = "/trainers/%s/stats/verify"
const VerifyPokemonsPath = "/trainers/%s/pokemons/verify"
const VerifyItemsPath = "/trainers/%s/bag/verify"
const GenerateAllTokensPath = "/trainers/%s/tokens"
const GenerateTrainerStatsTokenPath = "/trainers/%s/stats/token"
const GenerateItemsTokenPath = "/trainers/%s/items/token"
const GeneratePokemonsTokenPath = "/trainers/%s/pokemons/token"
