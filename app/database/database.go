package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/mobiquai/go_final_project/app/appsettings"
	"github.com/mobiquai/go_final_project/tests"
	_ "modernc.org/sqlite"
)

var db *sql.DB

func InitiateDb() {
	dbFilePath := getDbFilePath()
	_, err := os.Stat(dbFilePath)

	if err != nil {
		log.Printf("Файл БД не существует!")

		db, err = createDbFile(dbFilePath)
		if err != nil {
			log.Fatal(err)
		}

		createTable(db)

	} else {
		log.Printf("Файл БД успешно найден")

		db, err = sql.Open("sqlite3", dbFilePath)
		if err != nil {
			log.Fatal(err)
		}

	}

	//defer db.Close()

	db.SetMaxIdleConns(appsettings.MaxIdleConns)       // максимальное количество неактивных соединений — 2
	db.SetMaxOpenConns(appsettings.MaxOpenConns)       // максимальное общее количество открытых соединений (активных и неактивных) — 5
	db.SetConnMaxIdleTime(appsettings.MaxIdleTime)     // максимальное время, в течение которого соединение может оставаться неактивным в пуле = 5 минутам
	db.SetConnMaxLifetime(appsettings.ConnMaxLifetime) // время жизни всех соединений = 1 час

}

func getDbFilePath() string {
	dbFilePath := tests.DBFile

	envDbFilePath := appsettings.EnvDbfile // получаем значение переменной окружения
	if len(envDbFilePath) > 0 {
		dbFilePath = envDbFilePath
		log.Printf("Получен путь к файлу БД из переменной окружения TODO_DBFILE: %s", dbFilePath)
	} else {
		log.Printf("Получен путь к файлу БД из файла settings.go: %s", dbFilePath)
	}

	return dbFilePath
}

func createDbFile(dbFilePath string) (*sql.DB, error) {
	_, err := os.Create(dbFilePath)
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite3", dbFilePath)
	if err != nil {
		return nil, err
	}

	log.Printf("Файл БД успешно создан")

	return db, nil
}

func createTable(db *sql.DB) {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS scheduler (id INTEGER PRIMARY KEY AUTOINCREMENT, date VARCHAR(8) NULL, title VARCHAR(100) NOT NULL, comment VARCHAR(500) NULL, repeat VARCHAR(128) NULL)")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Таблица БД 'scheduler' успешно создана!")

	_, err = db.Exec("CREATE INDEX IF NOT EXISTS scheduler_date ON scheduler (date)")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Индекс scheduler_date таблицы 'scheduler' успешно создан!")

}

func AddTask(task Task) (int, error) {
	result, err := db.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)",
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat))

	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil

}

func TasksRead() ([]Task, error) {
	var tasks []Task

	rows, err := db.Query("SELECT * FROM scheduler ORDER BY date LIMIT :limit",
		sql.Named("limit", appsettings.SelectRowsLimit))
	if err != nil {
		return []Task{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			return []Task{}, err
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return []Task{}, err
	}

	if tasks == nil {
		tasks = []Task{}
	}

	return tasks, nil

}

func SearchTasks(search string) ([]Task, error) {
	var tasks []Task

	search = fmt.Sprintf("%%%s%%", search)
	rows, err := db.Query("SELECT * FROM scheduler WHERE title LIKE :search OR comment LIKE :search ORDER BY date LIMIT :limit",
		sql.Named("search", search),
		sql.Named("limit", appsettings.SelectRowsLimit))
	if err != nil {
		return []Task{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			return []Task{}, err
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return []Task{}, err
	}

	if len(tasks) == 0 {
		tasks = []Task{}
	}

	return tasks, nil

}

func SearchTasksByDate(date string) ([]Task, error) {
	var tasks []Task

	rows, err := db.Query("SELECT * FROM scheduler WHERE date = :date LIMIT :limit",
		sql.Named("date", date),
		sql.Named("limit", appsettings.SelectRowsLimit))
	if err != nil {
		return []Task{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			return []Task{}, err
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return []Task{}, err
	}

	if tasks == nil {
		tasks = []Task{}
	}

	return tasks, nil

}

func ReadTask(id string) (Task, error) {
	var task Task

	row := db.QueryRow("SELECT * FROM scheduler WHERE id = :id", sql.Named("id", id))
	if err := row.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
		return Task{}, errors.New("failed to read the task")
	}

	return task, nil

}

func UpdateTask(task Task) (Task, error) {
	result, err := db.Exec("UPDATE scheduler SET date = :date, title = :title, comment = :comment, repeat = :repeat WHERE id = :id",
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
		sql.Named("id", task.Id))
	if err != nil {
		return Task{}, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return Task{}, err
	}

	if rowsAffected == 0 {
		return Task{}, errors.New("failed to update the task")
	}

	return task, nil

}

func DeleteTask(id string) error {
	result, err := db.Exec("DELETE FROM scheduler WHERE id = :id", sql.Named("id", id))
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("failed to delete the task")
	}

	return err

}
