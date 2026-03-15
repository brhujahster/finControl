package models

import "fincontrol/db"

func CartaoTemDespesas(id int) (bool, error) {
    var count int
    err := db.DB.QueryRow(`SELECT COUNT(*) FROM despesas WHERE cartao_id = ?`, id).Scan(&count)
    if err != nil {
        return false, err
    }
    return count > 0, nil
}

func TerceiroTemDespesas(id int) (bool, error) {
    var count int
    err := db.DB.QueryRow(`SELECT COUNT(*) FROM despesas WHERE terceiro_id = ?`, id).Scan(&count)
    if err != nil {
        return false, err
    }
    return count > 0, nil
}

func TerceiroTemEmprestimos(id int) (bool, error) {
    var count int
    err := db.DB.QueryRow(`SELECT COUNT(*) FROM emprestimos WHERE terceiro_id = ?`, id).Scan(&count)
    if err != nil {
        return false, err
    }
    return count > 0, nil
}
