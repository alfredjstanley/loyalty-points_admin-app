// Switch tabs
document.getElementById("addUserTab").addEventListener("click", () => {
  showSection("addUserSection", "addUserTab");
});
document.getElementById("listUsersTab").addEventListener("click", () => {
  showSection("listUsersSection", "listUsersTab");
  loadUsers();
});

function showSection(sectionId, tabId) {
  document.querySelectorAll(".section").forEach((section) => section.classList.remove("active"));
  document.getElementById(sectionId).classList.add("active");

  document.querySelectorAll(".sidebar nav ul li a").forEach((link) => link.classList.remove("active"));
  document.getElementById(tabId).classList.add("active");

  const sectionTitle = tabId === "addUserTab" ? "Add User" : "User List";
  document.getElementById("sectionTitle").textContent = sectionTitle;
}

// Add User
document.getElementById("addUserForm").addEventListener("submit", async (e) => {
  e.preventDefault();

  const formData = new FormData(e.target);
  const user = {
    store_name: formData.get("storeName"),
    location: formData.get("location"),
    phone_number: formData.get("phoneNumber"),
    password: formData.get("password"),
  };

  try {
    const res = await fetch("/api/admin/onboard", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(user),
    });

    if (res.ok) {
      alert("User added successfully!");
      e.target.reset();
    } else {
      const error = await res.json();
      alert("Error: " + error.message);
    }
  } catch (error) {
    alert("Failed to add user.");
  }
});

// Load Users
async function loadUsers() {
  const res = await fetch("/api/admin/users");
  const users = await res.json();

  const tbody = document.getElementById("userTableBody");
  tbody.innerHTML = "";

  users.forEach((user) => {
    const row = `<tr>
      <td>${user.store_name}</td>
      <td>${user.location}</td>
      <td>${user.phone_number}</td>
    </tr>`;
    tbody.innerHTML += row;
  });
}
