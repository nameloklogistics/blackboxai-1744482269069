// API endpoint for blockchain interaction
const API_BASE_URL = 'http://localhost:3000/api';

// Toast notification handler
function showToast(message, isError = false) {
    const toast = document.getElementById('statusToast');
    const statusMessage = document.getElementById('statusMessage');
    statusMessage.textContent = message;
    statusMessage.className = isError ? 'text-danger' : 'text-success';
    const bsToast = new bootstrap.Toast(toast);
    bsToast.show();
}

// Shipper Management
document.getElementById('shipperForm').addEventListener('submit', async (e) => {
    e.preventDefault();
    try {
        const response = await fetch(`${API_BASE_URL}/shipper`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                id: document.getElementById('shipperId').value,
                name: document.getElementById('shipperName').value
            })
        });
        
        if (!response.ok) throw new Error('Failed to create shipper');
        showToast('Shipper created successfully');
        e.target.reset();
    } catch (error) {
        showToast(error.message, true);
    }
});

// Freight Quotation
document.getElementById('quotationForm').addEventListener('submit', async (e) => {
    e.preventDefault();
    try {
        const response = await fetch(`${API_BASE_URL}/quotation`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                id: document.getElementById('quotationId').value,
                carrier: document.getElementById('carrier').value,
                amount: document.getElementById('amount').value
            })
        });
        
        if (!response.ok) throw new Error('Failed to create quotation');
        showToast('Quotation created successfully');
        e.target.reset();
    } catch (error) {
        showToast(error.message, true);
    }
});

// Shipment Booking
document.getElementById('shipmentForm').addEventListener('submit', async (e) => {
    e.preventDefault();
    try {
        const response = await fetch(`${API_BASE_URL}/shipment`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                id: document.getElementById('shipmentId').value,
                shipperId: document.getElementById('shipmentShipperId').value,
                consigneeId: document.getElementById('shipmentConsigneeId').value
            })
        });
        
        if (!response.ok) throw new Error('Failed to book shipment');
        showToast('Shipment booked successfully');
        e.target.reset();
    } catch (error) {
        showToast(error.message, true);
    }
});

// Shipment Tracking
document.getElementById('trackingForm').addEventListener('submit', async (e) => {
    e.preventDefault();
    try {
        const trackingId = document.getElementById('trackingId').value;
        const response = await fetch(`${API_BASE_URL}/shipment/${trackingId}`);
        
        if (!response.ok) throw new Error('Failed to retrieve shipment');
        
        const shipment = await response.json();
        const trackingResult = document.getElementById('trackingResult');
        trackingResult.innerHTML = `
            <div class="card">
                <div class="card-body">
                    <h6>Shipment Details</h6>
                    <p><strong>Status:</strong> ${shipment.status}</p>
                    <p><strong>Shipper ID:</strong> ${shipment.shipperId}</p>
                    <p><strong>Consignee ID:</strong> ${shipment.consigneeId}</p>
                    ${shipment.route ? `<p><strong>Route:</strong> ${shipment.route}</p>` : ''}
                    ${shipment.customsInfo ? `<p><strong>Customs Info:</strong> ${shipment.customsInfo}</p>` : ''}
                </div>
            </div>
        `;
        showToast('Shipment details retrieved successfully');
    } catch (error) {
        showToast(error.message, true);
        document.getElementById('trackingResult').innerHTML = '';
    }
});

// Initialize Bootstrap components
document.addEventListener('DOMContentLoaded', () => {
    // Initialize all tooltips
    const tooltipTriggerList = [].slice.call(document.querySelectorAll('[data-bs-toggle="tooltip"]'));
    tooltipTriggerList.map(function (tooltipTriggerEl) {
        return new bootstrap.Tooltip(tooltipTriggerEl);
    });
});

// WebSocket connection for real-time updates
let ws;
try {
    ws = new WebSocket('ws://localhost:3000/ws');
    
    ws.onopen = () => {
        console.log('WebSocket connection established');
    };
    
    ws.onmessage = (event) => {
        const update = JSON.parse(event.data);
        showToast(`Update: ${update.message}`);
    };
    
    ws.onerror = (error) => {
        console.error('WebSocket error:', error);
    };
    
    ws.onclose = () => {
        console.log('WebSocket connection closed');
        // Attempt to reconnect after 5 seconds
        setTimeout(() => {
            ws = new WebSocket('ws://localhost:3000/ws');
        }, 5000);
    };
} catch (error) {
    console.error('Failed to establish WebSocket connection:', error);
}
