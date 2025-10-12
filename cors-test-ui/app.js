document.getElementById("getUsers").addEventListener("click", async () => {
    const output = document.getElementById("output");
    output.textContent = "Loading...";

    try {
        const response = await fetch("http://localhost:8080/users", {
            method: "GET",
            credentials: "include" // include cookies/sessions
        });

        if (!response.ok) {
            throw new Error(`HTTP error! Status: ${response.status}`); // at browser (client) console
        }

        const data = await response.json();
        output.textContent = JSON.stringify(data, null, 2);
    } catch (err) {
        output.textContent = "Error: " + err.message;  // at browser as response (🚫 CORS errors)
        console.error("Fetch error:", err); // at browser (client) console (🚫 CORS errors)
    }
});
