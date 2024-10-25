package test

import (
    "fmt"
    "os"

    "github.com/jackc/pgx"
)

func LoadTest1() {
    pgxConConfig := pgx.ConnConfig{
        Port:     5433,
        Host:     "localhost",
        Database: "postgres",
        User:     "postgres",
        Password: "password",
    }

    conn, err := pgx.Connect(pgxConConfig)
    if err != nil {
        panic(err)
    }
    defer conn.Close()

    // tables := []string{"table1", "table2", "table3",}
    tables := []string{"test"}

    import_dir := "test"
    // export_dir := "test"
    
    for _, t := range tables {
        // f, err := os.OpenFile(fmt.Sprintf("%s/table%s.csv", import_dir, t), os.O_RDONLY, 0777)
        f, err := os.OpenFile(fmt.Sprintf("test.csv", import_dir, t), os.O_RDONLY, 0777)
        if err != nil {
            return
        }
        f.Close()

        err = importer(conn, f, t)
        if err != nil {
            break
        }

        // fmt.Println("  Done with import and doing export")
        // ef, err := os.OpenFile(fmt.Sprintf("%s/table_%s.csv", export_dir, t), os.O_CREATE|os.O_WRONLY, 0777)
        // if err != nil {
        //     fmt.Println("error opening file:", err)
        //     return
        // }
        // ef.Close()
        //
        // err = exporter(conn, ef, t)
        // if err != nil {
        //     break
        // }
    }
}

func importer(conn *pgx.Conn, f *os.File, table string) error {
    res, err := conn.CopyFromReader(f, fmt.Sprintf("COPY %s FROM STDIN DELIMITER '|' CSV HEADER", table))
    if err != nil {
        return err
    }
    fmt.Println("==> import rows affected:", res.RowsAffected())

    return nil
}

func exporter(conn *pgx.Conn, f *os.File, table string) error {
    res, err := conn.CopyToWriter(f, fmt.Sprintf("COPY %s TO STDOUT DELIMITER '|' CSV HEADER", table))
    if err != nil {
        return fmt.Errorf("error exporting file: %+v", err)
    }
    fmt.Println("==> export rows affected:", res.RowsAffected())
    return nil
}
