let currentUsername = "";

// Function to submit the username
async function submitUsername() {
  const usernameInput = document.getElementById("username");
  const username = usernameInput.value.trim();

  if (!username) {
    displayMessage("Username cannot be empty", "error");
    return;
  }

  currentUsername = username;
  document.getElementById("user").textContent = username;
  document.getElementById("username-form-container").style.display = "none";
  document.getElementById("todos-container").style.display = "block";

  await logAction(
    "submitUsername",
    `User ${username} submitted their username`
  );
  fetchTodos(username);
}

// Function to fetch todos for the given username
async function fetchTodos(username) {
  try {
    const response = await fetch(`/todos?username=${username}`, {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
      },
    });

    if (!response.ok) {
      throw new Error(`Error: ${response.statusText}`);
    }

    const todos = await response.json();
    displayTodos(todos);
    await logAction("fetchTodos", `Fetched todos for user ${username}`);
  } catch (error) {
    console.error("Error fetching todos:", error);
    displayMessage("You have no todos, create one", "success");
  }
}

// Function to display the fetched todos
function displayTodos(todos) {
  const todosList = document.getElementById("todos-list");
  todosList.innerHTML = ""; // Clear the list
  todos.forEach((todo) => {
    const todoItem = document.createElement("li");

    const titleSpan = document.createElement("span");
    titleSpan.textContent = todo.title;
    titleSpan.className = "todo-title";

    const deleteButton = document.createElement("button");
    deleteButton.textContent = "Delete";
    deleteButton.onclick = () => deleteTodo(todo.id);

    const doneCheckbox = document.createElement("input");
    doneCheckbox.type = "checkbox";
    doneCheckbox.checked = todo.isDone;
    doneCheckbox.onchange = () => toggleTodoDone(todo.id, doneCheckbox.checked);

    todoItem.appendChild(titleSpan);
    todoItem.appendChild(doneCheckbox);
    todoItem.appendChild(deleteButton);
    todosList.appendChild(todoItem);
  });
}

// Function to add a new todo
async function addTodo() {
  const newTodoTitle = document.getElementById("new-todo-title").value.trim();
  undisplayMessage();
  if (!newTodoTitle) {
    displayMessage("Todo title cannot be empty", "error");
    return;
  }

  try {
    const response = await fetch("/todo", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ username: currentUsername, title: newTodoTitle }),
    });

    if (!response.ok) {
      throw new Error(`Error: ${response.statusText}`);
    }

    const todos = await response.json();
    displayTodos(todos);
    document.getElementById("new-todo-title").value = ""; // Clear the input
    await logAction("addTodo", `Added new todo: ${newTodoTitle}`);
  } catch (error) {
    console.error("Error adding todo:", error);
    displayMessage("Error adding todo", "error");
  }
}

// Function to delete a todo
async function deleteTodo(id) {
  try {
    const response = await fetch(`/todo?id=${id}&username=${currentUsername}`, {
      method: "DELETE",
      headers: {
        "Content-Type": "application/json",
      },
    });

    if (!response.ok) {
      throw new Error(`Error: ${response.statusText}`);
    }

    const todos = await response.json();
    displayTodos(todos);
    await logAction("deleteTodo", `Deleted todo with id: ${id}`);
  } catch (error) {
    console.error("Error deleting todo:", error);
    displayMessage("Error deleting todo", "error");
  }
}

// Function to toggle the done status of a todo
async function toggleTodoDone(id, isDone) {
  try {
    const response = await fetch(
      `/todo?id=${id}&username=${currentUsername}&isDone=${isDone}`,
      {
        method: "PUT",
        headers: {
          "Content-Type": "application/json",
        },
      }
    );

    if (!response.ok) {
      throw new Error(`Error: ${response.statusText}`);
    }

    const todos = await response.json();
    displayTodos(todos);
    await logAction(
      "toggleTodoDone",
      `Toggled todo with id: ${id} to ${isDone ? "done" : "not done"}`
    );
  } catch (error) {
    console.error("Error toggling todo status:", error);
    displayMessage("Error toggling todo status", "error");
  }
}

// Function to log user actions
async function logAction(action, message) {
  try {
    await fetch("/log", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        username: currentUsername,
        logMessage: `${action}: ${message}`,
      }),
    });
  } catch (error) {
    console.error("Error logging action:", error);
  }
}

// Function to display messages
function displayMessage(message, type) {
  const messageDiv = document.getElementById("message");
  messageDiv.textContent = message;
  messageDiv.className = type;
}
function undisplayMessage() {
  const messageDiv = document.getElementById("message");
  messageDiv.textContent = "";
  messageDiv.className = "";
}

document.addEventListener("DOMContentLoaded", () => {
  // Additional initialization if necessary
});
