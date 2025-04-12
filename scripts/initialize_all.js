const path = require('path');

async function initializeAll() {
    const scripts = [
        'initialize_categories.js',
        'initialize_shipping_modes.js',
        'initialize_packaging_modes.js',
        'initialize_container_types.js',
        'initialize_freight_rates.js',
        'initialize_transport_terminals.js',
        'initialize_customs_services.js',
        'initialize_local_charges.js',
        'initialize_payment_methods.js',
        'initialize_token_utilities.js'
    ];

    for (const scriptName of scripts) {
        console.log(`Executing ${scriptName}...`);
        try {
            const scriptPath = path.join(__dirname, scriptName);
            require(scriptPath);
            await new Promise(resolve => setTimeout(resolve, 1000));
        } catch (error) {
            console.error(`Error executing ${scriptName}:`, error);
        }
    }
}

initializeAll().then(() => {
    console.log('Initialization complete');
}).catch(error => {
    console.error('Initialization failed:', error);
    process.exit(1);
});