document.addEventListener("DOMContentLoaded", function () {
    const loginForm = document.getElementById("loginForm");

    loginForm.addEventListener("submit", async function (event) {
        event.preventDefault();

        const formData = new FormData(loginForm);
        const data = {
            username: formData.get("username"),
            password: formData.get("password"),
        };

        try {
            const response = await fetch("/login", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify(data),
            });

            if (response.ok) {
                window.location.href = "/";
            } else {
                const result = await response.json();
                document.getElementById("error-message").textContent = result.error;
            }
        } catch (error) {
            document.getElementById("error-message").textContent = "An error occurred. Please try again.";
        }
    });
});
