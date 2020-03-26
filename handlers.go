package main

import (
	"encoding/json"
	"fmt"
	"github.com/NOVAPokemon/utils"
	trainerdb "github.com/NOVAPokemon/utils/database/trainer"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

func AddTrainer(w http.ResponseWriter, r *http.Request) {

	var trainer = &utils.Trainer{}
	err := json.NewDecoder(r.Body).Decode(trainer)

	if err != nil {
		http.Error(w, "An error occurred decoding trainer", http.StatusBadRequest)
		return
	}

	_, err = trainerdb.AddTrainer(*trainer)

	if err != nil {
		http.Error(w, "An error occurred inserting trainer", http.StatusInternalServerError)
		return
	}

}

func GetAllTrainers(w http.ResponseWriter, r *http.Request) {

	trainers := trainerdb.GetAllTrainers()

	fmt.Println(trainers)

	toSend, err := json.Marshal(trainers)

	if err != nil {
		http.Error(w, "An error occurred marshaling trainers", http.StatusInternalServerError)
		return
	}

	_, _ = w.Write(toSend)

}

func GetTrainerByUsername(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	trainerUsername := vars["username"]

	trainer, err := trainerdb.GetTrainerByUsername(trainerUsername)

	if err != nil {
		http.Error(w, "Trainer Missing", http.StatusNotFound)
		return
	}

	toSend, err := json.Marshal(trainer)

	if err != nil {
		http.Error(w, "An error occurred marshaling trainer", http.StatusInternalServerError)
		return
	}

	_, _ = w.Write(toSend)
}

func AddPokemonToTrainer(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	trainerUsername := vars["username"]

	var pokemon = &utils.Pokemon{}

	err := json.NewDecoder(r.Body).Decode(pokemon)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = trainerdb.AddPokemonToTrainer(trainerUsername, *pokemon)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(200)
}

func RemovePokemonFromTrainer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	trainerUsername := vars["username"]
	pokemonId, err := primitive.ObjectIDFromHex(vars["pokemonId"])

	if err != nil {
		http.Error(w, "Pokemon id is not valid", http.StatusBadRequest)
		return
	}

	err = trainerdb.RemovePokemonFromTrainer(trainerUsername, pokemonId)

	if err != nil {
		http.Error(w, "An error occurred removing pokemon", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(200)
}

func AddItemToTrainer(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	trainerUsername := vars["username"]

	var item = &utils.Item{}

	err := json.NewDecoder(r.Body).Decode(item)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = trainerdb.AddItemToTrainer(trainerUsername, *item)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(200)

}

func RemoveItemToTrainer(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	trainerUsername := vars["username"]
	itemId, err := primitive.ObjectIDFromHex(vars["itemId"])

	if err != nil {
		http.Error(w, "Pokemon id is not valid", http.StatusBadRequest)
		return
	}

	err = trainerdb.RemoveItemFromTrainer(trainerUsername, itemId)

	if err != nil {
		http.Error(w, "An error occurred removing item", http.StatusNotFound)
		return
	}

	w.WriteHeader(200)
}
