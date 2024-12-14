// Switch tabs
document.getElementById("addUserTab").addEventListener("click", () => {
  showSection("addUserSection", "addUserTab");
});
document.getElementById("listUsersTab").addEventListener("click", () => {
  showSection("listUsersSection", "listUsersTab");
  loadUsers(); // Load users when navigating to the Merchant List
});

// Handle sidebar navigation
document.getElementById("transactionLogsTab").addEventListener("click", () => {
  showSection("transactionLogsSection", "transactionLogsTab");
  loadTransactionLogs(); // Load logs when section is clicked
});

// Close modal
document
  .querySelector("#viewModal .close-btn")
  .addEventListener("click", () => {
    document.getElementById("viewModal").style.display = "none";
  });

// Close Edit Modal
document.querySelector(".close-btn").addEventListener("click", () => {
  document.getElementById("editModal").style.display = "none";
});

// Submit Edit Form
document.getElementById("editForm").addEventListener("submit", async (e) => {
  e.preventDefault();
  const id = document.getElementById("editMerchantId").value;
  const store_name = document.getElementById("editStoreName").value;
  const location = document.getElementById("editLocation").value;
  const password = document.getElementById("editPassword").value;

  try {
    const adminToken = sessionStorage.getItem("adminToken");

    await fetch("/api/admin/edit-merchant", {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
        Authorization: adminToken,
      },
      body: JSON.stringify({ id, store_name, location, password }),
    });
    alert("Merchant updated successfully!");
    document.getElementById("editModal").style.display = "none";
    loadUsers();
  } catch (error) {
    alert("Failed to update merchant.");
  }
});

// Add Merchant
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
    const adminToken = sessionStorage.getItem("adminToken");

    await fetch("/api/admin/onboard", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: adminToken,
      },
      body: JSON.stringify(user),
    });
    alert("Merchant added successfully!");
    e.target.reset();
  } catch (error) {
    alert("Failed to add merchant.");
  }
});

// Add event listener for button click
document
  .getElementById("searchButton")
  .addEventListener("click", triggerSearch);

// Add event listener for 'Enter' key press
document.getElementById("searchInput").addEventListener("keydown", (event) => {
  if (event.key === "Enter") {
    triggerSearch();
  }
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
}

// Load Merchants with Pagination
async function loadUsers(page = 1) {
  try {
    const adminToken = sessionStorage.getItem("adminToken");

    const res = await fetch(`/api/admin/users?page=${page}&limit=10`, {
      method: "GET",
      headers: { Authorization: adminToken },
    });
    const { users, total, limit } = await res.json();

    const tbody = document.getElementById("userTableBody");
    tbody.innerHTML = ""; // Clear the table before rendering new data

    const totalCountElement = document.getElementById("totalMerchants");
    totalCountElement.textContent = `Total Merchants: ${total}`;

    users.forEach((user, index) => {
      tbody.innerHTML += `
        <tr>
          <td>${(page - 1) * limit + index + 1}</td>
          <td>${user.store_name}</td>
          <td>${user.location}</td>
          <td>${user.phone_number}</td>
          <td class="table-actions">
            <button class="btn view-btn" data-id="${
              user.phone_number
            }">View</button>
            <button class="btn edit-btn" data-id="${user.id}">Edit</button>
          </td>
        </tr>`;
    });

    renderPagination(total, limit, page);
  } catch (error) {
    alert("Failed to load users.");
  }
}

// Render Pagination
function renderPagination(total, limit, currentPage) {
  const pagination = document.getElementById("pagination");
  pagination.innerHTML = ""; // Clear existing pagination

  const totalPages = Math.ceil(total / limit);
  if (totalPages <= 1) return;

  const createButton = (text, disabled, onClick) => {
    const button = document.createElement("button");
    button.textContent = text;
    button.disabled = disabled;
    if (onClick) button.addEventListener("click", onClick);
    return button;
  };

  // Previous Button
  pagination.appendChild(
    createButton("Previous", currentPage === 1, () =>
      loadUsers(currentPage - 1)
    )
  );

  // Page Buttons
  for (let i = 1; i <= totalPages; i++) {
    const pageButton = createButton(i, false, () => loadUsers(i));
    if (i === currentPage) pageButton.className = "active";
    pagination.appendChild(pageButton);
  }

  // Next Button
  pagination.appendChild(
    createButton("Next", currentPage === totalPages, () =>
      loadUsers(currentPage + 1)
    )
  );
}

// search fn
function triggerSearch() {
  const query = document.getElementById("searchInput").value.trim();
  if (query) {
    searchMerchants(query);
  } else {
    loadUsers();
  }
}

// Edit Merchant
document.addEventListener("click", (e) => {
  if (e.target.classList.contains("edit-btn")) {
    const row = e.target.closest("tr");
    document.getElementById("editMerchantId").value = e.target.dataset.id;
    document.getElementById("editStoreName").value =
      row.children[1].textContent;
    document.getElementById("editLocation").value = row.children[2].textContent;
    document.getElementById("editModal").style.display = "block";
  }
});

document.addEventListener("click", async (e) => {
  if (e.target.classList.contains("view-btn")) {
    const mobileNumber = e.target.dataset.id; // Get the merchant ID
    const adminToken = sessionStorage.getItem("adminToken");

    try {
      // Fetch merchant details from the API
      const res = await fetch(
        `/api/admin/edit-merchant?mobileNumber=${mobileNumber}`,
        {
          method: "GET",
          headers: {
            "Content-Type": "application/json",
            Authorization: adminToken,
          },
        }
      );

      const { success, data, message } = await res.json();

      if (success) {
        // Populate the modal with merchant details
        const details = document.getElementById("merchantDetails");
        details.innerHTML = `
          <img src="${data.profile_picture}" alt="${
          data.store_name
        }" width="100" style="border-radius: 50%;" />
          <p><strong>Name:</strong> ${data.first_name} ${data.last_name}</p>
          <p><strong>Email:</strong> ${data.email}</p>
          <p><strong>Phone:</strong> ${data.mobile_number}</p>
          <p><strong>Store:</strong> ${data.store_name}</p>
          <p><strong>Account Type:</strong> ${data.account_type}</p>
          <p><strong>Location:</strong> ${data.address}</p>
          <p><strong>Working Hours:</strong> ${data.working_hours}</p>
          <p><strong>Balance:</strong> ${data.point_balance} points</p>
          <p><strong>Categories:</strong> ${data.product_categories.join(
            ", "
          )}</p>
          <p><strong>Days Open:</strong> ${data.working_days.join(", ")}</p>
        `;

        // Open the modal
        document.getElementById("viewModal").style.display = "block";
      } else {
        alert("Failed to fetch merchant details.");
      }
    } catch (error) {
      console.error("Error fetching merchant details:", error);
      alert("Something went wrong while fetching details.");
    }
  }
});

// Load transaction logs
async function loadTransactionLogs(page = 1) {
  try {
    const adminToken = sessionStorage.getItem("adminToken");

    const res = await fetch(
      `/api/admin/transaction-logs?page=${page}&limit=10`,
      {
        method: "GET",
        headers: { Authorization: adminToken },
      }
    );
    const data = await res.json(); // Parse the response as JSON

    // Confirm the structure of `data.logs` (if `logs` is nested)
    const logs = data.logs || [];

    const tbody = document.getElementById("logsTableBody");
    tbody.innerHTML = "";

    logs.forEach((log, index) => {
      const date = new Date(log.CreatedAt).toLocaleString(); // Format date
      const row = `
        <tr>
          <td>${(page - 1) * 10 + index + 1}</td>
          <td>${log.UserPhone}</td>
          <td>${log.MerchantPhone}</td>
          <td>${log.Amount}</td>
          <td>${log.InvoiceID}</td>
          <td>${log.Status}</td>
          <td>${date}</td>
        </tr>`;
      tbody.innerHTML += row;
    });

    renderLogsPagination(data.total, data.limit, page);
  } catch (error) {
    alert("Failed to load transaction logs.");
    console.error(error);
  }
}

// Render pagination buttons
function renderLogsPagination(total, limit, page) {
  const pagination = document.getElementById("logsPagination");
  pagination.innerHTML = ""; // Clear existing pagination

  const totalPages = Math.ceil(total / limit);

  if (totalPages <= 1) return; // No pagination needed for a single page

  // Add Previous button
  const prevButton = document.createElement("button");
  prevButton.textContent = "Previous";
  prevButton.disabled = page === 1; // Disable if on the first page
  prevButton.addEventListener("click", () => loadTransactionLogs(page - 1));
  pagination.appendChild(prevButton);

  // Add numbered page buttons
  for (let i = 1; i <= totalPages; i++) {
    const pageButton = document.createElement("button");
    pageButton.textContent = i;
    pageButton.className = i === page ? "active" : ""; // Highlight current page
    pageButton.addEventListener("click", () => loadTransactionLogs(i));
    pagination.appendChild(pageButton);
  }

  // Add Next button
  const nextButton = document.createElement("button");
  nextButton.textContent = "Next";
  nextButton.disabled = page === totalPages; // Disable if on the last page
  nextButton.addEventListener("click", () => loadTransactionLogs(page + 1));
  pagination.appendChild(nextButton);
}

async function searchMerchants(query) {
  try {
    const adminToken = sessionStorage.getItem("adminToken");

    const res = await fetch(
      `/api/admin/merchants/search?query=${encodeURIComponent(query)}`,
      {
        method: "GET",
        headers: { Authorization: adminToken },
      }
    );
    const { merchants, total } = await res.json();

    const tbody = document.getElementById("userTableBody");
    tbody.innerHTML = "";

    // Render search results
    merchants.forEach((merchant, index) => {
      const row = `
        <tr>
          <td>${index + 1}</td>
          <td>${merchant.store_name}</td>
          <td>${merchant.location}</td>
          <td>${merchant.phone_number}</td>
          <td class="table-actions">
            <button class="btn edit-btn" data-id="${merchant.id}">Edit</button>
            <button class="btn view-btn" data-id="${
              merchant.phone_number
            }")">View</button>
          </td>
        </tr>`;
      tbody.innerHTML += row;
    });

    document.getElementById(
      "totalMerchants"
    ).textContent = `Total Results: ${total}`;
  } catch (error) {
    alert("Failed to search merchants.");
    console.error(error);
  }
}

// View Merchant Details
function viewMerchant(id) {
  // Fetch merchant details and display in modal
  const adminToken = sessionStorage.getItem("adminToken");

  fetch(`/api/admin/merchants/${id}`, {
    method: "GET",
    headers: { Authorization: adminToken },
  })
    .then((response) => response.json())
    .then((data) => {
      const details = `
        <p><strong>Store Name:</strong> ${data.store_name}</p>
        <p><strong>Location:</strong> ${data.location}</p>
        <p><strong>Phone Number:</strong> ${data.phone_number}</p>
        <p><strong>Status:</strong> ${data.status}</p>
      `;
      document.getElementById("merchantDetails").innerHTML = details;
      document.getElementById("viewModal").style.display = "block";
    })
    .catch((error) => {
      alert("Failed to fetch merchant details.");
      console.error(error);
    });
}

// Edit Merchant
function editMerchant(id) {
  // Fetch merchant details and prefill the edit form
  const adminToken = sessionStorage.getItem("adminToken");

  fetch(`/api/admin/merchants/${id}`, {
    method: "GET",
    headers: { Authorization: adminToken },
  })
    .then((response) => response.json())
    .then((data) => {
      document.getElementById("editMerchantId").value = data.id;
      document.getElementById("editStoreName").value = data.store_name;
      document.getElementById("editLocation").value = data.location;
      document.getElementById("editPassword").value = ""; // Keep password empty
      document.getElementById("editModal").style.display = "block";
    })
    .catch((error) => {
      alert("Failed to fetch merchant details.");
      console.error(error);
    });
}

// Logout functionality
document.getElementById("logoutButton").addEventListener("click", () => {
  // Remove the JWT token from localStorage
  localStorage.removeItem("adminToken");

  // Redirect to the login page
  window.location.reload();
});

async function loadSidebarStats() {
  try {
    const adminToken = sessionStorage.getItem("adminToken");

    // Fetch transaction count
    const transactionCountResponse = await fetch(
      "/api/admin/transaction-count",
      {
        method: "GET",
        headers: { Authorization: adminToken },
      }
    );

    if (transactionCountResponse.ok) {
      const transactionCountData = await transactionCountResponse.json();
      document.getElementById("totalTransactionCount").textContent =
        transactionCountData.count;
    } else {
      console.error("Failed to fetch transaction count");
      document.getElementById("totalTransactionCount").textContent = "Error";
    }

    // Fetch total transaction amount
    const totalAmountResponse = await fetch(
      "/api/admin/total-transaction-amount",
      {
        method: "GET",
        headers: { Authorization: adminToken },
      }
    );

    if (totalAmountResponse.ok) {
      const totalAmountData = await totalAmountResponse.json();
      const totalAmount = totalAmountData.totalAmount.toFixed(2);
      document.getElementById(
        "totalTransactionAmount"
      ).textContent = `₹${totalAmount}`;

      // Calculate Total Points
      const totalPoints = (totalAmount / 25).toFixed(2);
      document.getElementById("totalTransactionPoints").textContent =
        totalPoints;
    } else {
      console.error("Failed to fetch total transaction amount");
      document.getElementById("totalTransactionAmount").textContent = "Error";
      document.getElementById("totalTransactionPoints").textContent = "Error";
    }
  } catch (error) {
    console.error("Error fetching sidebar stats:", error);
    document.getElementById("totalTransactionCount").textContent = "Error";
    document.getElementById("totalTransactionAmount").textContent = "Error";
    document.getElementById("totalTransactionPoints").textContent = "Error";
  }
}

// Load stats on page load
document.addEventListener("DOMContentLoaded", loadSidebarStats);

// Load stats on page load
document.addEventListener("DOMContentLoaded", loadSidebarStats);
