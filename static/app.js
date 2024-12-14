// Function to get URL query parameters

function getQueryParams() {
  const params = {};
  const queryString = window.location.search;
  const urlParams = new URLSearchParams(queryString);
  urlParams.forEach((value, key) => {
    params[key] = value;
  });
  return params;
}

// Populate Merchant Phone from URL
function populateMerchantPhone() {
  const params = getQueryParams();
  if (params.merchantPhone) {
    const merchantPhoneInput = document.getElementById("merchantPhone");
    merchantPhoneInput.value = params.merchantPhone;
  }
  const merchantPhoneInput = document.getElementById("merchantPhone");
  const phoneNumber = sessionStorage.getItem("phoneNumber");
  merchantPhoneInput.value = phoneNumber;

  const storeName = sessionStorage.getItem("storeName");
  document.getElementById("storeTitle").innerText = storeName;
}

// Show the loader
function showLoader() {
  document.getElementById("loader").style.display = "block";
}

// Hide the loader
function hideLoader() {
  document.getElementById("loader").style.display = "none";
}

// Show modal with a message
function showModal(message) {
  const modal = document.getElementById("modal");
  const modalMessage = document.getElementById("modalMessage");
  modalMessage.textContent = message;
  modal.style.display = "flex";

  // Close modal on clicking the close button
  document.getElementById("closeModal").onclick = function () {
    modal.style.display = "none";
  };

  // Close modal on clicking outside the modal content
  window.onclick = function (event) {
    if (event.target === modal) {
      modal.style.display = "none";
    }
  };
}

// Handle form submission
async function handleFormSubmit(event) {
  event.preventDefault();

  showLoader();

  const form = event.target;
  const formData = new FormData(form);
  // check user and merchant phone numbers are valid phone numbers
  const phoneRegex = /^\d{10}$/;
  if (!phoneRegex.test(formData.get("userPhone"))) {
    hideLoader();
    showModal(
      "Invalid user phone number. Please enter a valid 10 digit number."
    );
    return;
  }
  if (!phoneRegex.test(formData.get("merchantPhone"))) {
    hideLoader();
    showModal(
      "Invalid merchant phone number. Please enter a valid 10 digit number."
    );
    return;
  }
  const data = {
    user_mobile_number: formData.get("userPhone"),
    merchant_mobile_number: formData.get("merchantPhone"),
    amount: parseFloat(formData.get("amount")),
    invoice_id: formData.get("invoiceId"),
    payment_mode: formData.get("paymentMode"),
  };

  try {
    const token = sessionStorage.getItem("jwtToken");
    if (!token) {
      hideLoader();
      alert("You are not logged in. Please login to continue.");
      window.location.replace("/");
    }

    const response = await fetch("/api/points", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify(data),
    });

    hideLoader();

    if (response.ok) {
      showModal("Points transferred successfully!");
      form.reset();
      populateMerchantPhone();
    }
    if (response.status == 409) {
      showModal(
        "Duplicate entry: a successful transaction with this invoice ID already exists"
      );
    }
  } catch (error) {
    hideLoader();
    showModal("Invalid request. Please check the details and try again.");
  }
}

// Attach event listeners on page load
window.onload = () => {
  populateMerchantPhone();
  document
    .getElementById("paymentForm")
    .addEventListener("submit", handleFormSubmit);
};
