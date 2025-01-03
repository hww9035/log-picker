package main

import (
    "flag"
    "fmt"

    "gorm.io/driver/mysql"
    "gorm.io/gen"
    "gorm.io/gorm"
    "gorm.io/gorm/schema"
)

type dbConfig struct {
    Host     string
    Port     int32
    User     string
    Password string
    Dbname   string
}

var dbMap = map[string]dbConfig{
    "db1": {
        Host:     "",
        Port:     3307,
        User:     "",
        Password: "",
        Dbname:   "",
    },
    "db2": {
        Host:     "",
        Port:     3306,
        User:     "",
        Password: "",
        Dbname:   "",
    },
    "db3": {
        Host:     "",
        Port:     3307,
        User:     "",
        Password: "",
        Dbname:   "",
    },
}

var modelPath = ""

func main() {
    flag.StringVar(&modelPath, "p", "", "please set modelPath parameter.")
    flag.Parse()
    if modelPath == "" {
        modelPath = "."
    }
    // genDb()
    // genDbTongji()
    // genDbConfig()
}

func getDsn(dbName string) string {
    return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbMap[dbName].User, dbMap[dbName].Password, dbMap[dbName].Host, dbMap[dbName].Port, dbMap[dbName].Dbname)
}

func genDb() {
    dbName := "db1"
    dsn := getDsn(dbName)
    gormdb, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
        NamingStrategy: schema.NamingStrategy{
            TablePrefix:   "prefix_",
            SingularTable: true,
        }})
    if err != nil {
        panic(fmt.Errorf("connect db fail: %w", err))
    }

    // _, file, _, _ := runtime.Caller(0)
    // modelPath = path.Dir(path.Dir(file))
    fmt.Printf("genDb modelPath: %s\n", modelPath)

    g := gen.NewGenerator(gen.Config{
        OutPath:      modelPath + "/" + dbName + "/query",
        ModelPkgPath: modelPath + "/" + dbName + "/model",
        Mode:         gen.WithDefaultQuery,
    })

    g.UseDB(gormdb)
    g.ApplyBasic(
        g.GenerateModel("user", gen.FieldType("is_freeze", "int32")),
        g.GenerateModel("download_log", gen.FieldType("size_type", "int32"), gen.FieldType("size_list_count", "int32"), gen.FieldType("date", "string")),
    )
    g.Execute()
}

func genDbTongji() {
    dbName := "db2"
    dsn := getDsn(dbName)
    gormdb, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
        NamingStrategy: schema.NamingStrategy{
            TablePrefix:   "tongji_",
            SingularTable: true,
        }})
    if err != nil {
        panic(fmt.Errorf("connect db fail: %w", err))
    }

    fmt.Printf("genDbTongji modelPath: %s\n", modelPath)

    g := gen.NewGenerator(gen.Config{
        OutPath:      modelPath + "/db_tongji/query",
        ModelPkgPath: modelPath + "/db_tongji/model",
        Mode:         gen.WithDefaultQuery,
    })

    g.UseDB(gormdb)
    g.ApplyBasic(
        g.GenerateModel("table1"),
        g.GenerateModel("table2", gen.FieldType("date", "string")),
    )
    g.Execute()
}

func genDbConfig() {
    dbName := "db3"
    dsn := getDsn(dbName)
    gormdb, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
        NamingStrategy: schema.NamingStrategy{
            TablePrefix:   "prefix_",
            SingularTable: true,
        }})
    if err != nil {
        panic(fmt.Errorf("connect db fail: %w", err))
    }

    fmt.Printf("genDbConfig modelPath: %s\n", modelPath)

    g := gen.NewGenerator(gen.Config{
        OutPath:      modelPath + "/" + dbName + "/query",
        ModelPkgPath: modelPath + "/" + dbName + "/model",
        Mode:         gen.WithDefaultQuery,
    })

    g.UseDB(gormdb)
    g.ApplyBasic(
        g.GenerateModel("table1"),
    )
    g.Execute()
}
