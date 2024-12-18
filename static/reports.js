document.getElementById("reportsTab").addEventListener("click", () => {
  showSection("reportsSection", "reportsTab");
  loadReportData(); // Load report data when the section is active
});

// Load Reports
async function loadReportData() {
  try {
    const adminToken = sessionStorage.getItem("adminToken");
    const res = await fetch("/api/admin/reports", {
      method: "GET",
      headers: { Authorization: adminToken },
    });

    const response = await res.json();
    const tbody = document.getElementById("reportTableBody");
    tbody.innerHTML = "";

    response.reports.sort(
      (a, b) => b.total_transactions - a.total_transactions
    );
    response.reports.forEach((report, index) => {
      tbody.innerHTML += `
        <tr>
          <td>${index + 1}</td>
          <td>${report.store_name}</td>
          <td>${report.location}</td>
          <td>${report.phone_number}</td>

          <td>₹${report.total_sales.toFixed(2)}</td>
          <td>${report.total_transactions}</td>
          <td>${report.points_earned}</td>
        </tr>
      `;
    });
  } catch (error) {
    console.error("Failed to load reports:", error);
    alert("Failed to load reports.");
  }
}
