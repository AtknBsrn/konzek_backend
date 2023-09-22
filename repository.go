package main

import "fmt"

func createTask(task Task) error {
	sqlStatement := "INSERT INTO Task (ID,Title, Description, Status) VALUES (?, ?, ?,?)"
	_, err := db.Exec(sqlStatement, task.ID, task.Title, task.Description, task.Status)
	if err != nil {
		fmt.Println("Failed to execute statement:", err)
		return err
	}
	fmt.Println("Success!")
	return nil
}

func getTask(id int) (Task, error) {
	var t Task
	sqlStatement := "SELECT * FROM Task WHERE ID=?"
	row := db.QueryRow(sqlStatement, id)
	err := row.Scan(&t.ID, &t.Title, &t.Description, &t.Status)
	return t, err
}

func updateTask(t Task) error {
	sqlStatement := "UPDATE Task SET Title = ?, Description = ?, Status = ? WHERE ID = ?"
	_, err := db.Exec(sqlStatement, t.Title, t.Description, t.Status, t.ID)
	return err
}

func deleteTask(id int) error {
	sqlStatement := "DELETE FROM Task WHERE ID = ?"
	_, err := db.Exec(sqlStatement, id)
	return err
}
