document.getElementById("counterTab").addEventListener("click", () => {
  showSection("counterSection", "counterTab");
  loadInitialMerchants();
});

let initialMerchants = []; // Cache the initial merchant list

async function loadInitialMerchants() {
  const merchantDropdown = document.getElementById("merchantDropdown");

  try {
    const adminToken = sessionStorage.getItem("adminToken");
    const response = await fetch(`/api/admin/users?page=1&limit=5`, {
      headers: { Authorization: adminToken },
    });

    const data = await response.json();

    if (data.success) {
      initialMerchants = data.users; // Cache the initial list
      populateMerchantDropdown(initialMerchants);
    } else {
      console.error("Failed to load merchants:", data.message);
      merchantDropdown.innerHTML =
        '<option value="">Error loading merchants</option>';
    }
  } catch (error) {
    console.error("Error loading merchants:", error);
    merchantDropdown.innerHTML =
      '<option value="">Error loading merchants</option>';
  }
}

function populateMerchantDropdown(merchants) {
  const merchantDropdown = document.getElementById("merchantDropdown");
  merchantDropdown.innerHTML = merchants.length
    ? merchants
        .map(
          (merchant) =>
            `<option value="${merchant.phone_number}">${merchant.store_name} (${merchant.phone_number}, ${merchant.location})</option>`
        )
        .join("")
    : '<option value="">No merchants found</option>';
}

function filterMerchantDropdown() {
  const searchInput = document
    .getElementById("searchMerchantDropdown")
    .value.toLowerCase();
  const filteredMerchants = initialMerchants.filter(
    (merchant) =>
      merchant.store_name.toLowerCase().includes(searchInput) ||
      merchant.phone_number.includes(searchInput) ||
      merchant.location.toLowerCase().includes(searchInput)
  );

  if (filteredMerchants.length > 0) {
    populateMerchantDropdown(filteredMerchants);
  } else {
    // Fallback to API search if no local matches
    searchMerchantsAPI(searchInput);
  }
}

async function searchMerchantsAPI(query) {
  if (query.length < 3) return; // Avoid making API calls for short queries

  const merchantDropdown = document.getElementById("merchantDropdown");
  merchantDropdown.innerHTML = '<option value="">Searching...</option>';

  try {
    const adminToken = sessionStorage.getItem("adminToken");
    const response = await fetch(
      `/api/admin/merchants/search?query=${encodeURIComponent(query)}`,
      {
        headers: { Authorization: adminToken },
      }
    );

    const data = await response.json();

    if (data.merchants && data.merchants.length > 0) {
      populateMerchantDropdown(data.merchants);
    } else {
      merchantDropdown.innerHTML =
        '<option value="">No merchants found</option>';
    }
  } catch (error) {
    console.error("Error searching merchants:", error);
    merchantDropdown.innerHTML =
      '<option value="">Error searching merchants</option>';
  }
}

// Load counters for the selected merchant
async function loadCounters(merchantPhone) {
  const counterTableBody = document.getElementById("counterTableBody");
  const counterList = document.getElementById("counterList");

  if (!merchantPhone) {
    counterTableBody.innerHTML = "";
    counterList.style.display = "none";
    return;
  }

  try {
    const adminToken = sessionStorage.getItem("adminToken");
    const response = await fetch(
      `/api/admin/counters?merchant=${merchantPhone}`,
      {
        headers: { Authorization: adminToken },
      }
    );

    const data = await response.json();
    if (data.success) {
      counterTableBody.innerHTML = data.counters
        .map(
          (counter, index) => `
          <tr>
            <td>${index + 1}</td>
            <td>${counter.Name}</td>
            <td>${counter.Location}</td>
            <td>${counter.Username}</td>
            <td>${counter.Description}</td>
          </tr>
        `
        )
        .join("");
      counterList.style.display = "block";
    } else {
      counterTableBody.innerHTML =
        '<tr><td colspan="4">No counters found</td></tr>';
      counterList.style.display = "block";
    }
  } catch (error) {
    counterTableBody.innerHTML =
      '<tr><td colspan="4">No counters found.</td></tr>';
    counterList.style.display = "block";
  }
}

// Handle merchant selection
document
  .getElementById("merchantDropdown")
  .addEventListener("change", (event) => {
    const merchantPhone = event.target.value;
    document.getElementById("addCounterForm").style.display = merchantPhone
      ? "block"
      : "none";
    loadCounters(merchantPhone);
  });

// Handle counter addition
document
  .getElementById("addCounterForm")
  .addEventListener("submit", async (event) => {
    event.preventDefault();
    const formData = new FormData(event.target);
    const counterData = {
      merchantPhone: document.getElementById("merchantDropdown").value,
      name: formData.get("counterName"),
      location: formData.get("counterLocation"),
      description: formData.get("counterDescription"),
      username: formData.get("username"),
      password: formData.get("password"),
    };

    try {
      const adminToken = sessionStorage.getItem("adminToken");
      const response = await fetch("/api/admin/counters", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: adminToken,
        },
        body: JSON.stringify(counterData),
      });

      const data = await response.json();
      if (data.success) {
        alert("Counter added successfully!");
        event.target.reset();
        loadCounters(counterData.merchantPhone);
      } else {
        alert("Failed to add counter: " + data.message);
      }
    } catch (error) {
      console.error("Error adding counter:", error);
      alert("Something went wrong. Please try again.");
    }
  });

document
  .getElementById("merchantDropdown")
  .addEventListener("change", (event) => {
    const selectedOption = event.target.options[event.target.selectedIndex];
    const selectedMerchantDetails = document.getElementById(
      "selectedMerchantDetails"
    );

    selectedMerchantDetails.innerHTML = `
      <span>${selectedOption.textContent}</span>
    `;

    // Hide the dropdown and show the selected merchant
    document.getElementById("selectMerchantSection").style.display = "none";
    document.getElementById("selectedMerchantSection").style.display = "block";
  });

// Change Merchant Button
document
  .getElementById("changeMerchantButton")
  .addEventListener("click", () => {
    // Reset UI for selecting merchant
    document.getElementById("selectMerchantSection").style.display = "block";
    document.getElementById("selectedMerchantSection").style.display = "none";
  });
