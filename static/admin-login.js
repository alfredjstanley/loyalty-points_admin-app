async function checkLoginStatus() {
  const adminToken = sessionStorage.getItem("adminToken");

  if (adminToken) {
    try {
      // Use the token to request the admin page
      const adminResponse = await fetch("/admin/home", {
        method: "GET",
        headers: { Authorization: adminToken },
      });

      if (adminResponse.ok) {
        const html = await adminResponse.text();
        document.open();
        document.write(html);
        document.close();
      } else {
        // Token is invalid or expired, clear the token and show the login page
        sessionStorage.removeItem("adminToken");
        console.error("Token invalid or expired.");
      }
    } catch (err) {
      console.error("Error checking login status:", err);
      sessionStorage.removeItem("adminToken");
    }
  }
}

async function handleLogin(username, password) {
  try {
    const response = await fetch("/api/admin/login", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ username, password }),
    });

    if (response.ok) {
      const data = await response.json();
      const adminToken = data.token;
      sessionStorage.setItem("adminToken", adminToken);

      // Use the token to request and load the admin HTML page
      const adminResponse = await fetch("/admin/home", {
        method: "GET",
        headers: { Authorization: adminToken },
      });

      if (adminResponse.ok) {
        const html = await adminResponse.text();
        document.open();
        document.write(html);
        document.close();
      } else {
        throw new Error("Failed to load admin page.");
      }
    } else {
      const error = await response.json();
      document.getElementById("errorMessage").textContent =
        error.message || "Login failed";
    }
  } catch (err) {
    document.getElementById("errorMessage").textContent =
      "An error occurred. Please try again.";
    console.error(err);
  }
}

// Check login status on page load
document.addEventListener("DOMContentLoaded", checkLoginStatus);

// Login form submission handler
document.getElementById("loginForm").addEventListener("submit", async (e) => {
  e.preventDefault();

  const username = document.getElementById("username").value;
  const password = document.getElementById("password").value;

  await handleLogin(username, password);
});
