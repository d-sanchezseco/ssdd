package main

import (
	"fmt"
)

func crearCliente() {
	var c Cliente

	c.ID = nextClienteID
	nextClienteID++

	ingresarDatos("nombre", &c.Nombre)
	ingresarDatos("telefono", &c.Telefono)

	for !checktelefono(c.Telefono) {
		fmt.Println("Teléfono inválido. Ingrese un teléfono válido (9 dígitos):")
		ingresarDatos("telefono", &c.Telefono)
	}

	ingresarDatos("email", &c.Email)

	mu.Lock()
	clientes = append(clientes, c)
	mu.Unlock()

	fmt.Println("Cliente creado con éxito:", c.Nombre)
}

func listarClientes() {
	mu.Lock()
	defer mu.Unlock()
	if listavacia(clientes, "clientes") {
		return
	}

	fmt.Println("Lista de Clientes:")
	for _, c := range clientes {
		fmt.Printf("ID: %d, Nombre: %s, Teléfono: %s, Email: %s\n", c.ID, c.Nombre, c.Telefono, c.Email)
		fmt.Println("  Vehículos:")
		for _, v := range c.Vehiculos {
			fmt.Printf("    Matrícula: %s, Marca: %s, Modelo: %s\n",
				v.Matricula, v.Marca, v.Modelo)
		}
	}
}

// Placeholder
func modificarCliente() {
	clearScreen()
	fmt.Println("Funcionalidad NO IMPLEMENTADA:")
	fmt.Println("Aquí se solicitaría el ID de un cliente para modificar sus datos.")
}

// Placeholder
func eliminarCliente() {
	clearScreen()
	fmt.Println("Funcionalidad NO IMPLEMENTADA:")
	fmt.Println("Aquí se solicitaría el ID de un cliente para eliminarlo.")
}

// Placeholder
func listarVehiculosDeCliente() {
	clearScreen()
	var clienteID int
	if listavacia(clientes, "clientes") {
		return
	}

	listarClientes()
	ingresarDatos("ID del cliente para ver sus vehiculos", &clienteID)

	mu.Lock()
	defer mu.Unlock()

	for _, c := range clientes {
		if c.ID == clienteID {
			fmt.Printf("Vehículos de %s (ID: %d):\n", c.Nombre, c.ID)
			if len(c.Vehiculos) == 0 {
				fmt.Println("Este cliente no tiene vehículos registrados.")
				return
			}
			for _, v := range c.Vehiculos {
				fmt.Printf("Matrícula: %s, Marca: %s, Modelo: %s\n", v.Matricula, v.Marca, v.Modelo)
			}
			return
		}
	}

	fmt.Println("Cliente no encontrado.")
}

func menuClientes() {
	var option int
	clearScreen()
	for {
		fmt.Println("----- Menú Clientes -----")
		fmt.Println("1. Crear Cliente")
		fmt.Println("2. Listar Clientes")
		fmt.Println("3. Modificar Cliente (P)")
		fmt.Println("4. Eliminar Cliente (P)")
		fmt.Println("5. Listar Vehículos de un Cliente")
		fmt.Println("6. Volver al Menú Principal")
		fmt.Print("Seleccione una opción: ")
		fmt.Scanln(&option)

		switch option {
		case 1:
			clearScreen()
			crearCliente()
		case 2:
			clearScreen()
			listarClientes()
		case 3:
			modificarCliente()
		case 4:
			eliminarCliente()
		case 5:
			listarVehiculosDeCliente()
		case 6:
			return
		default:
			fmt.Println("Opción no válida. Intente de nuevo.")

		}
		fmt.Println("\nPresione ENTER para volver al menú de Clientes...")
		reader.ReadString('\n') 
		reader.ReadString('\n')
	}
}
