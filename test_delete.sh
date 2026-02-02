#!/bin/bash

# Test Product Delete Functionality with File Cleanup

BASE_URL="http://localhost:8004"

echo "==================================="
echo "Testing Product Delete Functionality"
echo "==================================="

# Step 1: List existing products
echo -e "\n1. Fetching existing products..."
PRODUCTS=$(curl -s "${BASE_URL}/products?limit=5")
echo "$PRODUCTS" | jq '.'

# Extract first product ID
PRODUCT_ID=$(echo "$PRODUCTS" | jq -r '.data[0].id // empty')

if [ -z "$PRODUCT_ID" ]; then
    echo "No products found to test delete"
    exit 1
fi

echo -e "\nTesting with Product ID: $PRODUCT_ID"

# Step 2: Get product details to see image paths
echo -e "\n2. Fetching product details..."
PRODUCT_DETAIL=$(curl -s "${BASE_URL}/products/${PRODUCT_ID}")
echo "$PRODUCT_DETAIL" | jq '.'

# Extract image paths
MAIN_IMAGE=$(echo "$PRODUCT_DETAIL" | jq -r '.data.image // empty')
ADDITIONAL_IMAGES=$(echo "$PRODUCT_DETAIL" | jq -r '.data.images[]?.path // empty')

echo -e "\n3. Checking if image files exist before delete..."
if [ -n "$MAIN_IMAGE" ]; then
    if [ -f "$MAIN_IMAGE" ]; then
        echo "✓ Main image exists: $MAIN_IMAGE"
    else
        echo "✗ Main image not found: $MAIN_IMAGE"
    fi
fi

while IFS= read -r img_path; do
    if [ -n "$img_path" ] && [ -f "$img_path" ]; then
        echo "✓ Additional image exists: $img_path"
    elif [ -n "$img_path" ]; then
        echo "✗ Additional image not found: $img_path"
    fi
done <<< "$ADDITIONAL_IMAGES"

# Step 4: Delete the product
echo -e "\n4. Deleting product ID: $PRODUCT_ID..."
DELETE_RESPONSE=$(curl -s -X DELETE "${BASE_URL}/products/${PRODUCT_ID}")
echo "$DELETE_RESPONSE" | jq '.'

# Step 5: Verify files are deleted
echo -e "\n5. Verifying image files are deleted..."
if [ -n "$MAIN_IMAGE" ]; then
    if [ -f "$MAIN_IMAGE" ]; then
        echo "✗ FAIL: Main image still exists: $MAIN_IMAGE"
    else
        echo "✓ PASS: Main image deleted: $MAIN_IMAGE"
    fi
fi

while IFS= read -r img_path; do
    if [ -n "$img_path" ]; then
        if [ -f "$img_path" ]; then
            echo "✗ FAIL: Additional image still exists: $img_path"
        else
            echo "✓ PASS: Additional image deleted: $img_path"
        fi
    fi
done <<< "$ADDITIONAL_IMAGES"

# Step 6: Verify product is deleted from database
echo -e "\n6. Verifying product is deleted from database..."
CHECK_RESPONSE=$(curl -s "${BASE_URL}/products/${PRODUCT_ID}")
ERROR_MSG=$(echo "$CHECK_RESPONSE" | jq -r '.error // empty')

if [ "$ERROR_MSG" = "Product not found" ]; then
    echo "✓ PASS: Product deleted from database"
else
    echo "✗ FAIL: Product still exists in database"
fi

echo -e "\n==================================="
echo "Test Complete"
echo "==================================="
