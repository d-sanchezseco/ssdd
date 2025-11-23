package main

import "fmt"

func crearVehiculo() {
	var v Vehiculo
	var clienteID int
	var encontrado bool = false

	ingresarDatos("matrícula", &v.Matricula)
	for !checkmatricula(v.Matricula) {
		fmt.Println("Matrícula inválida. Ingrese una matrícula válida (formato: 1234ABC):")
		ingresarDatos("matrícula", &v.Matricula)
	}
	ingresarDatos("marca", &v.Marca)
	ingresarDatos("modelo", &v.Modelo)
	ingresarDatos("fecha de entrada", &v.FechaEntrada)
	ingresarDatos("fecha de salida", &v.FechaSalida)
	ingresarDatos("incidencias", &v.Incidencias)

	fmt.Println("Ingrese el ID del cliente al que pertenece el vehículo:")
	listarClientes()
	fmt.Scanln(&clienteID)

	mu.Lock()
	defer mu.Unlock()
	for i := 0; i < len(clientes); i++ {
		if clientes[i].ID == clienteID {
			clientes[i].Vehiculos = append(clientes[i].Vehiculos, v)
			encontrado = true
			vehiculos = append(vehiculos, v)
			break
		}
	}

	if !encontrado {
		fmt.Println("Cliente no encontrado. Vehículo no asignado.")
	} else {
		fmt.Println("Vehículo creado y asignado con éxito.")
	}
}
func listarVehiculos() {
	mu.Lock()
	defer mu.Unlock()
	if listavacia(vehiculos, "vehículos") {
		return
	}
	fmt.Println("Lista de Vehículos (en base de datos):")
	for _, v := range vehiculos {
		fmt.Printf("Matrícula: %s, Marca: %s, Modelo: %s, Incidencias: %s\n",
			v.Matricula, v.Marca, v.Modelo, v.Incidencias)
	}
}

// Placeholder
func modificarVehiculo() {
	clearScreen()
	fmt.Println("Funcionalidad NO IMPLEMENTADA:")
	fmt.Println("Aquí se solicitaría la matrícula de un vehículo para modificar sus datos.")
}

// Placeholder
func eliminarVehiculo() {
	clearScreen()
	fmt.Println("Funcionalidad NO IMPLEMENTADA:")
	fmt.Println("Aquí se solicitaría la matrícula de un vehículo para eliminarlo.")
}

// Placeholder
func ListarIncidenciasVehiculo() {
	clearScreen()
	var matricula string
	ingresarDatos("matrícula del vehículo", &matricula)

	if !checkmatricula(matricula) {
		fmt.Println("Matrícula inválida.")
		return
	}

	mu.Lock()
	defer mu.Unlock()

	encontradas := 0
	fmt.Printf("Incidencias para el vehículo %s:\n", matricula)
	for _, inc := range incidencias {
		if inc.Vehiculo.Matricula == matricula {
			fmt.Printf("ID: %d, Tipo: %s, Prioridad: %s, Estado: %s\n",
				inc.ID, inc.Tipo, inc.Prioridad, inc.Estado)
			encontradas++
		}
	}
	if encontradas == 0 {
		fmt.Println("No se encontraron incidencias activas o históricas para este vehículo.")
	}
}

func menuVehiculos() {
	var option int
	clearScreen()
	for {
		fmt.Println("----- Menú Vehiculos -----")
		fmt.Println("1. Crear Vehiculo")
		fmt.Println("2. Listar Vehiculos (Base de Datos)")
		fmt.Println("3. Modificar Vehiculo (P)")
		fmt.Println("4. Eliminar Vehiculo (P)")
		fmt.Println("5. Listar Incidencias de un Vehículo")
		fmt.Println("6. Volver al Menú Principal")
		fmt.Print("Seleccione una opción: ")
		fmt.Scanln(&option)

		switch option {
		case 1:
			clearScreen()
			crearVehiculo()
		case 2:
			clearScreen()
			listarVehiculos()
		case 3:
			modificarVehiculo()
		case 4:
			eliminarVehiculo()
		case 5:
			clearScreen()
			ListarIncidenciasVehiculo()
		case 6:
			return
		default:
			fmt.Println("Opción no válida. Intente de nuevo.")
		}
		fmt.Println("\nPresione ENTER para volver al menú de Vehículos...")
		reader.ReadString('\n')
		reader.ReadString('\n')
	}
}
