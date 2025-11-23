package main

import (
	"fmt"
	"strings"
)

func crearIncidencia() {
	var inc Incidencia
	var matricula string
	var encontrado bool = false

	inc.ID = nextIncidenciaID
	nextIncidenciaID++

	ingresarDatos("tipo (mecánica/eléctrica/carrocería)", &inc.Tipo)
	ingresarDatos("prioridad", &inc.Prioridad)
	ingresarDatos("descripción", &inc.Descripcion)

	ingresarDatos("matricula del vehiculo", &matricula)
	for !checkmatricula(matricula) {
		fmt.Println("Matrícula inválida. Ingrese una matrícula válida (formato: 1234ABC):")
		ingresarDatos("matrícula", &matricula)
	}

	for _, v := range vehiculos {
		if v.Matricula == matricula {
			inc.Vehiculo = v
			encontrado = true
			break
		}
	}

	if !encontrado {
		fmt.Println("Vehículo no encontrado. La incidencia no se ha creado.")
		return
	}

	inc.Estado = "en espera"

	NuevaIncidencia <- inc

	mu.Lock()
	incidencias = append(incidencias, inc)
	mu.Unlock()

	fmt.Println("\nIncidencia creada y ENVIADA a la cola de atención del taller (Simulación).")
	fmt.Printf("ID: %d | Tipo: %s | Estado: %s\n", inc.ID, inc.Tipo, inc.Estado)
}

func listarIncidencias() {
	mu.Lock()
	defer mu.Unlock()
	if listavacia(incidencias, "incidencias") {
		return
	}

	fmt.Println("Lista de Incidencias:")
	for _, inc := range incidencias {
		mecanicosAsignados := ""
		for _, m := range inc.Mecanicos {
			mecanicosAsignados += fmt.Sprintf("%s (ID: %d), ", m.Nombre, m.ID)
		}

		fmt.Printf("ID: %d | Tipo: %s | Prioridad: %s | Estado: %s | Vehículo: %s | Mecánicos: %s | Tiempo Acumulado: %v\n",
			inc.ID, inc.Tipo, inc.Prioridad, inc.Estado, inc.Vehiculo.Matricula, strings.TrimSuffix(mecanicosAsignados, ", "), inc.TiempoAcumulado)
	}
}

// Placeholder
func modificarIncidencia() {
	clearScreen()
	fmt.Println("Funcionalidad NO IMPLEMENTADA:")
	fmt.Println("Aquí se solicitaría el ID de una incidencia para modificar sus datos.")
}

func cambiarEstadoIncidencia() {
	clearScreen()
	var id int
	ingresarDatos("ID de la incidencia a cambiar estado", &id)

	mu.Lock()
	defer mu.Unlock()

	for i, inc := range incidencias {
		if inc.ID == id {
			fmt.Printf("Estado actual de la incidencia %d: %s\n", inc.ID, inc.Estado)
			fmt.Println("Estados disponibles: abierta, en espera, en proceso, prioritaria, finalizada")

			var nuevoEstado string
			ingresarDatos("nuevo estado", &nuevoEstado)

			nuevoEstado = strings.ToLower(nuevoEstado)

			if nuevoEstado != "abierta" && nuevoEstado != "en espera" && nuevoEstado != "en proceso" && nuevoEstado != "prioritaria" && nuevoEstado != "finalizada" {
				fmt.Println("Estado no válido.")
				return
			}
			incidencias[i].Estado = nuevoEstado
			fmt.Printf("Estado de la incidencia con ID %d actualizado a: %s\n", id, incidencias[i].Estado)
			return
		}
	}
	fmt.Println("Incidencia no encontrada.")
}

func eliminarIncidencia() {
	clearScreen()
	fmt.Println("Funcionalidad NO IMPLEMENTADA:")
	fmt.Println("Aquí se solicitaría el ID de una incidencia para eliminarla.")
}

func menuIncidencias() {
	var option int
	clearScreen()
	for {
		fmt.Println("----- Menú Incidencias -----\n1. Crear Incidencia\n2. Listar Incidencia\n3. Modificar Incidencia (P)\n4. Cambiar Estado de Incidencia\n5. Eliminar Incidencia (P)\n6. Volver al Menú Principal")
		fmt.Print("Seleccione una opción: ")
		fmt.Scanln(&option)

		switch option {
		case 1:
			clearScreen()
			crearIncidencia()
		case 2:
			clearScreen()
			listarIncidencias()
		case 3:
			modificarIncidencia()
		case 4:
			clearScreen()
			cambiarEstadoIncidencia()
		case 5:
			eliminarIncidencia()
		case 6:
			return
		default:
			fmt.Println("Opción no válida. Intente de nuevo.")
		}
		fmt.Println("\nPresione ENTER para volver al menú de Incidencias...")
		reader.ReadString('\n')
		reader.ReadString('\n')
	}
}
