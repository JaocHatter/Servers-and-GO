package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"log"

	"github.com/gorilla/mux"
)

//Acá creo una structura llamada name, la inicializo con sus respectivos atributo y especifico el formato en que las deseo
// en este caso un json--->
/*La sintaxis "json:<fieldname>" le dice al paquete json de Go que cuando se codifique o se decode
una estructura de este tipo, use el campo <fieldname> como el nombre del campo en el formato JSON.*/
type task struct {
	ID      int    `json:"ID"`
	Name    string `json:"Name"`
	Content string `json:"Content"`
}

// creo un array de task's para guardar todas las tareas!
type allTask []task

// CREO UNA TASK inicializada
var tasks = allTask{
	{
		ID:      1,
		Name:    "Task One",
		Content: "Hola mundo",
	},
	{
		ID:      2,
		Name:    "Bombardear Bolivia",
		Content: "Exterminio",
	},
}

// por cierto, colocamos un puntero al lado de la clase de nuestro parametro para pasarlo por referencia
// es decir , le pasaremos la direccion de memoria de nuestra variable para que cada modificación se realize!
func MostrarTask(w http.ResponseWriter, r *http.Request) {
	// Los mostramos en tipo JSON
	w.Header().Set("Content-Type", "aplication/json")

	json.NewEncoder(w).Encode(tasks)
}

// Creo una nueva funcion que permita al Usuario agregar una tarea al array tasks
func createTask(w http.ResponseWriter, r *http.Request) {
	//Creo una variable tipo task, en esta alojaremos lo que el Usuario ingresará en formato Json, para posteriormente
	//agregarlo a slicen tasks
	var NewTask task
	reqBody, error := ioutil.ReadAll(r.Body)
	// nil == NULL
	if error != nil {
		fmt.Fprintf(w, "Introduzca una tarea válida!")
	}
	/*Unmarshal es una función del paquete encoding/json en Go, que se utiliza para convertir datos en formato JSON
	a una estructura o variable en Go. La función Unmarshal toma como entrada una secuencia de bytes que representa datos en formato JSON y
	un puntero a una variable o estructura en Go, y asigna los valores del JSON a los campos correspondientes de
	la variable o estructura en Go.*/
	//acá estoy pasando informacion tipo json a la clase task
	json.Unmarshal(reqBody, &NewTask)
	//agrego la tarea Nueva!
	NewTask.ID = len(tasks) + 1
	tasks = append(tasks, NewTask)
	//En esta parte se brinda informacion al cliente y se le informa que su info ha sido enviada
	//-------------------------------------------------------------------------------------------
	//Brinda información del tipo de dato
	w.Header().Set("Content-Type", "aplication/json")
	//
	w.WriteHeader(http.StatusCreated)
	//Encode(NewTask)--> Es un json
	json.NewEncoder(w).Encode(NewTask)
}
func IndexRoute(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "================ BIENVENIDO A MI API! ================")
}
func main() {
	//CREO UN ROUTER de prueba
	router := mux.NewRouter().StrictSlash(true)
	//cada vez que  un usuario ingrese a la pagina se ejecutará la función IndexRoute

	router.HandleFunc("/", IndexRoute)
	router.HandleFunc("/tasks", MostrarTask)
	router.HandleFunc("/Add", createTask)
	log.Fatal(http.ListenAndServe(":3000", router))
}
