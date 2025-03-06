document.addEventListener('DOMContentLoaded', () => {
    const queryInput = document.getElementById('queryInput');
    const runQueryBtn = document.getElementById('runQuery');
    const errorDiv = document.getElementById('error');
    const successDiv = document.getElementById('success');
    const resultsDiv = document.getElementById('results');
    const tableList = document.getElementById('tableList');

    // Helper function to show error messages
    const showError = (message) => {
        errorDiv.textContent = message;
        errorDiv.style.display = 'block';
        successDiv.style.display = 'none';
        resultsDiv.innerHTML = '';
    };

    // Helper function to show success messages
    const showSuccess = (message) => {
        successDiv.textContent = message;
        successDiv.style.display = 'block';
        errorDiv.style.display = 'none';
    };

    // Helper function to hide messages
    const hideMessages = () => {
        errorDiv.style.display = 'none';
        successDiv.style.display = 'none';
    };

    // Helper function to update table list
    const updateTableList = async () => {
        try {
            const response = await fetch('/tables');
            const tables = await response.json();
            
            tableList.innerHTML = '';
            tables.forEach(table => {
                const li = document.createElement('li');
                li.textContent = table;
                li.addEventListener('click', () => {
                    queryInput.value = `SELECT * FROM ${table};`;
                    executeQuery();
                });
                tableList.appendChild(li);
            });
        } catch (error) {
            console.error('Failed to fetch tables:', error);
        }
    };

    // Helper function to create table from results
    const createTable = (records) => {
        if (!records || records.length === 0) {
            resultsDiv.innerHTML = '<div class="empty-message">No results to display</div>';
            return;
        }

        const table = document.createElement('table');
        
        // Create table header
        const thead = document.createElement('thead');
        const headerRow = document.createElement('tr');
        
        // Get column names from the first record
        const columns = Object.keys(records[0].Columns || {});
        columns.forEach(key => {
            const th = document.createElement('th');
            th.textContent = key;
            headerRow.appendChild(th);
        });
        thead.appendChild(headerRow);
        table.appendChild(thead);

        // Create table body
        const tbody = document.createElement('tbody');
        records.forEach(record => {
            const row = document.createElement('tr');
            columns.forEach(key => {
                const td = document.createElement('td');
                const value = record.Columns[key];
                td.textContent = value === null ? 'NULL' : value;
                row.appendChild(td);
            });
            tbody.appendChild(row);
        });
        table.appendChild(tbody);

        resultsDiv.innerHTML = '';
        resultsDiv.appendChild(table);
    };

    // Handle query execution
    const executeQuery = async () => {
        const query = queryInput.value.trim();
        
        if (!query) {
            showError('Please enter a SQL query');
            return;
        }

        try {
            runQueryBtn.disabled = true;
            runQueryBtn.textContent = 'Running...';
            hideMessages();
            
            const response = await fetch('/query', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ query }),
            });

            const data = await response.json();

            if (!data.success) {
                showError(data.message);
                return;
            }

            // Show success message
            if (query.toUpperCase().startsWith('CREATE TABLE')) {
                showSuccess('Table created successfully');
                // Update table list after creating a new table
                await updateTableList();
            } else if (query.toUpperCase().startsWith('INSERT INTO')) {
                showSuccess('Record inserted successfully');
            } else {
                showSuccess('Query executed successfully');
            }
            
            createTable(data.records);
            
        } catch (error) {
            showError('Failed to execute query: ' + error.message);
        } finally {
            runQueryBtn.disabled = false;
            runQueryBtn.textContent = 'Run Query';
        }
    };

    // Event listeners
    runQueryBtn.addEventListener('click', executeQuery);
    
    queryInput.addEventListener('keydown', (e) => {
        // Execute query on Ctrl+Enter or Cmd+Enter
        if ((e.ctrlKey || e.metaKey) && e.key === 'Enter') {
            e.preventDefault();
            executeQuery();
        }
    });

    // Add some example queries to help users get started
    queryInput.placeholder = `Enter your SQL query here...

Example queries:
CREATE TABLE users (id INTEGER PRIMARY KEY, name TEXT, email TEXT);
INSERT INTO users (name, email) VALUES ('John Doe', 'john@example.com');
SELECT * FROM users;`;

    // Initialize table list
    updateTableList();
});
