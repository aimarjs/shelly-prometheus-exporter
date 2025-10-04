#!/bin/bash

# Test script for Shelly Pro3em power readings
# Usage: ./test-shelly-pro3em.sh [IP_ADDRESS]

SHELLY_IP=${1:-"10.69.20.19"}

echo "Testing Shelly Pro3em at $SHELLY_IP"
echo "=================================="

# Test RPC endpoint
echo "1. Testing RPC endpoint (/rpc/Shelly.GetStatus):"
echo "-----------------------------------------------"
RPC_RESPONSE=$(curl -s --connect-timeout 5 "http://$SHELLY_IP/rpc/Shelly.GetStatus")
if [ $? -eq 0 ] && [ -n "$RPC_RESPONSE" ]; then
    echo "✓ RPC endpoint accessible"
    echo "Power readings from RPC:"
    echo "$RPC_RESPONSE" | jq -r '
        if .em then
            "Phase A Power: " + (.em."a_act_power" | tostring) + " W",
            "Phase B Power: " + (.em."b_act_power" | tostring) + " W", 
            "Phase C Power: " + (.em."c_act_power" | tostring) + " W",
            "Total Power: " + (.em."total_act_power" | tostring) + " W"
        else
            "No EM data found in RPC response"
        end
    '
else
    echo "✗ RPC endpoint not accessible or empty response"
fi

echo ""

# Test legacy endpoint
echo "2. Testing legacy endpoint (/status):"
echo "------------------------------------"
LEGACY_RESPONSE=$(curl -s --connect-timeout 5 "http://$SHELLY_IP/status")
if [ $? -eq 0 ] && [ -n "$LEGACY_RESPONSE" ]; then
    echo "✓ Legacy endpoint accessible"
    echo "Power readings from legacy:"
    echo "$LEGACY_RESPONSE" | jq -r '
        if .meters and (.meters | length) > 0 then
            "Power: " + (.meters[0].power | tostring) + " W",
            "Total Energy: " + (.meters[0].total | tostring) + " Wh"
        else
            "No meter data found in legacy response"
        end
    '
else
    echo "✗ Legacy endpoint not accessible or empty response"
fi

echo ""

# Test exporter metrics
echo "3. Testing exporter metrics (if running):"
echo "-----------------------------------------"
EXPORTER_RESPONSE=$(curl -s --connect-timeout 5 "http://localhost:8080/metrics")
if [ $? -eq 0 ] && [ -n "$EXPORTER_RESPONSE" ]; then
    echo "✓ Exporter accessible"
    echo "Power metrics from exporter:"
    echo "$EXPORTER_RESPONSE" | grep "shelly_power_watts" | head -10
else
    echo "✗ Exporter not accessible (make sure it's running on localhost:8080)"
fi

echo ""

# Analysis
echo "4. Analysis:"
echo "-----------"
echo "If you see inflated power readings, possible causes:"
echo "1. Exporter is summing individual phases instead of using total power"
echo "2. Unit conversion issues (W vs kW)"
echo "3. Multiple devices being monitored"
echo "4. Incorrect parsing of Shelly API response"
echo ""
echo "Expected behavior:"
echo "- Shelly Pro3em should report total power consumption"
echo "- Individual phase readings are for monitoring balance"
echo "- Total power should equal sum of all phases"
