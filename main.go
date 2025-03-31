package main

// Operador := cria uma variavel, atribui um valor e define um tipo automaticamente

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

// Estrutura de um todo
type Task struct {
	ID   int    `json:"id"`
	Description string `json:"description"`
	Completed bool   `json:"completed"`
}

var tasks []Task 	// Lista de tarefas
var mu sync.Mutex 	// Mutex para garantir concorrência
var nextID = 1 		// ID para próxima tarefa

// Função para listar tarefas
func getTask(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

// Função para criar nova tarefa
func addTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	decoder := json.NewDecoder(r.Body) // transforma o body da requisição na estrutura definido para tarefas (linha 11)

	if err := decoder.Decode(&task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	task.ID = nextID 			// Atribui o ID da tarefa
	nextID++ 					// Incrementa o ID para a próxima tarefa
	tasks = append(tasks, task) // Adiciona a tarefa à lista

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task) // Retorna a tarefa criada
}

// Função para mudar status da tarefa
func maskTaskComplete(w http.ResponseWriter, r *http.Request) {
	var taskID map[string]int

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&taskID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	for i := range tasks {
		if tasks[i].ID == taskID["id"] {
			tasks[i].Completed = true // Marca a tarefa como completa
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(tasks[i]) // Retorna a tarefa atualizada
			return
		}
	}

	http.NotFound(w, r) // Retorna erro se não encontrar a tarefa
}

// Função para deletar tarefa
func deleteTask(w http.ResponseWriter, r *http.Request) {
	var taskID map[string]int

	decoder := json.NewDecoder(r.Body)
	
	if err := decoder.Decode(&taskID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	for i, task := range tasks {
		if task.ID == taskID["id"] {
			tasks = append(tasks[:i], tasks[i+1:]...) // Remove a tarefa da lista
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(task) // Retorna a tarefa deletada
			return
		}
	}

	http.NotFound(w, r) // Retorna erro se não encontrar a tarefa
}

// Função principal
func main() {
	http.HandleFunc("/tasks", getTask) // Rota para listar tarefas
	http.HandleFunc("/add", addTask)     // Rota para adicionar tarefa
	http.HandleFunc("/complete", maskTaskComplete) // Rota para marcar tarefa como completa
	http.HandleFunc("/delete", deleteTask) // Rota para deletar tarefa

	fmt.Println("Servidor rodando na porta 8080")
	http.ListenAndServe(":8080", nil) // Inicia o servidor na porta 8080
}

