// Flexplane - Progressive Enhancement JavaScript
// Adds interactivity to server-rendered content

document.addEventListener('DOMContentLoaded', function() {
    initializeDateDisplay();
    initializeTodoInteractivity();
});

// Update header with current date
function initializeDateDisplay() {
    const dateElement = document.getElementById('current-date');
    if (dateElement) {
        const now = new Date();
        dateElement.textContent = now.toLocaleDateString('en-US', {
            weekday: 'long',
            year: 'numeric',
            month: 'long',
            day: 'numeric'
        });
    }
}

// Todo interactivity
function initializeTodoInteractivity() {
    // Add todo form submission
    const addButton = document.getElementById('add-todo-btn');
    const newTodoInput = document.getElementById('new-todo');

    if (addButton && newTodoInput) {
        // Button click
        addButton.addEventListener('click', handleAddTodo);

        // Enter key in input
        newTodoInput.addEventListener('keypress', function(e) {
            if (e.key === 'Enter') {
                handleAddTodo();
            }
        });
    }

    // Todo checkbox toggles
    document.querySelectorAll('.todo-checkbox').forEach(checkbox => {
        checkbox.addEventListener('change', handleToggleTodo);
    });
}

async function handleAddTodo() {
    const input = document.getElementById('new-todo');
    const message = input.value.trim();

    if (!message) return;

    try {
        const response = await fetch('/api/todos', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ message: message })
        });

        if (response.ok) {
            // Optimistic UI update - reload page to get fresh data
            window.location.reload();
        } else {
            console.error('Failed to add todo');
        }
    } catch (error) {
        console.error('Error adding todo:', error);
    }
}

async function handleToggleTodo(event) {
    const checkbox = event.target;
    const todoItem = checkbox.closest('.todo-item');
    const index = todoItem.dataset.index;

    // Optimistic UI update
    todoItem.classList.toggle('completed', checkbox.checked);

    try {
        const response = await fetch(`/api/todos?index=${index}`, {
            method: 'PATCH'
        });

        if (!response.ok) {
            // Revert optimistic update on failure
            checkbox.checked = !checkbox.checked;
            todoItem.classList.toggle('completed', checkbox.checked);
            console.error('Failed to toggle todo');
        }
    } catch (error) {
        // Revert optimistic update on failure
        checkbox.checked = !checkbox.checked;
        todoItem.classList.toggle('completed', checkbox.checked);
        console.error('Error toggling todo:', error);
    }
}