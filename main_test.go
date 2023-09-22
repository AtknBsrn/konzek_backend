package main

import (
	"database/sql"
	"log"
	"testing"
)

const bearer_token = "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.F7cR-QcZLXSoznCh7Fj4Uwc5kiaiy5kC4umPeeNPX5U"

var test_task = Task{ID: 333, Title: "Test Title", Description: "description", Status: "start"}

func TestCreateTask(t *testing.T) {
	create_database()
	err := createTask(test_task)
	if err != nil {
		select {
		case ErrQueue <- err:
			log.Printf("Error creating task: %v", err)
		default:
			log.Printf("Error queue is full: %v", err)
		}
	}
}

func TestGetTask(t *testing.T) {
	create_database()
	t_, err := getTask(test_task.ID)

	if err != nil {
		t.Errorf("getTask returned error")
	}

	if t_ != test_task {
		t.Errorf("got %q, wanted %q", t_, test_task)
	}
}

func TestUpdateTask(t *testing.T) {
	create_database()
	test_task.Description = "Updated Description"
	test_task.Title = "Updated Title"

	err_updt := updateTask(test_task)
	if err_updt != nil {
		t.Errorf("updateTask returned error")
	}

	t_c, err_get := getTask(test_task.ID)

	if err_get != nil {
		t.Errorf("getTask returned error")
	}

	if t_c != test_task {
		t.Errorf("got %q, wanted %q", t_c, test_task)
	}
}

func TestDeleteTask(t *testing.T) {
	create_database()

	err := deleteTask(test_task.ID)

	if err != nil {
		t.Errorf("deleteTask returned error")
	}

	t_, err_get := getTask(test_task.ID)

	if err_get == sql.ErrNoRows {
		t.Skip()
	} else {
		t.Errorf("got %q, wanted EMPTY", t_)
	}
}
