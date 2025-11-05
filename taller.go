package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"unicode"
)

/*
	Taller de coches
*/

// Estructuras

type Cliente struct {
	ID        int
	Nombre    string
	Telefono  string
	Email     string
	Vehiculos []Vehiculo
}

type Vehiculo struct {
	Matricula    string
	Marca        string
	Modelo       string
	FechaEntrada string
	FechaSalida  string
	Incidencias  string
}

type Mecanico struct {
	ID           int
	Nombre       string
	Especialidad string
	Experiencia  int
	Activo       bool
}

type Incidencia struct {
	ID          int
	Mecanicos   []Mecanico
	Tipo        string
	Prioridad   string
	Descripcion string
	Estado      string
	Vehiculo    Vehiculo
}

// Variables globales
var clientes []Cliente
var vehiculos []Vehiculo
var mecanicos []Mecanico
var incidencias []Incidencia
var vehiculosEnTaller []string
var plazas map[int]string

var nextClienteID = 1
var nextMecanicoID = 1
var nextIncidenciaID = 1

var reader = bufio.NewReader(os.Stdin)

// Auxiliares

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
			if !unicode.IsLetter(c) || unicode.IsLower(c) {
				return false
			}
		}
	}
	return true
}

func isAsignado(idMecanico int) bool {
	for _, inc := range incidencias {
		for _, mec := range inc.Mecanicos {
			if mec.ID == idMecanico {
				return true
			}
		}
	}
	return false
}

// Cliente

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
	clientes = append(clientes, c)

	fmt.Println("Cliente creado con éxito:", c)
}

func listarClientes() {
	if listavacia(clientes, "clientes") {
		return
	}

	fmt.Println("Lista de Clientes:")
	for i := 0; i < len(clientes); i++ {
		var c Cliente = clientes[i]
		fmt.Printf("ID: %d, Nombre: %s, Teléfono: %s, Email: %s\n", c.ID, c.Nombre, c.Telefono, c.Email)
		fmt.Println("  Vehículos:")
		for j := 0; j < len(c.Vehiculos); j++ {
			var v Vehiculo = c.Vehiculos[j]
			fmt.Printf("    Matrícula: %s, Marca: %s, Modelo: %s, Fecha Entrada: %s, Fecha Salida: %s, Incidencias: %s\n",
				v.Matricula, v.Marca, v.Modelo, v.FechaEntrada, v.FechaSalida, v.Incidencias)
		}
	}
}

func modificarCliente() {
	var clienteID int
	var encontrado bool = false
	var opcion int

	if listavacia(clientes, "clientes") {
		return
	}

	fmt.Println("Ingrese el ID del cliente a modificar:")
	fmt.Scanln(&clienteID)

	for i := 0; i < len(clientes); i++ {
		if clientes[i].ID == clienteID {
			encontrado = true
			fmt.Println("Seleccione el campo a modificar:")
			fmt.Println("1. Nombre")
			fmt.Println("2. Teléfono")
			fmt.Println("3. Email")
			fmt.Scanln(&opcion)

			switch opcion {
			case 1:
				ingresarDatos("nuevo nombre", &clientes[i].Nombre)
			case 2:
				ingresarDatos("nuevo teléfono", &clientes[i].Telefono)
				for !checktelefono(clientes[i].Telefono) {
					fmt.Println("Teléfono inválido. Ingrese un teléfono válido (9 dígitos):")
					ingresarDatos("nuevo teléfono", &clientes[i].Telefono)
				}
			case 3:
				ingresarDatos("nuevo email", &clientes[i].Email)
			default:
				fmt.Println("Opción no válida.")
				return
			}
			fmt.Println("Cliente modificado con éxito:", clientes[i])
			break
		}
	}

	if !encontrado {
		fmt.Println("Cliente no encontrado.")
	}
}

func eliminarCliente() {
	var clienteID int
	var encontrado bool = false

	if listavacia(clientes, "clientes") {
		return
	}

	fmt.Println("Ingrese el ID del cliente a eliminar:")
	fmt.Scanln(&clienteID)
	for i := 0; i < len(clientes); i++ {
		if clientes[i].ID == clienteID {
			encontrado = true
			clientes = append(clientes[:i], clientes[i+1:]...)
			fmt.Println("Cliente eliminado con éxito.")
			break
		}
	}

	if !encontrado {
		fmt.Println("Cliente no encontrado.")
	}
}

// Vehículo

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
	}
}

func listarVehiculos() {
	if listavacia(clientes, "clientes") {
		return
	}

	fmt.Println("Lista de Vehículos:")
	for _, c := range clientes {
		for _, v := range c.Vehiculos {
			fmt.Printf("Cliente ID: %d, Nombre: %s -> Matrícula: %s, Marca: %s, Modelo: %s, Fecha Entrada: %s, Fecha Salida: %s, Incidencias: %s\n",
				c.ID, c.Nombre, v.Matricula, v.Marca, v.Modelo, v.FechaEntrada, v.FechaSalida, v.Incidencias)
		}
	}
}

func modificarVehiculo() {
	var matricula string
	var option int
	if listavacia(vehiculos, "vehículos") {
		return
	}

	fmt.Print("Introducir matricula del vehiculo a modificar: ")
	fmt.Scanln(&matricula)

	if !checkmatricula(matricula) {
		fmt.Println("Matrícula inválida. Ingrese una matrícula válida (formato: 1234ABC):")
		fmt.Scanln(&matricula)
	}

	for i := 0; i < len(vehiculos); i++ {
		if vehiculos[i].Matricula == matricula {
			fmt.Println("Vehículo encontrado.")
			fmt.Println("Datos actuales del vehículo:", vehiculos[i])
			fmt.Println("Seleccione el campo a modificar:")
			fmt.Println("1. Marca")
			fmt.Println("2. Modelo")
			fmt.Println("3. Fecha de Entrada")
			fmt.Println("4. Fecha de Salida")
			fmt.Println("5. Incidencias")
			fmt.Print("Opción: ")
			fmt.Scanln(&option)

			switch option {
			case 1:
				ingresarDatos("nueva marca", &vehiculos[i].Marca)
			case 2:
				ingresarDatos("nuevo modelo", &vehiculos[i].Modelo)
			case 3:
				ingresarDatos("nueva fecha de entrada", &vehiculos[i].FechaEntrada)
			case 4:
				ingresarDatos("nueva fecha de salida", &vehiculos[i].FechaSalida)
			case 5:
				ingresarDatos("nuevas incidencias", &vehiculos[i].Incidencias)
			default:
				fmt.Println("Opción no válida.")
				return
			}
			fmt.Println("Vehículo modificado con éxito:", vehiculos[i])

			return
		}
	}

	fmt.Println("Vehículo no encontrado.")

}

func eliminarVehiculo() {
	var matricula string
	var encontrado bool = false

	if listavacia(vehiculos, "vehículos") {
		return
	}
	fmt.Print("Introducir matricula del vehiculo a eliminar: ")
	fmt.Scanln(&matricula)

	if !checkmatricula(matricula) {
		fmt.Println("Matrícula inválida. Ingrese una matrícula válida (formato: 1234ABC):")
		fmt.Scanln(&matricula)
	}
	for i := 0; i < len(vehiculos); i++ {
		if vehiculos[i].Matricula == matricula {
			encontrado = true
			vehiculos = append(vehiculos[:i], vehiculos[i+1:]...)
			fmt.Println("Vehículo eliminado con éxito.")
			break
		}
	}

	if !encontrado {
		fmt.Println("Vehículo no encontrado.")
	}
}

// Mecanicos
func crearMecanico() {
	var m Mecanico
	var activo string

	m.ID = nextMecanicoID
	nextMecanicoID++

	ingresarDatos("nombre", &m.Nombre)
	ingresarDatos("especialidad", &m.Especialidad)
	ingresarDatos("experiencia (años)", &m.Experiencia)
	ingresarDatos("activo (si/no)", &activo)

	if strings.ToLower(activo) == "si" {
		m.Activo = true
	}

	mecanicos = append(mecanicos, m)
	fmt.Println("Mecánico creado con éxito:", m)

	actualizarPlazas()
}

func listarMecanicos() {
	if listavacia(mecanicos, "mecánicos") {
		return
	}

	fmt.Println("Lista de Mecánicos:")
	for i := 0; i < len(mecanicos); i++ {
		var m Mecanico = mecanicos[i]
		var activoStr string = "No"
		if m.Activo {
			activoStr = "Sí"
		}
		fmt.Printf("ID: %d, Nombre: %s, Especialidad: %s, Experiencia: %d años, Activo: %s\n",
			m.ID, m.Nombre, m.Especialidad, m.Experiencia, activoStr)
	}
}

func activarODesactivarMecanico() {
	var id int
	fmt.Print("Ingrese el ID del mecánico: ")
	fmt.Scanln(&id)

	for i := 0; i < len(mecanicos); i++ {
		if mecanicos[i].ID == id {
			mecanicos[i].Activo = !mecanicos[i].Activo
			estado := "inactivo"

			if mecanicos[i].Activo {
				estado = "activo"
			}

			fmt.Printf("El mecánico %s ahora está %s.\n", mecanicos[i].Nombre, estado)

			actualizarPlazas()
			return
		}
	}

	fmt.Println("Mecánico no encontrado.")
}

func modificarMecanico() {
	var mecanicoID int
	var encontrado bool = false
	var opcion int

	if listavacia(mecanicos, "mecánicos") {
		return
	}
	fmt.Println("Ingrese el ID del mecánico a modificar:")
	fmt.Scanln(&mecanicoID)

	for i := 0; i < len(mecanicos); i++ {
		if mecanicos[i].ID == mecanicoID {
			encontrado = true
			fmt.Println("Seleccione el campo a modificar:")
			fmt.Println("1. Nombre")
			fmt.Println("2. Especialidad")
			fmt.Println("3. Experiencia")
			fmt.Println("4. Activo")
			fmt.Scanln(&opcion)

			switch opcion {
			case 1:
				ingresarDatos("nuevo nombre", &mecanicos[i].Nombre)
			case 2:
				ingresarDatos("nueva especialidad", &mecanicos[i].Especialidad)
			case 3:
				ingresarDatos("nueva experiencia (años)", &mecanicos[i].Experiencia)
			case 4:
				activarODesactivarMecanico()
			default:
				fmt.Println("Opción no válida.")
				return
			}
			fmt.Println("Mecánico modificado con éxito:", mecanicos[i])
			break
		}
	}

	if !encontrado {
		fmt.Println("Mecánico no encontrado.")
	}

}

func eliminarMecanico() {
	var mecanicoID int
	var encontrado bool = false
	if listavacia(mecanicos, "mecánicos") {
		return
	}
	fmt.Println("Ingrese el ID del mecánico a eliminar:")
	fmt.Scanln(&mecanicoID)
	for i := 0; i < len(mecanicos); i++ {
		if mecanicos[i].ID == mecanicoID {
			encontrado = true
			mecanicos = append(mecanicos[:i], mecanicos[i+1:]...)
			fmt.Println("Mecánico eliminado con éxito.")
			break
		}
	}
	if !encontrado {
		fmt.Println("Mecánico no encontrado.")
	}
}

// Incidencias

func crearIncidencia() {
	var inc Incidencia
	var matricula string
	var encontrado bool = false

	inc.ID = nextIncidenciaID
	nextIncidenciaID++

	ingresarDatos("tipo", &inc.Tipo)
	ingresarDatos("prioridad", &inc.Prioridad)
	ingresarDatos("descripción", &inc.Descripcion)
	ingresarDatos("estado", &inc.Estado)
	ingresarDatos("matricula del vehiculo", &matricula) // Asumiendo que se guarda en Descripcion
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

	inc.Mecanicos = asignarMecanicos()
	incidencias = append(incidencias, inc)

	fmt.Println("\nIncidencia creada con éxito:")
	fmt.Printf("ID: %d | Tipo: %s | Prioridad: %s | Estado: %s\n",
		inc.ID, inc.Tipo, inc.Prioridad, inc.Estado)
	fmt.Printf("Incidencia del vehículo: %s\n", inc.Vehiculo.Matricula)
	fmt.Printf("Mecánicos asignados: %d\n", len(inc.Mecanicos))
}
func asignarMecanicos() []Mecanico {
	var asignados []Mecanico
	var idMecanico int
	var respuesta string
	var encontrado bool = false

	if listavacia(mecanicos, "mecánicos") {
		return asignados
	}

	for {
		fmt.Println("\nMecánicos disponibles para asignar:")
		listarMecanicos()

		fmt.Print("Ingrese el ID del mecánico a asignar: ")
		fmt.Scanln(&idMecanico)

		for i := 0; i < len(mecanicos); i++ {
			if mecanicos[i].ID == idMecanico {
				asignados = append(asignados, mecanicos[i])
				fmt.Println("Mecánico asignado:", mecanicos[i].Nombre)
				encontrado = true
				break
			}
		}

		if !encontrado {
			fmt.Println("Mecánico no encontrado.")
		}

		fmt.Print("¿Desea asignar otro mecánico? (si/no): ")
		respuesta, _ = reader.ReadString('\n')
		respuesta = strings.ToLower(strings.TrimSpace(respuesta))

		if respuesta != "si" {
			break
		}
	}
	return asignados
}
func listarIncidencias() {
	if listavacia(incidencias, "incidencias") {
		return
	}

	fmt.Println("Lista de Incidencias:")
	for i := 0; i < len(incidencias); i++ {
		var inc Incidencia = incidencias[i]
		fmt.Printf("------------ INCIDENCIA %d ------------\n", inc.ID)
		fmt.Printf("ID: %d, Tipo: %s, Prioridad: %s, Descripción: %s, Estado: %s\n",
			inc.ID, inc.Tipo, inc.Prioridad, inc.Descripcion, inc.Estado)
		fmt.Println("  Mecánicos asignados:")
		for j := 0; j < len(inc.Mecanicos); j++ {
			var m Mecanico = inc.Mecanicos[j]
			fmt.Printf("    ID: %d, Nombre: %s, Especialidad: %s\n",
				m.ID, m.Nombre, m.Especialidad)
		}
		fmt.Println("-------------------------------------")
	}
}
func modificarIncidencia() {
	var incidenciaID int
	var encontrado bool = false
	var opcion int

	if listavacia(incidencias, "incidencias") {
		return
	}

	fmt.Println("Ingrese el ID del cliente a modificar:")
	fmt.Scanln(&incidenciaID)

	for i := 0; i < len(incidencias); i++ {
		if incidencias[i].ID == incidenciaID {
			encontrado = true
			fmt.Println("Seleccione el campo a modificar:")
			fmt.Println("1. Tipo de incidencia")
			fmt.Println("2. Prioridad")
			fmt.Println("3. Descripción")
			fmt.Println("4. Estado")
			fmt.Scanln(&opcion)

			switch opcion {
			case 1:
				ingresarDatos("tipo", &incidencias[i].Tipo)
			case 2:
				ingresarDatos("prioridad", &incidencias[i].Prioridad)
			case 3:
				ingresarDatos("descripcion", &incidencias[i].Descripcion)
			case 4:
				ingresarDatos("estado", &incidencias[i].Estado)
			default:
				fmt.Println("Opción no válida.")
				return
			}
			fmt.Println("Incidencia modificada con éxito:", incidencias[i])
			break
		}
	}

	if !encontrado {
		fmt.Println("Cliente no encontrado.")
	}

}
func eliminarIncidencia() {
	var incidenciaID int
	var encontrado bool = false

	if listavacia(incidencias, "incidencias") {
		return
	}
	fmt.Println("Ingrese el ID de la incidencia a eliminar:")
	fmt.Scanln(&incidenciaID)
	for i := 0; i < len(incidencias); i++ {
		if incidencias[i].ID == incidenciaID {
			encontrado = true
			incidencias = append(incidencias[:i], incidencias[i+1:]...)
			fmt.Println("Incidencia eliminada con éxito.")
			break
		}
	}

	if !encontrado {
		fmt.Println("Incidencia no encontrada.")
	}
}
func CambiarEstadoIncidencia() {
	var id int
	var nuevoEstado string
	var encontrada bool = false
	var i int

	if len(incidencias) == 0 {
		fmt.Println("No hay incidencias registradas en el taller.")
		return
	}

	listarIncidencias()
	fmt.Print("Ingrese el ID de la incidencia que desea actualizar: ")
	fmt.Scanln(&id)

	for i = 0; i < len(incidencias); i++ {
		if incidencias[i].ID == id {
			encontrada = true
			break
		}
	}

	if !encontrada {
		fmt.Println("No se encontró ninguna incidencia con ese ID.")
		return
	}

	fmt.Printf("Incidencia encontrada: %s (Estado actual: %s)\n", incidencias[i].Tipo, incidencias[i].Estado)
	fmt.Println("Estados disponibles: abierta, en proceso, cerrada")

	fmt.Print("Ingrese el nuevo estado: ")
	fmt.Scanln(&nuevoEstado)

	nuevoEstado = strings.ToLower(strings.TrimSpace(nuevoEstado))

	if nuevoEstado != "abierta" && nuevoEstado != "en proceso" && nuevoEstado != "cerrada" {
		fmt.Println("Estado no válido. Los estados posibles son: abierta, en proceso o cerrada.")
		return
	}

	incidencias[i].Estado = nuevoEstado
	fmt.Printf("Estado de la incidencia con ID %d actualizado a: %s\n", incidencias[i].ID, incidencias[i].Estado)
}

// Taller

func actualizarPlazas() {
	totalPlazas := calcularplazasTotales()

	if totalPlazas == 0 {
		plazas = make(map[int]string)
		return
	}

	if len(plazas) == totalPlazas {
		return
	}

	plazas = make(map[int]string)
	for i := 1; i <= totalPlazas; i++ {
		plazas[i] = ""
	}
}
func calcularplazasTotales() int {
	var totalPlazas int = 0
	for _, m := range mecanicos {
		if m.Activo {
			totalPlazas += 2
		}
	}
	return totalPlazas

}
func asignarVehiculoAPlaza() bool {

	listarVehiculos()
	fmt.Print("Ingrese la matrícula del vehículo: ")
	matricula, _ := reader.ReadString('\n')
	matricula = strings.TrimSpace(matricula)

	if !checkmatricula(matricula) {
		fmt.Print("Matricula erronea")
		matricula, _ := reader.ReadString('\n')
		matricula = strings.TrimSpace(matricula)
	}

	encontrado := false
	for _, v := range vehiculos {
		if v.Matricula == matricula {
			encontrado = true
			break
		}
	}

	for _, v := range plazas {
		if v == matricula {
			fmt.Println("Ese vehiculo ya esta asignado a una plaza del taller")
			return false
		}
	}

	if !encontrado {
		fmt.Println("La matrícula ingresada no existe en el registro de vehículos.")
		return false
	}

	for i := 1; i <= len(plazas); i++ {
		if plazas[i] == "" {
			plazas[i] = matricula
			fmt.Printf("Vehiculo con matricula %s asignado a plaza %d", matricula, i)
			return true
		}
	}
	fmt.Println("No hay plazas disponibles en el taller")
	return false
}
func liberarPlaza() bool {
	var matricula string

	actualizarPlazas()

	if len(plazas) == 0 {
		fmt.Println("No hay mecánicos activos, por tanto no hay vehículos en el taller.")
		return false
	}

	ingresarDatos("matrícula", &matricula)
	if !checkmatricula(matricula) {
		fmt.Println("Matrícula inválida. Ingrese una matrícula válida (formato: 1234ABC):")
		ingresarDatos("matrícula", &matricula)

	}

	for i, v := range plazas {
		if v == matricula {
			plazas[i] = ""
			fmt.Printf("Vehículo con matrícula %s ha salido del taller (plaza %d liberada).\n", matricula, i)
			return true
		}
	}

	// Si no se encontró el vehículo
	fmt.Println("El vehículo no está actualmente en el taller.")
	return false
}

func estadoTaller() {
	fmt.Println("Estado del Taller:")
	total := calcularplazasTotales()
	ocupadas := 0

	for i := 1; i <= total; i++ {
		if plazas[i] != "" {
			fmt.Printf("Plaza %d: %s\n", i, plazas[i])
			ocupadas++
		} else {
			fmt.Printf("Plaza %d: (libre)\n", i)
		}
	}

	fmt.Printf("\nPlazas totales: %d\n", total)
	fmt.Printf("Plazas ocupadas: %d\n", ocupadas)
	fmt.Printf("Plazas libres: %d\n", total-ocupadas)
}

// Listados
func ListarIncidenciasVehiculo() {
	var matricula string
	ingresarDatos(matricula, &matricula)

	if !checkmatricula(matricula) {
		fmt.Println("Matrícula inválida. Ingrese una matrícula válida (formato: 1234ABC):")
		ingresarDatos(matricula, &matricula)

	}
	if listavacia(incidencias, "incidencias") {
		return
	}

	fmt.Printf("Incidencias del vehículo %s:\n", matricula)
	for i := 0; i < len(incidencias); i++ {
		var inc Incidencia = incidencias[i]
		if inc.Vehiculo.Matricula == matricula {
			fmt.Printf("ID: %d, Tipo: %s, Prioridad: %s, Descripción: %s, Estado: %s\n",
				inc.ID, inc.Tipo, inc.Prioridad, inc.Descripcion, inc.Estado)
		}
	}
}

func ListarVehiculosCliente() {
	var clienteID int

	if listavacia(clientes, "clientes") {
		return
	}

	listarClientes()
	ingresarDatos("ID del cliente", &clienteID)

	for i := 0; i < len(clientes); i++ {
		var c Cliente = clientes[i]
		if c.ID == clienteID {
			fmt.Printf("Vehículos del cliente %s (ID: %d):\n", c.Nombre, c.ID)
			for j := 0; j < len(c.Vehiculos); j++ {
				var v Vehiculo = c.Vehiculos[j]
				fmt.Printf("Matrícula: %s, Marca: %s, Modelo: %s, Fecha Entrada: %s, Fecha Salida: %s, Incidencias: %s\n",
					v.Matricula, v.Marca, v.Modelo, v.FechaEntrada, v.FechaSalida, v.Incidencias)
			}
			return
		}
	}

	fmt.Println("Cliente no encontrado.")
}

func ListarMecanicosDisponibles() {
	if listavacia(mecanicos, "mecánicos") {
		return
	}

	fmt.Println("Mecánicos disponibles:")

	disponibles := 0
	for _, m := range mecanicos {
		if m.Activo && !isAsignado(m.ID) {
			fmt.Printf("ID: %d, Nombre: %s, Especialidad: %s, Experiencia: %d años\n",
				m.ID, m.Nombre, m.Especialidad, m.Experiencia)
			disponibles++
		}
	}

	if disponibles == 0 {
		fmt.Println("No hay mecánicos disponibles (todos están asignados a incidencias activas).")
	}
}

func ListarIncidenciasMecanico() {
	var id int

	if listavacia(mecanicos, "mecánicos") || listavacia(incidencias, "incidencias") {
		return
	}

	listarMecanicos()

	fmt.Print("Ingrese el ID del mecánico: ")
	fmt.Scanln(&id)

	encontradas := 0
	fmt.Printf("Incidencias asignadas al mecánico con ID %d:\n", id)
	for _, inc := range incidencias {
		for _, m := range inc.Mecanicos {
			if m.ID == id {
				fmt.Printf("ID: %d | Tipo: %s | Prioridad: %s | Estado: %s | Vehículo: %s\n",
					inc.ID, inc.Tipo, inc.Prioridad, inc.Estado, inc.Vehiculo.Matricula)
				encontradas++
				break
			}
		}
	}

	if encontradas == 0 {
		fmt.Println("Este mecánico no tiene incidencias asignadas.")
	}
}

func ListarClientesConVehiculosEnTaller() {
	if len(plazas) == 0 {
		fmt.Println("No hay mecánicos activos, por tanto no hay plazas en el taller.")
		return
	}

	if len(clientes) == 0 {
		fmt.Println("No hay clientes registrados.")
		return
	}

	clientesMostrados := make(map[int]bool)
	hayVehiculos := false

	fmt.Println("\nClientes con vehículos actualmente en el taller:")

	for _, matricula := range plazas {
		if matricula == "" {
			continue
		}

		for _, cliente := range clientes {
			for _, v := range cliente.Vehiculos {
				if v.Matricula == matricula {
					if !clientesMostrados[cliente.ID] {
						fmt.Printf("Cliente: %s | Teléfono: %s | Email: %s\n",
							cliente.Nombre, cliente.Telefono, cliente.Email)
						clientesMostrados[cliente.ID] = true
					}
					fmt.Printf("Vehículo: %s (%s %s)\n", v.Matricula, v.Marca, v.Modelo)
					hayVehiculos = true
				}
			}
		}
	}

	if !hayVehiculos {
		fmt.Println("No hay vehículos actualmente en el taller.")
	}
}

func ListarIncidenciasEnTaller() {

	if len(incidencias) == 0 {
		fmt.Println("No hay incidencias registradas en el taller.")
		return
	}

	fmt.Println("Listado de incidencias en el taller:")
	fmt.Println("------------------------------------")

	for _, inc := range incidencias {
		fmt.Printf("ID: %d\n", inc.ID)
		fmt.Printf("Tipo: %s\n", inc.Tipo)
		fmt.Printf("Prioridad: %s\n", inc.Prioridad)
		fmt.Printf("Estado: %s\n", inc.Estado)
		fmt.Printf("Descripción: %s\n", inc.Descripcion)
		fmt.Printf("Vehículo: %s (%s %s)\n", inc.Vehiculo.Matricula, inc.Vehiculo.Marca, inc.Vehiculo.Modelo)

		if len(inc.Mecanicos) > 0 {
			fmt.Println("Mecánicos asignados:")
			for _, m := range inc.Mecanicos {
				fmt.Printf("   - %s (ID: %d, Especialidad: %s, Experiencia: %d años)\n",
					m.Nombre, m.ID, m.Especialidad, m.Experiencia)
			}
		} else {
			fmt.Println("Mecánicos asignados: Ninguno")
		}

		fmt.Println("------------------------------------")
	}
}

// Menús
func menuPrincipal() {
	var option int
	clearScreen()
	fmt.Println("----- Menú Principal -----")
	fmt.Println("1. Clientes")
	fmt.Println("2. Vehículos")
	fmt.Println("3. Mecánicos")
	fmt.Println("4. Incidencias")
	fmt.Println("5. Estado del Taller")
	fmt.Println("6. Salir")
	fmt.Print("Seleccione una opción: ")
	fmt.Scanln(&option)
	switch option {
	case 1:
		menuClientes()
	case 2:
		menuVehiculos()
	case 3:
		menuMecanicos()
	case 4:
		menuIncidencias()
	case 5:
		menuTaller()
	case 6:
		fmt.Println("Saliendo del programa...")
		return
	default:
		fmt.Println("Opción no válida. Intente de nuevo.")
		menuPrincipal()

	}

}

func menuClientes() {
	var option int
	clearScreen()
	for {
		fmt.Println("----- Menú Clientes -----")
		fmt.Println("1. Crear Cliente")
		fmt.Println("2. Listar Clientes")
		fmt.Println("3. Modificar Cliente")
		fmt.Println("4. Eliminar Cliente")
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
			clearScreen()
			modificarCliente()
		case 4:
			clearScreen()
			eliminarCliente()
		case 5:
			ListarVehiculosCliente()
		case 6:
			menuPrincipal()
		default:
			fmt.Println("Opción no válida. Intente de nuevo.")
			fmt.Scan(&option)
		}

		fmt.Println("\nPresione ENTER para volver al menu de clientes ...")
		reader.ReadString('\n')
	}

}

func menuVehiculos() {
	var option int
	clearScreen()
	for {
		fmt.Println("----- Menú Vehiculos -----")
		fmt.Println("1. Crear Vehiculo")
		fmt.Println("2. Listar Vehiculos")
		fmt.Println("3. Modificar Vehiculo")
		fmt.Println("4. Eliminar Vehiculo")
		fmt.Println("5. Listar Incidencias de un Vehículo")
		fmt.Println("6. Volver al Menú Principal")
		fmt.Print("Seleccione una opción: ")
		fmt.Scanln(&option)

		switch option {
		case 1:
			crearVehiculo()
		case 2:
			listarVehiculos()
		case 3:
			modificarVehiculo()
		case 4:
			eliminarVehiculo()
		case 5:
			ListarIncidenciasVehiculo()
		case 6:
			menuPrincipal()
		default:
			fmt.Println("Opción no válida. Intente de nuevo.")
			fmt.Scan(&option)
		}

		fmt.Println("\nPresione ENTER para volver al menu de vehiculos ...")
		reader.ReadString('\n')
	}
}

func menuIncidencias() {
	var option int
	clearScreen()
	for {
		fmt.Println("----- Menú Incidencias -----")
		fmt.Println("1. Crear Incidencia")
		fmt.Println("2. Listar Incidencia")
		fmt.Println("3. Modificar Incidencia")
		fmt.Println("4. Cambiar Estado de Incidencia")
		fmt.Println("5. Eliminar Incidencia")
		fmt.Println("6. Volver al Menú Principal")
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
			clearScreen()
			modificarIncidencia()
		case 4:
			clearScreen()
			CambiarEstadoIncidencia()
		case 5:
			clearScreen()
			eliminarIncidencia()
		case 6:
			menuPrincipal()
		default:
			fmt.Println("Opción no válida. Intente de nuevo.")
			fmt.Scan(&option)
		}

		fmt.Println("\nPresione ENTER para volver al menu de incidencias...")
		reader.ReadString('\n')
	}
}

func menuMecanicos() {
	var option int
	clearScreen()
	for {
		fmt.Println("----- Menú Mecanicos -----")
		fmt.Println("1. Crear Mecanico")
		fmt.Println("2. Listar Mecanicos")
		fmt.Println("3. Modificar Mecanico")
		fmt.Println("4. Eliminar Mecanico")
		fmt.Println("5. Listar Mecanicos Disponibles")
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
			clearScreen()
			modificarMecanico()
		case 4:
			clearScreen()
			eliminarMecanico()
		case 5:
			clearScreen()
			ListarMecanicosDisponibles()
		case 6:
			ListarIncidenciasMecanico()
		case 7:
			menuPrincipal()
		default:
			fmt.Println("Opción no válida. Intente de nuevo.")
			fmt.Scan(&option)
		}

		fmt.Println("\nPresione ENTER para volver al menu de mecanicos...")
		reader.ReadString('\n')
	}
}

func menuTaller() {
	var option int
	clearScreen()
	for {
		fmt.Println("----- Menú Taller -----")
		fmt.Println("1. Asignar plaza")
		fmt.Println("2. Liberar plaza")
		fmt.Println("3. Estado Taller")
		fmt.Println("4. Listar Clientes con Vehículos en Taller")
		fmt.Println("5. Listar incidencias en Taller")
		fmt.Println("6. Volver al Menú Principal")
		fmt.Print("Seleccione una opción: ")
		fmt.Scanln(&option)

		switch option {
		case 1:
			clearScreen()
			asignarVehiculoAPlaza()
		case 2:
			clearScreen()
			liberarPlaza()
		case 3:
			clearScreen()
			estadoTaller()
		case 4:
			clearScreen()
			ListarClientesConVehiculosEnTaller()
		case 5:
			clearScreen()
			ListarIncidenciasEnTaller()
		case 6:
			clearScreen()
			menuPrincipal()
		default:
			fmt.Println("Opción no válida. Intente de nuevo.")
			fmt.Scan(&option)
		}

		fmt.Println("\nPresione ENTER para volver al menu de mecanicos...")
		reader.ReadString('\n')
	}
}

// Programa principal

func main() {
	clearScreen()
	menuPrincipal()

}
