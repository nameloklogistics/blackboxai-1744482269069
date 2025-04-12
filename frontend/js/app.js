// API Base URL
const API_BASE_URL = 'http://localhost:3000/api';

// Wallet Connection Simulation
document.getElementById('connectWallet').addEventListener('click', async () => {
    const button = document.getElementById('connectWallet');
    button.disabled = true;
    button.innerHTML = '<i class="fas fa-spinner fa-spin mr-2"></i>Connecting...';
    
    // Simulate connection delay
    await new Promise(resolve => setTimeout(resolve, 1000));
    
    button.innerHTML = '<i class="fas fa-check-circle mr-2"></i>Connected';
    button.classList.remove('bg-blue-600', 'hover:bg-blue-700');
    button.classList.add('bg-green-600');
});

// Fetch and Display Service Categories
async function fetchServiceCategories() {
    try {
        const response = await fetch(`${API_BASE_URL}/service-categories`);
        const categories = await response.json();
        displayServices(categories);
    } catch (error) {
        console.error('Error fetching service categories:', error);
        showError('Failed to load service categories');
    }
}

// Display Services in the Service List
function displayServices(categories) {
    const serviceList = document.getElementById('serviceList');
    serviceList.innerHTML = '';

    categories.forEach(category => {
        category.subcategories.forEach(subcategory => {
            subcategory.items.forEach(item => {
                const serviceCard = createServiceCard(category, subcategory, item);
                serviceList.appendChild(serviceCard);
            });
        });
    });
}

// Create Service Card Element
function createServiceCard(category, subcategory, item) {
    const card = document.createElement('div');
    card.className = 'bg-gray-50 rounded-lg p-6 hover:shadow-md transition-shadow';
    
    card.innerHTML = `
        <div class="flex justify-between items-start mb-4">
            <div>
                <h3 class="text-lg font-semibold text-gray-900">${item.name}</h3>
                <p class="text-sm text-gray-600">${category.name} > ${subcategory.name}</p>
            </div>
            <span class="bg-blue-100 text-blue-800 text-sm font-medium px-3 py-1 rounded-full">
                ${item.rate} ${item.currency}
            </span>
        </div>
        <p class="text-gray-700 mb-4">${item.description}</p>
        <button onclick="bookService('${item.id}')" 
                class="bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700 transition-colors w-full">
            Book Service
        </button>
    `;
    
    return card;
}

// Book Service Handler
async function bookService(serviceId) {
    if (!isWalletConnected()) {
        showError('Please connect your wallet first');
        return;
    }

    try {
        // Here you would typically make an API call to book the service
        showSuccess('Service booked successfully!');
    } catch (error) {
        console.error('Error booking service:', error);
        showError('Failed to book service');
    }
}

// Service Provider Registration
document.getElementById('providerForm').addEventListener('submit', async (e) => {
    e.preventDefault();
    
    if (!isWalletConnected()) {
        showError('Please connect your wallet first');
        return;
    }

    const formData = new FormData(e.target);
    const providerData = {
        id: 'SP' + Date.now(), // Generate a simple ID
        companyName: formData.get('companyName'),
        contactPerson: formData.get('contactPerson'),
        telephone: formData.get('telephone'),
        country: formData.get('country')
    };

    try {
        const response = await fetch(`${API_BASE_URL}/service-seller`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(providerData)
        });

        if (response.ok) {
            showSuccess('Registration successful!');
            e.target.reset();
        } else {
            throw new Error('Registration failed');
        }
    } catch (error) {
        console.error('Error registering provider:', error);
        showError('Registration failed');
    }
});

// Utility Functions
function isWalletConnected() {
    const button = document.getElementById('connectWallet');
    return button.classList.contains('bg-green-600');
}

function showError(message) {
    // Create and show error notification
    const notification = document.createElement('div');
    notification.className = 'fixed top-4 right-4 bg-red-500 text-white px-6 py-3 rounded-lg shadow-lg';
    notification.textContent = message;
    document.body.appendChild(notification);
    setTimeout(() => notification.remove(), 3000);
}

function showSuccess(message) {
    // Create and show success notification
    const notification = document.createElement('div');
    notification.className = 'fixed top-4 right-4 bg-green-500 text-white px-6 py-3 rounded-lg shadow-lg';
    notification.textContent = message;
    document.body.appendChild(notification);
    setTimeout(() => notification.remove(), 3000);
}

// Initialize
document.addEventListener('DOMContentLoaded', () => {
    fetchServiceCategories();
});

// WebSocket Connection for Real-time Updates
const ws = new WebSocket('ws://localhost:8080');

ws.onmessage = (event) => {
    const update = JSON.parse(event.data);
    // Handle different types of updates
    if (update.type === 'new_service') {
        fetchServiceCategories(); // Refresh the service list
    }
};

ws.onerror = (error) => {
    console.error('WebSocket error:', error);
};

ws.onclose = () => {
    console.log('WebSocket connection closed');
    // Implement reconnection logic if needed
};
