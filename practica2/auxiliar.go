package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"unicode"
)

func clearScreen() {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "cls")
	default:
		cmd = exec.Command("clear")
	}

	cmd.Stdout = os.Stdout
	cmd.Run()
}
func listavacia[T any](lista []T, nombreLista string) bool {
	if len(lista) == 0 {
		fmt.Printf("No hay %s registrados.\n", nombreLista)
		return true
	}
	return false
}

func ingresarDatos[T any](campo string, valor *T) {
	var input string
	fmt.Print("Ingrese ", campo, ": ")
	input, _ = reader.ReadString('\n')
	input = strings.TrimSpace(input)

	switch v := any(valor).(type) {
	case *int:
		parsed, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println("Error: Ingrese un número válido")
			ingresarDatos(campo, valor)
			return
		}
		*v = parsed
	case *string:
		*v = input
	case *bool:
		*v = strings.ToLower(input) == "si" || strings.ToLower(input) == "sí"
	default:
		// Fallback para tipos no manejados explícitamente, asumiendo string de entrada
		*valor = any(input).(T)
	}
}

func checktelefono(telefono string) bool {
	if len(telefono) != 9 {
		return false
	}
	for _, c := range telefono {
		if !unicode.IsDigit(c) {
			return false
		}
	}
	return true
}

func checkmatricula(matricula string) bool {
	if len(matricula) != 7 {
		return false
	}
	for i, c := range matricula {
		if i < 4 {
			if !unicode.IsDigit(c) {
				return false
			}
		} else {
			if !unicode.IsLetter(c) {
				return false
			}
		}
	}
	return true
}
