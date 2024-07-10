package main

import (
    "fmt"
    "path"
    "runtime"

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
    "db_51miz": {
        Host:     "",
        Port:     3306,
        User:     "",
        Password: "",
        Dbname:   "",
    },
    "tongji_51miz": {
        Host:     "",
        Port:     3306,
        User:     "",
        Password: "",
        Dbname:   "",
    },
    "db_51miz_config": {
        Host:     "",
        Port:     3306,
        User:     "",
        Password: "",
        Dbname:   "",
    },
}

func main() {
    // genDb51miz()
    // genDbTongji()
    // genDb51mizConfig()
}

func genDb51miz() {
    dbName := "db_51miz"
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbMap[dbName].User, dbMap[dbName].Password, dbMap[dbName].Host, dbMap[dbName].Port, dbMap[dbName].Dbname)
    gormdb, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
        NamingStrategy: schema.NamingStrategy{
            TablePrefix:   "51miz_",
            SingularTable: true,
        }})
    if err != nil {
        panic(fmt.Errorf("connect db fail: %w", err))
    }

    _, file, _, _ := runtime.Caller(0)
    rootPath := path.Dir(path.Dir(file))

    g := gen.NewGenerator(gen.Config{
        OutPath:      rootPath + "/db_51miz/query",
        ModelPkgPath: rootPath + "/db_51miz/model",
        Mode:         gen.WithDefaultQuery,
    })

    g.UseDB(gormdb)
    g.ApplyBasic(
        g.GenerateModel("51miz_user", gen.FieldType("is_freeze", "int32")),
        g.GenerateModel("51miz_userinfo"),
        g.GenerateModel("51miz_user_account"),
        g.GenerateModel("51miz_user_account_log"),
        g.GenerateModel("51miz_designer"),
        g.GenerateModel("51miz_vip_user"),
        g.GenerateModel("51miz_vip_order"),
        g.GenerateModel("51miz_templet_single_pay_order"),
        g.GenerateModel("51miz_company_user"),
        g.GenerateModel("51miz_exchange_user", gen.FieldType("status", "int32")),
        g.GenerateModel("51miz_exchange_with_order"),
        g.GenerateModel("51miz_vip_first_down"),
        g.GenerateModel("51miz_vip_user_first_download_info"),
        g.GenerateModel("51miz_out_company"),
        g.GenerateModel("51miz_out_company_vip_order", gen.FieldType("pkgtype", "int32")),
        g.GenerateModel("51miz_out_company_vip_user", gen.FieldType("pkgtype", "int32")),
        g.GenerateModel("51miz_out_company_member"),
        g.GenerateModel("51miz_new_out_company_vip_identity", gen.FieldType("mark", "int32"), gen.FieldType("status", "int32"), gen.FieldType("current", "int32")),
        g.GenerateModel("51miz_new_out_company_vip_package"),
        g.GenerateModel("51miz_single_vip_user"),
        g.GenerateModel("51miz_single_vip_download"),
        g.GenerateModel("51miz_single_vip_download_first_end"),
        g.GenerateModel("51miz_member_vip_mapping"),
        g.GenerateModel("51miz_member_pool", gen.FieldType("status", "int32")),
        g.GenerateModel("51miz_member_power", gen.FieldType("is_proxy", "int32")),
        g.GenerateModel("51miz_member_power_relation"),
        g.GenerateModel("51miz_member_trigger", gen.FieldType("status", "int32")),
        g.GenerateModel("51miz_member_channel", gen.FieldType("status", "int32")),
        g.GenerateModel("51miz_member_package", gen.FieldType("type", "int32"), gen.FieldType("is_retain", "int32"), gen.FieldType("is_online", "int32"), gen.FieldType("is_alone", "int32")),
        g.GenerateModel("51miz_member_user_power_relation", gen.FieldType("status", "int32")),
        g.GenerateModel("51miz_member_vip_user"),
        g.GenerateModel("51miz_member_coupon", gen.FieldType("status", "int32")),
        g.GenerateModel("51miz_member_coupon_user_relation", gen.FieldType("status", "int32")),
        g.GenerateModel("51miz_member_download_count"),
        g.GenerateModel("51miz_aggregated_payment"),
        g.GenerateModel("51miz_plate"),
        g.GenerateModel("51miz_vip_plate_dlimit"),
        g.GenerateModel("51miz_vip_plate_team"),
        g.GenerateModel("51miz_multi_size_image"),
        g.GenerateModel("51miz_audio"),
        g.GenerateModel("51miz_photo"),
        g.GenerateModel("51miz_photo_info"),
        g.GenerateModel("51miz_photo_size"),
        g.GenerateModel("51miz_graphslice"),
        g.GenerateModel("51miz_video"),
        g.GenerateModel("51miz_video_info"),
        g.GenerateModel("51miz_templet"),
        g.GenerateModel("51miz_templetinfo"),
        g.GenerateModel("51miz_element"),
        g.GenerateModel("51miz_element_info"),
        g.GenerateModel("51miz_font_detail"),
        g.GenerateModel("51miz_sound"),
        g.GenerateModel("51miz_video_waitmove"),
        g.GenerateModel("51miz_photo_waitmove"),
        g.GenerateModel("51miz_element_waitmove"),
        g.GenerateModel("51miz_jianzhi_power"),
        g.GenerateModel("51miz_audio_download_log"),
        g.GenerateModel("51miz_element_download_log"),
        g.GenerateModel("51miz_font_download_log"),
        g.GenerateModel("51miz_gif_download_log"),
        g.GenerateModel("51miz_photo_download_log"),
        g.GenerateModel("51miz_sound_download_log"),
        g.GenerateModel("51miz_super_download_log"),
        g.GenerateModel("51miz_designer_download_log"),
        g.GenerateModel("51miz_templet_download_log"),
        g.GenerateModel("51miz_video_download_log"),
        g.GenerateModel("51miz_ai_download_log"),
        g.GenerateModel("51miz_download_user_log_2018", gen.FieldType("ThisDate", "string")),
        g.GenerateModel("51miz_download_ip_log_2018", gen.FieldType("ThisDate", "string")),
        g.GenerateModel("51miz_size_download_log", gen.FieldType("size_type", "int32"), gen.FieldType("size_list_count", "int32"), gen.FieldType("date", "string")),
        g.GenerateModel("51miz_designer_upload_for_salary_info"),
        g.GenerateModel("51miz_designer_charging_type"),
        g.GenerateModel("51miz_salary_min_upload_user"),
        g.GenerateModel("51miz_salary_detail"),
        g.GenerateModel("51miz_search_keyword"),
    )
    g.Execute()
}

func genDbTongji() {
    dbName := "tongji_51miz"
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbMap[dbName].User, dbMap[dbName].Password, dbMap[dbName].Host, dbMap[dbName].Port, dbMap[dbName].Dbname)
    gormdb, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
        NamingStrategy: schema.NamingStrategy{
            TablePrefix:   "tongji_",
            SingularTable: true,
        }})
    if err != nil {
        panic(fmt.Errorf("connect db fail: %w", err))
    }

    _, file, _, _ := runtime.Caller(0)
    rootPath := path.Dir(path.Dir(file))

    g := gen.NewGenerator(gen.Config{
        OutPath:      rootPath + "/db_51miz_tongji/query",
        ModelPkgPath: rootPath + "/db_51miz_tongji/model",
        Mode:         gen.WithDefaultQuery,
    })

    g.UseDB(gormdb)
    g.ApplyBasic(
        g.GenerateModel("tongji_search_count_2024_05_31"),
        g.GenerateModel("tongji_search_from_2024_05_31", gen.FieldType("date", "string")),
        g.GenerateModel("tongji_page_count_2024_05_31"),
    )
    g.Execute()
}

func genDb51mizConfig() {
    dbName := "db_51miz_config"
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbMap[dbName].User, dbMap[dbName].Password, dbMap[dbName].Host, dbMap[dbName].Port, dbMap[dbName].Dbname)
    gormdb, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
        NamingStrategy: schema.NamingStrategy{
            TablePrefix:   "51miz_",
            SingularTable: true,
        }})
    if err != nil {
        panic(fmt.Errorf("connect db fail: %w", err))
    }

    _, file, _, _ := runtime.Caller(0)
    rootPath := path.Dir(path.Dir(file))

    g := gen.NewGenerator(gen.Config{
        OutPath:      rootPath + "/db_51miz_config/query",
        ModelPkgPath: rootPath + "/db_51miz_config/model",
        Mode:         gen.WithDefaultQuery,
    })

    g.UseDB(gormdb)
    g.ApplyBasic(
        g.GenerateModel("51miz_config_download_qd"),
    )
    g.Execute()
}
