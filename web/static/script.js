// Global variables
let packageSizes = [];

// Initialize the application
document.addEventListener('DOMContentLoaded', function() {
    loadPackageSizes();
    setupEventListeners();
});

// Setup event listeners
function setupEventListeners() {
    const quantityInput = document.getElementById('quantity');
    const calculateBtn = document.getElementById('calculate-btn');
    
    // Calculate on Enter key
    quantityInput.addEventListener('keypress', function(e) {
        if (e.key === 'Enter') {
            calculate();
        }
    });
    
    // Calculate on button click
    calculateBtn.addEventListener('click', calculate);
}

// Load package sizes from the server
async function loadPackageSizes() {
    try {
        const response = await fetch('/api/health');
        if (response.ok) {
            // For now, we'll use default package sizes
            // In a real implementation, you might have an endpoint to get package sizes
            packageSizes = [250, 500, 1000, 2000];
            displayPackageSizes();
        }
    } catch (error) {
        console.error('Failed to load package sizes:', error);
        // Use default package sizes
        packageSizes = [250, 500, 1000, 2000];
        displayPackageSizes();
    }
}

// Display package sizes in the UI
function displayPackageSizes() {
    const container = document.getElementById('package-sizes-display');
    container.innerHTML = '';
    
    packageSizes.forEach(size => {
        const sizeElement = document.createElement('div');
        sizeElement.className = 'package-size';
        sizeElement.textContent = size;
        container.appendChild(sizeElement);
    });
}

// Main calculation function
async function calculate() {
    const quantityInput = document.getElementById('quantity');
    const quantity = parseInt(quantityInput.value);
    
    // Validate input
    if (!quantityInput.value || isNaN(quantity) || quantity < 0) {
        showError('Please enter a valid positive number');
        return;
    }
    
    // Show loading state
    showLoading();
    hideResults();
    hideError();
    
    try {
        const response = await fetch(`/api/calculate?qty=${quantity}`);
        
        if (!response.ok) {
            const errorData = await response.json();
            throw new Error(errorData.error || 'Failed to calculate');
        }
        
        const result = await response.json();
        displayResults(result);
        
    } catch (error) {
        console.error('Calculation error:', error);
        showError(error.message || 'Failed to calculate optimal packages');
    } finally {
        hideLoading();
    }
}

// Display calculation results
function displayResults(result) {
    // Update result values
    document.getElementById('requested').textContent = result.requested.toLocaleString();
    document.getElementById('total-delivered').textContent = result.total_delivered.toLocaleString();
    document.getElementById('over-delivery').textContent = result.over_delivery.toLocaleString();
    
    // Display packages
    const packagesContainer = document.getElementById('packages-display');
    packagesContainer.innerHTML = '';
    
    if (Object.keys(result.packages).length === 0) {
        packagesContainer.innerHTML = '<p>No packages needed</p>';
    } else {
        Object.entries(result.packages).forEach(([size, count]) => {
            const packageElement = document.createElement('div');
            packageElement.className = 'package-item';
            packageElement.innerHTML = `
                <div class="size">${size}</div>
                <div class="count">${count} package${count > 1 ? 's' : ''}</div>
            `;
            packagesContainer.appendChild(packageElement);
        });
    }
    
    // Show results
    showResults();
}

// Show/hide functions
function showLoading() {
    document.getElementById('loading-section').style.display = 'block';
}

function hideLoading() {
    document.getElementById('loading-section').style.display = 'none';
}

function showResults() {
    document.getElementById('result-section').style.display = 'block';
}

function hideResults() {
    document.getElementById('result-section').style.display = 'none';
}

function showError(message) {
    const errorSection = document.getElementById('error-section');
    const errorMessage = document.getElementById('error-message');
    errorMessage.textContent = message;
    errorSection.style.display = 'block';
}

function hideError() {
    document.getElementById('error-section').style.display = 'none';
}

// Utility function to format numbers with commas
function formatNumber(num) {
    return num.toLocaleString();
} 