package main

import (
	"bufio"
	"fmt"
	"os"
	"sync"
	"time"
)

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

	TiempoAcumulado time.Duration
	esPrioritaria   bool
}

var clientes []Cliente
var vehiculos []Vehiculo
var mecanicos []Mecanico
var incidencias []Incidencia

var nextIncidenciaID = 1
var nextMecanicoID = 1
var nextClienteID = 1
var reader = bufio.NewReader(os.Stdin)

var (
	NuevaIncidencia        chan Incidencia
	MecanicosDisponibles   chan Mecanico
	ColaDeEspera           chan Incidencia
	ContratarMecanicoCanal chan string

	mu sync.Mutex
)

const (
	TiempoMecanica   = 5
	TiempoElectrica  = 7
	TiempoCarroceria = 11
	LimitePrioridad  = 15 * time.Second
)

func inicializarDatosPrueba() {
	clientes = []Cliente{{ID: 1, Nombre: "Juan Pérez", Telefono: "123456789", Email: "juan@mail.com"}}
	nextClienteID = 2

	vehiculos = []Vehiculo{
		{Matricula: "1234ABC", Marca: "Seat", Modelo: "Ibiza", Incidencias: "Fallo motor"},
		{Matricula: "5678DEF", Marca: "Ford", Modelo: "Focus", Incidencias: "Luces"},
	}
	clientes[0].Vehiculos = append(clientes[0].Vehiculos, vehiculos...)

	mecanicos = []Mecanico{
		{ID: 1, Nombre: "Luis", Especialidad: "mecánica", Experiencia: 5, Activo: true},
		{ID: 2, Nombre: "Elena", Especialidad: "eléctrica", Experiencia: 7, Activo: true},
		{ID: 3, Nombre: "Pablo", Especialidad: "carrocería", Experiencia: 10, Activo: true},
	}
	nextMecanicoID = 4

	fmt.Println("Datos de prueba cargados: 1 Cliente, 2 Vehículos, 3 Mecánicos.")
}

func menuPrincipal() {
	var option int
	clearScreen()
	for {
		fmt.Println("----- Menú Principal -----")
		fmt.Println("1. Clientes")
		fmt.Println("2. Vehículos")
		fmt.Println("3. Mecánicos")
		fmt.Println("4. Incidencias")
		fmt.Println("5. Estado del Taller (Simulación)")
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
			estadoSimulacion()
		case 6:
			fmt.Println("Saliendo del programa...")
			return
		default:
			fmt.Println("Opción no válida. Intente de nuevo.")
		}
	}
}

func main() {
	inicializarDatosPrueba()
	iniciarSimulacion()
	clearScreen()
	menuPrincipal()
}
