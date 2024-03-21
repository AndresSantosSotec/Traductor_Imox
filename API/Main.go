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
    // Consulta las traducciones en la base de datos
    rows, err := db.Query("SELECT id, palabra_maya, palabra_esp FROM traducciones")
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()

    // Itera sobre las filas y construye un slice de Traduccion
    var traducciones []Traduccion
    for rows.Next() {
        var traduccion Traduccion
        if err := rows.Scan(&traduccion.ID, &traduccion.PalabraMaya, &traduccion.PalabraEsp); err != nil {
            log.Fatal(err)
        }
        traducciones = append(traducciones, traduccion)
    }

    // Codifica el slice de Traduccion a JSON y lo envía como respuesta
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(traducciones)
}
