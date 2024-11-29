// Switch tabs
document.getElementById("addUserTab").addEventListener("click", () => {
  showSection("addUserSection", "addUserTab");
});
document.getElementById("listUsersTab").addEventListener("click", () => {
  showSection("listUsersSection", "listUsersTab");
  loadUsers();
});

function showSection(sectionId, tabId) {
  document
    .querySelectorAll(".section")
    .forEach((section) => section.classList.remove("active"));
  document.getElementById(sectionId).classList.add("active");

  document
    .querySelectorAll(".sidebar nav ul li a")
    .forEach((link) => link.classList.remove("active"));
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

// // Load Users
// async function loadUsers() {
//   const res = await fetch("/api/admin/users");
//   const users = await res.json();

//   const tbody = document.getElementById("userTableBody");
//   tbody.innerHTML = "";

//   users.forEach((user) => {
//     const row = `<tr>
//       <td>${user.store_name}</td>
//       <td>${user.location}</td>
//       <td>${user.phone_number}</td>
//     </tr>`;
//     tbody.innerHTML += row;
//   });
// }

let currentPage = 1; // Initialize page

async function loadUsers(page = 1) {
  try {
    const res = await fetch(`/api/admin/users?page=${page}&limit=20`);
    const data = await res.json();

    const tbody = document.getElementById("userTableBody");
    tbody.innerHTML = "";

    const { users, total, limit, page: currentPage } = data;

    // Render table rows with serial numbers
    users.forEach((user, index) => {
      const row = `<tr>
        <td>${(page - 1) * limit + index + 1}</td>
        <td>${user.store_name}</td>
        <td>${user.location}</td>
        <td>${user.phone_number}</td>
      </tr>`;
      tbody.innerHTML += row;
    });

    // Update pagination UI
    renderPagination(total, limit, currentPage);
  } catch (error) {
    alert("Failed to load users.");
  }
}

function renderPagination(total, limit, page) {
  const pagination = document.getElementById("pagination");
  pagination.innerHTML = "";

  const totalPages = Math.ceil(total / limit);

  for (let i = 1; i <= totalPages; i++) {
    const btn = document.createElement("button");
    btn.textContent = i;
    btn.className = i === page ? "active" : "";
    btn.addEventListener("click", () => {
      loadUsers(i);
    });
    pagination.appendChild(btn);
  }
}
