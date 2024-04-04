package main

import (
    "database/sql"
    "encoding/json"
    "log"
    "net/http"

    "github.com/go-sql-driver/mysql"
    "github.com/gorilla/mux"
)

var db *sql.DB

type Traduccion struct {
    ID           int    `json:"id"`
    PalabraMaya  string `json:"palabra_maya"`
    PalabraEsp   string `json:"palabra_esp"`
    PalabraQeqchi string `json:"palabra_qeqchi"`
}

func main() {
    // Configura la conexión con la base de datos MySQL
    config := mysql.Config{
        User:   "root",
        Passwd: "",
        Net:    "tcp",
        Addr:   "127.0.0.1:3306",
        DBName: "Imox_bd",
    }

    // Abre una conexión con la base de datos
    var err error
    db, err = sql.Open("mysql", config.FormatDSN())
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // Inicializa el enrutador de la API
    router := mux.NewRouter()

    // Define rutas de la API
    router.HandleFunc("/traducciones", getTraducciones).Methods("GET")

    // Configura el servidor HTTP
    log.Fatal(http.ListenAndServe(":8000", router))
}

func getTraducciones(w http.ResponseWriter, r *http.Request) {
    // Obtiene el idioma solicitado de la consulta
    idioma := r.URL.Query().Get("idioma")

    // Consulta las traducciones en la base de datos según el idioma especificado
    var rows *sql.Rows
    var err error
    switch idioma {
    case "qeqchi":
        rows, err = db.Query("SELECT id, palabra_maya, palabra_qeqchi FROM traducciones")
    case "espanol":
        rows, err = db.Query("SELECT id, palabra_maya, palabra_esp FROM traducciones")
    default:
        http.Error(w, "Idioma no válido", http.StatusBadRequest)
        return
    }
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()

    // Itera sobre las filas y construye un slice de Traduccion
    var traducciones []Traduccion
    for rows.Next() {
        var traduccion Traduccion
        switch idioma {
        case "qeqchi":
            if err := rows.Scan(&traduccion.ID, &traduccion.PalabraMaya, &traduccion.PalabraQeqchi); err != nil {
                log.Fatal(err)
            }
        case "espanol":
            if err := rows.Scan(&traduccion.ID, &traduccion.PalabraMaya, &traduccion.PalabraEsp); err != nil {
                log.Fatal(err)
            }
        }
        traducciones = append(traducciones, traduccion)
    }

    // Codifica el slice de Traduccion a JSON y lo envía como respuesta
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(traducciones)
}