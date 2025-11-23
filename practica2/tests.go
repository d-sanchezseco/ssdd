package main

import (
	"fmt"
	"time"
)

const (
	INCIDENCIAS_BASE = 50
	TOTAL_TEST_RUNS  = 1
)

type TestResult struct {
	Nombre      string
	Duracion    time.Duration
	Mecanicos   int
	Incidencias int
}

func configurarBase(numIncidencias int) {

	incidencias = []Incidencia{}
	mecanicos = []Mecanico{
		{ID: 1, Nombre: "Base-Luis", Especialidad: "mecánica", Experiencia: 5, Activo: true},
		{ID: 2, Nombre: "Base-Elena", Especialidad: "eléctrica", Experiencia: 7, Activo: true},
		{ID: 3, Nombre: "Base-Pablo", Especialidad: "carrocería", Experiencia: 10, Activo: true},
	}
	nextMecanicoID = 4

	iniciarSimulacion()

	enviarIncidencias(numIncidencias)
}

func configurarDuplicarIncidencias(numBase int) {
	numIncidencias := numBase * 2

	incidencias = []Incidencia{}
	mecanicos = []Mecanico{
		{ID: 1, Nombre: "Test1-Luis", Especialidad: "mecánica", Experiencia: 5, Activo: true},
		{ID: 2, Nombre: "Test1-Elena", Especialidad: "eléctrica", Experiencia: 7, Activo: true},
		{ID: 3, Nombre: "Test1-Pablo", Especialidad: "carrocería", Experiencia: 10, Activo: true},
	}
	nextMecanicoID = 4

	iniciarSimulacion()
	enviarIncidencias(numIncidencias)
}

func configurarDuplicarMecanicos(numIncidencias int) {

	incidencias = []Incidencia{}
	mecanicos = []Mecanico{

		{ID: 1, Nombre: "Test2-Luis1", Especialidad: "mecánica", Experiencia: 5, Activo: true},
		{ID: 2, Nombre: "Test2-Elena1", Especialidad: "eléctrica", Experiencia: 7, Activo: true},
		{ID: 3, Nombre: "Test2-Pablo1", Especialidad: "carrocería", Experiencia: 10, Activo: true},

		{ID: 4, Nombre: "Test2-Luis2", Especialidad: "mecánica", Experiencia: 5, Activo: true},
		{ID: 5, Nombre: "Test2-Elena2", Especialidad: "eléctrica", Experiencia: 7, Activo: true},
		{ID: 6, Nombre: "Test2-Pablo2", Especialidad: "carrocería", Experiencia: 10, Activo: true},
	}
	nextMecanicoID = 7

	iniciarSimulacion()
	enviarIncidencias(numIncidencias)
}

func configurarMecanicos311(numIncidencias int) {

	incidencias = []Incidencia{}
	mecanicos = []Mecanico{

		{ID: 1, Nombre: "Test3-M1", Especialidad: "mecánica", Experiencia: 5, Activo: true},
		{ID: 2, Nombre: "Test3-M2", Especialidad: "mecánica", Experiencia: 5, Activo: true},
		{ID: 3, Nombre: "Test3-M3", Especialidad: "mecánica", Experiencia: 5, Activo: true},
		{ID: 4, Nombre: "Test3-E1", Especialidad: "eléctrica", Experiencia: 7, Activo: true},
		{ID: 5, Nombre: "Test3-C1", Especialidad: "carrocería", Experiencia: 10, Activo: true},
	}
	nextMecanicoID = 6

	iniciarSimulacion()
	enviarIncidencias(numIncidencias)
}

func configurarMecanicos133(numIncidencias int) {

	incidencias = []Incidencia{}
	mecanicos = []Mecanico{

		{ID: 1, Nombre: "Test4-M1", Especialidad: "mecánica", Experiencia: 5, Activo: true},
		{ID: 2, Nombre: "Test4-E1", Especialidad: "eléctrica", Experiencia: 7, Activo: true},
		{ID: 3, Nombre: "Test4-E2", Especialidad: "eléctrica", Experiencia: 7, Activo: true},
		{ID: 4, Nombre: "Test4-E3", Especialidad: "eléctrica", Experiencia: 7, Activo: true},
		{ID: 5, Nombre: "Test4-C1", Especialidad: "carrocería", Experiencia: 10, Activo: true},
		{ID: 6, Nombre: "Test4-C2", Especialidad: "carrocería", Experiencia: 10, Activo: true},
		{ID: 7, Nombre: "Test4-C3", Especialidad: "carrocería", Experiencia: 10, Activo: true},
	}
	nextMecanicoID = 8

	iniciarSimulacion()
	enviarIncidencias(numIncidencias)
}

func enviarIncidencias(count int) {
	tipos := []string{"mecánica", "eléctrica", "carrocería"}
	vehiculoBase := vehiculos[0]

	mu.Lock()
	defer mu.Unlock()

	fmt.Printf("\nEnviando %d incidencias al canal NuevaIncidencia...\n", count)

	for i := 0; i < count; i++ {
		inc := Incidencia{
			ID:          nextIncidenciaID,
			Tipo:        tipos[i%len(tipos)],
			Prioridad:   "normal",
			Descripcion: fmt.Sprintf("Incidencia automática %d", nextIncidenciaID),
			Estado:      "en espera",
			Vehiculo:    vehiculoBase,
		}
		nextIncidenciaID++
		incidencias = append(incidencias, inc)

		NuevaIncidencia <- inc
	}
	fmt.Println("Todas las incidencias enviadas. El Administrador del Taller está trabajando.")
}

func monitorearIncidencias(targetCount int) {
	start := time.Now()

	for {
		time.Sleep(500 * time.Millisecond)

		mu.Lock()
		terminadas := 0
		for _, inc := range incidencias {
			if inc.Estado == "finalizada" {
				terminadas++
			}
		}

		if terminadas >= targetCount {
			mu.Unlock()
			fmt.Printf("\n SIMULACIÓN FINALIZADA. Incidencias completadas: %d\n", targetCount)
			return
		}

		tiempoTranscurrido := time.Since(start)
		fmt.Printf("... Procesando: %d/%d completadas. Tiempo: %v\n", terminadas, targetCount, tiempoTranscurrido.Round(time.Second))

		if tiempoTranscurrido > 3*time.Minute {
			mu.Unlock()
			fmt.Println(" SIMULACIÓN INTERRUMPIDA: Se superó el límite de 3 minutos. Revise la configuración.")
			return
		}

		mu.Unlock()
	}
}

func ejecutarTests() {
	var resultados []TestResult

	fmt.Println("====================================================")
	fmt.Println("      INICIANDO BATERÍA DE TESTS DE CONCURRENCIA    ")
	fmt.Println("====================================================")

	fmt.Println("\n--- COMPARATIVA 1: N INCIDENCIAS VS 2N INCIDENCIAS ---")

	fmt.Printf("\n[C1. BASE] Configuración: 3 Mecánicos (1:1:1) | Incidencias: %d\n", INCIDENCIAS_BASE)
	configurarBase(INCIDENCIAS_BASE)
	start1 := time.Now()
	monitorearIncidencias(INCIDENCIAS_BASE)
	duracion1 := time.Since(start1)
	resultados = append(resultados, TestResult{"C1. Base (N Incid.)", duracion1, 3, INCIDENCIAS_BASE})

	fmt.Printf("\n[C1. TEST] Configuración: 3 Mecánicos (1:1:1) | Incidencias: %d\n", INCIDENCIAS_BASE*2)
	configurarDuplicarIncidencias(INCIDENCIAS_BASE)
	start2 := time.Now()
	monitorearIncidencias(INCIDENCIAS_BASE * 2)
	duracion2 := time.Since(start2)
	resultados = append(resultados, TestResult{"C1. Test (2N Incid.)", duracion2, 3, INCIDENCIAS_BASE * 2})

	fmt.Println("\n--- COMPARATIVA 2: 3 MECÁNICOS VS 6 MECÁNICOS ---")

	fmt.Printf("\n[C2. TEST] Configuración: 6 Mecánicos (2:2:2) | Incidencias: %d\n", INCIDENCIAS_BASE)
	configurarDuplicarMecanicos(INCIDENCIAS_BASE)
	start3 := time.Now()
	monitorearIncidencias(INCIDENCIAS_BASE)
	duracion3 := time.Since(start3)
	resultados = append(resultados, TestResult{"C2. Test (6 Mec.)", duracion3, 6, INCIDENCIAS_BASE})

	fmt.Println("\n--- COMPARATIVA 3: BALANCEO DE ESPECIALIDADES ---")

	fmt.Printf("\n[C3. TEST A] Configuración: 5 Mecánicos (3:1:1) | Incidencias: %d\n", INCIDENCIAS_BASE)
	configurarMecanicos311(INCIDENCIAS_BASE)
	start4 := time.Now()
	monitorearIncidencias(INCIDENCIAS_BASE)
	duracion4 := time.Since(start4)
	resultados = append(resultados, TestResult{"C3. Test (3:1:1)", duracion4, 5, INCIDENCIAS_BASE})

	fmt.Printf("\n[C3. TEST B] Configuración: 7 Mecánicos (1:3:3) | Incidencias: %d\n", INCIDENCIAS_BASE)
	configurarMecanicos133(INCIDENCIAS_BASE)
	start5 := time.Now()
	monitorearIncidencias(INCIDENCIAS_BASE)
	duracion5 := time.Since(start5)
	resultados = append(resultados, TestResult{"C3. Test (1:3:3)", duracion5, 7, INCIDENCIAS_BASE})

	fmt.Println("\n====================================================")
	fmt.Println("            RESUMEN DE TIEMPOS DE EJECUCIÓN           ")
	fmt.Println("====================================================")
	fmt.Printf("| %-20s | %-12s | %-12s | %-15s |\n", "Escenario", "Mecánicos", "Incidencias", "Duración Total")
	fmt.Println("----------------------|-------------|--------------|-----------------")
	for _, res := range resultados {
		fmt.Printf("| %-20s | %-12d| %-12d | %-15v |\n", res.Nombre, res.Mecanicos, res.Incidencias, res.Duracion.Round(time.Second))
	}
	fmt.Println("====================================================")
}

// Función que debe llamarse desde main.go para iniciar los tests
func mainTest() {
	inicializarDatosPrueba()
	ejecutarTests()
}
