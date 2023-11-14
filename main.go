package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"sync"
)

var unUserData UserData
var unUserDataLock sync.Mutex


var temp *template.Template

type Student struct {
	Nom    string
	Prenom string
	Age    int
	Genre  string
}

type Promo struct {
	NomPromo    string
	Filiere     string
	Niveau      string
	NbEtudiant  []Student
}

type UserData struct {
	Nom           string
	Prenom        string
	Datedenaissance string
	Sexe          string
}

func main() {
	var err error // Déclarez err pour qu'elle soit accessible dans toute la portée de la fonction main

	temp, err = template.ParseGlob("./*.html") // Utilisez l'affectation simple pour temp et err
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	Data := Promo{
		NomPromo: "Mentor'ac",
		Filiere:  "Informatique",
		Niveau:   "5",
		NbEtudiant: []Student{
			{Nom: "Rodrigue", Prenom: "Ciryl", Age: 22, Genre: "Male"},
			{Nom: "MEDRREG", Prenom: "Kheir-Eddine", Age: 22, Genre: "Male"},
			{Nom: "PHILIPERT", Prenom: "Alan", Age: 26, Genre: "Male"},
		},
	}

	http.HandleFunc("/promo", func(w http.ResponseWriter, r *http.Request) {
		temp.ExecuteTemplate(w, "index", Data)
	})

	http.HandleFunc("/user/init", func(w http.ResponseWriter, r *http.Request) {
		temp.ExecuteTemplate(w, "User", nil)
	})

	rootDoc, _ := os.Getwd()
	fileserver := http.FileServer(http.Dir(rootDoc + "/asset"))
	http.Handle("/static/", http.StripPrefix("/static/", fileserver))

	http.HandleFunc("/user/treatment", func(w http.ResponseWriter, r *http.Request) {
        nom := r.FormValue("Nom")
        prenom := r.FormValue("Prenom")
        datedenaissance := r.FormValue("Datedenaissance")
        sexe := r.FormValue("Sexe")

		unUserDataLock.Lock()
		defer unUserDataLock.Unlock()

        // Créez une instance de UserData avec les données du formulaire
        unUserData = UserData{
            Nom:           nom,
            Prenom:        prenom,
            Datedenaissance: datedenaissance,
            Sexe:          sexe,
        }

        // Redirection vers la page d'affichage des données UserData
        http.Redirect(w, r, "/user/display", http.StatusSeeOther)
    })


	http.HandleFunc("/user/display", func(w http.ResponseWriter, r *http.Request) {
    temp.ExecuteTemplate(w, "UserData", unUserData)
	unUserDataLock.Lock()
	defer unUserDataLock.Unlock()
	temp.Execute(w, unUserData)


})

	http.ListenAndServe(":8080", nil)
}
