// DOM Elements
const queryInput = document.getElementById('queryInput');
const runQueryBtn = document.getElementById('runQuery');
const resultsDiv = document.getElementById('results');
const successDiv = document.getElementById('success');
const errorDiv = document.getElementById('error');
const tableList = document.getElementById('tableList');

// Helper functions
const hideMessages = () => {
    successDiv.textContent = '';
    successDiv.style.display = 'none';
    errorDiv.textContent = '';
    errorDiv.style.display = 'none';
};

const showSuccess = (message) => {
    successDiv.textContent = message;
    successDiv.style.display = 'block';
    errorDiv.style.display = 'none';
};

const showError = (message) => {
    errorDiv.textContent = message;
    errorDiv.style.display = 'block';
    successDiv.style.display = 'none';
};

// Helper function to create table from results
const createTable = (records, columns) => {
    console.log('Creating table with:', { records, columns });

    // Clear previous results
    resultsDiv.innerHTML = '';

    if (!records || records.length === 0) {
        const emptyMsg = document.createElement('div');
        emptyMsg.className = 'empty-message';
        emptyMsg.textContent = 'No records found';
        resultsDiv.appendChild(emptyMsg);
        return;
    }

    const table = document.createElement('table');
    table.className = 'data-table';
    
    // Create table header
    const thead = document.createElement('thead');
    const headerRow = document.createElement('tr');
    
    // Get column names from the result
    const columnNames = columns || Object.keys(records[0]);
    columnNames.forEach(key => {
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
        columnNames.forEach(key => {
            const td = document.createElement('td');
            const value = record[key];
            td.textContent = value === null ? 'NULL' : String(value);
            row.appendChild(td);
        });
        tbody.appendChild(row);
    });
    table.appendChild(tbody);

    // Append table to results div
    resultsDiv.appendChild(table);
    console.log('Table created and appended');
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
        resultsDiv.innerHTML = '';
        
        const response = await fetch('/query', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ query }),
        });

        const data = await response.json();
        console.log('Query response:', data);

        if (!data.success) {
            showError(data.message);
            return;
        }

        // Show success message based on query type
        const upperQuery = query.toUpperCase().trim();
        if (upperQuery.startsWith('CREATE TABLE')) {
            showSuccess('Table created successfully');
            await updateTableList();
        } else if (upperQuery.startsWith('INSERT INTO')) {
            showSuccess('Record inserted successfully');
        } else if (upperQuery.startsWith('DELETE')) {
            showSuccess(data.message || 'Records deleted successfully');
        } else if (upperQuery.startsWith('SELECT')) {
            showSuccess(data.message || 'Query executed successfully');
            createTable(data.records, data.columns);
        } else {
            showSuccess(data.message || 'Query executed successfully');
        }
        
    } catch (error) {
        console.error('Query error:', error);
        showError('Failed to execute query: ' + error.message);
    } finally {
        runQueryBtn.disabled = false;
        runQueryBtn.textContent = 'Run Query';
    }
};

// Update table list
const updateTableList = async () => {
    try {
        const response = await fetch('/tables');
        const data = await response.json();
        
        tableList.innerHTML = '';
        data.tables.forEach(table => {
            const button = document.createElement('button');
            button.className = 'table-button';
            button.innerHTML = `
                <div class="table-button-content">
                    <span class="table-icon">📋</span>
                    <span class="table-name">${table}</span>
                </div>
            `;
            button.addEventListener('click', () => {
                queryInput.value = `SELECT * FROM ${table};`;
                executeQuery();
            });
            tableList.appendChild(button);
        });
    } catch (error) {
        console.error('Failed to update table list:', error);
    }
};

// Event listeners
runQueryBtn.addEventListener('click', executeQuery);
queryInput.addEventListener('keydown', (e) => {
    if (e.key === 'Enter' && e.ctrlKey) {
        executeQuery();
    }
});

// Initialize
document.addEventListener('DOMContentLoaded', () => {
    updateTableList();
    
    // Add example queries
    queryInput.placeholder = `Enter your SQL query here...

Example queries:
CREATE TABLE users (id INTEGER PRIMARY KEY, name TEXT, email TEXT);
INSERT INTO users (name, email) VALUES ('John Doe', 'john@example.com');
SELECT * FROM users;`;
});
