// Show loader
function showLoader() {
  document.getElementById("loader").style.display = "block";
}

// Hide loader
function hideLoader() {
  document.getElementById("loader").style.display = "none";
}

// Show modal with message
function showModal(message) {
  const modal = document.getElementById("modal");
  const modalMessage = document.getElementById("modalMessage");
  modalMessage.textContent = message;
  modal.style.display = "flex";

  // Close modal on click
  document.getElementById("closeModal").onclick = function () {
    modal.style.display = "none";
  };

  window.onclick = function (event) {
    if (event.target === modal) {
      modal.style.display = "none";
    }
  };
}

// Handle login form submission
async function handleLoginSubmit(event) {
  event.preventDefault();

  showLoader();

  const form = event.target;
  const formData = new FormData(form);
  const data = {
    mobile_number: formData.get("mobileNumber"),
    password: formData.get("password"),
  };

  try {
    const response = await fetch("/api/login", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(data),
    });

    hideLoader();

    if (response.ok) {
      const result = await response.json();
      sessionStorage.setItem("jwtToken", result.token);

      sessionStorage.setItem("storeName", result.store_name);
      sessionStorage.setItem("location", result.location);
      sessionStorage.setItem("phoneNumber", result.phone_number);

      showModal("Login successful! Redirecting...");
      window.location.replace("/form");
    } else {
      const error = await response.json();
      showModal(`Error: ${error.message}`);
    }
  } catch (error) {
    hideLoader();
    showModal("An error occurred during login.");
  }
}

// Attach event listener
document
  .getElementById("loginForm")
  .addEventListener("submit", handleLoginSubmit);
