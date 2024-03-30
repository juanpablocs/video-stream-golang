package util

func CalculeInterval(duracionMinutos float64) float64 {
	const (
		intervaloBase          = 1.5  // Intervalo inicial en segundos para el primer minuto
		incrementoIntervalo    = 0.5  // Cantidad de incremento por cada minuto adicional hasta 20 minutos
		cambioIntervaloMinutos = 20   // Duración en minutos después de la cual el intervalo se fija
		intervaloPostCambio    = 10.0 // Intervalo fijo en segundos después de 20 minutos
	)

	if duracionMinutos <= cambioIntervaloMinutos {
		// Incrementar el intervalo en medio segundo por cada minuto adicional
		return intervaloBase + float64(duracionMinutos-1)*incrementoIntervalo
	} else {
		// Fijar el intervalo en 10 segundos para videos más largos que 20 minutos
		return intervaloPostCambio
	}
}
