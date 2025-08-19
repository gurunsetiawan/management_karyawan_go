// assets/js/app.js
class EmployeeManager {
  constructor() {
    this.employees = [];
    this.filteredEmployees = [];
    this.isLoading = false;
    this.init();
  }

  async init() {
    await this.loadEmployees();
    this.setupEventListeners();
  }

  showLoading() {
    this.isLoading = true;
    const tableBody = document.getElementById("employee-table-body");
    tableBody.innerHTML = `
      <tr>
        <td colspan="7" class="has-text-centered">
          <div class="loading-spinner">
            <span class="icon is-large">
              <i class="fas fa-spinner fa-pulse"></i>
            </span>
            <p class="mt-2">Memuat data karyawan...</p>
          </div>
        </td>
      </tr>
    `;
  }

  hideLoading() {
    this.isLoading = false;
  }

  showError(message) {
    const tableBody = document.getElementById("employee-table-body");
    tableBody.innerHTML = `
      <tr>
        <td colspan="7" class="has-text-centered has-text-danger">
          <span class="icon is-large">
            <i class="fas fa-exclamation-triangle"></i>
          </span>
          <p class="mt-2">${message}</p>
          <button class="button is-small is-primary mt-2" onclick="location.reload()">
            Coba Lagi
          </button>
        </td>
      </tr>
    `;
  }

  async loadEmployees() {
    try {
      this.showLoading();
      const res = await fetch("/api/employees");
      
      if (!res.ok) {
        const errorData = await res.json();
        throw new Error(errorData.message || 'Gagal memuat data karyawan');
      }
      
      this.employees = await res.json();
      this.filteredEmployees = [...this.employees];
      this.updateEmployeeCount();
      this.renderEmployees();
    } catch (err) {
      console.error('Error loading employees:', err);
      this.showError(err.message || 'Gagal memuat data karyawan');
    } finally {
      this.hideLoading();
    }
  }

  updateEmployeeCount() {
    const totalCount = document.getElementById("total-count");
    if (totalCount) {
      totalCount.textContent = this.filteredEmployees.length;
    }
  }

  renderEmployees() {
    const tableBody = document.getElementById("employee-table-body");
    
    if (this.filteredEmployees.length === 0) {
      tableBody.innerHTML = `
        <tr>
          <td colspan="7" class="has-text-centered has-text-grey">
            <span class="icon is-large">
              <i class="fas fa-users"></i>
            </span>
            <p class="mt-2">Tidak ada data karyawan</p>
          </td>
        </tr>
      `;
      return;
    }

    tableBody.innerHTML = this.filteredEmployees
      .map(
        (emp, index) => `
        <tr class="fade-in" style="animation-delay: ${index * 0.1}s">
          <td><strong>${emp.id}</strong></td>
          <td>
            <div class="has-text-weight-semibold">${this.escapeHtml(emp.name)}</div>
          </td>
          <td>
            <a href="mailto:${this.escapeHtml(emp.email)}" class="has-text-info">
              ${this.escapeHtml(emp.email)}
            </a>
          </td>
          <td><span class="tag is-info bounce-in">${this.escapeHtml(emp.role)}</span></td>
          <td>
            <a href="tel:${this.escapeHtml(emp.phone)}" class="has-text-success">
              ${this.escapeHtml(emp.phone)}
            </a>
          </td>
          <td>
            <div class="has-text-grey-dark">
              <i class="fas fa-map-marker-alt mr-2"></i>
              ${this.escapeHtml(emp.alamat)}
            </div>
          </td>
          <td>
            <div class="buttons are-small">
              <button class="button is-warning slide-in-right edit-btn" 
                      data-id="${emp.id}"
                      style="animation-delay: ${index * 0.1 + 0.2}s">
                <span class="icon is-small">
                  <i class="fas fa-edit"></i>
                </span>
                <span>Edit</span>
              </button>
              <button class="button is-danger slide-in-right delete-btn" 
                      data-id="${emp.id}"
                      style="animation-delay: ${index * 0.1 + 0.3}s">
                <span class="icon is-small">
                  <i class="fas fa-trash"></i>
                </span>
                <span>Hapus</span>
              </button>
            </div>
          </td>
        </tr>
      `,
      )
      .join("");

    this.setupActionButtons();
  }

  escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
  }

  setupActionButtons() {
    // Delete buttons
    document.querySelectorAll(".delete-btn").forEach((btn) => {
      btn.addEventListener("click", async (e) => {
        e.preventDefault();
        const id = btn.getAttribute("data-id");
        const employee = this.employees.find(emp => emp.id == id);
        
        if (confirm(`Yakin ingin menghapus karyawan "${employee?.name}"?`)) {
          await this.deleteEmployee(id);
        }
      });
    });

    // Edit buttons
    document.querySelectorAll(".edit-btn").forEach((btn) => {
      btn.addEventListener("click", (e) => {
        e.preventDefault();
        const id = btn.getAttribute("data-id");
        window.location.href = `/edit.html?id=${id}`;
      });
    });
  }

  async deleteEmployee(id) {
    try {
      const res = await fetch(`/api/employees/${id}`, { 
        method: "DELETE" 
      });
      
      if (!res.ok) {
        const errorData = await res.json();
        throw new Error(errorData.message || 'Gagal menghapus karyawan');
      }
      
      // Remove from local arrays
      this.employees = this.employees.filter(emp => emp.id != id);
      this.filteredEmployees = this.filteredEmployees.filter(emp => emp.id != id);
      
      this.updateEmployeeCount();
      this.renderEmployees();
      
      // Show success message
      this.showNotification('Karyawan berhasil dihapus!', 'success');
    } catch (err) {
      console.error('Error deleting employee:', err);
      this.showNotification(err.message || 'Gagal menghapus karyawan', 'danger');
    }
  }

  setupEventListeners() {
    // Search functionality
    const searchInput = document.getElementById("search-input");
    if (searchInput) {
      searchInput.addEventListener("input", (e) => {
        this.filterEmployees(e.target.value);
      });
    }

    // Filter by role
    const roleFilter = document.getElementById("role-filter");
    if (roleFilter) {
      roleFilter.addEventListener("change", (e) => {
        this.filterByRole(e.target.value);
      });
    }
  }

  filterEmployees(searchTerm) {
    if (!searchTerm.trim()) {
      this.filteredEmployees = [...this.employees];
    } else {
      const term = searchTerm.toLowerCase();
      this.filteredEmployees = this.employees.filter(emp => 
        emp.name.toLowerCase().includes(term) ||
        emp.email.toLowerCase().includes(term) ||
        emp.role.toLowerCase().includes(term) ||
        emp.phone.includes(term)
      );
    }
    
    this.updateEmployeeCount();
    this.renderEmployees();
  }

  filterByRole(role) {
    if (!role || role === "all") {
      this.filteredEmployees = [...this.employees];
    } else {
      this.filteredEmployees = this.employees.filter(emp => emp.role === role);
    }
    
    this.updateEmployeeCount();
    this.renderEmployees();
  }

  showNotification(message, type = 'info') {
    // Create notification element
    const notification = document.createElement('div');
    notification.className = `notification is-${type} is-light`;
    notification.innerHTML = `
      <button class="delete" onclick="this.parentElement.remove()"></button>
      ${message}
    `;
    
    // Add to page
    const container = document.querySelector('.container');
    if (container) {
      container.insertBefore(notification, container.firstChild);
      
      // Auto remove after 5 seconds
      setTimeout(() => {
        if (notification.parentElement) {
          notification.remove();
        }
      }, 5000);
    }
  }
}

// Initialize when DOM is loaded
document.addEventListener("DOMContentLoaded", () => {
  new EmployeeManager();
});
