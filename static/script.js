// Parse form inputs into a clean JS object
function parseFormData(form) {
    const formData = new FormData(form);
    const obj = {};

    formData.forEach((value, key) => {
        const val = value.trim();
        if (val === '') { // skip empty fields
            return;
        }

        // Convert numeric values
        if (!isNaN(val) && val !== '') {
            // Converts types: If the value looks like a number → store as Number(val) (e.g. "42" → 42).
            obj[key] = Number(val);
        } else {
            obj[key] = val; // Converts types: keep it as string.
        }
    });

    return obj;
}

// Send an AJAX request based on form data
async function sendRequest(form, urlTemplate, method) {
    const output = form.nextElementSibling; // element to show response
    const btn = form.querySelector('button[type="submit"]');
    btn.disabled = true;
    output.textContent = 'Loading...';

    const obj = parseFormData(form, method);
    let url = urlTemplate;
    let body = null;

    // Replace path parameters like {id} in the URL
    for (const key in obj) {
        if (url.includes(`{${key}}`)) {
            url = url.replace(`{${key}}`, encodeURIComponent(obj[key]));
            delete obj[key]; // remove from body/query params
        }
    }

    // Build URL query string for GET, or JSON body for others
    if (method === 'GET') {
        const params = new URLSearchParams();
        for (const key in obj) {
            params.append(key, obj[key]);
        }
        if (params.toString()) {
            url += (url.includes('?') ? '&' : '?') + params.toString();
        }
    } else {
        body = JSON.stringify(obj);
    }

    try {
        const res = await fetch(url, {
            method,
            headers: { 'Content-Type': 'application/json' },
            body: method === 'GET' ? null : body
        });

        const text = await res.text();
        output.textContent = `Status: ${res.status} ${res.statusText}\n\n${text}`;
    
    } catch (err) {
        output.textContent = 'Error: ' + err;
    } finally {
        btn.disabled = false;
    }
}

// Bind form submission to AJAX request
function bindForm(id, url, method) {
    const form = document.getElementById(id);
    if (form) {
        form.addEventListener('submit', e => {
            e.preventDefault();
            sendRequest(form, url, method);
        });
    }
}

// Bind all forms
bindForm('get-user-by-id-form', '/users/{id}', 'GET');
bindForm('create-user-form', '/users', 'POST');
bindForm('update-user-form', '/users/{id}', 'PUT');
bindForm('delete-user-form', '/users/{id}', 'DELETE');
bindForm('get-book-by-id-form', '/books/{id}', 'GET');
bindForm('create-book-form', '/books', 'POST');
bindForm('update-book-form', '/books/{id}', 'PUT');
bindForm('delete-book-form', '/books/{id}', 'DELETE');
bindForm('get-album-by-id-form', '/albums/{id}', 'GET');
bindForm('can-purchase-form', '/albums/{id}/can-purchase', 'GET');
bindForm('create-album-form', '/albums', 'POST');
bindForm('create-order-form', '/orders', 'POST');
bindForm('customer-name-form', '/customer-name', 'GET');
bindForm('json-encode-form', '/json/encode', 'POST');
bindForm('json-decode-form', '/json/decode', 'POST');
bindForm('pathfinder-form', '/pathfinder', 'GET');
bindForm('runtime-errors-form', '/runtime-errors', 'GET');

// Form validation: stops empty required fields from submitting
document.querySelectorAll("form").forEach(form => {
    form.addEventListener("submit", e => {
        let invalid = false;
        form.querySelectorAll("input[required], textarea[required]").forEach(input => {
            if (!input.value.trim()) {
                input.classList.add("invalid");
                invalid = true;
            } else {
                input.classList.remove("invalid");
            }
        });
        if (invalid) e.preventDefault();
    });
});
