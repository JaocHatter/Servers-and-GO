package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

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
		Name:    "Tarea 1",
		Content: "Hola mundo!",
	},
	{
		ID:      2,
		Name:    "Bombardear Bolivia",
		Content: "Exterminio",
	},
}

// por cierto, colocamos un puntero al lado de la clase de nuestro parametro para pasarlo por referencia
// es decir , le pasaremos la direccion de memoria de nuestra variable para que cada modificación se realize!
func MostrarTasks(w http.ResponseWriter, r *http.Request) {

	//y con esta función descargamos un archivo con un arreglo de Jsons;
	w.Header().Set("Content-Type", "aplication/json")
	// Los mostramos en tipo JSON
	json.NewEncoder(w).Encode(tasks)
}

// Creo una nueva funcion que permita al Usuario agregar una tarea al array tasks
func createTask(w http.ResponseWriter, r *http.Request) {
	//Creo una variable tipo task, en esta alojaremos lo que el Usuario ingresará en formato Json, para posteriormente
	//agregarlo a slicen tasks
	var NewTask task
	//Uso la libreria ioutil para leer el fichero r.Body
	//Si la lectura es exitosa, me retorna el arreglo de bits y los deposita en reqBody, tambien retorna un error nulo.
	//caso contrario solo retorna un error
	reqBody, error := ioutil.ReadAll(r.Body)
	//reqBody es un .JSON
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

	w.WriteHeader(http.StatusCreated)
	//Encode(NewTask)--> Es un json
	json.NewEncoder(w).Encode(NewTask)
}

func IndexRoute(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "================ BIENVENIDO A MI API! ================\n")
	fmt.Fprint(w, "Lista de Funciones\n")
	fmt.Fprintf(w, "Ver todos los json :/tasks\n")
	fmt.Fprintf(w, "Agregar un json: /Add\n")
}

//Proximamente
/*func EliminarTask(w http.ResponseWriter,r *http.Request){
	fmt.Fprintf(w,"Que tarea desea eliminar: ")

}*/
func GetTask(w http.ResponseWriter, r *http.Request) {
	//la función mux.Vars(r) me permitira crear un  mapa, conjunto de llaves y valores
	//las llaves serán los ID
	vars := mux.Vars(r)
	taskID, err := strconv.Atoi(vars["id"])
	if err != nil {
		fmt.Fprintf(w, "ID NO VALIDA")
		return
	}
	for i := 0; i < len(tasks); i++ {
		if taskID == tasks[i].ID {
			w.Header().Set("Content-Type", "aplication/json")
			// Los mostramos en tipo JSON
			json.NewEncoder(w).Encode(tasks[i])
			return
		}
	}
	fmt.Fprintf(w, "NO SE ENCONTRÓ LA TAREA")
}

// Aquí ya ni le agruegue comentarios, porque además de ya haberle manyado
// las funciones se repiten y el procedimiento es similar!
func EliminarTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskID, err := strconv.Atoi(vars["id"])
	if err != nil {
		fmt.Fprintf(w, "ID NO VALIDA")
		return
	}
	for i, task := range tasks {
		if task.ID == taskID {
			tasks = append(tasks[:i], tasks[i+1:]...)
			fmt.Fprintf(w, "La id %v ha sido removida de la base de datos!", taskID)
			return
		}
	}
	fmt.Fprintf(w, "NO SE ENCONTRÓ LA TAREA")
}
func ActualizarTask(w http.ResponseWriter, r *http.Request) {
	//transformo los json en un diccionario o "map"
	vars := mux.Vars(r)
	taskID, err := strconv.Atoi(vars["id"])
	var UpTask task
	if err != nil {
		fmt.Fprintf(w, "ID NO VALIDA")
		return
	}
	//itero en la base de datos hasta encontrar un task con el mismo ID
	for i, t := range tasks {
		if t.ID == taskID {
			//LEO EL CONTENIDO MANDADO POR EL CLIENTE
			reqBody, err := ioutil.ReadAll(r.Body)
			if err != nil {
				fmt.Fprintf(w, "Task no valida, Recargue la página")
				return
			}
			//Necesito transformar mi variable UpTask en un json
			json.Unmarshal(reqBody, &UpTask)
			UpTask.ID = i + 1
			//Remlazo ese task que buscaba por el modificado
			tasks[i] = UpTask
			fmt.Fprintf(w, "La tarea ha sido actualizada exitosamente")
			return
		}
	}
	fmt.Fprintf(w, "No se encontró esa ID, Recargue la página")
}
func main() {
	//Creo un enroutador, para qué? me sirve
	router := mux.NewRouter().StrictSlash(true)
	//cada vez que  un usuario ingrese a la pagina se ejecutará la función IndexRoute
	router.HandleFunc("/", IndexRoute)
	router.HandleFunc("/tasks", MostrarTasks).Methods("GET")
	router.HandleFunc("/Add", createTask).Methods("POST")
	router.HandleFunc("/tasks/{id}", GetTask).Methods("GET")
	router.HandleFunc("/delete/{id}", EliminarTask).Methods("DELETE")
	router.HandleFunc("/Update/{id}", ActualizarTask).Methods("PUT")
	log.Fatal(http.ListenAndServe(":3000", router))
	//thats ALL :)
}

