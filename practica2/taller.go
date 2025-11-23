package main

import (
	"fmt"
	"strings"
)

func estadoSimulacion() {
	mu.Lock()
	defer mu.Unlock()

	totalMecanicosActivos := 0
	for _, m := range mecanicos {
		if m.Activo {
			totalMecanicosActivos++
		}
	}

	totalIncidencias := len(incidencias)
	disponibles := len(MecanicosDisponibles)
	enCola := len(ColaDeEspera)
	enProceso := totalMecanicosActivos - disponibles
	fmt.Println("===================================")
	fmt.Println("       ESTADO DEL TALLER (CONCURRENTE)   ")
	fmt.Println("===================================")
	fmt.Printf("Total de Incidencias registradas: %d\n", totalIncidencias)
	fmt.Printf("Coches en Cola de Espera: %d\n", enCola)
	fmt.Printf("Total de Mecánicos activos: %d\n", totalMecanicosActivos)
	fmt.Printf("Mecánicos actualmente LIBRES (disponibles en canal): %d\n", disponibles)
	fmt.Printf("Mecánicos en PROCESO de trabajo (aproximado): %d\n", enProceso)

	fmt.Println("\nIncidencias Activas/Esperando:")
	for _, inc := range incidencias {
		if inc.Estado != "finalizada" {
			mecanicosAsignados := ""
			for _, m := range inc.Mecanicos {
				mecanicosAsignados += fmt.Sprintf("ID:%d, ", m.ID)
			}
			estado := inc.Estado
			if inc.esPrioritaria {
				estado = strings.ToUpper(estado)
			}
			fmt.Printf("ID: %d | Matrícula: %s | Tipo: %s | Estado: %s | Mecánicos: %s | Tiempo Acumulado: %v\n",
				inc.ID, inc.Vehiculo.Matricula, inc.Tipo, estado, strings.TrimSuffix(mecanicosAsignados, ", "), inc.TiempoAcumulado)
		}
	}
}
