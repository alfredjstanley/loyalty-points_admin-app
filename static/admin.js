// Switch tabs
document.getElementById("addUserTab").addEventListener("click", () => {
  showSection("addUserSection", "addUserTab");
});
document.getElementById("listUsersTab").addEventListener("click", () => {
  showSection("listUsersSection", "listUsersTab");
  loadUsers(); // Load users when navigating to the Merchant List
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
    await fetch("/api/admin/onboard", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(user),
    });
    alert("Merchant added successfully!");
    e.target.reset();
  } catch (error) {
    alert("Failed to add merchant.");
  }
});

// Load Merchants with Pagination
async function loadUsers(page = 1) {
  try {
    const res = await fetch(`/api/admin/users?page=${page}&limit=10`);
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
          <td>
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
    await fetch("/api/admin/edit-merchant", {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ id, store_name, location, password }),
    });
    alert("Merchant updated successfully!");
    document.getElementById("editModal").style.display = "none";
    loadUsers();
  } catch (error) {
    alert("Failed to update merchant.");
  }
});

document.addEventListener("click", async (e) => {
  if (e.target.classList.contains("view-btn")) {
    const mobileNumber = e.target.dataset.id; // Get the merchant ID
    try {
      // Fetch merchant details from the API
      const res = await fetch(
        `/api/admin/edit-merchant?mobileNumber=${mobileNumber}`,
        {
          method: "GET",
          headers: { "Content-Type": "application/json" },
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

// Close the modal
document
  .querySelector("#viewModal .close-btn")
  .addEventListener("click", () => {
    document.getElementById("viewModal").style.display = "none";
  });

// Handle sidebar navigation
document.getElementById("transactionLogsTab").addEventListener("click", () => {
  showSection("transactionLogsSection", "transactionLogsTab");
  loadTransactionLogs(); // Load logs when section is clicked
});

// Load transaction logs
async function loadTransactionLogs(page = 1) {
  try {
    const res = await fetch(
      `/api/admin/transaction-logs?page=${page}&limit=10`
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
