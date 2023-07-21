const data = { errorId: '', path: '/home', size: '1', sortType: 'ASC' };
const url = '/dir-sizes';

var folders = document.querySelectorAll('.folder-list li a');

// Add a double click event listener to each link
folders.forEach(function(link) {
    link.addEventListener('dblclick', function(e) {
        e.preventDefault();
        createRequest(url);
        createFolders(arayfold);
    });
});

function createRequest(url = '', data = {}, headers = {}) {
    return fetch(url, {
        method: 'POST',
        headers: headers,
        body: JSON.stringify(data)
    })
        .then(response => response.json())
        .catch(error => console.error('Error:', error));
}

fetch(url, {
    method: 'POST',
    body: JSON.stringify(data),
    headers: {
        'Content-Type': 'application/json'
    }
}).then(response => {
    console.log('Response:', response);
}).catch(error => {
    console.error('Error:', error);
});

