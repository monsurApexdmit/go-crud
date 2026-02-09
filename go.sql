-- phpMyAdmin SQL Dump
-- version 5.2.3
-- https://www.phpmyadmin.net/
--
-- Host: central_mysql
-- Generation Time: Feb 09, 2026 at 09:02 AM
-- Server version: 9.5.0
-- PHP Version: 8.3.26

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
START TRANSACTION;
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- Database: `go`
--

-- --------------------------------------------------------

--
-- Table structure for table `attributes`
--

CREATE TABLE `attributes` (
  `id` bigint UNSIGNED NOT NULL,
  `name` varchar(255) NOT NULL,
  `display_name` varchar(255) DEFAULT NULL,
  `option_type` varchar(50) DEFAULT NULL,
  `values` text,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

--
-- Dumping data for table `attributes`
--

INSERT INTO `attributes` (`id`, `name`, `display_name`, `option_type`, `values`, `created_at`, `updated_at`, `deleted_at`) VALUES
(1, 'color', 'Color', 'dropdown', 'red,green,blue', '2026-01-28 09:06:52', '2026-01-28 09:06:52', NULL),
(2, 'size', 'Size', 'dropdown', 'S,M,L,XL', '2026-01-28 09:07:53', '2026-01-28 09:07:53', NULL),
(3, 'update', 'Size', 'dropdown', 'S,M,L,XL', '2026-01-28 09:08:14', '2026-01-28 09:09:18', NULL);

-- --------------------------------------------------------

--
-- Table structure for table `authors`
--

CREATE TABLE `authors` (
  `id` bigint UNSIGNED NOT NULL,
  `name` varchar(150) NOT NULL,
  `email` varchar(150) NOT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- --------------------------------------------------------

--
-- Table structure for table `books`
--

CREATE TABLE `books` (
  `id` bigint UNSIGNED NOT NULL,
  `title` varchar(255) NOT NULL,
  `author_id` bigint UNSIGNED NOT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- --------------------------------------------------------

--
-- Table structure for table `categories`
--

CREATE TABLE `categories` (
  `id` bigint UNSIGNED NOT NULL,
  `category_name` varchar(255) NOT NULL,
  `parent_id` bigint UNSIGNED DEFAULT NULL,
  `status` tinyint(1) DEFAULT '0',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

--
-- Dumping data for table `categories`
--

INSERT INTO `categories` (`id`, `category_name`, `parent_id`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES
(1, 'Electronics udpate', NULL, 0, '2026-01-27 12:56:10', '2026-01-27 12:57:53', NULL),
(2, 'Mobile Phones', 1, 0, '2026-01-27 12:56:56', '2026-01-27 12:56:56', NULL),
(3, 'test', 1, 0, '2026-01-27 12:58:23', '2026-01-27 12:59:06', '2026-01-27 12:59:07');

-- --------------------------------------------------------

--
-- Table structure for table `coupons`
--

CREATE TABLE `coupons` (
  `id` int NOT NULL,
  `campaign_name` varchar(255) NOT NULL,
  `code` varchar(50) NOT NULL,
  `discount` decimal(10,2) NOT NULL,
  `type` enum('percentage','fixed') NOT NULL,
  `start_date` datetime NOT NULL,
  `end_date` datetime NOT NULL,
  `status` tinyint(1) DEFAULT '0',
  `image` varchar(255) DEFAULT NULL,
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

--
-- Dumping data for table `coupons`
--

INSERT INTO `coupons` (`id`, `campaign_name`, `code`, `discount`, `type`, `start_date`, `end_date`, `status`, `image`, `created_at`, `updated_at`, `deleted_at`) VALUES
(1, 'New Year', 'sf23', 10.00, 'percentage', '2026-01-01 00:00:00', '2026-01-31 23:59:59', 0, 'uploads/coupons/1769684357490002994.jpg', '2026-01-29 08:57:14', '2026-01-29 11:15:49', NULL),
(3, 'New Year656565', '6565', 10.00, 'percentage', '2026-01-01 00:00:00', '2026-01-31 23:59:59', 1, 'uploads/coupons/1769679873575066600.jpg', '2026-01-29 09:44:34', '2026-01-29 10:12:27', '2026-01-29 10:12:27');

-- --------------------------------------------------------

--
-- Table structure for table `customers`
--

CREATE TABLE `customers` (
  `id` bigint UNSIGNED NOT NULL,
  `user_id` bigint UNSIGNED DEFAULT NULL,
  `name` varchar(150) NOT NULL,
  `email` varchar(150) NOT NULL,
  `phone` varchar(20) DEFAULT NULL,
  `address` varchar(255) DEFAULT NULL,
  `city` varchar(100) DEFAULT NULL,
  `state` varchar(100) DEFAULT NULL,
  `zip_code` varchar(20) DEFAULT NULL,
  `country` varchar(100) DEFAULT NULL,
  `customer_type` enum('retail','wholesale') DEFAULT 'retail',
  `status` enum('active','inactive') DEFAULT 'active',
  `notes` text,
  `store_credit` decimal(10,2) DEFAULT '0.00',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

--
-- Dumping data for table `customers`
--

INSERT INTO `customers` (`id`, `user_id`, `name`, `email`, `phone`, `address`, `city`, `state`, `zip_code`, `country`, `customer_type`, `status`, `notes`, `store_credit`, `created_at`, `updated_at`, `deleted_at`) VALUES
(2, 4, 'John Smith', 'john@example.com', '555-1234', '', '', '', '', '', 'retail', 'active', '', 50.00, '2026-02-05 07:24:10', '2026-02-05 07:25:24', '2026-02-05 07:25:25'),
(3, 8, 'customer', 'customer@example.com', '555-1234', '123 Main St', 'New York', 'NY', '10001', 'USA', 'retail', 'active', 'VIP customer', 100.00, '2026-02-05 09:01:49', '2026-02-05 09:01:49', NULL),
(4, 9, 'customer2', 'customer2@example.com', '555-1234', '123 Main St', 'New York', 'NY', '10001', 'USA', 'retail', 'active', 'VIP customer', 100.00, '2026-02-05 09:02:30', '2026-02-05 09:02:30', NULL),
(5, 10, 'customer3', 'customer3@example.com', '555-1234', '123 Main St', 'New York', 'NY', '10001', 'USA', 'wholesale', 'active', 'VIP customer', 100.00, '2026-02-05 09:02:40', '2026-02-05 09:05:31', NULL),
(6, 11, 'Test Customer', 'test5@example.com', '555-0005', '', '', '', '', '', 'wholesale', 'active', '', 0.00, '2026-02-05 09:04:26', '2026-02-05 09:05:48', '2026-02-05 09:05:48'),
(7, 23, 'Test Customer', 'testcust@example.com', '1234567890', '123 Main St', '', '', '', '', 'retail', 'active', '', 0.00, '2026-02-08 08:31:31', '2026-02-08 08:31:31', NULL);

-- --------------------------------------------------------

--
-- Table structure for table `customer_returns`
--

CREATE TABLE `customer_returns` (
  `id` bigint UNSIGNED NOT NULL,
  `return_number` varchar(50) NOT NULL,
  `customer_id` bigint UNSIGNED NOT NULL,
  `customer_name` varchar(255) NOT NULL,
  `order_id` bigint UNSIGNED DEFAULT NULL,
  `order_number` varchar(50) DEFAULT NULL,
  `total_amount` decimal(10,2) NOT NULL DEFAULT '0.00',
  `status` enum('pending','approved','rejected','completed') NOT NULL DEFAULT 'pending',
  `request_date` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `processed_date` timestamp NULL DEFAULT NULL,
  `refund_method` enum('cash','store_credit','original_payment') NOT NULL,
  `notes` text,
  `processed_by` varchar(255) DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

--
-- Dumping data for table `customer_returns`
--

INSERT INTO `customer_returns` (`id`, `return_number`, `customer_id`, `customer_name`, `order_id`, `order_number`, `total_amount`, `status`, `request_date`, `processed_date`, `refund_method`, `notes`, `processed_by`, `created_at`, `updated_at`, `deleted_at`) VALUES
(2, 'RET-48739', 2, 'John Smith', NULL, 'INV-12345', 599.98, 'approved', '2026-02-08 11:05:39', '2026-02-08 11:05:39', 'store_credit', 'Customer ordered wrong size', 'Admin', '2026-02-08 11:05:39', '2026-02-08 11:05:39', NULL),
(3, 'RET-49561', 2, 'John Smith', NULL, '', 500.00, 'approved', '2026-02-08 11:19:22', '2026-02-08 11:19:22', 'cash', '', 'Admin', '2026-02-08 11:19:22', '2026-02-08 11:19:22', NULL),
(4, 'RET-1770550208772566', 2, 'John Smith', NULL, 'INV-12345', 599.98, 'approved', '2026-02-08 11:30:09', '2026-02-08 11:34:44', 'store_credit', 'Customer ordered wrong size', 'Admin', '2026-02-08 11:30:09', '2026-02-08 11:34:44', NULL);

-- --------------------------------------------------------

--
-- Table structure for table `customer_return_items`
--

CREATE TABLE `customer_return_items` (
  `id` bigint UNSIGNED NOT NULL,
  `return_id` bigint UNSIGNED NOT NULL,
  `product_id` bigint UNSIGNED DEFAULT NULL,
  `product_name` varchar(255) NOT NULL,
  `variant_id` bigint UNSIGNED DEFAULT NULL,
  `variant_name` varchar(255) DEFAULT NULL,
  `quantity` int NOT NULL DEFAULT '1',
  `price` decimal(10,2) NOT NULL DEFAULT '0.00',
  `reason` varchar(255) NOT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

--
-- Dumping data for table `customer_return_items`
--

INSERT INTO `customer_return_items` (`id`, `return_id`, `product_id`, `product_name`, `variant_id`, `variant_name`, `quantity`, `price`, `reason`, `created_at`, `updated_at`) VALUES
(1, 2, NULL, 'Premium T-Shirt', NULL, 'Large / Blue', 2, 299.99, 'Wrong size', '2026-02-08 11:05:39', '2026-02-08 11:05:39'),
(2, 3, 28, 'Test Product for Returns', NULL, '', 5, 100.00, 'Defective product', '2026-02-08 11:19:22', '2026-02-08 11:19:22'),
(3, 4, NULL, 'Premium T-Shirt', NULL, 'Large / Blue', 2, 299.99, 'Wrong size', '2026-02-08 11:30:09', '2026-02-08 11:30:09');

-- --------------------------------------------------------

--
-- Table structure for table `locations`
--

CREATE TABLE `locations` (
  `id` bigint UNSIGNED NOT NULL,
  `name` varchar(255) NOT NULL,
  `address` varchar(255) DEFAULT NULL,
  `contact_person` varchar(255) DEFAULT NULL,
  `is_default` tinyint(1) DEFAULT '0',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

--
-- Dumping data for table `locations`
--

INSERT INTO `locations` (`id`, `name`, `address`, `contact_person`, `is_default`, `created_at`, `updated_at`, `deleted_at`) VALUES
(1, 'Main Office', '123 Main Street, Dhaka', 'John Doe', 1, '2026-02-02 05:27:35', '2026-02-02 05:27:35', NULL),
(2, 'wirehouse 2', '123 Main Street, Netrokona', 'New Contact person', 0, '2026-02-05 09:44:23', '2026-02-05 09:44:23', NULL);

-- --------------------------------------------------------

--
-- Table structure for table `order_items`
--

CREATE TABLE `order_items` (
  `id` bigint UNSIGNED NOT NULL,
  `sell_id` bigint UNSIGNED NOT NULL,
  `product_id` bigint UNSIGNED DEFAULT NULL,
  `variant_id` bigint UNSIGNED DEFAULT NULL,
  `product_name` varchar(255) NOT NULL,
  `variant_name` varchar(255) DEFAULT NULL,
  `quantity` int NOT NULL DEFAULT '1',
  `unit_price` decimal(10,2) NOT NULL DEFAULT '0.00',
  `total_price` decimal(10,2) NOT NULL DEFAULT '0.00',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

--
-- Dumping data for table `order_items`
--

INSERT INTO `order_items` (`id`, `sell_id`, `product_id`, `variant_id`, `product_name`, `variant_name`, `quantity`, `unit_price`, `total_price`, `created_at`, `updated_at`) VALUES
(1, 11, NULL, NULL, 'Fresh Mustard Oil', '', 2, 99.99, 199.98, '2026-02-08 09:09:00', '2026-02-08 09:09:00'),
(2, 11, NULL, NULL, 'Organic Honey', '', 1, 119.99, 119.99, '2026-02-08 09:09:00', '2026-02-08 09:09:00'),
(5, 14, 17, NULL, 'Fresh Mustard Oil', '', 2, 99.99, 199.98, '2026-02-08 09:13:37', '2026-02-08 09:13:37'),
(6, 14, 18, NULL, 'Organic Honey', '', 1, 119.99, 119.99, '2026-02-08 09:13:37', '2026-02-08 09:13:37'),
(8, 16, 16, NULL, 'New product 1', '', 1, 0.00, 0.00, '2026-02-09 07:54:58', '2026-02-09 07:54:58'),
(9, 17, 17, NULL, 'Fresh Mustard Oil', '', 2, 0.00, 0.00, '2026-02-09 08:38:54', '2026-02-09 08:38:54'),
(10, 17, 18, NULL, 'Organic Honey', '', 1, 0.00, 0.00, '2026-02-09 08:38:54', '2026-02-09 08:38:54'),
(11, 18, 16, NULL, 'Product A', '', 1, 0.00, 0.00, '2026-02-09 08:45:05', '2026-02-09 08:45:05'),
(12, 19, 17, NULL, 'Product', '', 1, 0.00, 0.00, '2026-02-09 08:45:44', '2026-02-09 08:45:44'),
(13, 20, 18, NULL, 'Product C', '', 1, 0.00, 0.00, '2026-02-09 08:46:48', '2026-02-09 08:46:48'),
(14, 21, 19, NULL, 'Test Product', '', 2, 0.00, 0.00, '2026-02-09 08:46:54', '2026-02-09 08:46:54'),
(15, 22, 19, NULL, 'Test Product', '', 2, 0.00, 0.00, '2026-02-09 08:55:26', '2026-02-09 08:55:26');

-- --------------------------------------------------------

--
-- Table structure for table `order_shipments`
--

CREATE TABLE `order_shipments` (
  `id` bigint UNSIGNED NOT NULL,
  `sell_id` bigint UNSIGNED NOT NULL,
  `tracking_number` varchar(100) NOT NULL,
  `carrier` varchar(100) NOT NULL,
  `shipping_method` varchar(50) DEFAULT NULL,
  `status` enum('pending','picked_up','in_transit','out_for_delivery','delivered','failed','returned') DEFAULT 'pending',
  `shipped_at` timestamp NULL DEFAULT NULL,
  `estimated_delivery` timestamp NULL DEFAULT NULL,
  `delivered_at` timestamp NULL DEFAULT NULL,
  `shipping_cost` decimal(10,2) DEFAULT '0.00',
  `weight` decimal(10,2) DEFAULT NULL COMMENT 'Weight in kg',
  `dimensions` varchar(50) DEFAULT NULL COMMENT 'LxWxH in cm',
  `notes` text,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

--
-- Dumping data for table `order_shipments`
--

INSERT INTO `order_shipments` (`id`, `sell_id`, `tracking_number`, `carrier`, `shipping_method`, `status`, `shipped_at`, `estimated_delivery`, `delivered_at`, `shipping_cost`, `weight`, `dimensions`, `notes`, `created_at`, `updated_at`) VALUES
(1, 14, 'DHL-BD-123456789', 'DHL Express', 'Express Delivery', 'in_transit', NULL, '2026-02-11 18:00:00', NULL, 150.00, 2.50, '30x20x15', 'Fragile items - handle with care', '2026-02-09 04:48:52', '2026-02-09 06:18:01'),
(2, 14, 'DHL-BD-123456789', 'DHL Express', '', 'pending', NULL, NULL, NULL, 150.00, 0.00, '', '', '2026-02-09 06:14:33', '2026-02-09 06:14:33');

-- --------------------------------------------------------

--
-- Table structure for table `permissions`
--

CREATE TABLE `permissions` (
  `id` bigint UNSIGNED NOT NULL,
  `name` varchar(50) NOT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

--
-- Dumping data for table `permissions`
--

INSERT INTO `permissions` (`id`, `name`, `created_at`, `updated_at`) VALUES
(1, 'Dashboard', '2026-02-05 11:50:28', '2026-02-05 11:50:28'),
(2, 'Products', '2026-02-05 11:50:28', '2026-02-05 11:50:28'),
(3, 'Categories', '2026-02-05 11:50:28', '2026-02-05 11:50:28'),
(4, 'Attributes', '2026-02-05 11:50:28', '2026-02-05 11:50:28'),
(5, 'Coupons', '2026-02-05 11:50:28', '2026-02-05 11:50:28'),
(6, 'Customers', '2026-02-05 11:50:28', '2026-02-05 11:50:28'),
(7, 'Orders', '2026-02-05 11:50:28', '2026-02-05 11:50:28'),
(8, 'POS', '2026-02-05 11:50:28', '2026-02-05 11:50:28'),
(9, 'Sells', '2026-02-05 11:50:28', '2026-02-05 11:50:28'),
(10, 'Staff', '2026-02-05 11:50:28', '2026-02-05 11:50:28'),
(11, 'Settings', '2026-02-05 11:50:28', '2026-02-05 11:50:28'),
(12, 'International', '2026-02-05 11:50:28', '2026-02-05 11:50:28'),
(13, 'Store', '2026-02-05 11:50:28', '2026-02-05 11:50:28'),
(14, 'Pages', '2026-02-05 11:50:28', '2026-02-05 11:50:28');

-- --------------------------------------------------------

--
-- Table structure for table `products`
--

CREATE TABLE `products` (
  `id` bigint UNSIGNED NOT NULL,
  `name` varchar(255) NOT NULL,
  `description` text,
  `category_id` bigint UNSIGNED DEFAULT NULL,
  `price` decimal(10,2) NOT NULL DEFAULT '0.00',
  `sale_price` decimal(10,2) NOT NULL DEFAULT '0.00',
  `stock` int NOT NULL DEFAULT '0',
  `sku` varchar(255) DEFAULT NULL,
  `barcode` varchar(255) DEFAULT NULL,
  `published` tinyint(1) NOT NULL DEFAULT '0',
  `vendor_id` bigint UNSIGNED DEFAULT NULL,
  `receipt_number` varchar(255) DEFAULT NULL,
  `location_id` bigint UNSIGNED DEFAULT NULL,
  `image` varchar(500) DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

--
-- Dumping data for table `products`
--

INSERT INTO `products` (`id`, `name`, `description`, `category_id`, `price`, `sale_price`, `stock`, `sku`, `barcode`, `published`, `vendor_id`, `receipt_number`, `location_id`, `image`, `created_at`, `updated_at`, `deleted_at`) VALUES
(16, 'New product 1', 'Soft cotton 1', 2, 20.00, 25.00, 50, 'New-product-1', 'h577e9971', 1, NULL, 'RCPT-004343-1', 1, '', '2026-02-02 11:14:50', '2026-02-02 11:48:16', '2026-02-02 11:48:16'),
(17, 'New product 1', 'Soft cotton 1', 2, 20.00, 25.00, 30, 'New-product-2', 'h577e9973', 1, NULL, 'RCPT-004343-3', 1, '', '2026-02-02 11:18:39', '2026-02-05 10:17:07', NULL),
(18, 'New product 34', 'Soft cotton 4', 2, 20.00, 25.00, 50, 'New-product-4', 'h577e9977', 1, NULL, 'RCPT-004343-477', 1, '', '2026-02-02 11:20:41', '2026-02-02 11:20:41', NULL),
(19, 'Test With Slash', '', NULL, 10.00, 0.00, 0, 'WITH-SLASH-001', 'WS001', 1, NULL, '', NULL, '', '2026-02-02 11:25:45', '2026-02-02 11:25:45', NULL),
(20, 'Test No Slash', '', NULL, 10.00, 0.00, 0, 'NO-SLASH-003', 'NS003', 1, NULL, '', NULL, '', '2026-02-02 11:27:09', '2026-02-02 11:27:09', NULL),
(21, 'Test With Slash', '', NULL, 10.00, 0.00, 0, 'WITH-SLASH-003', 'WS003', 1, NULL, '', NULL, '', '2026-02-02 11:27:09', '2026-02-02 12:29:36', '2026-02-02 12:29:37'),
(22, 'No Slash Product', '', NULL, 50.00, 0.00, 0, 'NOSLASH-IMG-003', 'NSIMG003', 1, NULL, '', NULL, '', '2026-02-02 11:27:55', '2026-02-02 11:48:33', '2026-02-02 11:48:33'),
(23, 'update', 'update description', 1, 19.00, 25.00, 50, 'TSHIRT-RED-update', '1234567554622', 1, NULL, 'RCPT-8955', 1, '', '2026-02-02 11:30:38', '2026-02-02 11:35:35', NULL),
(24, 'Test Smartphone PRO', 'A test product for CRUD validation', 2, 1299.99, 899.99, 30, 'TEST-SKU-001', 'TEST-BAR-001', 1, NULL, 'REC-001', 1, 'uploads/products/1770113800586539125.jpg', '2026-02-03 10:15:46', '2026-02-03 10:17:08', '2026-02-03 10:17:08'),
(25, 'Smoke Test Item', '', NULL, 10.00, 0.00, 0, 'SMOKE-001', 'SMOKE-BAR-001', 0, NULL, '', NULL, '', '2026-02-03 10:25:18', '2026-02-03 10:25:18', '2026-02-03 10:25:18'),
(26, 'Test CRUD Product UPDATED', 'Created by CRUD test', 2, 149.99, 79.99, 50, 'CRUD-TEST-001', 'CRUD-BAR-001', 1, NULL, 'RCPT-CRUD-001', 1, '', '2026-02-05 05:43:13', '2026-02-05 05:51:44', '2026-02-05 05:51:44'),
(27, 'Clean Code Test UPDATED', '', NULL, 75.00, 0.00, 0, 'CLEAN-001', 'CLN-001', 0, NULL, '', NULL, '', '2026-02-05 06:02:45', '2026-02-05 06:02:44', '2026-02-05 06:02:45'),
(28, 'Test Product for Returns', '', 1, 100.00, 0.00, 45, '', '', 0, NULL, '', 1, '', '2026-02-08 11:19:22', '2026-02-08 11:19:37', NULL);

-- --------------------------------------------------------

--
-- Table structure for table `product_attributes`
--

CREATE TABLE `product_attributes` (
  `product_id` bigint UNSIGNED NOT NULL,
  `attribute_id` bigint UNSIGNED NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

--
-- Dumping data for table `product_attributes`
--

INSERT INTO `product_attributes` (`product_id`, `attribute_id`) VALUES
(16, 1),
(17, 1),
(18, 1),
(23, 1),
(24, 1),
(26, 1),
(27, 1),
(16, 2),
(17, 2),
(18, 2),
(23, 2),
(26, 2);

-- --------------------------------------------------------

--
-- Table structure for table `product_images`
--

CREATE TABLE `product_images` (
  `id` bigint UNSIGNED NOT NULL,
  `product_id` bigint UNSIGNED NOT NULL,
  `path` varchar(255) NOT NULL,
  `position` int NOT NULL DEFAULT '0',
  `is_primary` tinyint(1) NOT NULL DEFAULT '0',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

--
-- Dumping data for table `product_images`
--

INSERT INTO `product_images` (`id`, `product_id`, `path`, `position`, `is_primary`, `created_at`) VALUES
(3, 16, 'uploads/products/1770030890042874776.webp', 0, 1, '2026-02-02 11:14:50'),
(4, 16, 'uploads/products/1770030890043066077.png', 1, 0, '2026-02-02 11:14:50'),
(5, 17, 'uploads/products/1770031118892914482.png', 0, 1, '2026-02-02 11:18:39'),
(6, 17, 'uploads/products/1770031118893206160.jpg', 1, 0, '2026-02-02 11:18:39'),
(7, 18, 'uploads/products/1770031241092737048.pdf', 0, 1, '2026-02-02 11:20:41'),
(8, 18, 'uploads/products/1770031241100060267.pdf', 1, 0, '2026-02-02 11:20:41'),
(9, 22, 'uploads/products/1770031674717412202.jpg', 0, 1, '2026-02-02 11:27:55'),
(10, 22, 'uploads/products/1770031674717556064.jpg', 1, 0, '2026-02-02 11:27:55'),
(19, 23, 'uploads/products/1770032135400721350.png', 0, 1, '2026-02-02 11:35:35'),
(23, 24, 'uploads/products/1770113800586699954.jpg', 0, 1, '2026-02-03 10:16:41');

-- --------------------------------------------------------

--
-- Table structure for table `product_variants`
--

CREATE TABLE `product_variants` (
  `id` bigint UNSIGNED NOT NULL,
  `product_id` bigint UNSIGNED NOT NULL,
  `name` varchar(255) NOT NULL,
  `attributes` json DEFAULT NULL,
  `price` decimal(10,2) NOT NULL DEFAULT '0.00',
  `sale_price` decimal(10,2) NOT NULL DEFAULT '0.00',
  `stock` int NOT NULL DEFAULT '0',
  `sku` varchar(255) DEFAULT NULL,
  `barcode` varchar(255) DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

--
-- Dumping data for table `product_variants`
--

INSERT INTO `product_variants` (`id`, `product_id`, `name`, `attributes`, `price`, `sale_price`, `stock`, `sku`, `barcode`, `created_at`, `updated_at`, `deleted_at`) VALUES
(45, 16, 'Red / S', '{\"Size\": \"S\", \"Color\": \"red\"}', 19.99, 14.99, 10, 'TS-RED-S', '111111111111', '2026-02-02 11:14:50', '2026-02-02 11:14:50', NULL),
(46, 16, 'Red / M', '{\"Size\": \"M\", \"Color\": \"red\"}', 19.99, 14.99, 8, 'TS-RED-M', '111111111112', '2026-02-02 11:14:50', '2026-02-02 11:14:50', NULL),
(47, 16, 'Green / S', '{\"Size\": \"S\", \"Color\": \"green\"}', 19.99, 14.99, 6, 'TS-GRN-S', '111111111113', '2026-02-02 11:14:50', '2026-02-02 11:14:50', NULL),
(48, 16, 'Blue / L', '{\"Size\": \"L\", \"Color\": \"blue\"}', 19.99, 14.99, 4, 'TS-BLU-L', '111111111114', '2026-02-02 11:14:50', '2026-02-02 11:14:50', NULL),
(49, 17, 'Red / S', '{\"Size\": \"S\", \"Color\": \"red\"}', 19.99, 14.99, 10, 'TS-RED-S', '111111111111', '2026-02-02 11:18:39', '2026-02-02 11:18:39', NULL),
(50, 17, 'Red / M', '{\"Size\": \"M\", \"Color\": \"red\"}', 19.99, 14.99, 8, 'TS-RED-M', '111111111112', '2026-02-02 11:18:39', '2026-02-02 11:18:39', NULL),
(51, 17, 'Green / S', '{\"Size\": \"S\", \"Color\": \"green\"}', 19.99, 14.99, 6, 'TS-GRN-S', '111111111113', '2026-02-02 11:18:39', '2026-02-02 11:18:39', NULL),
(52, 17, 'Blue / L', '{\"Size\": \"L\", \"Color\": \"blue\"}', 19.99, 14.99, 4, 'TS-BLU-L', '111111111114', '2026-02-02 11:18:39', '2026-02-02 11:18:39', NULL),
(53, 18, 'Red / S', '{\"Size\": \"S\", \"Color\": \"red\"}', 19.99, 14.99, 10, 'TS-RED-S', '111111111111', '2026-02-02 11:20:41', '2026-02-02 11:20:41', NULL),
(54, 18, 'Red / M', '{\"Size\": \"M\", \"Color\": \"red\"}', 19.99, 14.99, 8, 'TS-RED-M', '111111111112', '2026-02-02 11:20:41', '2026-02-02 11:20:41', NULL),
(55, 18, 'Green / S', '{\"Size\": \"S\", \"Color\": \"green\"}', 19.99, 14.99, 6, 'TS-GRN-S', '111111111113', '2026-02-02 11:20:41', '2026-02-02 11:20:41', NULL),
(56, 18, 'Blue / L', '{\"Size\": \"L\", \"Color\": \"blue\"}', 19.99, 14.99, 4, 'TS-BLU-L', '111111111114', '2026-02-02 11:20:41', '2026-02-02 11:20:41', NULL),
(57, 23, 'Red / S', '{\"Size\": \"S\", \"Color\": \"red\"}', 19.99, 14.99, 10, 'TS-RED-S', '111111111111', '2026-02-02 11:35:35', '2026-02-02 11:35:35', NULL),
(58, 23, 'Red / M', '{\"Size\": \"M\", \"Color\": \"red\"}', 19.99, 14.99, 8, 'TS-RED-M', '111111111112', '2026-02-02 11:35:35', '2026-02-02 11:35:35', NULL),
(59, 23, 'Green / S', '{\"Size\": \"S\", \"Color\": \"green\"}', 19.99, 14.99, 6, 'TS-GRN-S', '111111111113', '2026-02-02 11:35:35', '2026-02-02 11:35:35', NULL),
(60, 23, 'Blue / L', '{\"Size\": \"L\", \"Color\": \"blue\"}', 19.99, 14.99, 4, 'TS-BLU-L', '111111111114', '2026-02-02 11:35:35', '2026-02-02 11:35:35', NULL),
(61, 24, 'Small / Red', '{\"Size\": \"S\", \"Color\": \"red\"}', 950.00, 850.00, 20, 'VAR-001', 'VBAR-001', '2026-02-03 10:15:46', '2026-02-03 10:16:40', '2026-02-03 10:16:41'),
(62, 24, 'Large / Blue', '{\"Size\": \"L\", \"Color\": \"blue\"}', 1050.00, 950.00, 10, 'VAR-002', 'VBAR-002', '2026-02-03 10:15:46', '2026-02-03 10:16:40', '2026-02-03 10:16:41'),
(63, 24, 'XL / Green', '{\"Size\": \"XL\", \"Color\": \"green\"}', 1200.00, 1100.00, 15, 'VAR-UPD-001', 'VBAR-UPD-001', '2026-02-03 10:16:41', '2026-02-03 10:16:41', NULL),
(64, 26, 'Red / S', '{\"Size\": \"S\", \"Color\": \"red\"}', 95.00, 75.00, 20, 'CRUD-RED-S', 'CB-RED-S', '2026-02-05 05:43:13', '2026-02-05 05:43:43', '2026-02-05 05:43:44'),
(65, 26, 'Blue / L', '{\"Size\": \"L\", \"Color\": \"blue\"}', 140.00, 110.00, 15, 'CRUD-BLU-L', 'CB-BLU-L', '2026-02-05 05:43:44', '2026-02-05 05:43:44', NULL),
(66, 27, 'S', '{\"Size\": \"S\"}', 48.00, 0.00, 10, 'CLN-S', 'CLN-S-1', '2026-02-05 06:02:45', '2026-02-05 06:02:45', NULL);

-- --------------------------------------------------------

--
-- Table structure for table `roles`
--

CREATE TABLE `roles` (
  `id` bigint UNSIGNED NOT NULL,
  `title` varchar(255) NOT NULL,
  `status` tinyint(1) NOT NULL DEFAULT '0',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

--
-- Dumping data for table `roles`
--

INSERT INTO `roles` (`id`, `title`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES
(1, 'Admin', 1, '2026-01-26 11:39:39', '2026-01-26 11:39:39', NULL),
(2, 'Manager', 1, '2026-01-26 11:39:39', '2026-01-26 11:39:39', NULL),
(3, 'Customer', 1, '2026-01-26 11:39:39', '2026-01-26 11:39:39', NULL),
(4, 'Vendor', 1, '2026-01-26 11:39:39', '2026-01-26 11:39:39', NULL),
(5, 'Staff', 1, '2026-01-26 11:39:39', '2026-01-26 11:39:39', NULL);

-- --------------------------------------------------------

--
-- Table structure for table `role_permissions`
--

CREATE TABLE `role_permissions` (
  `id` bigint UNSIGNED NOT NULL,
  `role_id` bigint UNSIGNED NOT NULL,
  `permission_id` bigint UNSIGNED NOT NULL,
  `read` tinyint(1) NOT NULL DEFAULT '0',
  `write` tinyint(1) NOT NULL DEFAULT '0',
  `delete` tinyint(1) NOT NULL DEFAULT '0',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

--
-- Dumping data for table `role_permissions`
--

INSERT INTO `role_permissions` (`id`, `role_id`, `permission_id`, `read`, `write`, `delete`, `created_at`, `updated_at`) VALUES
(7, 5, 1, 1, 0, 0, '2026-02-05 12:59:38', '2026-02-05 12:59:38'),
(8, 5, 2, 1, 1, 0, '2026-02-05 12:59:38', '2026-02-05 12:59:38'),
(9, 5, 3, 1, 1, 1, '2026-02-05 12:59:38', '2026-02-05 12:59:38'),
(10, 5, 4, 0, 0, 0, '2026-02-05 12:59:38', '2026-02-05 12:59:38'),
(11, 5, 5, 0, 0, 0, '2026-02-05 12:59:38', '2026-02-05 12:59:38'),
(12, 5, 6, 0, 0, 0, '2026-02-05 12:59:38', '2026-02-05 12:59:38'),
(13, 5, 7, 0, 0, 0, '2026-02-05 12:59:38', '2026-02-05 12:59:38'),
(14, 5, 8, 0, 0, 0, '2026-02-05 12:59:38', '2026-02-05 12:59:38'),
(15, 5, 9, 0, 0, 0, '2026-02-05 12:59:38', '2026-02-05 12:59:38'),
(16, 5, 10, 0, 0, 0, '2026-02-05 12:59:38', '2026-02-05 12:59:38'),
(17, 5, 11, 0, 0, 0, '2026-02-05 12:59:38', '2026-02-05 12:59:38'),
(18, 5, 12, 0, 0, 0, '2026-02-05 12:59:38', '2026-02-05 12:59:38'),
(19, 5, 13, 0, 0, 0, '2026-02-05 12:59:38', '2026-02-05 12:59:38'),
(20, 5, 14, 0, 0, 0, '2026-02-05 12:59:38', '2026-02-05 12:59:38'),
(21, 1, 1, 1, 0, 0, '2026-02-08 06:33:13', '2026-02-08 06:33:13'),
(22, 1, 2, 1, 1, 0, '2026-02-08 06:33:13', '2026-02-08 06:33:13'),
(23, 1, 3, 1, 1, 1, '2026-02-08 06:33:13', '2026-02-08 06:33:13'),
(24, 1, 4, 0, 0, 0, '2026-02-08 06:33:13', '2026-02-08 06:33:13'),
(25, 1, 5, 0, 0, 0, '2026-02-08 06:33:13', '2026-02-08 06:33:13'),
(26, 1, 6, 0, 0, 0, '2026-02-08 06:33:13', '2026-02-08 06:33:13'),
(27, 1, 7, 0, 0, 0, '2026-02-08 06:33:13', '2026-02-08 06:33:13'),
(28, 1, 8, 0, 0, 0, '2026-02-08 06:33:13', '2026-02-08 06:33:13'),
(29, 1, 9, 0, 0, 0, '2026-02-08 06:33:13', '2026-02-08 06:33:13'),
(30, 1, 10, 0, 0, 0, '2026-02-08 06:33:13', '2026-02-08 06:33:13'),
(31, 1, 11, 0, 0, 0, '2026-02-08 06:33:13', '2026-02-08 06:33:13'),
(32, 1, 12, 0, 0, 0, '2026-02-08 06:33:13', '2026-02-08 06:33:13'),
(33, 1, 13, 0, 0, 0, '2026-02-08 06:33:13', '2026-02-08 06:33:13'),
(34, 1, 14, 0, 0, 0, '2026-02-08 06:33:13', '2026-02-08 06:33:13');

-- --------------------------------------------------------

--
-- Table structure for table `salary_payments`
--

CREATE TABLE `salary_payments` (
  `id` bigint UNSIGNED NOT NULL,
  `staff_id` bigint UNSIGNED NOT NULL,
  `month` varchar(10) NOT NULL,
  `amount` decimal(12,2) NOT NULL DEFAULT '0.00',
  `paid_amount` decimal(12,2) NOT NULL DEFAULT '0.00',
  `status` enum('Paid','Pending','Partial') NOT NULL DEFAULT 'Pending',
  `payment_date` varchar(20) DEFAULT NULL,
  `payment_method` varchar(50) DEFAULT NULL,
  `notes` text,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

--
-- Dumping data for table `salary_payments`
--

INSERT INTO `salary_payments` (`id`, `staff_id`, `month`, `amount`, `paid_amount`, `status`, `payment_date`, `payment_method`, `notes`, `created_at`, `updated_at`) VALUES
(1, 4, 'Jan 2026', 5000.00, 5000.00, 'Paid', '2026-01-15', 'Bank Transfer', 'Full payment', '2026-02-05 11:35:49', '2026-02-05 11:43:49'),
(2, 4, 'Feb 2026', 5000.00, 5000.00, 'Paid', '2026-02-15', 'Cash', 'Remaining balance', '2026-02-05 11:35:49', '2026-02-05 11:43:49');

-- --------------------------------------------------------

--
-- Table structure for table `schema_migrations`
--

CREATE TABLE `schema_migrations` (
  `version` bigint NOT NULL,
  `dirty` tinyint(1) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

--
-- Dumping data for table `schema_migrations`
--

INSERT INTO `schema_migrations` (`version`, `dirty`) VALUES
(37, 0);

-- --------------------------------------------------------

--
-- Table structure for table `sells`
--

CREATE TABLE `sells` (
  `id` bigint UNSIGNED NOT NULL,
  `invoice_no` varchar(50) NOT NULL,
  `order_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `customer_id` bigint UNSIGNED DEFAULT NULL,
  `customer_name` varchar(255) NOT NULL,
  `shipping_address_id` bigint UNSIGNED DEFAULT NULL,
  `shipping_full_name` varchar(255) DEFAULT NULL,
  `shipping_phone` varchar(20) DEFAULT NULL,
  `shipping_email` varchar(255) DEFAULT NULL,
  `shipping_address_line1` varchar(255) DEFAULT NULL,
  `shipping_address_line2` varchar(255) DEFAULT NULL,
  `shipping_city` varchar(100) DEFAULT NULL,
  `shipping_state` varchar(100) DEFAULT NULL,
  `shipping_postal_code` varchar(20) DEFAULT NULL,
  `shipping_country` varchar(100) DEFAULT NULL,
  `shipping_address_type` enum('home','office','other') DEFAULT NULL,
  `method` varchar(50) NOT NULL DEFAULT 'Cash',
  `amount` decimal(10,2) NOT NULL DEFAULT '0.00',
  `shipping_cost` decimal(10,2) NOT NULL DEFAULT '0.00',
  `shipping_method` varchar(50) DEFAULT NULL,
  `discount` decimal(10,2) NOT NULL DEFAULT '0.00',
  `status` varchar(20) NOT NULL DEFAULT 'Pending',
  `payment_status` enum('pending','paid','partially_paid','refunded','failed') DEFAULT 'pending',
  `fulfillment_status` enum('unfulfilled','processing','shipped','delivered','cancelled') DEFAULT 'unfulfilled',
  `tracking_number` varchar(100) DEFAULT NULL,
  `carrier` varchar(100) DEFAULT NULL,
  `shipped_at` timestamp NULL DEFAULT NULL,
  `delivered_at` timestamp NULL DEFAULT NULL,
  `notes` text,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

--
-- Dumping data for table `sells`
--

INSERT INTO `sells` (`id`, `invoice_no`, `order_time`, `customer_id`, `customer_name`, `shipping_address_id`, `shipping_full_name`, `shipping_phone`, `shipping_email`, `shipping_address_line1`, `shipping_address_line2`, `shipping_city`, `shipping_state`, `shipping_postal_code`, `shipping_country`, `shipping_address_type`, `method`, `amount`, `shipping_cost`, `shipping_method`, `discount`, `status`, `payment_status`, `fulfillment_status`, `tracking_number`, `carrier`, `shipped_at`, `delivered_at`, `notes`, `created_at`, `updated_at`, `deleted_at`) VALUES
(11, 'INV-WITH-ITEMS-001', '2026-02-08 09:09:00', NULL, 'John Doe', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 'Card', 329.97, 20.00, NULL, 10.00, 'Processing', 'pending', 'unfulfilled', NULL, NULL, NULL, NULL, 'Express delivery', '2026-02-08 09:09:00', '2026-02-08 09:09:00', NULL),
(12, 'INV-NO-ITEMS-001', '2026-02-08 09:09:00', NULL, 'Jane Smith', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 'Cash', 159.99, 0.00, NULL, 0.00, 'Pending', 'pending', 'unfulfilled', NULL, NULL, NULL, NULL, '', '2026-02-08 09:09:00', '2026-02-08 09:09:00', NULL),
(14, 'INV-001', '2026-02-08 09:13:37', 7, 'John Doe', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 'Card', 329.97, 20.00, NULL, 10.00, 'Processing', 'pending', 'shipped', 'DHL-BD-123456789', 'DHL Express', '2026-02-09 06:18:01', NULL, 'Express delivery', '2026-02-08 09:13:37', '2026-02-09 06:18:01', NULL),
(16, 'INV-1770623697', '2026-02-09 07:54:58', 2, 'Test Customer', 6, 'Jane Smith', '+8801812345678', '', '456 Commerce Road', '', 'Chittagong', 'Chittagong Division', '4000', 'Bangladesh', 'office', 'Card', 1500.00, 0.00, '', 0.00, 'Pending', 'pending', 'unfulfilled', '', '', NULL, NULL, '', '2026-02-09 07:54:58', '2026-02-09 07:54:58', NULL),
(17, 'INV-TEST-001', '2026-02-09 08:38:54', 7, 'John Doe', NULL, 'John Doe', '+8801912345678', 'john@example.com', '789 Custom Street', 'Suite 100', 'Sylhet', 'Sylhet Division', '3100', 'Bangladesh', 'office', 'Card', 329.97, 20.00, '', 10.00, 'Processing', 'pending', 'unfulfilled', '', '', NULL, NULL, 'Express delivery', '2026-02-09 08:38:54', '2026-02-09 08:38:54', NULL),
(18, 'INV-1770626705', '2026-02-09 08:45:05', 2, 'Test Auto Default', 6, 'Jane Smith', '+8801812345678', '', '456 Commerce Road', '', 'Chittagong', 'Chittagong Division', '4000', 'Bangladesh', 'office', 'Card', 500.00, 0.00, '', 0.00, 'Pending', 'pending', 'unfulfilled', '', '', NULL, NULL, '', '2026-02-09 08:45:05', '2026-02-09 08:45:05', NULL),
(19, 'INV-1770626743', '2026-02-09 08:45:44', 2, 'Test', 4, 'John Smith', '+8801712345678', 'john@example.com', '123 Main Street', 'Apartment 4B', 'Dhaka', 'Dhaka Division', '1205', 'Bangladesh', 'home', 'Online', 750.00, 0.00, '', 0.00, 'Pending', 'pending', 'unfulfilled', '', '', NULL, NULL, '', '2026-02-09 08:45:44', '2026-02-09 08:45:44', NULL),
(20, 'INV-1770626807', '2026-02-09 08:46:48', 2, 'Test Custom', NULL, 'Custom Name', '+8801912345678', '', '123 Custom Street', '', 'Rajshahi', '', '', 'Bangladesh', 'other', 'Card', 1000.00, 0.00, '', 0.00, 'Pending', 'pending', 'unfulfilled', '', '', NULL, NULL, '', '2026-02-09 08:46:48', '2026-02-09 08:46:48', NULL),
(21, 'INV-1770626814', '2026-02-09 08:46:54', 2, 'Full Custom Test', NULL, 'Jane Recipient', '+8801812345678', 'jane@test.com', '456 Full Address Street', 'Building 2, Floor 3', 'Sylhet', 'Sylhet Division', '3100', 'Bangladesh', 'office', 'Online', 1500.00, 0.00, '', 0.00, 'Pending', 'pending', 'unfulfilled', '', '', NULL, NULL, '', '2026-02-09 08:46:54', '2026-02-09 08:46:54', NULL),
(22, 'INV-1770627326', '2026-02-09 08:55:26', 2, 'Full Custom Test', NULL, 'Jane Recipient', '+8801812345678', 'jane@test.com', '456 Full Address Street', 'Building 2, Floor 3', 'Sylhet', 'Sylhet Division', '3100', 'Bangladesh', 'office', 'Online', 1500.00, 0.00, '', 0.00, 'Pending', 'pending', 'unfulfilled', '', '', NULL, NULL, '', '2026-02-09 08:55:26', '2026-02-09 08:55:26', NULL);

-- --------------------------------------------------------

--
-- Table structure for table `shipment_tracking_history`
--

CREATE TABLE `shipment_tracking_history` (
  `id` bigint UNSIGNED NOT NULL,
  `shipment_id` bigint UNSIGNED NOT NULL,
  `status` varchar(50) NOT NULL,
  `location` varchar(255) DEFAULT NULL,
  `description` text,
  `event_time` timestamp NOT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

--
-- Dumping data for table `shipment_tracking_history`
--

INSERT INTO `shipment_tracking_history` (`id`, `shipment_id`, `status`, `location`, `description`, `event_time`, `created_at`) VALUES
(1, 1, 'pending', '', 'Shipment created with DHL Express', '2026-02-09 04:48:52', '2026-02-09 04:48:52'),
(2, 1, 'in_transit', 'Dhaka Distribution Center', 'Package in transit to destination', '2026-02-09 04:48:56', '2026-02-09 04:48:56'),
(3, 2, 'pending', '', 'Shipment created with DHL Express', '2026-02-09 06:14:33', '2026-02-09 06:14:33'),
(4, 1, 'in_transit', 'Dhaka Distribution Center', 'Status changed from in_transit to in_transit', '2026-02-09 06:18:01', '2026-02-09 06:18:01'),
(5, 1, 'out_for_delivery', 'Gulshan Delivery Hub', 'Out for delivery - arriving today', '2026-02-10 09:00:00', '2026-02-09 06:23:08');

-- --------------------------------------------------------

--
-- Table structure for table `shipping_addresses`
--

CREATE TABLE `shipping_addresses` (
  `id` bigint UNSIGNED NOT NULL,
  `customer_id` bigint UNSIGNED DEFAULT NULL,
  `full_name` varchar(255) NOT NULL,
  `phone` varchar(20) NOT NULL,
  `email` varchar(255) DEFAULT NULL,
  `address_line1` varchar(255) NOT NULL,
  `address_line2` varchar(255) DEFAULT NULL,
  `city` varchar(100) NOT NULL,
  `state` varchar(100) NOT NULL,
  `postal_code` varchar(20) NOT NULL,
  `country` varchar(100) NOT NULL DEFAULT 'Bangladesh',
  `is_default` tinyint(1) DEFAULT '0',
  `address_type` enum('home','office','other') DEFAULT 'home',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

--
-- Dumping data for table `shipping_addresses`
--

INSERT INTO `shipping_addresses` (`id`, `customer_id`, `full_name`, `phone`, `email`, `address_line1`, `address_line2`, `city`, `state`, `postal_code`, `country`, `is_default`, `address_type`, `created_at`, `updated_at`, `deleted_at`) VALUES
(4, 2, 'John Smith', '+8801712345678', 'john@example.com', '123 Main Street', 'Apartment 4B', 'Dhaka', 'Dhaka Division', '1205', 'Bangladesh', 0, 'home', '2026-02-09 04:48:40', '2026-02-09 04:56:22', NULL),
(6, 2, 'Jane Smith', '+8801812345678', '', '456 Commerce Road', '', 'Chittagong', 'Chittagong Division', '4000', 'Bangladesh', 1, 'office', '2026-02-09 04:56:22', '2026-02-09 04:56:22', NULL);

-- --------------------------------------------------------

--
-- Table structure for table `staff`
--

CREATE TABLE `staff` (
  `id` bigint UNSIGNED NOT NULL,
  `user_id` bigint UNSIGNED DEFAULT NULL,
  `name` varchar(150) NOT NULL,
  `email` varchar(150) NOT NULL,
  `contact` varchar(20) DEFAULT NULL,
  `joining_date` varchar(20) DEFAULT NULL,
  `role` varchar(100) DEFAULT NULL,
  `status` enum('Active','Inactive') DEFAULT 'Active',
  `published` tinyint(1) DEFAULT '0',
  `avatar` varchar(500) DEFAULT NULL,
  `salary` decimal(12,2) DEFAULT '0.00',
  `bank_account` varchar(100) DEFAULT NULL,
  `payment_method` enum('Bank Transfer','Cash','Check') DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

--
-- Dumping data for table `staff`
--

INSERT INTO `staff` (`id`, `user_id`, `name`, `email`, `contact`, `joining_date`, `role`, `status`, `published`, `avatar`, `salary`, `bank_account`, `payment_method`, `created_at`, `updated_at`, `deleted_at`) VALUES
(2, 6, 'Alice Manager', 'alice@store.com', '+1-555-9999', '2025-03-15', 'Manager', 'Active', 1, '', 4500.00, '', 'Bank Transfer', '2026-02-05 07:24:10', '2026-02-05 07:25:24', '2026-02-05 07:25:25'),
(3, 14, 'staff', 'staff@store.com', '+1-555-9999', '2025-03-15', 'Manager', 'Active', 1, 'https://example.com/avatar.png', 4500.00, 'ACC-12345', 'Bank Transfer', '2026-02-05 09:35:15', '2026-02-05 09:35:55', '2026-02-05 09:35:55'),
(4, 17, 'staff update', 'staff1@store.com', '+1-555-9999', '2025-03-15', 'Manager', 'Active', 1, 'https://example.com/avatar.png', 4500.00, 'ACC-12345', 'Bank Transfer', '2026-02-05 09:36:13', '2026-02-05 09:37:19', NULL),
(5, 18, 'staff2', 'staff2@store.com', '+1-555-9999', '2025-03-15', 'Manager', 'Active', 1, 'https://example.com/avatar.png', 4500.00, 'ACC-12345', 'Bank Transfer', '2026-02-05 09:36:20', '2026-02-05 09:36:20', NULL);

-- --------------------------------------------------------

--
-- Table structure for table `staff_roles`
--

CREATE TABLE `staff_roles` (
  `id` bigint UNSIGNED NOT NULL,
  `name` varchar(100) NOT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

--
-- Dumping data for table `staff_roles`
--

INSERT INTO `staff_roles` (`id`, `name`, `created_at`, `updated_at`, `deleted_at`) VALUES
(1, 'Manager', '2026-02-05 11:23:40', '2026-02-08 06:33:13', NULL),
(5, 'string', '2026-02-05 12:59:38', '2026-02-05 12:59:38', NULL);

-- --------------------------------------------------------

--
-- Table structure for table `stock_transfers`
--

CREATE TABLE `stock_transfers` (
  `id` bigint UNSIGNED NOT NULL,
  `product_id` bigint UNSIGNED NOT NULL,
  `variant_id` bigint UNSIGNED DEFAULT NULL,
  `from_location_id` bigint UNSIGNED NOT NULL,
  `to_location_id` bigint UNSIGNED NOT NULL,
  `quantity` int NOT NULL,
  `status` enum('Pending','Completed','Cancelled') NOT NULL DEFAULT 'Pending',
  `notes` text,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

--
-- Dumping data for table `stock_transfers`
--

INSERT INTO `stock_transfers` (`id`, `product_id`, `variant_id`, `from_location_id`, `to_location_id`, `quantity`, `status`, `notes`, `created_at`, `updated_at`) VALUES
(1, 17, 49, 1, 2, 5, 'Cancelled', 'Moving to wirehouse 2', '2026-02-05 10:14:18', '2026-02-05 10:14:36'),
(2, 17, NULL, 1, 2, 20, 'Completed', 'Partial move to wirehouse', '2026-02-05 10:17:07', '2026-02-05 10:17:07'),
(3, 17, 49, 1, 2, 5, 'Cancelled', 'Moving to wirehouse 2', '2026-02-05 11:09:18', '2026-02-05 11:10:22');

-- --------------------------------------------------------

--
-- Table structure for table `users`
--

CREATE TABLE `users` (
  `id` bigint UNSIGNED NOT NULL,
  `username` varchar(100) NOT NULL,
  `role_id` bigint UNSIGNED NOT NULL,
  `email` varchar(150) NOT NULL,
  `password` varchar(255) NOT NULL,
  `address` varchar(255) DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

--
-- Dumping data for table `users`
--

INSERT INTO `users` (`id`, `username`, `role_id`, `email`, `password`, `address`, `created_at`, `updated_at`, `deleted_at`) VALUES
(1, 'Jhon', 1, 'john2@example.com', '$2a$10$maBNXkR795w3bM8LMLuABu03gBl6c0Mjm127.3S6ffQ/MitvwAY4e', 'Singapore', '2026-01-26 13:03:14', '2026-01-27 06:15:42', NULL),
(2, 'Admin', 1, 'admin@gmail.com', '$2a$10$k/jAiUrUBMehogU1qIzFR.bNFyVGVn0Wc5MGbxBQvYSSV3.rEWBqu', 'New York', '2026-01-27 05:58:46', '2026-01-27 05:58:46', NULL),
(4, 'John Smith', 3, 'john@example.com', '$2a$10$ZHH9kLj9LziEvqTyAyF0zOn9w8P68YQq.BvF6yS3klKaOZyRvrDQa', '', '2026-02-05 07:24:10', '2026-02-05 07:25:24', '2026-02-05 07:25:25'),
(5, 'Acme Corp', 4, 'acme@example.com', '$2a$10$E63Rk4mfqjJQR6TFins6vOPlO6/X3ijwTVD99izAW3ZgzAjUYz9l2', '', '2026-02-05 07:24:10', '2026-02-05 07:25:24', '2026-02-05 07:25:25'),
(6, 'Alice Manager', 5, 'alice@store.com', '$2a$10$MIBWDPeoq9HUxAUXFNGeOuc4ghDLrIaiIoFcN5h1AOAxM9ryo.ioS', '', '2026-02-05 07:24:10', '2026-02-05 07:25:24', '2026-02-05 07:25:25'),
(8, 'customer', 3, 'customer@example.com', '$2a$10$3UErCuGxFA5VmtmnXK002u9s0BgJhgjr.anWWqkrEIWbEFbaZhXOK', '', '2026-02-05 09:01:49', '2026-02-05 09:01:49', NULL),
(9, 'customer2', 3, 'customer2@example.com', '$2a$10$gwl.L7HKZvxeSLerdZ0ZGuzyAhWj4JuK8XN4/Vj6DdkK4/mjlv11G', '', '2026-02-05 09:02:30', '2026-02-05 09:02:30', NULL),
(10, 'customer3', 3, 'customer3@example.com', '$2a$10$Pv5Rc9oyMFCYE1WXhPICzet2QZJANpEXx8ZB7Xf6pVLYOHpgGk61O', '', '2026-02-05 09:02:40', '2026-02-05 09:02:40', NULL),
(11, 'Test Customer', 3, 'test5@example.com', '$2a$10$7rlvvle1a6q..Xh4hPKnkO3iGttdeKsCCDsw0rql7TG7xRagq/2iq', '', '2026-02-05 09:04:26', '2026-02-05 09:05:48', '2026-02-05 09:05:48'),
(12, 'Acme Corp', 4, 'vendor@example.com', '$2a$10$yllnpeytuQaQDtxkbNlIVuYLjvYO.0c0ibn47ofzVI1TtDMDxwzu6', '', '2026-02-05 09:07:55', '2026-02-05 09:07:55', NULL),
(13, 'Acme Corp 2', 4, 'vendor2@example.com', '$2a$10$2AAUnR5aUr0iPZQjV9XSVOmmBEf0OXy2CQ/m72KjmGVGUI2.fo95i', '', '2026-02-05 09:08:15', '2026-02-05 09:09:37', '2026-02-05 09:09:38'),
(14, 'staff', 5, 'staff@store.com', '$2a$10$xJlqJfYGqyo7dWdvUi5NRuafZ.l0MzNxMheyDPzS1I6gEUFFPGWUG', '', '2026-02-05 09:35:15', '2026-02-05 09:35:55', '2026-02-05 09:35:55'),
(17, 'staff update', 5, 'staff1@store.com', '$2a$10$W1WJ0bk70mKzCT5SXgjDoOs7pWhMpV8hkxR9VXcE6fJE9NUGZ0IiO', '', '2026-02-05 09:36:13', '2026-02-05 09:37:19', NULL),
(18, 'staff2', 5, 'staff2@store.com', '$2a$10$gsDUf99c9qwIBJhXcpEAge3HaldGeIo6r0A9Q407YIBcAeIDTu/Aq', '', '2026-02-05 09:36:20', '2026-02-05 09:36:20', NULL),
(23, 'Test Customer', 3, 'testcust@example.com', '$2a$10$5LIGO574suMxQkGX5k0.UuX2BY51bnnaaURqKmVXa3q4486Oay4q2', '', '2026-02-08 08:31:31', '2026-02-08 08:31:31', NULL);

-- --------------------------------------------------------

--
-- Table structure for table `variant_inventory`
--

CREATE TABLE `variant_inventory` (
  `id` bigint UNSIGNED NOT NULL,
  `variant_id` bigint UNSIGNED NOT NULL,
  `location_id` bigint UNSIGNED NOT NULL,
  `quantity` int NOT NULL DEFAULT '0',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

--
-- Dumping data for table `variant_inventory`
--

INSERT INTO `variant_inventory` (`id`, `variant_id`, `location_id`, `quantity`, `created_at`, `updated_at`) VALUES
(1, 45, 1, 10, '2026-02-05 10:13:57', '2026-02-05 10:13:57'),
(2, 46, 1, 8, '2026-02-05 10:13:57', '2026-02-05 10:13:57'),
(3, 49, 1, 15, '2026-02-05 10:14:13', '2026-02-05 11:10:22'),
(4, 50, 1, 8, '2026-02-05 10:14:13', '2026-02-05 10:14:13'),
(5, 49, 2, 0, '2026-02-05 10:14:18', '2026-02-05 11:10:22');

-- --------------------------------------------------------

--
-- Table structure for table `vendors`
--

CREATE TABLE `vendors` (
  `id` bigint UNSIGNED NOT NULL,
  `user_id` bigint UNSIGNED DEFAULT NULL,
  `name` varchar(150) NOT NULL,
  `email` varchar(150) NOT NULL,
  `phone` varchar(20) DEFAULT NULL,
  `address` varchar(255) DEFAULT NULL,
  `logo` varchar(500) DEFAULT NULL,
  `status` enum('Active','Inactive','Blocked') DEFAULT 'Active',
  `description` text,
  `total_paid` decimal(12,2) DEFAULT '0.00',
  `amount_payable` decimal(12,2) DEFAULT '0.00',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

--
-- Dumping data for table `vendors`
--

INSERT INTO `vendors` (`id`, `user_id`, `name`, `email`, `phone`, `address`, `logo`, `status`, `description`, `total_paid`, `amount_payable`, `created_at`, `updated_at`, `deleted_at`) VALUES
(2, 5, 'Acme Corp', 'acme@example.com', '555-5678', '', '', 'Active', 'A test vendor', 0.00, 0.00, '2026-02-05 07:24:10', '2026-02-05 07:25:24', '2026-02-05 07:25:25'),
(3, 12, 'Acme Corp', 'vendor@example.com', '555-5678', '456 Business Ave', 'https://example.com/logo.png', 'Active', 'A reliable supplier', 0.00, 0.00, '2026-02-05 09:07:55', '2026-02-05 09:07:55', NULL),
(4, 13, 'Acme Corp 2', 'vendor2@example.com', '555-5678', '456 Business Ave', 'https://example.com/logo.png', 'Active', 'A reliable supplier', 0.00, 0.00, '2026-02-05 09:08:15', '2026-02-05 09:09:37', '2026-02-05 09:09:38');

-- --------------------------------------------------------

--
-- Table structure for table `vendor_returns`
--

CREATE TABLE `vendor_returns` (
  `id` bigint UNSIGNED NOT NULL,
  `return_number` varchar(50) NOT NULL,
  `vendor_id` bigint UNSIGNED NOT NULL,
  `vendor_name` varchar(255) NOT NULL,
  `total_amount` decimal(10,2) NOT NULL DEFAULT '0.00',
  `status` enum('pending','shipped','received_by_vendor','completed') NOT NULL DEFAULT 'pending',
  `return_date` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `completed_date` timestamp NULL DEFAULT NULL,
  `credit_type` enum('refund','credit_note','replacement') NOT NULL,
  `notes` text,
  `created_by` varchar(255) NOT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

--
-- Dumping data for table `vendor_returns`
--

INSERT INTO `vendor_returns` (`id`, `return_number`, `vendor_id`, `vendor_name`, `total_amount`, `status`, `return_date`, `completed_date`, `credit_type`, `notes`, `created_by`, `created_at`, `updated_at`, `deleted_at`) VALUES
(2, 'VRT-49364', 2, 'Acme Corp', 1127.20, 'completed', '2026-02-08 11:16:04', '2026-02-08 11:16:11', 'credit_note', '', 'Admin', '2026-02-08 11:16:04', '2026-02-08 11:16:11', NULL),
(3, 'VRT-49577', 2, 'Acme Corp', 1000.00, 'completed', '2026-02-08 11:19:37', '2026-02-08 11:41:13', 'refund', '', 'Admin', '2026-02-08 11:19:37', '2026-02-08 11:41:13', NULL),
(6, 'VRT-1770550827744490', 2, 'Acme Corp', 1127.20, 'pending', '2026-02-08 11:40:28', NULL, 'credit_note', '', 'Admin', '2026-02-08 11:40:28', '2026-02-08 11:40:28', NULL);

-- --------------------------------------------------------

--
-- Table structure for table `vendor_return_items`
--

CREATE TABLE `vendor_return_items` (
  `id` bigint UNSIGNED NOT NULL,
  `return_id` bigint UNSIGNED NOT NULL,
  `product_id` bigint UNSIGNED DEFAULT NULL,
  `product_name` varchar(255) NOT NULL,
  `variant_id` bigint UNSIGNED DEFAULT NULL,
  `variant_name` varchar(255) DEFAULT NULL,
  `quantity` int NOT NULL DEFAULT '1',
  `unit_price` decimal(10,2) NOT NULL DEFAULT '0.00',
  `total_price` decimal(10,2) NOT NULL DEFAULT '0.00',
  `reason` varchar(255) NOT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

--
-- Dumping data for table `vendor_return_items`
--

INSERT INTO `vendor_return_items` (`id`, `return_id`, `product_id`, `product_name`, `variant_id`, `variant_name`, `quantity`, `unit_price`, `total_price`, `reason`, `created_at`, `updated_at`) VALUES
(1, 2, NULL, 'Test Product', NULL, '', 10, 112.72, 1127.20, 'Defective/Damaged', '2026-02-08 11:16:04', '2026-02-08 11:16:04'),
(2, 3, 28, 'Test Product for Returns', NULL, '', 10, 100.00, 1000.00, 'Overstocked items', '2026-02-08 11:19:37', '2026-02-08 11:19:37'),
(4, 6, NULL, 'Green Leaf Lettuce', NULL, '', 10, 112.72, 1127.20, 'Damaged during shipping', '2026-02-08 11:40:28', '2026-02-08 11:40:28');

--
-- Indexes for dumped tables
--

--
-- Indexes for table `attributes`
--
ALTER TABLE `attributes`
  ADD PRIMARY KEY (`id`);

--
-- Indexes for table `authors`
--
ALTER TABLE `authors`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `email` (`email`);

--
-- Indexes for table `books`
--
ALTER TABLE `books`
  ADD PRIMARY KEY (`id`),
  ADD KEY `fk_books_author` (`author_id`);

--
-- Indexes for table `categories`
--
ALTER TABLE `categories`
  ADD PRIMARY KEY (`id`),
  ADD KEY `fk_parent` (`parent_id`);

--
-- Indexes for table `coupons`
--
ALTER TABLE `coupons`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `code` (`code`);

--
-- Indexes for table `customers`
--
ALTER TABLE `customers`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `email` (`email`),
  ADD KEY `idx_customers_deleted_at` (`deleted_at`),
  ADD KEY `fk_customers_user` (`user_id`);

--
-- Indexes for table `customer_returns`
--
ALTER TABLE `customer_returns`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `return_number` (`return_number`),
  ADD KEY `order_id` (`order_id`),
  ADD KEY `idx_customer_returns_status` (`status`),
  ADD KEY `idx_customer_returns_customer_id` (`customer_id`),
  ADD KEY `idx_customer_returns_request_date` (`request_date`),
  ADD KEY `idx_customer_returns_deleted_at` (`deleted_at`);

--
-- Indexes for table `customer_return_items`
--
ALTER TABLE `customer_return_items`
  ADD PRIMARY KEY (`id`),
  ADD KEY `variant_id` (`variant_id`),
  ADD KEY `idx_customer_return_items_return_id` (`return_id`),
  ADD KEY `idx_customer_return_items_product_id` (`product_id`);

--
-- Indexes for table `locations`
--
ALTER TABLE `locations`
  ADD PRIMARY KEY (`id`);

--
-- Indexes for table `order_items`
--
ALTER TABLE `order_items`
  ADD PRIMARY KEY (`id`),
  ADD KEY `idx_order_items_sell_id` (`sell_id`),
  ADD KEY `idx_order_items_product_id` (`product_id`),
  ADD KEY `idx_order_items_variant_id` (`variant_id`);

--
-- Indexes for table `order_shipments`
--
ALTER TABLE `order_shipments`
  ADD PRIMARY KEY (`id`),
  ADD KEY `idx_order_shipments_sell_id` (`sell_id`),
  ADD KEY `idx_order_shipments_tracking_number` (`tracking_number`),
  ADD KEY `idx_order_shipments_status` (`status`);

--
-- Indexes for table `permissions`
--
ALTER TABLE `permissions`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `name` (`name`);

--
-- Indexes for table `products`
--
ALTER TABLE `products`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `sku` (`sku`),
  ADD UNIQUE KEY `barcode` (`barcode`),
  ADD KEY `fk_products_category` (`category_id`),
  ADD KEY `fk_products_vendor` (`vendor_id`),
  ADD KEY `fk_products_location` (`location_id`);

--
-- Indexes for table `product_attributes`
--
ALTER TABLE `product_attributes`
  ADD PRIMARY KEY (`product_id`,`attribute_id`),
  ADD KEY `fk_product_attributes_attribute` (`attribute_id`);

--
-- Indexes for table `product_images`
--
ALTER TABLE `product_images`
  ADD PRIMARY KEY (`id`),
  ADD KEY `idx_product_images_product` (`product_id`),
  ADD KEY `idx_product_images_position` (`product_id`,`position`);

--
-- Indexes for table `product_variants`
--
ALTER TABLE `product_variants`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `idx_product_sku` (`product_id`,`sku`),
  ADD UNIQUE KEY `idx_product_barcode` (`product_id`,`barcode`);

--
-- Indexes for table `roles`
--
ALTER TABLE `roles`
  ADD PRIMARY KEY (`id`),
  ADD KEY `idx_roles_deleted_at` (`deleted_at`);

--
-- Indexes for table `role_permissions`
--
ALTER TABLE `role_permissions`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `uq_role_permission` (`role_id`,`permission_id`),
  ADD KEY `permission_id` (`permission_id`);

--
-- Indexes for table `salary_payments`
--
ALTER TABLE `salary_payments`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `idx_staff_month` (`staff_id`,`month`);

--
-- Indexes for table `schema_migrations`
--
ALTER TABLE `schema_migrations`
  ADD PRIMARY KEY (`version`);

--
-- Indexes for table `sells`
--
ALTER TABLE `sells`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `invoice_no` (`invoice_no`),
  ADD KEY `idx_sells_invoice_no` (`invoice_no`),
  ADD KEY `idx_sells_customer_id` (`customer_id`),
  ADD KEY `idx_sells_status` (`status`),
  ADD KEY `idx_sells_method` (`method`),
  ADD KEY `idx_sells_order_time` (`order_time`),
  ADD KEY `idx_sells_deleted_at` (`deleted_at`),
  ADD KEY `idx_sells_shipping_address_id` (`shipping_address_id`),
  ADD KEY `idx_sells_payment_status` (`payment_status`),
  ADD KEY `idx_sells_fulfillment_status` (`fulfillment_status`),
  ADD KEY `idx_sells_tracking_number` (`tracking_number`),
  ADD KEY `idx_sells_shipping_city` (`shipping_city`),
  ADD KEY `idx_sells_shipping_postal_code` (`shipping_postal_code`),
  ADD KEY `idx_sells_shipping_country` (`shipping_country`);

--
-- Indexes for table `shipment_tracking_history`
--
ALTER TABLE `shipment_tracking_history`
  ADD PRIMARY KEY (`id`),
  ADD KEY `idx_shipment_tracking_shipment_id` (`shipment_id`),
  ADD KEY `idx_shipment_tracking_event_time` (`event_time`);

--
-- Indexes for table `shipping_addresses`
--
ALTER TABLE `shipping_addresses`
  ADD PRIMARY KEY (`id`),
  ADD KEY `idx_shipping_addresses_customer_id` (`customer_id`),
  ADD KEY `idx_shipping_addresses_is_default` (`is_default`),
  ADD KEY `idx_shipping_addresses_deleted_at` (`deleted_at`);

--
-- Indexes for table `staff`
--
ALTER TABLE `staff`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `email` (`email`),
  ADD KEY `idx_staff_deleted_at` (`deleted_at`),
  ADD KEY `fk_staff_user` (`user_id`);

--
-- Indexes for table `staff_roles`
--
ALTER TABLE `staff_roles`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `name` (`name`),
  ADD KEY `idx_staff_roles_deleted_at` (`deleted_at`);

--
-- Indexes for table `stock_transfers`
--
ALTER TABLE `stock_transfers`
  ADD PRIMARY KEY (`id`),
  ADD KEY `fk_transfer_product` (`product_id`),
  ADD KEY `fk_transfer_variant` (`variant_id`),
  ADD KEY `fk_transfer_from_location` (`from_location_id`),
  ADD KEY `fk_transfer_to_location` (`to_location_id`);

--
-- Indexes for table `users`
--
ALTER TABLE `users`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `email` (`email`),
  ADD KEY `fk_roles` (`role_id`);

--
-- Indexes for table `variant_inventory`
--
ALTER TABLE `variant_inventory`
  ADD PRIMARY KEY (`id`),
  ADD KEY `fk_variant_inventory_variant` (`variant_id`),
  ADD KEY `fk_variant_inventory_location` (`location_id`);

--
-- Indexes for table `vendors`
--
ALTER TABLE `vendors`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `email` (`email`),
  ADD KEY `idx_vendors_deleted_at` (`deleted_at`),
  ADD KEY `fk_vendors_user` (`user_id`);

--
-- Indexes for table `vendor_returns`
--
ALTER TABLE `vendor_returns`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `return_number` (`return_number`),
  ADD KEY `idx_vendor_returns_status` (`status`),
  ADD KEY `idx_vendor_returns_vendor_id` (`vendor_id`),
  ADD KEY `idx_vendor_returns_return_date` (`return_date`),
  ADD KEY `idx_vendor_returns_deleted_at` (`deleted_at`);

--
-- Indexes for table `vendor_return_items`
--
ALTER TABLE `vendor_return_items`
  ADD PRIMARY KEY (`id`),
  ADD KEY `variant_id` (`variant_id`),
  ADD KEY `idx_vendor_return_items_return_id` (`return_id`),
  ADD KEY `idx_vendor_return_items_product_id` (`product_id`);

--
-- AUTO_INCREMENT for dumped tables
--

--
-- AUTO_INCREMENT for table `attributes`
--
ALTER TABLE `attributes`
  MODIFY `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=4;

--
-- AUTO_INCREMENT for table `authors`
--
ALTER TABLE `authors`
  MODIFY `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `books`
--
ALTER TABLE `books`
  MODIFY `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `categories`
--
ALTER TABLE `categories`
  MODIFY `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=4;

--
-- AUTO_INCREMENT for table `coupons`
--
ALTER TABLE `coupons`
  MODIFY `id` int NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=4;

--
-- AUTO_INCREMENT for table `customers`
--
ALTER TABLE `customers`
  MODIFY `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=8;

--
-- AUTO_INCREMENT for table `customer_returns`
--
ALTER TABLE `customer_returns`
  MODIFY `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=5;

--
-- AUTO_INCREMENT for table `customer_return_items`
--
ALTER TABLE `customer_return_items`
  MODIFY `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=4;

--
-- AUTO_INCREMENT for table `locations`
--
ALTER TABLE `locations`
  MODIFY `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=3;

--
-- AUTO_INCREMENT for table `order_items`
--
ALTER TABLE `order_items`
  MODIFY `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=16;

--
-- AUTO_INCREMENT for table `order_shipments`
--
ALTER TABLE `order_shipments`
  MODIFY `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=3;

--
-- AUTO_INCREMENT for table `permissions`
--
ALTER TABLE `permissions`
  MODIFY `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=15;

--
-- AUTO_INCREMENT for table `products`
--
ALTER TABLE `products`
  MODIFY `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=29;

--
-- AUTO_INCREMENT for table `product_images`
--
ALTER TABLE `product_images`
  MODIFY `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=24;

--
-- AUTO_INCREMENT for table `product_variants`
--
ALTER TABLE `product_variants`
  MODIFY `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=67;

--
-- AUTO_INCREMENT for table `roles`
--
ALTER TABLE `roles`
  MODIFY `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=6;

--
-- AUTO_INCREMENT for table `role_permissions`
--
ALTER TABLE `role_permissions`
  MODIFY `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=35;

--
-- AUTO_INCREMENT for table `salary_payments`
--
ALTER TABLE `salary_payments`
  MODIFY `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=4;

--
-- AUTO_INCREMENT for table `sells`
--
ALTER TABLE `sells`
  MODIFY `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=23;

--
-- AUTO_INCREMENT for table `shipment_tracking_history`
--
ALTER TABLE `shipment_tracking_history`
  MODIFY `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=6;

--
-- AUTO_INCREMENT for table `shipping_addresses`
--
ALTER TABLE `shipping_addresses`
  MODIFY `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=7;

--
-- AUTO_INCREMENT for table `staff`
--
ALTER TABLE `staff`
  MODIFY `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=6;

--
-- AUTO_INCREMENT for table `staff_roles`
--
ALTER TABLE `staff_roles`
  MODIFY `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=6;

--
-- AUTO_INCREMENT for table `stock_transfers`
--
ALTER TABLE `stock_transfers`
  MODIFY `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=4;

--
-- AUTO_INCREMENT for table `users`
--
ALTER TABLE `users`
  MODIFY `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=24;

--
-- AUTO_INCREMENT for table `variant_inventory`
--
ALTER TABLE `variant_inventory`
  MODIFY `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=6;

--
-- AUTO_INCREMENT for table `vendors`
--
ALTER TABLE `vendors`
  MODIFY `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=5;

--
-- AUTO_INCREMENT for table `vendor_returns`
--
ALTER TABLE `vendor_returns`
  MODIFY `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=7;

--
-- AUTO_INCREMENT for table `vendor_return_items`
--
ALTER TABLE `vendor_return_items`
  MODIFY `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=5;

--
-- Constraints for dumped tables
--

--
-- Constraints for table `books`
--
ALTER TABLE `books`
  ADD CONSTRAINT `fk_books_author` FOREIGN KEY (`author_id`) REFERENCES `authors` (`id`) ON DELETE RESTRICT ON UPDATE CASCADE;

--
-- Constraints for table `categories`
--
ALTER TABLE `categories`
  ADD CONSTRAINT `fk_parent` FOREIGN KEY (`parent_id`) REFERENCES `categories` (`id`);

--
-- Constraints for table `customers`
--
ALTER TABLE `customers`
  ADD CONSTRAINT `fk_customers_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE SET NULL ON UPDATE CASCADE;

--
-- Constraints for table `customer_returns`
--
ALTER TABLE `customer_returns`
  ADD CONSTRAINT `customer_returns_ibfk_1` FOREIGN KEY (`customer_id`) REFERENCES `customers` (`id`) ON DELETE RESTRICT,
  ADD CONSTRAINT `customer_returns_ibfk_2` FOREIGN KEY (`order_id`) REFERENCES `sells` (`id`) ON DELETE SET NULL;

--
-- Constraints for table `customer_return_items`
--
ALTER TABLE `customer_return_items`
  ADD CONSTRAINT `customer_return_items_ibfk_1` FOREIGN KEY (`return_id`) REFERENCES `customer_returns` (`id`) ON DELETE CASCADE,
  ADD CONSTRAINT `customer_return_items_ibfk_2` FOREIGN KEY (`product_id`) REFERENCES `products` (`id`) ON DELETE SET NULL,
  ADD CONSTRAINT `customer_return_items_ibfk_3` FOREIGN KEY (`variant_id`) REFERENCES `product_variants` (`id`) ON DELETE SET NULL;

--
-- Constraints for table `order_items`
--
ALTER TABLE `order_items`
  ADD CONSTRAINT `order_items_ibfk_1` FOREIGN KEY (`sell_id`) REFERENCES `sells` (`id`) ON DELETE CASCADE,
  ADD CONSTRAINT `order_items_ibfk_2` FOREIGN KEY (`product_id`) REFERENCES `products` (`id`) ON DELETE SET NULL,
  ADD CONSTRAINT `order_items_ibfk_3` FOREIGN KEY (`variant_id`) REFERENCES `product_variants` (`id`) ON DELETE SET NULL;

--
-- Constraints for table `order_shipments`
--
ALTER TABLE `order_shipments`
  ADD CONSTRAINT `order_shipments_ibfk_1` FOREIGN KEY (`sell_id`) REFERENCES `sells` (`id`) ON DELETE CASCADE;

--
-- Constraints for table `products`
--
ALTER TABLE `products`
  ADD CONSTRAINT `fk_products_category` FOREIGN KEY (`category_id`) REFERENCES `categories` (`id`) ON DELETE SET NULL,
  ADD CONSTRAINT `fk_products_location` FOREIGN KEY (`location_id`) REFERENCES `locations` (`id`) ON DELETE SET NULL,
  ADD CONSTRAINT `fk_products_vendor` FOREIGN KEY (`vendor_id`) REFERENCES `users` (`id`) ON DELETE SET NULL;

--
-- Constraints for table `product_attributes`
--
ALTER TABLE `product_attributes`
  ADD CONSTRAINT `fk_product_attributes_attribute` FOREIGN KEY (`attribute_id`) REFERENCES `attributes` (`id`) ON DELETE CASCADE,
  ADD CONSTRAINT `fk_product_attributes_product` FOREIGN KEY (`product_id`) REFERENCES `products` (`id`) ON DELETE CASCADE;

--
-- Constraints for table `product_images`
--
ALTER TABLE `product_images`
  ADD CONSTRAINT `fk_product_images_product` FOREIGN KEY (`product_id`) REFERENCES `products` (`id`) ON DELETE CASCADE;

--
-- Constraints for table `product_variants`
--
ALTER TABLE `product_variants`
  ADD CONSTRAINT `fk_product_variants_product` FOREIGN KEY (`product_id`) REFERENCES `products` (`id`) ON DELETE CASCADE;

--
-- Constraints for table `role_permissions`
--
ALTER TABLE `role_permissions`
  ADD CONSTRAINT `role_permissions_ibfk_1` FOREIGN KEY (`role_id`) REFERENCES `staff_roles` (`id`) ON DELETE CASCADE,
  ADD CONSTRAINT `role_permissions_ibfk_2` FOREIGN KEY (`permission_id`) REFERENCES `permissions` (`id`) ON DELETE CASCADE;

--
-- Constraints for table `salary_payments`
--
ALTER TABLE `salary_payments`
  ADD CONSTRAINT `fk_salary_payment_staff` FOREIGN KEY (`staff_id`) REFERENCES `staff` (`id`) ON DELETE CASCADE;

--
-- Constraints for table `sells`
--
ALTER TABLE `sells`
  ADD CONSTRAINT `sells_ibfk_1` FOREIGN KEY (`customer_id`) REFERENCES `customers` (`id`) ON DELETE SET NULL,
  ADD CONSTRAINT `sells_ibfk_2` FOREIGN KEY (`shipping_address_id`) REFERENCES `shipping_addresses` (`id`) ON DELETE SET NULL;

--
-- Constraints for table `shipment_tracking_history`
--
ALTER TABLE `shipment_tracking_history`
  ADD CONSTRAINT `shipment_tracking_history_ibfk_1` FOREIGN KEY (`shipment_id`) REFERENCES `order_shipments` (`id`) ON DELETE CASCADE;

--
-- Constraints for table `shipping_addresses`
--
ALTER TABLE `shipping_addresses`
  ADD CONSTRAINT `shipping_addresses_ibfk_1` FOREIGN KEY (`customer_id`) REFERENCES `customers` (`id`) ON DELETE CASCADE;

--
-- Constraints for table `staff`
--
ALTER TABLE `staff`
  ADD CONSTRAINT `fk_staff_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE SET NULL ON UPDATE CASCADE;

--
-- Constraints for table `stock_transfers`
--
ALTER TABLE `stock_transfers`
  ADD CONSTRAINT `fk_transfer_from_location` FOREIGN KEY (`from_location_id`) REFERENCES `locations` (`id`) ON DELETE CASCADE,
  ADD CONSTRAINT `fk_transfer_product` FOREIGN KEY (`product_id`) REFERENCES `products` (`id`) ON DELETE CASCADE,
  ADD CONSTRAINT `fk_transfer_to_location` FOREIGN KEY (`to_location_id`) REFERENCES `locations` (`id`) ON DELETE CASCADE,
  ADD CONSTRAINT `fk_transfer_variant` FOREIGN KEY (`variant_id`) REFERENCES `product_variants` (`id`) ON DELETE SET NULL;

--
-- Constraints for table `users`
--
ALTER TABLE `users`
  ADD CONSTRAINT `fk_roles` FOREIGN KEY (`role_id`) REFERENCES `roles` (`id`) ON DELETE RESTRICT ON UPDATE CASCADE;

--
-- Constraints for table `variant_inventory`
--
ALTER TABLE `variant_inventory`
  ADD CONSTRAINT `fk_variant_inventory_location` FOREIGN KEY (`location_id`) REFERENCES `locations` (`id`) ON DELETE CASCADE,
  ADD CONSTRAINT `fk_variant_inventory_variant` FOREIGN KEY (`variant_id`) REFERENCES `product_variants` (`id`) ON DELETE CASCADE;

--
-- Constraints for table `vendors`
--
ALTER TABLE `vendors`
  ADD CONSTRAINT `fk_vendors_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE SET NULL ON UPDATE CASCADE;

--
-- Constraints for table `vendor_returns`
--
ALTER TABLE `vendor_returns`
  ADD CONSTRAINT `vendor_returns_ibfk_1` FOREIGN KEY (`vendor_id`) REFERENCES `vendors` (`id`) ON DELETE RESTRICT;

--
-- Constraints for table `vendor_return_items`
--
ALTER TABLE `vendor_return_items`
  ADD CONSTRAINT `vendor_return_items_ibfk_1` FOREIGN KEY (`return_id`) REFERENCES `vendor_returns` (`id`) ON DELETE CASCADE,
  ADD CONSTRAINT `vendor_return_items_ibfk_2` FOREIGN KEY (`product_id`) REFERENCES `products` (`id`) ON DELETE SET NULL,
  ADD CONSTRAINT `vendor_return_items_ibfk_3` FOREIGN KEY (`variant_id`) REFERENCES `product_variants` (`id`) ON DELETE SET NULL;
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
