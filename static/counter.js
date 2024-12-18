document.getElementById("counterTab").addEventListener("click", () => {
  showSection("counterSection", "counterTab");
  loadMerchantsForDropdown();
});

async function loadMerchantsForDropdown() {
  const merchantDropdown = document.getElementById("merchantDropdown");
  try {
    const adminToken = sessionStorage.getItem("adminToken");
    const response = await fetch("/api/admin/users", {
      headers: { Authorization: adminToken },
    });

    const data = await response.json();

    merchantDropdown.innerHTML = '<option value="">Select Merchant</option>';
    data.users.forEach((merchant) => {
      const option = document.createElement("option");
      option.value = merchant.phone_number;
      option.textContent = `${merchant.store_name} (${merchant.phone_number})`;
      merchantDropdown.appendChild(option);
    });
  } catch (error) {
    console.error("Error loading merchants:", error);
    merchantDropdown.innerHTML =
      '<option value="">Error loading merchants</option>';
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
    console.log(`data.counters knri `, data.counters);
    if (data.success) {
      counterTableBody.innerHTML = data.counters
        .map(
          (counter, index) => `
          <tr>
            <td>${index + 1}</td>
            <td>${counter.Name}</td>
            <td>${counter.Location}</td>
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
