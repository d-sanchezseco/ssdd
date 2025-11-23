package main

import (
	"fmt"
	"strings"
)

func iniciarMecanico(m Mecanico) {
	mu.Lock()
	if m.ID == 0 {
		m.ID = nextMecanicoID
		nextMecanicoID++
		m.Activo = true
		mecanicos = append(mecanicos, m)
		fmt.Printf("Mecánico %s (ID: %d, Especialidad: %s) CONTRATADO DE INMEDIATO.\n", m.Nombre, m.ID, m.Especialidad)
	}
	mu.Unlock()

	go goroutineMecanico(m)
	MecanicosDisponibles <- m
}

func crearMecanico() {
	var m Mecanico
	var activo string

	m.ID = nextMecanicoID
	nextMecanicoID++

	ingresarDatos("nombre", &m.Nombre)
	ingresarDatos("especialidad", &m.Especialidad)
	ingresarDatos("experiencia (años)", &m.Experiencia)
	ingresarDatos("activo (si/no)", &activo)

	m.Activo = strings.ToLower(activo) == "si" || strings.ToLower(activo) == "sí"

	mu.Lock()
	mecanicos = append(mecanicos, m)
	mu.Unlock()
	fmt.Println("Mecánico creado con éxito:", m.Nombre)

	if m.Activo {
		go goroutineMecanico(m)
		MecanicosDisponibles <- m
	}
}
func listarMecanicos() {
	mu.Lock()
	defer mu.Unlock()
	if listavacia(mecanicos, "mecánicos") {
		return
	}

	fmt.Println("Lista de Mecánicos:")
	for _, m := range mecanicos {
		var activoStr string = "No"
		if m.Activo {
			activoStr = "Sí"
		}
		fmt.Printf("ID: %d, Nombre: %s, Especialidad: %s, Experiencia: %d años, Activo: %s\n",
			m.ID, m.Nombre, m.Especialidad, m.Experiencia, activoStr)
	}
}

// Placeholder
func modificarMecanico() {
	clearScreen()
	fmt.Println("Funcionalidad NO IMPLEMENTADA:")
	fmt.Println("Aquí se solicitaría el ID de un mecánico para modificar sus datos.")
}

// Placeholder
func eliminarMecanico() {
	clearScreen()
	fmt.Println("Funcionalidad NO IMPLEMENTADA:")
	fmt.Println("Aquí se solicitaría el ID de un mecánico para eliminarlo.")
}


func listarMecanicosActivos() {
	clearScreen()
	mu.Lock()
	defer mu.Unlock()

	fmt.Println("Mecánicos actualmente activos:")
	encontrados := 0
	for _, m := range mecanicos {
		if m.Activo {
			fmt.Printf("ID: %d, Nombre: %s, Especialidad: %s\n", m.ID, m.Nombre, m.Especialidad)
			encontrados++
		}
	}
	if encontrados == 0 {
		fmt.Println("No hay mecánicos activos en este momento.")
	}
}


func consultarIncidenciasMecanico() {
	clearScreen()
	var id int
	ingresarDatos("ID del Mecánico", &id)

	mu.Lock()
	defer mu.Unlock()

	encontradas := 0
	fmt.Printf("Incidencias asignadas al mecánico con ID %d:\n", id)
	for _, inc := range incidencias {
		for _, m := range inc.Mecanicos {
			if m.ID == id {
				fmt.Printf("ID: %d | Tipo: %s | Estado: %s | Vehículo: %s\n",
					inc.ID, inc.Tipo, inc.Estado, inc.Vehiculo.Matricula)
				encontradas++
				break
			}
		}
	}

	if encontradas == 0 {
		fmt.Println("Este mecánico no tiene incidencias asignadas o el ID no existe.")
	}
}

func menuMecanicos() {
	var option int
	clearScreen()
	for {
		fmt.Println("----- Menú Mecanicos -----")
		fmt.Println("1. Crear Mecanico")
		fmt.Println("2. Listar Mecanicos")
		fmt.Println("3. Modificar Mecanico (P)")
		fmt.Println("4. Eliminar Mecanico (P)")
		fmt.Println("5. Listar Mecanicos Activos")
		fmt.Println("6. Consultar incidencias de un Mecanico")
		fmt.Println("7. Volver al Menú Principal")
		fmt.Print("Seleccione una opción: ")
		fmt.Scanln(&option)

		switch option {
		case 1:
			clearScreen()
			crearMecanico()
		case 2:
			clearScreen()
			listarMecanicos()
		case 3:
			modificarMecanico()
		case 4:
			eliminarMecanico()
		case 5:
			listarMecanicosActivos()
		case 6:
			consultarIncidenciasMecanico()
		case 7:
			return
		default:
			fmt.Println("Opción no válida. Intente de nuevo.")
		}
		fmt.Println("\nPresione ENTER para volver al menú de Mecánicos...")
		reader.ReadString('\n')
		reader.ReadString('\n')
	}
}
