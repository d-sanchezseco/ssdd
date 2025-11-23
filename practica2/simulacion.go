package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

func obtenerDuracionMedia(tipo string) time.Duration {
	var media int
	switch strings.ToLower(tipo) {
	case "mecánica":
		media = TiempoMecanica
	case "eléctrica":
		media = TiempoElectrica
	case "carrocería":
		media = TiempoCarroceria
	default:
		media = 5
	}
	desviacion := rand.Intn(media/2) + 1
	return time.Duration(media+desviacion) * time.Second
}

func administradorTaller() {
	for {
		select {
		case inc := <-NuevaIncidencia:
			go gestionarLlegada(inc)

		case disponible := <-MecanicosDisponibles:
			go gestionarDisponibilidad(disponible)

		case especialidad := <-ContratarMecanicoCanal:
			contratarMecanicoUrgente(especialidad)
		}
	}
}

func gestionarLlegada(inc Incidencia) {

	select {
	case mecanico := <-MecanicosDisponibles:
		fmt.Printf("[LLEGADA] Incidencia %d (Tipo: %s) ASIGNADA a Mecánico %d (%s).\n",
			inc.ID, inc.Tipo, mecanico.ID, mecanico.Nombre)
		go goroutineMecanico(mecanico, inc)

	default:
		ColaDeEspera <- inc
		mu.Lock()
		for i := range incidencias {
			if incidencias[i].ID == inc.ID {
				incidencias[i].Estado = "en espera"
				break
			}
		}
		mu.Unlock()
		fmt.Printf("[LLEGADA] Incidencia %d (Tipo: %s) enviada a la COLA de espera (Tamaño: %d).\n",
			inc.ID, inc.Tipo, len(ColaDeEspera))
	}
}

func gestionarDisponibilidad(disponible Mecanico) {
	select {
	case inc := <-ColaDeEspera:
		mu.Lock()
		var idx int = -1
		for i := range incidencias {
			if incidencias[i].ID == inc.ID {
				incidencias[i].Estado = "en proceso"
				idx = i
				break
			}
		}
		mu.Unlock()

		if idx != -1 {
			fmt.Printf("[DISPONIBLE] Mecánico %d (%s) ha tomado Incidencia %d de la cola.\n",
				disponible.ID, disponible.Nombre, inc.ID)
			go goroutineMecanico(disponible, inc)
		} else {
			MecanicosDisponibles <- disponible
		}

	default:
		MecanicosDisponibles <- disponible
	}
}

func goroutineMecanico(m Mecanico, incs ...Incidencia) {
	if len(incs) > 0 {
		inc := incs[0]

		mu.Lock()
		var idx int = -1
		for i, existingInc := range incidencias {
			if existingInc.ID == inc.ID {
				idx = i
				break
			}
		}
		if idx == -1 {
			mu.Unlock()
			return
		}

		if !estaMecanicoAsignado(incidencias[idx].Mecanicos, m.ID) {
			incidencias[idx].Mecanicos = append(incidencias[idx].Mecanicos, m)
		}
		incidencias[idx].Estado = "en proceso"
		mu.Unlock()

		duracionTrabajo := obtenerDuracionMedia(inc.Tipo)
		fmt.Printf("[TRABAJO] Mecánico %d (%s) inicia Incidencia %d (Tipo: %s) por %v.\n",
			m.ID, m.Nombre, inc.ID, inc.Tipo, duracionTrabajo)
		time.Sleep(duracionTrabajo)

		mu.Lock()

		var currentIdx int = -1
		for i, existingInc := range incidencias {
			if existingInc.ID == inc.ID {
				currentIdx = i
				break
			}
		}

		if currentIdx != -1 && incidencias[currentIdx].Estado != "finalizada" {

			incidencias[currentIdx].TiempoAcumulado += duracionTrabajo
			if incidencias[currentIdx].TiempoAcumulado > LimitePrioridad {

				if !incidencias[currentIdx].esPrioritaria {
					fmt.Printf("[PRIORIDAD] Incidencia %d supera el límite (%v). Buscando ayuda!\n",
						inc.ID, incidencias[currentIdx].TiempoAcumulado)
					incidencias[currentIdx].esPrioritaria = true
					incidencias[currentIdx].Estado = "prioritaria"

					select {
					case mecanicoExtra := <-MecanicosDisponibles:
						fmt.Printf("[PRIORIDAD] Mecánico extra %d (%s) asignado a Incidencia %d.\n",
							mecanicoExtra.ID, mecanicoExtra.Nombre, inc.ID)
						go goroutineMecanico(mecanicoExtra, inc)

					default:
						fmt.Printf("[PRIORIDAD] No hay mecánicos libres. Contratando un especialista en %s.\n", inc.Tipo)
						ContratarMecanicoCanal <- inc.Tipo
					}
				}

				ColaDeEspera <- incidencias[currentIdx]

			} else {
				incidencias[currentIdx].Estado = "finalizada"
				fmt.Printf("[FINALIZADA] Incidencia %d completada por Mecánico %d. Tiempo total: %v.\n",
					inc.ID, m.ID, incidencias[currentIdx].TiempoAcumulado)
			}
		}
		mu.Unlock()
	}

	MecanicosDisponibles <- m
	fmt.Printf("[DISPONIBLE] Mecánico %d (%s) ha terminado y está libre.\n", m.ID, m.Nombre)
}

func estaMecanicoAsignado(mecanicos []Mecanico, id int) bool {
	for _, m := range mecanicos {
		if m.ID == id {
			return true
		}
	}
	return false
}

func contratarMecanicoUrgente(especialidad string) {
	nuevoMecanico := Mecanico{
		ID:           0,
		Nombre:       fmt.Sprintf("Contratado-%s-%d", strings.Title(especialidad), nextMecanicoID),
		Especialidad: especialidad,
		Experiencia:  1,
		Activo:       true,
	}

	iniciarMecanico(nuevoMecanico)
}

func iniciarSimulacion() {
	ColaDeEspera = make(chan Incidencia, 1000)
	MecanicosDisponibles = make(chan Mecanico, 100)
	NuevaIncidencia = make(chan Incidencia)
	ContratarMecanicoCanal = make(chan string)

	go administradorTaller()

	for _, m := range mecanicos {
		if m.Activo {
			go goroutineMecanico(m)
			MecanicosDisponibles <- m
		}
	}
	fmt.Println("Simulación del Taller iniciada con", len(mecanicos), "mecánicos.")
}
